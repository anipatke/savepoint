---
id: E32-audit-register-tui-review/T005-work-item-finding-backlinks
status: done
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

- [x] Task detail shows a "Linked Findings" section listing findings whose `tasks` include the task ID, with finding ID, title, status, and severity.
- [x] Epic detail surfaces findings whose `epics` include the epic, consistent with its existing tab/section layout.
- [x] Work items with no linked findings render a clear empty state, not a missing section.
- [x] `enter` on a linked finding opens the existing read-only finding detail renderer.
- [x] `esc` and `q` return from finding detail to the originating epic or task detail, restoring its prior scroll/cursor.
- [x] No key mutates finding status, links, or audit files.
- [x] Lookups tolerate an absent audit register and never block detail rendering.

## Implementation Plan

- [x] Add a reverse-lookup helper that maps a work-item ID to its linked findings from loaded audit data.
- [x] Render the linked-findings section in task detail and epic detail using existing list styles.
- [x] Add cursor state for the linked-findings list and clamp it.
- [x] Wire enter/escape between work-item detail and the shared finding detail renderer, preserving origin state.
- [x] Add tests for link matching, empty states, navigation round-trip, and no-op mutation keys.

## Context Log

### Read
- `.savepoint/router.md`, epic tasks T005 (this file)
- `internal/board/detail.go`, `epic_panel.go`, `audit_detail.go`, `audit_overlay.go`, `model.go`, `update.go`, `view.go`
- `internal/data/audit.go`, `audit_finding.go`, `audit_validate.go`, `dependency.go`

### Edited / added
- `internal/data/audit_backlinks.go` (+ test): `FindingsForTask`/`FindingsForEpic` reverse lookup. Matching mirrors `ValidateAuditFindings`/`ResolveDependency` (full ID or short `T###`/`E##`), preserves the loader's finding sort order, and returns nil for an empty ID or register so an absent audit tree never blocks rendering.
- `internal/board/audit_backlinks.go` (+ test): shared `linkedFindingsSection`/`renderLinkedFindingRow` (severity tag, ID, title, `[status]`, cursor marker) reused by both detail overlays; `detailFooterHint` adds `enter:open` only when findings exist; model helpers `focusedTaskFindings`, `openEpicFindings`, `moveLinkedFindingCursor`, `activeFinding`, `findingDetailReturnOverlay`.
- `internal/board/model.go`: added `LinkedFindingCursor` and `FindingDetailOrigin` to `AuditState`.
- `internal/board/detail.go` / `epic_panel.go`: `RenderDetail`/`RenderEpicDetail` take `findings []data.AuditFinding, findingCursor int` and render the section (epic: Detail tab only).
- `internal/board/update.go`: task/epic detail up/down drive the finding cursor when findings are present (else scroll); `enter` opens finding detail and records origin; finding-detail `esc`/`q` return to origin; cursor resets on detail open and epic tab switch; audit-tab `enter` now records `OverlayAudit` origin.
- `internal/board/view.go`: overlays pass linked findings + cursor; finding detail renders `activeFinding()` resolved from origin.
- Updated existing `RenderDetail`/`RenderEpicDetail` call sites in `detail_test.go`, `epic_panel_test.go`, `render_policy_test.go`.

### Quality gates
- `make build` — ok
- `make test` — all packages pass; `go vet ./internal/board ./internal/data` clean.

## Drift Notes

- Added two files within existing modules for existing responsibilities: `internal/data/audit_backlinks.go` (audit reverse-lookup, extends the `internal/data` audit model) and `internal/board/audit_backlinks.go` (linked-findings rendering + navigation, extends the board overlays). No new module or architectural change versus the Codebase Map; the board still reads all audit data through `data.AuditRegisterSet`.
