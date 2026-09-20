---
id: E08-board-workflow-cleanup/T005-ink-board-layout-cleanup
status: done
objective: Rework the Ink board into a stable release-level screen with an epic selector, readable task cards, and bounded terminal layout.
depends_on:
  - E08-board-workflow-cleanup/T004-board-selection-state
---

# T005: ink-board-layout-cleanup

## Acceptance Criteria

- The TTY board renders as an immersive full-screen application using the alternate screen buffer (`\x1b[?1049h`), with a header, epic selector (fixed left rail, min 25 columns), task surface, warning area, and footer.
- Epic items are visually distinct from task cards and explicitly show status (`✓`, `▶`, `○`) and active-router state separately from selected/browsed state. Focus switches between panes via `Tab`/`Shift+Tab`.
- Task cards prioritize readable labels from slug/objective/title-style text while keeping full IDs visible as dim metadata.
- The backlog status remains represented but no longer consumes one fifth of the horizontal screen when that space is better used for epic selection.
- Layout uses terminal-safe widths, truncation, and fixed regions so arrow-key navigation does not fragment or wrap the board.
- Component tests cover epic selector rendering, readable labels, active-vs-selected markers, no-color output, and empty states.

## Implementation Plan

- [x] Refactor `src/tui/Board.tsx` into smaller render parts if needed so epic selection, task grouping, and task cards each have one job.
- [x] Implement alternate screen buffer entry/exit in `src/commands/board.ts` for full-screen takeover.
- [x] Add a release-level header and a fixed-width left rail epic selector that follow Atari-Noir terminal constraints and the Ink guide's screen/part split.
- [x] Replace full-ID-first task rendering with readable card labels, status/count metadata, and full IDs as secondary detail.
- [x] Rebalance the task surface so backlog is still visible but does not monopolize equal-width Kanban space when an epic rail is present.
- [x] Add terminal-width truncation and safe margins for long task labels, objectives, criteria, and epic names.
- [x] Update `src/tui/App.tsx` footer/help text for the new focus and navigation model without adding long in-app instructions.
- [x] Read and adhere to `agent-skills/ink-tui-design/SKILL.md` for component structure and visual constraints.
- [x] Update Ink component tests for the new layout and readable labels.
- [x] Run focused board component tests before handing off.
- [x] Update this task's context log with implementation files read before moving it to review.

## Context Log

- Files read: this task; E08 Design; T003/T004 context logs; `src/tui/Board.tsx`; `src/tui/App.tsx`; `src/tui/board-data.ts`; `src/tui/state/view-state.ts`; `src/tui/state/reducer.ts`; `src/tui/state/app-reducer.ts`; `src/commands/board.ts`; `src/tui/theme/palette.ts`; `agent-skills/ink-tui-design/SKILL.md`; `.savepoint/visual-identity.md`; `test/tui/components/Board.test.tsx`; `test/tui/components/App.test.tsx`.
- Implementation: refactored `Board.tsx` into `EpicRail` (React.memo, fixed 25-col left rail, glyph status + browsed cursor `›`) + `TaskCard` (React.memo, label primary, dim full ID secondary) + flex task surface. `App.tsx` gains header, Tab→switch-focus, rail up/down→move-epic, board left/right guard, mutation guard when browsedEpic≠activeEpic. `board.ts` wraps render in `\x1b[?1049h`/`\x1b[?1049l` alt screen buffer with cleanup on process exit. 755 tests pass, typecheck clean.
