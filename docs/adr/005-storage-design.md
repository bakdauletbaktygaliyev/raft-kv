# ADR-005: Own WAL for the Raft log; in-memory state machine in MVP; LSM with its own WAL in milestone 2

- Status: superseded by ADR-007 (2026-11-04)
- Date: 2026-11-04
- Component label: Deep
- Decision type: optimal (5a is the safer first engine)

## Context
The storage engine is a Deep component. All reads and writes go through the Raft log (ADR-003), so the log is the durable record of every operation. All persistent code writes through the `FS` interface (ADR-002).

## Alternatives considered
| Option | Pros | Cons | Evidence |
|---|---|---|---|
| Library for the Raft log (bbolt, Pebble) | Fast to finish | Contradicts Deep goal | — |
| MVP: own WAL for Raft log + HardState; KV map in memory, rebuilt by replay | Smallest correct MVP; no second WAL yet | Startup replays the whole log until snapshots exist | [S1] |
| 5a. M2 LSM with its own WAL (reusing `wal/`) | Standalone, testable and benchmarkable engine; simple recovery; code reuse | Each write hits disk twice | — |
| 5b. M2 LSM without WAL, rely on Raft log replay from the last flushed applied index | Fewer writes; closer to optimized systems | Couples LSM to Raft; subtle recovery; LSM cannot be shown alone | — |

## Decision
- `wal/`: segment files; records `[len][crc32c][type][payload]`; `Sync` before messages leave the node (persist-before-send, ADR-001); truncation on log conflict.
- MVP state machine: in-memory map + dedup table, rebuilt by log replay.
- Milestone 2: snapshots (applied index + KV data + dedup table) and log compaction; LSM engine (memtable, SSTables, compaction) with its own WAL behind a `Storage` interface (5a); torn-write and bit-flip faults (ADR-002, 2c). Recovery stops at the first bad record.

## Why
- Caller must persist HardState and entries before sending messages `[source: S1]`
- A standalone LSM with its own crash-recovery suite is a stronger README section `[assumption: check by README review at M2]`
- CRC per record detects torn and corrupted records `[assumption: check by 2c fault tests in M2]`

## Cost (what we pay)
- Write amplification in M2 (Raft WAL + LSM WAL)
- Full log replay on startup until snapshots arrive (fine for test-size data)

## Revisit if
- M2 benchmarks show the double write dominates latency; then consider 5b as an optimization ADR

## Sources
| ID | Link | Date checked | Supports which claim |
|---|---|---|---|
| S1 | https://pkg.go.dev/go.etcd.io/raft/v3 | 2026-11-04 | Persist-before-send rule |

## Uncertainty
- `fsync` on macOS does not guarantee a flush to the device unless `F_FULLFSYNC` is used; whether Go's `File.Sync` on darwin does this must be checked in Verification. Linux in Docker is the reference platform.
