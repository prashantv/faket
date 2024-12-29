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

func (r TestResult) Skipped() bool {
	// If a test is failed and skipped, failed takes precedence, from SkipNow docs:
	// If a test fails (see Error, Errorf, Fail) and is then skipped,
	// it is still considered to have failed.
	if r.Failed() {
		return false
	}
	return r.res.Skipped()
}

// FailedAndSkipped allows determine whether a test marked as Failed was also skipped.
// This is not possible to determine using the stdlib.
func (r TestResult) FailedAndSkipped() bool {
	return r.res.Failed() && r.res.Skipped()
}

func (r TestResult) Failed() bool {
	return r.res.Failed()
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
		line := fmt.Sprintf("%v: %v", getFileLine(l.callers), l.entry)
		logs = append(logs, line)
	}
	return logs
}

func getFileLine(callers []uintptr) string {
	frames := runtime.CallersFrames(callers)
	f, _ := frames.Next()
	if f == (runtime.Frame{}) {
		return ""
	}

	return fmt.Sprintf("%v:%v", filepath.Base(f.File), f.Line)
}
