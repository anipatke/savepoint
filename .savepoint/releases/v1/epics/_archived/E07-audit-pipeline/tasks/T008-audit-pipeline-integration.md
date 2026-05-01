---
id: E07-audit-pipeline/T008-audit-pipeline-integration
status: done
objective: Verify the complete audit lifecycle from command execution through review application with end-to-end coverage and full quality gates.
depends_on:
  - E07-audit-pipeline/T007-audit-review-ui
---

# T008: audit-pipeline-integration

## Implementation Plan

- [x] Add end-to-end tests that cover audit skip logging, gate failure short-circuiting, successful snapshot/prompt generation, router `audit-pending` transition, proposal review, and approved proposal application.
- [x] Verify audit artifacts remain bounded by asserting snapshots contain metadata and paths but not source/test code bodies.
- [x] Confirm explicit `--epic` and active-router epic selection work consistently across command, snapshot, proposal, and review paths.
- [x] Tighten command help, error messages, and fixture cleanup discovered by the integration tests.
- [x] Run focused audit integration tests first, then run `npm run typecheck`, `npm run lint`, `npm run format:check`, `npm run build`, and `npm test`.
- [x] Update this task's context log with implementation files read before moving it to review.

## Context Log

- Files read: `.savepoint/router.md`; `.savepoint/releases/v1/epics/E07-audit-pipeline/Design.md`; `.savepoint/releases/v1/epics/E07-audit-pipeline/tasks/T008-audit-pipeline-integration.md`; `src/commands/audit.ts`; `src/audit/quality-gates.ts`; `src/audit/snapshot.ts`; `src/audit/router-state.ts`; `src/audit/prompts.ts`; `src/audit/proposals.ts`; `src/audit/apply-proposals.ts`; `src/audit/log.ts`; `test/commands/audit.test.ts`; `test/cli/run.test.ts`; `test/audit/quality-gates.test.ts`; `test/audit/snapshot.test.ts`; `test/audit/proposals.test.ts`; `test/audit/apply-proposals.test.ts`; `test/audit/prompts.test.ts`; `src/tui/audit-review/state.ts`; `src/tui/audit-review/AuditReviewApp.tsx`; `src/tui/audit-review/summary.ts`; `test/tui/audit-review/state.test.ts`; `test/tui/audit-review/AuditReviewApp.test.tsx`; `test/commands/board.test.ts`; `src/cli/args.ts`; `src/cli/run.ts`; `src/cli/help.ts`; `src/readers/router.ts`; `.savepoint/metrics/context-bench.md`
- Estimated input tokens: ~45,844
- Notes: Full integration context: all E07 audit source and test files, CLI dispatch layer, TUI audit-review state/components, board test patterns for fixture setup, and context bench. Prettier formatting applied to all E07-scoped files with existing warnings.
