package faket

import (
	"path/filepath"
	"strings"
	"testing"
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
		if pc != 0 {
			t.Fatalf("expected empty pc, got %v", pc)
		}
	})

	t.Run("skip 0", func(t *testing.T) {
		s := pcToFunction(getCallerCaller(0))
		want := "faket.getCallerCaller"
		if got := filepath.Base(s); got != want {
			t.Fatalf("with no skip, got %v, want %v", got, want)
		}
	})

	t.Run("skip 1", func(t *testing.T) {
		got := pcToFunction(getCallerCaller(1))
		want := "TestGetCaller"
		if !strings.Contains(got, want) {
			t.Fatalf("with skip, got %v, want contains %v", got, want)
		}
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
