# Audit Proposals: v1.1 E04 Epic Navigation

## Target File

`.savepoint/Design.md`

## Replace

```yaml
last_audited: E01-tui-optimisation (2026-05-02)
```

## With

```yaml
last_audited: E04-epic-navigation (2026-05-02)
```

## Target File

`.savepoint/Design.md`

## Replace

```md
- **Board command** (`savepoint board`) reads project state, renders the Atari-Noir TUI board, supports release/epic filtering, detail overlays, task status transitions with mtime-guarded writes, release/epic-scoped router priority markers, fsnotify-based task auto-refresh (epic E06), header Next Activity display, height-aware column/detail viewport scrolling, and stable focused/unfocused column border geometry (v1.1 E01).
```

## With

```md
- **Board command** (`savepoint board`) reads project state, renders the Atari-Noir TUI board, supports release/epic filtering, detail overlays, task status transitions with mtime-guarded writes, release/epic-scoped router priority markers, fsnotify-based task auto-refresh (epic E06), header Next Activity display, height-aware column/detail viewport scrolling, stable focused/unfocused column border geometry (v1.1 E01), and a focusable wide-screen epic sidebar with purple epic focus, epic detail overlays, and status glyphs loaded from epic detail frontmatter (v1.1 E04).
```

## Target File

`.savepoint/Design.md`

## Replace

```md
**Layout:** single screen with a 3-column task board (`planned`, `in_progress`, `done`), optional epic sidebar on wide terminals, centered overlays for release/epic/help/task detail, static Atari-Noir header/footer, full-width dividers, uniform black TUI backgrounds, and navigation hints. The header can show a compact right-aligned Next Activity value from router state. Columns and detail overlays use height-aware viewport slicing with subtle above/more scroll indicators. Focused and unfocused columns preserve the same rounded-border geometry so focus changes do not shift content. Non-TTY output remains a plain table fallback.
```

## With

```md
**Layout:** single screen with a 3-column task board (`planned`, `in_progress`, `done`), optional epic sidebar on wide terminals, centered overlays for release/epic/help/task/epic-detail views, static Atari-Noir header/footer, full-width dividers, uniform black TUI backgrounds, and navigation hints. The header can show a compact right-aligned Next Activity value from router state. Columns and detail overlays use height-aware viewport slicing with subtle above/more scroll indicators. Focused and unfocused columns preserve the same rounded-border geometry so focus changes do not shift content. On terminals at least 120 columns wide, the epic sidebar is focusable from the Planned column; it uses the purple epic accent for focused panel borders, focused epic labels, and epic detail overlays while task-column focus remains orange. Non-TTY output remains a plain table fallback.
```

## Target File

`.savepoint/Design.md`

## Replace

```md
**Board persistence and refresh:** task status transitions write canonical task frontmatter through `internal/data.WriteTaskStatus` with mtime conflict checks. The board treats `Model.Root` as the `.savepoint` directory, watches `.savepoint/releases/` recursively with fsnotify, adds watches for newly-created release/epic/task directories, and reloads task plus release/epic index data after debounced file changes. Router priority markers match release + epic + task, not only the short `T###` value; completed cards render with the orange build glyph even if they previously matched router priority.
```

## With

```md
**Board persistence and refresh:** task status transitions write canonical task frontmatter through `internal/data.WriteTaskStatus` with mtime conflict checks. The board treats `Model.Root` as the `.savepoint` directory, watches `.savepoint/releases/` recursively with fsnotify, adds watches for newly-created release/epic/task directories, and reloads task plus release/epic index data plus epic status metadata after debounced file changes. Router priority markers match release + epic + task, not only the short `T###` value; completed cards render with the orange build glyph even if they previously matched router priority. Epic status glyphs are cached from each epic's `E##-Detail.md` frontmatter and shown in the wide epic sidebar only.
```

## Target File

`AGENTS.md`

## Replace

```md
| `internal/board/`                    | TUI board models, layout, rendering, overlays, task transitions, router priority markers, and fsnotify refresh |
| `internal/data/`                     | Task/router/config models, frontmatter parsing, checklist state parsing, mtime-guarded writes, discovery, and generic file readers |
| `internal/styles/`                   | Atari-Noir palette constants, terminal color fallbacks, shared TUI styles, stable column border styles, scroll indicators, semantic glyph/tag styles, and footer/header styling |
```

## With

```md
| `internal/board/`                    | TUI board models, layout, rendering, overlays, focusable epic sidebar navigation, epic detail overlays, epic status glyph loading, task transitions, router priority markers, and fsnotify refresh |
| `internal/data/`                     | Task/router/config models, frontmatter parsing, checklist state parsing, mtime-guarded writes, discovery, and generic file readers |
| `internal/styles/`                   | Atari-Noir palette constants, terminal color fallbacks, shared TUI styles, stable column border styles, scroll indicators, purple epic navigation/detail styles, semantic glyph/tag styles, and footer/header styling |
```

## Target File

`.savepoint/releases/v1.1/epics/E04-epic-navigation/E04-Detail.md`

## Replace

```yaml
status: planned
```

## With

```yaml
status: audited
```

## Target File

`.savepoint/releases/v1.1/epics/E04-epic-navigation/E04-Detail.md`

## Replace

```md
## Architectural notes

- `EpicPanelFocus bool` is a lightweight flag replacing the need for a new column type
- The epic detail overlay content is read from `E##-Detail.md` on the filesystem when Enter is pressed (not pre-loaded)
- Detail file path is deterministic: `{root}/releases/{release}/epics/{epic-slug}/{shortID}-Detail.md`
- If the detail file is missing, the overlay shows a "(no detail available)" message
- Column selection, task selection, and existing overlay behavior are completely unchanged
```

## With

```md
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
```

## Quality Review

## Must Fix Before Close

None.

## Carry Forward

- Accepted design decision: `RenderEpicDetail` intentionally omits the components/files table from the overlay for readability. The epic detail overlay shows the narrative sections and keeps implementation file mapping in the source `E##-Detail.md`.
- `readEpicDetailFile` falls back to `E##-Design.md` and then any `E##-*.md`, which is more permissive than the epic design's deterministic path. The fallback is harmless and user-friendly, but it should be documented if kept as product behavior.
- `epicPanelPageSize()` was added for eventual PgUp/PgDown support but is not currently used. This is minor dead code carried from the task plan and can be removed or wired in during a future sidebar paging task.

## Already Fixed

- Global `q`, `e`, `r`, and `?` keys are processed before epic-panel-specific key handling, so panel focus does not trap quit, release, epic dropdown, or help commands.
- Focused epic sidebar rendering uses purple focused border/title/item styles, making the colour-scheme shift from task orange to epic purple visible in the TUI.
- Focused board tests, `go build -o savepoint main.go`, and `go test ./...` pass.
