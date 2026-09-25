---
id: E08-board-workflow-cleanup/T003-board-data-and-plain-output
status: done
objective: Build release-level board data and update non-TTY board output so epics and readable task labels are visible without changing router state.
depends_on:
  - E08-board-workflow-cleanup/T002-release-task-set-reader
---

# T003: board-data-and-plain-output

## Acceptance Criteria

- `BoardData` includes release epic metadata (including status mapped to `✓`, `▶`, `○` markers), active router epic/task, selected epic, selected epic columns, readable task labels, full task IDs, acceptance criteria summaries, and per-epic warnings.
- `runBoard` loads release-level board data for refresh and non-TTY rendering while keeping router state read-only.
- Status writes derive the target task path from the selected task ID's epic segment, not only from the router's active epic.
- Non-TTY output separates epic navigation from task status output and uses readable task labels before full IDs.
- Tests cover release-level board data, readable labels, acceptance criteria summaries, selected-vs-active epic metadata, and non-TTY output.

## Implementation Plan

- [x] Extend `src/tui/board-data.ts` types to represent release epics, selected epic metadata, readable task labels, full IDs, and acceptance criteria summaries.
- [x] Update board-data construction to consume the release-level reader from T002 while preserving current single-epic behavior through compatible defaults.
- [x] Update `src/commands/board.ts` loading and refresh callbacks to use release-level data without mutating `.savepoint/router.md`.
- [x] Update status-write path resolution to derive `release`, `epic`, and task filename from the selected task ID and current release.
- [x] Update `src/tui/render/plain-table.ts` so plain output shows an epic section, selected/active markers, readable task lines, and full IDs as metadata.
- [x] Add or update focused tests in `test/commands/board.test.ts`, `test/tui/render/plain-table.test.ts`, and board-data tests for the new release-level contract.
- [x] Run focused board-data, command, and plain-render tests before handing off.
- [x] Update this task's context log with implementation files read before moving it to review.

## Context Log

- Files read: `src/tui/board-data.ts`; `src/commands/board.ts`; `src/tui/render/plain-table.ts`; `src/readers/tasks.ts`; `src/domain/task.ts`; `src/domain/ids.ts`; `src/fs/markdown.ts`; `src/domain/router.ts`; `src/tui/state/view-state.ts`; `test/tui/render/plain-table.test.ts`; `test/commands/board.test.ts`; `test/tui/state/reducer.test.ts`; `test/tui/state/app-reducer.test.ts`; `test/tui/components/Board.test.tsx`; `.savepoint/releases/v1/epics/E08-board-workflow-cleanup/Design.md`; `tsconfig.json`.
- Also fixed pre-existing bug: `src/readers/tasks.ts` was not calling `parseTaskDocument` after parsing markdown, so `acceptanceCriteria` was always undefined at runtime.
- 727 tests pass (54 test files). New tests: `test/tui/board-data.test.ts` (14), `test/tui/render/plain-table.test.ts` +7 cases, `test/commands/board.test.ts` +4 cases.
