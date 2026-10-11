---
id: T-138
title: Remove the unused run command
objective: O-048
status: planned
depends_on: []
owner_validation: {required: false}
planned_by: {role: planner, session: planning-o048-animate-tiles-20261011}
---

# Remove the unused run command

## Outcome

The `run` subcommand, the `internal/run` package, and doctor's lost-runner report, all added by O-045 in commit ced16e7, are removed. `agent_launcher` returns to the parsed-but-unused config block it was before O-045, and existing configs that set it still load.

## User Check

Not required. Owner may run the built `savepoint` with no arguments and confirm `run` is no longer listed.

## Done When

1. Deleted: `cmd/run.go`, `cmd/run_test.go`, `internal/run/` (all files and `testdata`), `internal/doctor/lost_runs.go`, and `internal/doctor/lost_runs_test.go`. The `run` dispatch and usage in `main.go`, and the lost-runs call in `internal/doctor/v2_runtime.go`, are removed. Nothing else imports `internal/run`.
2. `internal/data/launcher_config.go` and its test go back to their pre-ced16e7 behavior for anything only the run path used (`DefaultAgentBuilder`, `RunTimeout`, the headless defaults). Validation keeps accepting every config that loaded before ced16e7 and every config this repository and the scaffold ship (CFG-01, DATA-03). If a later non-run caller now depends on one of those additions, keep it and record why.
3. `internal/init/claude_guard_test.go`: the three `run` command cases go. `TestGuard_repositoryCopyMatchesTheScaffold` stays, because it guards the hook copy and not the run command.
4. AGENTS.md: the Codebase Map rows for `internal/run` and the doctor lost-runner sentence are removed or reworded. TPL-02 applies to shipped guidance: scan `templates/` and `agent-skills/` for mentions of the run command, and correct any found.
5. `.savepoint/runs/` is not created by any remaining code. `docs/reviews/agent-orchestration-*.md` and O-045's records stay as history.
6. `make build && make test-fast` and fresh `make test-full` pass, because process-control code for Unix and Windows is removed. Per-criterion evidence follows AGENTS.md's Verification Policy.

## Context Files

`main.go`, `cmd/run.go`, `cmd/run_test.go`, `internal/doctor/lost_runs.go`, `internal/doctor/lost_runs_test.go`, `internal/doctor/v2_runtime.go`, `internal/data/launcher_config.go`, `internal/data/launcher_config_test.go`, `internal/init/claude_guard_test.go`, `AGENTS.md`, `.savepoint/objectives/O-048-animate-tiles/Objective.md`. The whole `internal/run/` directory is deleted, not read.

## Design References

Design §8 "Task tile run animation" decision note. O-048 Confirmed Design (run code).

## Guardrails

ARCH-04, CFG-01, CFG-03, DATA-03, TPL-02, TEST-01, TEST-06, TEST-08.

## Implementation Plan

1. `git show ced16e7 --stat` to confirm the full list of O-045 code changes, and record it in the evidence.
2. Delete the run files and remove the dispatch, doctor hook, and guard cases.
3. Trim the launcher config back, keeping existing configs loading.
4. Update AGENTS.md and any shipped guidance.
5. Gates.

## Boundaries

No change to the board, `ResolveNext`, gates, or skills. Do not delete O-045's records, Checks, or Issues, or the review documents. No Task, Objective, Check, or router write.

## Technical Verification

Focused tests during iteration, then fresh `make test-full` at handoff. The native Windows CI job for the O-048 Full Objective Check covers the removal of the Windows process code. Evaluated later under `agent-skills/references/check-method.md`.

## Technical Evidence

Pending execution.

## Drift Notes

Record architecture deltas and reconcile through the planner before Check.
