package faket

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestEvent matches https://pkg.go.dev/cmd/test2json#hdr-Output_Format
type TestEvent struct {
	Time    time.Time // encodes as an RFC3339-format string
	Action  string
	Package string
	Test    string
	Elapsed float64 // seconds
	Output  string
}

var testEvents map[string][]TestEvent

func init() {
	resultsJSON, err := ioutil.ReadFile("testdata/results.json")
	if err != nil {
		log.Fatalf("failed to read test results: %v", err)
	}

	testEvents = make(map[string][]TestEvent)
	dec := json.NewDecoder(bytes.NewReader(resultsJSON))
	for dec.More() {
		var ev TestEvent
		if err := dec.Decode(&ev); err != nil {
			log.Fatalf("failed to unmarshal test event: %v", err)
		}

		testEvents[ev.Test] = append(testEvents[ev.Test], ev)
	}
}

func testWrapper(t *testing.T, f func(testing.TB)) {
	if os.Getenv("RUN_ACTUAL_TEST") != "" {
		f(t)
		return
	}

	// Otherwise, compare the real result vs our recorded result.
	res := RunTest(f)
	// compare res to something

	// run    - the test has started running
	// pause  - the test has been paused
	// cont   - the test has continued running
	// pass   - the test passed
	// bench  - the benchmark printed log output but did not fail
	// fail   - the test or benchmark failed
	// output - the test printed output
	// skip   - the test was skipped or the package contained no tests

	var wantOutput []string
	realTestEvents := testEvents[t.Name()]
	var resultVerified bool
	for _, ev := range realTestEvents {
		switch ev.Action {

		case "pass":
			assert.False(t, res.Failed(), "pass test got Failed")
			assert.False(t, res.Skipped(), "pass test got Skipped")
			resultVerified = true
		case "fail":
			assert.True(t, res.Failed(), "expect fail test")
			assert.False(t, res.Skipped(), "fali test got Skipped")
			resultVerified = true
		case "skip":
			assert.False(t, res.Failed(), "skip test got Failed")
			assert.True(t, res.Skipped(), "expect skip test")
			resultVerified = true

		case "output":
			// Only space prefixed lines are t.Log output
			if strings.HasPrefix(ev.Output, "    ") {
				wantOutput = append(wantOutput, ev.Output)
			}
		}
	}

	assert.True(t, resultVerified, "Result comparison missing")

	// compare output
	// assert.Equal(t, wantOutput, res.Logs())
}

func TestSuccess(t *testing.T) {
	testWrapper(t, func(t testing.TB) {
		t.Log("log1")
		t.Log("log2")
	})
}

func TestSkip(t *testing.T) {
	testWrapper(t, func(t testing.TB) {
		t.Log("pre-skip")
		t.Skip("skip")
		t.Log("post-skip")
	})
}

func TestFailure(t *testing.T) {
	testWrapper(t, func(t testing.TB) {
		t.Log("pre-fail log")
		t.Error("error log")
		t.Log("post-fail log")
	})
}

func TestPanic(t *testing.T) {
	testWrapper(t, func(t testing.TB) {
		panic("panic")
	})
}

func TestFailThenSkip(t *testing.T) {
	testWrapper(t, func(t testing.TB) {
		t.Error("error")
		t.Skip("skipped")
	})
}

func TestSkipThenFail(t *testing.T) {
	testWrapper(t, func(t testing.TB) {
		t.Skip("skip")
		t.Error("skipped error")
	})
}

func TestFatal(t *testing.T) {
	testWrapper(t, func(t testing.TB) {
		t.Log("pre-fatal")
		t.Fatal("fatal")
		t.Log("post-fatal")
	})
}

func TestCleanup(t *testing.T) {
	testWrapper(t, func(t testing.TB) {
		t.Log("log 1")
		t.Cleanup(func() {
			t.Log("log in cleanup")
		})
		t.Log("log 2")
	})
}

func TestCleanupError(t *testing.T) {
	testWrapper(t, func(t testing.TB) {
		t.Log("log 1")
		t.Cleanup(func() {
			t.Error("error in cleanup")
		})
		t.Log("log 2")
	})
}

func TestCleanupSkip(t *testing.T) {
	testWrapper(t, func(t testing.TB) {
		t.Log("log 1")
		t.Cleanup(func() {
			t.Log("cleanup 1")
		})
		t.Cleanup(func() {
			t.Skip("skip in cleanup")
			t.Log("log after skip in cleanup")
		})
		t.Cleanup(func() {
			t.Log("cleanup 2")
		})
		t.Log("log 2")
	})
}
