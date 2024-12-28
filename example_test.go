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

func ExampleFail() {
	res := faket.RunTest(func(t testing.TB) {
		StrContainsInOrder(t, "help test foo", "test", "helper")
	})

	fmt.Println("Failed:", res.Failed())
	fmt.Println("Logs:", res.Logs())

	// Output:
	// Failed: true
	// Logs: failed to find "helper" in remaining string " foo"
}

func TestStrContainsInOrder(t *testing.T) {
	const str = "test helper validation"
	tests := []struct {
		name string
		strs []string
		fail bool
	}{
		{
			name: "substring",
			strs: []string{"helper"},
		},
		{
			name: "substrings in order",
			strs: []string{"test", "helper", "validation"},
		},
		{
			name: "substrings in wrong order",
			strs: []string{"test", "validation", "helper"}, // wrong order
			fail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := faket.RunTest(func(t testing.TB) {
				StrContainsInOrder(t, str, tt.strs...)
			})
			if res.Failed() != tt.fail {
				t.Errorf("Failed: got %v, want %v", res.Failed(), tt.fail)
			}
		})
	}
}
