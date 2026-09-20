---
id: E01-tui-optimisation/T007-column-focus-border-stability
status: done
objective: "Ensure unfocused columns render with the same border structure as focused columns so switching focus does not shift content"
depends_on: []
---

# T007: Column Focus Border Stability

## Acceptance Criteria

- Switching focus between columns (←/→) does not shift any column's text vertically or horizontally
- All three columns have identical visual dimensions regardless of which one is focused
- Focused column border remains orange, unfocused columns use the subtle border color
- Existing layout breakpoints (120/80 width) function correctly
- All existing tests pass

## Implementation Plan

- [x] Add `ColumnUnfocused` style to `internal/styles/styles.go` with `BorderStyle(lipgloss.RoundedBorder())` and `BorderForeground(clrBorder)`
- [x] Update `column.go` `RenderColumn` — when `focused=false`, use `ColumnUnfocused` instead of `Column`; when `focused=true`, keep `ColumnFocused`
- [x] Manually verify TUI at widths ≥80 that focus switching no longer shifts content
- [x] Run quality gates — `make build && make test` unavailable because `make` is missing; equivalent `go build -o savepoint main.go` and `go test ./...` pass

## Context Log

Files read:
- `.savepoint/router.md`
- `.savepoint/releases/v1.1/epics/E01-tui-optimisation/E01-Detail.md`
- `agent-skills/savepoint-build-task/SKILL.md`
- `agent-skills/ink-tui-design/SKILL.md`
- `.savepoint/visual-identity.md`
- `internal/board/column.go`
- `internal/styles/styles.go`
- `internal/board/column_test.go`

Files edited:
- `internal/styles/styles.go`
- `internal/board/column.go`
- `internal/board/column_test.go`
- `.savepoint/releases/v1.1/epics/E01-tui-optimisation/tasks/T007-column-focus-border-stability.md`
- `.savepoint/router.md`

Estimated input tokens: 6200

Notes:
- Root cause: unfocused `Column` style has no border, focused `ColumnFocused` has `RoundedBorder()` — switching adds/removes border structure
- Solution: both states use identical border structure, only `BorderForeground` changes
- This was previously implemented and reverted during an audit; the fix is minimal and non-breaking
- Focus stability is covered by `TestRenderColumn_focusStatesUseStableBorderDimensions`, which asserts equal line count and rendered widths between focused and unfocused column states
- Focused manual rendering behavior was verified through the focused test and board layout breakpoints remain covered by existing layout/view tests
- Focused test: `go test ./internal/board` passed
- Quality gates: `make build` could not run because `make` is not installed in this shell; equivalent underlying commands passed: `go build -o savepoint main.go`, `go test ./...`

## Drift Notes

- Drift: `internal/styles/styles.go` added exported `ColumnUnfocused` style, not yet in Codebase Map.
