---
id: E23-board-launch-actions/T003-defect-build-and-audit-actions
status: planned
objective: Add defect Build and optional Audit actions with canonical defect lifecycle handling.
depends_on:
  - E23-board-launch-actions/T001-board-launch-availability
  - E23-board-launch-actions/T002-task-build-and-audit-actions
complexity_tier: medium
complexity_reason: Defect actions reuse launch orchestration but require separate lifecycle and overlay behavior.
---

# T003: Defect Build and Audit Actions

## Problem

Defects can be inspected and resolved from the board, but cannot start a bounded builder or auditor session from their detail view.

## Context Files

- `internal/board/defect_detail.go`
- `internal/board/defect_detail_test.go`
- `internal/board/defect_overlay.go`
- `internal/board/defect_overlay_test.go`
- `internal/board/update.go`
- `internal/board/io.go`
- `internal/board/launch_actions.go`
- `internal/data/defect.go`
- `internal/data/defect_test.go`
- `internal/data/write.go`
- `internal/launcher/request.go`

## Acceptance Criteria

- [ ] The defect detail view shows Build and eligible Audit actions only when enabled.
- [ ] Build moves an open defect to `status: in_progress` and `stage: build` before launch.
- [ ] Audit changes an in-progress defect to `stage: audit` and never resolves it.
- [ ] Resolved defects cannot dispatch Build or Audit.
- [ ] Defect router priority and launch prompts identify the exact release-level defect file.
- [ ] Task action regressions remain covered while shared orchestration is reused.

## Implementation Plan

- [ ] Add defect action labels and key handling to the detail overlay.
- [ ] Extend shared launch orchestration with canonical defect status writes and router state.
- [ ] Build launch requests from the selected defect ID, path, release, and optional reference.
- [ ] Reuse in-flight suppression and typed launch result messages.
- [ ] Add tests for open, in-progress, resolved, disabled, mtime-conflict, and launch-error cases.

## Context Log

Pending.
