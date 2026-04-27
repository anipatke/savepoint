---
id: E02-data-model/T001-domain-ids-status
status: done
objective: "Define pure domain primitives for release, epic, task IDs, task statuses, and allowed status transitions."
depends_on: []
---

# T001: Domain IDs and Statuses

## Scope

Introduce the core domain types and validation helpers that later readers and validators can share.

## Acceptance Criteria

- Release, epic, and task ID parsing accepts the E02 design's expected ID shapes.
- Task status values are represented from one source of truth.
- Status transition validation covers allowed and disallowed transitions.
- Branching validation behavior has focused unit tests.

## Implementation Plan

- [x] Inspect existing package/test conventions needed to place domain modules and focused unit tests.
- [x] Add pure ID parsing and formatting helpers for release IDs, epic IDs, task numbers, and full task IDs.
- [x] Add task status constants/types and one transition rule table for `backlog`, `planned`, `in_progress`, `review`, and `done`.
- [x] Add validation helpers that return typed success/failure results without filesystem access.
- [x] Add unit tests for valid and invalid ID shapes plus allowed and disallowed status transitions.
