package sweep

import (
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
		cache.Put(v.Key, v.Value)

		actualVal, err := cache.Get(v.Key)
		assert.NoErrorf(t, err, "get should be successful for %s", v.Key)

		assert.Equalf(t, v.Value, actualVal, "expected %s, got %s",
			v.Value, actualVal)
	}
}

func TestSweepEntryExpiration(t *testing.T) {
	cache := New(time.Second)

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
			cache.Put(v.Key, v.Value)
		}

		time.Sleep(2 * time.Second)

		for _, v := range tcs {
			_, err := cache.Get(v.Key)
			assert.NotNilf(t, err, "get should not be successful for %s", v.Key)
		}
	})

	t.Run("cache should contain entries which doesn't expired yet", func(t *testing.T) {
		for _, v := range tcs {
			cache.Put(v.Key, v.Value)

			actualVal, err := cache.Get(v.Key)
			assert.NoErrorf(t, err, "get should be successful for %s", v.Key)

			assert.Equalf(t, v.Value, actualVal, "expected %s, got %s",
				v.Value, actualVal)
		}
	})
}
