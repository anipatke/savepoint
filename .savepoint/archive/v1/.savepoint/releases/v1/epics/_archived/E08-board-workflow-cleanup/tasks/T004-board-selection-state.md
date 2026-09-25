---
id: E08-board-workflow-cleanup/T004-board-selection-state
status: done
objective: Model board selection state for epic browsing, task focus, refresh preservation, and task detail popup visibility.
depends_on:
  - E08-board-workflow-cleanup/T003-board-data-and-plain-output
---

# T004: board-selection-state

## Acceptance Criteria

- Board state distinguishes epic focus, task focus, selected epic, selected task, and detail-popup visibility.
- Browsing another epic never changes the active router epic/task in state. Advancing (`Space`) or retreating (`Backspace`) task status is blocked unless the browsed epic matches the active router epic.
- Refresh preserves the selected epic and task when they still exist, and falls back predictably when an epic or task disappears.
- Navigation clamps at list boundaries where wrapping would make the screen feel unstable.
- Reducer tests cover epic switching, task movement within an epic, detail open/close, refresh preservation, and fallback behavior.

## Implementation Plan

- [x] Extend `src/tui/state/view-state.ts` with selected-epic, selected-task, focus-area, and detail-popup state derived from the new `BoardData`.
- [x] Update `src/tui/state/reducer.ts` actions for moving across epics, moving across tasks/status groups, switching focus, opening detail, closing detail, and refreshing release-level data.
- [x] Update `src/tui/state/app-reducer.ts` to preserve warning/message/loading behavior while carrying the expanded view state.
- [x] Replace wraparound movement with explicit clamp behavior where needed to avoid disorienting jumps and layout fragmentation.
- [x] Update state tests for the new release-level board shape and add coverage for detail visibility and selected-vs-active router metadata.
- [x] Run focused TUI state tests before handing off.
- [x] Update this task's context log with implementation files read before moving it to review.

## Context Log

- Files read: `.savepoint/router.md`; `.savepoint/releases/v1/epics/E08-board-workflow-cleanup/Design.md`; this task file; `src/tui/state/view-state.ts`; `src/tui/state/reducer.ts`; `src/tui/state/app-reducer.ts`; `src/tui/board-data.ts`; `test/tui/state/reducer.test.ts`; `test/tui/state/app-reducer.test.ts`.
- Estimated input tokens: ~14,000
- Notes: Added `FocusArea`, `epicIndex`, `selectedEpicRaw`, `focusArea`, `detailVisible` to `BoardViewState`. New actions: `move-epic`, `switch-focus`, `open-detail`, `close-detail`. Vertical task nav changed from wrap to clamp. 174 TUI tests pass.
