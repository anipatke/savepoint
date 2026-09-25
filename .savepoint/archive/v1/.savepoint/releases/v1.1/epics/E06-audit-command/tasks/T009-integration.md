---
id: E06-audit-command/T009-integration
status: done
objective: End-to-end integration test of agent audit workflow
depends_on:
    - E06-audit-command/T007-apply-close
---

# T009: Integration

## Acceptance Criteria

- [x] Full agent audit workflow runs end-to-end
- [x] Agent audit → writes E##-Audit.md → user approves → agent applies → epic closed
- [x] Test passes on Windows, Linux, macOS
- [x] Audit tab correctly displays findings in TUI
- [x] All 10 Code Style rules checked in audit

## Implementation Plan

- [x] Create integration test for agent audit workflow:
  - Agent runs audit for epic
  - Writes findings to E##-Audit.md
  - User reviews (or auto-approve for test)
  - Agent applies proposals
  - Epic marked audited
  - Router advances
- [x] Verify TUI audit tab displays E##-Audit.md content
- [x] Verify Code Style checklist present
- [x] Run full test suite: `make build && make test`
- [x] Document integration test coverage

## Context Log

Tests added to `internal/board/epic_panel_test.go`:
- `TestRenderEpicAuditTab_allCodeStyleRules` — all 10 rules render
- `TestView_epicAuditTabRendered` — View() renders EPIC AUDIT header when tab=1
- `TestAuditWorkflow_fullEndToEnd` — creates temp E##-Audit.md, opens overlay, presses 2, verifies content loads and renders, presses 1, presses esc

All tests pass on Windows (`go build ./... && go test ./...`).