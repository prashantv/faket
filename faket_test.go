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
