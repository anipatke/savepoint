---
id: E45-safe-migration/T011-run-migration-from-the-command-line
title: Run migration from the command line
status: done
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

**Files read:** `cmd/doctor.go`, `cmd/upgrade-assets.go`, `cmd/init.go`, `main.go`, `main_test.go`, `internal/migrate/plan.go`, `internal/migrate/apply.go`, `internal/migrate/operation.go`, `internal/migrate/decisions.go`, `internal/migrate/manifest.go`, `internal/migrate/classify.go`, `internal/migrate/fixture_test.go`, `internal/migrate/apply_test.go`, `internal/migrate/decisions_test.go`, `internal/init/validate.go`, `internal/data/discover.go` (`FindSavepointRoot`), `internal/data/config.go` (`ReadSchemaVersion`), `README.md`, `.savepoint/Guardrails.md`, `E45-Detail.md`.

**Files added:**
- `cmd/migrate.go`: `MigrateOptions` (`Dir`, `Apply`, `DryRun`, `DecisionsFile`, `Recover`), `WillWrite()` (apply-without-dry-run), `ParseMigrateArgs`, `RunMigrate` — pure argument parsing and dispatch, no filesystem or domain imports, matching `cmd/doctor.go`'s `(int, error)` runner shape. Unknown-flag and parse-error paths print `migrateUsage` to stdout before returning exit code 2, satisfying the "usage text rather than silently ignored" AC.
- `cmd/migrate_test.go`: help text, every flag combination (including `--apply --dry-run` together resolving to preview), `--decisions` value parsing and its missing-value error, unknown-flag/multiple-directory rejection with usage text asserted, and runner-code passthrough — all against the injected fake runner, no real filesystem or `internal/migrate` interaction (matching the existing `cmd` test style).
- `internal/migrate/preview.go`: `FormatPreview(plan *ConversionPlan) string` — the domain-aware rendering the Implementation Plan called for, kept in `internal/migrate` per ARC-01. Renders the already-computed plan value only (no IO): planned records with their legacy-to-global identity mapping, documents, archives (with archived-`Health-Check.md` candidate commands surfaced for owner review only), legacy prerequisites, waived references, conflicts, and every ambiguity with its blocking/advisory status, detail, choices, and resolution — sorted defensively by stable keys so the report's structure never depends on `Plan`'s internal build order.

**Files edited:**
- `main.go`: added the `migrate` dispatch case (mirrors `doctor`'s `(code, err)` shape); `resolveMigrateTarget` (three distinct named errors — `ErrMigrateTargetMissing`, `ErrMigrateTargetNotSavepoint`, `ErrMigrateTargetUnwritable` — checking `dir` itself rather than walking up parents the way `FindSavepointRoot` does, so migrate can never silently operate on an unrelated ancestor project); `newMigrationOperationID` (the one production `migrate.OperationIDSource`, timestamp + `crypto/rand` suffix — tests inject their own deterministic source, this is the real one); `migrateRunner`, which owns every filesystem interaction and every call into `internal/migrate` (`PendingOperation`, `ReadDecisionsFile`, `Plan`, `Apply`), keeping `cmd/migrate.go` itself free of all of it. Exit codes: `2` for parse/internal errors (mirroring `doctor`), `1` for named refusals (bad target, plan conflict, unresolved blocking ambiguities, `Apply` error), `0` for a successful preview or apply. `--decisions` is read before `Plan` so an unresolved-ID or disallowed-value error from `ReadDecisionsFile`/`validateDecisions` surfaces as a normal named error. `--recover` prints `PendingOperationReport.RecoveryGuidance()` (or "no incomplete migration operation found") and returns immediately unless `--apply` was also given, in which case it falls through to the normal `Plan`/`Apply` call — `Apply` already auto-detects and resumes the pending operation via its own `PendingOperation` check, so no separate resume code path was needed.
- `main_test.go`: added `copyMigrateFixture` (copies the frozen `internal/data/testdata/migration/v1-basic/project` fixture into a temp dir, mirroring `internal/migrate/apply_test.go`'s `copyFixtureProject`), `snapshotDir`/`assertSameSnapshot` (content-hash + mtime per file, for the write-free proof), `writeMigrateMinimalProject`/`writeMigrateAmbiguousProject` (a minimal V1 project whose one task carries an unrecognized status, mirroring `internal/migrate/decisions_test.go`'s `writeMinimalV1Project`/`writeTaskWithStatus`, to exercise the unresolved-blocking-ambiguity exit path end to end), and one test per AC (listed below). Also added `TestMainUpgradeAssetsStillWorksAfterMigrateAdded` as an explicit regression guard for the "every existing command behaves unchanged" AC, on top of the pre-existing `--version`/`init --help`/`upgrade-assets` tests continuing to pass unmodified.
- `README.md`: new "Migrating a Legacy Project (`savepoint migrate`)" section — states preview is the default and writes nothing, shows the `--apply` invocation, and documents `--decisions`, `--recover`, and the three named target-directory errors.

**Tests added (`main_test.go`, run against the real compiled binary via the existing `runMainForTest` helper):** `TestMainMigrateHelpStillUsesNormalDispatch`, `TestMainMigrateRejectsUnknownFlag`, `TestMainMigratePreviewDefaultWritesNothing` (full-tree snapshot before/after), `TestMainMigrateDryRunIsSynonymOfDefault`, `TestMainMigrateApplyAndDryRunTogetherPreviews`, `TestMainMigratePreviewIsDeterministic` (two runs, `operation:`/`generated:` lines normalized out before comparing — those two vary legitimately against the real clock and op-id source; everything else must match), `TestMainMigrateApplyWritesAndActivatesSchema`, `TestMainMigrateAmbiguousPlanBlocksAndNamesTheID`, `TestMainMigrateDecisionsFileResolvesAmbiguity`, `TestMainMigrateRecoverWithNoPendingOperation`, `TestMainMigrateRecoverReportsPendingOperation` (via a real `migrate.CreateOperation` call), `TestMainMigrateMissingDirectory`, `TestMainMigrateNotASavepointProject`, `TestMainMigrateUnwritableDirectory` (skipped on Windows/root, matching the existing `TestMainUpgradeAssetsPrintsPartialWorkOnFailure` pattern), `TestMainUpgradeAssetsStillWorksAfterMigrateAdded`.

**Quality gates:** `go build ./...`, `go vet ./...`, `gofmt -l` (clean) on every new/edited file, `go test ./cmd/... -v`, `go test . -run TestMainMigrate -v`, `go test ./... -count=1`, and `make build && make test` all pass. No `.savepoint/Health-Check.md` in this project, so the Quick check step was skipped per the build-task skill.

**Not done / out of scope:** No end-to-end test drives an actual interrupted-then-resumed apply through the `migrate` command itself — `internal/migrate`'s own `apply_test.go` (T008/T009) already covers every injected-failure and resume-convergence case against `Apply` directly, and the command layer's obligation is only that `--recover --apply` reaches `Apply`, which it does through the same code path as a plain `--apply` (see `migrateRunner`'s fallthrough above). Router priority (`p` in the TUI) was not exercised interactively; no TUI session was run in this environment.
