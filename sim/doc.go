// Package sim is a deterministic simulator for a raft-kv cluster.
// It provides a virtual clock, a seeded RNG, and a simulated network
// and disk that can delay, drop, duplicate and reorder messages, and
// crash and restart nodes. Nothing here may use real time, goroutines
// or global randomness: the same seed must replay the same run.
package sim
