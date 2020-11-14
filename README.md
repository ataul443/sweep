# Sweep

Sweep is an in-memory cache. It uses sharding techinque to store millions of entries inside it. Even with millions of entries inside it the GC pause takes only a few milliseconds.

**Note**: It is a toy project, not intended to be used in production. I was practicing for an LLD interview and I ended up writing this.

## Usage
```go
package main

import (
	"fmt"
	"github.com/ataul443/sweep"
)

func main() {
	cache := sweep.Default()
	defer cache.Close()
	
	err := cache.Put("pikachu", []byte("value of pikachu"))
	if err != nil {
		panic(err)
	}
	
	val, err := cache.Get("pikachu")
	if err != nil {
		panic(err)
	}
	
	fmt.Println(string(val))
}

```

## How GC Pause minimized ?
```
Starting GC Pause benchmark....
Going to put 20000000 entries in sweep.
Time took in inserting 20000000 entries in sweep: 18.93 seconds
GC Pause took: 85 milliseconds
Total entries in cache: 20000000
```

Instead of storing _(key, val)_ pair directly in `map[string]*entry`, Sweep stores the _val_ in a queue and stores the index in a `map[uint64]uint64` as an _(key, index)_ pair. Large map containing pointers causes huge GC pause time, however this is not same with map containing integers. You can read up more about this at [issue-9477](https://github.com/golang/go/issues/9477).
