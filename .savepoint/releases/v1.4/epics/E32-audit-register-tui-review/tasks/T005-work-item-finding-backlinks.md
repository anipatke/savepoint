---
id: E32-audit-register-tui-review/T005-work-item-finding-backlinks
status: planned
objective: Surface audit findings linked to the focused epic or task inside their detail overlays.
depends_on:
  - E32-audit-register-tui-review/T003-finding-list-and-detail
complexity_tier: medium
complexity_reason: Adds reverse-lookup rendering plus cross-overlay navigation back into finding detail.
---

# T005: Work-Item Finding Backlinks

## Problem

Findings link forward to epics and tasks, but a reviewer reading an epic or task cannot see which audit findings touch it. The relationship is only navigable from the audit overlay, so audit context is invisible exactly when someone is reviewing the work item it belongs to.

## Context Files

- `internal/board/detail.go`
- `internal/board/detail_test.go`
- `internal/board/epic_panel.go`
- `internal/board/epic_panel_test.go`
- `internal/board/audit_detail.go`
- `internal/board/update.go`
- `internal/board/model.go`
- `internal/data/audit.go`

## Acceptance Criteria

- [ ] Task detail shows a "Linked Findings" section listing findings whose `tasks` include the task ID, with finding ID, title, status, and severity.
- [ ] Epic detail surfaces findings whose `epics` include the epic, consistent with its existing tab/section layout.
- [ ] Work items with no linked findings render a clear empty state, not a missing section.
- [ ] `enter` on a linked finding opens the existing read-only finding detail renderer.
- [ ] `esc` and `q` return from finding detail to the originating epic or task detail, restoring its prior scroll/cursor.
- [ ] No key mutates finding status, links, or audit files.
- [ ] Lookups tolerate an absent audit register and never block detail rendering.

## Implementation Plan

- [ ] Add a reverse-lookup helper that maps a work-item ID to its linked findings from loaded audit data.
- [ ] Render the linked-findings section in task detail and epic detail using existing list styles.
- [ ] Add cursor state for the linked-findings list and clamp it.
- [ ] Wire enter/escape between work-item detail and the shared finding detail renderer, preserving origin state.
- [ ] Add tests for link matching, empty states, navigation round-trip, and no-op mutation keys.

## Context Log

Pending.
