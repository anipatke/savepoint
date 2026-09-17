---
id: E45-safe-migration/T008-make-an-interrupted-migration-recoverable
title: Make an interrupted migration recoverable
status: planned
objective: Record the operation with verified backups, staging, and a journal so any interruption leaves recoverable bytes and a truthful report.
depends_on:
    - E45-safe-migration/T001-decide-how-windows-replaces-a-file
    - E45-safe-migration/T003-plan-the-conversion-before-touching-anything
complexity_tier: high
complexity_reason: Owns the recovery guarantee: backup verification, journal states, and interruption behavior at every step.
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

- [ ] An operation directory at `.savepoint/.migration/<op-id>/` holds `backup/`, `staging/`, and `operation.yml`, and no path the operation writes ever falls inside it.
- [ ] The operation directory is create-only: an existing `<op-id>` is refused with a named diagnostic rather than reused or overwritten, so a prior backup or user sidecar can never be destroyed.
- [ ] `operation.yml` records the operation ID, its creation time, the plan's source hashes, and for every path an action, a planned content hash, an installed hash once written, and a state among `planned`, `backed_up`, `staged`, `installed`, and `verified`.
- [ ] Every file the operation will replace or remove is copied to `backup/` and each copy is re-hashed against the recorded source hash before any replacement begins; a backup that fails verification aborts before the first replacement.
- [ ] Staged target files are written complete in `staging/` before any live path is touched.
- [ ] The journal is written before and after each step, so a crash at any point leaves a state that names what had happened and what had not.
- [ ] `PendingOperation(root)` is a read-only detector returning the incomplete operation and its recovery guidance, and returning nothing when none exists.
- [ ] More than one incomplete operation produces a named diagnostic rather than a heuristic choice between them.
- [ ] A test proves `LoadV2Index` ignores `.savepoint/.migration/` entirely, including an operation directory containing staged files that look like Objectives, Checks, or Issues.
- [ ] Injected failures at each journal step leave every affected path recoverable from either its original location or its verified backup; a test asserts there is no step at which both the original and its backup are absent.
- [ ] The failure report names the operation, the step, the affected paths, and the recovery command, and never reports a partial state as a success.
- [ ] Recovery distinguishes FS-06's preflight no-partial-write expectation from a recorded interruption after a valid apply began, and says which one occurred.

## Implementation Plan

- [ ] Add `operation.go` with the operation directory layout, the create-only guard, and the journal model and its states.
- [ ] Implement backup capture and per-copy hash verification against the plan's recorded source hashes.
- [ ] Implement staging writes through T001's replacement primitive, never through `AtomicWrite`'s truncating fallback.
- [ ] Implement journal persistence around each step, flushing before the step it describes.
- [ ] Implement `PendingOperation` and the recovery-guidance report, including the multiple-operation diagnostic.
- [ ] Add the discovery-isolation test asserting a staged Objective-shaped file inside the operation directory is invisible to `LoadV2Index`.
- [ ] Add injected-failure tests at every journal step, each asserting recoverable bytes, a truthful report, and the original-or-backup invariant.
- [ ] Run the focused `internal/migrate` and `internal/data` suites, then `make build && make test`.

## Context Log

Pending.
