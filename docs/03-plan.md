# Plan

> Status: approved (2026-11-04)

Total estimate (MVP, M1-M6): **~336 h** (258 h of tasks + ~30% buffer, lesson time included) → at 20 h/week from 2026-11-09, plus 2 weeks for the holidays → **projected MVP ~2027-03-22**, target 2027-04-15 (~3.5 weeks slack). Post-MVP: M7 +59 h (~2027-05-20), M8 +110 h (~2027-07-01).

Estimates for M2-M8 are bottom-up per task group (number of 15-20 min tasks). They are re-baselined at the M1 retro (M1-T78) using your real pace.

| Milestone | Goal | Demo (what can be run and shown) | Definition of Done | Estimate | Risks | ADRs |
|---|---|---|---|---|---|---|
| M1 | Walking skeleton + spikes | Green CI; `go test ./sim` runs 1000 seeds of a toy protocol twice with identical trace hashes; Porcupine flags a hand-made stale read with a visualization | Repo, Makefile, lint (forbidigo), CI (lint, race, govulncheck, image build); Docker image; spikes 1-3 reported in `02-verification.md`; reading done; plan re-baselined | 25 h → **33 h** (79 tasks) | Toolchain friction; sim design wrong | 001, 002, 007 |
| M2 | Leader election in the simulator | `go test ./chaos -seed=N`: 3 nodes under random partitions elect leaders; election-safety invariant holds; a seed replays with a readable trace | Election + heartbeats in `raft/`; `server/` drives Raft in `sim/`; election-safety invariant every tick; scenario tests; CI 100 seeds | 38 h → **49 h** (~115 tasks) | State-machine style too hard (ADR-001 40 h trigger) | 001, 002 |
| M3 | Replicated KV, linearizability-checked (in memory) | Chaos run with client workloads (retries, sessions) → Porcupine OK; invariants and progress checks pass; history visualization | Log replication, commit, no-op, batch cap; `kv/` with sessions/dedup; pending requests in `server/`; workload + history recording rules; Porcupine per key; progress check; CI 100 seeds < 3 min | 64 h → **83 h** (~195 tasks) | Subtle replication bugs; history recording mistakes | 001, 002, 003, 006 |
| M4 | **Checkpoint A:** persistence + crash faults + mutation suite | Same chaos runs now crash/restart nodes and lose unsynced writes; still green; turning on each `mutation` bug turns runs red | `FS` interface + simulated FS; WAL (ADR-007) with recovery rules + unit/fuzz tests; HardState; persist-before-send; `--bootstrap`; crash faults; mutation suite (≥ 3 bugs caught in ≤ 100 seeds); seed-twice hash check in CI | 50 h → **65 h** (~150 tasks) | WAL recovery bugs; simulated FS too lenient | 002, 007 |
| M5 | Real 3-node cluster | `docker compose up`; `kvctl put/get` via any node; kill the leader, writes resume | Concurrency warm-up; `buf` + protos; protobuf round-trip in sim; gRPC transport; `node/` runtime with `-race` tests; KV gRPC server + health; `kvctl`; Compose; smoke test script | 59 h → **77 h** (~180 tasks) | Concurrency bugs in glue; Docker Desktop Sync latency | 001, 004, 006, 007 |
| M6 | Polish and portfolio (**MVP, v0.1.0**) | README with diagram, demo GIF, Porcupine visualization, bug log; release tag | README (what, why, how it works, how to run, how testing works); diagrams; GIF; bug log with seeds; nightly 2000 seeds + release check; smoke test in CI; understanding check with mentor; tag v0.1.0 | 22 h → **29 h** (~65 tasks) | Scope creep in polish | all |
| M7 | Snapshots + torn-write/bit-flip faults (optional) | Long chaos runs with compaction stay green; InstallSnapshot catches up a lagging node | Snapshots (KV + dedup + applied index), InstallSnapshot, log compaction, WAL segments, 2c faults; tag v0.2.0 | 45 h → **59 h** | Snapshot/log interplay | 001, 002, 007 |
| M8 | Own LSM engine (optional) | Same chaos suite on the LSM; amplification + latency report | Memtable, SSTables, compaction, own WAL, atomic apply batch; benchmarks in README; tag v0.3.0 | 85 h → **110 h** | Largest unknown; may be cut | 007 |

