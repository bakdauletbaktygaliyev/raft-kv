// Package server glues Raft, the KV state machine and the WAL together.
// It handles one event at a time, persists state before sending
// messages, and tracks pending client requests, including what happens
// to them when leadership is lost. The simulator and the real node
// runtime both drive it; it uses no goroutines or real time.
package server
