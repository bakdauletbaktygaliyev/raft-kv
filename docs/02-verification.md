# Verification

> Round 2, 2026-11-04. Approved 2026-11-04. Result: **no blocking issues remain**. Round 1: 5 blocking fixed. Round 2: 4 blocking fixed, 1 moved to Planning as a required output. Round 3 fresh-eyes skipped (user agreed): round 2 refined decisions without reversing any.

## Checks
| # | Check | Result | Notes |
|---|---|---|---|
| 1 | ADR form | pass (after fixes) | All ADRs have ≥ 2 alternatives, cost, revisit-if, sources. ADR-005 superseded by ADR-007 |
| 2 | Claim tags | pass (after fixes) | Untagged claims tagged; weak sources replaced (Go spec, Raft paper, Porcupine README) |
| 3 | Checklist completeness | pass | All 10 categories; added timing, limits, disk errors, shutdown, versioning |
| 4 | Traceability | pass | Table below; no orphans |
| 5 | Real-world assumptions | pass with fixes | Versions/licenses verified; time budget re-baselined (~460 h MVP incl. buffer); toolchain problems found |
| 6 | Over-engineering | pass (after fixes) | Single-file WAL in MVP, 3 nodes, M2/M3 split |
| 7 | Spikes | pass | 4 spikes with pass criteria (`01-architecture.md`) |
| 8 | Fresh-eyes review | 2 rounds done | R1: 5 blocking, 26 minor. R2: 5 blocking, 24 minor. Triaged below |

## Real-world checks (round 1)
| Item | Finding | Source |
|---|---|---|
| Go | 1.27.3 installed (Homebrew) and released; second old install `/usr/local/go` 1.24.12 on PATH | `go version` |
| `File.Sync` on macOS | Uses `F_FULLFSYNC` | `internal/poll/fd_fsync_darwin.go` in Go 1.27.3 |
| Porcupine | v1.3.1 (2026-09-21), MIT, partition support, used by etcd + MemoryDB | module proxy, README, `model.go` |
| grpc-go | v1.85.0, Apache-2.0, needs Go ≥ 1.26 | module proxy, LICENSE, go.mod |
| protobuf-go | v1.36.12, BSD-3 | module proxy, LICENSE |
| `buf` / `protoc` | `buf` 1.75.0 works; `protoc` 33.4 fails (missing abseil dylib) | local run |
| Plugins / tools | `protoc-gen-go(-grpc)` in `~/go/bin` (not on PATH); `golangci-lint`, `govulncheck` missing | `which` |
| stdlib `uuid` | Present in Go 1.27 | `go doc uuid` |
| Time budget | ~355 h + 30% ≈ 460 h ≈ 23 weeks at 20 h → MVP ~2027-04-15 | estimate, re-baseline after spikes |

## Traceability
| Requirement (vision) | Decision | Component | Task / milestone |
|---|---|---|---|
| 3-node cluster via Compose, `kvctl` through any node | ADR-004, 006 | `deploy/`, `node/`, `cmd/kvctl` | M5 |
| Failover < 2 s | ADR-001 (timing) | `raft/`, `node/` | M2 (election), M5 (smoke test); spike 3 (M1-T75..T77) |
| Restart loses no acknowledged write | ADR-007, ADR-002 (2b) | `wal/`, `server/` | M4; reading M1-T41 |
| 200 chaos runs linearizable + invariants + progress | ADR-002, 003, 006 | `sim/`, `chaos/`, `kv/` | spikes M1-T43..T74; M2, M3; release check M6 |
| Mutation suite ≥ 3 bugs caught | ADR-002 | `chaos/` | M4 |
| CI budgets | ADR-002, checklist | `.github/` | M1-T08, T09, T12; M3; M4; M6 |
| README, GIF, bug log | checklist | docs | M1-T13; M6 |
| Snapshots + 2c faults | ADR-007, 002 | `wal/`, `raft/`, `sim/` | M7 |
| LSM | ADR-007 | `lsm/` | M8 |