## Cut list (ordered: cut from the top first)
| # | Item | Effect if cut | Type |
|---|---|---|---|
| 1 | M8 LSM | Storage story = WAL only | nice |
| 2 | M7 snapshots + 2c faults | Log grows forever (fine for a demo) | nice |
| 3 | Blog post (M6) | — | nice |
| 4 | Nightly 2000 seeds + release check | Only per-push 100 seeds | nice |
| 5 | Polished GIF → plain asciinema recording | Less shiny README | nice |
| 6 | gRPC health service → TCP healthcheck | — | nice |
| 7 | Smoke test in CI → manual script | Real-cluster regressions caught later | nice |
| 8 | Duplicate-message fault, simulated Sync latency | Fewer fault classes | nice |
| 9 | Docker Compose → 3 local processes via Makefile | Less "backend" look | nice |
| 10 | M5 entirely (gRPC, node, kvctl) | **Checkpoint A becomes the deliverable**; the angle stays intact | last resort |
| — | Raft core, WAL + crash faults, Porcupine + invariants + progress + mutation suite, determinism check, CI, README | Never cut | must |

## M1: Walking skeleton + spikes
Order: setup (T01-T13) → Go warm-up (T14-T32) → reading (T33-T42, can interleave) → spike 1 (T43-T63) → spike 2 (T64-T74) → spike 3 (T75-T77) → retro (T78-T79).

