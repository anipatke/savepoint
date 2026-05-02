---
type: epic-design
status: planned
---

# Epic E03: Board TUI Core

## Purpose

Build the core Bubble Tea model, update loop, view rendering, styles, and layout system. This is the engine — no components yet, just the skeleton that renders a responsive grid.

## What this epic adds

- `Model` struct with tasks, focus state, selected epic/release, overlay mode, dimensions.
- `Init()`, `Update(msg)`, `View()` methods.
- Atari-Noir color palette in `internal/styles/`.
- Manual layout calculation: column widths, epic panel width, responsive breakpoints.
- `tea.WindowSizeMsg` handling for terminal resize.

## Definition of Done

- `go run main.go` launches and shows a bordered 3-column grid.
- Columns resize correctly on terminal resize.
- At 3 breakpoints: >=120 cols shows epic panel + 3 columns; 80-119 shows 3 columns; <80 shows 1 column.
- Colors match Atari-Noir palette in truecolor; degrade gracefully in 256/16-color modes.
- `q` quits cleanly.
- No layout break at any width.

## Components and files

| Path | Purpose |
|------|---------|
| `internal/board/model.go` | Model struct and Init |
| `internal/board/update.go` | Update loop, message handling |
| `internal/board/view.go` | View composition |
| `internal/styles/palette.go` | Atari-Noir hex colors |
| `internal/styles/styles.go` | Lip Gloss style definitions |
| `internal/board/layout.go` | Width/height calculations |
