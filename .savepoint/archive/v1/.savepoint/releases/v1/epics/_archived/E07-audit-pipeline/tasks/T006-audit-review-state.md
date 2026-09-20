---
id: E07-audit-pipeline/T006-audit-review-state
status: done
objective: Model audit proposal review decisions so approve, reject, edit, and apply flows are testable before the Ink UI.
depends_on:
  - E07-audit-pipeline/T005-proposal-validation-apply
---

# T006: audit-review-state

## Implementation Plan

- [x] Add non-rendering audit review state under `src/tui/audit-review/` for loaded proposals, selected operation, review decision, edited replacement text, warnings, and apply readiness.
- [x] Implement reducer actions for navigation, approve, reject, edit, reset edit, show divergence warning, and mark apply result.
- [x] Derive concise review summaries that the UI can render without reparsing proposal documents.
- [x] Keep all filesystem writes out of the state layer; call proposal/apply services through injected callbacks only in later integration.
- [x] Add deterministic reducer tests for navigation bounds, decision changes, edit preservation, ready-to-apply calculation, and warning state.
- [x] Run the focused audit review state tests before handing off the task.

## Context Log

- Files read: `.savepoint/router.md`; `.savepoint/releases/v1/epics/E07-audit-pipeline/Design.md`; `.savepoint/releases/v1/epics/E07-audit-pipeline/tasks/T006-audit-review-state.md`; `src/tui/state/reducer.ts`; `src/tui/state/app-reducer.ts`; `src/tui/state/view-state.ts`; `test/tui/state/reducer.test.ts`; `test/tui/state/app-reducer.test.ts`; `src/audit/proposals.ts`; `src/audit/apply-proposals.ts`; `.savepoint/metrics/context-bench.md`
- Estimated input tokens: ~11,975
- Notes: E06 state/test conventions for reducer shape and deterministic testing; E07 audit proposal/apply types for state model alignment.
