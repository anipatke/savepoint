---
id: E32-audit-register-tui-review/T002-audit-overlay-tabs
status: done
objective: Add the read-only Audit Register overlay with Prompt, Findings, and Runs tabs.
depends_on:
    - E32-audit-register-tui-review/T001-board-audit-data-loading
complexity_tier: medium
complexity_reason: The overlay reuses existing patterns but adds a new top-level interaction path.
---

# T002: Audit Overlay Tabs

## Problem

Users need a dedicated TUI section for audit prompt, current findings, and run history instead of hunting through markdown files.

## Context Files

- `internal/board/model.go`
- `internal/board/update.go`
- `internal/board/view.go`
- `internal/board/audit_overlay.go`
- `internal/board/audit_overlay_test.go`
- `internal/board/epic_panel.go`
- `internal/board/release_docs_overlay_test.go`

## Acceptance Criteria

- [x] Pressing `A` opens the Audit Register overlay.
- [x] Overlay tabs are Prompt, Findings, and Runs.
- [x] `[`/`]`, left/right, and `h`/`l` switch tabs consistently with existing document overlays.
- [x] `up`/`down`, `k`/`j`, `pgup`, and `pgdown` scroll the selected tab body.
- [x] `esc` and `q` close the overlay.
- [x] Missing prompt, findings, or runs render useful empty states.

## Implementation Plan

- [x] Add audit overlay constants and navigation state.
- [x] Add top-level `A` key handling.
- [x] Render the tab strip and selected body using existing overlay styles.
- [x] Add scroll offset handling per audit tab.
- [x] Add tests for tab switching, scrolling, empty states, and close behavior.

## Context Log

Read: `router.md`, `E32-Detail.md`, this task plus T001/T003 for scope boundaries, and board
`model.go`/`update.go`/`view.go`/`io.go`/`interfaces.go`/`detail.go`/`epic_panel.go`/`help.go` with
`release_docs_overlay_test.go` for the document-overlay pattern; `internal/data`
`audit.go`/`audit_register.go`/`audit_run.go`/`audit_finding.go` for the data set.

Edited:
- `internal/board/model.go` — added the `OverlayAudit` constant and extended `AuditState` with `AuditTab`
  and the per-tab `AuditOffsets` scroll map.
- `internal/board/audit_overlay.go` (new) — `auditTab` type + labels, the `selectAuditTab`/`scrollAuditTab`/
  `selectedAuditOffset` nav helpers, and `RenderAuditOverlay` with the tab strip and Prompt/Findings/Runs
  bodies (prompt reuses `renderReleaseDocBody`; findings/runs render read-only rows) plus empty states.
- `internal/board/update.go` — `A` opens the overlay (resets tab/offsets, triggers `loadAuditRegisterCmd`),
  `OverlayAudit` routes to the new read-only `handleAuditOverlay` (tab switch `[`/`]`/arrows/`h`/`l`,
  scroll, close).
- `internal/board/view.go` — render branch for `OverlayAudit`.
- `internal/board/audit_overlay_test.go` (new) — renderer (header/tabs/bodies/footer/empty states), `A`
  open + load wiring, tab switching + clamp, per-tab scroll preservation + pgup/pgdown, and esc/`q` close.

Scope note: the Findings tab renders read-only summary rows only; the grouped finding list with a cursor
and the finding detail drill-in belong to T003. Help/footer hint updates belong to T004.

Quality gates: `make build && make test` pass; `go vet ./internal/board` clean.
