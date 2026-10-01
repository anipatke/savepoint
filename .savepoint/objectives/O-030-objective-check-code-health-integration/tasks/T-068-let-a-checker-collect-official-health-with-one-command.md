---
id: T-068
title: Let a checker collect official health with one command
objective: O-030
status: planned
depends_on: [{task: T-066, requires: clear}]
owner_validation: {required: false}
planned_by: {role: planner, session: planning-o030-20261001}
complexity_tier: medium
complexity_reason: New CLI subcommand that runs external tools through existing collection and must stay thin and platform-safe.
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

Pending execution.

## Drift Notes

The CLI table, CLI Rules exception, and Codebase Map are reconciled in the skill-and-docs Task.
