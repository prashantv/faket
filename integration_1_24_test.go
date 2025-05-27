//go:build go1.24

package faket_test

import (
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/prashantv/faket/internal/cmptest"
	"github.com/prashantv/faket/internal/want"
)

func TestCmp_Chdir(t *testing.T) {
	t.Run("to non-existent", func(t *testing.T) {
		// ChDir failures are reported by testing.go in stdlib
		// and similarly fake_tb.go with faket.
		opts := cmptest.Opts{
			LogReplace: func(s string) string {
				s = regexp.MustCompile("testing.go:[0-9]+").ReplaceAllString(s, "testing.go:Chdir")
				s = regexp.MustCompile("fake_tb.go:[0-9]+").ReplaceAllString(s, "testing.go:Chdir")
				return s
			},
		}
		cmptest.CompareOpts(t, opts, func(t testing.TB) {
			t.Chdir("./does/not/exist")
			t.Log("skipped (due to fatal in ChDir)")
		})
	})

	mustGetWd := func() string {
		wd, err := os.Getwd()
		want.NoErr(t, err)
		return wd
	}

	initialWD := mustGetWd()
	initialPWD := os.Getenv("PWD")

	stripWD := func(s string) string {
		return strings.ReplaceAll(s, initialWD, "/path/to/faket")
	}

	t.Run("success", func(t *testing.T) {
		cmptest.Compare(t, func(t testing.TB) {
			t.Chdir("./testdata")

			wd, err := os.Getwd()
			want.NoErr(t, err)
			t.Log("wd is", stripWD(wd))
			t.Log("PWD is", stripWD(os.Getenv("PWD")))
		})
	})

	postWD := mustGetWd()
	postPWD := os.Getenv("PWD")
	want.Equal(t, "Getwd", postWD, initialWD)
	want.Equal(t, "PWD", postPWD, initialPWD)
}

func TestCmp_Context(t *testing.T) {
	cmptest.Compare(t, func(t testing.TB) {
		ctx := t.Context()
		t.Log("no ctx.Err initially", ctx.Err())

		deadline, ok := ctx.Deadline()
		t.Log("ctx has no deadline", ok, deadline)

		t.Cleanup(func() {
			t.Log("but err in cleanup", ctx.Err())
		})
	})
}
