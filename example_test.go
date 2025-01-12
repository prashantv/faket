package faket_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/prashantv/faket"
)

func StrContainsInOrder(t testing.TB, s string, contains ...string) {
	for _, v := range contains {
		i := strings.Index(s, v)
		if i < 0 {
			// Failed to find element, report error.
			t.Errorf("failed to find %q in remaining string %q", v, s)
			return
		}

		// Continue searching for remaining elements after the current found value.
		s = s[i+len(v):]
	}
}

func ExampleRunTest_failure() {
	res := faket.RunTest(func(t testing.TB) {
		StrContainsInOrder(t, "help test foo", "test", "helper")
	})

	fmt.Println("Failed:", res.Failed())
	fmt.Println("Logs:", res.Logs().String())

	// Output:
	// Failed: true
	// Logs: example_test.go:16: failed to find "helper" in remaining string " foo"
}

func TestStrContainsInOrder(t *testing.T) {
	t.Run("correct order", func(t *testing.T) {
		faket.RunTest(func(t testing.TB) {
			StrContainsInOrder(t, "test helper function", "test", "helper")
		}).MustPass(t)
	})

	t.Run("incorrect order", func(t *testing.T) {
		faket.RunTest(func(t testing.TB) {
			StrContainsInOrder(t, "test helper function", "helper", "test")
		}).MustFail(t, `failed to find "test" in remaining string " function"`)
	})
}
