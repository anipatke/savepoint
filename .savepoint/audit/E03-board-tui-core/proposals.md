# Audit Proposals — E03-board-tui-core

## Target File

`.savepoint/Design.md`

## Replace

```md
## 8. TUI

**Theming:** Atari-Noir is the default theme. **For full design tokens, palette, and rendering rules, see `.savepoint/visual-identity.md`** (loaded conditionally for TUI tasks). Live values in `config.yml` `theme:` section.

Acknowledged terminal limits: fonts, scanlines, glows, letter-spacing, mouse-driven motion don't translate. Lean on color discipline + box-drawing geometry + uppercase headings.

**Render fallbacks:** 256-color → 16-color hard-coded → `NO_COLOR=1` monochrome with glyphs → non-TTY plain table.

**Layout:** single screen with a 5-column Kanban board and detail pane. Non-TTY output uses `src/tui/render/plain-table.ts`.

**Implementation modules:** see AGENTS.md Codebase Map (E06 and E07 epic rows).

**Keybindings:** arrow/vim navigation, enter advances, backspace retreats, r/R refreshes, a/A exits toward audit review when proposals exist, q quits.
```

## With

```md
## 8. TUI

**Theming:** Atari-Noir is the default theme. **For full design tokens, palette, and rendering rules, see `.savepoint/visual-identity.md`** (loaded conditionally for TUI tasks). Live values in `config.yml` `theme:` section.

Acknowledged terminal limits: fonts, scanlines, glows, letter-spacing, mouse-driven motion don't translate. Lean on color discipline + box-drawing geometry + uppercase headings.

**Go board core:** established in epic `E03-board-tui-core` (2026-05-01). `internal/board` owns the Bubble Tea model, update loop, responsive layout calculation, and Lip Gloss view composition. `internal/styles` owns the Atari-Noir terminal palette and reusable Lip Gloss styles with truecolor, 256-color, and 16-color fallbacks.

**Layout:** the Go board core currently renders a responsive 3-column grid. At `>=120` columns it includes a 28-cell epic panel plus 3 task columns; at `80-119` columns it renders 3 task columns; below `80` columns it renders only the focused column and supports left/right or h/l column switching.

**Render fallbacks:** truecolor → 256-color → 16-color are encoded in Go styles. A future non-TTY/plain renderer remains outside this epic.

**Implementation modules:** see AGENTS.md Codebase Map.

**Keybindings:** `q` and `ctrl+c` quit. Left/right arrows and h/l move the focused column in narrow layouts.
```

## Target File

`AGENTS.md`

## Insert After

```md
| `internal/data/errors.go`            | [E02-data-readers](.savepoint/releases/v1/epics/E02-data-readers/Design.md)                   | Shared data-reader boundary error sentinels                                                          |
```

## With

