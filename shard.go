package sweep

import (
	"sweep/internal/entry"
	"sync"
	"time"
)

type shard struct {
	hashIndexBucket map[uint64]int

	queue *entry.Queue

	maxSize int

	mu *sync.RWMutex
}

func newShard(maxSize int) *shard {
	return &shard{
		hashIndexBucket: make(map[uint64]int),
		queue:           entry.NewQueue(maxSize),
		maxSize:         maxSize,
		mu:              &sync.RWMutex{},
	}
}

func (sh *shard) put(hashedKey uint64, timestamp int64, val []byte) error {
	sh.mu.Lock()
	defer sh.mu.Unlock()

	spaceExist := sh.queue.SpaceAvailable(entry.FrameLen(val))
	if !spaceExist {
		err := sh.queue.Grow()
		if err != nil {
			return err
		}
	}

	idx, err := sh.queue.Push(hashedKey, timestamp, val)
	if err != nil {
		return err
	}

	sh.hashIndexBucket[hashedKey] = idx
	return nil
}

func (sh *shard) get(hashedKey uint64) ([]byte, error) {
	sh.mu.RLock()
	defer sh.mu.RUnlock()

	idx, ok := sh.hashIndexBucket[hashedKey]
	if !ok {
		return nil, ErrEntryNotFound
	}

	frame, err := sh.queue.PeekAt(idx)
	if err != nil {
		return nil, err
	}

	val, err := entry.ValFromFrame(frame)
	if err != nil {
		return nil, err
	}

	return val, nil
}

func (sh *shard) cleanupExpiredEntries(entryLifetime time.Duration) (int, error) {
	sh.mu.Lock()
	defer sh.mu.Unlock()

	poppedCount := 0

	for {
		frame, err := sh.queue.Front()
		if err != nil {
			if err == entry.ErrQueueEmpty {
				return poppedCount, nil
			}

			return poppedCount, err
		}

		hk, tm, _, err := entry.GetEntryFromFrame(frame)
		if err != nil {
			return poppedCount, err
		}

		if time.Since(time.Unix(tm, 0)) > entryLifetime {
			_, err = sh.queue.Pop()
			if err != nil {
				return poppedCount, err
			}

			// delete the key from map
			delete(sh.hashIndexBucket, hk)
			poppedCount += 1

		} else {
			break
		}
	}

	return poppedCount, nil
}
