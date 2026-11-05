# raft-kv

[![CI](https://github.com/bakdauletbaktygaliyev/raft-kv/actions/workflows/ci.yml/badge.svg)](https://github.com/bakdauletbaktygaliyev/raft-kv/actions/workflows/ci.yml)

A distributed, strongly consistent key-value store in Go, with Raft and a
write-ahead log written from scratch. Every chaos test run is a deterministic
simulation that replays byte-for-byte from a single seed, including network
partitions and lost disk writes, and is checked for linearizability with
[Porcupine](https://github.com/anishathalye/porcupine), Raft safety invariants
and progress checks. A mutation suite plants known bugs to prove the checks
catch them.

## Status

Early: walking skeleton. CI, lint (with determinism rules), Docker image and
node config are in place; Raft, the WAL and the simulator are not implemented
yet. See the [plan](docs/03-plan.md) for milestones.

## Build and run

Requires Go 1.27.

```sh
make            # build, lint, test (with -race)
make vuln       # govulncheck

go run ./cmd/raftkv --id 1 --peers 1=n1:7000,2=n2:7000,3=n3:7000
```

With Docker:

```sh
docker build -t raftkv .
docker run --rm raftkv --id 1 --peers 1=n1:7000,2=n2:7000,3=n3:7000
```

## Docs

- [Vision](docs/00-vision.md): goals, MVP, how it differs from similar projects
- [Architecture](docs/01-architecture.md) and [ADRs](docs/adr/)

## License

[MIT](LICENSE)
