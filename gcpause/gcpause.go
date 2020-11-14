package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/ataul443/sweep"
)

func main() {
	fmt.Println("Starting GC Pause benchmark....")
	cfg := sweep.Configuration{
		ShardsCount:   512,
		EntryLifetime: 20 * time.Minute,
		MaxShardSize:  0,
	}

	cache := sweep.New(cfg)

	val := []byte("cool cache, brother.")

	entriesCount := 20000000

	fmt.Printf("Going to put %d entries in sweep.\n", entriesCount)

	cachePutStartedAt := time.Now()
	for i := 0; i < entriesCount; i++ {
		key := fmt.Sprintf("key_%d", i)
		err := cache.Put(key, val)
		if err != nil {
			panic(err)
		}
	}

	timeTookInPut := time.Since(cachePutStartedAt).Seconds()
	fmt.Printf("Time took in inserting %d entries in sweep: %.2f seconds\n",
		entriesCount, timeTookInPut)

	gcStartedAt := time.Now()
	runtime.GC()
	fmt.Printf("GC Pause took: %d milliseconds\n", time.Since(gcStartedAt).Milliseconds())

	_, err := cache.Get("key_1234")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Total entries in cache: %d\n", cache.EntriesCount())
}
