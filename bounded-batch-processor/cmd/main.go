package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

func main() {
	batchSize := 10
	fmt.Printf("Starting batch process for %d items...\n", batchSize)

	results := processBatch(batchSize)

	fmt.Printf("Completed processed values: %v\n", results)
}

func processBatch(batchSize int) []int {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var processed []int

	// Using the Go 1.22 integer ranging syntax
	for i := range batchSize {
		wg.Add(1)

		// Because of Go 1.22's per-iteration variable scoping,
		// 'i' is safely captured by the goroutine. We no longer
		// need to manually shadow it (e.g., i := i).
		go func() {
			defer wg.Done()

			time.Sleep(rand.N(500 * time.Millisecond))

			mu.Lock()
			processed = append(processed, i)
			mu.Unlock()
		}()
	}

	wg.Wait()

	return processed
}
