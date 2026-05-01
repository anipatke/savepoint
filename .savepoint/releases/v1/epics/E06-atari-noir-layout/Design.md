---
type: epic-design
status: planned
---

# Epic E07: Atari-Noir Layout Uplift

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
