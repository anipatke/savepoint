---
type: epic-design
status: planned
---

# E23: Board Launch Actions

## Purpose

Expose launcher actions from task, defect, and epic detail views while preserving canonical lifecycle rules, router priority, and existing disabled behavior.

## What this epic adds

- Launcher configuration and service dependencies in the board model.
- Action availability rules based on item type, lifecycle, router state, and configured role.
- Task Build and optional Audit actions.
- Defect Build and optional Audit actions.
- Epic Audit at the active `audit-pending` gate.
- Non-blocking Bubble Tea commands, in-flight suppression, and actionable status messages.

## Components and files

| Module | Purpose |
|--------|---------|
| `internal/board/model.go` | Store launcher state and in-flight action state |
| `internal/board/interfaces.go` | Inject launcher behavior into board tests |
| `internal/board/update.go` | Route action keys and handle launch results |
| `internal/board/detail.go` | Render task action affordances |
| `internal/board/defect_detail.go` | Render defect action affordances |
| `internal/board/epic_panel.go` | Render epic Audit affordance |
| `internal/data` | Apply canonical lifecycle and router writes before launch |

## Architectural delta

The board gains outbound launch commands using the existing Bubble Tea message pattern. Lifecycle and router writes remain in data-owned helpers; the launcher receives a request only after action eligibility and required writes are resolved.

## Boundaries

**In scope:**
- Keyboard actions and detail-view labels
- Eligibility rules and disabled-state invisibility
- Canonical task/defect stage and router updates
- In-flight duplicate suppression and status messages

**Out of scope:**
- Monitoring terminal processes after startup
- Completing or resolving items
- Editing launcher configuration from the TUI
- Launching from non-interactive plain output

## Quality gates

- Board reducer and rendering tests cover the complete enabled/disabled action matrix.
- Existing navigation, status transitions, and overlay behavior remain unchanged.
- `go test ./internal/board ./internal/data` passes.

## Open decisions

None.
