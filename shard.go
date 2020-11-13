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
