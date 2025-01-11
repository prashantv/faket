// Package retryt can retry test assertions
// that may need multiple attempts to succeed.
//
// This is useful for tests that are waiting for some background work.
package retryt

import (
	"testing"

	"github.com/prashantv/faket"
)

// N runs `testFn“ upto N times, with [RetryBackoff] for retries.
// See [Run] for more details.
func N(t testing.TB, n int, testFn func(testing.TB)) {
	t.Helper()

	Opts{
		Attempts: n,
		Retry:    RetryBackoff,
	}.Run(t, testFn)
}

// Run tries `testFn` till it succeeds, or fails too many times.
// `opts` is used to customize retry behaviour.
//
// A last attempt is run against `t` so logs and failures are reported.
func Run(t testing.TB, opts Opts, testFn func(testing.TB)) {
	t.Helper()

	opts.setDefaults()
	for attempt := 1; attempt < opts.Attempts; attempt++ {
		tr := faket.RunTest(testFn)

		if !tr.Failed() {
			opts.Passed(t, attempt, tr)
			break
		}

		opts.Retry(t, attempt, tr)
	}

	testFn(t)
}

// Run is an alias for [Run].
func (os Opts) Run(t testing.TB, testFn func(testing.TB)) {
	Run(t, os, testFn)
}
