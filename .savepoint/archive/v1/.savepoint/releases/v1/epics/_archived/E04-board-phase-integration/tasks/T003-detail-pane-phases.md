---
id: E04-board-phase-integration/T003-detail-pane-phases
status: planned
objective: "Render detail pane as centered overlay with board visible behind; add Esc close"
depends_on: [E04-board-phase-integration/T001-board-data-phases]
---

# T003: Detail Pane Overlay

## Acceptance Criteria

- Detail pane renders as a centered overlay on top of the board.
- Board remains visible and dimmed behind the overlay.
- Overlay shows: task ID, title, epic, release, status, current phase, description, acceptance criteria.
- Pressing `Esc` closes the overlay.
- Overlay width is calculated from terminal size (max 40 chars, or terminal width - 4).
- No layout breaks at any terminal width.

## Implementation Plan

- [ ] Read `DetailPane.tsx` and `test/tui/components/`.
- [ ] Convert DetailPane from replacing board to overlay mode.
- [ ] Add board dimming behind overlay (reduce opacity or dark tint).
- [ ] Add `Esc` handler in `App.tsx` to close overlay.
- [ ] Calculate overlay width from `useWindowSize`.
- [ ] Update component tests.
- [ ] Run `npm test`.
