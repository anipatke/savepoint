---
id: E04-board-phase-integration/T001-board-data-phases
status: planned
objective: "Add phase to BoardTask; remove auditProposalsAvailable from board-data.ts"
depends_on: [E02-domain-phase-model/T002-phase-frontmatter]
---

# T001: Board Data Phases

## Acceptance Criteria

- `BoardTask` interface includes `phase?: TaskPhase`.
- `buildBoardData` carries phase from task frontmatter to board tasks.
- `BoardData` interface has no `auditProposalsAvailable` field.
- `toBoardTask` copies phase from task document.
- `test/tui/board-data.test.ts` passes with phase data.

## Implementation Plan

- [ ] Read `board-data.ts` and `test/tui/board-data.test.ts`.
- [ ] Add `phase?: TaskPhase` to `BoardTask`.
- [ ] Remove `auditProposalsAvailable` from `BoardData`.
- [ ] Update `toBoardTask` to copy phase.
- [ ] Update `buildBoardData` signature and body.
- [ ] Update board-data tests.
- [ ] Run `npm test`.
