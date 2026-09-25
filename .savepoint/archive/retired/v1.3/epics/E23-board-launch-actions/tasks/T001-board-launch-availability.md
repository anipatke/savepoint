---
id: E23-board-launch-actions/T001-board-launch-availability
status: planned
objective: Load launcher state into the board and centralize action eligibility rules.
depends_on:
  - E22-agent-launcher-foundation/T004-launcher-service
complexity_tier: medium
complexity_reason: Availability combines configuration, lifecycle, selection, and router state across board views.
---

# T001: Board Launch Availability

## Problem

The board needs one source of truth for whether Build or Audit is visible and allowed for the currently selected task, defect, or epic.

## Context Files

- `internal/board/model.go`
- `internal/board/model_test.go`
- `internal/board/interfaces.go`
- `internal/board/interfaces_test.go`
- `internal/board/tui.go`
- `internal/board/launch_actions.go`
- `internal/board/launch_actions_test.go`
- `internal/data/config.go`
- `internal/launcher/launcher.go`

## Acceptance Criteria

- [ ] Launcher state and an injected launcher service are available to board update commands.
- [ ] No action is visible or dispatchable when the launcher is disabled.
- [ ] Build requires a configured builder and a non-complete task or defect.
- [ ] Item Audit requires a configured auditor and an in-progress task or defect.
- [ ] Epic Audit requires a configured auditor plus matching release/epic router state at `audit-pending`.
- [ ] Eligibility logic is centralized and table-tested rather than repeated in each overlay.

## Implementation Plan

- [ ] Add launcher configuration and launch dependency fields to focused board state structs.
- [ ] Load config once during project model startup and preserve existing theme behavior.
- [ ] Define action descriptors and centralized eligibility functions for all item types.
- [ ] Add table-driven tests for disabled, missing-role, lifecycle, router mismatch, and valid cases.
- [ ] Keep plain non-TTY output unchanged.

## Context Log

Pending.
