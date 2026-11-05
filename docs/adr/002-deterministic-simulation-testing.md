# ADR-002: Deterministic simulation with crash faults; Porcupine + invariants + mutation suite

- Status: accepted (2026-11-04), revised in verification rounds 1-2 (see `../02-verification.md`)
- Date: 2026-11-04
- Component label: Learn (built from scratch; effort counted honestly in the plan)
- Decision type: optimal (the project angle)

## Context
Success criterion: randomized chaos runs, all linearizable, in CI; every failure reproducible. A green check must not be possible for a dead or broken cluster, and the harness itself must be shown to catch bugs (verification B2). Requires ADR-001.

## Alternatives considered
| Option | Pros | Cons | Evidence |
|---|---|---|---|
| 2a. Simulate network only | Simplest | Cannot test persistence ("restart loses no committed write") | [S3] |
| 2b. Network + crash/restart + disk that loses unsynced writes | Catches persistence bugs (double vote, lost commit) | Needs an `FS` interface | [S1] |
| 2c. 2b + torn writes, bit flips | Tests WAL checksums and recovery | Mostly storage recovery | [S1] |
| Real processes + Docker/proxy faults | Realistic | Slow, flaky, not reproducible | — |

## Decision
MVP: 2b. M7: add 2c. One seed drives virtual clock, RNG, network faults (drop, delay, reorder, **duplicate**, partition), disk faults and simulated `Sync` latency (0-20 ms virtual). Each run:
1. **Linearizability:** client history checked by Porcupine against a KV model **partitioned per key** `[source: S4]`. See *History recording* below.
2. **Progress:** after the final heal, a fault-free quiet period of 1000 ticks; within its first 600 ticks (10 × max election timeout) every client completes at least one op. A run with no progress fails.
3. **Invariants every tick:** election safety (≤ 1 leader per term), log matching, state machine safety (no two nodes apply different commands at one index).
4. **Mutation suite (CI):** at least 3 planted bugs behind the `mutation` build tag (no runtime flags in Deep code): no persist-before-send, no dedup, reads from local state. **Caught** = fails in ≥ 1 of the first 100 seeds.
5. **Protobuf round-trip:** every simulated message is encoded and decoded through `raft.proto`, so conversion code is tested.
A failing seed is printed; rerunning it is byte-identical (Go map iteration is randomized, so iterate sorted slices `[source: S6]`). Determinism is enforced by **running each CI seed twice and comparing trace hashes**; `forbidigo` additionally bans `time.Now` and global `math/rand` in `raft/ kv/ wal/ server/ sim/ chaos/`. Simulated clients get IDs from the seeded RNG (ADR-006 injectable generator).

### History recording (Porcupine)
- One **logical operation** spans all its retries: call time = first attempt, return time = first definite answer.
- `NOT_LEADER`, `UNAVAILABLE` or timeout after a request was sent is **never** recorded as "not applied"; the client retries the same seq.
- Ops still pending at the end of the run are recorded as **unknown outcome** (open return time) `[assumption: check by spike 2]`.
- Spike 2 includes a hand-made history for each rule.

**CI budget:** chaos runs **without** `-race` (unit tests and 10 chaos seeds run with `-race`). Per push: 100 seeds × 300 ops < 3 min. By MVP release (M6): nightly 2000 seeds < 30 min; release check 200 seeds × ≥ 1000 ops. All budgets are targets until measured in spikes 1-2 `[assumption: check by spikes 1-2]`.

## Why
- One seed controlling clock, randomness and faults makes runs reproducible `[source: S1]`
- Porcupine is used by etcd robustness tests and Amazon MemoryDB; supports per-key partitioning and visualization `[source: S4]`
- Linearizability checking is NP-complete; partitioning keeps it fast `[source: S4]`
- A checker that is never shown to fail proves little; the mutation suite demonstrates it `[assumption: check by M4 mutation suite]`

## Cost (what we pay)
- Build the simulator ourselves: ~30 h estimated `[assumption: check by M1 spike 1 actuals]`
- `node/` and `transport/` are outside the simulator; covered by keeping them thin (pending requests live in `server/`) and one real-process smoke test in CI (ADR-004)

## Revisit if
- A per-key Porcupine check takes > 1 s at 1000 ops
- A bug class appears that only real processes reproduce

## Sources
| ID | Link | Date checked | Supports which claim |
|---|---|---|---|
| S1 | https://notes.eatonphil.com/2024-08-20-deterministic-simulation-testing.html | 2026-11-04 | DST concept, single seed |
| S2 | https://pkg.go.dev/github.com/arshnah/detsim | 2026-11-04 | Example only, not evidence: seeded network faults in Go (README, unverified) |
| S3 | https://pdos.csail.mit.edu/6.824/labs/lab-raft1.html | 2026-11-04 | Network-fault testing of Raft (network-only model) |
| S4 | Porcupine v1.3.1 README + `model.go` (github.com/anishathalye/porcupine) | 2026-11-04 | Users, partitions, NP-completeness, MIT license |
| S6 | https://go.dev/ref/spec#For_statements | 2026-11-04 | Map iteration order is not specified |

## Uncertainty
- Exact Porcupine API for unknown-outcome ops: verified in spike 2.
- The simulated FS is a model of `fsync` semantics; real-process tests cannot verify durability (killing a process keeps the OS page cache). The simulator is the only durability evidence; this is accepted and stated in the README.
