package sweep

import (
	"encoding/binary"
	"errors"
)

const (
	frameLenLegth = 4 // bytes

	timestampLength = 8 // bytes
)

var (
	ErrShortBuffer = errors.New("short buffer to read frame into")

	ErrShortWrite = errors.New("short buffer to write from")
)

type entry struct {
	val       []byte
	createdAt int64
}

func (e *entry) Read(b []byte) (int, error) {
	frameLenNeeded := frameLenLegth + timestampLength + len(e.val)

	if frameLenNeeded > len(b) {
		return 0, ErrShortBuffer
	}

	binary.LittleEndian.PutUint32(b, uint32(frameLenNeeded))
	binary.LittleEndian.PutUint64(b[frameLenLegth:], uint64(e.createdAt))
	copy(b[frameLenLegth+timestampLength:], e.val)

	return frameLenNeeded, nil
}

func (e *entry) Write(b []byte) (int, error) {
	frameLen := binary.LittleEndian.Uint32(b)

	if frameLen > uint32(len(b)) {
		return 0, ErrShortWrite
	}

	timestamp := binary.LittleEndian.Uint64(b[frameLenLegth:])

	e.createdAt = int64(timestamp)

	if e.val == nil {
		e.val = make([]byte, frameLen-(frameLenLegth+timestampLength))
	}

	copy(e.val, b[frameLenLegth+timestampLength:frameLen])
	return int(frameLen), nil
}
