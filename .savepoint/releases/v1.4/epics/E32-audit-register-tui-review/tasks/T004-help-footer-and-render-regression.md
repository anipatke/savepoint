---
id: E32-audit-register-tui-review/T004-help-footer-and-render-regression
status: done
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

- [x] Footer includes `A:audits` without making existing hints wrap poorly at supported widths.
- [x] Help overlay includes the Audit Register shortcut.
- [x] Existing `d:defects` and `D:docs` hints remain unchanged.
- [x] Render policy tests include the new audit overlay.
- [x] Existing board, release docs, defect, epic detail, and task detail tests pass.

## Implementation Plan

- [x] Add the audit shortcut to footer hints.
- [x] Add the audit shortcut to help rendering.
- [x] Add render policy coverage for the audit overlay and finding detail.
- [x] Add width-focused tests for footer readability.
- [x] Run the board package tests.

## Context Log

- Read: `view.go`, `view_test.go`, `help.go`, `help_test.go`, `render_policy_test.go`,
  `audit_overlay.go`, `audit_overlay_test.go`, `audit_detail.go`, `update.go`, `layout.go`.
- Footer (`view.go`): extracted the hint bar to a named `footerHints` constant (content out of
  logic) and added `A:audits`. The prior bar filled exactly 80 columns with two-space separators,
  so adding a 10-column hint would overflow and truncate `q:quit` at the 80-column three-column
  minimum. Switched the hint bar to single-space separators and abbreviated `p:priority`→`p:prio`,
  landing at width 78 with every hint (including the unchanged `d:defects`/`D:docs`) visible at 80.
  The overlay hint footers keep their own two-space style; only the denser board hint bar changed.
- Help (`help.go`): added `A: open audit register` immediately after the `D` docs row, matching the
  footer ordering.
- Render policy (`render_policy_test.go`): added the audit overlay (Findings tab) and finding detail
  to the no-background-escape case map, reusing the shared `sampleAuditSet` fixture.
- Tests: updated `TestView_containsFooterHints` to assert the `footerHints` constant; added
  `TestView_footerIncludesAuditHint`, `TestView_footerKeepsDefectAndDocHintsUnchanged`, and
  `TestView_footerHintsFitSupportedWidths` (80/100/120 stay one untruncated line ending in
  `q:quit`); added `TestRenderHelp_containsAuditShortcut`.
- Quality gates: `make build && make test` pass; `go vet ./internal/board/` and `gofmt -l` clean for
  the touched files (`release.go` is a pre-existing, untouched gofmt nit).
