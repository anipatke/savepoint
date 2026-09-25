---
id: E02-domain-phase-model/T001-phase-types
status: planned
objective: "Add TaskPhase type, constants, transitions, and accent color mapping to status.ts"
depends_on: []
---

# T001: Phase Types

## Acceptance Criteria

- `status.ts` exports `TaskPhase` type (`build | test | audit`).
- Phase order constant exists: `['build', 'test', 'audit']`.
- `getNextPhase(current)` and `getPreviousPhase(current)` functions exist.
- `PHASE_ACCENTS` maps each phase to a hex color string.
- Phase transition validation: `isPhaseTransitionAllowed(from, to)` returns boolean.
- All phase functions have unit tests.

## Implementation Plan

- [ ] Read existing `status.ts` to understand current structure.
- [ ] Add `TaskPhase` type and `TASK_PHASES` constant array.
- [ ] Add `getNextPhase` and `getPreviousPhase` helpers.
- [ ] Add `PHASE_ACCENTS` mapping.
- [ ] Add `isPhaseTransitionAllowed` with no-skipping rule.
- [ ] Write/update `test/domain/status.test.ts` with phase tests.
- [ ] Run `npm test` and verify phase tests pass.
