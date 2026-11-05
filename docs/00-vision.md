# Vision: raft-kv

> Status: approved (2026-11-04), revised in verification rounds 1-2 (2026-11-04)

## Idea
A distributed, strongly consistent key-value store in Go with **Raft and a write-ahead log written from scratch** (an own LSM engine is planned for M8). What sets it apart: **every chaos run is a deterministic simulation replayable from one seed, including disk faults**, and is checked by Porcupine, Raft invariants and progress checks; a mutation suite shows the checks catch planted bugs. For the author it is a deep dive into distributed systems and a portfolio piece for junior Go backend roles.

## Goals
| Goal | Concretely |
|---|---|
| Portfolio | Deterministic, seed-replayable chaos runs with disk faults, checked for linearizability, plus a mutation suite, in CI. Backend-service basics: gRPC API with health service, Docker Compose, CI (metrics later) |
| Learning | Go (basics → concurrency, `context`), gRPC/protobuf, testing distributed systems (fault injection, Porcupine), GitHub Actions |
| Deep dive | Raft (election, replication, persistence, later snapshots) and storage (WAL, later LSM tree), both without libraries |

## Portfolio angle
"Raft + LSM in Go" is a common portfolio topic. Checking linearizability with Porcupine is **not** new by itself: the MIT 6.5840 KV lab tests do it under partitions and crashes, and many public lab solutions inherit those tests; etcd's robustness tests do it too. Our difference is narrower and stated precisely: **one seed replays a whole run byte-for-byte, including simulated disk faults (lost unsynced writes), and a mutation suite demonstrates that the checks catch planted bugs.** Comparison below is a sample, not a survey.

| Similar project | Link | What it does | Our difference |
|---|---|---|---|
| DistroKV | https://github.com/ThaRealJozef/DistroKV | Custom Raft + LSM + gRPC | Linearizability checking under injected faults |
| SisyphusDB | https://github.com/srija-pixel/SisyphusDB (canonical URL assumed; found via a mirror) | Raft KV, Prometheus/Grafana, chaos tests, EKS (per its README, unverified) | Seed-replayable runs with disk faults; mutation suite |
| MIT 6.5840 KV lab | https://pdos.csail.mit.edu/6.824/labs/lab-raft1.html | Course tests: Porcupine under partitions/crashes (real goroutine timing) | Deterministic replay, disk-fault model, mutation suite |
| golt | https://github.com/aunchagaonkar/golt | Raft, snapshots, WAL, Docker Compose | Same, plus our own LSM (M8) |
| eatonphil/goraft | https://github.com/eatonphil/goraft | Minimal teaching Raft + KV | Reference only |
| eliben/raft | https://github.com/eliben/raft | Instructional Raft in Go | Reference only |
| Porcupine | https://github.com/anishathalye/porcupine | Linearizability checker; used by etcd robustness tests and Amazon MemoryDB (per its README) | We use it as a tool |

## MVP
A **3-node** cluster with `docker compose up`:
- Raft from scratch: election (with no-op entry), log replication, persistence
- Single-file WAL with recovery rules (ADR-007); in-memory KV map rebuilt by replay
- **gRPC** client API (`Get` / `Put` / `Delete`) with sessions (ADR-006) + CLI client
- Deterministic chaos harness: Porcupine + invariants + progress checks + mutation suite, in CI
- One real-process smoke test in CI

30-second demo: start the cluster, write keys, kill the leader, see a new leader and no data lost, show the Porcupine visualization of a chaos run and a planted bug being caught.

## Success criteria
- [ ] `docker compose up` starts a 3-node cluster; `kvctl` can `put`/`get` through any node
- [ ] After the leader is killed, a new leader is elected and writes succeed again in < 2 s
- [ ] A full cluster restart (process restart) loses no acknowledged write. Note: this checks replay logic, not `fsync`; durability is evidenced by the simulated disk
- [ ] 200 seeded chaos runs (partitions, crashes, drops/delays; ≥ 1000 ops each) are all linearizable, pass invariant and progress checks
- [ ] Every planted bug in the mutation suite (≥ 3) is caught (fails in ≥ 1 of the first 100 seeds)
- [ ] Each CI seed run twice produces the same trace hash
- [ ] CI on every push: lint, `go test -race`, 100 chaos seeds < 3 min. By MVP release: real-process smoke test (failover < 5 s in CI; < 2 s measured locally) and nightly 2000 seeds
- [ ] README: architecture diagram, demo GIF, bug log ("bugs our tests caught", with seeds)

**M7 (optional for "finished"):**
- [ ] Snapshots + compaction pass the chaos suite with torn-write/bit-flip faults

**M8 (optional for "finished"):**
- [ ] LSM engine passes the same suite; README reports its write/read amplification and latency percentiles

## Non-goals
- No sharding / multi-Raft
- No dynamic membership changes (fixed node set)
- No transactions or multi-key operations
- No auth / TLS (nodes listen on the Compose network; host ports published on 127.0.0.1 only)
- No safe single-node reset/rejoin (whole-cluster reset only)
- Not for production use

## Later / out of scope
- M7: snapshots, InstallSnapshot, log compaction, segment files, torn-write/bit-flip faults
- M8: LSM engine (memtable, SSTables, compaction, own WAL) behind `Storage`
- Prometheus metrics, TUI cluster view, ReadIndex reads, session expiry, 5-node demo

## Constraints
| Item | Value |
|---|---|
| Time per week | 20+ h |
| Target finish | No hard date. Soft targets (incl. 30% buffer, re-baselined after spikes): **MVP (M1-M6) ~2027-04-15**, Checkpoint A (end of M4) ~2027-02-10, M7 ~2027-05-20, M8 ~2027-07-01; see `03-plan.md` |
| Current level in stack | Go: beginner. Built a REST API (`event-booking`, SQLite) following a course; no real use of goroutines/channels/`sync`/`context`. Java Spring Boot projects. Docker: has experience. gRPC: none |

## Components (labels)
| Component | Label | Why |
|---|---|---|
| Raft core + KV semantics | **Deep** | The heart of the project |
| Storage (WAL, later LSM) | **Deep** | Second chosen deep dive |
| Simulator + chaos harness | **Learn** | Built from scratch, but a tool; effort counted in the plan |
| gRPC transport and API | **Learn** | Standard in Go backend roles |
| Go basics warm-up (start), concurrency warm-up (before `node/`) | **Learn** | Prerequisites, placed where they are needed |
| Docker Compose, CI, CLI | **Commodity** | Proven tools |
