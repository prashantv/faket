package faket

import (
	"context"
	"testing"
)

type Result struct {
	res *recordingTB
}

func (r Result) Skipped() bool {
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
func (r Result) FailedAndSkipped() bool {
	return r.res.Failed() && r.res.Skipped()
}

func (r Result) Failed() bool {
	return r.res.Failed()
}

func (r Result) Logs() []string {
	// TODO
	return nil
}

func RunTest(testFn func(t testing.TB)) Result {
	tb := newRecordingTB()

	go func() {
		defer tb.postTest()

		testFn(tb)
	}()

	tb.waitForCompleted(context.Background())
	return Result{tb}
}
