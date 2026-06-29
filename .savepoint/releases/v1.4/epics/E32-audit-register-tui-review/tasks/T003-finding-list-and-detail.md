---
id: E32-audit-register-tui-review/T003-finding-list-and-detail
status: planned
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

- [ ] Findings render grouped by status and severity with ID, title, confidence, and linked work item.
- [ ] Findings cursor moves with `up`/`down` and `k`/`j`.
- [ ] `enter` on a finding opens read-only finding detail.
- [ ] Detail shows ID, title, status, severity, confidence, links, locations, proof needed, first seen, last seen, and markdown body sections.
- [ ] `esc` and `q` return from detail to the findings tab.
- [ ] No key mutates finding status or writes audit files.

## Implementation Plan

- [ ] Render finding rows with deterministic grouping.
- [ ] Add selected finding cursor state and clamping.
- [ ] Add read-only finding detail renderer.
- [ ] Wire enter/escape behavior between list and detail.
- [ ] Add tests for row content, cursor behavior, detail content, and no-op mutation keys.

## Context Log

Pending.
