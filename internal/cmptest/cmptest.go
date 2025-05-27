// Package cmptest is used to compare faket against go test.
package cmptest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/prashantv/faket"
	"github.com/prashantv/faket/internal/want"
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

var (
	runActual  = os.Getenv("RUN_ACTUAL_TEST") != ""
	testEvents = mustReadTestEvents("cmp_test_results.json")
)

func mustReadTestEvents(file string) map[string][]TestEvent {
	if runActual {
		return nil
	}

	resultsJSON, err := os.ReadFile(filepath.Join("testdata", file))
	if err != nil {
		panic(fmt.Errorf("failed to read test results: %v", err))
	}

	events := make(map[string][]TestEvent)
	dec := json.NewDecoder(bytes.NewReader(resultsJSON))
	for dec.More() {
		var ev TestEvent
		if err := dec.Decode(&ev); err != nil {
			log.Fatalf("failed to unmarshal test event: %v", err)
		}

		events[ev.Test] = append(events[ev.Test], ev)
	}

	return events
}

// Compare compares the result of running the given test function
// using `faket.RunTest` against `go test`.
func Compare(t *testing.T, f func(testing.TB)) {
	CompareOpts(t, Opts{}, f)
}

// Opts are options for comparing faket to a real test run.
type Opts struct {
	WantPanic  bool
	LogReplace func(string) string
}

// CompareOpts is the same as Compare, but supports options for customizing comparisons.
func CompareOpts(t *testing.T, opts Opts, f func(t testing.TB)) {
	// Verify the name, since the Makefile uses this to only run Cmp tests when generating the test data.
	if !strings.HasPrefix(t.Name(), "TestCmp_") {
		t.Fatalf("test %v is a Cmp test, and should be named TestCmp_*", t.Name())
	}

	if runActual {
		f(t)
		return
	}

	res := faket.RunTest(f)

	var wantOutput strings.Builder
	realTestEvents := testEvents[t.Name()]
	var resultEvent bool
	for _, ev := range realTestEvents {
		switch ev.Action {
		case "run", "pause":
			// Events that don't need any processing
			// Started running a specific test.
		case "pass":
			want.Equal(t, "pass test Failed", res.Failed(), false)
			want.Equal(t, "pass test Skipped", res.Skipped(), false)
			resultEvent = true
		case "fail":
			want.Equal(t, "fail test Failed", res.Failed(), true)
			want.Equal(t, "fail test Skipped", res.Skipped(), false)
			resultEvent = true
		case "skip":
			want.Equal(t, "skip test Failed", res.Failed(), false)
			want.Equal(t, "skip test Skipped", res.Skipped(), true)
			resultEvent = true
		case "output":
			trimmed, ok := strings.CutPrefix(ev.Output, "    ")
			if !ok {
				// Only space prefixed lines are t.Log output
				continue
			}

			if strings.HasPrefix(trimmed, "--- ") {
				// Lines starting with '    ---' are subtest pass/skip/fail lines, skip them.
				continue
			}

			wantOutput.WriteString(trimmed)
		default:
			t.Fatal("unknown action", ev.Action)
		}
	}

	logs := res.Logs()
	if res.Panicked() {
		// Drop the "panic:" line
		last := len(logs) - 1
		if strings.HasPrefix(logs[last].Message, "panic:") {
			logs = logs[:last]
		}
	}

	wantLogs := wantOutput.String()
	gotLogs := logs.String()
	if opts.LogReplace != nil {
		wantLogs = opts.LogReplace(wantLogs)
		gotLogs = opts.LogReplace(gotLogs)
	}

	want.Equal(t, "result event", resultEvent, true)
	want.Equal(t, "log output", gotLogs, wantLogs)
	want.Equal(t, "panicked", res.Panicked(), opts.WantPanic)
}
