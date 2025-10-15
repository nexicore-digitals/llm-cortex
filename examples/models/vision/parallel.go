package vision

import (
	"sync"
)

func ParallelLoadExample() {
	var wg sync.WaitGroup

	// Define the list of example functions to run in parallel.
	examples := []func(){
		BlipExample,
		ClipExample,
		CliptionExample,
	}

	// Launch each example in its own goroutine.
	for _, example := range examples {
		wg.Add(1)
		go func(f func()) {
			defer wg.Done()
			f()
		}(example)
	}

	// Wait for all examples to complete.
	wg.Wait()
}