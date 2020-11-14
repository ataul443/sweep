package sweep

import (
	"time"
)

const (
	defaultShardsCount = 1024

	defaultShardSize = 4 * 1024 // 4KB

	defaultEntryLifeTime = 10 // minutes

	defaultCleanupInterval = 1 // minute

	defaultMaxEntrySize = 1024 // bytes
)

type Configuration struct {
	// ShardsCount represents a fixed number shards sweep will have.
	// This should be power of two. If it is not, then it will be set
	// to next power of two greater than current value.
	ShardsCount int

	// MaxShardSize represents the upper bound limit of a shard size in bytes.
	// This can be either 0 or should be power of two. If it is not,
	// then it will be set to next power of two greater than current
	// value. A zero value means no restriction on shard size.
	MaxShardSize int

	// EntryLifetime represents lifetime of an Entry in the sweep.
	EntryLifetime time.Duration

	// MaxEntrySize represents maximum size of Entry in bytes
	// in sweep can be stored
	MaxEntrySize int

	// CleanupInterval represents the waiting period between cleanup
	// cycles in sweep.
	CleanupInterval time.Duration
}

func setupVacantDefaultsInConfig(cfg Configuration) Configuration {
	if cfg.ShardsCount <= 0 {
		cfg.ShardsCount = defaultShardsCount
	}

	if !isPowerOfTwo(cfg.ShardsCount) {
		cfg.ShardsCount = getNextPowerOfTwo(cfg.ShardsCount)
	}

	if cfg.MaxShardSize < 0 {
		cfg.MaxShardSize = defaultShardSize
	}

	if cfg.MaxShardSize != 0 {
		if !isPowerOfTwo(cfg.MaxShardSize) {
			cfg.MaxShardSize = getNextPowerOfTwo(cfg.MaxShardSize)
		}
	}

	if cfg.EntryLifetime == 0 {
		cfg.EntryLifetime = defaultEntryLifeTime
	}

	if cfg.CleanupInterval == 0 {
		cfg.CleanupInterval = defaultCleanupInterval
	}

	if cfg.MaxEntrySize == 0 {
		cfg.MaxEntrySize = defaultMaxEntrySize
	}

	return cfg
}

func getNextPowerOfTwo(n int) int {
	k := 1

	for k < n {
		k = k << 1
	}

	return k
}

func isPowerOfTwo(n int) bool {
	return n&(n-1) == 0
}
