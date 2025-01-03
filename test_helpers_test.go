package faket

import (
	"testing"

	"github.com/prashantv/faket/internal/want"
)

func TestMustPass(t *testing.T) {
	passTR := RunTest(func(testing.TB) {})
	failTR := RunTest(func(t testing.TB) {
		t.Error("failed")
	})

	t.Run("MustPass on passing test", func(t *testing.T) {
		passTR.MustPass(t)
	})

	t.Run("MustPass on failing test", func(t *testing.T) {
		tr := RunTest(func(t testing.TB) {
			failTR.MustPass(t)
		})
		want.Equal(t, "Failed", tr.Failed(), true)
	})

	t.Run("MustFail on passing test", func(t *testing.T) {
		tr := RunTest(func(t testing.TB) {
			passTR.MustFail(t, "failed")
		})
		want.Equal(t, "Failed", tr.Failed(), true)
	})

	t.Run("MustFail on failed test", func(t *testing.T) {
		failTR.MustFail(t, "failed")
	})

	t.Run("MustFail with wrong message", func(t *testing.T) {
		tr := RunTest(func(t testing.TB) {
			failTR.MustFail(t, "incorrect message")
		})
		want.Equal(t, "Failed", tr.Failed(), true)
	})
}
