package faket

import "testing"

func TestFakeT_Success(t *testing.T) {
	res := RunTest(func(t testing.TB) {
		t.Log("this", "is", "log", 1)
	})
	wantEqual(t, "Failed", res.Failed(), false)
	wantEqual(t, "Skipped", res.Skipped(), false)

	wantEqual(t, "Logs", res.Logs(), "this is log 1")
}
