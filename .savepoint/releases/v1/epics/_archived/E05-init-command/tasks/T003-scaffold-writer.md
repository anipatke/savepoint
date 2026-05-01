---
id: E05-init-command/T003-scaffold-writer
status: done
objective: "Write the project scaffold from template assets while preserving existing user files."
depends_on:
  - E05-init-command/T002-target-validation
---

# T003: scaffold-writer

## Implementation Plan

- [x] Create `src/init/write-scaffold.ts` to load project templates through `src/templates` APIs and render project metadata variables.
- [x] Write `AGENTS.md` and `.savepoint/` files to the target directory with parent directory creation handled in one place.
- [x] Refuse to overwrite conflicting existing files unless the caller has already selected the explicit overwrite mode.
- [x] Use temporary-file-plus-rename writes for file content so existing files are not partially corrupted.
- [x] Add integration tests covering successful scaffold creation, rendered project metadata, conflict refusal, and explicit overwrite behavior.
- [x] Update `.savepoint/metrics/context-bench.md` with this task's measured context before setting status to review.

## Context Log

- Files read: 17 (router.md, E05 Design.md, T003 task, AGENTS.md, T001, T002, validate-target.ts, manifest.ts, render.ts, load.ts, paths.ts, index.ts, validate-target.test.ts, context-bench.md, init.ts, config.yml, AGENTS.md template)
- Estimated input tokens: ~7,225
- Notes: All 9 new tests pass; 338 total suite green.
