package sweep

import (
	"sync/atomic"
	"time"

	"github.com/cespare/xxhash"
)

type Sweep struct {
	cfg Configuration

	shards []*shard

	closeCh chan struct{}

	entriesCount uint64
}

// Get retrieves value associated with the key from the sweep.
func (s *Sweep) Get(key string) (value []byte, err error) {
	if s.isClosed() {
		err = ErrClosed
		return
	}

	keyHash := s.hashKey(key)
	shardAlloted := s.shards[s.getShardIndex(keyHash)]

	val, err := shardAlloted.get(keyHash)
	if err != nil {
		return nil, err
	}

	return val, nil
}

// Put inserts the value associated with the key into the sweep.
func (s *Sweep) Put(key string, value []byte) error {
	if s.isClosed() {
		return ErrClosed
	}

	if len(value) > s.cfg.MaxEntrySize {
		return ErrEntryTooLarge
	}

	keyHash := s.hashKey(key)
	shardAllotted := s.shards[s.getShardIndex(keyHash)]

	err := shardAllotted.put(keyHash, time.Now().Unix(), value)
	if err != nil {
		return err
	}

	atomic.AddUint64(&s.entriesCount, 1)
	return nil
}

// EntriesCount returns number of current entries stored.
// This count includes those entries too which are expired
// but not cleaned up yet.
func (s *Sweep) EntriesCount() int {
	return int(atomic.LoadUint64(&s.entriesCount))
}

// Close closes the sweep and removes all entries.
func (s *Sweep) Close() error {
	select {
	case <-s.closeCh:
		return ErrClosed
	default:
		close(s.closeCh)
	}

	return nil
}

// Default return's sweep with default Entry lifetime
// of 10 minutes and 1000 shards.
func Default() *Sweep {
	cfg := setupVacantDefaultsInConfig(Configuration{})

	return newSweep(cfg)
}

// New return a sweep instance configured to given configuration.
func New(cfg Configuration) *Sweep {
	cfg = setupVacantDefaultsInConfig(cfg)

	return newSweep(cfg)
}

func newSweep(cfg Configuration) *Sweep {
	s := &Sweep{
		cfg:     cfg,
		closeCh: make(chan struct{}),
	}

	// Initialize the shards
	s.shards = make([]*shard, cfg.ShardsCount)
	for i := 0; i < cfg.ShardsCount; i++ {
		s.shards[i] = newShard(cfg.MaxShardSize)
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

func (s *Sweep) getShardIndex(hashedKey uint64) uint64 {
	return hashedKey & (uint64(s.cfg.ShardsCount - 1))
}
