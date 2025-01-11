package want

import (
	"reflect"
	"strings"
	"testing"
)

// NoErr asserts that the given err is nil.
func NoErr(t testing.TB, err error) {
	t.Helper()

	if err == nil {
		return
	}

	t.Fatalf("expected no error, got: %v", err)
}

// Equal asserts that the given got/want are equal using ==.
func Equal[T comparable](t testing.TB, msg string, got, want T) {
	t.Helper()

	if want == got {
		return
	}

	t.Fatalf("%s: expected equal\ngot:  %v\nwant: %v", msg, got, want)
}

// Contains asserts that the given substring is contained.
func Contains(t testing.TB, msg string, got, contains string) {
	t.Helper()

	if strings.Contains(got, contains) {
		return
	}

	t.Fatalf("%s: expected contains\ngot: %s\nwant contains: %s", msg, got, contains)
}

func NotContains(t testing.TB, msg string, got, notContains string) {
	t.Helper()

	if !strings.Contains(got, notContains) {
		return
	}

	t.Fatalf("%s: expected not contains\ngot: %s\nshould not contain: %s", msg, got, notContains)
}

// DeepEqual assers that the given got/want are equal using reflect.DeepEqual.
func DeepEqual[T any](t testing.TB, msg string, got, want T) {
	t.Helper()

	if reflect.DeepEqual(want, got) {
		return
	}

	t.Fatalf("%s: expected DeepEqual\ngot:  %v\nwant: %v", msg, got, want)
}
