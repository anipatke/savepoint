---
 type: epic-design
 status: planned
 ---

 # Epic E04: Board Phase Integration

 ## Purpose

 Integrate the phase model into the board TUI. Render phase letters on in_progress tasks, enable phase stepping via keyboard, enforce phase gating, and write phase back to task frontmatter. Remove the audit-review flow entirely.

 ## What this epic adds

 - Phase letter rendering (B/T/A) with phase-derived colors on in_progress task cards.
 - Phase stepping: space advances (build→test→audit→done), backspace retreats.
 - Phase gating: cannot advance to done unless phase is audit; cannot skip phases.
 - Phase write-back: `phase` field written alongside `status` in task frontmatter.
 - Complete removal of audit proposal loading, AuditReviewApp, and audit event logging.

 ## Definition of Done

 - Board renders phase letters on in_progress task cards with correct colors.
 - Space/backspace walks through phases; done requires audit phase.
 - Detail pane shows current phase and phase-aware transition preview.
 - Task files are updated with `phase` field on transition.
 - No audit-review code remains in `src/tui/audit-review/` or `src/commands/board.ts`.
 - All board tests pass with phase model.
 - Build and typecheck pass.

 ## Components and files

 | Path | Purpose |
 |------|---------|
 | `src/tui/board-data.ts` | Phase on BoardTask; remove auditProposalsAvailable |
 | `src/tui/Board.tsx` | Phase letter rendering on task cards |
 | `src/tui/DetailPane.tsx` | Phase display and transition preview |
 | `src/tui/App.tsx` | Phase stepping in input handler; remove audit request |
 | `src/tui/io/gates.ts` | Phase-aware gate enforcement |
 | `src/tui/io/write-status.ts` | Phase write-back to frontmatter |
 | `src/tui/io/write-active-router-epic.ts` | Updated for 3-state router |
 | `src/tui/io/write-active-router-release.ts` | Updated for 3-state router |
 | `src/tui/state/app-reducer.ts` | Remove auditProposalsAvailable from state |
 | `src/commands/board.ts` | Remove audit flow |
 | `test/tui/board-data.test.ts` | Phase data model tests |
 | `test/tui/state/` | Phase-aware state tests |
 | `test/tui/io/` | Phase write and gate tests |
 | `test/commands/board.test.ts` | Updated board command tests |
 | (deleted) `src/tui/audit-review/` | Removed (5 files) |
