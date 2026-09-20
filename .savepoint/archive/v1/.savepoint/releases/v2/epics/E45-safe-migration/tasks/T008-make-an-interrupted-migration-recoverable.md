---
id: E45-safe-migration/T008-make-an-interrupted-migration-recoverable
title: Make an interrupted migration recoverable
status: done
objective: Record the operation with verified backups, staging, and a journal so any interruption leaves recoverable bytes and a truthful report.
depends_on:
    - E45-safe-migration/T001-decide-how-windows-replaces-a-file
    - E45-safe-migration/T003-plan-the-conversion-before-touching-anything
complexity_tier: high
complexity_reason: "Owns the recovery guarantee: backup verification, journal states, and interruption behavior at every step."
---

# T008: Make an interrupted migration recoverable

## Problem

Migration changes many files, and no filesystem this project targets can change many files atomically. The epic therefore promises something weaker and honest: whatever happens, the bytes are recoverable and the report is true.

That promise needs somewhere to live that the operation itself cannot destroy. A backup stored among the files being replaced is not a backup. `.savepoint/.migration/<op-id>/` is chosen because it is inside the project — so it survives with it, stays writable, and travels with a copied directory — while sitting outside every path the operation writes. It is dot-prefixed like the existing `.upgrade-manifest.yml`, and V2 discovery is already confined to `objectives/`, `checks/`, and `issues/`, so it is invisible to the index.

The journal is what makes the difference between recovery and archaeology. Without a per-path record of what was planned, backed up, staged, and installed, a half-finished migration is just a directory full of files in an unknown state.

## Context Files

- `internal/migrate/operation.go`
- `internal/migrate/operation_test.go`
- `internal/migrate/replace.go`
- `internal/migrate/plan.go`
- `internal/migrate/manifest.go`
- `internal/init/upgrade.go`
- `internal/init/manifest.go`
- `internal/init/upgrade_failure_test.go`
- `internal/data/discover.go`
- `.savepoint/Guardrails.md`

## Acceptance Criteria

- [x] An operation directory at `.savepoint/.migration/<op-id>/` holds `backup/`, `staging/`, and `operation.yml`, and no path the operation writes ever falls inside it.
- [x] The operation directory is create-only: an existing `<op-id>` is refused with a named diagnostic rather than reused or overwritten, so a prior backup or user sidecar can never be destroyed.
- [x] `operation.yml` records the operation ID, its creation time, the plan's source hashes, and for every path an action, a planned content hash, an installed hash once written, and a state among `planned`, `backed_up`, `staged`, `installed`, and `verified`.
- [x] Every file the operation will replace or remove is copied to `backup/` and each copy is re-hashed against the recorded source hash before any replacement begins; a backup that fails verification aborts before the first replacement.
- [x] Staged target files are written complete in `staging/` before any live path is touched.
- [x] The journal is written before and after each step, so a crash at any point leaves a state that names what had happened and what had not.
- [x] `PendingOperation(root)` is a read-only detector returning the incomplete operation and its recovery guidance, and returning nothing when none exists.
- [x] More than one incomplete operation produces a named diagnostic rather than a heuristic choice between them.
- [x] A test proves `LoadV2Index` ignores `.savepoint/.migration/` entirely, including an operation directory containing staged files that look like Objectives, Checks, or Issues.
- [x] Injected failures at each journal step leave every affected path recoverable from either its original location or its verified backup; a test asserts there is no step at which both the original and its backup are absent.
- [x] The failure report names the operation, the step, the affected paths, and the recovery command, and never reports a partial state as a success.
- [x] Recovery distinguishes FS-06's preflight no-partial-write expectation from a recorded interruption after a valid apply began, and says which one occurred.

## Implementation Plan

