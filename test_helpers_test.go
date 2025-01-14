package faket

import (
	"testing"

	"github.com/prashantv/faket/internal/want"
)

func TestMustHelpers(t *testing.T) {
	tests := []struct {
		name      string
		fn        func(testing.TB)
		mustPass  bool
		mustFail  bool
		mustPanic bool
	}{
		{
			name:     "passing test",
			fn:       func(testing.TB) {},
			mustPass: true,
		},
		{
			name: "failing test",
			fn: func(t testing.TB) {
				t.Error("failed with Error")
			},
			mustFail: true,
		},
		{
			name: "panic test",
			fn: func(testing.TB) {
				panic("failed with panic")
			},
			mustFail:  true,
			mustPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fnTR := RunTest(tt.fn)

			t.Run("MustPass", func(t *testing.T) {
				tr := RunTest(func(t testing.TB) {
					fnTR.MustPass(t)
				})
				want.Equal(t, "Failed", tr.Failed(), !tt.mustPass)
			})

			t.Run("MustFail", func(t *testing.T) {
				tr := RunTest(func(t testing.TB) {
					fnTR.MustFail(t, "failed")
				})
				want.Equal(t, "Failed", tr.Failed(), !tt.mustFail)
			})

			if tt.mustFail {
				t.Run("MustFail wrong message", func(t *testing.T) {
					tr := RunTest(func(t testing.TB) {
						fnTR.MustFail(t, "unknown")
					})
					want.Equal(t, "Failed", tr.Failed(), true)
					want.Contains(t, "Message", tr.Logs().String(), "missing expected log")
				})
			}

			t.Run("MustPanic", func(t *testing.T) {
				tr := RunTest(func(t testing.TB) {
					fnTR.MustPanic(t, "panic")
				})
				want.Equal(t, "Failed", tr.Failed(), !tt.mustPanic)
			})

			if tt.mustPanic {
				t.Run("MustPanic wrong contains", func(t *testing.T) {
					tr := RunTest(func(t testing.TB) {
						fnTR.MustPanic(t, "unknown")
					})
					want.Equal(t, "Failed", tr.Failed(), true)
					want.Contains(t, "Message", tr.Logs().String(), "panic string doesn't contain")
				})
			}
		})
	}
}
