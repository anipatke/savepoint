---
id: T-083
title: Keep the report current and add a command to rewrite it
objective: O-036
status: done
depends_on: [{task: T-082, requires: clear}]
owner_validation:
    required: true
    accepted_check: ""
planned_by: {role: planner, session: planning-o036-20261002}
check_waiver:
    task: T-083
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-02T02:14:38Z"
---

# Keep the report current and add a command to rewrite it

## Outcome

The report file is rewritten whenever a snapshot is saved, and `savepoint health report` rewrites it on demand from the newest saved snapshot without running any tool.

## User Check

Run `savepoint health check O-036`, then open `.savepoint/health/report.md`: it describes the snapshot just saved. Delete the file and run `savepoint health report`: it comes back with the same content and no new snapshot appears. In a project with no snapshot the command says so and writes nothing. Press `R` in the popover and let it finish: the file is rewritten again.

## Done When

- Saving a snapshot from a check or a popover refresh rewrites the report for it, from the one place both already go through. A failure to write the report is a warning only: the snapshot and the check's verdict are unchanged and the command still succeeds.
- `savepoint health report [dir]` renders and writes the report from the newest saved snapshot, prints the file path, creates no snapshot and runs no tool. With no snapshot, or Code Health not set up, it prints a plain message, writes nothing and exits non-zero. It takes no Objective argument; the re-run line uses the router's Objective, or `O-###`.
- The command appears in the usage text and `--help`, and the argument parser rejects unknown flags and extra arguments with the usage line.
- Tests cover the rewrite after a check and after a refresh, the command with a snapshot, with none, with Code Health not set up, with an unwritable directory (warning on save, error on the command), the usage text, and that no snapshot is created by the command.

## Context Files

`.savepoint/objectives/O-036-hand-code-health-to-an-agent/Objective.md`; `internal/codehealth/collect.go`, `internal/codehealth/collect_test.go`, `internal/codehealth/report.go` (read only), `internal/codehealth/storage.go` (read only); `cmd/health.go`, `cmd/health_test.go`; `main.go`, `main_health_test.go`; `internal/healthcheck/healthcheck.go` (read only).

## Design References

O-036 Confirmed Design Decisions; Design section 6 (CLI surface).

## Guardrails

ARCH-01, ARCH-03, FS-04, FS-06, CFG-01, STYLE-05, TEST-01, TEST-02, TEST-04.

## Implementation Plan

1. Confirm the report renderer and store write method from the previous Task exist; return REPLAN REQUIRED otherwise.
2. Call the store write after the snapshot is saved in the collection path; report a failure as a warning.
3. Add the `report` subcommand to the health parser and runners, thin, with the work in `codehealth`.
4. Document it in the usage text, then add the tests in Done When.

## Boundaries

No popover change and no change to what a check measures or when snapshots are created. The command takes no Objective argument.

## Technical Verification

Focused tests during iteration; `make build && make test-fast` for ordinary handoff. Later evaluated per `agent-skills/references/check-method.md`.

## Technical Evidence

Executor handoff (not a Check; no waiver recorded). Commands: `make build && make test-fast` (ok); focused `go test ./cmd ./internal/healthcheck ./internal/codehealth . ./internal/board/...` (ok).

Per criterion:
- Rewrite on save: `Collect` (the one path used by `health check` and the popover refresh) calls the new `RefreshReport` after `SaveSnapshot`; failure lands in `Collection.ReportErr`, never an error. `healthcheck.Run` prints it to stderr as a warning. `CollectRequest.Objective` carries the re-run Objective (check passes its own, popover passes the router's). Tests: `TestCollectRewritesReportForTheSavedSnapshot`, `TestCollectReportWriteFailureKeepsTheSnapshot`, `TestRun_rewritesReportForTheSavedSnapshot`, `TestRun_reportWriteFailureIsOnlyAWarning`.
- Command: `savepoint health report [dir]` -> `healthcheck.RunReport` -> `codehealth.RefreshReport`; prints the path, no snapshot, no tool. No snapshot or not set up returns `ErrReportNoSnapshot` / `ErrReportNotConfigured` (non-zero exit via main), writing nothing. Objective from router, else `O-###`. Tests: `TestRunReport_rewritesDeletedReportWithoutNewSnapshot`, `TestRunReport_nothingToReportWritesNothing`, `TestRunReport_unwritableReportIsAnError`, `TestRunReport_namesTheRoutersObjective`, `TestRefreshReportCreatesNoSnapshot`.
- Usage and parser: usage text, `--help`, and main help list the command; unknown flags and extra arguments are rejected with the usage line. Tests: `TestParseHealthReportArgs`, `TestParseHealthReportRejectionShowsUsage`, `TestRunHealthDispatchesReport`, `TestParseHealthArgsRejectsBadInput`, `TestMainHelpListsHealthSetup`.

Files read: Context Files, plus `AGENTS.md`, `agent-skills/savepoint-task/SKILL.md`, `internal/data/router_v2.go`, `internal/board/v2/io.go` (extra reads, for the router Objective and the popover refresh path). Files changed: `internal/codehealth/collect.go`, `collect_test.go`, `internal/healthcheck/healthcheck.go`, `healthcheck_test.go`, `cmd/health.go`, `cmd/health_test.go`, `main.go`, `main_health_test.go`, `internal/board/v2/io.go` (passes the router Objective; no popover change), `AGENTS.md` and `.savepoint/Design.md` (command documented as human-only).

Limitations: no dedicated test drives the board's `refreshHealth`; the rewrite after a popover refresh is covered through `Collect`, which it calls. A report warning during a popover refresh is not shown (the popover has no channel for it). `make test-full` not run. The tree also holds uncommitted work from earlier Tasks.

## Drift Notes

New module or architecture delta beyond the documented Codebase Map, reconciled through the planner before Check.
