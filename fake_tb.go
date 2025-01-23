package faket

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"testing"

	"github.com/prashantv/faket/internal/sliceutil"
)

var _ testing.TB = (*fakeTB)(nil)

type fakeTB struct {
	// embedded to include testing.TB.private, which we can't implement.
	// Since this is an interface, unimplemented methods will panic.
	testing.TB

	mu sync.Mutex // protects all of the below fields.

	cleanups []cleanup
	helpers  map[uintptr]struct{}
	logs     []logEntry

	completed chan struct{}
	failed    bool
	skipped   bool
	panicked  bool

	// only set during a cleanup
	cleanupRoot  string
	curCleanupPC []uintptr

	// panic metadata
	recovered      any
	recoverCallers []uintptr
}

type logEntry struct {
	callers        []uintptr // callers[0] is the tb function that logged
	cleanupCallers []uintptr // for logs within a cleanup function
	entry          string
}

type cleanup struct {
	fn      func()
	callers []uintptr
}

// RunTest runs the given test using a fake [testing.TB] and returns
// the result of running the test.
func RunTest(testFn func(t testing.TB)) TestResult {
	tb := newFakeTB()

	go func() {
		defer close(tb.completed)
		defer tb.checkPanic()
		defer tb.runCleanups()

		testFn(tb)
	}()

	<-tb.completed
	return TestResult{tb}
}

func newFakeTB() *fakeTB {
	return &fakeTB{
		completed: make(chan struct{}),
		helpers:   make(map[uintptr]struct{}),
	}
}

func (tb *fakeTB) Name() string {
	return "faket-no-name"
}

// Cleaup and post-test methods.

func (tb *fakeTB) Cleanup(f func()) {
	callerPCs := getCallers(skipSelf)

	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.cleanups = append(tb.cleanups, cleanup{
		callers: callerPCs,
		fn:      f,
	})
}

func (tb *fakeTB) checkPanic() {
	if r := recover(); r != nil {
		tb.mu.Lock()
		defer tb.mu.Unlock()

		tb.panicked = true
		tb.recovered = r
		tb.recoverCallers = getCallers(skipSelf)
		tb.logfLocked(tb.recoverCallers, "panic: %v", r)
	}
}

func (tb *fakeTB) runCleanups() {
	// Set cleanupRoot so log callers can use cleanup's callers.
	if self := getCaller(withSelf); self != 0 {
		f := pcToFunction(self)
		func() {
			tb.mu.Lock()
			defer tb.mu.Unlock()

			tb.cleanupRoot = f
		}()
	}

	// If defer runs with !finished, then a cleanup must have panicked
	// (which could be a Skip/Fatal). Continue running remaining cleanups.
	var finished bool
	defer func() {
		if !finished {
			tb.runCleanups()
		}
	}()

	// Run cleanups in last-first order, similar to defers.
	// Don't iterate by index, as the slice can grow (cleanups can add cleanups).
	for {
		c, ok := func() (cleanup, bool) {
			tb.mu.Lock()
			defer tb.mu.Unlock()

			if len(tb.cleanups) == 0 {
				return cleanup{}, false
			}

			last := len(tb.cleanups) - 1
			c := tb.cleanups[last]
			tb.cleanups = tb.cleanups[:last]
			return c, true
		}()
		if !ok {
			finished = true
			break
		}

		// Set the caller of cleanup for logs with all `t.Helper()` frames.
		func() {
			tb.mu.Lock()
			defer tb.mu.Unlock()

			tb.curCleanupPC = c.callers
		}()

		c.fn()
	}
}

// Logging methods

func (tb *fakeTB) Error(args ...interface{}) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.logLocked(getCallers(withSelf), args...)
	tb.failLocked()
}

func (tb *fakeTB) Errorf(format string, args ...interface{}) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.logfLocked(getCallers(withSelf), format, args...)
	tb.failLocked()
}

func (tb *fakeTB) Fatal(args ...interface{}) {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.logLocked(getCallers(withSelf), args...)
	tb.failNowLocked()
}

func (tb *fakeTB) Fatalf(format string, args ...interface{}) {
	tb.logfLocked(getCallers(withSelf), format, args...)
	tb.FailNow()
}

func (tb *fakeTB) Log(args ...interface{}) {
	tb.logLocked(getCallers(withSelf), args...)
}

func (tb *fakeTB) Logf(format string, args ...interface{}) {
	tb.logfLocked(getCallers(withSelf), format, args...)
}

func (tb *fakeTB) Skip(args ...interface{}) {
	tb.logLocked(getCallers(withSelf), args...)
	tb.skipNowLocked()
}

