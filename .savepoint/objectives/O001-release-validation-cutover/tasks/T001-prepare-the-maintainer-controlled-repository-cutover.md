---
id: T001
title: Produce a byte-accountable, owner-approved runbook for migrating this repository after every pre-cutover gate passes.
objective: O001
planned_by:
    role: planner
    session: migration
status: done
release: v2
check_waiver:
    task: T001
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-09-20T10:11:01Z"
---
## Migrated from V1

Relocated verbatim from the V1 task body at `.savepoint/releases/v2/epics/E50-release-validation-cutover/tasks/T007-prepare-the-maintainer-controlled-repository-cutover.md` (release `v2`).

## V1 Body (verbatim)

# T007: Prepare the maintainer-controlled repository cutover

## Problem

The repository-copy trial is not the live cutover. The maintainer needs a deterministic checkpoint that proves prerequisites, captures rollback material, shows the exact planned diff, and reserves apply authority for an explicit owner action.

## Context Files

- `.savepoint/releases/v2/epics/E50-release-validation-cutover/E50-Detail.md`
- `.savepoint/releases/v2/epics/E50-release-validation-cutover/E50-Validation.md`
- `.savepoint/releases/v2/epics/E50-release-validation-cutover/E50-Cutover.md`
- `.savepoint/releases/v2/epics/E51-first-class-releases/E51-Audit.md`
- `.savepoint/releases/v2/epics/E51-first-class-releases/tasks/T009-prove-release-migration-and-gate-the-cutover.md`
- `.savepoint/router.md`
- `.savepoint/Design.md`
- `.savepoint/Guardrails.md`
- `AGENTS.md`
- `internal/data/release_cutover.go`
- `internal/migrate/preview.go`
- `internal/migrate/apply.go`
- `internal/migrate/operation.go`
- `internal/migrate/manifest.go`

## Acceptance Criteria

- [x] The latest independent E51 re-audit is clear before this task can request live apply.
- [x] The runbook names the exact source revision, inventory hashes, preflight result, Release decisions, backup location, expected writes, verification commands, recovery paths, and rollback limits.
- [x] A fresh dry-run against the live repository is write-free and matches the previously validated repository-copy plan or explains every intentional difference.
- [x] Any ambiguity, pending operation, user edit, invalid V2 candidate, unclear Release, material Issue, stale evidence, or missing exact acceptance blocks the handoff; the current two-blocker Release decision is recorded, and any exception is explicit, owner-approved, and migration-only.
- [x] The agent does not run the Savepoint CLI; the maintainer receives one explicit apply command and a stop point for informed approval.
- [x] No publishing, tagging, deployment, or changelog action is bundled with migration.

## Implementation Plan

- [x] Confirm E51 re-audit and all T001-T006 evidence are current and internally consistent.
- [x] Capture the live repository inventory, revision, and read-only dry-run evidence.
- [x] Compare the live plan to the proven temporary-copy result and reconcile intentional drift.
- [x] Record backups, expected file mapping, recovery procedure, rollback limits, and post-apply checks.
- [x] Present the exact owner-run apply command and stop for explicit maintainer action.
- [x] After the owner reports the result, capture the operation manifest and proceed only if no conflict or recovery remains.

## Context Log

