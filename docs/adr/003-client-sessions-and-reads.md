# ADR-003: Client sessions for deduplication; reads go through the log

- Status: accepted (2026-11-04); contract detailed in ADR-006
- Date: 2026-11-04
- Component label: Deep
- Decision type: optimal (for MVP)

## Context
Porcupine will reject histories where a retried write is applied twice or a deposed leader serves a stale read.

## Alternatives considered
| Option | Pros | Cons | Evidence |
|---|---|---|---|
| No dedup | Simple | Retries after leader crash apply twice; not linearizable | [S1], [S2] |
| Sessions: `clientID` + sequence number, state machine stores last seq + response per client | Standard Raft answer; exactly-once effect | Dedup table must be persisted and included in snapshots | [S1], [S2] |
| Reads through the log | Trivially linearizable; no extra code | Slower (each read is replicated) | [S1] |
| ReadIndex (quorum heartbeat, then local read) | Faster reads | More code and more subtle | [S3] |
| Read from local state, no check | Fastest | Stale reads; not linearizable | [S3] |

## Decision
Client sessions with deduplication; `Get` goes through the log in the MVP. ReadIndex is in "Later" as a measured optimization.

## Why
- Each client tags commands with a unique serial number; the state machine skips duplicates `[source: S1]`, `[source: S2]`
- Without extra precautions a leader can serve stale reads; going through the log avoids them `[source: S3]`

## Cost (what we pay)
- Read latency equals write latency
- Dedup table is part of the state machine; must be in snapshots (milestone 2)
- Session expiry is not handled in the MVP (the table grows with clients; fine for tests)

## Revisit if
- Not tracked by a trigger: ReadIndex stays in "Later" as an optional improvement (no read benchmark is planned for the MVP)

## Sources
| ID | Link | Date checked | Supports which claim |
|---|---|---|---|
| S1 | https://github.com/ongardie/dissertation | 2026-11-04 | Raft dissertation ch. 6 (client interaction) |
| S2 | https://github.com/lorin/raft-pluscal | 2026-11-04 | Unique id per command so it is not executed twice |
| S3 | https://emptysqua.re/blog/review-rethink-linearizability-constraints-raft/ | 2026-11-04 | ReadIndex and stale-read precautions |

## Uncertainty
- Chapter 6 details were read via secondary summaries; read the dissertation section itself during the milestone.
