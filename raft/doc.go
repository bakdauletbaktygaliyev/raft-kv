// Package raft implements the Raft consensus algorithm as a pure,
// deterministic state machine. Callers drive it with Tick and Step and
// act on the Ready it returns, then call Advance. It performs no I/O and
// uses no goroutines, clocks or global randomness.
package raft
