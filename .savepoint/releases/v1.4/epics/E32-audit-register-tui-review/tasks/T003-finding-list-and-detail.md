---
id: E32-audit-register-tui-review/T003-finding-list-and-detail
status: done
objective: Add finding list navigation and read-only finding detail.
depends_on:
    - E32-audit-register-tui-review/T002-audit-overlay-tabs
complexity_tier: medium
complexity_reason: Detail navigation adds a nested overlay state and rendering branch.
---

# T003: Finding List and Detail

## Problem

The Findings tab needs quick scanning and a detail view for proof, locations, linked work items, and body notes.

## Context Files

- `internal/board/audit_overlay.go`
- `internal/board/audit_overlay_test.go`
- `internal/board/audit_detail.go`
- `internal/board/audit_detail_test.go`
- `internal/board/update.go`
- `internal/board/view.go`
- `internal/board/defect_detail.go`

## Acceptance Criteria

- [x] Findings render grouped by status and severity with ID, title, confidence, and linked work item.
- [x] Findings cursor moves with `up`/`down` and `k`/`j`.
- [x] `enter` on a finding opens read-only finding detail.
- [x] Detail shows ID, title, status, severity, confidence, links, locations, proof needed, first seen, last seen, and markdown body sections.
- [x] `esc` and `q` return from detail to the findings tab.
- [x] No key mutates finding status or writes audit files.

## Implementation Plan

- [x] Render finding rows with deterministic grouping.
- [x] Add selected finding cursor state and clamping.
- [x] Add read-only finding detail renderer.
- [x] Wire enter/escape behavior between list and detail.
- [x] Add tests for row content, cursor behavior, detail content, and no-op mutation keys.

## Context Log

- Read: `audit_overlay.go`, `audit_overlay_test.go`, `update.go`, `view.go`, `model.go`,
  `defect_detail.go`, `defect_overlay.go`, `detail.go`, `epic_panel.go`, `util.go`,
  and `internal/data/audit_finding.go` / `audit.go` / `audit_register.go` for the finding model and sort order.
- Rendering (`audit_overlay.go`): replaced the flat findings body with `auditFindingLayout`,
  a single grouping pass that emits status-header sections (findings arrive pre-sorted by
  status → severity → ID) plus a per-line finding index. Rows show severity tag, ID, title,
  `conf:<confidence>`, and the linked work item; the cursor row is highlighted with the shared
  `releaseActiveMarker`. `findingStatusLabel` derives headers from the status constants so there
  is one source of truth. `RenderAuditOverlay` gained a `findingCursor` parameter.
- Detail (`audit_detail.go`, new): `RenderFindingDetail` shows ID, title, status, severity,
  confidence, work-item/related-record links, locations, proof needed, first/last seen, and the
  markdown body via `renderReleaseDocBody`; absent links/locations/proof rows are omitted.
- State/keys (`model.go`, `update.go`, `view.go`): added `OverlayFindingDetail`, `FindingCursor`,
  and `FindingDetailOffset`. On the Findings tab up/down/k/j move the cursor (clamped, with
  follow-scroll via `ensureFindingCursorVisible`); enter opens the detail; the detail scrolls and
  returns to the Findings tab on esc/q. Cursor is reset on `A` and clamped on audit reload. The
  finding paths are read-only and emit no write commands.
- Tests (`audit_detail_test.go`, new): grouping/ordering, cursor highlight, cursor move + clamp,
  empty-list no-ops, enter→detail, esc/q→Findings tab, detail scroll clamp, detail content,
  absent-link omission, and a read-only guarantee that space/backspace/p/a never mutate a finding
  or emit a command. Updated existing `audit_overlay_test.go` call sites for the new signature.
- Quality gates: `make build && make test` pass; `gofmt` and `go vet ./internal/board/` clean.