```md
| `internal/board/board.go`            | [E03-board-tui-core](.savepoint/releases/v1/epics/E03-board-tui-core/Design.md)               | Bubble Tea board entrypoint and program launch                                                        |
| `internal/board/model.go`            | [E03-board-tui-core](.savepoint/releases/v1/epics/E03-board-tui-core/Design.md)               | Board model state, constructor, task grouping, and initialization                                     |
| `internal/board/update.go`           | [E03-board-tui-core](.savepoint/releases/v1/epics/E03-board-tui-core/Design.md)               | Board key handling, quit behavior, column focus navigation, and terminal resize updates               |
| `internal/board/view.go`             | [E03-board-tui-core](.savepoint/releases/v1/epics/E03-board-tui-core/Design.md)               | Lip Gloss view composition for header, responsive board columns, epic panel, tasks, and status bar    |
| `internal/board/layout.go`           | [E03-board-tui-core](.savepoint/releases/v1/epics/E03-board-tui-core/Design.md)               | Responsive board breakpoint and column width calculations                                            |
| `internal/styles/palette.go`         | [E03-board-tui-core](.savepoint/releases/v1/epics/E03-board-tui-core/Design.md)               | Atari-Noir terminal color constants and 256-color/16-color fallbacks                                 |
| `internal/styles/styles.go`          | [E03-board-tui-core](.savepoint/releases/v1/epics/E03-board-tui-core/Design.md)               | Reusable Lip Gloss styles for board chrome, focused state, task rows, status, tags, and epic panel   |
| `internal/board/*_test.go`           | [E03-board-tui-core](.savepoint/releases/v1/epics/E03-board-tui-core/Design.md)               | Unit tests for board model, update loop, layout breakpoints, and rendered output                     |
```

## Target File

`.savepoint/releases/v1/epics/E03-board-tui-core/Design.md`

## Replace

```md
type: epic-design
status: planned
```

## With

```md
type: epic-design
status: implemented
```

## Target File

`.savepoint/releases/v1/epics/E03-board-tui-core/Design.md`

## Insert After

```md
## Components and files

| Path | Purpose |
|------|---------|
| `internal/board/model.go` | Model struct and Init |
| `internal/board/update.go` | Update loop, message handling |
| `internal/board/view.go` | View composition |
| `internal/styles/palette.go` | Atari-Noir hex colors |
| `internal/styles/styles.go` | Lip Gloss style definitions |
| `internal/board/layout.go` | Width/height calculations |
```

## With

```md

## Implemented as

- `internal/board/model.go` defines exported board state, groups tasks by `data.ColumnType`, defaults empty task columns to `planned`, and initializes with no startup command.
- `internal/board/update.go` handles `q`, `ctrl+c`, terminal resize messages, and left/right or h/l focused-column navigation.
- `internal/board/layout.go` implements the requested 120/80 column breakpoints and returns explicit layout metadata for view rendering.
- `internal/board/view.go` composes the header, responsive task columns, optional epic panel, and status bar using Lip Gloss.
- `internal/styles/palette.go` and `internal/styles/styles.go` define Atari-Noir truecolor values plus 256-color and 16-color fallbacks.

## Audit deltas

- The implemented board core is wired into `internal/board/board.go`; `Run()` now launches the exported board model rather than the pre-epic placeholder screen.
- The layout currently uses rounded borders, while the visual identity prefers single-line quiet borders by default. This is acceptable if intentional, but should be normalized before the visual system hardens.
- Layout width coverage now lives on `CalculateLayout`; the dead `columnContentWidth()` helper was removed.
```

## Quality Review

## Must Fix Before Close

None.

## Must Fix Before Next Epic

- `internal/board/layout.go:28` accepts height but ignores it. The epic originally calls for width/height calculations and no layout break at any width. Either incorporate height into visible task capacity or narrow the documented contract before components depend on it.
- `internal/board/view.go:65` renders task labels without truncation or width-aware clipping. Long IDs/titles can wrap and break the grid, conflicting with the visual guardrail that accidental wrapping is a bug.

## Carry Forward

- `internal/styles/styles.go` uses rounded borders for columns and the epic panel. The visual identity says single-line borders are the default and rounded borders are not called out as a signature style; decide whether rounded borders are an intentional product choice in E04.
- The Go implementation has replaced the older TypeScript architecture in active epics. Continue reconciling `.savepoint/Design.md` and AGENTS.md away from stale `src/` rows as each Go epic audits.

## Already Fixed

- Board state, resize handling, responsive breakpoints, Atari-Noir palette constants, color fallbacks, focused-column rendering, and focused-column navigation have focused tests.
- `internal/board/board.go` now launches `newProgramModel()`, which returns the exported `Model` created by `NewModel(...)`; the private placeholder model was removed.
- `internal/board/board_test.go` verifies the program model renders board columns and does not render the old welcome screen.
- The dead `columnContentWidth()` helper and its stale tests were removed; layout width coverage remains in `internal/board/layout_test.go`.
- `go test ./...` and `go build ./...` passed during audit.
