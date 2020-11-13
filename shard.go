package sweep

import "sync"

const defaultEntryBufferSize = 4 * 1024 // 4KB

type shard struct {
	hashIndexBucket map[uint64]uint64

	entryBuffer []byte

	mu *sync.Mutex
}

func newShard() *shard {
	return &shard{
		hashIndexBucket: make(map[uint64]uint64),
		entryBuffer:     make([]byte, defaultEntryBufferSize),
		mu:              &sync.Mutex{},
	}
}

func (sh *shard) put(hashedKey uint64, value []byte) error {
	return nil
}

func (sh *shard) get(hashedKey uint64) ([]byte, error) {
	return nil, nil
}

func (s *shard) cleanupExpiredEntries() {
	// Unimplemented
}
