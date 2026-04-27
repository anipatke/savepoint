---
id: E05-init-command/T002-target-validation
status: planned
objective: "Validate whether a target directory can safely receive a Savepoint scaffold."
depends_on: []
---

# T002: target-validation

## Implementation Plan

- [ ] Create `src/init/validate-target.ts` with typed validation results for empty, compatible, already-initialized, and conflicting target directories.
- [ ] Treat existing `.savepoint/` and root `AGENTS.md` paths as protected unless an explicit overwrite mode allows compatible replacement.
- [ ] Normalize filesystem boundary errors into clear validation failures that include the affected path.
- [ ] Add unit tests with temporary directories for empty targets, compatible existing files, incompatible files, missing directories, and inaccessible filesystem errors where practical.
- [ ] Export only the validation types and function needed by later init orchestration.
