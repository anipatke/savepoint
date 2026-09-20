---
id: E07-audit-pipeline/T002-quality-gate-runner
status: done
objective: Run configured audit quality gates with deterministic pass/fail results and no snapshot side effects.
depends_on: []
---

# T002: quality-gate-runner

## Implementation Plan

- [x] Introduce `src/audit/quality-gates.ts` with a small typed model for configured lint, typecheck, test, and future gate commands.
- [x] Implement an injectable command runner that executes gates in project-root context and captures exit code, stdout, stderr, duration, and spawn failures.
- [x] Normalize empty or missing gate configuration into the project defaults without hard-coding command execution inside the CLI command.
- [x] Return an aggregate result that clearly distinguishes pass, blocking failure, and runner error so orchestration can stop before snapshot generation.
- [x] Add unit tests for all-pass gates, nonzero failures, spawn errors, empty gate lists, output capture, and deterministic execution order.
- [x] Run the focused audit quality-gate tests before handing off the task.

## Context Log

- Files read: `.savepoint/router.md` (via AGENTS.md routing); `AGENTS.md`; `.savepoint/releases/v1/epics/E07-audit-pipeline/Design.md`; `T001-audit-cli-contract.md`; `T002-quality-gate-runner.md`; `src/commands/audit.ts`; `src/audit/log.ts`; `src/domain/config.ts`; `src/readers/config.ts`; `src/fs/project.ts`; `src/cli/args.ts`; `test/commands/audit.test.ts`
- Estimated input tokens: ~8,500
- Notes: Created `src/audit/quality-gates.ts` with `Gate`, `GateResult`, `QualityGateOutcome`, `CommandRunner` types, `normalizeGates()`, `defaultRunner()`, and `runQualityGates()`. Created `test/audit/quality-gates.test.ts` with 13 tests covering all-pass, blocking failure, non-blocking failure, spawn errors, empty gates, output capture, duration recording, and deterministic order. All 13 tests pass.
