package p1

import (
	"testing"

	"github.com/prashantv/faket/internal/cmptest"
)

func TestCmp_Panic(t *testing.T) {
	opts := cmptest.Opts{WantPanic: true}
	cmptest.CompareOpts(t, opts, func(t testing.TB) {
		defer t.Log("defer log")
		t.Log("normal log")

		panic("fatal")
	})
}
