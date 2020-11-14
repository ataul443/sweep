package sweep

import (
	"fmt"
	"testing"
)

func Benchmark1EntriesPut(b *testing.B) {
	cache := Default()
	defer cache.Close()

	val := []byte("pikachu")
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		p := fmt.Sprintf("%d", i)
		err := cache.Put(p, val)
		if err != nil {
			panic(err)
		}
	}
}
