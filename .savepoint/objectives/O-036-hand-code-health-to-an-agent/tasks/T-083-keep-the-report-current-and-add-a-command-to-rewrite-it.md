---
id: T-083
title: Keep the report current and add a command to rewrite it
objective: O-036
status: planned
depends_on: [{task: T-082, requires: clear}]
owner_validation: {required: true}
planned_by: {role: planner, session: planning-o036-20261002}
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

Pending execution: named cases, results, reviewed source basis, files read/changed, and limitations, recorded after the work lands.

## Drift Notes

New module or architecture delta beyond the documented Codebase Map, reconciled through the planner before Check.
