---
type: epic-design
status: audited
---

# Epic E04: Epic Navigation

## Purpose

Make the epic sidebar an active, focusable UI component. The sidebar currently lists epics as static text — the only way to switch epics is through the `E` key dropdown popup. This epic gives the sidebar parity with task columns: ↑/↓ navigation, focus highlighting, and an Enter-triggered "Epic Detail" overlay.

## Definition of Done

- Epic sidebar is focusable — left arrow from `Planned` column moves focus to sidebar
- In sidebar focus: ↑/↓ navigates through epics, right arrow returns to column view
- Focused epic filters tasks in all three columns (same filtering as current `SelectedEpic`)
- A focused epic entry in the sidebar is visually distinct from unfocused ones
- Enter on a focused epic opens an "Epic Detail" overlay showing the epic's `E##-Detail.md` content
- Overlay supports scrolling (↑/↓, PgUp/PgDown) identical to task detail overlay
- `E` key dropdown behavior is preserved unchanged
- All above works only when epic panel is visible (≥120 terminal width)

## Components and files

| Path | Purpose |
|------|---------|
| `internal/board/model.go` | Add `EpicPanelFocus`, `EpicPanelCursor`, `EpicDetailOffset`, `EpicDetailContent` state |
| `internal/board/update.go` | Handle left arrow to enter panel, ↑/↓ within, Enter to open detail, right arrow to exit |
| `internal/board/view.go` | Pass focus/cursor to `RenderEpicSidebar`, render epic detail overlay |
| `internal/board/epic_panel.go` | Update `RenderEpicSidebar` for focus/cursor, add `RenderEpicDetail` |
| `internal/board/detail.go` | Reuse `visibleDetailLines`, `WrapText`, scroll helpers for epic detail |
| `internal/board/board.go` | Load epic status from detail files in `loadBoardData` |
| `internal/board/watch.go` | Carry `epicStatuses` in `reloadMsg` for dynamic reloads |
| `internal/data/parser.go` | Reuse `ParseFrontmatter` to extract status YAML field |

## Architectural notes

- `EpicPanelFocus bool` is a lightweight flag replacing the need for a new column type
- The epic detail overlay content is read from `E##-Detail.md` on the filesystem when Enter is pressed (not pre-loaded)
- Detail file path is deterministic: `{root}/releases/{release}/epics/{epic-slug}/{shortID}-Detail.md`
- If the detail file is missing, the overlay shows a "(no detail available)" message
- Column selection, task selection, and existing overlay behavior are completely unchanged

## Implemented as

- `internal/board/model.go` adds `EpicPanelFocus`, `EpicPanelCursor`, `EpicDetailOffset`, `EpicDetailContent`, and `EpicStatus` model state.
- `internal/board/update.go` handles global keys before epic-panel routing, focuses the panel from the Planned column on wide layouts, changes the selected epic during panel cursor movement, writes router release/epic state, and opens the epic detail overlay on Enter.
- `internal/board/epic_panel.go` renders the purple-accented epic sidebar focus state, purple epic detail overlay, markdown detail body, and side-panel-only status glyph prefixes.
- `internal/board/board.go` loads epic status frontmatter during board-data loading; `internal/board/watch.go` carries that status map through reloads.
- Epic navigation deliberately uses `VibePurple` (`#B1A1DF`) for focused epic panel borders, focused epic labels, epic detail overlays, and epic/audit accents, while task-column focus remains Atari Orange.
- Implementation deviation: T001 originally said Enter in epic-panel focus selected the focused epic. The final behavior from T002 is that up/down selects and filters immediately, while Enter opens the epic detail overlay.