func (tb *fakeTB) Skipf(format string, args ...interface{}) {
	tb.logfLocked(getCallers(withSelf), format, args...)
	tb.SkipNow()
}

func (tb *fakeTB) logfLocked(callers []uintptr, format string, args ...interface{}) {
	formatted := fmt.Sprintf(format, args...)
	tb.logs = append(tb.logs, logEntry{
		callers: callers,
		entry:   formatted,
	})
}

func (tb *fakeTB) logLocked(callers []uintptr, args ...interface{}) {
	// Log args are formatted using Sprintln in the testing package
	// but we drop the trailing newline as we store an array of lines.
	formatted := fmt.Sprintln(args...)
	formatted = strings.TrimSuffix(formatted, "\n")

	tb.logs = append(tb.logs, logEntry{
		callers:        callers,
		cleanupCallers: tb.curCleanupPC,
		entry:          formatted,
	})
}

// Fail-related methods.

func (tb *fakeTB) Fail() {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	tb.failLocked()
}

func (tb *fakeTB) Failed() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	return tb.failed || tb.panicked
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

func (tb *fakeTB) failLocked() {
	if tb.done() {
		panic("Fail in goroutine after test completed")
	}
	tb.failed = true
}

func (tb *fakeTB) done() bool {
	select {
	case <-tb.completed:
		return true
	default:
		return false
	}
}

// Skip-related methods.

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

// Helper tracking methods.

func (tb *fakeTB) Helper() {
	callerPC := getCaller(skipSelf)
	if callerPC == 0 {
		// no callers, ignore.
		// Note: real testing.TB would panic here, but we avoid panics in faket.
		return
	}

	tb.mu.Lock()
	defer tb.mu.Unlock()

	if _, ok := tb.helpers[callerPC]; !ok {
		tb.helpers[callerPC] = struct{}{}
	}
}

func (tb *fakeTB) helperFuncs() []string {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	funcs := make([]string, 0, len(tb.helpers))
	for pc := range tb.helpers {
		if fn := pcToFunction(pc); fn != "" {
			funcs = append(funcs, fn)
		}
	}
	return funcs
}

// Convert internal logEntry (using PCs) to exported Log (no PCs).
func (tb *fakeTB) toLog(e logEntry) Log {
	l := Log{
		Message: e.entry,
	}

	skipSet := sliceutil.ToSet(tb.helperFuncs())
	// When a defer is triggered by a panic, it's added to the trace
	// but panic is not shown as a log caller.
	skipSet["runtime.gopanic"] = struct{}{}

	frames := runtime.CallersFrames(e.callers)

	f, _ := frames.Next()
	if f == (runtime.Frame{}) {
		return l
	}

	// First frame is the tb caller.
	l.TBFunc = f.Function

	skip := true
	for skip {
		f, _ = frames.Next()
		if f == (runtime.Frame{}) {
			return l
		}

		// If we hit the cleanup root, then use the callers of the t.Cleanup.
		if f.Function == tb.cleanupRoot {
			frames = runtime.CallersFrames(e.cleanupCallers)
			continue
		}

		_, skip = skipSet[f.Function]
	}

	l.CallerFile = f.File
	l.CallerLine = f.Line
	l.CallerFunc = f.Function
	return l
}

// Helpers which aren't core to testing.TB

func (tb *fakeTB) Setenv(key, value string) {
	prevVal, prevSet := os.LookupEnv(key)

	if err := os.Setenv(key, value); err != nil {
		// Match the error from stdlib.
		tb.Fatalf("cannot set environment variable: %v", err)
	}

	tb.Cleanup(func() {
		var err error
		if !prevSet {
			err = os.Unsetenv(key)
		} else {
			err = os.Setenv(key, prevVal)
		}
		if err != nil {
			// This error is ignored by stdlib, but let's report it.
			tb.Fatalf("cannot revert environment variable: %v", err)
		}
	})
}

func (tb *fakeTB) TempDir() string {
	pattern := strings.Map(func(r rune) rune {
		// Similar to stdlib, but with more restricted allowed set.
		const allowed = " -_"
		if '0' <= r && r <= '9' ||
			'a' <= r && r <= 'z' ||
			'A' <= r && r <= 'Z' {
			return r
		}
		if strings.ContainsRune(allowed, r) {
			return r
		}
		return -1
	}, tb.Name())

	d, err := os.MkdirTemp("", pattern)
	if err != nil {
		tb.Fatalf("TempDir: %v", err)
	}

	tb.Cleanup(func() {
		if err := os.RemoveAll(d); err != nil {
			tb.Errorf("TempDir RemoveAll cleanup: %v", err)
		}
	})

	return d
}
