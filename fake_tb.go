package faket

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"testing"
)

type expectedPanic string

const (
	panicFatal expectedPanic = "fatal"
	panicSkip  expectedPanic = "skip"
)

var _ testing.TB = (*fakeTB)(nil)

type fakeTB struct {
	// embedded to include testing.TB.private, which we can't implement.
	// Since this is an interface, unimplemented methods will panic.
	testing.TB

	mu sync.Mutex // protects all of the below fields.

	cleanups []func()
	Logs     []logEntry

	completed chan struct{}
	failed    bool
	skipped   bool
}

type logEntry struct {
	callers []uintptr // callers[0] is the tb function that logged
	entry   string
}

// RunTest runs the given test using a fake [testing.TB] and returns
// the result of running the test.
func RunTest(testFn func(t testing.TB)) TestResult {
	tb := newRecordingTB()

	go func() {
		defer tb.postTest()

		testFn(tb)
	}()

	tb.waitForCompleted(context.Background())
	return TestResult{tb}
}

func newRecordingTB() *fakeTB {
	return &fakeTB{
		completed: make(chan struct{}),
	}
}

func (tb *fakeTB) postTest() {
	defer close(tb.completed)

	defer func() {
		if r := recover(); r != nil {
			panic(r)
		}
	}()

	for _, f := range tb.cleanups {
		defer f()
	}

	// TODO(prashant): Handle nested Cleanups
}

func (tb *fakeTB) done() bool {
	select {
	case <-tb.completed:
		return true
	default:
		return false
	}
}

func (tb *fakeTB) waitForCompleted(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-tb.completed:
		return nil
	}
}

func (tb *fakeTB) Cleanup(f func()) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.cleanups = append(tb.cleanups, f)
}

func (tb *fakeTB) logfLocked(callers []uintptr, format string, args ...interface{}) {
	formatted := fmt.Sprintf(format, args...)
	tb.Logs = append(tb.Logs, logEntry{
		callers: callers,
		entry:   formatted,
	})
}

func (tb *fakeTB) logLocked(callers []uintptr, args ...interface{}) {
	// Log args are formatted using Sprintln in the testing package
	// but we drop the trailing newline as we store an array of lines.
	formatted := fmt.Sprintln(args...)
	formatted = strings.TrimSuffix(formatted, "\n")

	tb.Logs = append(tb.Logs, logEntry{
		callers: callers,
		entry:   formatted,
	})
}

func (tb *fakeTB) Error(args ...interface{}) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.logLocked(getCallers(), args...)
	tb.failLocked()
}

func (tb *fakeTB) Errorf(format string, args ...interface{}) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.logfLocked(getCallers(), format, args...)
	tb.failLocked()
}

func (tb *fakeTB) Fail() {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.failLocked()
}

func (tb *fakeTB) failLocked() {
	if tb.done() {
		panic("Fail in goroutine after test completed")
	}
	tb.failed = true
}

func (tb *fakeTB) FailNow() {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.failNowLocked()
}

func (tb *fakeTB) failNowLocked() {
	tb.failLocked()
	runtime.Goexit()
}

func (tb *fakeTB) Failed() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	return tb.failed
}

func (tb *fakeTB) Fatal(args ...interface{}) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.logLocked(getCallers(), args...)
	tb.failNowLocked()
}

func (tb *fakeTB) Fatalf(format string, args ...interface{}) {
	tb.logfLocked(getCallers(), format, args...)
	tb.FailNow()
}

func (tb *fakeTB) Helper() {
	// TODO(prashant): Implement Helper, this should result in the helper function frame being skipped
	// in any caller file:name resolution.
}

func (tb *fakeTB) Log(args ...interface{}) {
	tb.logLocked(getCallers(), args...)
}

func (tb *fakeTB) Logf(format string, args ...interface{}) {
	tb.logfLocked(getCallers(), format, args...)
}

func (tb *fakeTB) Name() string {
	return "faket-no-name"
}

func (tb *fakeTB) Setenv(key, value string) {
	// Set the environment, but clear it on cleanup

}

func (tb *fakeTB) Skip(args ...interface{}) {
	tb.logLocked(getCallers(), args...)
	tb.skipNowLocked()
}

func (tb *fakeTB) Skipf(format string, args ...interface{}) {
	tb.logfLocked(getCallers(), format, args...)
	tb.SkipNow()
}

func (tb *fakeTB) SkipNow() {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.skipNowLocked()
}

func (tb *fakeTB) skipNowLocked() {
	tb.skipped = true
	runtime.Goexit()
}

func (tb *fakeTB) Skipped() bool {
	return tb.skipped
}

func (tb *fakeTB) TempDir() string {
	return "tmp"
}

func getCallers() []uintptr {
	depth := 32
	for {
		pc := make([]uintptr, depth)
		// runtime.Callers returns itself, so skip that, and this function.
		n := runtime.Callers(2, pc)
		if n < len(pc) {
			return pc[:n]
		}
		depth *= 2
	}
}
