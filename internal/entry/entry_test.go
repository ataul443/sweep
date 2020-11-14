package entry

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEntry(t *testing.T) {
	hardCodedFrame := []byte{27, 0, 0, 0, 161, 183, 175, 95, 0, 0, 0, 0, 78, 97, 188,
		0, 0, 0, 0, 0, 112, 105, 107, 97, 99, 104, 117,}

	var hardCodedTimeStamp int64 = 1605351329
	var hardCodedHashKey uint64 = 12345678
	hardCodedVal := "pikachu"

	t.Run("return non nil err when buff provided is short for Read", func(t *testing.T) {
		buff := make([]byte, 1)

		_, err := ReadEntryIntoBuffer(hardCodedHashKey, hardCodedTimeStamp, []byte(hardCodedVal), buff)
		assert.EqualError(t, err, ErrEntryShortBuffer.Error(), "err should be short buffer")
	})

	t.Run("return non nil err when buff provided is short for write", func(t *testing.T) {
		buff := []byte{19, 0, 0, 0}

		_, _ , _, err := GetEntryFromFrame(buff)
		assert.EqualError(t, err, ErrEntryShortWrite.Error(), "err should be short write")
	})

	t.Run("return nil err when buff provided is adequate for Read", func(t *testing.T) {
		buff := make([]byte, FrameLen([]byte(hardCodedVal)))

		_, err := ReadEntryIntoBuffer(hardCodedHashKey, hardCodedTimeStamp, []byte(hardCodedVal), buff)
		assert.NoError(t, err, "err should be nil")
	})

	t.Run("read valid frame into provided buff", func(t *testing.T) {
		fl := FrameLen([]byte(hardCodedVal))
		buff := make([]byte, fl)

		n, err := ReadEntryIntoBuffer(hardCodedHashKey, hardCodedTimeStamp, []byte(hardCodedVal), buff)
		assert.NoError(t, err, "err should be nil")
		assert.Equalf(t, fl, n, "expected %d, got %d", fl, n)
		assert.Equal(t, hardCodedFrame, buff, "frame should match")
	})

	t.Run("write valid frame from provided buff", func(t *testing.T) {
		hk, tm, val, err := GetEntryFromFrame(hardCodedFrame)
		assert.NoError(t, err, "err should be nil")

		assert.Equalf(t, hardCodedVal, string(val), "expected `%s`, got `%s`",
			hardCodedVal, val)
		assert.Equalf(t, hardCodedTimeStamp, tm, "expected %d, got %d",
			hardCodedTimeStamp, tm)
		assert.Equalf(t, hardCodedHashKey, hk, "expected hashed key %d, got %d",
			hardCodedHashKey, hk)
	})

}
