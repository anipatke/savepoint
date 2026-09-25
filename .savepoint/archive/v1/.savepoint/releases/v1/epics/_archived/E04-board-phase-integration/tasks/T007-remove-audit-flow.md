---
id: E04-board-phase-integration/T007-remove-audit-flow
status: planned
objective: "Strip audit proposal loading, AuditReviewApp, and audit event logging from board.ts"
depends_on: [E04-board-phase-integration/T004-phase-transitions]
---

# T007: Remove Audit Flow

## Acceptance Criteria

- `board.ts` has no imports from `src/audit/` or `src/tui/audit-review/`.
- `board.ts` does not check for audit proposals.
- `board.ts` does not launch `AuditReviewApp` after board exit.
- `board.ts` does not log audit events.
- Board command tests pass.

## Implementation Plan

- [ ] Read `board.ts` and `test/commands/board.test.ts`.
- [ ] Remove audit-related imports.
- [ ] Remove `hasAuditProposals` function.
- [ ] Remove audit-requested flag and post-exit audit flow.
- [ ] Remove `auditProposalsAvailable` from `loadBoardData`.
- [ ] Simplify `runBoard` to exit cleanly after Ink TUI closes.
- [ ] Update board command tests.
- [ ] Run `npm test`.
