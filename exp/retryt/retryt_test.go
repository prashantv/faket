package retryt_test

import (
	"testing"

	"github.com/prashantv/faket"
	"github.com/prashantv/faket/exp/retryt"
	"github.com/prashantv/faket/internal/want"
)

func TestN(t *testing.T) {
	const n = 5
	tests := []struct {
		name            string
		innerFn         func(t testing.TB, count int)
		wantSkipped     bool
		wantFailed      bool
		wantCount       int
		containsLogs    []string
		notContainsLogs []string
	}{
		{
			name:      "pass noop",
			innerFn:   func(testing.TB, int) {},
			wantCount: 2, /* faket run + real testing.TB run */
		},
		{
			name: "pass 2nd attempt",
			innerFn: func(t testing.TB, count int) {
				want.Equal(t, "count", count, 2)
			},
			wantFailed: true,  // needs to pass consistently, but fails on final attempt.
			wantCount:  2 + 1, // 2 runs to pass on faket + final attempt,
			containsLogs: []string{
				"retryt attempt 1 failed",
				"count: expected equal",
			},
			notContainsLogs: []string{
				"retryt attempt 2 failed", // 2nd attempt succeeds
				"retryt attempt 3 failed", // No retries after 2nd attempt succeeds
			},
		},
		{
			name: "pass last attempt",
			innerFn: func(t testing.TB, count int) {
				t.Log("count =", count)
				want.Equal(t, "count", count, n)
			},
			wantCount: n,
			containsLogs: []string{
				"retryt attempt 1 failed, retrying in 1ms",
				"retryt attempt 4 failed, retrying in 4ms",
				"count = 5",
			},
			notContainsLogs: []string{
				"retry attempt 5 failed", // last attempt passes
			},
		},
		{
			name: "skip",
			innerFn: func(t testing.TB, count int) {
				t.Skip("skip in testFn")
			},
			wantSkipped: true,
			wantCount:   2, /* faket run + real testing.TB run */
			containsLogs: []string{
				"skip in testFn",
			},
			notContainsLogs: []string{
				"retry attempt 1 failed", // no retries on skip
			},
		},
		{
			name: "fail",
			innerFn: func(t testing.TB, count int) {
				want.Equal(t, "count", 0, count)
			},
			wantFailed: true,
			wantCount:  n,
			containsLogs: []string{
				"retryt attempt 1 failed, retrying in 1ms",
				"retryt attempt 4 failed, retrying in 4ms",
				"got:  0\nwant: 5",
			},
			notContainsLogs: []string{
				"retryt attempt 5 failed", // last attempt is not logged
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var count int
			tr := faket.RunTest(func(t testing.TB) {
				retryt.N(t, n, func(t testing.TB) {
					count++
					tt.innerFn(t, count)
				})
			})
			want.Equal(t, "Skipped", tr.Skipped(), tt.wantSkipped)
			want.Equal(t, "Failed", tr.Failed(), tt.wantFailed)

			gotLogs := tr.Logs().String()
			for _, log := range tt.containsLogs {
				want.Contains(t, "logs", gotLogs, log)
			}
			for _, log := range tt.notContainsLogs {
				want.NotContains(t, "logs", gotLogs, log)
			}
			want.Equal(t, "run count", count, tt.wantCount)
		})
	}
}

func TestRun_Defaults(t *testing.T) {
	var count int
	tr := faket.RunTest(func(t testing.TB) {
		retryt.Opts{}.Run(t, func(t testing.TB) {
			count++
			t.Fatal("fail")
		})
	})
	want.Equal(t, "Failed", tr.Failed(), true)
	want.Equal(t, "run count", count, 10)

	gotLogs := tr.Logs().String()
	want.Contains(t, "logs", gotLogs, "fail")
	want.NotContains(t, "logs", gotLogs, "retryt attempt") // no logs by default
}
