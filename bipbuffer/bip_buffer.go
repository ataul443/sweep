package bipbuffer

import "errors"

// BipBuffer is a spsc circular non-thread safe buffer that always
// supports writing a contiguous chunk of data. Write requests that
// cannot fit in an available contiguous area will be failed with
// an error.
type BipBuffer struct {
	buf []byte

	idxRegionA    int
	sizeOfRegionA int

	idxRegionB    int
	sizeOfRegionB int

	idxReserve    int
	sizeOfReserve int

}

func New(size uint64) *BipBuffer {
	return &BipBuffer{buf: make([]byte, size)}
}

// Reserve reserves a space for writing at an index in the buffer,
// of length equal to size, and returns that index and space as byte slice.
// It will return nil slice if requested size is not available in buffer.
func (bb *BipBuffer) Reserve(size int) (int, []byte) {
	if bb.sizeOfRegionB != 0 {

		availableSpace := bb.getFreeSpaceInRegionB()
		if availableSpace == 0 {
			return 0, nil
		}

		if size > availableSpace {
			return 0, nil
		}

		bb.sizeOfReserve = size
		bb.idxReserve = bb.idxRegionB + bb.sizeOfRegionB

	} else {
		// Only region A is present
		availableSpace := bb.getFreeSpaceAfterRegionA()

		if availableSpace >= bb.idxRegionA {

			if availableSpace == 0 {
				return 0, nil
			}

			if size > availableSpace {
				return 0, nil
			}

			bb.sizeOfReserve = size
			bb.idxReserve = bb.idxRegionA + bb.sizeOfRegionA
		} else {

			if bb.idxRegionA == 0 {
				return 0, nil
			}

			if bb.sizeOfRegionA < size {
				return 0, nil
			}

			bb.sizeOfReserve = size
			bb.idxReserve = 0
		}
	}

	return bb.idxReserve, bb.buf[bb.idxReserve : bb.idxReserve+bb.sizeOfReserve]
}

// Commit commits reserved space of length equal to size in the buffer.
// It makes the data in the reserved space permanent in buffer. If the
// asked size is more than current reserved space, it will simply
// commits the current reserved space.
func (bb *BipBuffer) Commit(size int) {
	if size > bb.sizeOfReserve {
		size = bb.sizeOfReserve
	}

	if bb.sizeOfRegionA == 0 && bb.sizeOfRegionB == 0 {
		bb.idxRegionA = bb.idxReserve
		bb.sizeOfRegionA = size

		bb.idxReserve = 0
		bb.sizeOfReserve = 0
		return
	}

	if bb.idxReserve == (bb.idxRegionA + bb.sizeOfRegionA) {
		bb.sizeOfRegionA += size
	} else {
		bb.sizeOfRegionB += size
	}

	bb.idxReserve = 0
	bb.sizeOfReserve = 0
	return
}

// Decommit frees already committed space of length size in the buffer.
func (bb *BipBuffer) Decommit(size int) {
	if size >= bb.sizeOfRegionA {
		bb.idxRegionA = bb.idxRegionB
		bb.sizeOfRegionA = bb.sizeOfRegionB

		bb.idxRegionB = 0
		bb.sizeOfRegionB = 0
	} else {
		bb.idxRegionA += size
		bb.sizeOfRegionA -= size
	}
}

// GetContiguousBlock returns byte slice representing single
// contiguous region in the buffer. To read all data out of the buffer
// call this method in loop. It will return nil slice when
// there will be no committed region.
func (bb *BipBuffer) GetContiguousBlock() []byte {
	if bb.sizeOfRegionA == 0 {
		return nil
	}

	return bb.buf[bb.idxRegionA : bb.idxRegionA+bb.sizeOfRegionA]
}

// PeekAt returns a byte slice representing a region starting at
// index idx of length size. It will throw error when idx doesn't
// belong to any region inside the buffer.
func (bb *BipBuffer) PeekAt(idx, size int) ([]byte, error) {
	if !bb.isAreaInRegionA(idx, size) && !bb.isAreaInRegionB(idx, size) {
		return nil, errors.New("invalid index")
	}

	return bb.buf[idx : idx+size], nil
}

// Capacity returns capacity of the buffer.
func (bb *BipBuffer) Capacity() int {
	return cap(bb.buf)
}

// CommittedSize returns total committed space in the buffer.
func (bb *BipBuffer) CommittedSize() int {
	return bb.sizeOfRegionA + bb.sizeOfRegionB
}

// Grow will increase the underlying buffer size to twice
// of the current size.
func (bb *BipBuffer) Grow() {
	newBuf := make([]byte, 2*cap(bb.buf))

	n := 0
	for {
		b := bb.GetContiguousBlock()
		if b == nil {
			break
		}

		k := copy(newBuf[n:], b)
		n += k

		bb.Decommit(k)
	}

	bb.buf = newBuf[:n]
	bb.idxReserve = 0
	bb.sizeOfReserve = n

	bb.Commit(n)
}

func (bb *BipBuffer) isAreaInRegionA(idx, size int) bool {
	if bb.sizeOfRegionA == 0 {
		return false
	}

	return idx >= bb.idxRegionA &&
		idx <= (bb.idxRegionA+bb.sizeOfRegionA) &&
		bb.sizeOfRegionA >= size
}

func (bb *BipBuffer) isAreaInRegionB(idx, size int) bool {
	if bb.sizeOfRegionB == 0 {
		return false
	}

	return idx >= bb.idxRegionB &&
		idx <= (bb.idxRegionB+bb.sizeOfRegionB) &&
		bb.sizeOfRegionB >= size
}

func (bb *BipBuffer) getFreeSpaceAfterRegionA() int {
	return cap(bb.buf) - (bb.idxRegionA + bb.sizeOfRegionA)
}

func (bb *BipBuffer) getFreeSpaceInRegionB() int {
	// Region B stands before Region A
	return bb.idxRegionA - (bb.idxRegionB + bb.sizeOfRegionB)
}
