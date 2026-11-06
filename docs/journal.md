# Journal

Newest first. Entry types: decision, change, retro, note.

## 2026-11-05: note: Setup block done (M1-T01..T13)
- What: toolchain fixed, repo public on GitHub, Makefile, golangci-lint with forbidigo determinism bans, CI (lint, race tests, govulncheck, image build + smoke test) with per-job cache keys, raftkv config parsing, 8.45 MB distroless image, README + MIT.
- Planned vs reality: 13 tasks in one day (estimate ~4 h). Most review rounds caught typos that "worked by accident" (`ling` skipped lint, `golangci.yml` without dot disabled bans, `CGO_ENABLE`) — pattern: verify the check actually ran, not just that it's green.
- Process: user skips comprehension questions; for Commodity tasks code/config is given directly, Learn tasks are written by the user.
- Affected: `tasks/M1-T01..T13`, `STATUS.md`

## 2026-11-04: decision: Planning approved
- What: 8 milestones (MVP = M1-M6, ~336 h, projected ~2027-03-22, target 2027-04-15); Checkpoint A at end of M4; ordered cut list; 79 M1 task cards; M2-M8 as task groups.
- Why: Bottom-up estimate replaces the rough 460 h; cards beyond the next milestone would go stale.
- Affected: `03-plan.md`, `tasks/`, `knowledge.md`

## 2026-11-04: change: Milestones renumbered for the plan
- What: The vision's "M2 snapshots" and "M3 LSM" became **M7** and **M8**; the MVP is split into vertical slices M1-M6. Old "M1" references in ADRs were mapped to the specific milestone (e.g. smoke test → M5, mutation suite → M4). "M1a checkpoint" became **Checkpoint A (end of M4)**.
- Why / planned vs reality: Planning needs small milestones with demos; one "M1" for the whole MVP was too large.
- Affected: 00-vision, 01-architecture, ADR-001, 002, 006, 007

## 2026-11-04: decision: Verification approved
- What: 2 rounds of checks and fresh-eyes review; 9 blocking issues fixed (one moved to Planning: estimate, cut list, Checkpoint A). ADR-006 and ADR-007 added; ADR-005 superseded.
- Why: Correctness angle needed precise history recording, WAL recovery, node identity and an honest prior-art statement.
- Affected: all docs; `02-verification.md`

## 2026-11-04: decision: Architecture approved
- What: State-machine Raft (ADR-001), deterministic simulation + Porcupine (ADR-002), sessions + reads through the log (ADR-003), gRPC (ADR-004), storage design (ADR-005, later superseded).
- Why: Reproducible chaos runs are the project's angle.
- Affected: ADR-001..005, `01-architecture.md`

## 2026-11-04: decision: Discovery approved
- What: raft-kv: Raft + storage from scratch in Go; angle: correctness under injected faults; Deep: Raft and storage.
- Why: Deep learning goal + junior Go backend job hunting; 20+ h/week, no hard deadline.
- Affected: `00-vision.md`
