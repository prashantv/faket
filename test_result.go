package faket

import (
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
)

// TestResult is the result of runnming a test against a fake [testing.TB].
type TestResult struct {
	res *fakeTB
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

// Failed reports if a test failed.
func (r TestResult) Failed() bool {
	return r.res.Failed()
}

// Panicked reports if a test panicked.
func (r TestResult) Panicked() bool {
	return r.res.panicked
}

// LogsList returns a list of logs
func (r TestResult) LogsList() []string {
	logs := make([]string, 0, len(r.res.Logs))
	for _, l := range r.res.Logs {
		logs = append(logs, l.entry)
	}
	return logs
}

// Logs returns the log output of the test.
// It merges all logs strings into a single string.
func (r TestResult) Logs() string {
	return strings.Join(r.LogsList(), "\n")
}

// LogsWithCaller returns the log output of the test with caller information
// similar to the output of `go test`.
func (r TestResult) LogsWithCaller() []string {
	logs := make([]string, 0, len(r.res.Logs))
	for _, l := range r.res.Logs {
		ci := getCallerInfo(l.callers)
		line := fmt.Sprintf("%s:%d: %v", filepath.Base(ci.callerFile), ci.callerLine, l.entry)
		logs = append(logs, line)
	}
	return logs
}

type callerInfo struct {
	logFn string

	callerFile string
	callerLine int
}

func getCallerInfo(callers []uintptr) callerInfo {
	frames := runtime.CallersFrames(callers)

	f, _ := frames.Next()
	if f == (runtime.Frame{}) {
		return callerInfo{}
	}

	// First frame is the fake_tb caller.
	ci := callerInfo{
		logFn: f.Function,
	}

	skip := true
	for skip {
		// TODO: Skip t.Helper() frames.
		f, _ = frames.Next()
		if f == (runtime.Frame{}) {
			return ci
		}

		// When a defer is triggered by a panic, it's added to the trace
		// but panic is not shown as a log caller.
		skip = f.Function == "runtime.gopanic"
	}

	ci.callerFile = f.File
	ci.callerLine = f.Line
	return ci
}
