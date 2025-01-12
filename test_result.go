package faket

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/prashantv/faket/internal/sliceutil"
)

// TestResult is the result of runnming a test against a fake [testing.TB].
type TestResult struct {
	res *fakeTB
}

// Logs represents a list of logged entries.
type Logs []Log

// Log is a single log entry, along with caller information.
type Log struct {
	Message    string
	CallerFile string
	CallerLine int
	CallerFunc string

	// TBFunc is the testing.TB function that generated this log message.
	TBFunc string
}

// Failed reports if a test failed.
func (r TestResult) Failed() bool {
	return r.res.Failed()
}

// Panicked reports if a test panicked.
func (r TestResult) Panicked() bool {
	return r.res.panicked
}

// Skipped reports if a test was skipped.
//
// If a test failed before it was skipped, then Failed takes precedence
// and Skipped returns false. To check if the test was skipped after a failure
// see [FailedAndSkipped].
func (r TestResult) Skipped() bool {
	// Above behaviour is for consistency with testing.TB, from SkipNow docs:
	// > If a test fails (see Error, Errorf, Fail) and is then skipped,
	// > it is still considered to have failed.
	if r.Failed() {
		return false
	}
	return r.res.Skipped()
}

// FailedAndSkipped reports if a test failed, and then was skipped.
//
// See [Skipped] for more details for how this differs from using
// [Failed] and [Skipped].
func (r TestResult) FailedAndSkipped() bool {
	return r.res.Failed() && r.res.Skipped()
}

// Helpers returns a list of functions that have called [testing.TB].Helper.
// The returned list is sorted by the full package+function.
func (r TestResult) Helpers() []string {
	funcs := r.res.helperFuncs()
	sort.Strings(funcs)
	return funcs
}

// Logs returns a list of log entries logged by the test.
func (r TestResult) Logs() Logs {
	return sliceutil.Map(r.res.logs, r.res.toLog)
}

// Messages returns a list of individual logs.
func (ls Logs) Messages() []string {
	return sliceutil.Map(ls, func(l Log) string {
		return l.Message
	})
}

// String returns the log output, as it would be printed by `go test`
// with caller information.
func (ls Logs) String() string {
	var buf strings.Builder
	for _, l := range ls {
		fmt.Fprintf(&buf, "%s:%d: %v\n", filepath.Base(l.CallerFile), l.CallerLine, l.Message)
	}
	return buf.String()
}
