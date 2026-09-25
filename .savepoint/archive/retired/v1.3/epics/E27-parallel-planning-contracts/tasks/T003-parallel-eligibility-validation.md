---
id: E27-parallel-planning-contracts/T003-parallel-eligibility-validation
status: planned
objective: Centralize dependency-aware parallel-group eligibility and exact write-scope conflict detection.
depends_on:
  - E27-parallel-planning-contracts/T001-advanced-parallel-config
  - E27-parallel-planning-contracts/T002-task-execution-metadata
complexity_tier: high
complexity_reason: Eligibility combines lifecycle, dependencies, path ownership, config limits, and deterministic diagnostics.
---

# T003: Parallel Eligibility Validation

## Problem

Board and doctor need one deterministic answer for whether a set of tasks may launch in parallel.

## Context Files

- `internal/data/parallel.go`
- `internal/data/parallel_test.go`
- `internal/data/dependency.go`
- `internal/data/task.go`
- `internal/data/lifecycle.go`

## Acceptance Criteria

- [ ] Eligibility requires advanced mode, at least two planned tasks in one group, and no completed or active task.
- [ ] Every selected task has satisfied dependencies and a valid non-empty exact write scope.
- [ ] Duplicate paths and file-versus-parent-directory collisions produce task-specific conflicts.
- [ ] Selection respects `max_agents` and returns a stable eligible subset or an explicit rejection.
- [ ] The result exposes structured reasons suitable for board messages and doctor diagnostics.
- [ ] Tasks without advanced metadata are ignored rather than made invalid globally.

## Implementation Plan

- [ ] Define typed parallel candidate, conflict, and eligibility result models.
- [ ] Reuse canonical dependency resolution and lifecycle constants.
- [ ] Implement normalized exact-path ownership comparison with stable ordering.
- [ ] Add table tests for readiness, group mismatch, limits, duplicate ownership, parent collisions, and valid groups.
- [ ] Keep the validator free of Git, board, and filesystem process behavior.

## Context Log

Pending.
