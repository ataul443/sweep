package sweep

import (
	"strconv"
	"testing"
)

func Benchmark1MillionEntriesWrite(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		writeToSweep(1000000)
	}
}

func Benchmark1EntriesWrite(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		writeToSweep(1)
	}
}

func writeToSweep(entriesCount int) {
	cache := Default()
	defer cache.Close()

	for i := 0; i < entriesCount; i++ {
		p := strconv.Itoa(i)
		err := cache.Put(p, []byte(p))
		if err != nil {
			panic(err)
		}
	}
}
