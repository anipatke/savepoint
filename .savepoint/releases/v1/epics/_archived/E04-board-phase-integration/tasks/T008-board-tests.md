---
id: E04-board-phase-integration/T008-board-tests
status: planned
objective: "Update all board/data/state/io/render tests for phase model"
depends_on: [E04-board-phase-integration/T007-remove-audit-flow]
---

# T008: Board Tests

## Acceptance Criteria

- All board-related tests pass: `test/tui/**/*.test.ts`.
- All command tests pass: `test/commands/*.test.ts`.
- No test references audit proposals or old 6-state model.
- `npm test` passes in full.

## Implementation Plan

- [ ] Run `npm test` and catalog failures.
- [ ] Update board-data tests for phase field.
- [ ] Update state/reducer tests for phase transitions.
- [ ] Update IO tests for phase write-back.
- [ ] Update render tests for phase display.
- [ ] Update command tests for simplified board.
- [ ] Run `npm test` until green.
