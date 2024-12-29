# faket [![GoDoc][doc-img]][doc] [![Build Status][ci-img]][ci]

Are your test helpers correct -- test them with `faket`.

## Introduction

Test helper functions play a crucial role in
encapsulating common validation patterns
and keeping your tests DRY (Don't Repeat Yourself).

These helpers often accept a `testing.TB` to report failures
but it can be tricky to test failure scenarios.

As these helpers grow in complexity and become widely used,
verifying their behaviour becomes important--
your tests' correctness depends upon them!

Failure handling by test helpers is particularly critical:
 * Do they indicate all errors, or fail on required checks?
 * Is the output formatted correctly
   and does it include all relevant information?

`faket` addresses these challenges by providing a fake `testing.TB`
so you can validate correctness of your test helpers
and ensure the reported messages are accurate and meaningful.

Let's take an example helper used to verify
if a string contains expected substrings in order:

```go
func StrContainsInOrder(t testing.TB, got string, contains ...string)
```

`faket` can be used to test that the helper behaves as expected:
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
   implementations to ensure correct behaviour
   * E.g., ensuring skips within cleanups work the same way.
 * Panic handling.

There are still some features to implement, see [TODO](TODO.md).

## Installation

`go get -u github.com/prashantv/faket`

faket is only supported and tested against the 2 most recent minor versions of Go.

## Development Status: Development

This library is still in development, and is not API stable.

[doc-img]: https://pkg.go.dev/badge/github.com/prashantv/zap
[doc]: https://pkg.go.dev/github.com/prashantv/faket
[ci-img]: https://github.com/prashantv/faket/actions/workflows/go.yml/badge.svg
[ci]: https://github.com/prashantv/faket/actions/workflows/go.yml
