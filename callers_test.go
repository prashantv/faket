package faket

import (
	"path/filepath"
	"testing"

	"github.com/prashantv/faket/internal/want"
)

func TestGetCallers(t *testing.T) {
	t.Run("skip all", func(t *testing.T) {
		pc := getCallers(1000)
		if len(pc) > 0 {
			t.Fatalf("expected empty pc, got:\n%v", pc)
		}
	})

	t.Run("large depth", func(t *testing.T) {
		const depth = 128
		pc := recurse(depth, func() []uintptr {
			return getCallers(0)
		})
		if len(pc) < depth {
			t.Fatalf("len(pc) = %v < %v", len(pc), depth)
		}
	})
}

func TestGetCaller(t *testing.T) {
	t.Run("skip all", func(t *testing.T) {
		pc := getCaller(1000)
		want.Equal(t, "pc", pc, 0)
	})

	t.Run("skip 0", func(t *testing.T) {
		got := pcToFunction(getCallerCaller(0))
		want.Equal(t, "caller function", filepath.Base(got), "faket.getCallerCaller")
	})

	t.Run("skip 1", func(t *testing.T) {
		got := pcToFunction(getCallerCaller(1))
		want.Contains(t, "caller function", got, "TestGetCaller")
	})
}

func getCallerCaller(skip int) uintptr {
	return getCaller(skip)
}

func recurse[T any](n int, fn func() T) T {
	if n == 0 {
		return fn()
	}
	return recurse(n-1, fn)
}
