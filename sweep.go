package sweep

import (
	"time"

	"github.com/cespare/xxhash"
)

type Sweep struct {
	shardsCount int

	shards []*shard

	entryLifetime time.Duration

	cleanupInterval time.Duration

	closeCh chan struct{}
}

// Get retrieves value associated with the key from the sweep.
func (s *Sweep) Get(key string) (value []byte, err error) {
	if s.isClosed() {
		err = ErrClosed
		return
	}

	keyHash := s.hashKey(key)
	shardAlloted := s.shards[keyHash&uint64(s.shardsCount)]

	value, err = shardAlloted.get(keyHash)
	return
}

// Put inserts the value associated with the key into the sweep.
func (s *Sweep) Put(key string, value []byte) error {
	if s.isClosed() {
		return ErrClosed
	}

	keyHash := s.hashKey(key)
	shardAlloted := s.shards[keyHash&uint64(s.shardsCount)]

	return shardAlloted.put(keyHash, value)
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
		entryLifetime:   cfg.EntryLifetime,
		cleanupInterval: cfg.CleanupInterval,
		closeCh:         make(chan struct{}),
	}

	// Initialize the shards
	s.shards = make([]*shard, s.shardsCount)
	for i := 0; i < s.shardsCount; i++ {
		s.shards[i] = newShard()
	}

	s.startBackgroundCleanupLoop()

	return s
}

func (s *Sweep) cleanupExpiredEntries() {
	for _, sh := range s.shards {
		sh.cleanupExpiredEntries()
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

func (s *Sweep) hashKey(key string) uint64 {
	return xxhash.Sum64([]byte(key))
}

func (s *Sweep) isClosed() bool {
	select {
	case <-s.closeCh:
		return true
	default:
		return false
	}
}
