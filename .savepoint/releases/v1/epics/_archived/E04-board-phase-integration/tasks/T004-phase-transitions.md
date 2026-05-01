---
id: E04-board-phase-integration/T004-phase-transitions
status: planned
objective: "Adopt brief navigation model (←/→ columns, ↑/↓ tasks, Enter detail, Space advance, Backspace retreat, e/r selectors)"
depends_on: [E04-board-phase-integration/T002-phase-rendering, E04-board-phase-integration/T003-detail-pane-phases]
---

# T004: Navigation & Phase Transitions

## Acceptance Criteria

- `←`/`→` change focused column (planned / in_progress / done).
- `↑`/`↓` change focused task within the current column.
- `Enter` opens the task detail overlay.
- `Esc` closes the overlay.
- `Space` advances: planned → in_progress:build → in_progress:test → in_progress:audit → done.
- `Backspace` retreats: done → in_progress:audit → in_progress:test → in_progress:build → planned.
- `e` opens epic selector dropdown.
- `r` opens release selector dropdown.
- `?` shows help overlay.
- `q` quits.
- Cannot advance to done unless phase is audit.
- `onTransition` callback accepts phase alongside status.
- No `onAuditRequest` prop or audit-related state in App.
- App reducer tests pass.

## Implementation Plan

- [ ] Read `App.tsx`, `state/app-reducer.ts`, `state/view-state.ts`.
- [ ] Update `onTransition` prop signature to include optional phase.
- [ ] Replace Tab focus switching with `←`/`→` column navigation.
- [ ] Bind `↑`/`↓` to task navigation within column.
- [ ] Bind `Enter` to open detail overlay.
- [ ] Bind `Esc` to close overlay.
- [ ] Bind `Space` to advance phase/status.
- [ ] Bind `Backspace` to retreat phase/status.
- [ ] Bind `e` to epic selector, `r` to release selector.
- [ ] Bind `?` to help overlay, `q` to quit.
- [ ] Remove `onAuditRequest` prop and all audit-related UI.
- [ ] Update reducer state types (remove auditProposalsAvailable).
- [ ] Update App and reducer tests.
- [ ] Run `npm test`.
