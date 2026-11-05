// Package node runs a server in a real process. It owns the real
// ticker, feeds gRPC requests and peer messages into a single event
// loop, wakes waiting client requests, and shuts down via context.
// It is the only place where goroutines and wall-clock time are used.
package node
