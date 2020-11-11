package sweep

import (
	"fmt"
	"time"
)

type Sweep struct {
	cache map[string][]byte

	entryLifetime time.Duration
}

// Get retrieves value associated with the key from the sweep.
func (s *Sweep) Get(key string) (value []byte, err error) {
	b, ok := s.cache[key]
	if !ok {
		err = fmt.Errorf("key %s, not found in sweep", key)
		return
	}

	value = b
	return
}

// Put inserts the value associated with the key into the sweep.
func (s *Sweep) Put(key string, value []byte) {
	s.cache[key] = value
	return
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
	return &Sweep{
		cache:         make(map[string][]byte),
		entryLifetime: entryLifetime,
	}
}
