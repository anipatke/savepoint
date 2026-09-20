---
id: E28-parallel-worktree-launch/T004-worktree-status-and-cleanup
status: planned
objective: Show active parallel worktrees and provide explicit cleanup that refuses unsafe removal.
depends_on:
  - E28-parallel-worktree-launch/T003-board-parallel-launch
complexity_tier: high
complexity_reason: Cleanup must reconcile manifests, Git, filesystem state, and dirty branches without destructive shortcuts.
---

# T004: Worktree Status and Cleanup

## Problem

Users need to understand and retire local parallel worktrees without Savepoint deleting unintegrated work.

## Context Files

- `internal/worktree/service.go`
- `internal/worktree/service_test.go`
- `internal/worktree/manifest.go`
- `internal/board/parallel_actions.go`
- `internal/board/parallel_actions_test.go`
- `internal/board/detail.go`
- `internal/board/help.go`

## Acceptance Criteria

- [ ] Board detail shows task, branch, path, base commit, and launch state for active parallel entries.
- [ ] Cleanup requires an exact manifest, worktree path, and branch match plus a clean worktree.
- [ ] Dirty, missing, locked, or mismatched worktrees are refused with a recoverable explanation.
- [ ] Cleanup removes the worktree and runtime entry but never deletes the task branch automatically.
- [ ] Turning advanced mode off hides actions while preserving status data for later re-enable or doctor inspection.
- [ ] Tests prove no force-removal or branch-deletion command is issued.

## Implementation Plan

- [ ] Add read-only active-entry summaries to task detail and advanced action views.
- [ ] Implement cleanup preflight against manifest and live Git state.
- [ ] Remove only clean matching worktrees through structured Git arguments.
- [ ] Update the coordinator manifest atomically after confirmed removal.
- [ ] Add tests for clean, dirty, missing, mismatched, disabled, and already-removed cases.

## Context Log

Pending.
