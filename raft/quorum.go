package raft

// Quorum returns the number of votes that forms a majority in a cluster of n nodes.
func Quorum(n int) int {
	return n/2 + 1
}
