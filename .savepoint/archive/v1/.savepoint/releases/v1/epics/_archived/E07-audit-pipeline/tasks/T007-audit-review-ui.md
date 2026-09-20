---
id: E07-audit-pipeline/T007-audit-review-ui
status: done
objective: "Add the Ink audit review mode that lets users inspect proposals, approve or reject them, edit replacements, and apply approved changes."
depends_on:
  - E07-audit-pipeline/T004-audit-orchestration-router
  - E07-audit-pipeline/T006-audit-review-state
---

# T007: audit-review-ui

## Implementation Plan

- [x] Create `src/tui/audit-review/*.tsx` components that follow the existing terminal palette and the Atari-Noir visual identity within terminal constraints.
- [x] Update the board audit entry point so proposal directories produced by E07 open the review mode for the active epic.
- [x] Render proposal targets, operation summaries, selected diffs, approval state, edit mode, and divergence warnings without dumping full unrelated file contents.
- [x] Wire apply actions through `src/audit/apply-proposals.ts`, require confirmation for high-divergence cases, and surface apply errors without losing review decisions.
- [x] Write audit log entries for applied, rejected-only, and user-aborted review outcomes.
- [x] Add TUI tests with mocked Ink rendering and service callbacks for navigation, approve/reject/edit, apply success, and apply failure.
- [x] Run the focused audit review UI and board integration tests before handing off the task.

## Context Log

- Files read: `.savepoint/router.md`; `.savepoint/releases/v1/epics/E07-audit-pipeline/Design.md`; `.savepoint/visual-identity.md`; `.savepoint/metrics/context-bench.md`; `src/tui/App.tsx`; `src/tui/Board.tsx`; `src/commands/board.ts`; `src/tui/theme/palette.ts`; `src/tui/DetailPane.tsx`; `src/audit/apply-proposals.ts`; `src/audit/proposals.ts`; `src/audit/log.ts`; `src/domain/config.ts`; `src/tui/state/app-reducer.ts`; `src/tui/board-data.ts`; `test/commands/board.test.ts`; `test/commands/board-tty.test.ts`; `test/tui/components/App.test.tsx`; `src/tui/state/reducer.ts`; `src/tui/audit-review/state.ts`; `src/tui/audit-review/summary.ts`; `test/tui/audit-review/state.test.ts`; `vitest.config.ts`; `vitest.config.js`
- Estimated input tokens: ~16,000
- Notes: Implementation required understanding E06 TUI patterns for reducer shape, keyboard input, Ink testing, and board integration; extended T006 state model with confirm-divergence and edit helpers; wired apply through existing E07 audit services.
