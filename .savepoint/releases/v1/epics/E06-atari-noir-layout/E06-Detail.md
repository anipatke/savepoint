---
type: epic-design
status: audited
---

# Epic E06: Atari-Noir Layout Uplift

## Purpose

Update the existing TUI to align with the **Atari-Noir SavePoint design system**. This is a visual and layout uplift focused on restraint, clarity, and control, preserving existing functionality while eliminating visual noise and strictly adhering to the new color and layout rules.

## What this epic adds

- Centralized Atari-Noir color system in `palette.go`.
- Restructured `styles.go` mapped to exact usage rules.
- Static header component: `▣  S A V E P O I N T`.
- Static footer component: `PLAN      │      BUILD      │      AUDIT`.
- Full-width subtle horizontal dividers.
- Refined component styling (columns, cards, epic panel) to remove unnecessary borders and improve spacing.

## Definition of Done

- All colors match the Atari-Noir strict specification.
- Header renders correctly with no wrapping.
- Footer displays PLAN / BUILD / AUDIT cleanly without dynamic status text.
- Full-width dividers frame the content.
- Cards/panels feel like surfaces, not text blocks.
- Spacing is consistent and components breathe.
- Focus/selection is obvious through position and Atari Orange borders, not relying solely on color text changes.
- Existing functionality remains 100% intact.

## Components and files

| Path | Purpose |
|------|---------|
| `internal/styles/palette.go` | Define strict hex values |
| `internal/styles/styles.go` | Map styles to new usage rules |
| `internal/board/view.go` | Implement header, footer, dividers |
| `internal/board/layout.go` | Adjust layout height calculations |
| `internal/board/column.go` | Refine column styling and spacing |
| `internal/board/card.go` | Refine card surface and focus state |
| `internal/board/epic_panel.go` | Refine epic panel layout |

## Implemented As

- `internal/styles/palette.go` defines the Atari-Noir hex palette plus ANSI256/ANSI16 fallbacks.
- `internal/styles/styles.go` centralizes header, divider, footer, column, card, detail, glyph, and semantic tag styles.
- `internal/board/view.go` renders the static SavePoint header, explicit full-width top/bottom dividers, the static `PLAN │ BUILD │ AUDIT` footer, subdued navigation hints, and overlay composition.
- `internal/board/card.go` wraps long task titles instead of truncating them, uses a green `▣` marker only for the release/epic/task-scoped router-priority task, and keeps all done cards on the orange build glyph.
- `internal/board/detail.go` adds spacing below Acceptance Criteria and Implementation Plan headings, renders checklist state with `☑`/`□`, and labels only the release/epic/task-scoped router-priority task.
- `internal/data/task.go` and `internal/data/parser.go` represent Implementation Plan items as `CheckItem{Text, Done}` and parse `- [x]`, `- [ ]`, and legacy `- ` items.
- `internal/data/write.go` persists task status and phase changes with mtime conflict protection.
- `internal/board/watch.go` adds fsnotify-based recursive release directory watching from the `.savepoint/releases/` root, dynamically adds watches for newly-created subdirectories, applies a 100ms debounce, and sends reload messages used by `internal/board/update.go`.

## Implementation Deltas

- The original visual-only scope expanded to include board usability fixes requested during E06: footer navigation hints, detail overlay spacing, card title word-wrap, checklist state rendering, router priority markers, and auto-refresh on task file changes.
- Audit closeout removed unnecessary borders from unfocused header, board, column, card, and epic-panel surfaces while preserving Atari Orange borders for focused cards/columns and detail overlays.
- User-approved visual guardrail: Background, Surface, and Surface 2 are intentionally the same black value in the terminal TUI. Do not reintroduce subtly different background fills for panels/cards.
- The main view now renders explicit full-width divider lines with `styles.Divider` rather than relying on component frames as separators.
- Auto-refresh reload now refreshes both task data and release/epic indexes so newly-added task files and newly-added epic directories can appear without restarting the board.
- Router-priority matching is intentionally scoped by release and epic before task number; a router task of `T001` must not mark every epic's `T001` as priority.
