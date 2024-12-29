package faket

import (
	"strings"
	"testing"
)

// MustPass will ensure the test passed or `t.Fatal` otherwise.
func (tr TestResult) MustPass(t testing.TB) {
	t.Helper()

	if tr.Failed() {
		t.Fatalf("test failed, logs:\n%v", tr.Logs())
	}
}

// MustFail ensures that the test failed, and the given log message
// is found in the test logs.
func (tr TestResult) MustFail(t testing.TB, wantLog string) {
	t.Helper()

	if !tr.Failed() {
		t.Fatal("test passed, but expected to fail")
	}

	if !strings.Contains(tr.Logs(), wantLog) {
		t.Fatalf("test expected to fail, missing expected log %q. logs:\n%v", wantLog, tr.Logs())
	}
}
