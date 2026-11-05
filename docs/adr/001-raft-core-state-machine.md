# ADR-001: Raft core as a pure state machine (Step / Tick / Ready)

- Status: accepted (2026-11-04), revised in verification rounds 1-2 (see `../02-verification.md`)
- Date: 2026-11-04
- Component label: Deep
- Decision type: optimal

## Context
Raft is written from scratch (Deep). The project angle is linearizability checking under injected faults, with thousands of reproducible runs. The core must be testable without real time, real network or real disk.

## Alternatives considered
| Option | Pros | Cons | Evidence |
|---|---|---|---|
| A. Pure state machine: `Tick()`, `Step(msg)`, `Ready()` / `Advance()`; no I/O, no goroutines, no clock | Deterministic; drivable by a simulator; bugs replay from a seed | Fewer tutorials in this style; virtual time is new | [S1], [S2] |
| B. MIT-style: goroutines, mutexes and timers inside Raft | Most tutorials and labs use it | Timing-dependent bugs, hard to reproduce; races in critical code | [S3], [S4] |
| C. B tested only with real processes + Docker/proxy faults | Most realistic | Slow, flaky in CI, barely reproducible | — |

## Decision
We choose A. Inputs: `Tick()` (logical time) and `Step(msg)`. Output: `Ready` with messages, entries + HardState to persist, committed entries to apply. The caller persists before sending (ADR-007), then calls `Advance()`.

Protocol details fixed now:
- A new leader appends a **no-op entry** in its term, so entries from earlier terms can commit `[source: S5]`
- **Timing:** 1 tick = 10 ms (real runtime). Heartbeat = 5 ticks (50 ms). Election timeout = random 30-60 ticks (300-600 ms). Satisfies broadcastTime ≪ electionTimeout `[source: S5]`; failover < 2 s is then expected `[assumption: check by M5 real-cluster smoke test]`
- InstallSnapshot is added in M7 with snapshots (ADR-007)
- AppendEntries carries at most 64 entries or 1 MiB; the simulator enforces the same cap
- Real runtime coalesces missed ticks (at most one `Tick()` per loop iteration) so a slow `Sync` cannot cause a burst of ticks

## Why
- The etcd raft design treats Raft as a deterministic state machine; I/O is the caller's job `[source: S1]`
- Persist-before-send gives a clear contract to test `[source: S2]`
- Subtle Raft bugs appear under long delays and reordering `[source: S3]`, `[source: S4]`
- Stepping through a replayed failing run helps a beginner `[assumption: check by M2 retro]`

## Cost (what we pay)
- Ideas from MIT-style tutorials must be translated, not copied
- Determinism rules (no `time.Now`, `go`, global rand, map-order iteration) in `raft/`, `kv/`, `wal/`, `server/`, enforced by `forbidigo` in CI
- Go concurrency is learned in `node/` and gRPC, not inside Raft

## Revisit if
- After **40 hours** on the Raft election milestone (M2), the election tests (single leader per term under partitions) still do not pass. Fallback: keep explicit ticks and a single goroutine, but restructure the code along the MIT lab's algorithm (handler per RPC type) instead of etcd's `Ready` batching. This keeps byte-identical replay. Putting goroutines inside Raft is **not** an acceptable fallback without a new ADR: it loses reproducible seeds, and the vision's angle must then be restated.

## Sources
| ID | Link | Date checked | Supports which claim |
|---|---|---|---|
| S1 | https://github.com/etcd-io/raft | 2026-11-04 | State-machine design, caller-owned I/O |
| S2 | https://pkg.go.dev/go.etcd.io/raft/v3 | 2026-11-04 | Persist-before-send, Advance |
| S3 | https://pdos.csail.mit.edu/6.824/labs/lab-raft1.html | 2026-11-04 | labrpc delays/drops/reorders to expose bugs |
| S4 | https://minhmannh2001.github.io/2026/04/13/mit-6.5840-distributed-systems-lab-3-part-3c-en.html | 2026-11-04 | Bug only with long reordering (one student report) |
| S5 | https://raft.github.io/raft.pdf | 2026-11-04 | Raft paper: §5.4.2 commit of earlier terms, §5.6 timing, §8 no-op |

## Uncertainty
- S4 is a single student report; S3 (official lab) is the main evidence.
