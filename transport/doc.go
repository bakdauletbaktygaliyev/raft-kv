// Package transport carries Raft messages between nodes and serves the
// client API over gRPC. It only converts and delivers messages; all
// consensus logic lives in raft and server.
package transport