Gaps: none. Every decision maps to a requirement or a quality rule.

## Issues found
| # | Issue | Source | Severity | Resolution |
|---|---|---|---|---|
| B1 | Schedule did not fit scope and level | fresh eyes, check 5 | blocking | Fixed: re-baselined MVP ~2027-04-15 with 30% buffer; 3 nodes; single-file WAL; M2/M3 split; cut list in Planning |
| B2 | Green check possible for broken cluster; harness not validated; unknown ops undefined | fresh eyes | blocking | Fixed in ADR-002: progress check, invariants, mutation suite, unknown-outcome ops, per-key model |
| B3 | Real-cluster code not exercised by simulator | fresh eyes | blocking | Fixed: pending requests in `server/`, protobuf round-trip in sim, smoke test in CI (ADR-002, 004), timing in ADR-001 |
| B4 | Client/session contract undefined | fresh eyes | blocking | Fixed: new ADR-006 |
| B5 | WAL recovery missing / unsafe | fresh eyes | blocking | Fixed: new ADR-007 supersedes ADR-005 |
| M1 | gRPC/HTTP vs gRPC | fresh eyes | minor | Fixed: gRPC only |
| M2 | 3-5 vs 5 nodes | fresh eyes | minor | Fixed: 3 nodes |
| M3 | Metrics goal vs deferred | fresh eyes | minor | Fixed: removed from MVP goals |
| M4 | Concurrency warm-up placed before code that needs it | fresh eyes | minor | Fixed: Go basics first, concurrency warm-up before `node/` |
| M5 | Determinism scope inconsistent, unenforced | fresh eyes | minor | Fixed: aligned; `forbidigo` in CI |
| M6 | Simulator effectively Deep | fresh eyes | minor | Partly: stays Learn (keep Deep at 2 areas), effort counted honestly |
| M7 | Approved docs with checkable open facts | fresh eyes | minor | Fixed: facts resolved in this round |
| M8 | "Proven" overclaims | fresh eyes | minor | Fixed: "linearizability-checked under injected faults" |
| M9 | Mirror URL, unsourced Porcupine users | fresh eyes | minor | Fixed: canonical URL (assumed), Porcupine README cited |
| M10 | Weak/mismatched sources in ADR-002 | fresh eyes | minor | Fixed: Go spec, MIT lab; detsim marked unverified |
| M11 | ADR-001 thin evidence, vague revisit | fresh eyes | minor | Fixed: Raft paper cited; 40 h trigger + fallback |
| M12 | ADR-004 popularity sources | fresh eyes | minor | Fixed: rests on learning choice; module proxy as source |
| M13 | ADR-005 "optimal" on assumptions | fresh eyes | minor | Fixed: ADR-007 marks 5a as learning/presentation choice |
| M14 | fsync spike aimed at wrong question | fresh eyes | minor | Fixed: macOS verified; spike re-scoped to Docker Desktop volumes |
| M15 | No risk register | fresh eyes | minor | Fixed: added |
| M16 | No-op entry, timing, InstallSnapshot missing | fresh eyes | minor | Fixed: ADR-001, 007 |
| M17 | M2/M3 storage atomicity | fresh eyes | minor | Fixed: ADR-007 |
| M18 | CI budget undefined | fresh eyes | minor | Fixed: ADR-002 |
| M19 | Limits and format versioning | fresh eyes | minor | Fixed: ADR-006, 007, checklist |
| M20 | Operations items | fresh eyes | minor | Fixed: checklist (healthchecks, shutdown, volumes, 127.0.0.1) |
| M21 | "≥ 3 real bugs" not under control | fresh eyes | minor | Fixed: bug log + mutation suite |
| O1 | Segments in MVP WAL | fresh eyes | minor | Fixed: single file in MVP |
| O2 | LSM own WAL | fresh eyes | minor | Declined: user chose 5a knowingly; moved to M3, off MVP path |
| O3 | Map vs LSM benchmark meaningless | fresh eyes | minor | Fixed: LSM amplification + latency percentiles |
| O4 | GIF, nightly, govulncheck, layers in MVP | fresh eyes | minor | Partly: 3 nodes accepted; GIF/nightly/govulncheck kept (cheap, high value); `server/` needed for B3 |
| O5 | M2 bundles four features | fresh eyes | minor | Fixed: M2 snapshots + 2c, M3 LSM |
| T1 | Two Go installs; `~/go/bin` not on PATH; `protoc` broken; missing lint tools | check 5 | minor | Spike 4 (setup task in M1) |

