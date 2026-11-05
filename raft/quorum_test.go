package raft

import "testing"

func TestQuorum(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want int
	}{
		{"single node", 1, 1},
		{"two nodes", 2, 2},
		{"three nodes", 3, 2},
		{"four nodes", 4, 3},
		{"five nodes", 5, 3},
		{"six nodes", 6, 4},
		{"large odd cluster", 101, 51},
		{"large even cluster", 100, 51},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Quorum(tt.n); got != tt.want {
				t.Errorf("Quorum(%d) = %d, want %d", tt.n, got, tt.want)
			}
		})
	}
}