| Task | Title | Label | Estimate | Depends on |
|---|---|---|---|---|
| [M1-T01](tasks/M1-T01.md) | Fix the Go toolchain | Commodity | 15 min | — |
| [M1-T02](tasks/M1-T02.md) | Install Go tools and fix PATH | Commodity | 20 min | M1-T01 |
| [M1-T03](tasks/M1-T03.md) | Create the repo and module | Commodity | 15 min | M1-T02 |
| [M1-T04](tasks/M1-T04.md) | Create the package layout | Commodity | 15 min | M1-T03 |
| [M1-T05](tasks/M1-T05.md) | Makefile with build/test/lint targets | Commodity | 20 min | M1-T04 |
| [M1-T06](tasks/M1-T06.md) | golangci-lint config with determinism bans | Commodity | 20 min | M1-T05 |
| [M1-T07](tasks/M1-T07.md) | First table-driven test | Learn | 20 min | M1-T05 |
| [M1-T08](tasks/M1-T08.md) | GitHub Actions: lint and test | Commodity | 20 min | M1-T06, M1-T07 |
| [M1-T09](tasks/M1-T09.md) | CI: govulncheck and caching | Commodity | 20 min | M1-T08 |
| [M1-T10](tasks/M1-T10.md) | Hello server binary with flags | Commodity | 20 min | M1-T04 |
| [M1-T11](tasks/M1-T11.md) | Multi-stage Dockerfile | Commodity | 20 min | M1-T10 |
| [M1-T12](tasks/M1-T12.md) | CI: build the Docker image | Commodity | 15 min | M1-T09, M1-T11 |
| [M1-T13](tasks/M1-T13.md) | README skeleton and MIT license | Commodity | 15 min | M1-T03 |
| [M1-T14](tasks/M1-T14.md) | Structs and methods: value vs pointer receivers | Learn | 20 min | M1-T07 |
| [M1-T15](tasks/M1-T15.md) | Slices: len, cap, append and aliasing | Learn | 20 min | M1-T14 |
| [M1-T16](tasks/M1-T16.md) | Maps and random iteration order | Learn | 20 min | M1-T15 |
| [M1-T17](tasks/M1-T17.md) | Interfaces: implicit implementation and fakes | Learn | 20 min | M1-T14 |
| [M1-T18](tasks/M1-T18.md) | Interfaces 2: small interfaces and io.Writer | Learn | 20 min | M1-T17 |
| [M1-T19](tasks/M1-T19.md) | Errors: sentinels and wrapping | Learn | 20 min | M1-T07 |
| [M1-T20](tasks/M1-T20.md) | Errors 2: custom error types; panic vs return | Learn | 20 min | M1-T19 |
| [M1-T21](tasks/M1-T21.md) | Subtests and t.Helper | Learn | 20 min | M1-T07 |
| [M1-T22](tasks/M1-T22.md) | t.TempDir and golden files | Learn | 20 min | M1-T21 |
| [M1-T23](tasks/M1-T23.md) | Seeded randomness with math/rand/v2 | Learn | 20 min | M1-T16 |
| [M1-T24](tasks/M1-T24.md) | container/heap priority queue | Learn | 20 min | M1-T18 |
| [M1-T25](tasks/M1-T25.md) | Stable ordering with a tie-breaker | Learn | 20 min | M1-T24 |
| [M1-T26](tasks/M1-T26.md) | encoding/binary little-endian | Learn | 20 min | M1-T15 |
| [M1-T27](tasks/M1-T27.md) | CRC32C checksums | Learn | 20 min | M1-T26 |
| [M1-T28](tasks/M1-T28.md) | Decode a framed record stream | Learn | 20 min | M1-T27 |
| [M1-T29](tasks/M1-T29.md) | Hashing a trace | Learn | 15 min | M1-T18 |
| [M1-T30](tasks/M1-T30.md) | Generics basics | Learn | 20 min | M1-T24 |
| [M1-T31](tasks/M1-T31.md) | Packages: exported names and internal/ | Learn | 15 min | M1-T04 |
| [M1-T32](tasks/M1-T32.md) | Warm-up retro | Learn | 15 min | M1-T14..T31 |
| [M1-T33](tasks/M1-T33.md) | Raft paper: overview (§1-4) and visualization | Deep | 20 min | — |
| [M1-T34](tasks/M1-T34.md) | Raft paper §5.1-5.2: terms and election | Deep | 20 min | M1-T33 |
| [M1-T35](tasks/M1-T35.md) | Raft paper §5.3: log replication | Deep | 20 min | M1-T34 |
| [M1-T36](tasks/M1-T36.md) | Raft paper §5.4: safety and Figure 8 | Deep | 20 min | M1-T35 |
| [M1-T37](tasks/M1-T37.md) | Figure 2 as your own checklist | Deep | 20 min | M1-T36 |
| [M1-T38](tasks/M1-T38.md) | etcd raft README: the Ready loop | Deep | 20 min | M1-T37 |
| [M1-T39](tasks/M1-T39.md) | Dissertation ch. 6: client sessions | Deep | 20 min | M1-T35 |
| [M1-T40](tasks/M1-T40.md) | Dissertation ch. 6: read-only queries | Deep | 20 min | M1-T39 |
| [M1-T41](tasks/M1-T41.md) | PostgreSQL wiki: Fsync Errors | Deep | 20 min | M1-T22 |
| [M1-T42](tasks/M1-T42.md) | Deterministic simulation testing post | Learn | 20 min | M1-T25 |
| [M1-T43](tasks/M1-T43.md) | Simulator design sketch | Learn | 20 min | M1-T42 |
| [M1-T44](tasks/M1-T44.md) | Virtual clock | Learn | 15 min | M1-T43 |
| [M1-T45](tasks/M1-T45.md) | Event queue | Learn | 20 min | M1-T25, M1-T44 |
| [M1-T46](tasks/M1-T46.md) | Scheduler loop | Learn | 20 min | M1-T45 |
| [M1-T47](tasks/M1-T47.md) | RNG plumbing from one seed | Learn | 20 min | M1-T23, M1-T46 |
| [M1-T48](tasks/M1-T48.md) | Simulated Node interface | Learn | 15 min | M1-T46 |
| [M1-T49](tasks/M1-T49.md) | Network: delivery with random delay | Learn | 20 min | M1-T47, M1-T48 |
| [M1-T50](tasks/M1-T50.md) | Network: drop | Learn | 15 min | M1-T49 |
| [M1-T51](tasks/M1-T51.md) | Network: duplicate | Learn | 15 min | M1-T49 |
| [M1-T52](tasks/M1-T52.md) | Network: reordering | Learn | 15 min | M1-T49 |
| [M1-T53](tasks/M1-T53.md) | Network: partitions and heal | Learn | 20 min | M1-T50 |
| [M1-T54](tasks/M1-T54.md) | Random fault schedule | Learn | 20 min | M1-T53 |
| [M1-T55](tasks/M1-T55.md) | Trace recorder | Learn | 20 min | M1-T46 |
| [M1-T56](tasks/M1-T56.md) | Trace hash | Learn | 15 min | M1-T29, M1-T55 |
| [M1-T57](tasks/M1-T57.md) | Toy ping protocol node | Learn | 20 min | M1-T48 |
| [M1-T58](tasks/M1-T58.md) | Run 3 toy nodes | Learn | 20 min | M1-T54, M1-T57 |
| [M1-T59](tasks/M1-T59.md) | Determinism test | Learn | 20 min | M1-T56, M1-T58 |
| [M1-T60](tasks/M1-T60.md) | 1000 seeds and timing | Learn | 20 min | M1-T59 |
| [M1-T61](tasks/M1-T61.md) | Replay a failing seed | Learn | 20 min | M1-T59 |
| [M1-T62](tasks/M1-T62.md) | Break determinism on purpose | Learn | 20 min | M1-T60 |
| [M1-T63](tasks/M1-T63.md) | Spike 1 report | Learn | 15 min | M1-T62 |
| [M1-T64](tasks/M1-T64.md) | Porcupine API overview | Learn | 20 min | M1-T21 |
| [M1-T65](tasks/M1-T65.md) | KV model: input and output types | Learn | 20 min | M1-T64 |
| [M1-T66](tasks/M1-T66.md) | KV model: Step function | Learn | 20 min | M1-T65 |
| [M1-T67](tasks/M1-T67.md) | KV model: per-key partition | Learn | 20 min | M1-T66 |
| [M1-T68](tasks/M1-T68.md) | Linearizable history passes | Learn | 20 min | M1-T67 |
| [M1-T69](tasks/M1-T69.md) | Stale read is caught | Learn | 20 min | M1-T68 |
| [M1-T70](tasks/M1-T70.md) | Unknown-outcome operations | Learn | 20 min | M1-T68 |
| [M1-T71](tasks/M1-T71.md) | Retry spanning one logical operation | Learn | 20 min | M1-T70 |
| [M1-T72](tasks/M1-T72.md) | Visualization output | Learn | 15 min | M1-T69 |
| [M1-T73](tasks/M1-T73.md) | Porcupine timing on 300-op histories | Learn | 20 min | M1-T67 |
| [M1-T74](tasks/M1-T74.md) | Spike 2 report | Learn | 15 min | M1-T71, M1-T73 |
| [M1-T75](tasks/M1-T75.md) | Sync latency program | Commodity | 20 min | M1-T22 |
| [M1-T76](tasks/M1-T76.md) | Sync latency in Docker | Commodity | 20 min | M1-T11, M1-T75 |
| [M1-T77](tasks/M1-T77.md) | Spike 3 report | Commodity | 15 min | M1-T76 |
| [M1-T78](tasks/M1-T78.md) | M1 retro and re-baseline | — | 20 min | M1-T63, M1-T74, M1-T77 |
| [M1-T79](tasks/M1-T79.md) | M1 demo | — | 15 min | M1-T78 |

