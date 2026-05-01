---
id: E07-audit-pipeline/T001-audit-cli-contract
status: done
objective: Define the audit command contract, explicit epic selection, and skip-with-reason logging path without running the full audit pipeline.
depends_on: []
---

# T001: audit-cli-contract

## Implementation Plan

- [x] Extend CLI parsing and help text for `savepoint audit --epic <epic>`, `--skip`, and `--reason <text>` while rejecting invalid flag combinations.
- [x] Replace the audit stub with a boundary-safe command handler that resolves the Savepoint project root and active epic from `--epic` or the router state.
- [x] Add audit path and log helpers that can create `.savepoint/audit/{epic}/` and record skipped audits with timestamp, epic, and reason.
- [x] Keep the non-skip path explicit and non-destructive until the pipeline tasks are available.
- [x] Add focused CLI tests for help output, option parsing, explicit epic selection, skip success, missing reason, and invalid reason usage.
- [x] Run the focused audit CLI tests before handing off the task.

## Context Log

- Files read: `.savepoint/router.md`; `.savepoint/releases/v1/epics/E07-audit-pipeline/Design.md`; `.savepoint/releases/v1/epics/E07-audit-pipeline/tasks/T001-audit-cli-contract.md`; `src/cli/args.ts`; `src/cli/help.ts`; `src/cli/run.ts`; `src/cli/exit-codes.ts`; `src/commands/audit.ts`; `src/commands/board.ts`; `src/commands/result.ts`; `src/domain/router.ts`; `src/fs/project.ts`; `src/fs/markdown.ts`; `src/readers/router.ts`; `test/cli/args.test.ts`; `test/cli/run.test.ts`; `test/cli/commands.test.ts`; `test/commands/board.test.ts`
- Estimated input tokens: ~18,000
- Notes: Added `AuditOptions` interface and `parseAuditArgs()` to `src/cli/args.ts`; audit now has its own `"audit-command"` and `"audit-invalid-options"` ParseResult kinds. Created `src/audit/log.ts` for skip logging. Replaced audit stub in `src/commands/audit.ts` with async `runAudit(ctx: AuditContext)`. Updated `run.ts` to dispatch audit-command and audit-invalid-options. 500 tests passing, 0 failures.
