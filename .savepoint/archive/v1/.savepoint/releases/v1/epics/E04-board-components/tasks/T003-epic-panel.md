---
id: E04-board-components/T003-epic-panel
status: done
objective: "Render epic sidebar on wide screens, dropdown on narrow"
depends_on: [E04-board-components/T002-card]
---

# T003: Epic Panel

## Acceptance Criteria

- >=120 cols: left sidebar with epic list, active indicator.
- <120 cols: `e` key opens epic dropdown overlay.
- Active epic marked with indicator.
- Selecting an epic updates board tasks.

## Implementation Plan

- [x] Create `internal/board/epic_panel.go`.
- [x] Implement sidebar rendering (`RenderEpicSidebar`).
- [x] Implement dropdown overlay (`RenderEpicDropdown`).
- [x] Wire into Update loop (`e` key, overlay nav, enter/esc).
- [x] Write tests (`epic_panel_test.go`, 25 tests).
- [x] Run `go test`.

## Context Log

**Files read:** model.go, view.go, update.go, layout.go, card.go, column.go, board.go, styles/styles.go, styles/palette.go, card_test.go, view_test.go, update_test.go, model_test.go

**Estimated input tokens:** ~6000

**Notes:**
- Added `OverlayEpic` to `OverlayType` constants in model.go.
- Added `Epics []string` and `EpicCursor int` to Model struct; `NewModel` signature unchanged — callers set Epics directly.
- `RenderEpicSidebar` falls back to showing `selected` as sole entry when `Epics` is nil; shows `(none)` if both empty.
- `e` key only opens dropdown when `m.Width < breakpointWide` (120); sidebar is always visible on wide screens.
- Overlay blocks column nav (`l`/`h`) while open; `esc`/`q` both close it.
- `lipgloss.Place` centers dropdown on full terminal; board not visible behind (limitation of lipgloss — noted as drift).
- `defaultTermH = 24` added to view.go for zero-Height fallback.
- `renderEpicPanel` stub in view.go replaced with `RenderEpicSidebar` call.
- `fmt` import in view.go retained (still used by `renderStatus`).

**Quality gates:** `go build ./...` ✓  `go vet ./...` ✓  `go test ./...` ✓ (all pass)

## Drift Notes

- Drift: `internal/board/epic_panel.go` added, not yet in Codebase Map.
- Drift: Model.Epics and Model.EpicCursor fields added — Design.md data model section may need update.
- Drift: Board-visible-behind overlay not implemented (lipgloss limitation); Design.md "board visible behind" note is aspirational.
