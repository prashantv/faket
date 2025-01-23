package syncutil

import (
	"sync/atomic"
	"testing"

	"github.com/prashantv/faket/internal/want"
)

func TestRunN(t *testing.T) {
	const n = 10
	var counter atomic.Int32
	done := make(chan struct{})
	RunN(n, func(i int) {
		if counter.Add(1) == n {
			close(done)
		}
		<-done
	})
	want.Equal(t, "counter", counter.Load(), 10)
}
