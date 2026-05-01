## Target File

`.savepoint/Design.md`

## Replace

```md
**Layout:** single screen with a 5-column Kanban board and detail pane. Non-TTY output uses `src/tui/render/plain-table.ts`.

**Implementation modules:** see AGENTS.md Codebase Map (E06 and E07 epic rows).

**Keybindings:** arrow/vim navigation, enter advances, backspace retreats, r/R refreshes, a/A exits toward audit review when proposals exist, q quits.
```

## With

```md
**Layout:** Go Bubble Tea board with responsive columns: wide terminals (`>=120` cols) show an epic sidebar plus planned / in-progress / done columns; medium terminals show three columns; narrow terminals show the focused column only.

**Board components:** E04 adds Lip Gloss renderers for columns, task cards, epic sidebar/dropdown, release dropdown, task detail overlay, and help overlay. Selector, detail, and help overlays dim the current board and composite a centered panel over it.

**Implementation modules:** see AGENTS.md Codebase Map (`internal/board/*`, `internal/styles/*`, and `internal/data/*` rows).

**Keybindings:** h/left and l/right move columns; enter opens task detail or selects an overlay item; e opens the epic selector on narrow screens; r opens the release selector; ? opens help; esc/q close overlays; q/ctrl+c quit from the base board.
```

## Target File

`AGENTS.md`

## Insert After

```md
| `internal/data/errors.go`            | [E02-data-readers](.savepoint/releases/v1/epics/E02-data-readers/Design.md)                   | Shared data-reader boundary error sentinels                                                          |
```

## With

```md
| `internal/board/model.go`            | [E03-board-tui-core](.savepoint/releases/v1/epics/E03-board-tui-core/Design.md)                | Bubble Tea board state for grouped tasks, focused column/task, selected release/epic, responsive size, and active overlay |
| `internal/board/update.go`           | [E03-board-tui-core](.savepoint/releases/v1/epics/E03-board-tui-core/Design.md)                | Keyboard and window-size update handling for board navigation, overlays, selectors, and quit behavior |
| `internal/board/layout.go`           | [E03-board-tui-core](.savepoint/releases/v1/epics/E03-board-tui-core/Design.md)                | Responsive board geometry breakpoints and column/sidebar width calculation                            |
| `internal/board/view.go`             | [E03-board-tui-core](.savepoint/releases/v1/epics/E03-board-tui-core/Design.md)                | Board view composition, status bar, responsive columns, and overlay placement helpers                 |
| `internal/board/column.go`           | [E04-board-components](.savepoint/releases/v1/epics/E04-board-components/Design.md)            | Column renderer with status header, task count, empty state, and focused border styling               |
| `internal/board/card.go`             | [E04-board-components](.savepoint/releases/v1/epics/E04-board-components/Design.md)            | Task card renderer with phase glyphs, metadata dimming, truncation, and focus styling                 |
| `internal/board/epic_panel.go`       | [E04-board-components](.savepoint/releases/v1/epics/E04-board-components/Design.md)            | Wide epic sidebar and narrow epic selector dropdown renderers                                         |
| `internal/board/detail.go`           | [E04-board-components](.savepoint/releases/v1/epics/E04-board-components/Design.md)            | Task detail overlay renderer with task metadata, description, acceptance criteria, and phase labels   |
| `internal/board/release.go`          | [E04-board-components](.savepoint/releases/v1/epics/E04-board-components/Design.md)            | Release selector dropdown renderer and current-release cursor lookup                                  |
| `internal/board/help.go`             | [E04-board-components](.savepoint/releases/v1/epics/E04-board-components/Design.md)            | Keyboard shortcut help overlay renderer                                                               |
| `internal/styles/*.go`               | [E04-board-components](.savepoint/releases/v1/epics/E04-board-components/Design.md)            | Atari-Noir Lip Gloss palette and shared board, card, overlay, glyph, status, and tag styles          |
| `internal/board/*_test.go`           | [E03-board-tui-core](.savepoint/releases/v1/epics/E03-board-tui-core/Design.md) / [E04-board-components](.savepoint/releases/v1/epics/E04-board-components/Design.md) | Go board model, layout, update, view, component, selector, detail, and help tests                     |
```

## Target File

`.savepoint/releases/v1/epics/E04-board-components/Design.md`

## Insert After

```md
## Definition of Done

- Cards render with correct phase glyphs and colors.
- Focused card has orange accent border.
- Epic panel renders on wide screens; dropdown on narrow.
- Detail overlay opens with Enter, closes with Esc.
- Epic dropdown opens with `e`, release with `r`.
- Help overlay opens with `?`.
- All components render without wrapping or layout breaks.
- Component tests pass.
```

## With

```md
## Implemented As

- `internal/board/column.go` renders status columns with uppercase label, task count, empty state, focused border, and task cards.
- `internal/board/card.go` renders standalone task cards with build/test/audit glyphs, dim metadata, truncation, and focused styling.
- `internal/board/epic_panel.go` renders the wide epic sidebar and narrow epic dropdown; `e` opens the dropdown only below the wide breakpoint.
- `internal/board/release.go` renders the release dropdown; `r` opens it at all widths.
- `internal/board/detail.go` renders task detail panels with wrapped long text; `internal/board/view.go` dims the board and composites selector/detail/help overlays over the base board.
- `internal/board/help.go` renders the keyboard shortcut overlay.
- `internal/styles/styles.go` and `internal/styles/palette.go` add shared Lip Gloss styles for cards, overlays, metadata, glyphs, and dim text.

## Audit Deltas

- Planned `internal/board/dropdown.go` was not created; epic and release dropdowns are separate renderers.
- Task cards are integrated into column rendering; board columns now display card borders and phase glyphs.
- Selecting an epic or release updates `SelectedEpic` / `SelectedRelease` and regroups visible board tasks from `Model.AllTasks`.
- Selector/detail/help overlays keep the dimmed board visible behind the panel.
- Base board up/down navigation moves the focused task within the focused column.
```

## Quality Review

## Must Fix Before Close

No open findings.

## Must Fix Before Next Epic

No open findings.

## Carry Forward

- The current project architecture and Codebase Map still contain substantial TypeScript-era rows. This proposal adds E04 Go rows as a delta, but a later audit should reconcile the stale TypeScript implementation sections once the Go transition is complete.

## Already Fixed

- Board columns render `RenderCard`, so phase glyphs, card borders, truncation, and focused card styling are visible in the actual board.
- Epic and release selections regroup visible tasks from `Model.AllTasks`.
- Epic and release dropdowns use the same dimmed-base overlay composition as detail/help.
- Base board up/down navigation moves `FocusedTask` and clamps at column bounds.
- Detail rows, descriptions, and acceptance criteria wrap long text within overlay width.
- Focused tests passed during audit: `go test ./internal/board/...`.
- Full Go suite passed during audit: `go test ./...`.
- Help and detail overlays use the dimmed-base compositing path and close on `esc` / `q`.
