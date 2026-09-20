---
type: epic-design
status: audited
---

# Epic E01: TUI Optimisation

## Purpose

Performance, layout robustness, and structural improvements for the board TUI. Focused on fixing rendering edge-cases (resize clipping, wasted space) and improving the codebase's architectural hygiene around file watching and data refresh.

## Definition of Done

- Right-border clipping is eliminated at all terminal widths ≥ 40
- Resize handling is robust — no corruption, no artifacts when growing/shrinking
- The board auto-refreshes when task files change on disk via fsnotify
- Board columns use virtual viewport scrolling — 4-5 cards visible with `↑/↓` indicators
- Detail overlay is height-capped (~70%) with scroll indicators for overflow content
- All scroll indicators use dim/subtle styling consistent with Atari Noir aesthetic

## Implemented As

- Header Next Activity rendering is implemented in `internal/board/view.go` using `FormatNextActivity` and `styles.HeaderRight`.
- Column and detail scrolling are implemented through `ColumnOffsets`, `DetailOffset`, height-aware layout, viewport slicing, and subtle scroll indicators in `internal/board/model.go`, `internal/board/layout.go`, `internal/board/update.go`, `internal/board/column.go`, and `internal/board/detail.go`.
- Focus border stability is implemented by rendering unfocused columns with `styles.ColumnUnfocused`, matching focused column border dimensions while using the subtle border color.
- Naming conventions were reconciled across workflow docs: per-epic `Design.md` files became `E##-Detail.md`, and release `PRD.md` became `{release}-PRD.md`.
- The originally listed fsnotify watcher scope was already present from the prior TUI epic and was reviewed as existing behavior rather than newly introduced in this epic.

## Components and files

| Path | Purpose |
|------|---------|
| `internal/board/layout.go` | Layout arithmetic and resize guards, height-aware ContentHeight |
| `internal/board/view.go` | Minimum width clamping, height passthrough to renderers |
| `internal/board/watch.go` | File watcher and reload commands |
| `internal/board/model.go` | Watcher lifecycle, ColumnOffsets, DetailOffset state |
| `internal/board/update.go` | Reload message handling, auto-scroll offsets, PageUp/PageDown |
| `internal/board/column.go` | Virtual viewport slicing, scroll indicator rendering |
| `internal/board/detail.go` | Height-capped overlay, detail scroll indicators |
| `internal/styles/styles.go` | ScrollIndicator style (dim/subtle) |
