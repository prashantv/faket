package faket

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

var testEvents = mustReadTestEvents("cmp_test_results.json")

func mustReadTestEvents(file string) map[string][]TestEvent {
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

func compareTest(t *testing.T, f func(testing.TB)) {
	// Verify the name, since the Makefile uses this to only run Cmp tests when generating the test data.
	if !strings.HasPrefix(t.Name(), "TestCmp_") {
		t.Fatalf("test %v is a Cmp test, and should be named TestCmp_*", t.Name())
	}

	if os.Getenv("RUN_ACTUAL_TEST") != "" {
		f(t)
		return
	}

	res := RunTest(f)

	var wantOutput []string
	realTestEvents := testEvents[t.Name()]
	var resultEvent bool
	for _, ev := range realTestEvents {
		switch ev.Action {
		case "run", "pause":
			// Events that don't need any processing
			// Started running a specific test.
		case "pass":
			assert.False(t, res.Failed(), "pass test got Failed")
			assert.False(t, res.Skipped(), "pass test got Skipped")
			resultEvent = true
		case "fail":
			assert.True(t, res.Failed(), "expect Failed for failing test")
			assert.False(t, res.Skipped(), "unexpected Skipped for failing test")
			resultEvent = true
		case "skip":
			assert.False(t, res.Failed(), "skip test got Failed")
			assert.True(t, res.Skipped(), "expect skip test")
			resultEvent = true
		case "output":
			// Only space prefixed lines are t.Log output
			if trimmed, ok := strings.CutPrefix(ev.Output, "    "); ok {
				wantOutput = append(wantOutput, strings.TrimSuffix(trimmed, "\n"))
			}
		default:
			t.Fatal("unknown action", ev.Action)

		}
	}

	assert.True(t, resultEvent, "Result event missing")
	assert.Equal(t, wantOutput, res.testingLogOutput())
}
