package sweep

import "fmt"

type Sweep struct {
	cache map[string][]byte
}

// Get retrieves value associated with the key from the sweep.
func (s *Sweep) Get(key string) (value []byte, err error) {
	b, ok := s.cache[key]
	if !ok {
		err = fmt.Errorf("key %s, not found in sweep", key)
		return
	}

	value = b
	return
}

// Put inserts the value associated with the key into the sweep.
func (s *Sweep) Put(key string, value []byte) {
	s.cache[key] = value
	return
}

func Default() *Sweep {
	return &Sweep{make(map[string][]byte)}
}
