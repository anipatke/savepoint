---
id: E08-board-workflow-cleanup/T008-board-workflow-integration
status: done
objective: Close out the board workflow cleanup by reconciling E06/E07 board rework, updating integration coverage, and running the full quality gates.
depends_on:
  - E08-board-workflow-cleanup/T006-task-detail-popup
  - E08-board-workflow-cleanup/T007-templates-acceptance-criteria
---

# T008: board-workflow-integration

## Acceptance Criteria

- E06 board command behavior still works after release-level browsing, including full-screen alternate buffer setup/teardown (restoring shell history cleanly), non-TTY fallback, TTY launch, refresh, transition gates, and mtime conflict handling.
- E07 audit review entry still works for the router's active epic and is not confused with a merely browsed epic.
- Existing E07 T008 audit-pipeline integration scope is respected; this task does not change E07 task status or audit lifecycle semantics.
- Focused tests cover cross-epic board browsing through command/state/component layers.
- Full quality gates pass or any blocker is documented in this task before handoff.
- The E08 task set remains independently buildable and ready for normal implementation after E07 audit closeout.

## Implementation Plan

- [x] Review the E06 and E07 tests affected by board data shape, App keybindings, audit entry, and non-TTY output, then update only assertions that reflect intentional E08 behavior.
- [x] Add integration coverage for browsing a non-active epic, refreshing after underlying task-file changes, transitioning a selected task by ID-derived path, and preserving active-router audit review behavior.
- [x] Run focused command, state, component, reader, template, and audit-review tests that cover E08 changes.
- [x] Run `npm run typecheck`, `npm run lint`, `npm run format:check`, `npm run build`, and `npm test`.
- [x] Update E08 task context logs and `.savepoint/metrics/context-bench.md` with measured implementation context before moving this task to review.
- [x] Leave router advancement to the normal post-E07 workflow; do not mark E08 implementation active while E07 is still incomplete or unaudited.
- [x] Simplify the canonical task workflow to `planned`, `in_progress`, and `done`, while preserving read-time compatibility for historical `backlog` and `review` task files.

## Rework Notes

- E06 rework expected: board data moves from active-epic-only to release-level, reducers gain epic focus, the board collapses to a calmer three-lane workflow plus the epic rail, and plain output changes from ID-first rows to readable task lines.
- E07 rework expected: audit proposal entry remains tied to the router active epic, and audit-review TUI patterns should be reused for mode-gated input and two-pane detail behavior where helpful.
- E04/E05 rework expected: scaffolded workflow docs and prompts must learn acceptance criteria; historical task files should remain compatible rather than being hand-migrated.

## Context Log

- Files read: `.savepoint/router.md`; `.savepoint/releases/v1/epics/E08-board-workflow-cleanup/Design.md`; `.savepoint/releases/v1/epics/E08-board-workflow-cleanup/tasks/T008-board-workflow-integration.md`; `agent-skills/ink-tui-design/SKILL.md`; `.savepoint/visual-identity.md`; `src/commands/board.ts`; `src/tui/App.tsx`; `src/tui/board-data.ts`; `src/tui/state/app-reducer.ts`; `src/tui/state/view-state.ts`; `src/tui/state/reducer.ts`; `src/tui/Board.tsx`; `src/tui/DetailPane.tsx`; `src/tui/io/write-active-router-epic.ts`; `src/tui/io/write-active-router-release.ts`; `src/tui/render/plain-table.ts`; `src/tui/theme/palette.ts`; `src/domain/config.ts`; `.savepoint/config.yml`; `test/commands/board.test.ts`; `test/commands/board-tty.test.ts`; `test/tui/components/App.test.tsx`; `test/tui/components/Board.test.tsx`; `test/tui/components/DetailPane.test.tsx`; `test/tui/state/reducer.test.ts`; `test/tui/state/app-reducer.test.ts`; `test/tui/board-data.test.ts`; `test/tui/render/plain-table.test.ts`; `.savepoint/metrics/context-bench.md`.
- Estimated input tokens: ~43,206
- Notes: Integration closeout for E08. Fixed the release-workflow gaps found in the manual audit: rail browsing now previews the selected epic immediately, `Enter` on the rail activates that epic in the router, transition gates evaluate against release-level task data so cross-epic dependencies still block invalid moves, and horizontal board navigation clamps instead of wrapping. Added explicit release navigation with `[` and `]`, backed by router release rewrites and reload of the first available epic/task in the selected release. Reworked the board layout to a calmer four-surface workflow (`Epics`, `Planned`, `In Progress`, `Done`) with deterministic terminal-width sizing, stronger panel outlines, and detail sections that better match `.savepoint/visual-identity.md` through dark panel surfaces, warm text, orange focus treatment, purple epic accents, and retro iconography. Simplified the canonical task status model to `planned`, `in_progress`, and `done`, while preserving read-time compatibility for historical `backlog` and `review` task files and updating write paths to persist only canonical values. Tightened the workflow docs and templates so any explicit user-requested audit must persist `.savepoint/audit/{E##-epic}/snapshot.md` and `.savepoint/audit/{E##-epic}/proposals.md` before replying, and so task-building closes at `done` instead of `review`. Focused verification for this pass: targeted board/status/theme/template tests now pass after the simplification. `.savepoint/metrics/context-bench.md` could not be updated in this pass because the file is currently deleted in the workspace.
