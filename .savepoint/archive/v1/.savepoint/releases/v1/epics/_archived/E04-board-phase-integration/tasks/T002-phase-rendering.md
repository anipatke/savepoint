---
id: E04-board-phase-integration/T002-phase-rendering
status: planned
objective: "Render phase glyphs with phase-derived color on in_progress task cards in Board.tsx"
depends_on: [E04-board-phase-integration/T001-board-data-phases]
---

# T002: Phase Rendering

## Acceptance Criteria

- In_progress task cards show phase glyphs inline: `▣` (build), `◇` (test), `◆` (audit).
- Phase glyph uses the accent color mapped by `PHASE_ACCENTS`.
- Planned and done task cards show no phase glyph.
- Cards are fixed width, titles are truncated (never wrap), metadata is dimmed.
- Focused card has accent border (Atari Orange) and slight background tint.
- Rendering works for 24-bit, 256-color, 16-color, and `NO_COLOR=1` terminals.
- Board component tests pass.

## Implementation Plan

- [ ] Read `Board.tsx` and `test/tui/components/`.
- [ ] Import `PHASE_ACCENTS` and `TaskPhase` in Board.
- [ ] Add `PHASE_GLYPHS` map: build→▣, test→◇, audit→◆.
- [ ] Update `TaskCardComponent` to accept phase and render glyph with phase color.
- [ ] Apply card styling rules: fixed width, truncate title, dimmed metadata, accent border on focus.
- [ ] Update component tests.
- [ ] Run `npm test`.
