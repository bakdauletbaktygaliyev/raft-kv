# ADR-007: WAL format, recovery rules and storage roadmap

- Status: accepted (2026-11-04)
- Date: 2026-11-04
- Component label: Deep
- Decision type: optimal for MVP; LSM own WAL (5a) is a **learning / presentation choice**
- Supersedes: ADR-005

## Context
ADR-005 left the MVP recovery rules undefined, and "stop at the first bad record" could drop synced, acknowledged entries. It also had segment files the MVP does not need and bundled too much into one milestone. Round 2 added node identity and a safe reset rule (`../02-verification.md`, B5, O1, O5, R2-3).

## Alternatives considered
| Option | Pros | Cons | Evidence |
|---|---|---|---|
| Library for the Raft log (bbolt, Pebble) | Fast | Contradicts Deep goal | — |
| MVP: single append-only file | Simplest; enough until compaction exists | Grows forever until snapshots (M7) | — |
| MVP: segment files | Ready for compaction | No benefit before M7 | — |
| On any bad record: stop reading, keep what was read | Simple | Silently loses synced entries / votes on mid-log corruption → breaks Raft safety | [S2] |
| Truncate only an incomplete tail; mid-log corruption → fail-stop | Never silently loses synced data | Node with corruption needs manual reset (re-sync from peers) | [S2], [S3] |
| M8 LSM with own WAL (5a) | Standalone, testable engine; reuses `wal/` | Double writes | — |
| M8 LSM without WAL, rely on Raft log (5b) | Fewer writes | Coupled to Raft; subtle recovery | — |

## Decision
**MVP (M4):**
- One file `raft.wal` with a header `[magic "RKVW"][version u16][cluster ID 16 B][node ID u64]`. A node whose WAL header does not match its configured cluster/node ID refuses to start.
- **Empty WAL:** a node refuses to start with no WAL unless started with `--bootstrap` (first start of the whole cluster). Wiping one node's WAL and restarting it would erase its vote and allow double voting; **in the MVP the only allowed reset is the whole cluster** (`make reset`).
- Record: `[len u32][crc32c u32][type u8][payload]`, all integers **little-endian**; CRC covers `len`, `type` and `payload`; `len` ≤ 1 MiB (larger = corrupt). Payloads are protobuf messages from `raft.proto` (Entry, HardState) or a varint index (Truncate). A byte-layout example is written with the code in M4.
- Record types: `Entry`, `HardState` (term, vote, commit; last one wins on replay, as in etcd), `Truncate` (from index N, for log conflicts).
- Persist before send: all records of a `Ready` are written, then one `Sync`, then messages go out (ADR-001).
- `fsync` the directory after creating the file.
- Recovery: read records in order. When a record fails (bad CRC or `len` out of range):
  - if its declared end reaches or passes EOF, or fewer than 9 header bytes remain → **torn tail**: truncate at its start, `fsync` the file, continue;
  - otherwise → **fail-stop**, exit code 3, with a clear error.
  - **Residual risk accepted:** a bit flip in the last synced record looks like a torn tail and is truncated, losing that record (same trade-off as etcd's WAL).
- Any `Sync` / write error (including ENOSPC) → crash the process with exit code 2; never retry a failed `fsync`.
- Compose uses `restart: "no"`, so a fail-stopped node does not crash-loop.

**M7:** snapshots (applied index + KV data + dedup table), InstallSnapshot RPC, log compaction, segment files; torn-write and bit-flip faults in the simulator (ADR-002, 2c).

**M8:** LSM engine (memtable, SSTables, compaction) with its own WAL reusing `wal/` (5a), behind a `Storage` interface. KV data, applied index and dedup table must become durable **atomically** (one LSM batch), or replay double-applies / skips entries.

## Why
- Persist HardState and entries before sending messages `[source: S1]`
- After `fsync` fails, the page cache state is undefined; retrying can falsely report success → crash instead `[source: S2]`, `[assumption: confirm with the PostgreSQL wiki "Fsync Errors" as an M1 reading task]`
- Crash consistency depends on fsync ordering, including the directory `[source: S3]`
- CRC covering the length prevents misframing from a flipped length `[assumption: check by M7 bit-flip tests]`
- 5a is chosen to keep the LSM standalone and reuse `wal/` `[learning choice]`

## Cost (what we pay)
- MVP startup replays the whole log
- Fail-stop on corruption needs a whole-cluster reset in the MVP (documented); safe single-node rejoin is out of scope
- M8 double writes

## Revisit if
- M8 benchmarks show double writes dominate latency → new ADR for 5b

## Sources
| ID | Link | Date checked | Supports which claim |
|---|---|---|---|
| S1 | https://pkg.go.dev/go.etcd.io/raft/v3 | 2026-11-04 | Persist-before-send, HardState incl. commit |
| S2 | https://lwn.net/Articles/667788/ | 2026-11-04 | fsync/crash consistency discussion ("Files are hard") |
| S3 | https://www.usenix.org/conference/osdi14/technical-sessions/presentation/pillai | 2026-11-04 | Pillai et al. OSDI'14: fsync ordering, FS differences |

## Uncertainty
- S2 is a discussion thread, not a primary source; see the assumption tag above.
- **macOS:** Go's `File.Sync` on darwin uses `F_FULLFSYNC` (verified in Go 1.27.3 source, `internal/poll/fd_fsync_darwin.go`). Docker Desktop runs Linux in a VM; durability of named volumes is **not** testable by killing containers (page cache survives). Spike 3 measures `Sync` latency only; durability evidence comes from the simulated FS (ADR-002).
