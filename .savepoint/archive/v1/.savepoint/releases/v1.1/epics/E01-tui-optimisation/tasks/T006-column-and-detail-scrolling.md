---
id: E01-tui-optimisation/T006-column-and-detail-scrolling
status: done
objective: "Add virtual viewport scrolling to board columns and detail overlay so content adapts to terminal height"
depends_on:
  - E01-tui-optimisation/T001-border-resize-fix
---

# T006: Column and Detail Scroll Viewport

## Acceptance Criteria

- Each column shows only 4-5 task cards by default (adapts to terminal height) and renders `↑ N above` / `↓ N more` indicators when content extends beyond the viewport
- The column viewport auto-scrolls to keep the focused task visible when navigating with j/k/up/down
- PageUp/PageDown keys scroll a full page of cards in the focused column
- The detail overlay is capped at ~70% of terminal height and shows the same `↑ N above` / `↓ N more` indicators when content exceeds available space
- j/k/up/down keys scroll the detail overlay content when it is open
- All scroll indicators use dim/subtle styling matching the Atari Noir aesthetic
- Existing layout breakpoints (120/80 width) still function correctly with height-aware layout
- No visual corruption on terminal resize — viewport recalculates available height
- All existing tests pass; new tests cover viewport offset calculations, indicator rendering, and auto-scroll behavior

## Implementation Plan

- [x] Add `ColumnOffsets map[data.ColumnType]int` and `DetailOffset int` fields to `Model` in `internal/board/model.go`
- [x] Extend `Layout` struct in `internal/board/layout.go` with `ContentHeight int` — computed from terminal height minus header (~3), footer (~3), board border (2), column border (2)
- [x] Update `CalculateLayout(width, height)` to accept height and populate `ContentHeight`
- [x] Refactor `RenderColumn` in `internal/board/column.go` to accept `maxHeight` param, compute visible task slice from offset, and append scroll indicators (`↑ N above` / `↓ N more`) styled with dim/subtle text
- [x] Add `ScrollIndicator` style to `internal/styles/styles.go` (dim, faint)
- [x] Refactor `RenderDetail` in `internal/board/detail.go` to accept `maxHeight` param (capped at 70% terminal height), count total lines, and render only the visible slice from `DetailOffset` with scroll indicators at edges
- [x] Update `view.go` — pass `layout.ContentHeight` through render chain to `renderColumn` and `RenderDetail`
- [x] Update `update.go` — auto-scroll `ColumnOffsets` when focused task moves beyond visible range; handle PageUp/PageDown keybindings in both normal and detail overlay modes; add j/k scroll handling in detail overlay
- [x] Add tests in `internal/board/column_test.go` for scroll indicator rendering and viewport slicing
- [x] Add tests in `internal/board/layout_test.go` for height-aware `CalculateLayout` including ContentHeight computation and minimum height floors
- [x] Run quality gates — `make` unavailable locally; equivalent `go build ./...` and `go test ./...` pass

## Context Log

Files read:
- `internal/board/model.go`
- `internal/board/layout.go`
- `internal/board/view.go`
- `internal/board/column.go`
- `internal/board/detail.go`
- `internal/board/card.go`
- `internal/board/update.go`
- `internal/styles/styles.go`
- `internal/board/column_test.go`
- `internal/board/detail_test.go`
- `internal/board/layout_test.go`
- `internal/board/update_test.go`
- `.savepoint/router.md`
- `.savepoint/releases/v1.1/epics/E01-tui-optimisation/E01-Detail.md`
- `.savepoint/releases/v1.1/epics/E01-tui-optimisation/tasks/T006-column-and-detail-scrolling.md`
- `agent-skills/savepoint-build-task/SKILL.md`
- `agent-skills/ink-tui-design/SKILL.md`
- `.savepoint/visual-identity.md`

Estimated input tokens: 12500

Notes:
- Pure rendering concern — no data model or card rendering changes needed
- `card.go` read only — column decides which cards to show
- Scroll uses virtual viewport (offset tracking) not literal scrollbar glyphs
- Auto-scroll-to-focus means explicit scrolling is rarely needed by user
- Focused board tests: `go test ./internal/board` passed
- Required quality gate `make build && make test` could not run because `make` is not installed in this PowerShell environment
- Equivalent gates passed: `go build ./...` and `go test ./...`
