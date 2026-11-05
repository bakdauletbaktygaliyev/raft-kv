# Knowledge map

Status: unknown → seen → understand → can apply. Starting statuses come from the vision (Go beginner via a course; Docker experience; no concurrency or gRPC in real code). Topics group the concepts listed on task cards.

| Topic | Status | Where needed | Evidence (task / answer / review) | Explained in |
|---|---|---|---|---|
| Go toolchain, modules, `go install`, PATH | seen | M1-T01 | event-booking (course) | |
| Packages, exported names, `internal/`, doc comments | seen | M1-T04 | event-booking | |
| Make, `go test` flags (`-race`, `-run`, custom flags) | unknown | M1-T05 | | |
| golangci-lint, forbidigo | unknown | M1-T06 | | |
| `testing`: table-driven tests, subtests, `t.Helper`, `t.TempDir`, golden files | seen | M1-T07 | event-booking has API tests | |
| GitHub Actions, CI caching, govulncheck | unknown | M1-T08 | | |
| `flag` package, env config | unknown | M1-T10 | | |
| Docker: multi-stage builds, named volumes | understand | M1-T11 | has Docker experience | |
| Structs, methods, pointer vs value receivers | seen | M1-T14 | event-booking | |
| Slices: len/cap, append, aliasing | seen | M1-T15 | | |
| Maps and random iteration order; `slices.Sorted(maps.Keys(...))` | seen | M1-T16 | | |
| Interfaces, small interfaces, fakes, dependency injection, `io.Writer`/`io.Reader` | seen | M1-T17 | | |
| Errors: sentinels, `%w`, `errors.Is/As`, custom types; panic vs return; exit codes | seen | M1-T19 | | |
| `math/rand/v2`, seeded PCG, passing the RNG explicitly | unknown | M1-T23 | | |
| `container/heap`, priority queues, tie-breakers | unknown | M1-T24 | | |
| `encoding/binary`, endianness | unknown | M1-T26 | | |
| CRC32C, hashing (`hash.Hash`, fnv/sha256) | unknown | M1-T27 | | |
| Framing records in a byte stream | unknown | M1-T28 | | |
| Generics: type parameters, constraints | unknown | M1-T30 | | |
| Replicated state machine, consensus | unknown | M1-T33 | | |
| Raft: terms, roles, RequestVote, election timeout | unknown | M1-T34 | | |
| Raft: AppendEntries, log matching, commit index | unknown | M1-T35 | | |
| Raft: election restriction, Figure 8, commitment rules, no-op entry | unknown | M1-T36 | | |
| Raft Figure 2 rules | unknown | M1-T37 | | |
| Ready/Advance, persist-before-send | unknown | M1-T38 | | |
| Linearizability, sessions, dedup | unknown | M1-T39 | | |
| Stale reads, ReadIndex, leases | unknown | M1-T40 | | |
| `fsync` semantics, page cache, fsync error handling | unknown | M1-T41 | | |
| Deterministic simulation: virtual time, seeds, event loop | unknown | M1-T42 | | |
| Fault injection: delay, drop, duplicate, reorder, partitions, schedules | unknown | M1-T49 | | |
| Trace recording and hashing; reproducibility; nondeterminism sources | unknown | M1-T55 | | |
| Porcupine: Model, Operation, partitions, visualization | unknown | M1-T64 | | |
| Sequential specification of a KV store | unknown | M1-T65 | | |
| Concurrent histories, indeterminate ops, recording rules | unknown | M1-T68 | | |
| Latency percentiles, simple benchmarks | unknown | M1-T60 | | |
| Estimation and retros | unknown | M1-T78 | | |
| Goroutines, channels, `select`, `context`, `sync`, race detector | unknown | M5 (concurrency warm-up) | not used in real code | |
| Protobuf, `buf`, gRPC services, status codes, deadlines | unknown | M5 | | |
| gRPC health service | unknown | M5 | | |
| Write-ahead log design, torn writes, fail-stop | unknown | M4 | | |
| Fuzz testing (`go test -fuzz`) | unknown | M4 | | |
| Mutation testing | unknown | M4 | | |
| Snapshots, InstallSnapshot, log compaction | unknown | M7 | | |
| LSM trees: memtable, SSTables, compaction, amplification | unknown | M8 | | |

## Later (interesting but off-plan)
- PreVote / CheckQuorum, membership changes, multi-Raft
