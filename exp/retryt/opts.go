package retryt

import (
	"testing"
	"time"

	"github.com/prashantv/faket"
)

// Opts are options to customize retryt.
type Opts struct {
	// Attempts is the maximum number of times the test function will be run.
	// If <= 0, defaults to 10.
	Attempts int

	// Retry is run after a failed run that will be retried.
	// By default, Retry logs the attempt number, and sleeps for
	// <attempt> milliseconds as a backoff.
	Retry func(t testing.TB, attempt int, tr faket.TestResult)

	// Passed is run after a successful run of the function passed to [Test].
	Passed func(t testing.TB, attempt int, tr faket.TestResult)
}

func (o *Opts) setDefaults() {
	if o.Attempts <= 0 {
		o.Attempts = 10
	}
	if o.Retry == nil {
		o.Retry = noop
	}
	if o.Passed == nil {
		o.Passed = noop
	}
}

// RetryBackoff logs `attempt`, and sleeps for `attempt` milliseconds.
//
// It's intended to be used as [Opts].Retry to log and backoff on failure.
func RetryBackoff(t testing.TB, attempt int, tr faket.TestResult) {
	t.Helper()

	sleepFor := time.Duration(attempt) * time.Millisecond
	t.Logf("retryt attempt %d failed, retrying in %v", attempt, sleepFor)
	time.Sleep(sleepFor)
}

func noop(testing.TB, int, faket.TestResult) {}
