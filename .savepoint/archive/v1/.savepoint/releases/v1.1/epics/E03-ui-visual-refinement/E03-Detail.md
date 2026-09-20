---
type: epic-design
status: audited
---

# Epic E03: UI Visual Refinement

## Purpose

Polish the Atari-Noir TUI rendering across Next Activity indicator repositioning below the header with phase-aligned styling, checkbox rendering at sentence boundaries (not line breaks), unify icon/glyph determination across task cards and the epic sidebar, and force a consistent 256-color terminal profile for reliable cross-terminal rendering.

## Definition of Done

- Next Activity renders as a dedicated line below the header with phase-aligned coloring
- Checkboxes render at sentence starts, not at line break positions
- ~~Total content width equals available space at every breakpoint with no underflow~~ (absorbed into T001, which is marked done; existing layout tests pass)
- Task card and epic sidebar glyphs are determined by a shared helper, unified on the same status→glyph mapping
- Terminal color profile is forced to 256-color, ensuring consistent background rendering on all terminals including PowerShell

## Components and files

| Path | Purpose |
|------|---------|
| `internal/board/layout.go` | Width arithmetic, layout breakpoints, min width guard (stable — existing implementation verified) |
| `internal/board/view.go` | Header rendering, Next Activity line, checkbox rendering point |
| `internal/board/update.go` | Resize handling, terminal dimension clamping |
| `internal/styles/styles.go` | Phase-aligned styles for Next Activity |
| `internal/data/router.go` | RouterState model |
| `internal/data/task.go` | Task markdown parsing for checkbox placement |
| `internal/board/status.go` | Shared status glyph mapping helper (new) |
| `internal/board/card.go` | Task card glyph determination (updated for shared helper) |
| `internal/board/board.go` | Board startup; sets the Lipgloss color profile to ANSI256 before model initialization |
| `internal/board/view_test.go` | Header, formatting, and checkbox tests |

## Implemented As

- `internal/board/view.go` renders `next_action` as a separate line below the header through `renderNextActivityLine`, using existing footer phase styles for `PLAN`, `BUILD`, and `AUDIT`.
- `internal/board/layout.go` accounts for the optional Next Activity line when calculating board chrome and content height.
- `internal/data/parser.go` joins hard-wrapped checklist continuation lines before rendering so markdown wrap points do not create duplicate checklist items.
- `internal/board/detail.go` splits checklist item text on semantic sentence boundaries and emits one `[ ]` or `[x]` marker per sentence.
- `internal/data/task.go` adds `Task.Status` plus status constants, including `audited`, for shared board glyph rendering.
- `internal/board/status.go` centralizes the planned, in-progress, done, and audited status glyph mapping used by task cards and the epic sidebar.
- `internal/board/card.go` uses explicit `Task.Status` when present, while preserving the legacy column/stage glyph fallback for older task data.
- `internal/board/epic_panel.go` delegates epic status glyph rendering to the shared status helper.
- `internal/board/board.go` sets Lipgloss to the ANSI256 color profile at board startup. This satisfies the deterministic 256-color rendering intent, although the implementation lives at the board boundary instead of `main.go`.
