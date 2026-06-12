---
id: E23-board-launch-actions/T002-task-build-and-audit-actions
status: planned
objective: Add task Build and optional Audit actions with canonical lifecycle and router updates.
depends_on:
  - E23-board-launch-actions/T001-board-launch-availability
complexity_tier: high
complexity_reason: Task launch coordinates overlay input, two file writes, async startup, and failure reporting.
---

# T002: Task Build and Audit Actions

## Problem

A focused task cannot currently start its configured builder or auditor directly from the detail view.

## Context Files

- `internal/board/detail.go`
- `internal/board/detail_test.go`
- `internal/board/update.go`
- `internal/board/update_test.go`
- `internal/board/io.go`
- `internal/board/transitions.go`
- `internal/board/transitions_test.go`
- `internal/board/launch_actions.go`
- `internal/data/lifecycle.go`
- `internal/data/router.go`
- `internal/data/write.go`
- `internal/launcher/request.go`

## Acceptance Criteria

- [ ] The task detail view shows Build and eligible Audit actions only when enabled.
- [ ] Build moves a planned task to `status: in_progress` and `stage: build`, then sets matching router priority.
- [ ] Build on an already in-progress task preserves valid lifecycle state and launches the selected task.
- [ ] Audit changes an in-progress task to `stage: audit` and never changes it to `done`.
- [ ] Required writes use canonical data helpers and mtime conflict handling before launch dispatch.
- [ ] Repeated keys are ignored while the same action is in flight, and failures produce actionable status text.

## Implementation Plan

- [ ] Render action labels and assign non-conflicting detail-overlay keys for Build and Audit.
- [ ] Add a Bubble Tea command sequence for lifecycle write, router write, and launch dispatch.
- [ ] Build the launch request from the selected task's release, epic, ID, path, and project root.
- [ ] Track the active launch key until a success or error message returns.
- [ ] Add reducer and rendering tests for planned, in-progress, done, disabled, unavailable-auditor, conflict, and launch-error cases.

## Context Log

Pending.