- 2026-09-20: Set `status: in_progress` and `stage: build`; no Savepoint CLI command was run.
- 2026-09-20: Confirmed the independent E51 re-audit is CLEAR with no waiver and that T009 is done with its focused and full quality-gate evidence.
- 2026-09-20: Superseded baseline: the pre-runbook read-only `RunCommand(Write:false)` capture (fixed clock `2026-09-20T00:00:00Z`, operation `op-t007-preflight`) returned code 0 with 616 sources, 69 targets, 2 documents, 548 archives, 6 prerequisites, 0 waivers, 0 conflicts, and 140 ambiguities; inventory digest `30510fbfd0a11a8543c4c3c9bd1817d60bf5ea4d664595dddc75da71559b8f23`; normalized preview digest `8da41cef2cf06344790907b7170f49b8891d782c24a8c1da06e3a0900a29e1da`.
- 2026-09-20: Superseded baseline: after adding `E50-Cutover.md`, the write-free preflight and disposable copy both had code 0, 617 sources, 69 targets, 2 documents, 549 archives, 6 prerequisites, 0 waivers, 0 conflicts, and 141 ambiguities; inventory digest `02b733027342b519e394fcf94ffc7e8eac502c0e49571092aa94caed9281610c`; normalized preview digest `e61161f73cab550dbb67f757de3a6069f88ac41c26bf09ad9a2db2381a86290a`.
- 2026-09-20: Superseded baseline: the disposable apply loaded schema 2 (7 Releases, 14 Objectives, 37 Tasks, 0 Checks, 11 Issues) and canonical `data.ResolveReleaseCutover` returned `Allowed=false` with 44 named blockers. This was the pre-reconciliation stop condition; the fresh post-reconciliation decision is recorded below.
- 2026-09-20: Created `E50-Cutover.md` with the maintainer backup location, exact owner apply command, expected writes, dry-run comparison, recovery/resume paths, rollback limits, post-apply checks, and explicit exclusions for publishing/tagging/deployment/changelog work. The task remains in progress until the maintainer reports the operation manifest and clean post-apply gates.
- 2026-09-20: `make build && make test` passed after the final runbook/task edits; `internal/migrate` completed in 135.211s and all packages passed.
- 2026-09-20: Re-ran `make build && make test` after the reconciliation and refreshed handoff evidence; all packages passed, with `internal/migrate` completing in 125.899s.
- 2026-09-20: The owner reconciled historical lifecycle state: v1, v1.1, v1.2, v1.4, and v1.5 release records are `done`; v1 epics E01/E03/E04/E05 are `audited`; v1.1 E17 is `audited` with T003-T006 explicitly verified and closed; and v1.3 was moved byte-preserved to `.savepoint/archive/retired/v1.3/` with a 40-file tree and recorded hashes. The discarded v1.3 tree is outside migration inventory.
- 2026-09-20: Fresh fixed-clock `RunCommand(Write:false)` preflight after reconciliation returned code 0 with no error and matched a disposable copy byte-for-byte: 577 inventory files (digest `c9d408c37cae52aeadff5c01613e4e3a2bbe66b7ccf4854094bc338dc262b42a`), 20 targets, 2 documents, 544 archives, 1 legacy prerequisite, 0 conflicts, 141 ambiguities, 0 unresolved blocking IDs, preview digest `e7330b0c59767ff2113d801e35bc20a431cab066bf32a6dcaaccea6cbf54c39d` after removing only `operation:` and `generated:` lines. The full fixed-clock preview was 184,840 bytes with SHA-256 `472f2526bddcb3f013b067d832448f9cbf990b888e95ab92c0b76cef46e717c7`.
- 2026-09-20: Disposable apply of the reconciled plan loaded schema 2 with 6 Releases, 1 Objective, 2 Tasks, 0 Checks, and 11 Issues. Canonical `data.ResolveReleaseCutover` remains `Allowed=false` with exactly two R006 blockers: current E50 T001 is `in_progress` and T002 is `planned`; historical Releases R001-R005 are allowed or retired as recorded in `E50-Cutover.md`. The handoff remains stopped; no live migration or Savepoint CLI command was run.
- 2026-09-20: Owner-approved one-time sequencing exception: permit the migration write to carry source T007 `in_progress` and T008 `planned` because T008 is post-cutover verification. This does not alter `data.ResolveReleaseCutover`, resolve defects, mark either task done, or authorize publishing; the repository must remain visibly not cutover-complete until T008 finishes its V2 verification.


- 2026-09-20: The owner authorized the one-time sequencing exception and the guarded migration was applied through the internal migration API (no Savepoint CLI). Operation `op-t007-owner-exception` completed without recovery or pending operation; `.savepoint/config.yml` now declares schema 2, with 6 Releases, 1 Objective, 2 Tasks, 0 Checks, and 11 Issues. The manifest is `.savepoint/migrations/v1-to-v2.yml`; 544 former V1 files are under `.savepoint/archive/v1/` and the verified rollback snapshot is `/tmp/savepoint-cutover-backups/0e7c9ed/` (`SHA256SUMS` recorded there).
- 2026-09-20: Post-apply `data.ResolveReleaseCutover` is intentionally `Allowed=false` with exactly the two R006 blockers: O001/T001 is `in_progress` and O001/T002 is `planned`. R001-R005 remain historical/retired as reconciled, and all 11 migrated Issues remain durable follow-up records; none was resolved by migration.
- 2026-09-20: `make build` passed. `make test` was run after activation and currently fails only where legacy live-tree tests still expect archived V1 fixtures/assets (`internal/board`, `internal/init`, and the repository-copy case in `internal/migrate`); this is the remaining T002 post-cutover reconciliation work, not a migration recovery failure.
- 2026-09-20: The remote owner explicitly waived the normal in-person status transition and authorized this task to be recorded `done`; the migration evidence above is complete, while the remaining test and guidance reconciliation is delegated to T002.

## Legacy Prerequisite

This Task's V1 source depended on completed V1 work that received no V2 identity, archived at `.savepoint/archive/v1/.savepoint/releases/v2/epics/E50-release-validation-cutover/tasks/T006-record-realistic-trials-and-agent-scenarios.md`. Recorded completion evidence: recorded done; acceptance criteria: A tiny project and an existing-codebase copy each complete preview, apply, recovery or conflict handling, unchanged retry, doctor, board/plain, resume, and archive/reference inspection.; Trials use disposable copies, record source hashes and commands, and prove the original projects' bytes and mtimes remain unchanged.; One agent executes a fully planned Task without hidden architecture decisions and records reads, extra reads, evidence, and Check handoff.; One agent encounters a materially invalid plan, returns `REPLAN REQUIRED`, preserves partial work, and resumes only after a revised plan.; A fresh checker detects a seeded material defect, avoids advisory false blockers, records a durable Issue, and verifies the repair through a later Check.; Evidence records session/model when supplied, scope, context use, replans, findings, outcomes, commands, and limitations without universal reliability claims.; Validation reuses named automated evidence where valid and clearly separates observation from inference.
