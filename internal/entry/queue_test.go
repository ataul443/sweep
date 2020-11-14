package entry

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	hardCodedTimeStamp int64  = 1605351329
	hardCodedHashKey   uint64 = 12345678
	hardCodedVal              = []byte("pikachu")
)

func TestQueue_Capacity(t *testing.T) {
	q := NewQueue(defaultEntryQueueSize)

	assert.Equalf(t, defaultEntryQueueSize, q.Capacity(),
		"expected default capacity %d, got %d", defaultEntryQueueSize,
		q.Capacity())
}

func TestQueue_SpaceAvailable(t *testing.T) {
	q := NewQueue(defaultEntryQueueSize)

	spaceAvailable := q.SpaceAvailable(defaultEntryQueueSize)
	assert.Equal(t, true, spaceAvailable, "space should be available")
}

func TestQueue_Push(t *testing.T) {
	q := NewQueue(defaultEntryQueueSize)
	idx, err := q.Push(hardCodedHashKey, hardCodedTimeStamp, hardCodedVal)
	assert.NoError(t, err, "push should be successful")

	frame, err := q.PeekAt(idx)
	assert.NoError(t, err, "peek should be successful")

	_, tm, val, err := GetEntryFromFrame(frame)
	assert.NoError(t, err, "entry write should be successful")

	assert.Equalf(t, hardCodedVal, val, "expected val %s, got %s",
		hardCodedVal, val)

	assert.Equalf(t, hardCodedTimeStamp, tm,
		"expected timestamp %d, got %d", hardCodedTimeStamp, tm)
}

func TestQueue_Pop(t *testing.T) {
	q := NewQueue(defaultEntryQueueSize)
	_, err := q.Push(hardCodedHashKey, hardCodedTimeStamp, hardCodedVal)
	assert.NoError(t, err, "push should be successful")

	_, err = q.Push(9876543, hardCodedTimeStamp+9, []byte("pikachu"))
	assert.NoError(t, err, "push should be successful")

	frame, err := q.Pop()
	assert.NoError(t, err, "pop should be successful")

	_, tm, val, err := GetEntryFromFrame(frame)
	assert.NoError(t, err, "entry write should be successful")

	assert.Equalf(t, hardCodedVal, val, "expected val %v, got %v",
		hardCodedVal, val)

	assert.Equalf(t, hardCodedTimeStamp, tm,
		"expected timestamp %d, got %d", hardCodedTimeStamp, tm)
}
