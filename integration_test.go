package faket_test

import (
	"math"
	"os"
	"path"
	"testing"

	"github.com/prashantv/faket/internal/cmptest"
	"github.com/prashantv/faket/internal/sliceutil"
)

// These integration-style tests are used to compare the real output of running
// a test to the result of `RunTest`.

var goLatest = os.Getenv("GO_NOT_LATEST") != "true"

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
		t.Fatal("fatal") //nolint:revive // skip unreachable skip
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
		onlyLatest bool
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
			onlyLatest:  true, // log caller of panic.go:<line> changes line across versions.
		},
		{
			name:        "mix nesting and skips",
			iterations:  3,
			returnOuter: iter(2),
			skipInner:   iter(1) | iter(2),
			onlyLatest:  true, // log caller of panic.go:<line> changes line across versions.
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.onlyLatest && !goLatest {
				t.Skip("can only run with latest go")
			}

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

func TestCmp_CleanupHelper(t *testing.T) {
	cmptest.Compare(t, func(t testing.TB) {
		helperWithCleanup(t)
	})
}

func helperWithCleanup(t testing.TB) {
	t.Helper()

	t.Cleanup(func() {
		t.Helper()
		t.Log("cleanup func log")
	})
}

func TestCmp_Setenv(t *testing.T) {
	const k = "FAKET_CMP_SETENV_TEST_KEY"

	// k shouldn't be set initially, or tests will fail.
	_, ok := os.LookupEnv(k)
	if ok {
		t.Fatalf("environment key %v cannot be set for test to run", k)
	}

	tests := []struct {
		name  string
		setup func(t testing.TB)
	}{
		{
			name:  "unset initially",
			setup: func(t testing.TB) {},
		},
		{
			name: "set initially",
			setup: func(t testing.TB) {
				if err := os.Setenv(k, "initial"); err != nil {
					t.Fatal("Setenv err", err)
				}
				t.Cleanup(func() {
					if err := os.Unsetenv(k); err != nil {
						t.Fatal("Unsetenv err", err)
					}
				})
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmptest.Compare(t, func(t testing.TB) {
				t.Cleanup(func() {
					v, ok := os.LookupEnv(k)
					t.Logf("LookupEnv in cleanup got %v %v", v, ok)
				})

				t.Setenv(k, "s1")
				t.Log("Getenv s1:", os.Getenv(k))

				t.Setenv(k, "s2")
				t.Log("Getenv s2:", os.Getenv(k))
			})
		})
	}
}

func TestCmp_TmpDir(t *testing.T) {
	cmptest.Compare(t, func(t testing.TB) {
		listFiles := func(msg, dir string) {
			// Note: dir is not included in messages for deterministic test output.
			entries, err := os.ReadDir(dir)
			if os.IsNotExist(err) {
				t.Logf("missing dir %s", msg)
				return
			} else if err != nil {
				t.Fatalf("ReadDir %s failed: %v", msg, err)
			}

			t.Logf("ReadDir %s: %v", msg, sliceutil.Map(entries, os.DirEntry.Name))
		}

		var d1, d2 string

		t.Cleanup(func() {
			listFiles("d1 post-cleanup", d1)
			listFiles("d2 post-cleanup", d2)
		})

		d1 = t.TempDir()
		d2 = t.TempDir()
		listFiles("d1 initial", d1)
		listFiles("d2 initial", d2)

		createFile := func(dir, f string) {
			if err := os.WriteFile(path.Join(dir, f), []byte("dummy"), 0o666); err != nil {
				t.Fatalf("WriteFile %s in %s failed: %v", f, dir, err)
			}
		}

		createFile(d1, "f1")
		createFile(d1, "f2")
		createFile(d2, "f3")
		listFiles("d1 after create", d1)
		listFiles("d2 after create", d2)
	})
}

func TestCmp_Fatalf(t *testing.T) {
	cmptest.Compare(t, func(t testing.TB) {
		t.Error("pre-fail log")
		t.Fatalf("log: %v", "fatal") //nolint:revive // skip unreachable skip
		t.Error("post-fail log")
	})
}

func TestCmp_FailNow(t *testing.T) {
	cmptest.Compare(t, func(t testing.TB) {
		t.Error("pre-fail log")
		t.FailNow() //nolint:revive // skip unreachable skip
		t.Error("post-fail log")
	})
}

func TestCmp_SkipNow(t *testing.T) {
	cmptest.Compare(t, func(t testing.TB) {
		t.Error("pre-skip log")
		t.SkipNow()
		t.Error("post-skip log")
	})
}

func TestCmp_Fail(t *testing.T) {
	cmptest.Compare(t, func(t testing.TB) {
		t.Fail()
	})
}

func TestCmp_Skipf(t *testing.T) {
	cmptest.Compare(t, func(t testing.TB) {
		t.Skipf("skip %s", "test")
	})
}
