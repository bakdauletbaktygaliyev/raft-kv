package sliceutil

import (
	"slices"
	"testing"
)

func TestHead(t *testing.T) {
	log := []int{1, 2, 3, 4, 5}
	got := Head(log, 3)
	want := []int{1, 2, 3}
	if !slices.Equal(got, want) {
		t.Errorf("Head(log, 3) = %v, want %v", got, want)
	}
}

func TestHeadDoesNotAlias(t *testing.T) {
	log := []int{1, 2, 3, 4, 5}
	got := Head(log, 3)
	_ = append(got, 99)

	want := []int{1, 2, 3, 4, 5}
	if !slices.Equal(log, want) {
		t.Errorf("log = %v, want %v", log, want)
	}
}
