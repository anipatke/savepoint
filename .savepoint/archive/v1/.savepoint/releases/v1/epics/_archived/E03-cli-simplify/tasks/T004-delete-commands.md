---
id: E03-cli-simplify/T004-delete-commands
status: planned
objective: "Delete dead command files and supporting directories"
depends_on: [E03-cli-simplify/T003-strip-run]
---

# T004: Delete Commands

## Acceptance Criteria

- `src/commands/init.ts`, `audit.ts`, `doctor.ts`, `result.ts` do not exist.
- `src/init/` directory does not exist.
- `src/templates/` directory does not exist.
- Build still compiles (no broken imports from deleted files).

## Implementation Plan

- [ ] Verify no remaining imports reference files to be deleted.
- [ ] Delete `src/commands/init.ts`, `audit.ts`, `doctor.ts`, `result.ts`.
- [ ] Delete `src/init/` directory and all files within.
- [ ] Delete `src/templates/` directory and all files within.
- [ ] Run `npm run build` to confirm no broken imports.
- [ ] Run `npm run typecheck`.
