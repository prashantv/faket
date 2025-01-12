// Package sliceutil slice utilities that extend the stdlib slices package.
package sliceutil

// Map maps every element in `xs to a new slice using `fn`.
func Map[X, Y any](xs []X, fn func(X) Y) []Y {
	if xs == nil {
		return nil
	}

	ys := make([]Y, len(xs))
	for i := range xs {
		ys[i] = fn(xs[i])
	}
	return ys
}

// ToSet converts a slice to a map with the elements as keys.
func ToSet[X comparable](xs []X) map[X]struct{} {
	set := make(map[X]struct{})
	for _, x := range xs {
		set[x] = struct{}{}
	}
	return set
}
