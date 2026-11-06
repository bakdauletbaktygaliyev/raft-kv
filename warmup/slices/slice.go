// Package sliceutil demonstrates slice aliasing and how to avoid it with a copy.
package sliceutil

import "slices"

// Head returns a copy of the first n entries of log. The result does not
// share memory with log, so appending to it or changing it leaves log untouched.
// It panics if n is negative or greater than cap(log). If n is greater than
// len(log) but within cap(log), the result includes entries past the end of log.
func Head(log []int, n int) []int {
	return slices.Clone(log[:n])
}
