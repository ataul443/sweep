package sweep

import (
	"testing"

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
		assert.NoError(t, err, "put should be successful for key1")

		actualVal, err := cache.Get(v.Key)
		assert.NoErrorf(t, err, "get should be successful for %s", v.Key)

		assert.Equalf(t, v.Value, actualVal, "expected %s, got %s",
			v.Value, actualVal)
	}
}
