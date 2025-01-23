// Package syncutil has helpers for running concurrent code.
package syncutil

import "sync"

// RunN runs n items in parallel, with no concurrency limiting.
func RunN(n int, fn func(i int)) {
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			fn(i)
		}(i)
	}
	wg.Wait()
}
