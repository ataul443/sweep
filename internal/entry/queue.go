package entry

import (
	"encoding/binary"
	"errors"

	"github.com/ataul443/sweep/bipbuffer"
)

const defaultEntryQueueSize = 4 * 1024 // 4KB

var (
	ErrQueueSpaceNotAvailable = errors.New("no space available in queue")

	ErrQueueMaxSizeReaced = errors.New("max queue size reached")

	ErrQueueEmpty = errors.New("queue is empty")
)

type Queue struct {
	bipbuf  *bipbuffer.BipBuffer
	maxSize int
}

func NewQueue(maxSize int) *Queue {
	return &Queue{
		bipbuf:  bipbuffer.New(defaultEntryQueueSize),
		maxSize: maxSize,
	}
}

// Push attempt to return an index where the queue is pushed otherwise error.
func (q *Queue) Push(hashedKey uint64, timestamp int64, val []byte) (int, error) {
	frameSize := FrameLen(val)

	idx, b := q.bipbuf.Reserve(frameSize)
	if b == nil || len(b) < frameSize {
		return 0, ErrQueueSpaceNotAvailable
	}

	k, err := ReadEntryIntoBuffer(hashedKey, timestamp, val, b)
	if err != nil {
		// This should never happen
		return 0, err
	}

	q.bipbuf.Commit(k)
	return idx, nil
}

// Pop attempt to return an entry frame from the front of the queue and returns it
// otherwise error.
func (q *Queue) Pop() (Frame, error) {
	frame, err := q.Front()
	if err != nil {
		return nil, err
	}

	q.bipbuf.Decommit(len(frame))
	return frame, nil
}

// Front attempt to return an entry frame from the front of the queue without
// removing it otherwise error.
func (q *Queue) Front() (Frame, error) {
	b := q.bipbuf.GetContiguousBlock()
	if b == nil {
		return nil, ErrQueueEmpty
	}

	frameSize := binary.LittleEndian.Uint32(b)
	return b[:frameSize], nil
}

// Peek attempt to return an entry frame at an index in the queue otherwise
// error.
func (q *Queue) PeekAt(idx int) (Frame, error) {
	if idx < 0 || idx >= q.bipbuf.Capacity() {
		return nil, errors.New("invalid index")
	}

	b, err := q.bipbuf.PeekAt(idx, frameLenLegth)
	if err != nil {
		return nil, err
	}

	frameLen := binary.LittleEndian.Uint32(b)

	frame, err := q.bipbuf.PeekAt(idx, int(frameLen))
	if err != nil {
		return nil, err
	}

	return frame, nil
}

// Grow will increase the queue size to twice of current size with all data
// intact. It throws error, if queue size reached it max limits.
func (q *Queue) Grow() error {
	if q.maxSize != 0 && (2*q.Capacity()) > q.maxSize {
		return ErrQueueMaxSizeReaced
	}

	q.bipbuf.Grow()
	return nil
}

// Capacity returns the total capacity of the queue.
func (q *Queue) Capacity() int {
	return q.bipbuf.Capacity()
}

// SpaceAvailable returns true if queue has space for write equal to size.
func (q *Queue) SpaceAvailable(size int) bool {
	_, b := q.bipbuf.Reserve(size)
	if b == nil || len(b) < size {
		return false
	}

	return true
}
