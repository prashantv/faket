# faket [![GoDoc][doc-img]][doc] [![Build Status][ci-img]][ci]

Test helpers have bugs too -- who watches the watchmen?

faket lets you test your test helpers that rely on a `testing.TB`.

## Installation

`go get -u github.com/prashantv/faket`

faket is only supported and tested against the 2 most recent minor versions of Go.

## Quick Start

Let's say we wrote a test helper that checks if a string contains given fragments in order,
```go
// Signature:
func ContainsInOrder(t testing.TB, s string, contains ...string)

// Usage:
strtest.ContainsInOrder(t, "how to test this helper works?", "test", "helper")
```

It's easy to test the success cases, by adding normal tests that call the `ContainsInOrder` with
correct arguments, but how do you ensure this helper fails when it doesn't contain the given values
in order? A normal test will fail, but what you really want is a test that is expected to fail.

faket makes it easy to verify that the test fails, and fails in the expected manner,
```go

res := faket.RunTest(func(t testing.TB) {
	strtest.ContainsInOrder(t, "help test", "test", "helper")
})


assert.True(t, res.Failed(), "expected ContainsInOrder to fail when substrings are not in order")
assert.Contains(t, res.Logs(), `did not find "helper" after "test" in`, "unexpected failure logs")
```

## Development Status: Development

This library is still in development, and is not API stable.

[doc-img]: https://pkg.go.dev/badge/github.com/prashantv/zap
[doc]: https://pkg.go.dev/github.com/prashantv/faket
[ci-img]: https://github.com/prashantv/faket/actions/workflows/go.yml/badge.svg
[ci]: https://github.com/prashantv/faket/actions/workflows/go.yml

