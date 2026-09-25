---
id: T-024
title: Check the project schema without the migration engine
objective: O-021
planned_by: {role: planner, session: o021-design-20260923}
status: done
complexity_tier: medium
complexity_reason: "Touches every runtime entry point (board, resume, doctor, upgrade-assets) and removes a Next rung, but each change replaces a call with a smaller one behind an existing interface."
depends_on: [{task: T-023, requires: clear}]
owner_validation:
    required: false
    accepted_check: ""
check_waiver:
    task: T-024
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-09-23T11:11:35Z"
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

**Per-criterion outcomes** (against `## Done When`):

1. `internal/data` exports the project-root helpers and one runtime schema
   check — DONE. New file `internal/data/runtime.go`: `ResolveTarget`,
   `FindProjectRoot`, `ErrTargetMissing`, `ErrTargetNotSavepoint` (moved
   verbatim from `internal/migrate/command.go`), and `CheckRuntimeSchema`
   built on `ReadSchemaVersion` — returns `nil` for schema 2,
   `ErrSchemaMigrationRequired` ("schema_version 1: run `savepoint migrate
   --dry-run`, then `savepoint migrate --apply`") for schema 1, and the
   existing `ErrUnsupportedSchemaVersion`/`ErrMalformedSchemaVersion`
   sentinels (or a wrapped read error) otherwise.
2. `main.go` (`runResume`, `runDoctorChecks`), `internal/board/board.go`,
   `internal/board/v2/load.go`, `internal/doctor/checks.go`, and
   `internal/init/upgrade.go` use that check and import no `internal/migrate`
   symbol — DONE. Verified with `go list -deps` on each package: 0 hits for
   `opencode/savepoint/internal/migrate` in `internal/board`,
   `internal/board/v2`, `internal/doctor`, `internal/init`; the root `.`
   package still depends on it transitively only through `main.go`'s
   `migrate` dispatch case (`migrateRunner`), which the Task explicitly
   permits. `internal/migrate/command.go`'s `ResolveTarget`/`FindProjectRoot`
   are now one-line wrappers over the moved `internal/data` functions, and
   its `ErrTargetMissing`/`ErrTargetNotSavepoint` are aliases to
   `data`'s, so migrate's own callers (`cutover.go`, its own tests) are
   unaffected — no change to migrate's apply/journal/cutover bodies.
3. No double load; the discarded `ResolveReleaseCutover` evaluation is gone
   from runtime paths — DONE. `runResume`, `runWithFilters`, and
   `board/v2.loadProject` each call `data.LoadV2Index` exactly once (or, for
   `runWithFilters`, not at all — it hands off to the V2 board, which loads
   once). `grep -rn ResolveReleaseCutover` shows the only call site left is
   `internal/migrate/cutover.go`'s own `PreflightCutover`, and `grep -rn
   PreflightCutover` (excluding migrate's own package) shows no remaining
   caller anywhere in the runtime or its tests.
4. `data.MigrationState`, `NextPendingMigration`, the board's
   pending-migration screen, and resume's pending-migration phrasing are
   removed — DONE. Removed from `internal/data/next.go` (type, `NextKind`
   value, `Next.Migration`/`NextInput.Migration` fields, the
   `input.Migration.Pending` branch in `ResolveNext`),
   `internal/board/v2/load.go` (`ProjectState.Migration`/`MigrationGuidance`,
   the `pendingMigration` helper and its use in `loadProject`),
   `internal/board/v2/io.go` (`writeGuard` and its 6 call sites — a pending
   operation can no longer coexist with schema 2 in production, since
   `migrate.Apply` only ever activates `schema_version: 2` as its last write,
   after every journalled path is already installed; see `apply.go:124-134`),
   `internal/board/v2/next_panel.go` and `view.go`/`plain.go` (the
   `NextPendingMigration` line and `renderMigration`/`MigrationGuidance`
   rendering), `internal/resume/resume.go` (`EvidenceLines`/`ActionPhrase`
   cases). `internal/doctor`'s `CheckMigration` (which called
   `migrate.PendingOperation`) is also removed, since `checks.go` can no
   longer import `internal/migrate`; `report.go`'s `Migration` field and
   "Migration Check" section go with it.
5. Exit codes and stdout/stderr placement for V1, malformed, and V2 cases
   match current behavior; the V1 message no longer lists plan conflicts or
   ambiguities — DONE. Proven by the existing (and updated) suites in
   `main_test.go`, `main_resume_test.go`, `main_board_test.go`,
   `main_resume_matrix_test.go`, `main_board_next_parity_test.go`,
   `cmd/resume_test.go`, `cmd/board_test.go`, `cmd/doctor_test.go`, and the
   board/doctor package tests — all pass. Manually verified end-to-end
   against the real `internal/data/testdata/migration/v1-basic/project`
   fixture: `resume`, `doctor`, and `board` each print `schema_version 1: run
   \`savepoint migrate --dry-run\`, then \`savepoint migrate --apply\`` and
   exit 1, with no plan/conflict/ambiguity detail (previously sourced from a
   built `ConversionPlan`). Also verified against a fresh `init`-scaffolded
   V2 project: `resume` exits 0 with "Plan the next Objective...", `doctor`
   exits 0 "ALL CLEAN", `board` (non-TTY) exits 0 with "Nothing selected
   yet".
6. `git diff --check` and `make build && make test-fast` pass — DONE, both
   run clean (0 failures) after a `gofmt -w` pass on the two touched files
   `gofmt -l` flagged.

**Commands run:** `go build ./...`, `go vet ./...`, `git diff --check`,
`make build && make test-fast` (all clean, 0 failures), `gofmt -l .` /
`gofmt -w` on the two files it flagged that this Task touched, `go list
-deps` per-package migrate-import checks, and the manual `resume`/`doctor`/
`board` runs above against both the v1-basic fixture and a fresh V2 init
scaffold.

**Files read beyond the Task's Context Files** (all directly required to
carry the removal through to compiling, passing code — logged per the
skill's Extra Reads rule):

- `internal/migrate/apply.go` (Apply's publish order) — to prove the
  pending-operation-implies-schema-1 invariant that justifies dropping every
  runtime/board pending-migration check without a behavior regression.
- `internal/doctor/report.go` — `DiagnosticReport.Migration` and the
  "Migration Check" section had to be removed once `checks.go` could no
  longer import `internal/migrate`.
- `internal/board/v2/io.go`, `plain.go` — the board's write guard and
  non-TTY renderer both read `pendingMigration`/`state.Migration` and had to
  be updated in lockstep with `load.go`.
- `internal/board/v2/fixture_test.go`, `actions_test.go`, `watch_test.go`,
  `view_test.go`, `releases_test.go` — all built or asserted on the removed
  pending-migration board behavior (`createPendingOperation` fixture and its
  five call sites).
- `internal/doctor/checks_test.go`, `report_test.go` — direct tests of the
  removed `CheckMigration` and the "Migration Check" report section.
- `internal/init/upgrade_test.go` — direct tests of the removed
  `refusePendingMigration` write-boundary guard.
- `main_resume_matrix_test.go`, `main_board_next_parity_test.go`,
  `main_resume_test.go`, `main_board_test.go` — the top-level rung-matrix and
  board/resume parity suites each carried a pending-migration case or
  guard clause.

**Files changed:** `internal/data/runtime.go` (new), `internal/data/next.go`,
`internal/migrate/command.go`, `internal/board/board.go`,
`internal/board/v2/{load,io,next_panel,view,plain}.go`,
`internal/doctor/{checks,v2_runtime,report}.go`, `internal/init/upgrade.go`,
`internal/resume/resume.go`, `main.go`, plus the test files listed above.

**Limitations:** The `savepoint-task` Context Files list did not name several
files this removal had to touch (see Extra Reads); all were read and edited
in full, not sampled. `deadcode` was not run against the built binary as part
of this Task — the Objective's dead-code success condition is O-021-wide, not
this Task's own acceptance criterion, and T-024's own Done bullets are all
covered above. No optional Task Check was requested; this evidence is handed
to the mandatory Full Objective Check per the router.

## Drift Notes

Codebase Map change: `internal/data` gains the project-root and schema
check responsibility; reconciled in T-026.
