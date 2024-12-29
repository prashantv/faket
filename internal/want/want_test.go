package want_test

import (
	"errors"
	"testing"

	"github.com/prashantv/faket"
	"github.com/prashantv/faket/internal/want"
)

func TestNoErr(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantFail string
	}{
		{
			name: "nil",
			err:  nil,
		},
		{
			name:     "error",
			err:      errors.New("err"),
			wantFail: "expected no error, got: err",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := faket.RunTest(func(t testing.TB) {
				want.NoErr(t, tt.err)
			})
			if tt.wantFail == "" {
				tr.MustPass(t)
			} else {
				tr.MustFail(t, tt.wantFail)
			}
		})
	}
}

func TestEqual(t *testing.T) {
	tests := []struct {
		name     string
		got      string
		want     string
		wantFail string
	}{
		{
			name: "equal",
			got:  "a",
			want: "a",
		},
		{
			name:     "not equal",
			got:      "a",
			want:     "b",
			wantFail: "got==want: expected equal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := faket.RunTest(func(t testing.TB) {
				want.Equal(t, "got==want", tt.got, tt.want)
			})
			if tt.wantFail == "" {
				tr.MustPass(t)
			} else {
				tr.MustFail(t, tt.wantFail)
			}
		})
	}
}

func TestDeepEqual(t *testing.T) {
	tests := []struct {
		name     string
		got      []string
		want     []string
		wantFail string
	}{
		{
			name: "equal",
			got:  []string{"a", "b", "c"},
			want: []string{"a", "b", "c"},
		},
		{
			name:     "not equal",
			got:      []string{"a", "c"},
			want:     []string{"a", "b", "c"},
			wantFail: "got equals want: expected DeepEqual",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := faket.RunTest(func(t testing.TB) {
				want.DeepEqual(t, "got equals want", tt.got, tt.want)
			})
			if tt.wantFail == "" {
				tr.MustPass(t)
			} else {
				tr.MustFail(t, tt.wantFail)
			}
		})
	}
}
