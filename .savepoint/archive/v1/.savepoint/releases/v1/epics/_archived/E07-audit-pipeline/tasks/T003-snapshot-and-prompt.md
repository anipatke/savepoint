---
id: E07-audit-pipeline/T003-snapshot-and-prompt
status: done
objective: Generate a bounded audit snapshot and AI handoff prompt for an epic without embedding source contents.
depends_on: []
---

# T003: snapshot-and-prompt

## Implementation Plan

- [x] Add `src/audit/snapshot.ts` to generate `.savepoint/audit/{epic}/snapshot.md` from router state, epic/task metadata, gitignore-respecting file listings, and changed-file paths.
- [x] Ensure changed-file output records path and status only, never source or test file contents.
- [x] Create `.savepoint/audit/{epic}/proposals/` alongside the snapshot so later proposal review has a stable artifact location.
- [x] Add `src/audit/prompts.ts` to render the bounded handoff prompt that points an AI agent at the snapshot and proposal directory.
- [x] Add tests for gitignore exclusions, snapshot bounds, changed-file metadata, proposal directory creation, and prompt content.
- [x] Run the focused audit snapshot and prompt tests before handing off the task.

## Context Log

- Files read: `AGENTS.md`; `.savepoint/releases/v1/epics/E07-audit-pipeline/Design.md`; `T003-snapshot-and-prompt.md`; `src/audit/quality-gates.ts`; `src/audit/log.ts`; `src/fs/project.ts`; `src/fs/markdown.ts`; `src/readers/router.ts`; `src/readers/epic.ts`; `src/readers/tasks.ts`; `src/domain/router.ts`; `src/domain/task.ts`; `src/domain/ids.ts`; `src/domain/status.ts`; `src/domain/config.ts`; `templates/prompts/audit-reconciliation.prompt.md`; `test/audit/quality-gates.test.ts`; `test/init/write-scaffold.test.ts`; `package.json`
- Estimated input tokens: ~12,500
- Notes: Implemented `buildSnapshotMarkdown` (pure), `generateSnapshot` (filesystem + injected GitReader), `defaultGitReader` (spawn-based), and `buildAuditPrompt` (pure). 31 tests, 555 total passing. Typecheck and lint clean.
