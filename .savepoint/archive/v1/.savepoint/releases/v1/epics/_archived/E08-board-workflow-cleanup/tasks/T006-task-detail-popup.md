---
id: E08-board-workflow-cleanup/T006-task-detail-popup
status: done
objective: Add a keyboard-driven task detail popup that shows acceptance criteria, implementation checklist state, dependencies, status, and transition affordances.
depends_on:
  - E08-board-workflow-cleanup/T005-ink-board-layout-cleanup
---

# T006: task-detail-popup

## Acceptance Criteria

- A selected task can open a full-screen detail modal using `Enter` (when focused on the board); `Escape` closes the modal. The implementation checklist remains read-only for v0.1.0.
- While the popup is open, underlying board key handlers are disabled so input is handled once.
- The popup shows objective, full ID, status, dependencies, acceptance criteria, implementation checklist progress, transition gate preview, and current warning/message text.
- Enter/backspace status transitions remain predictable and refresh board data after successful writes.
- Detail rendering is compact, no-color compatible, and safe for narrow terminals.
- Ink tests cover open/close, disabled underlying navigation, content sections, transition affordances, loading/error messages, and long text truncation.

## Implementation Plan

- [x] Extend board task data with implementation checklist summary parsed from the task body.
- [x] Add a focused detail popup component or refactor `src/tui/DetailPane.tsx` into popup-oriented parts while preserving one job per file.
- [x] Wire `src/tui/App.tsx` input handling so `Enter` opens detail, `Escape` closes it, and board-level movement/transition keys are disabled in popup mode.
- [x] Render criteria, checklist progress, dependency status, full ID, and transition gate previews with no hidden database state.
- [x] Ensure status transitions across selected epics refresh release-level data and surface mtime conflicts without closing or corrupting detail state.
- [x] Read and adhere to `agent-skills/ink-tui-design/SKILL.md` for modal component structure and visual constraints.
- [x] Update Ink tests for popup keyboard flow and rendered sections.
- [x] Run focused App and detail component tests before handing off.
- [x] Update this task's context log with implementation files read before moving it to review.

## Context Log

- Files read: `.savepoint/router.md`; `.savepoint/releases/v1/epics/E08-board-workflow-cleanup/Design.md`; `.savepoint/releases/v1/epics/E08-board-workflow-cleanup/tasks/T006-task-detail-popup.md`; `.savepoint/visual-identity.md`; `agent-skills/ink-tui-design/SKILL.md`; `agent-skills/ink-tui-design/references/component-patterns.md`; `agent-skills/ink-tui-design/references/ink-gotchas.md`; `agent-skills/ink-tui-design/references/testing-patterns.md`; `src/domain/task.ts`; `src/tui/App.tsx`; `src/tui/Board.tsx`; `src/tui/DetailPane.tsx`; `src/tui/board-data.ts`; `src/tui/state/app-reducer.ts`; `src/tui/state/view-state.ts`; `src/tui/state/reducer.ts`; `test/domain/task.test.ts`; `test/tui/board-data.test.ts`; `test/tui/components/App.test.tsx`; `test/tui/components/Board.test.tsx`; `test/tui/components/DetailPane.test.tsx`; `test/tui/state/app-reducer.test.ts`; `.savepoint/metrics/context-bench.md`.
- Estimated input tokens: ~41,158
- Notes: Implemented Enter/Escape modal flow per acceptance criteria and retained Space/backspace status transitions. Focused tests, typecheck, and lint pass; full format check is still blocked by unrelated pre-existing unformatted files outside this task's edit set.
