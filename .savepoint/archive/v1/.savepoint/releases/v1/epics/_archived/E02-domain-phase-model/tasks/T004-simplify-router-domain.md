---
id: E02-domain-phase-model/T004-simplify-router-domain
status: planned
objective: "Collapse router states from 6 to 3 in router.ts"
depends_on: []
---

# T004: Simplify Router Domain

## Acceptance Criteria

- `router.ts` defines 3 states: `planning`, `building`, `reviewing`.
- `isRouterStateValue` validates only the 3 states.
- `validateRouterState` accepts the 3-state model.
- `RouterState` interface unchanged (still has `state`, `release`, `epic`, `task`, `next_action`).
- All router tests pass.

## Implementation Plan

- [ ] Read existing `router.ts` and `test/domain/router.test.ts`.
- [ ] Replace `ROUTER_STATES` array with 3 states.
- [ ] Update `isRouterStateValue` and `validateRouterState`.
- [ ] Update `test/domain/router.test.ts` with 3-state cases.
- [ ] Run `npm test`.
