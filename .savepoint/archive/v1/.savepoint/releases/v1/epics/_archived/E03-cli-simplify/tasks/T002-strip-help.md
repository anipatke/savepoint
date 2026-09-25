---
id: E03-cli-simplify/T002-strip-help
status: planned
objective: "Remove init, audit, doctor help text from help.ts"
depends_on: [E03-cli-simplify/T001-strip-args]
---

# T002: Strip Help

## Acceptance Criteria

- `help.ts` top-level help shows only `board` command.
- `commandHelp('board')` returns board-specific help text.
- No help text for init, audit, or doctor exists.
- All help tests pass.

## Implementation Plan

- [ ] Read existing `help.ts` and `test/cli/help.test.ts`.
- [ ] Remove init/audit/doctor help text blocks.
- [ ] Update top-level help to list only `board`.
- [ ] Update `test/cli/help.test.ts`.
- [ ] Run `npm test`.
