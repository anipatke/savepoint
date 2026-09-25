---
id: E05-project-cleanup/T002-delete-dead-tests
status: planned
objective: "Delete test/audit/, test/init/, test/templates/, test/commands/audit.test.ts, test/tui/audit-review/"
depends_on: [E04-board-phase-integration/T008-board-tests]
---

# T002: Delete Dead Tests

## Acceptance Criteria

- `test/audit/` directory does not exist.
- `test/init/` directory does not exist.
- `test/templates/` directory does not exist.
- `test/commands/audit.test.ts` does not exist.
- `test/tui/audit-review/` directory does not exist.
- `npm test` still passes.

## Implementation Plan

- [ ] Delete `test/audit/` and all files within.
- [ ] Delete `test/init/` and all files within.
- [ ] Delete `test/templates/` and all files within.
- [ ] Delete `test/commands/audit.test.ts`.
- [ ] Delete `test/tui/audit-review/` and all files within.
- [ ] Run `npm test`.
