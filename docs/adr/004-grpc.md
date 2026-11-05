# ADR-004: gRPC + protobuf for peer transport and client API

- Status: accepted (2026-11-04), revised in verification round 1
- Date: 2026-11-04
- Component label: Learn
- Decision type: learning choice (chosen to learn gRPC for Go backend roles; not claimed optimal)

## Context
Raft returns messages in `Ready` (ADR-001); a `Transport` adapter delivers them. In the simulator the adapter is replaced, so this choice affects the real cluster and the client API.

## Alternatives considered
| Option | Pros | Cons | Evidence |
|---|---|---|---|
| grpc-go + protobuf | Widely used in Go infra (etcd, Kubernetes); deadlines, status codes, interceptors | Codegen step; not `curl`-able; no semver | [S1], [S3] |
| connect-go | gRPC and HTTP/JSON; slimmer | Less widely used | [S2] |
| `net/http` + JSON | Nothing new | Misses the learning goal | — |

## Decision
grpc-go **v1.85.0** (pinned) with `raft.proto` and `kv.proto` (ADR-006). Codegen with **`buf`** (installed 1.75.0; local `protoc` is broken) using `protoc-gen-go` and `protoc-gen-go-grpc`. Generated code is committed.

Real-cluster coverage: CI runs one **real-process smoke test** (Docker Compose, 3 nodes: write, kill leader, check new leader < 2 s and no lost writes).

## Why
- Learning gRPC for Go backend roles `[learning choice]`
- grpc-go v1.85.0 exists, Apache-2.0, requires Go ≥ 1.26 `[source: S1]`
- grpc-go does not follow semver; pin and read "Behavior Changes" `[source: S3]`
- `buf` works locally, `protoc` fails (missing abseil dylib) `[experiment: verification round 1]`

## Cost (what we pay)
- Protobuf toolchain in build and CI
- Manual testing needs `grpcurl` or `kvctl`

## Revisit if
- A browser UI becomes a goal → connect-go

## Sources
| ID | Link | Date checked | Supports which claim |
|---|---|---|---|
| S1 | `go list -m -versions google.golang.org/grpc` + module LICENSE/go.mod | 2026-11-04 | v1.85.0, Apache-2.0, go 1.26 |
| S2 | https://github.com/connectrpc/connect-go | 2026-11-04 | connect-go features |
| S3 | https://buf.build/blog/connect-a-better-grpc | 2026-11-04 | No semver in grpc-go |
