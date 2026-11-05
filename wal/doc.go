// Package wal implements the write-ahead log that persists Raft entries
// and hard state. Records are CRC-framed in a single file; recovery
// discards a torn tail and stops on any other corruption. All file
// access goes through an FS interface so the simulator can inject
// disk faults.
package wal
