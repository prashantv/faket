package faket

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFakeT_Success(t *testing.T) {
	res := RunTest(func(t testing.TB) {
		t.Log("this", "is", "log", 1)
	})
	assert.False(t, res.Failed())
	assert.False(t, res.Skipped())

	want := "this is log 1"
	assert.Equal(t, want, res.Logs())
}
