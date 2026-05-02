# Audit Proposals: v1.1 E01 TUI Optimisation

## Target File

`.savepoint/Design.md`

## Replace

```md
last_audited: E06-atari-noir-layout (2026-05-02)
```

## With

```md
last_audited: E01-tui-optimisation (2026-05-02)
```

## Target File

`.savepoint/Design.md`

## Replace

```md
- **Board command** (`savepoint board`) reads project state, renders the Atari-Noir TUI board, supports release/epic filtering, detail overlays, task status transitions with mtime-guarded writes, router priority markers, and fsnotify-based task auto-refresh (epic E06).
```

## With

```md
- **Board command** (`savepoint board`) reads project state, renders the Atari-Noir TUI board, supports release/epic filtering, detail overlays, task status transitions with mtime-guarded writes, router priority markers, fsnotify-based task auto-refresh (epic E06), header Next Activity display, height-aware column/detail viewport scrolling, and stable focused/unfocused column border geometry (v1.1 E01).
```

## Target File

`.savepoint/Design.md`

## Replace

```md
    │   └── {E##-epic}/
    │       ├── snapshot.md
    │       └── proposals/
    │           └── proposals.md
```

## With

```md
    │   └── {E##-epic}/
    │       ├── snapshot.md
    │       └── proposals.md
```

## Target File

`.savepoint/Design.md`

## Replace

```md
                └── E##-{epic-name}/
                    ├── Design.md   ← epic delta
                    └── tasks/
                        └── T001-slug.md
```

## With

```md
                └── E##-{epic-name}/
                    ├── E##-Detail.md   ← epic delta
                    └── tasks/
                        └── T001-slug.md
```

## Target File

`.savepoint/Design.md`

## Replace

```md
| **Epic**     | A major feature within a release. Has its own Design.md (delta from project Design).   |
```

## With

```md
| **Epic**     | A major feature within a release. Has its own E##-Detail.md (delta from project Design). |
```

## Target File

`.savepoint/Design.md`

## Replace

```md
**Layout:** single screen with a 3-column task board (`planned`, `in_progress`, `done`), optional epic sidebar on wide terminals, centered overlays for release/epic/help/task detail, static Atari-Noir header/footer, full-width dividers, uniform black TUI backgrounds, and navigation hints. Non-TTY output remains a plain table fallback.
```

## With

```md
**Layout:** single screen with a 3-column task board (`planned`, `in_progress`, `done`), optional epic sidebar on wide terminals, centered overlays for release/epic/help/task detail, static Atari-Noir header/footer, full-width dividers, uniform black TUI backgrounds, and navigation hints. The header can show a compact right-aligned Next Activity value from router state. Columns and detail overlays use height-aware viewport slicing with subtle above/more scroll indicators. Focused and unfocused columns preserve the same rounded-border geometry so focus changes do not shift content. Non-TTY output remains a plain table fallback.
```

## Target File

`AGENTS.md`

## Replace

```md
| `internal/styles/`                   | Atari-Noir palette constants, terminal color fallbacks, shared TUI styles, semantic glyph/tag styles, and footer/header styling |
```

## With

```md
| `internal/styles/`                   | Atari-Noir palette constants, terminal color fallbacks, shared TUI styles, stable column border styles, scroll indicators, semantic glyph/tag styles, and footer/header styling |
```

## Target File

`.savepoint/releases/v1.1/epics/E01-tui-optimisation/E01-Detail.md`

## Replace

```md
## Components and files
```

## With

```md
## Implemented As

- Header Next Activity rendering is implemented in `internal/board/view.go` using `FormatNextActivity` and `styles.HeaderRight`.
- Column and detail scrolling are implemented through `ColumnOffsets`, `DetailOffset`, height-aware layout, viewport slicing, and subtle scroll indicators in `internal/board/model.go`, `internal/board/layout.go`, `internal/board/update.go`, `internal/board/column.go`, and `internal/board/detail.go`.
- Focus border stability is implemented by rendering unfocused columns with `styles.ColumnUnfocused`, matching focused column border dimensions while using the subtle border color.
- Naming conventions were reconciled across workflow docs: per-epic `Design.md` files became `E##-Detail.md`, and release `PRD.md` became `{release}-PRD.md`.
- The originally listed fsnotify watcher scope was already present from the prior TUI epic and was reviewed as existing behavior rather than newly introduced in this epic.

## Components and files
```

## Quality Review

## Must Fix Before Close

No remaining must-fix items.

## Carry Forward

- `make build && make test` could not be verified directly in this environment. The audit verified the underlying Go gates with `go build ./...` and `go test ./...`, both passing.
- `internal/board/update.go` allows `DetailOffset` to grow beyond the rendered detail body and relies on render-time clamping. This is not currently user-visible breakage, but a future cleanup could clamp detail scrolling against actual visible body height to keep model state closer to rendered state.

## Already Fixed

- Stale unchecked implementation checklist items in `T001-next-activity-header.md`, `T004-update-instruction-files.md`, and `T005-update-cross-references.md` were reconciled after audit verification.
- Column viewport slicing and subtle scroll indicators are covered by `internal/board/column_test.go`.
- Detail overlay scrolling is covered by `internal/board/detail_test.go`.
- Height-aware layout is covered by `internal/board/layout_test.go`.
- Header Next Activity rendering and truncation are covered by `internal/board/view_test.go`.
- Focused and unfocused column border dimensions are covered by `TestRenderColumn_focusStatesUseStableBorderDimensions`.
