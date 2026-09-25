---
id: E28-parallel-worktree-launch/T003-board-parallel-launch
status: planned
objective: Add an advanced board action that creates eligible worktrees and launches one builder per task.
depends_on:
  - E22-agent-launcher-foundation/T004-launcher-service
  - E23-board-launch-actions/T002-task-build-and-audit-actions
  - E28-parallel-worktree-launch/T002-runtime-manifests-and-scope
complexity_tier: high
complexity_reason: One TUI action coordinates eligibility, Git mutations, manifests, multiple launches, and partial results.
---

# T003: Board Parallel Launch

## Problem

Eligible parallel groups cannot be launched or inspected from the board without manually coordinating Git and terminals.

## Context Files

- `internal/board/model.go`
- `internal/board/interfaces.go`
- `internal/board/update.go`
- `internal/board/view.go`
- `internal/board/help.go`
- `internal/board/parallel_actions.go`
- `internal/board/parallel_actions_test.go`
- `internal/worktree/service.go`
- `internal/launcher/launcher.go`

## Acceptance Criteria

- [ ] Parallel launch is hidden when advanced mode is disabled or fewer than two tasks are eligible.
- [ ] The action shows the selected group and task count before creation and respects `max_agents`.
- [ ] Worktree creation completes before one configured interactive builder is launched in each successful worktree.
- [ ] Per-task creation and launch outcomes are reported without process monitoring or completion claims.
- [ ] Active manifest entries prevent duplicate dispatch for the same task.
- [ ] Sequential Build and Audit actions retain their existing behavior.

## Implementation Plan

- [ ] Inject parallel eligibility, worktree, and launcher dependencies into board state.
- [ ] Add centralized action availability and confirmation data for the focused task group.
- [ ] Implement Bubble Tea commands for worktree creation followed by bounded parallel terminal launch.
- [ ] Record launch outcomes in the coordinator manifest and surface partial failures.
- [ ] Add fake-service tests for disabled, ineligible, successful, duplicate, and partial-failure paths.

## Context Log

Pending.
