// Package kv is the replicated key-value state machine. It applies
// committed commands in log order and uses client sessions to detect
// and drop duplicate retries. It is deterministic: no I/O, goroutines,
// clocks or global randomness.
package kv
