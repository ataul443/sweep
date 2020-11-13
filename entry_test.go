package sweep

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEntry(t *testing.T) {
	hardCodedFrame := []byte{19, 0, 0, 0, 131, 79, 174, 95, 0, 0, 0, 0, 112, 105, 107, 97, 99, 104, 117}
	var hardCodedTimeStamp int64 = 1605259139
	hardCodedVal := "pikachu"

	t.Run("return non nil err when buff provided is short for Read", func(t *testing.T) {
		buff := make([]byte, 1)

		e := &entry{
			val:       []byte(hardCodedVal),
			createdAt: hardCodedTimeStamp,
		}

		_, err := e.Read(buff)
		assert.EqualError(t, err, ErrShortBuffer.Error(), "err should be short buffer")
	})

	t.Run("return non nil err when buff provided is short for write", func(t *testing.T) {
		buff := []byte{19, 0, 0, 0}

		e := &entry{}
		_, err := e.Write(buff)
		assert.EqualError(t, err, ErrShortWrite.Error(), "err should be short write")
	})

	t.Run("return nil err when buff provided is adequate for Read", func(t *testing.T) {
		buff := make([]byte, 19)

		e := &entry{
			val:       []byte(hardCodedVal),
			createdAt: hardCodedTimeStamp,
		}

		_, err := e.Read(buff)
		assert.NoError(t, err, "err should be nil")
	})



	t.Run("read valid frame into provided buff", func(t *testing.T) {
		buff := make([]byte, 19)

		e := &entry{val: []byte(hardCodedVal), createdAt: hardCodedTimeStamp}

		n, err := e.Read(buff)
		assert.NoError(t, err, "err should be nil")
		assert.Equalf(t, 19, n, "expected 19, got %d", n)
		assert.Equal(t, hardCodedFrame, buff, "frame should match")
	})

	t.Run("write valid frame from provided buff", func(t *testing.T) {
		e := &entry{}
		n, err := e.Write(hardCodedFrame)
		assert.NoError(t, err, "err should be nil")
		assert.Equalf(t, 19, n, "expected 19, got %d", n)

		assert.Equalf(t, hardCodedVal, string(e.val), "expected `%s`, got `%s`", hardCodedVal, e.val)
		assert.Equalf(t, hardCodedTimeStamp, e.createdAt, "expected %d, got %d", hardCodedTimeStamp, e.createdAt)
	})

}
