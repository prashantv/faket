package faket_test

import (
	"math"
	"testing"

	"github.com/prashantv/faket/internal/cmptest"
)

// These integration-style tests are used to compare the real output of running
// a test to the result of `RunTest`.

func TestCmp_Success(t *testing.T) {
	cmptest.Compare(t, func(t testing.TB) {
		t.Log("log1")
		t.Log("log2")
	})
}

func TestCmp_LogFormatting(t *testing.T) {
	cmptest.Compare(t, func(t testing.TB) {
		t.Log("a", 1, "b", 2, "c", "d")
		t.Logf("a: %v b: %v", 1, 2)
	})
}

func TestCmp_Skip(t *testing.T) {
	cmptest.Compare(t, func(t testing.TB) {
		t.Log("pre-skip")
		t.Skip("skip")
		t.Log("post-skip")
	})
}

func TestCmp_Failure(t *testing.T) {
	cmptest.Compare(t, func(t testing.TB) {
		t.Log("pre-fail log")
		t.Error("error log")
		t.Log("post-fail log")
	})
}

// TODO(prashant): Panic stops remaining test execution, so move test
// to a separate package.
// func TestCmp_Panic(t *testing.T) {
// 	compareTest(t, func(t testing.TB) {
// 		panic("panic")
// 	})
// }

func TestCmp_FailThenSkipCmp(t *testing.T) {
	cmptest.Compare(t, func(t testing.TB) {
		t.Error("error")
		t.Skip("skipped")
	})
}

func TestCmp_SkipThenFail(t *testing.T) {
	cmptest.Compare(t, func(t testing.TB) {
		t.Skip("skip")
		t.Error("skipped error")
	})
}

func TestCmp_Fatal(t *testing.T) {
	cmptest.Compare(t, func(t testing.TB) {
		t.Log("pre-fatal")
		//nolint:revive // skip unreachable skip
		t.Fatal("fatal")
		t.Log("post-fatal")
	})
}

func TestCmp_Cleanup(t *testing.T) {
	cmptest.Compare(t, func(t testing.TB) {
		t.Log("log 1")
		t.Cleanup(func() {
			t.Log("log in cleanup")
		})
		t.Log("log 2")
	})
}

func TestCmp_CleanupError(t *testing.T) {
	cmptest.Compare(t, func(t testing.TB) {
		t.Log("log 1")
		t.Cleanup(func() {
			t.Error("error in cleanup")
		})
		t.Log("log 2")
	})
}

func TestCmp_CleanupSkip(t *testing.T) {
	cmptest.Compare(t, func(t testing.TB) {
		t.Log("log 1")
		t.Cleanup(func() {
			t.Log("cleanup 1")
		})
		t.Cleanup(func() {
			t.Log("cleanup 2")
			t.Skip("skip in cleanup")
			t.Log("log after skip in cleanup")
		})
		t.Cleanup(func() {
			t.Log("cleanup 3")
		})
		t.Log("log 2")
	})
}

func TestCmp_Helper(t *testing.T) {
	cmptest.Compare(t, func(t testing.TB) {
		t.Log("call log directly")
		log(t)

		t.Log("log in helper")
		logHelper(t, 1, func() {})
		logHelper(t, 3, func() {})

		t.Log("log helper then log")
		logHelper(t, 3, func() { log(t) })
	})
}

func logHelper(t testing.TB, n int, last func()) {
	t.Helper()

	t.Log("logHelper", n)

	if n == 0 {
		last()
		return
	}
	logHelper(t, n-1, last)
}

func log(t testing.TB) {
	t.Log("log")
}

func TestCmp_NestedCleanup(t *testing.T) {
	// bit mask helpers
	const never = 0
	const all = math.MaxUint32
	iter := func(i int) uint32 {
		return (1 << i)
	}
	isSet := func(mask uint32, i int) bool {
		return mask&iter(i) > 0
	}

	tests := []struct {
		name       string
		iterations int

		// bit masks of which iterations to do the action on.
		returnOuter uint32
		skipInner   uint32
	}{
		{
			name:        "always nest without skip",
			iterations:  2,
			returnOuter: never,
			skipInner:   never,
		},
		{
			name:        "always nest with skips",
			iterations:  2,
			returnOuter: never,
			skipInner:   all,
		},
		{
			name:        "mix nesting and skips",
			iterations:  3,
			returnOuter: iter(2),
			skipInner:   iter(1) | iter(2),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmptest.Compare(t, func(t testing.TB) {
				t.Log("log 1")
				defer t.Log("defer 1")

				for i := 1; i <= tt.iterations; i++ {
					t.Cleanup(func() {
						defer t.Log("defer cleanup", i)
						t.Log("cleanup", i)

						if isSet(tt.returnOuter, i) {
							return
						}

						t.Cleanup(func() {
							defer t.Log("defer nested cleanup", i)
							t.Log("nested cleanup", i)

							if isSet(tt.skipInner, i) {
								t.Skip("skip in nested cleanup", i)
							}
						})
					})
				}
			})
		})
	}
}
