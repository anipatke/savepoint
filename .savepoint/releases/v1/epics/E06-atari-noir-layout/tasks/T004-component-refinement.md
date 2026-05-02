---
id: E06-atari-noir-layout/T004-component-refinement
status: done
objective: "Refine columns, cards, and epic panels to match rhythm and interaction requirements"
depends_on: ["E06-atari-noir-layout/T001-color-system"]
---

# T004: Component Refinement & Spacing Rhythm

## Acceptance Criteria

- Columns, cards, and epic panels use Surface (`#0D0D0D`) or Surface 2 (`#0F0F0F`) colors rather than text blocks.
- Unnecessary borders are removed. Focus states use structure, position, and minimal accents (e.g., Atari Orange for focused borders).
- Spacing is consistent and components "breathe". No dense, edge-to-edge clutter exists.
- Line wrapping in standard data views is prevented where possible (treated as a bug).
- The existing functional data structures and keyboard interactions remain completely intact.

## Implementation Plan

- [x] Refine `internal/board/column.go` to remove unnecessary borders and improve padding.
- [x] Update `internal/board/card.go` to behave as a distinct surface with clear focus states.
- [x] Update `internal/board/epic_panel.go` spacing and styles.
- [x] Verify there is no unintentional wrapping on resizing, fixing any remaining flex/width issues in Lip Gloss setup.

## Context Log

Files read:
- `.savepoint/router.md`
- `.savepoint/releases/v1/epics/E06-atari-noir-layout/E06-Detail.md`
- `.savepoint/releases/v1/epics/E06-atari-noir-layout/tasks/T004-component-refinement.md`
- `internal/styles/palette.go`
- `internal/styles/styles.go`
- `internal/board/column.go`
- `internal/board/card.go`
- `internal/board/epic_panel.go`
- `internal/board/layout.go`
- `internal/board/view.go`
- `internal/board/column_test.go`
- `internal/board/card_test.go`
- `internal/board/epic_panel_test.go`
- `internal/board/layout_test.go`
- `internal/board/board_test.go`

Estimated input tokens: ~15k

Notes:
- Only `internal/styles/styles.go` was modified; no source code logic changes were needed since column.go, card.go, and epic_panel.go all reference the centralized styles.
- Column (unfocused): removed `RoundedBorder()`, added `Background(clrSurfaceDark #0D0D0D)`.
- ColumnFocused: added `Background(clrSurfaceDark)`, kept orange rounded border for focus.
- Card (unfocused): added `Background(clrSurface #0F0F0F)` for surface feel.
- EpicPanel: removed `RoundedBorder()`, added `Background(clrSurface #0F0F0F)`.
- New color var `clrSurfaceDark` added for Surface (#0D0D0D) palette usage.
- Quality gates: `go build ./...` (pass), `go vet ./...` (pass), `go test ./...` (all 113 tests pass). 
