package faket

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"testing"
)

var _ testing.TB = (*recordingTB)(nil)

type recordingTB struct {
	testing.TB // embedded to include testing.TB.private, which we can't implement.

	mu sync.Mutex // protects all of the below fields.

	Cleanups []func()
	Logs     []logEntry

	completed chan struct{}
	failed    bool
	skipped   bool
}

type logEntry struct {
	callers []uintptr
	entry   string
}

func newRecordingTB() *recordingTB {
	return &recordingTB{
		completed: make(chan struct{}),
	}
}

func (tb *recordingTB) postTest() {
	close(tb.completed)
}

func (tb *recordingTB) done() bool {
	select {
	case <-tb.completed:
		return true
	default:
		return false
	}
}

func (tb *recordingTB) waitForCompleted(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-tb.completed:
		return nil
	}
}

func (tb *recordingTB) Cleanup(f func()) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.Cleanups = append(tb.Cleanups, f)
}

func (tb *recordingTB) logfLocked(skip int, format string, args ...interface{}) {
	tb.logLocked(skip+1, fmt.Sprint(args...))
}

func (tb *recordingTB) logLocked(skip int, args ...interface{}) {
	callers := getCallers(skip + 1)
	tb.Logs = append(tb.Logs, logEntry{
		callers: callers,
		entry:   fmt.Sprint(args...),
	})
}

func (tb *recordingTB) Error(args ...interface{}) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.logLocked(1 /* skip */, args...)
	tb.failLocked()
}

func (tb *recordingTB) Errorf(format string, args ...interface{}) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.logfLocked(1 /* skip */, format, args...)
	tb.failLocked()
}

func (tb *recordingTB) Fail() {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.failLocked()
}

func (tb *recordingTB) failLocked() {
	if tb.done() {
		panic("Fail in goroutine after test completed")
	}
	tb.failed = true
}

func (tb *recordingTB) FailNow() {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.failNowLocked()
}

func (tb *recordingTB) failNowLocked() {
	tb.failLocked()
	runtime.Goexit()
}

func (tb *recordingTB) Failed() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	return tb.failed
}

func (tb *recordingTB) Fatal(args ...interface{}) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.logLocked(1 /* skip */, args...)
	tb.failNowLocked()
}

func (tb *recordingTB) Fatalf(format string, args ...interface{}) {
	tb.logfLocked(1 /* skip */, format, args...)
	tb.FailNow()
}

func (tb *recordingTB) Helper() {
	// indicates stack should be skipped.
}

func (tb *recordingTB) Log(args ...interface{}) {
	tb.logLocked(1 /* skip */, args...)
}

func (tb *recordingTB) Logf(format string, args ...interface{}) {
	tb.logfLocked(1 /* skip */, format, args...)
}

func (tb *recordingTB) Name() string {
	return "faket-no-name"
}

func (tb *recordingTB) Setenv(key, value string) {
	// Set the environment, but clear it on cleanup

}

func (tb *recordingTB) Skip(args ...interface{}) {
	tb.logLocked(1 /* depth */, args...)
	tb.skipNowLocked()
}

func (tb *recordingTB) SkipNow() {
	tb.skipped = true
	runtime.Goexit()
}

func (tb *recordingTB) skipNowLocked() {
	tb.skipped = true
	runtime.Goexit()
}

func (tb *recordingTB) Skipf(format string, args ...interface{}) {
	tb.logfLocked(1 /* skip */, format, args...)
	tb.SkipNow()
}

func (tb *recordingTB) Skipped() bool {
	return tb.skipped
}

func (tb *recordingTB) TempDir() string {
	return "tmp"
}

func getCallers(skip int) []uintptr {
	depth := 32
	for {
		pc := make([]uintptr, depth)
		n := runtime.Callers(skip, pc)
		if n < len(pc) {
			return pc[:n]
		}
		depth *= 2
	}
}
