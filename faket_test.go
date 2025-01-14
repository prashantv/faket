package faket

import (
	"testing"

	"github.com/prashantv/faket/internal/want"
)

func TestFakeT_Success(t *testing.T) {
	res := RunTest(func(t testing.TB) {
		t.Log("this", "is", "log", 1)
	})
	want.Equal(t, "Failed", res.Failed(), false)
	want.Equal(t, "Skipped", res.Skipped(), false)

	want.DeepEqual(t, "Logs", res.Logs().Messages(), []string{"this is log 1"})
}

func TestFakeT_FailSkip(t *testing.T) {
	tr := RunTest(func(t testing.TB) {
		t.Error("about to skip")
		t.Skipf("skip %s", t.Name())
	})
	want.Equal(t, "Skipped", tr.Skipped(), false)
	want.Equal(t, "Failed", tr.Failed(), true)
	want.Equal(t, "FailedAndSkipped", tr.FailedAndSkipped(), true)
}

func TestFakeT_Helpers(t *testing.T) {
	tr := RunTest(func(t testing.TB) {
		testHelper1(t)
		testHelper2(t)
		testHelper3(t)
	})
	want.DeepEqual(t, "Helpers", tr.Helpers(), []string{
		"github.com/prashantv/faket.testHelper1",
		"github.com/prashantv/faket.testHelper3",
	})
}

func testHelper1(t testing.TB) { t.Helper() }
func testHelper2(t testing.TB) {}
func testHelper3(t testing.TB) { t.Helper() }
