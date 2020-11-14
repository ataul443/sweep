package entry

import (
	"encoding/binary"
	"errors"
)

const (
	frameLenLegth = 4 // bytes

	timestampLength = 8 // bytes

	hashedKeyLength = 8 // bytes
)

var (
	ErrEntryShortBuffer = errors.New("short buffer to read frame into")

	ErrEntryShortWrite = errors.New("short buffer to write from")
)

// Frame is a binary representation of entry.
type Frame []byte

func (f Frame) Len() int {
	return len(f)
}

func ReadEntryIntoBuffer(hashedKey uint64, timestamp int64, val []byte, buf []byte) (int, error) {
	frameLenNeeded := FrameLen(val)

	if frameLenNeeded > len(buf) {
		return 0, ErrEntryShortBuffer
	}

	binary.LittleEndian.PutUint32(buf, uint32(frameLenNeeded))
	binary.LittleEndian.PutUint64(buf[frameLenLegth:], uint64(timestamp))
	binary.LittleEndian.PutUint64(buf[frameLenLegth+timestampLength:], hashedKey)

	copy(buf[frameLenLegth+timestampLength+hashedKeyLength:], val)

	return frameLenNeeded, nil
}

func GetEntryFromFrame(frame Frame) (hashedKey uint64, timestamp int64, val []byte, err error) {
	frameLen := binary.LittleEndian.Uint32(frame)

	if frameLen > uint32(len(frame)) {
		err =  ErrEntryShortWrite
		return
	}

	timestamp = int64(binary.LittleEndian.Uint64(frame[frameLenLegth:]))


	hashedKey = binary.LittleEndian.Uint64(frame[frameLenLegth+timestampLength:])

	val = make([]byte, frameLen-(frameLenLegth+timestampLength+hashedKeyLength))

	copy(val, frame[frameLenLegth+timestampLength+hashedKeyLength:frameLen])
	return
}

func ValFromFrame(frame Frame) ([]byte, error) {
	_, _, val, err := GetEntryFromFrame(frame)
	if err != nil {
		return nil, err
	}

	return val, err
}

func TimestampFromFrame(frame Frame) (int64, error) {
	_, tm, _, err := GetEntryFromFrame(frame)
	if err != nil {
		return 0, err
	}

	return tm, err
}

func FrameLen(val []byte) int {
	return frameLenLegth + timestampLength + hashedKeyLength + len(val)
}