## M2-M8: task groups
Cards are written when a milestone starts (they go stale otherwise); counts are the basis of the estimates.

| Milestone | Task group | Tasks (≈) | Label |
|---|---|---|---|
| M2 | Raft types and message structs | 8 | Deep |
| M2 | Roles and term rules (become follower/candidate/leader) | 12 | Deep |
| M2 | Tick, randomized election timeout from injected RNG | 8 | Deep |
| M2 | RequestVote send/handle, vote rules, election restriction stub | 16 | Deep |
| M2 | Heartbeats (empty AppendEntries) | 8 | Deep |
| M2 | Ready/Advance for messages + HardState | 10 | Deep |
| M2 | `server/` skeleton driving Raft in `sim/` | 14 | Deep |
| M2 | Election-safety invariant per tick | 5 | Learn |
| M2 | Scenario tests: leader partition, split votes, rejoin with higher term | 16 | Deep |
| M2 | Trace dump of Raft state for replay; CI 100 seeds | 8 | Learn |
| M2 | Debug sessions + retro | 10 | — |
| M3 | Log, AppendEntries consistency check, conflict truncation | 30 | Deep |
| M3 | Commit index, majority match, no-op entry, Figure 8 test | 22 | Deep |
| M3 | Ready entries/committed, apply loop, batch cap | 12 | Deep |
| M3 | `kv/` commands, state machine, sessions/dedup | 24 | Deep |
| M3 | `server/` pending requests, leadership loss → NOT_LEADER | 18 | Deep |
| M3 | Simulated clients: retries, backoff, injectable IDs | 18 | Learn |
| M3 | History recording rules → Porcupine per key | 16 | Learn |
| M3 | Invariants: log matching, state machine safety; progress check | 15 | Learn |
| M3 | CI chaos 100 seeds × 300 ops, timing | 6 | Commodity |
| M3 | Debug sessions + retro | 34 | — |
| M4 | `FS` interface + real impl | 9 | Deep |
| M4 | Simulated FS: unsynced loss, dir fsync, Sync latency | 18 | Learn |
| M4 | WAL header, record encode/decode, CRC | 18 | Deep |
| M4 | WAL append/sync/truncate, replay, torn tail, fail-stop, exit codes | 26 | Deep |
| M4 | HardState + persist-before-send in `server/`; `--bootstrap` | 16 | Deep |
| M4 | Crash/restart faults in sim | 12 | Learn |
| M4 | Mutation suite (3 bugs, build tag) | 12 | Learn |
| M4 | Seed-twice hash check in CI; WAL fuzz test | 12 | Learn |
| M4 | Debug sessions + retro + Checkpoint A demo | 27 | — |
| M5 | Concurrency warm-up: goroutines, channels, select, context, sync, -race | 24 | Learn |
| M5 | `buf` setup, `raft.proto`, `kv.proto`, codegen; protobuf round-trip in sim | 18 | Learn |
| M5 | gRPC peer transport | 18 | Learn |
| M5 | `node/` event loop, ticker + coalescing, shutdown, `-race` tests | 36 | Learn |
| M5 | KV gRPC server: handlers → loop, wakeups, error model, health | 21 | Learn |
| M5 | `kvctl`: session, retry, redirect | 15 | Commodity |
| M5 | Config, Compose (3 nodes, volumes, healthchecks, `make reset`) | 15 | Commodity |
| M5 | Smoke test script | 12 | Commodity |
| M5 | Debug sessions + retro | 21 | — |
| M6 | README sections + diagrams | 15 | Commodity |
| M6 | Demo GIF, Porcupine visualization export, bug log page | 15 | Commodity |
| M6 | Nightly + release-check CI, smoke test in CI | 9 | Commodity |
| M6 | Deep component story (Raft, WAL, simulator) | 12 | Deep |
| M6 | Understanding check with mentor; release v0.1.0 | 8 | — |
| M6 | Optional blog post | 6 | nice |

## Later / out of scope
- Prometheus metrics, TUI cluster view, ReadIndex reads, session expiry, 5-node demo, safe single-node rejoin
