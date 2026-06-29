---
id: E32-audit-register-tui-review/T004-help-footer-and-render-regression
status: planned
objective: Integrate Audit Register hints into board help, footer, and render policy coverage.
depends_on:
  - E32-audit-register-tui-review/T003-finding-list-and-detail
complexity_tier: medium
complexity_reason: Help and render updates must avoid regressions across existing overlays.
---

# T004: Help, Footer, and Render Regression

## Problem

The new Audit Register section needs discoverable shortcuts and regression coverage so existing board interactions stay stable.

## Context Files

- `internal/board/help.go`
- `internal/board/help_test.go`
- `internal/board/view.go`
- `internal/board/view_test.go`
- `internal/board/render_policy_test.go`
- `internal/board/audit_overlay_test.go`

## Acceptance Criteria

- [ ] Footer includes `A:audits` without making existing hints wrap poorly at supported widths.
- [ ] Help overlay includes the Audit Register shortcut.
- [ ] Existing `d:defects` and `D:docs` hints remain unchanged.
- [ ] Render policy tests include the new audit overlay.
- [ ] Existing board, release docs, defect, epic detail, and task detail tests pass.

## Implementation Plan

- [ ] Add the audit shortcut to footer hints.
- [ ] Add the audit shortcut to help rendering.
- [ ] Add render policy coverage for the audit overlay and finding detail.
- [ ] Add width-focused tests for footer readability.
- [ ] Run the board package tests.

## Context Log

Pending.
