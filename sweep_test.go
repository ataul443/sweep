package sweep

import (
	"fmt"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type keyValPayload struct {
	Key   string
	Value []byte
}

func TestSweep(t *testing.T) {
	cache := Default()

	tcs := []keyValPayload{
		{"putKey1", []byte("valueofputKey1")},
		{"putKey100", []byte("valueofputKey100")},
		{"putKey985", []byte("valueofputKey985")},
	}

	for _, v := range tcs {
		err := cache.Put(v.Key, v.Value)
		assert.NoErrorf(t, err, "put should be successful with key %s", v.Key)

		actualVal, err := cache.Get(v.Key)
		assert.NoErrorf(t, err, "get should be successful for %s", v.Key)

		assert.Equalf(t, v.Value, actualVal, "expected %s, got %s",
			v.Value, actualVal)
	}
}

func TestSweepEntryExpiration(t *testing.T) {
	cfg := Configuration{ShardsCount: 1000, EntryLifetime: time.Second}
	cache := New(cfg)

	tcs := []keyValPayload{
		{
			"putKeyWith1SecondExpiry1",
			[]byte("valueofputKeyWith1SecondExpiry1"),
		},
		{
			"putKeyWith1SecondExpiry100",
			[]byte("valueofputKeyWith1SecondExpiry100"),
		},
		{
			"putKeyWith1SecondExpiry985",
			[]byte("valueofputKeyWith1SecondExpiry985"),
		},
	}

	t.Run("cache should not have any expired entries", func(t *testing.T) {
		for _, v := range tcs {
			err := cache.Put(v.Key, v.Value)
			assert.NoErrorf(t, err, "put should be successful with key %s", v.Key)
		}

		time.Sleep(2 * time.Second)

		for _, v := range tcs {
			_, err := cache.Get(v.Key)
			assert.NotNilf(t, err, "get should not be successful for %s", v.Key)
		}
	})

	t.Run("cache should contain entries which doesn't expired yet", func(t *testing.T) {
		for _, v := range tcs {
			err := cache.Put(v.Key, v.Value)
			assert.NoErrorf(t, err, "put should be successful with key %s", v.Key)

			actualVal, err := cache.Get(v.Key)
			assert.NoErrorf(t, err, "get should be successful for %s", v.Key)

			assert.Equalf(t, v.Value, actualVal, "expected %s, got %s",
				v.Value, actualVal)
		}
	})
}

func TestGCPause(t *testing.T) {

	start := time.Now()

	cache := Default()
	for i := 0; i < 1e7; i++ {
		p := strconv.Itoa(i)
		err := cache.Put(p, []byte(p))
		if err != nil {
			panic(err)
		}

		// fmt.Println("Iteration:", p)
	}

	runtime.GC()

	val, _ := cache.Get("100")
	fmt.Println("100", string(val))

	fmt.Printf("GC Pause took: %.2f seconds\n.", time.Since(start).Seconds())
}
