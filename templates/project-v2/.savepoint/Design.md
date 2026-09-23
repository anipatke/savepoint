---
type: project-design
status: active
---

# {{PROJECT_NAME}} — Design

> Design describes implemented reality; it is not a backlog. Planned change lives in the active Objective's deltas until the work lands, at which point this file is reconciled to match. Fill each section with concrete, falsifiable claims — a generic Design produces a generic agent.

## Architecture

List the architecture commitments this project follows, each as a rule the system actually follows, not a preference.

## Components/Codebase Map

Table of module → purpose, kept current as code lands. A new responsibility means a new row, not a silent addition to an existing one.

## Interfaces and Data Flow

The settled interfaces and data ownership between components — what `savepoint-design`'s Readiness Gate checks before an Objective's Tasks can be detailed.

## Boundaries

What this project explicitly does and does not do.

## Decisions

Durable technical decisions and the reasoning behind them, recorded so a later reader does not have to rediscover why.

## Current Technical State

What exists today, concretely — not what's planned. When adopting Savepoint into an existing codebase, this is the section to reconstruct first: read targeted evidence from the code, ask the owner for intent on anything the code alone cannot answer, and record what's actually there. Record only what was verified by reading, and name what has not yet been examined as unknown rather than leaving it implied. Never invent this section from assumption, and never let an automatic scan write it.

## Verification Contract

This contract applies to every implementation:

- Every Task records implementation evidence and configured quality-gate
  results before handoff.
- A Task Check is optional. If the owner skips it, record an explicit waiver
  in the Task evidence naming the Task, reason, actor, and time. The waiver is
  not technical `CLEAR` and does not waive acceptance criteria or guardrails.
  It satisfies a Task dependency that requires `clear` — the waiver stands in
  as the owner's own completion decision — but never one that requires
  `accepted`, since there is no Check for the owner to have accepted.
  Pressing Space on the V2 board to complete a Task at stage check with no
  recorded Check at all is itself that explicit owner action — the board
  auto-records the waiver rather than requiring it written by hand first —
  but never when a recorded Check actually found a problem; that result
  stands.
- The Full Objective Check is mandatory before Objective closure. It is the V2
  equivalent of the epic-level integration gate and reviews every owned Task,
  including waived Tasks, cross-Task integration, and reconciliation against
  this Design.
- A Goal Check is mandatory whenever a Goal exists, followed by exact owner
  acceptance of the current Check.
