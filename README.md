# faket [![GoDoc][doc-img]][doc] [![Build Status][ci-img]][ci] [![Coverage][codecov-img]][codecov]


Are your test helpers correct -- be confident with faket.

## Introduction

Test helper functions taking a `testing.TB` are
helpful for encapsulating common validation patterns
and DRY'ing up tests (Don't Repeat Yourself).

These helpers are not immune from bugs and
as they grow in complexity and become widely used.
Verifying their behaviour becomes important--
your tests' correctness depends upon them!

How do you test scenarios where the helper is expected to fail?
Using a real `testing.T` would fail the test,
when you want the opposite:
the test should fail if the helper passes unexpectedly
but pass if the helper fails with the expected failures!

faket lets you validate correctness of your helpers
using a fake `testing.TB`, which provides insight
into the full behaviour of the test function:

 * Are all the expected errors reported?
 * Did the helper `Fatal` the test?
 * Is the output formatted correctly
   with all the relevant information?

### Example

Let's test a helper used to verify
a string contains expected substrings in order:

```go
func StrContainsInOrder(t testing.TB, got string, contains ...string)
```

With faket, both successful and failure scenarios can be tested:
```go
func TestStrContainsInOrder(t *testing.T) {
  t.Run("correct order", func(t *testing.T) {
    faket.RunTest(func(t testing.TB) {
      StrContainsInOrder(t, "test helper function", "test", "helper")
    }).MustPass(t)
  })

  t.Run("incorrect order", func(t *testing.T) {
    faket.RunTest(func(t testing.TB) {
      StrContainsInOrder(t, "test helper function", "helper", "test")
    }).MustFail(t, `failed to find "test" in remaining string " function"`)
  })
}
```

## Features

 * Fully-featured implementation of `testing.TB`
   including `Skip`, `Fatal`, and `Cleanup`.
 * Thoroughly tested against the real `testing.TB`
   to ensure correct behaviour.
 * Panic handling to allow validation of expected panics.

## Installation

`go get -u github.com/prashantv/faket`

faket is only supported and tested against the 2 most recent minor versions of Go.

## Development Status: Alpha

This library is ready for use **in tests**,
but it's _not_ API stable.

[doc-img]: https://pkg.go.dev/badge/github.com/prashantv/zap
[doc]: https://pkg.go.dev/github.com/prashantv/faket
[ci-img]: https://github.com/prashantv/faket/actions/workflows/go.yml/badge.svg
[ci]: https://github.com/prashantv/faket/actions/workflows/go.yml
[codecov-img]: https://codecov.io/github/prashantv/faket/graph/badge.svg?token=RUXXMHOX4Q
[codecov]: https://codecov.io/github/prashantv/faket
