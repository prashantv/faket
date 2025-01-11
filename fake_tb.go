package faket

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
	"testing"
)

var _ testing.TB = (*fakeTB)(nil)

type fakeTB struct {
	// embedded to include testing.TB.private, which we can't implement.
	// Since this is an interface, unimplemented methods will panic.
	testing.TB

	mu sync.Mutex // protects all of the below fields.

	cleanups []func()
	helpers  map[uintptr]struct{}
	Logs     []logEntry

	completed chan struct{}
	failed    bool
	skipped   bool
	panicked  bool

	// panic metadata
	recovered      any
	recoverCallers []uintptr
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
		defer tb.checkPanic()
		defer tb.postTest()

		testFn(tb)
	}()

	<-tb.completed
	return TestResult{tb}
}

func newRecordingTB() *fakeTB {
	return &fakeTB{
		completed: make(chan struct{}),
		helpers:   make(map[uintptr]struct{}),
	}
}

func (tb *fakeTB) checkPanic() {
	if r := recover(); r != nil {
		tb.panicked = true
		tb.recovered = r
		tb.recoverCallers = getCallers()
	}
}

func (tb *fakeTB) postTest() {
	defer close(tb.completed)

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

	return tb.failed || tb.panicked
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
	tb.mu.Lock()
	defer tb.mu.Unlock()

	const skip = 2 // runtime.Callers + Helper
	var pc [1]uintptr
	n := runtime.Callers(skip, pc[:])
	if n == 0 {
		// no callers, ignore.
		// Note: real testing.TB would panic here, but we avoid panics in faket.
		return
	}

	if _, ok := tb.helpers[pc[0]]; !ok {
		tb.helpers[pc[0]] = struct{}{}
	}
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
