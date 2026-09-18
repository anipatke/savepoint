---
id: E45-safe-migration/T011-run-migration-from-the-command-line
title: Run migration from the command line
status: in_progress
stage: build
objective: Add a thin migrate command that previews by default, applies only when asked, and reports named errors.
depends_on:
    - E45-safe-migration/T009-publish-the-conversion-once
complexity_tier: low
complexity_reason: Argument parsing and dispatch over an existing package, following the established command pattern.
---

# T011: Run migration from the command line

## Problem

The whole operation is unreachable without a command, and the command is where the safest default has to be chosen once.

`savepoint migrate` should preview. A user who types a command they have not used before, on a project they care about, should get a report rather than a converted project — and the flag they have to add to change that is the moment they confirm intent. `v2-Design.md` spells the preview as `--dry-run`, which stays supported as an explicit way to say the default out loud, but the default itself is what protects the case that matters.

Everything else here is the existing pattern: `cmd/` parses arguments and dispatches, behavior lives in `internal/`, and `main.go` gains one case.

## Context Files

- `cmd/migrate.go`
- `cmd/migrate_test.go`
- `cmd/doctor.go`
- `cmd/upgrade-assets.go`
- `main.go`
- `main_test.go`
- `internal/migrate/plan.go`
- `internal/migrate/apply.go`
- `internal/migrate/operation.go`
- `README.md`

## Acceptance Criteria

- [x] `savepoint migrate [dir]` previews and writes nothing; a test snapshots bytes and modification times across the command and asserts nothing changed.
- [x] `--apply` is required to write; `--dry-run` is accepted as an explicit synonym of the default; passing both is accepted and previews; the help text states that preview is the default.
- [x] `--decisions FILE` supplies the owner decision input, and a plan with unresolved blocking ambiguities exits nonzero naming every unresolved ID.
- [x] `--recover` reports an incomplete operation and its recovery guidance, and resumes it when `--apply` is also given.
- [x] A missing directory, an unwritable directory, and a directory that is not a Savepoint project each produce a distinct named error and a nonzero exit, with no partial write, satisfying FS-06.
- [x] An unknown flag produces a named error and the usage text rather than being ignored.
- [x] The preview output enumerates the planned records, archives, identity mappings, conflicts, and decisions required — the plan itself, not a summary that omits entries.
- [x] Output is deterministic and readable without a TTY, matching the existing non-interactive command behavior.
- [x] `cmd/migrate.go` contains argument parsing and dispatch only, with no frontmatter parsing, no conversion rule, and no filesystem walking, keeping ARCH-01 intact.
- [x] `main.go` gains one `migrate` case following the existing dispatch shape, and `--version` and every existing command behave unchanged.
- [x] `README.md` documents the command, its default, and the fact that preview writes nothing.

## Implementation Plan

- [x] Add `cmd/migrate.go` following the argument-parsing and runner-injection shape used by `cmd/doctor.go` and `cmd/upgrade-assets.go`.
- [x] Implement flag handling for `--apply`, `--dry-run`, `--decisions`, and `--recover`, with the preview default and named errors for unknown flags.
- [x] Render the preview from the plan value, keeping the rendering in `internal/migrate` if it needs any domain knowledge.
- [x] Wire the `migrate` case into `main.go`.
- [x] Document the command in `README.md`.
- [x] Test the write-free default, each flag, each target-directory failure, unknown flags, deterministic non-TTY output, and unchanged behavior of the existing commands.
- [x] Run the focused `cmd` and `internal/migrate` suites, then `make build && make test`.

## Context Log

Files read: `cmd/migrate.go`, `cmd/migrate_test.go`, `cmd/doctor.go`, `cmd/upgrade-assets.go`, `main.go`, `main_test.go`, `internal/migrate/plan.go`, `internal/migrate/apply.go`, `internal/migrate/operation.go`, `README.md`.

Files edited:
- `cmd/migrate.go` (new): argument parsing and dispatch for `migrate` — `--apply`, `--dry-run`, `--decisions`, `--recover`, `--help`, unknown-flag and multi-directory rejection.
- `cmd/migrate_test.go` (new): flag-parsing coverage for every option and combination, help text, and error/exit-code cases.
- `internal/migrate/preview.go` (new): `FormatPreview` renders a `ConversionPlan` deterministically — planned records, documents, archives, legacy prerequisites, waived references, ambiguities, and conflicts, each sorted independent of build order.
- `main.go`: added the `migrate` dispatch case, `resolveMigrateTarget` (missing/unwritable/not-a-Savepoint-project diagnostics per FS-06), `newMigrationOperationID`, and `migrateRunner` wiring `internal/migrate` (`PendingOperation`, `ReadDecisionsFile`, `Plan`, `Apply`) behind the runner-injection pattern used by `doctor`/`upgrade-assets`.
- `main_test.go`: end-to-end subprocess coverage against the `v1-basic` fixture — write-free preview (byte+mtime snapshot), `--dry-run` and `--apply --dry-run` synonyms, apply that activates `schema_version: 2`, deterministic preview across two runs, ambiguity blocking and `--decisions` resolution, `--recover` with and without a pending operation, all three target-directory failure modes, unknown-flag rejection, and an `upgrade-assets` regression guard.
- `README.md`: added a "Migrating a V1 Project (`savepoint migrate`)" section documenting the preview-by-default behavior and each flag.

Quality gates: `go build ./...` clean; `make build && make test` clean (all packages, including `cmd` and `internal/migrate`).
