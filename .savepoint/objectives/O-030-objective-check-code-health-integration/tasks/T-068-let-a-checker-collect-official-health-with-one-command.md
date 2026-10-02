---
id: T-068
title: Let a checker collect official health with one command
objective: O-030
status: done
depends_on: [{task: T-066, requires: clear}]
owner_validation:
    required: false
    accepted_check: ""
planned_by: {role: planner, session: planning-o030-20261001}
complexity_tier: medium
complexity_reason: New CLI subcommand that runs external tools through existing collection and must stay thin and platform-safe.
check_waiver:
    task: T-068
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-01T19:57:40Z"
---

# Let a checker collect official health with one command

## Outcome

`savepoint health check O-### [dir]` collects one official Code Health snapshot using the project's confirmed tools, reusing the reports the full gate already wrote, and prints the snapshot ID with the plain blocking verdict for the Check record.

## User Check

None beyond the Full Objective Check.

## Done When

- `health check` parses one required `O-###` and an optional directory; a missing or malformed ID, extra arguments, and unknown flags are named errors, and `health --help` lists both subcommands.
- The named Objective must exist in the strictly loaded V2 index; otherwise the command fails before running any tool.
- With no health configuration, the command prints that Code Health is not configured for this project, saves nothing, and exits 0.
- With configuration, it runs `codehealth.Collect` with `OriginOfficial` and `DefaultReaders()`, then prints the Objective, the snapshot ID, whether a new snapshot was created, and the rendered verdict from the gate. It exits 0 whenever a snapshot was saved, including when the verdict blocks; it exits non-zero only when no snapshot was saved.
- There is no flag or path that produces a manual snapshot.
- Ctrl-C cancels collection through the context.
- Tests with a fake runner and temporary Git projects cover: not configured, a blocking verdict, a non-blocking verdict, an unknown Objective, and a collection error. They show that test and coverage instances read existing reports and that no test tool is executed.
- `main.go` stays thin dispatch (ARCH-01), and the help text lists the command.

## Context Files

`.savepoint/objectives/O-030-objective-check-code-health-integration/Objective.md`; `cmd/health.go`, `cmd/health_test.go`; `main.go`, `main_health_test.go`; `internal/codehealth/collect.go`, `internal/codehealth/readers.go`, `internal/codehealth/storage.go`, `internal/codehealth/gate.go`; `internal/data/v2_index.go` (or the index loader `main.go` already uses for resume).

## Design References

Design section 6 (CLI surface) and section 1 (Code Health readers); O-030 Confirmed Design Decisions.

## Guardrails

ARCH-01, ARCH-03, FS-05, FS-06, CFG-01, CFG-02, CFG-03, TEST-01, TEST-02, TEST-04, TEST-08.

## Implementation Plan

1. Confirm the gate evaluation and renderer from T-066 exist; return REPLAN REQUIRED if not.
2. Extend argument parsing to a subcommand switch with a typed options value for `check`.
3. Add the production runner in `main.go`: resolve the root, strict-load the index, check the Objective, load the config, collect, evaluate, and print.
4. Add tests through the injected runner pattern used by `health setup`.

## Boundaries

No manual refresh, no pruning, no Check record writes, no skill or AGENTS wording, and no change to `Collect` or the readers.

## Technical Verification

It runs external processes and handles paths, so it is platform-sensitive: fresh `make test-full` at handoff.

## Technical Evidence

Per-criterion outcomes:

- Parsing: `health check` needs one `O-###` and an optional directory; missing/malformed ID, extra arguments, and unknown flags (including `--manual`) are named errors; `health --help` and `health check --help` list both subcommands. `cmd/health_test.go` (`TestParseHealthCheckArgs`, `TestParseHealthArgsRejectsBadInput`, `TestRunHealthHelpDoesNotRun`) and `main_health_test.go`.
- Objective validation first: `internal/healthcheck.Run` strict-loads the V2 index and rejects an unknown Objective before reading config or running any tool. `TestRun_unknownObjectiveRunsNothing`.
- Not configured: prints that Code Health is not configured, saves nothing, creates no health directory, returns nil (exit 0). `TestRun_notConfiguredSavesNothing`.
- Configured: `codehealth.Collect` with `OriginOfficial` and `DefaultReaders()`, then prints Objective, snapshot ID, created/already-stored, and `Verdict.Render()`. A blocking verdict returns nil; only a no-snapshot failure returns an error (exit 1). `TestRun_blockingVerdictStillSucceeds`, `TestRun_nonBlockingVerdict`, `TestRun_collectionErrorSavesNothing`.
- No manual path: `Request` and `HealthCheckOptions` have no origin or manual field; `TestRun_neverProducesManualSnapshot`.
- Ctrl-C: `main.go` runs `health` under `signal.NotifyContext(os.Interrupt)`, and Collect honours the context; `TestRun_cancelledContextSavesNothing` covers a cancelled context. Real signal delivery is not exercised by a test.
- Fake runner and temporary Git projects: every test asserts the only executed tool is `lizard`, so the go-test and coverage instances read existing reports and no test tool runs.
- Thin main: `main.go` only dispatches (`healthCheckRunner` is one call); behaviour is in `internal/healthcheck`. Help lists the command (`TestMainHelpListsHealthCheck`).

Commands run: `go test ./cmd ./internal/healthcheck`, `make build`, `make test-full` (fresh, platform-sensitive handoff; all passed, including linux/darwin/windows cross-builds).

Files changed: `cmd/health.go`, `cmd/health_test.go`, `main.go`, `main_health_test.go`, new `internal/healthcheck/healthcheck.go`, `internal/healthcheck/healthcheck_test.go`.

Deviation and extra reads, recorded for transparency:

- ARCH-01 (and the `cmd` import test that allows only context/errors/fmt/io) forbids behaviour in `cmd` or `main.go`, so the collection logic lives in a new package `internal/healthcheck` that is not in the Context Files. The plan's "production runner in main.go" is a one-line wrapper. The Codebase Map needs a line for it in the skill-and-docs Task (ARCH-04).
- `RunHealth` now takes a `HealthRunners` struct (Setup and Check) instead of one function, and `ParseHealthArgs` returns a `HealthInvocation`; the existing setup tests were adapted, with unchanged setup behaviour.
- Extra reads: `internal/codehealth/runner.go`, `config.go`, `collect_test.go`, `repository_test.go`, `testdata/readers/*` fixtures, `cmd/resume_test.go` (import rule), `internal/doctor/checks_test.go` (project fixture shape), `internal/data/project.go` (`LoadV2Index`).

Limitations: the not-configured and collection-error paths use real temporary Git repos but no real tools. A repeated run's "already stored" label is not tested because timestamps make each snapshot distinct. No Task Check requested and no owner waiver recorded.

## Drift Notes

The CLI table, CLI Rules exception, and Codebase Map are reconciled in the skill-and-docs Task.
