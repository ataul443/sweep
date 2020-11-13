package sweep

import (
	"hash"
	"time"
)

const (
	defaultShardsCount = 1000

	defaultEntryLifeTime = 10 // minutes

	defaultCleanupInterval = 1 // minute
)

type Configuration struct {
	// ShardsCount represents a fixed number shards sweep will have.
	ShardsCount int

	// Hasher represents a hash implementation for hashing keys.
	Hasher hash.Hash64

	// EntryLifetime represents lifetime of an entry in the sweep.
	EntryLifetime time.Duration

	// CleanupInterval represents the waiting period between cleanup
	// cycles in sweep.
	CleanupInterval time.Duration
}

func setupVacantDefaultsInConfig(cfg Configuration) Configuration {
	if cfg.ShardsCount <= 0 {
		cfg.ShardsCount = defaultShardsCount
	}

	if cfg.EntryLifetime == 0 {
		cfg.EntryLifetime = defaultEntryLifeTime
	}

	if cfg.CleanupInterval == 0 {
		cfg.CleanupInterval = defaultCleanupInterval
	}

	return cfg
}
