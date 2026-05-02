---
type: epic-design
status: planned
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
| `internal/board/view_test.go` | Header, formatting, and checkbox tests |
