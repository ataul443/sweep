package sweep

import (
	"sync"
	"time"
)

type Sweep struct {
	cache map[string]*entry

	entryLifetime time.Duration

	closeCh chan struct{}

	mu *sync.Mutex
}

type entry struct {
	Val       []byte
	CreatedAt time.Time
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

	value = e.Val
	return
}

// Put inserts the value associated with the key into the sweep.
func (s *Sweep) Put(key string, value []byte) error {
	if s.isClosed() {
		return ErrClosed
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.cache[key] = &entry{Val: value, CreatedAt: time.Now()}

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

// Default return's sweep with default entry lifetime of 10 minutes.
func Default() *Sweep {
	defaultEntryLifeTime := 10 * time.Minute
	return new(defaultEntryLifeTime)
}

func New(entryLifetime time.Duration) *Sweep {
	return new(entryLifetime)
}

func new(entryLifetime time.Duration) *Sweep {
	s := &Sweep{
		cache:         make(map[string]*entry),
		entryLifetime: entryLifetime,
		closeCh:       make(chan struct{}),
		mu:            &sync.Mutex{},
	}

	s.startBackgroundCleanupLoop()

	return s
}

func (s *Sweep) cleanupExpiredEntries() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for k, e := range s.cache {
		if time.Since(e.CreatedAt) > s.entryLifetime {
			delete(s.cache, k)
		}
	}
}

func (s *Sweep) startBackgroundCleanupLoop() {
	ticker := time.NewTicker(time.Second)

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
