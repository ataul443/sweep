package bipbuffer

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBipBuffer_Reserve(t *testing.T) {

	t.Run("return nil slice when requested space is too large",
		func(t *testing.T) {
			bb := New(2)
			n, b := bb.Reserve(8)
			assert.Equalf(t, 0, n, "expected index 0, got %d", n)
			assert.Nil(t, b, "expected nil, got %s", b)
		})

	t.Run("return slice len of requested space when requested space available",
		func(t *testing.T) {
			bb := New(32)
			requestedSize := 8
			n, b := bb.Reserve(requestedSize)
			assert.Equalf(t, 0, n, "expected index %d, got %d",
				requestedSize, n)

			assert.Equalf(t, requestedSize, len(b), "expected %d, got d", requestedSize, b)
		})
}

func TestBipBuffer_Commit(t *testing.T) {
	t.Run("commit space equal to all reserved space in the buffer", func(t *testing.T) {
		bb := New(32)
		_, _ = bb.Reserve(16)
		bb.Commit(28)
		k := bb.CommittedSize()
		assert.Equalf(t, 16, k,
			"expected committed size %d, got %d", 16, k)
	})

	t.Run("commit space equal to reserved space in the buffer", func(t *testing.T) {
		bb := New(32)
		requestedSize := 12
		_, _ = bb.Reserve(requestedSize)
		bb.Commit(requestedSize)
		k := bb.CommittedSize()
		assert.Equalf(t, requestedSize, k,
			"expected committed size %d, got %d", requestedSize, k)
	})

	t.Run("commit space less than reserved space in the buffer", func(t *testing.T) {
		bb := New(32)
		requestedSize := 12
		commitSize := 8
		_, _ = bb.Reserve(requestedSize)
		bb.Commit(commitSize)
		k := bb.CommittedSize()
		assert.Equalf(t, commitSize, k,
			"expected committed size %d, got %d", commitSize, k)
	})
}

func TestBipBuffer_Decommit(t *testing.T) {
	t.Run("decommit space equal to committed space in the buffer", func(t *testing.T) {
		bb := New(32)
		requestedSpace := 16
		_, _ = bb.Reserve(requestedSpace)
		bb.Commit(requestedSpace)
		bb.Decommit(requestedSpace)
		k := bb.CommittedSize()

		assert.Equalf(t, 0, k,
			"expected committed size %d, got %d", 0, k)
	})

	t.Run("decommit space less than committed space in the buffer", func(t *testing.T) {
		bb := New(32)
		requestedSize := 12
		decommitSize := 8
		_, _ = bb.Reserve(requestedSize)
		bb.Commit(requestedSize)
		bb.Decommit(decommitSize)

		k := bb.CommittedSize()
		assert.Equalf(t, 4, k,
			"expected committed size %d, got %d", 4, k)
	})

}

func TestBipBuffer_Grow(t *testing.T) {
	bb := New(8)
	bb.Grow()
	assert.Equalf(t, 16, bb.Capacity(),
		"expected capacity %d, got %d", 16, bb.Capacity())
}

func TestBipBuffer_PeekAt(t *testing.T) {
	bb := New(64)

	// Create Region A
	bb.idxRegionA = 40
	bb.sizeOfRegionA = 24

	// Create Region B
	bb.idxRegionB = 0
	bb.sizeOfRegionB = 16

	t.Run("return non nill error provided index doesn't lie on either region",
		func(t *testing.T) {
			requestedIdx := 32
			_, err := bb.PeekAt(requestedIdx, 2)
			assert.Error(t, err, "err should be non nil")
		})

	t.Run("return slice of requested size provided idx falls in region A", func(t *testing.T) {
		requestedIdx := 46
		b, err := bb.PeekAt(requestedIdx, 2)
		assert.NoError(t, err, "err should be nil")
		assert.Equalf(t, 2, len(b), "expected len %d, got %d", 2, len(b))
	})

	t.Run("return slice of requested size provided idx falls in region B", func(t *testing.T) {
		requestedIdx := 8
		b, err := bb.PeekAt(requestedIdx, 2)
		assert.NoError(t, err, "err should be nil")
		assert.Equalf(t, 2, len(b), "expected len %d, got %d", 2, len(b))
	})

}

func TestBipBuffer_Capacity(t *testing.T) {
	bb := New(8)
	assert.Equalf(t, 8, bb.Capacity(),
		"expected capacity %d, got %d", 8, bb.Capacity())
}

func TestBipBuffer_CommittedSize(t *testing.T) {
	bb := New(64)
	_, reservedSlice := bb.Reserve(15)
	assert.NotNil(t, reservedSlice, "reserved slice should not be nil")

	bb.Commit(7)
	assert.Equalf(t, 7, bb.CommittedSize(),
		"expected committed size %d, got %d")
}
