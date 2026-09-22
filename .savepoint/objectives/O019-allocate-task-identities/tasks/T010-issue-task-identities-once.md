---
id: T010
title: Issue each Task number once
objective: O019
planned_by: {role: planner, session: i024-allocation-design-20260923}
status: planned
complexity_tier: high
complexity_reason: Durable allocation and concurrent writers must behave consistently across supported platforms.
depends_on: []
owner_validation: {required: false}
---

# T010: Issue each Task number once

## Outcome

A project-owned allocation operation issues monotonically increasing Task IDs under a project lock, including when Tasks are deleted or creation fails.

## Done When

- Strict-load before allocation; refuse invalid projects without changing records.
- Initialize a durable high-water mark from all active Tasks; never lower it after deletion or failure.
- Concurrent callers cannot receive the same number. Lock failure and interrupted reservation have explicit Linux and Windows behavior.
- Existing duplicate-ID validation still fails closed and names both paths.

## Context Files

`internal/data/project.go`, `internal/data/discover.go`, `internal/data/task_v2.go`, `internal/data/write.go`, `internal/data/project_test.go`, `internal/data/discover_test.go`, `internal/data/task_v2_test.go`, `internal/data/write_test.go`, `internal/migrate/replace.go`, `internal/migrate/replace_unix.go`, `internal/migrate/replace_windows.go`, `AGENTS.md`.

## Design References

Design sections 1, 2, 5, and 6; O019 Architectural Considerations.

## Guardrails

FS-01, FS-04, FS-05, FS-06, DATA-03, CFG-02, ARCH-03, ARCH-04, TEST-01, TEST-02, TEST-03, TEST-04, TEST-05, TEST-08.

## Implementation Plan

1. Confirm O018's grammar and existing platform-safe writes; return REPLAN REQUIRED if interfaces contradict this plan.
2. Add durable Task-number high-water state and a project lock; initialize existing projects from their complete active Task set.
3. Persist the next reservation while holding the lock. Never roll it back after failure.
4. Test bootstrap, deletion, failure, concurrency, interruption, invalid projects, and the unchanged two-path duplicate diagnostic.

## Boundaries

No allocator for other record kinds, Task lifecycle change, or weakening of strict loading.

## Technical Verification

Focused `internal/data` tests, `git diff --check`, and `make build && make test`.

## Technical Evidence

Pending execution: per-criterion outcomes, command results, files read/changed, and limitations.

## Drift Notes

If cross-platform locking or crash recovery requires a different storage contract, return REPLAN REQUIRED.
