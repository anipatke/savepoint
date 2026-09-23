---
id: T-024
title: Check the project schema without the migration engine
objective: O-021
planned_by: {role: planner, session: o021-design-20260923}
status: planned
complexity_tier: medium
complexity_reason: "Touches every runtime entry point (board, resume, doctor, upgrade-assets) and removes a Next rung, but each change replaces a call with a smaller one behind an existing interface."
depends_on: [{task: T-023, requires: clear}]
owner_validation:
    required: false
    accepted_check: ""
---

# T-024: Check the project schema without the migration engine

## Outcome

`board`, `resume`, `doctor`, `init`, and `upgrade-assets` decide whether a
project is usable through one small `internal/data` check instead of the
migration package's cutover preflight. A V1 project gets one clear message
to run `savepoint migrate`. A V2 project loads once instead of twice.
Only `main.go`'s `migrate` dispatch imports `internal/migrate`.

## User Check

On a V2 project, `savepoint board`, `resume`, and `doctor` behave as
before. On a copy of the `v1-basic` fixture, each prints the same
"run `savepoint migrate`" guidance and exits nonzero with its current exit
code. A malformed `schema_version` still gets a named diagnostic.

## Done When

- `internal/data` exports the project-root helpers the runtime needs
  (moved from `migrate.FindProjectRoot` and `migrate.ResolveTarget`, with
  their `ErrTargetMissing`/`ErrTargetNotSavepoint` errors) and one
  runtime schema check built on `ReadSchemaVersion`. It returns `nil` for
  schema 2, and named errors for schema 1 ("schema_version 1: run
  `savepoint migrate --dry-run`, then `savepoint migrate --apply`"),
  unsupported, and malformed/unreadable config.
- `main.go` (`runResume`, `runDoctorChecks`), `internal/board/board.go`,
  `internal/board/v2/load.go`, `internal/doctor/checks.go`, and
  `internal/init/upgrade.go` use that check. None of them imports
  `internal/migrate`; `go list -deps` evidence proves it.
- The runtime no longer loads the project inside a preflight and then
  again for rendering. The discarded `ResolveReleaseCutover` evaluation on
  every startup is gone from runtime paths.
- `data.MigrationState`, `NextPendingMigration`, the board's pending-
  migration screen, and resume's pending-migration phrasing are removed.
- Exit codes and stdout/stderr placement for V1, malformed, and V2 cases
  match current behavior, proven by `main_test.go`, `cmd/*_test.go`, and
  board/doctor tests. The V1 message no longer lists plan conflicts or
  ambiguities; it points at `migrate --dry-run`, which reports them.
- `git diff --check`, `make build && make test-fast` pass.

## Context Files

`main.go`, `main_test.go`, `internal/board/board.go`,
`internal/board/board_test.go`, `internal/board/v2/load.go`,
`internal/board/v2/load_test.go`, `internal/board/v2/next_panel.go`,
`internal/board/v2/next_panel_test.go`, `internal/board/v2/view.go`,
`internal/doctor/checks.go`, `internal/doctor/v2_runtime.go`,
`internal/doctor/v2_runtime_test.go`, `internal/init/upgrade.go`,
`internal/init/upgrade_schema_test.go`, `internal/data/config.go`,
`internal/data/next.go`, `internal/data/next_test.go`,
`internal/resume/resume.go`, `internal/resume/resume_test.go`,
`internal/migrate/command.go`, `internal/migrate/cutover.go`,
`internal/migrate/operation.go`.

## Design References

Design sections 2, 4 (runtime load), 8 (Next), and 10 (commands).

## Guardrails

ARCH-01, ARCH-03, ARCH-04, DATA-02, DATA-03, FS-06, TEST-01, TEST-02,
TEST-08.

## Implementation Plan

1. Add the root helpers and schema check to `internal/data` with tests
   for V2, V1, unsupported, malformed, missing config, and missing target.
2. Switch each runtime caller; keep exit-code behavior identical.
3. Remove the pending-migration rung, state, board screen, and resume
   phrasing, with their tests.
4. Have `internal/migrate` call the moved root helpers from `data` so there
   is one implementation.
5. Prove the import graph and run `make build && make test-fast`.

## Boundaries

Do not change `internal/migrate`'s apply, journal, or cutover code beyond
switching it to the moved helpers (T-025 deletes them). No change to V2
loading rules, `ResolveNext` precedence other than removing the migration
rung, or O-014's router scope.

## Technical Verification

Focused `internal/data`, board, doctor, resume, and `main` tests;
`make build && make test-fast` at handoff.

## Technical Evidence

Pending execution.

## Drift Notes

Codebase Map change: `internal/data` gains the project-root and schema
check responsibility; reconciled in T-026.
