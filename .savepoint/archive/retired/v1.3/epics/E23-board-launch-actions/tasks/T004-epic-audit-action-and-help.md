---
id: E23-board-launch-actions/T004-epic-audit-action-and-help
status: planned
objective: Add gated epic Audit launching and complete board help and status integration.
depends_on:
  - E23-board-launch-actions/T002-task-build-and-audit-actions
  - E23-board-launch-actions/T003-defect-build-and-audit-actions
complexity_tier: medium
complexity_reason: Epic gating and shared help/status updates touch several established board interaction paths.
---

# T004: Epic Audit Action and Help

## Problem

The epic detail view cannot start the required fresh audit session, and users need one coherent explanation of all launcher keys and availability rules.

## Context Files

- `internal/board/epic_panel.go`
- `internal/board/epic_panel_test.go`
- `internal/board/update.go`
- `internal/board/update_test.go`
- `internal/board/help.go`
- `internal/board/help_test.go`
- `internal/board/view.go`
- `internal/board/view_test.go`
- `internal/board/launch_actions.go`
- `internal/data/router.go`
- `internal/launcher/request.go`

## Acceptance Criteria

- [ ] Epic detail shows Audit only for the router-selected epic at `audit-pending` with an auditor configured.
- [ ] Epic Audit starts a new terminal session and does not mutate epic status or create an audit file itself.
- [ ] Mismatched release, epic, router state, disabled config, and missing auditor cannot dispatch the action.
- [ ] Help and detail footers explain Build/Audit keys only when the launcher is enabled.
- [ ] Success and failure status messages identify role, action, and selected item without exposing the full prompt.
- [ ] Existing epic tabs, scrolling, and navigation remain unchanged.

## Implementation Plan

- [ ] Add the gated Audit affordance to the epic detail interaction path.
- [ ] Build an epic audit request from router and selected epic state.
- [ ] Reuse shared launch messages and in-flight suppression without lifecycle writes.
- [ ] Update contextual help and footer rendering for enabled action keys.
- [ ] Add integration tests for all epic gate outcomes and existing overlay navigation.

## Context Log

Pending.
