---
 type: epic-design
 status: planned
 ---

 # Epic E05: Project Cleanup

 ## Purpose

 Delete all dead files, tests, templates, and assets. Rewrite AGENTS.md as a lean agent guide. Clean package.json of unused dependencies. Verify everything builds, typechecks, lints, and tests.

 ## What this epic adds

 - Removal of dead source directories (`src/audit/`, `src/tui/audit-review/`).
 - Removal of dead test directories (`test/audit/`, `test/init/`, `test/templates/`).
 - Removal of dead assets (`templates/`, `agent-skills/`, stale markdown files).
 - Cleanup of `.savepoint/` data (delete audit/, metrics/, Design.md, PRD.md, visual-identity.md).
 - Lean AGENTS.md: router read order, phase workflow, build commands, code style.
 - Clean package.json with no unused devDependencies.

 ## Definition of Done

 - No dead files remain in `src/`, `test/`, `templates/`, `agent-skills/`.
 - `.savepoint/` contains only `router.md`, `config.yml`, and `releases/`.
 - AGENTS.md is under 100 lines, covers router, phases, build, and style.
 - `package.json` has no unused dependencies.
 - `npm run build && npm run typecheck && npm run lint && npm test` all pass.
 - No references to deleted commands, audit pipeline, or old state model remain anywhere.

 ## Components and files

 | Path | Purpose |
 |------|---------|
 | `AGENTS.md` | Rewritten lean agent guide |
 | `package.json` | Cleaned dependencies |
 | (deleted) `src/audit/` | Removed (7 files) |
 | (deleted) `src/tui/audit-review/` | Removed (5 files) |
 | (deleted) `test/audit/` | Removed (6 files) |
 | (deleted) `test/init/` | Removed |
 | (deleted) `test/templates/` | Removed |
 | (deleted) `test/commands/audit.test.ts` | Removed |
 | (deleted) `test/tui/audit-review/` | Removed |
 | (deleted) `templates/` | Removed |
 | (deleted) `agent-skills/` | Removed |
 | (deleted) `CLAUDE.md` | Removed |
 | (deleted) `GEMINI.md` | Removed |
 | (deleted) `ink-cli-ui-design.zip` | Removed |
 | (deleted) `vitest.config.ts` | Removed |
 | (deleted) `.savepoint/audit/` | Removed |
 | (deleted) `.savepoint/metrics/` | Removed |
 | (deleted) `.savepoint/Design.md` | Removed |
 | (deleted) `.savepoint/PRD.md` | Removed |
 | (deleted) `.savepoint/visual-identity.md` | Removed |
