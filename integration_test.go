package faket

import "testing"

// These integration-style tests are used to compare the real output of running
// a test to the result of `RunTest`.

func TestCmp_Success(t *testing.T) {
	compareTest(t, func(t testing.TB) {
		t.Log("log1")
		t.Log("log2")
	})
}

func TestCmp_LogFormatting(t *testing.T) {
	compareTest(t, func(t testing.TB) {
		t.Log("a", 1, "b", 2, "c", "d")
		t.Logf("a: %v b: %v", 1, 2)
	})
}

func TestCmp_Skip(t *testing.T) {
	compareTest(t, func(t testing.TB) {
		t.Log("pre-skip")
		t.Skip("skip")
		t.Log("post-skip")
	})
}

func TestCmp_Failure(t *testing.T) {
	compareTest(t, func(t testing.TB) {
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
	compareTest(t, func(t testing.TB) {
		t.Error("error")
		t.Skip("skipped")
	})
}

func TestCmp_SkipThenFail(t *testing.T) {
	compareTest(t, func(t testing.TB) {
		t.Skip("skip")
		t.Error("skipped error")
	})
}

func TestCmp_Fatal(t *testing.T) {
	compareTest(t, func(t testing.TB) {
		t.Log("pre-fatal")
		t.Fatal("fatal")
		t.Log("post-fatal")
	})
}

func TestCmp_Cleanup(t *testing.T) {
	compareTest(t, func(t testing.TB) {
		t.Log("log 1")
		t.Cleanup(func() {
			t.Log("log in cleanup")
		})
		t.Log("log 2")
	})
}

func TestCmp_CleanupError(t *testing.T) {
	compareTest(t, func(t testing.TB) {
		t.Log("log 1")
		t.Cleanup(func() {
			t.Error("error in cleanup")
		})
		t.Log("log 2")
	})
}

func TestCmp_CleanupSkip(t *testing.T) {
	compareTest(t, func(t testing.TB) {
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
