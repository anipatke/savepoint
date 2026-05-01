---
id: E03-cli-simplify/T005-update-cli-tests
status: planned
objective: "Update CLI tests to reflect board-only interface"
depends_on: [E03-cli-simplify/T003-strip-run]
---

# T005: Update CLI Tests

## Acceptance Criteria

- All CLI tests pass with board-only interface.
- No test cases reference init, audit, or doctor commands.
- `test/cli/`, `test/commands/` contain only board-relevant tests.
- `npm test` passes.

## Implementation Plan

- [ ] Review all files in `test/cli/` and `test/commands/`.
- [ ] Update or delete tests referencing init/audit/doctor.
- [ ] Ensure board command tests still cover TTY and non-TTY paths.
- [ ] Run `npm test`.
