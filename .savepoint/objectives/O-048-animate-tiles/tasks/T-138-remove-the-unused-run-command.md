---
id: T-138
title: Remove the unused run command
objective: O-048
status: done
depends_on: []
owner_validation:
    required: false
    accepted_check: ""
planned_by: {role: planner, session: planning-o048-animate-tiles-20261011}
check_waiver:
    task: T-138
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-11T05:38:57Z"
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

Executor evidence (claims for a fresh Check to verify).

Scope confirmed with `git show ced16e7 --stat`: O-045 code changes were `cmd/run.go`, `cmd/run_test.go`, `internal/run/*` (incl. testdata), `internal/doctor/lost_runs*.go`, `internal/doctor/v2_runtime.go`, `internal/data/launcher_config*.go`, `internal/init/claude_guard_test.go`, `main.go`, `AGENTS.md`. Everything else in that commit is records/docs and was left alone.

Per criterion:
1. Met. The listed files and `internal/run/` deleted (`git rm`); `run` case, usage line, `runRunner` and the `internal/run` import removed from `main.go`; `lostRunProblems` call removed from `v2_runtime.go`. `go build ./...` succeeds, so nothing else imports `internal/run`.
2. Met. `launcher_config.go` and its test restored to their `ced16e7^` content (removes `DefaultAgentBuilder`, `RunTimeout`, headless defaults, `Timeout` field). No remaining caller of them (grep over templates, agent-skills, cmd, internal, main.go). Restored code is what every config loaded under before ced16e7; repo and scaffold config ship `agent_launcher` only as a comment, so they still load.
3. Met. Three `run` guard cases removed from `claude_guard_test.go`; `TestGuard_repositoryCopyMatchesTheScaffold` kept, its comment reworded to drop the run-command mention.
4. Met. AGENTS.md: `run` dropped from the `cmd/` row and the `internal/run/` row removed. The doctor lost-runner sentence lived only in that row. grep of `templates/` and `agent-skills/` for the run command found no mentions (TPL-02).
5. Met. No remaining code references `.savepoint/runs` (grep). `docs/reviews/*` and O-045 records untouched.
6. Met. `make build && make test-fast` passed; fresh `make test-full` exited 0 after the final edit.

Extra reads (beyond Context Files): `git show ced16e7` diff of main.go, v2_runtime.go, claude_guard_test.go, AGENTS.md, launcher_config.go, to restore exact prior state. `.savepoint/config.yml` and the scaffold config grep, to confirm configs still load.

Limitations: Windows process-control removal is covered only by the O-048 Full Objective Check's native Windows CI job, not run here. Running the built program with no arguments to eyeball the usage list was blocked by the agent command guard, so "run no longer listed" is verified by the diff and tests, not by that run. The working tree also holds unrelated uncommitted T-137 board changes, which this Task did not touch.

## Drift Notes

Record architecture deltas and reconcile through the planner before Check.
