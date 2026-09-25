---
id: E28-parallel-worktree-launch/T001-worktree-service
status: planned
objective: Add a testable Git worktree service with deterministic preflight, creation, and inspection.
depends_on:
  - E27-parallel-planning-contracts/T003-parallel-eligibility-validation
complexity_tier: high
complexity_reason: Git state, cross-platform process execution, naming collisions, and partial failure require strict boundaries.
---

# T001: Worktree Service

## Problem

Savepoint has no safe boundary for inspecting Git state or creating isolated task branches and worktrees.

## Context Files

- `internal/worktree/service.go`
- `internal/worktree/service_test.go`
- `internal/worktree/git.go`
- `internal/worktree/git_test.go`
- `internal/data/parallel.go`

## Acceptance Criteria

- [ ] Preflight requires a Git worktree, a clean primary checkout, a resolvable base commit, and eligible tasks.
- [ ] Branch and worktree paths are deterministic, sanitized, and collision-checked before mutation.
- [ ] Git commands use structured arguments with explicit working directories and captured stderr.
- [ ] Creation returns per-task results and rolls back only empty worktrees created by the failed operation.
- [ ] Existing branches or paths are never overwritten or force-deleted.
- [ ] Tests run in temporary repositories and cover success, dirty state, collisions, and partial failure.

## Implementation Plan

- [ ] Define Git runner and filesystem interfaces for testable process boundaries.
- [ ] Implement repository, cleanliness, base-commit, branch, and path preflight.
- [ ] Define deterministic branch and worktree naming from release, epic, and task IDs.
- [ ] Implement sequential worktree creation with structured per-task outcomes.
- [ ] Add temporary-repository integration tests without invoking the agent launcher.

## Context Log

Pending.