- [x] Add `operation.go` with the operation directory layout, the create-only guard, and the journal model and its states.
- [x] Implement backup capture and per-copy hash verification against the plan's recorded source hashes.
- [x] Implement staging writes through T001's replacement primitive, never through `AtomicWrite`'s truncating fallback.
- [x] Implement journal persistence around each step, flushing before the step it describes.
- [x] Implement `PendingOperation` and the recovery-guidance report, including the multiple-operation diagnostic.
- [x] Add the discovery-isolation test asserting a staged Objective-shaped file inside the operation directory is invisible to `LoadV2Index`.
- [x] Add injected-failure tests at every journal step, each asserting recoverable bytes, a truthful report, and the original-or-backup invariant.
- [x] Run the focused `internal/migrate` and `internal/data` suites, then `make build && make test`.

## Context Log

**Files read:** `internal/migrate/replace.go`, `plan.go`, `manifest.go`, `decisions.go`, `inventory.go`, `internal/data/project.go` (`LoadV2Index`, `DiscoverV2Records`/`Checks`/`Issues`), `.savepoint/Guardrails.md`, E45-Detail.md, T009 (dependent task, to confirm the operation/apply boundary).

**Files written:**
- `internal/migrate/operation.go` (new) — operation directory layout (`backup/`, `staging/`, `operation.yml`) under `.savepoint/.migration/<op-id>/`, reusing the existing unexported `migrationStateDir` constant from `inventory.go` rather than a second definition of the same path; `Journal`/`JournalEntry`/`StepState`/`EntryAction` types; `CreateOperation` (create-only, self-cleaning on partial failure); `LoadOperation`; per-path `Backup`, `WriteStaged`, `Install`, `Remove`, `Verify` methods enforcing the `planned → backed_up → staged → installed → verified` sequence in code (`requireState`); `writeCompleteFile` (temp+rename, distinct from `ReplaceFile` — used only for the operation's own staging/backup copies and for `ActionCreate` destinations, never for replacing a live file a user already has, which goes through `ReplaceFile`); `PendingOperation` and `PendingOperationReport.RecoveryGuidance()`.
- `internal/migrate/operation_test.go` (new) — layout/create-only guard tests; backup capture + hash-mismatch abort; full replace and remove lifecycle tests asserting the recoverability invariant (`assertRecoverable`) at every step, including a simulated-crash reload via `LoadOperation`; out-of-order step refusal; injected failures at staging (wrong content), install (read-only destination directory), and verify (post-install corruption); `PendingOperation` none/found/multiple cases and `RecoveryGuidance` content; `TestLoadV2Index_ignoresMigrationOperationDirEntirely` using deliberately malformed Objective/Check/Issue-shaped files staged under `.migration/` to prove `data.LoadV2Index` never reads into it.

**Design decisions / boundary notes:**
- `operation.go` is deliberately plan-agnostic: it never reads a `ConversionPlan`. Deciding which paths get which action, in what order, and schema activation is T009's `apply.go`; this task only had to prove the per-path state machine is safe when driven directly, which the tests do without any orchestration layer existing yet.
- `JournalEntry.PlannedHash` is the hash of the *new* content a path will end up with (used to verify staged/installed content); the *original* source hash backups are checked against lives in `Journal.SourceHashes`, reusing `ManifestSource` — these are two different hashes and AC3/AC4 each refer to a different one.
- Staging writes and `ActionCreate` installs go through a new `writeCompleteFile` helper (temp file + fsync + rename), not the exported `ReplaceFile`: `ReplaceFile`'s own contract refuses a missing destination by design, and staging/create destinations never exist yet on first write. `ActionReplace` installs — the safety-critical case of overwriting a live user file — do go through `ReplaceFile` directly.

**Quality gates:**
- `go build ./...` — clean.
- `go vet ./...` — clean.
- `gofmt -l internal/migrate/operation.go internal/migrate/operation_test.go` — no output.
- `go test ./internal/migrate/... ./internal/data/...` — pass.
- `make build && make test` — pass (all packages).
- No `.savepoint/Health-Check.md` in this project, so the Quick health-check evidence step is skipped per the build-task skill.
