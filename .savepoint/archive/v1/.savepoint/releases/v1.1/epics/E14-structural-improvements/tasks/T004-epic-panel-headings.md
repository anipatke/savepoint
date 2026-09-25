---
id: E14-structural-improvements/T004-epic-panel-headings
status: done
objective: Replace substring heading matching in epic_panel with exact heading match and configurable blocklist
depends_on: []
---

# T004: Replace Ad-hoc Heading Matching in epic_panel

## Context Files

- `internal/board/epic_panel.go:131-136` — substring heading detection in epicAuditBody

## Acceptance Criteria

- [x] Heading filtering uses exact match against a blocklist (not substring)
- [x] Blocklist is configurable via a named variable
- [x] Existing behavior preserved for all known audit headings
- [x] `go test ./...` passes with no regressions

## Implementation Plan

- [x] Identify current substring logic and blocklist
- [x] Replace with exact heading matching
- [x] Rename/restructure the blocklist variable
- [x] Update tests if applicable
- [x] Run `make build && make test`

## Context Log

- Files read: `internal/board/epic_panel.go`, `internal/board/epic_panel_test.go`.
- Files edited: `internal/board/epic_panel.go`, `internal/board/epic_panel_test.go`.
- Token estimate: ~6k.
- Quality gates: `go test ./internal/board` passed; `make build && make test` passed (`go test ./...` all packages).
- TUI note: could not press `p` from the agent environment; router already points at this task as current priority.
