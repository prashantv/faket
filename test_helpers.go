package faket

import (
	"fmt"
	"strings"
	"testing"
)

// MustPass ensures the test passed.
// Otherwise, it will report a fatal failure to `t`.
func (tr TestResult) MustPass(t testing.TB) {
	t.Helper()

	if tr.Failed() {
		t.Fatalf("test failed, logs:\n%v", tr.Logs())
	}
}

// MustFail ensures that the test failed, and the given
// log message is found in the test logs.
// Otherwise, it will report a fatal failure to `t`.
func (tr TestResult) MustFail(t testing.TB, wantLog string) {
	t.Helper()

	if !tr.Failed() {
		t.Fatal("test passed, but expected to fail")
	}

	if !strings.Contains(tr.Logs().String(), wantLog) {
		t.Fatalf("test expected to fail, missing expected log %q. logs:\n%v", wantLog, tr.Logs())
	}
}

// MustPanic ensures that the test panicked, and the given
// substring is found in the recovered's value as a string.
// Otherwise, it will report a fatal failure to `t`.
func (tr TestResult) MustPanic(t testing.TB, contains string) {
	t.Helper()

	if !tr.Panicked() {
		t.Fatal("test did not panic, but expected to panic")
	}

	rec := tr.res.recovered
	if !strings.Contains(fmt.Sprint(rec), contains) {
		t.Fatalf("test expected to panic, panic string doesn't contain %q. got:\n%v", contains, rec)
	}
}
