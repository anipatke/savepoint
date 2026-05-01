---
type: epic-design
status: planned
---

# Epic E01: TUI Optimisation

## Purpose

Performance, layout robustness, and structural improvements for the board TUI. Focused on fixing rendering edge-cases (resize clipping, wasted space) and improving the codebase's architectural hygiene around file watching and data refresh.

## Definition of Done

- Right-border clipping is eliminated at all terminal widths ≥ 40
- Resize handling is robust — no corruption, no artifacts when growing/shrinking
- The board auto-refreshes when task files change on disk via fsnotify

## Components and files

| Path | Purpose |
|------|---------|
| `internal/board/layout.go` | Layout arithmetic and resize guards |
| `internal/board/view.go` | Minimum width clamping |
| `internal/board/watch.go` | File watcher and reload commands |
| `internal/board/model.go` | Watcher lifecycle |
| `internal/board/update.go` | Reload message handling |
