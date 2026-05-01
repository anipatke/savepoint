---
type: epic-design
status: planned
---

# Epic E04: Board Components

## Purpose

Build all visual components: columns, task cards with phase glyphs, epic panel, detail overlay, epic/release dropdowns, and help overlay.

## What this epic adds

- Column component: header with status label + count, bordered container.
- Task card component: fixed width, truncated title, dimmed metadata, phase glyphs (▣/◇/◆), accent border on focus.
- Epic panel component: left sidebar with epic list, active/selected indicators.
- Detail overlay: centered modal showing full task info, board visible behind.
- Epic selector dropdown: overlay list of epics.
- Release selector dropdown: overlay list of releases.
- Help overlay: keyboard shortcuts reference.

## Definition of Done

- Cards render with correct phase glyphs and colors.
- Focused card has orange accent border.
- Epic panel renders on wide screens; dropdown on narrow.
- Detail overlay opens with Enter, closes with Esc.
- Epic dropdown opens with `e`, release with `r`.
- Help overlay opens with `?`.
- All components render without wrapping or layout breaks.
- Component tests pass.

## Components and files

| Path | Purpose |
|------|---------|
| `internal/board/column.go` | Column rendering |
| `internal/board/card.go` | Task card rendering |
| `internal/board/epic_panel.go` | Epic sidebar / dropdown |
| `internal/board/detail.go` | Task detail overlay |
| `internal/board/dropdown.go` | Generic dropdown component |
| `internal/board/help.go` | Help overlay |
| `internal/board/card_test.go` | Card rendering tests |
| `internal/board/column_test.go` | Column rendering tests |
