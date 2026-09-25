---
id: E04-board-phase-integration/T005-phase-gates
status: planned
objective: "Update gates.ts: must reach audit before done; cannot skip phases"
depends_on: [E04-board-phase-integration/T001-board-data-phases]
---

# T005: Phase Gates

## Acceptance Criteria

- `evaluateTransitionGate` blocks advance to `done` unless `phase === 'audit'`.
- `evaluateTransitionGate` blocks phase skips (e.g., build -> audit).
- Retreat transitions (done -> audit, audit -> test, etc.) are always allowed.
- Gate tests cover all phase transition combinations.

## Implementation Plan

- [ ] Read `gates.ts` and `test/tui/io/gates.test.ts`.
- [ ] Update `GateBlockReason` to include phase skip reason.
- [ ] Add phase check before allowing done transition.
- [ ] Add phase adjacency check for in_progress -> in_progress transitions.
- [ ] Update gate tests.
- [ ] Run `npm test`.
