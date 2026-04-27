---
id: E05-init-command/T003-scaffold-writer
status: planned
objective: "Write the project scaffold from template assets while preserving existing user files."
depends_on:
  - E05-init-command/T002-target-validation
---

# T003: scaffold-writer

## Implementation Plan

- [ ] Create `src/init/write-scaffold.ts` to load project templates through `src/templates` APIs and render project metadata variables.
- [ ] Write `AGENTS.md` and `.savepoint/` files to the target directory with parent directory creation handled in one place.
- [ ] Refuse to overwrite conflicting existing files unless the caller has already selected the explicit overwrite mode.
- [ ] Use temporary-file-plus-rename writes for file content so existing files are not partially corrupted.
- [ ] Add integration tests covering successful scaffold creation, rendered project metadata, conflict refusal, and explicit overwrite behavior.
