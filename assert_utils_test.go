package faket

import (
	"reflect"
	"testing"
)

func wantNoErr(t testing.TB, err error) {
	if err == nil {
		return
	}

	t.Fatalf("expected no error, got %v", err)
}

func wantEqual[T comparable](t testing.TB, msg string, got, want T) {
	if want == got {
		return
	}

	t.Fatalf("expected == %s: got:  %v\nwant: %v", msg, want, got)
}

func wantDeepEqual[T any](t testing.TB, msg string, got, want T) {
	if reflect.DeepEqual(want, got) {
		return
	}

	t.Fatalf("expected DeepEqual %s: got:  %v\nwant: %v", msg, want, got)
}
