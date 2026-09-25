---
id: E08-board-workflow-cleanup/T002-release-task-set-reader
status: done
objective: Read release epics into a release-level task collection that preserves per-epic errors and resolves cross-epic task dependencies.
depends_on:
  - E08-board-workflow-cleanup/T001-acceptance-criteria-model
---

# T002: release-task-set-reader

## Acceptance Criteria

- A release-level reader lists epic directories under `.savepoint/releases/{release}/epics` in deterministic order.
- The reader returns task documents grouped by epic while preserving filesystem/schema errors for each epic.
- Cross-epic dependencies are validated at release scope so valid cross-epic links do not appear as missing dependencies.
- Existing `readEpicTaskSet` behavior remains available for callers that only need one epic.
- Tests cover multiple epics, empty task directories, missing task directories, cross-epic dependencies, duplicate IDs, and true missing dependencies.

## Implementation Plan

- [x] Add a release-level task-set result type in `src/readers/tasks.ts` or a new focused reader module if that keeps the file to one job.
- [x] Implement deterministic epic directory discovery from the release `epics/` directory, ignoring non-directories.
- [x] Reuse existing markdown/frontmatter parsing and acceptance-criteria extraction from T001 for every task file.
- [x] Add release-scope graph validation that treats cross-epic dependencies as valid when the target task exists anywhere in the release.
- [x] Add reader tests in `test/readers/tasks.test.ts` for multi-epic success, partial epic errors, cross-epic dependency success, and release-level graph failures.
- [x] Run focused reader tests before handing off.
- [x] Update this task's context log with implementation files read before moving it to review.

## Context Log

- Files read: `.savepoint/router.md`; `.savepoint/Design.md`; `.savepoint/releases/v1/epics/E08-board-workflow-cleanup/Design.md`; `.savepoint/releases/v1/epics/E06-tui-board/Design.md`; `.savepoint/releases/v1/epics/E07-audit-pipeline/Design.md`; `.savepoint/visual-identity.md`; `agent-skills/ink-tui-design/SKILL.md` and focused references; current task/domain/reader/board/template/test files for E08 planning.
- Estimated input tokens: ~58,547
- Notes: Task-breakdown planning context only. Implementation should re-read this task plus directly touched source/test files.
