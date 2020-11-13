package sweep

import (
	"hash"
	"sync"
	"time"
)

type Sweep struct {
	shardsCount uint64

	hasher hash.Hash64

	cache map[string]*entry

	entryLifetime time.Duration

	cleanupInterval time.Duration

	closeCh chan struct{}

	mu *sync.Mutex
}

// Get retrieves value associated with the key from the sweep.
func (s *Sweep) Get(key string) (value []byte, err error) {
	if s.isClosed() {
		err = ErrClosed
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	e, ok := s.cache[key]
	if !ok {
		err = ErrEntryNotFound
		return
	}

	value = e.val
	return
}

// Put inserts the value associated with the key into the sweep.
func (s *Sweep) Put(key string, value []byte) error {
	if s.isClosed() {
		return ErrClosed
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.cache[key] = &entry{val: value, createdAt: time.Now().Unix()}

	return nil
}

// Close closes the sweep and removes all entries.
func (s *Sweep) Close() error {
	select {
	case <-s.closeCh:
		return ErrClosed
	default:
		s.closeCh <- struct{}{}
	}

	return nil
}

// Default return's sweep with default entry lifetime
// of 10 minutes and 1000 shards.
func Default() *Sweep {
	cfg := setupVacantDefaultsInConfig(Configuration{})

	return new(cfg)
}

func New(cfg Configuration) *Sweep {
	cfg = setupVacantDefaultsInConfig(cfg)

	return new(cfg)
}

func new(cfg Configuration) *Sweep {
	s := &Sweep{
		shardsCount:     cfg.ShardsCount,
		cache:           make(map[string]*entry),
		entryLifetime:   cfg.EntryLifetime,
		cleanupInterval: cfg.CleanupInterval,
		closeCh:         make(chan struct{}),
		mu:              &sync.Mutex{},
	}

	s.startBackgroundCleanupLoop()

	return s
}

func (s *Sweep) cleanupExpiredEntries() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for k, e := range s.cache {
		if time.Since(time.Unix(e.createdAt, 0)) > s.entryLifetime {
			delete(s.cache, k)
		}
	}
}

func (s *Sweep) startBackgroundCleanupLoop() {
	ticker := time.NewTicker(s.cleanupInterval)

	go func() {
		for {
			select {
			case <-s.closeCh:
				return
			case <-ticker.C:
				s.cleanupExpiredEntries()
			}
		}
	}()
}

func (s *Sweep) isClosed() bool {
	select {
	case <-s.closeCh:
		return true
	default:
		return false
	}
}
