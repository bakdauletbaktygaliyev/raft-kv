// Package chaos runs fault-injection tests on top of sim. It drives
// client workloads, records histories, checks them for linearizability
// with Porcupine, checks safety invariants and progress, and runs
// the mutation suite. A failure prints the seed needed to replay it.
package chaos
