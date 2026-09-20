---
id: E05-init-command/T002-target-validation
status: done
objective: "Validate whether a target directory can safely receive a Savepoint scaffold."
depends_on: []
---

# T002: target-validation

## Implementation Plan

- [x] Create `src/init/validate-target.ts` with typed validation results for empty, compatible, already-initialized, and conflicting target directories.
- [x] Treat existing `.savepoint/` and root `AGENTS.md` paths as protected unless an explicit overwrite mode allows compatible replacement.
- [x] Normalize filesystem boundary errors into clear validation failures that include the affected path.
- [x] Add unit tests with temporary directories for empty targets, compatible existing files, incompatible files, missing directories, and inaccessible filesystem errors where practical.
- [x] Export only the validation types and function needed by later init orchestration.
- [x] Update `.savepoint/metrics/context-bench.md` with this task's measured context before setting status to review.

## Context Log

- Files read: not captured before context tracking was added to task files.
- Estimated input tokens: n/a
- Notes: Use `.savepoint/metrics/context-bench.md` for measured rows going forward.
