---
id: E07-audit-pipeline/T004-audit-orchestration-router
status: done
objective: Wire `savepoint audit` to run gates, write audit artifacts, and move the router into audit-pending only after blocking gates pass.
depends_on:
  - E07-audit-pipeline/T001-audit-cli-contract
  - E07-audit-pipeline/T002-quality-gate-runner
  - E07-audit-pipeline/T003-snapshot-and-prompt
---

# T004: audit-orchestration-router

## Implementation Plan

- [x] Replace the non-skip audit placeholder with orchestration that resolves the target epic, reads config, runs quality gates, and reports gate results.
- [x] Stop before snapshot generation when any blocking quality gate fails, preserving clear stdout/stderr and exit-code behavior.
- [x] Generate the snapshot, proposal directory, and handoff prompt after gates pass.
- [x] Add `src/audit/router-state.ts` to transition `.savepoint/router.md` into `audit-pending` for the audited epic without disturbing unrelated router content.
- [x] Write audit log entries for started, failed-gate, and snapshot-created outcomes.
- [x] Add integration tests for happy path, failed gates stopping snapshots, explicit `--epic`, active-router epic fallback, and router transition.
- [x] Run the focused audit command integration tests before handing off the task.

## Context Log

- Files read: `.savepoint/router.md`; `E07/Design.md`; `T004`; `T001`; `T002`; `T003`; `src/commands/audit.ts`; `src/audit/log.ts`; `src/audit/quality-gates.ts`; `src/audit/snapshot.ts`; `src/audit/prompts.ts`; `src/domain/config.ts`; `src/readers/config.ts`; `src/readers/tasks.ts`; `src/domain/router.ts`; `src/readers/router.ts`; `src/fs/project.ts`; `src/cli/exit-codes.ts`; `test/audit/quality-gates.test.ts`; `test/audit/snapshot.test.ts` (partial); `test/commands/audit.test.ts`; `.savepoint/metrics/context-bench.md`
- Estimated input tokens: ~15,300
- Notes: Created `src/audit/router-state.ts` with `transitionToAuditPending`. Added `logAuditEntry` to `src/audit/log.ts`. Replaced placeholder in `src/commands/audit.ts` with full orchestration (gates → snapshot → router transition → prompt). Added injectable `commandRunner` and `gitReader` to `AuditContext` for testing. 20 audit command tests, 564 total passing. Typecheck and lint clean.
