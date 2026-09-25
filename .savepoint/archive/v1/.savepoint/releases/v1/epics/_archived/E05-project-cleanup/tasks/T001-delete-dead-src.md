---
id: E05-project-cleanup/T001-delete-dead-src
status: planned
objective: "Delete src/audit/ and src/tui/audit-review/ directories"
depends_on: [E04-board-phase-integration/T008-board-tests]
---

# T001: Delete Dead Source

## Acceptance Criteria

- `src/audit/` directory does not exist.
- `src/tui/audit-review/` directory does not exist.
- No remaining imports reference files in deleted directories.
- Build compiles successfully.

## Implementation Plan

- [ ] Search codebase for imports from `src/audit/` and `src/tui/audit-review/`.
- [ ] Confirm all references were removed in E04.
- [ ] Delete `src/audit/` and all files within.
- [ ] Delete `src/tui/audit-review/` and all files within.
- [ ] Run `npm run build`.