## Issues found: round 2
| # | Issue | Severity | Resolution |
|---|---|---|---|
| R2-1 | No bottom-up estimate; cut list missing | blocking | **Moved to Planning as required output**: per-component estimate, ordered cut list (must/nice), M1a checkpoint (Raft + sim + Porcupine, no gRPC/Docker) |
| R2-2 | Spike 3 cannot detect fsync bugs (kill keeps page cache) | blocking | Fixed: spike 3 measures Sync latency only; durability evidence = simulated FS, stated in ADR-002/007 and vision |
| R2-3 | Single-node reset erases vote → double voting | blocking | Fixed in ADR-007: cluster/node ID in WAL header, `--bootstrap`, whole-cluster reset only; non-goal added |
| R2-4 | History recording for retries / NOT_LEADER undefined | blocking | Fixed: "History recording" section in ADR-002; hand-made histories in spike 2 |
| R2-5 | Fallback in ADR-001 loses reproducible seeds | blocking | Fixed: fallback keeps ticks + single goroutine; goroutines in Raft need a new ADR and a restated angle |
| R2-6 | `forbidigo` cannot detect `go`/map range/select | minor | Fixed: run each seed twice, compare trace hashes |
| R2-7 | Client UUIDs break reproducibility | minor | Fixed: injectable ID generator; rules extend to `sim/`, `chaos/` |
| R2-8 | Progress bound too tight, N undefined | minor | Fixed: quiet period 1000 ticks, ≥ 1 op per client within 600 ticks |
| R2-9 | "Caught" undefined; flags in Deep code | minor | Fixed: ≥ 1 of first 100 seeds; `mutation` build tag |
| R2-10 | Prior art (MIT lab, etcd) not named; angle overstated | minor | Fixed: vision names prior art; angle = deterministic replay incl. disk faults + mutation suite |
| R2-11 | Torn-tail rule ambiguous; can drop a synced record | minor | Fixed: precise rule, fsync after truncation, residual risk documented |
| R2-12 | AppendEntries batch unbounded | minor | Fixed: ≤ 64 entries / 1 MiB, enforced in sim |
| R2-13 | Client deadline/backoff undecided; 2 s flaky in CI | minor | Fixed: 500 ms deadline, jittered backoff; CI < 5 s, 2 s locally |
| R2-14 | Missing-key semantics | minor | Fixed in ADR-006 |
| R2-15 | WAL encoding underspecified | minor | Fixed: little-endian, protobuf payloads; byte layout example with code in M1 |
| R2-16 | macOS fsync note dropped from ADR-007 | minor | Fixed: carried into ADR-007 Uncertainty (closed with evidence) |
| R2-17 | Blocking Sync stalls ticks | minor | Fixed: simulated Sync latency; missed ticks coalesced; spike 3 reports p99 |
| R2-18 | No message duplication fault | minor | Fixed: added |
| R2-19 | Simulator estimate/label understated | minor | Partly: re-estimate bottom-up in Planning; label stays Learn (user decision) |
| R2-20 | CI budget unmeasured; `-race` unclear | minor | Fixed: chaos without race; budgets are targets until spikes 1-2 measure them |
| R2-21 | Headline promises storage engine not in MVP | minor | Fixed: "Raft and a WAL from scratch; LSM planned"; criteria split by milestone |
| R2-22 | "Real backend service" vs deferred observability | minor | Fixed: gRPC health service in MVP; goal wording softened |
| R2-23 | Concurrency goal vs thin `node/` | minor | Fixed: concurrent pieces named; `-race` unit tests |
| R2-24 | Bind address ambiguous | minor | Fixed: listen on Compose network, publish on 127.0.0.1 |
| R2-25 | Healthcheck / restart policy undefined | minor | Fixed: gRPC health; `restart: "no"`; exit codes 2/3 |
| R2-26 | ReadIndex trigger cannot fire | minor | Fixed: trigger dropped |
| R2-27 | ADRs reference verification items without link | minor | Fixed: links to `02-verification.md` |
| R2-28 | Key semantics rest on unread primary sources | minor | Fixed: reading dissertation ch. 6 + PostgreSQL "Fsync Errors" is the first M1 task |
| R2-29 | Over-engineering for first milestone | minor | Partly: smoke test, nightly, release check move to "by MVP release"; protobuf round-trip and `govulncheck` kept (cheap, R1 fix); must/nice tags in Planning |

