package maporder

import (
	"slices"
	"testing"
)

func TestRawKeysOrderVaries(t *testing.T) {
	m := map[string]int{}
	for _, k := range []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"} {
		m[k] = 1
	}

	first := RawKeys(m)
	for i := 0; i < 200; i++ {
		if got := RawKeys(m); !slices.Equal(got, first) {
			return // found a second, different order
		}
	}
	t.Errorf("RawKeys returned the same order %v in 201 calls, want at least two different orders", first)
}

func TestSortedKeysStable(t *testing.T) {
	m := map[string]int{}
	for _, k := range []string{"e", "d", "c", "b", "a", "j", "i", "h", "g", "f"} {
		m[k] = 1
	}

	want := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	for i := 0; i < 200; i++ {
		if got := SortedKeys(m); !slices.Equal(got, want) {
			t.Fatalf("run %d: SortedKeys(m) = %v, want %v", i, got, want)
		}
	}
}
