# E06 Atari-Noir Layout Audit Proposals

## Target File

`.savepoint/Design.md`

## Replace

```md
- **Board command** (`savepoint board`) reads project, non-TTY fallback, Ink TUI, transition gates, mtime writes, audit signaling (epic E06).
```

## With

```md
- **Board command** (`savepoint board`) reads project state, renders the Atari-Noir TUI board, supports release/epic filtering, detail overlays, task status transitions with mtime-guarded writes, router priority markers, and fsnotify-based task auto-refresh (epic E06).
```

## Target File

`.savepoint/Design.md`

## Replace

```md
**Layout:** single screen with a 5-column Kanban board and detail pane. Non-TTY output uses `src/tui/render/plain-table.ts`.

**Implementation modules:** see AGENTS.md Codebase Map (E06 and E07 epic rows).
```

## With

```md
**Layout:** single screen with a 3-column task board (`planned`, `in_progress`, `done`), optional epic sidebar on wide terminals, centered overlays for release/epic/help/task detail, static Atari-Noir header/footer, and navigation hints. Non-TTY output remains a plain table fallback.

**Board persistence and refresh:** task status transitions write canonical task frontmatter through `internal/data.WriteTaskStatus` with mtime conflict checks. The board watches `.savepoint/releases/` recursively with fsnotify and reloads task data after debounced file changes.

**Implementation modules:** see AGENTS.md Codebase Map.
```

## Target File

`AGENTS.md`

## Replace

```md
| `internal/board/`                    | TUI board components, models, layouts, transitions, and rendering logic                              |
| `internal/data/`                     | Task data models, frontmatter parsing, project configuration, routing, and generic file readers      |
| `internal/styles/`                   | Shared visual design system, TUI styling, and palettes                                               |
```

## With

```md
| `internal/board/`                    | TUI board models, layout, rendering, overlays, task transitions, router priority markers, and fsnotify refresh |
| `internal/data/`                     | Task/router/config models, frontmatter parsing, checklist state parsing, mtime-guarded writes, discovery, and generic file readers |
| `internal/styles/`                   | Atari-Noir palette constants, terminal color fallbacks, shared TUI styles, semantic glyph/tag styles, and footer/header styling |
```

## Target File

`.savepoint/releases/v1/epics/E06-atari-noir-layout/E06-Detail.md`

## Replace

```md
# Epic E07: Atari-Noir Layout Uplift
```

## With

```md
# Epic E06: Atari-Noir Layout Uplift
```

## Target File

`.savepoint/releases/v1/epics/E06-atari-noir-layout/E06-Detail.md`

## Insert After

```md
| `internal/board/epic_panel.go` | Refine epic panel layout |
```

## With

```md

## Implemented As

- `internal/styles/palette.go` defines the Atari-Noir hex palette plus ANSI256/ANSI16 fallbacks.
- `internal/styles/styles.go` centralizes header, divider, footer, column, card, detail, glyph, and semantic tag styles.
- `internal/board/view.go` renders the static SavePoint header, the static `PLAN │ BUILD │ AUDIT` footer, subdued navigation hints, and overlay composition.
- `internal/board/card.go` wraps long task titles instead of truncating them and uses a green `▣` marker for the router-priority task.
- `internal/board/detail.go` adds spacing below Acceptance Criteria and Implementation Plan headings, renders checklist state with `☑`/`□`, and labels the router-priority task.
- `internal/data/task.go` and `internal/data/parser.go` represent Implementation Plan items as `CheckItem{Text, Done}` and parse `- [x]`, `- [ ]`, and legacy `- ` items.
- `internal/data/write.go` persists task status and phase changes with mtime conflict protection.
- `internal/board/watch.go` adds fsnotify-based recursive release directory watching, 100ms debounce, and reload messages used by `internal/board/update.go`.

## Implementation Deltas

- The original visual-only scope expanded to include board usability fixes requested during E06: footer navigation hints, detail overlay spacing, card title word-wrap, checklist state rendering, router priority markers, and auto-refresh on task file changes.
- The current implementation still uses rounded frames/borders on header, board, columns, cards, and epic panel. That differs from the design language of "unnecessary borders removed" and should be reconciled before closing the visual uplift.
- `styles.Divider` exists but the main view does not render explicit full-width top/bottom divider lines; divider-like output currently comes from framed containers and component separators.
```

## Quality Review

## Must Fix Before Close

- None remaining after approval closeout.

## Carry Forward

- `internal/board/watch.go:17`: watcher errors are consumed silently. This is acceptable for v1 audit as non-blocking resilience, but a future task should surface watcher failures in `StatusMessage` or diagnostics.
- `internal/board/update.go:68` and `internal/board/update.go:84`: `ErrMtimeConflict` is intentionally non-destructive, but the user receives no visible conflict message because that error is ignored. Consider showing a manual-refresh message so conflicts are understandable.
- `go.mod`: `github.com/fsnotify/fsnotify` is currently listed as indirect even though project code imports it directly. `go mod tidy` may fix this metadata.
- Root instruction drift: `AGENTS.md` references `agent-skills/audit/SKILL.md`, but the repository contains `agent-skills/savepoint-audit/SKILL.md`.

## Already Fixed

- `go build ./...`: PASS during audit.
- `go test ./...`: PASS during audit.
- Approval closeout added explicit full-width divider lines rendered through `styles.Divider`.
- Approval closeout removed unnecessary borders from unfocused header, board, column, card, and epic panel surfaces while preserving focused Atari Orange borders.
- T008 checklist parsing preserves checked/unchecked state in `data.CheckItem` and renders checked items distinctly.
- T009 router priority marker is wired through model, card, column, and detail rendering.
- T010 auto-refresh watcher is isolated in `internal/board/watch.go` and cleanly disabled in tests when `Model.Watcher` is nil.