## Deep components: understanding criteria
| Component | User must be able to explain | Study sources |
|---|---|---|
| Raft core | Why a candidate needs a majority; why logs match by (index, term); why a leader may not commit earlier-term entries by counting (Figure 8) and how the no-op fixes it; what must be persisted before replying and why; how randomized timeouts avoid split votes | [Raft paper](https://raft.github.io/raft.pdf); [Ongaro dissertation](https://github.com/ongardie/dissertation); [etcd-io/raft](https://github.com/etcd-io/raft) source; [eliben/raft](https://github.com/eliben/raft) + blog series; [MIT 6.5840 Raft lab](https://pdos.csail.mit.edu/6.824/labs/lab-raft1.html) |
| KV semantics (sessions) | What linearizability means and why a retried write breaks it without dedup; why reads through the log are linearizable and local reads are not; why dedup state must be in snapshots | Dissertation ch. 6; [Porcupine post](https://anishathalye.com/testing-distributed-systems-for-linearizability); [ReadIndex review](https://emptysqua.re/blog/review-rethink-linearizability-constraints-raft/) |
| WAL | What `fsync` does and does not guarantee; why directory fsync is needed; torn tail vs mid-log corruption and why they are handled differently; why a failed fsync must not be retried | [Pillai et al., OSDI'14](https://www.usenix.org/conference/osdi14/technical-sessions/presentation/pillai); ["Files are hard" discussion](https://lwn.net/Articles/667788/); etcd `wal` package source |
| LSM (M3) | Write path (WAL → memtable → SSTable); read path and bloom filters; compaction strategies and write/read amplification | [mini-lsm](https://skyzh.github.io/mini-lsm/) (Rust, concepts transfer); LevelDB `doc/impl.md`; *Database Internals* (Petrov), ch. 7 |

## Unknowns
| Unknown | Why it can't be resolved now | How / when we'll find out |
|---|---|---|
| Real effort for a Go beginner | No actuals yet | Re-baseline after spikes and after M1 first 40 h |
| Porcupine API for unknown-outcome ops and retries | Needs code | Spike 2 |
| `Sync` latency on Docker Desktop volumes | Needs experiment | Spike 3 |
| Real `fsync` durability on Docker Desktop | Not testable by killing containers | Accepted: simulated FS is the evidence (ADR-002/007) |
| Real simulator effort and CI time per seed | Needs code | Spikes 1-2 record hours and timings |
| Whether the state-machine Raft style suits the user | Needs experience | ADR-001 revisit trigger (40 h) |
| SisyphusDB canonical URL | Found only via a mirror | Check on GitHub when writing README |
| "Failed fsync must not be retried" primary source | Not re-checked | Read PostgreSQL wiki "Fsync Errors" in M1 (ADR-007) |
