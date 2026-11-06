// Package maporder shows that map iteration order is not guaranteed and how
// to get a stable order by sorting the keys.
package maporder

import (
	"maps"
	"slices"
)

// RawKeys returns the keys of m in map iteration order, which is not
// guaranteed and can differ between calls.
func RawKeys(m map[string]int) []string {
	keys := []string{}
	for key := range m {
		keys = append(keys, key)
	}

	return keys
}

// SortedKeys returns the keys of m in sorted order. Use it wherever the
// order of a loop affects behaviour.
func SortedKeys(m map[string]int) []string {
	return slices.Sorted(maps.Keys(m))
}
