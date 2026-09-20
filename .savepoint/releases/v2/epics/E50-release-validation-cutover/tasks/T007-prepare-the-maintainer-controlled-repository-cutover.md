---
id: E50-release-validation-cutover/T007-prepare-the-maintainer-controlled-repository-cutover
status: in_progress
stage: build
objective: Produce a byte-accountable, owner-approved runbook for migrating this repository after every pre-cutover gate passes.
depends_on:
    - E50-release-validation-cutover/T006-record-realistic-trials-and-agent-scenarios
complexity_tier: high
complexity_reason: Coordinates audit, migration, backups, Release gates, and an irreversible workflow boundary.
---

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
- [x] Any ambiguity, pending operation, user edit, invalid V2 candidate, unclear Release, material Issue, stale evidence, or missing exact acceptance blocks the handoff; the current 44-blocker Release decision is recorded as a hard stop.
- [x] The agent does not run the Savepoint CLI; the maintainer receives one explicit apply command and a stop point for informed approval.
- [x] No publishing, tagging, deployment, or changelog action is bundled with migration.

## Implementation Plan

- [x] Confirm E51 re-audit and all T001-T006 evidence are current and internally consistent.
- [x] Capture the live repository inventory, revision, and read-only dry-run evidence.
- [x] Compare the live plan to the proven temporary-copy result and reconcile intentional drift.
- [x] Record backups, expected file mapping, recovery procedure, rollback limits, and post-apply checks.
- [x] Present the exact owner-run apply command and stop for explicit maintainer action.
- [ ] After the owner reports the result, capture the operation manifest and proceed only if no conflict or recovery remains.

## Context Log

- 2026-09-20: Set `status: in_progress` and `stage: build`; no Savepoint CLI command was run.
- 2026-09-20: Confirmed the independent E51 re-audit is CLEAR with no waiver and that T009 is done with its focused and full quality-gate evidence.
- 2026-09-20: Read-only live `RunCommand(Write:false)` preflight (fixed clock `2026-09-20T00:00:00Z`, operation `op-t007-preflight`) returned code 0 with no error. The pre-runbook plan had 616 sources, 69 targets, 2 documents, 548 archives, 6 prerequisites, 0 waivers, 0 conflicts, 140 ambiguities, and no unresolved blocking IDs; inventory digest `30510fbfd0a11a8543c4c3c9bd1817d60bf5ea4d664595dddc75da71559b8f23`; normalized preview digest `8da41cef2cf06344790907b7170f49b8891d782c24a8c1da06e3a0900a29e1da`.
- 2026-09-20: After adding `E50-Cutover.md`, repeated the write-free preflight and disposable repository-copy comparison. Both had code 0, 617 sources, 69 targets, 2 documents, 549 archives, 6 prerequisites, 0 waivers, 0 conflicts, 141 ambiguities, and no unresolved blocking IDs; inventory digest `02b733027342b519e394fcf94ffc7e8eac502c0e49571092aa94caed9281610c`; normalized preview digest `e61161f73cab550dbb67f757de3a6069f88ac41c26bf09ad9a2db2381a86290a`.
- 2026-09-20: The disposable apply loaded schema 2 (7 Releases, 14 Objectives, 37 Tasks, 0 Checks, 11 Issues) but canonical `data.ResolveReleaseCutover` returned `Allowed=false` with 44 named blockers: R001 4, R002 4, R003 1, R004 31, R005 1, R006 1, and R007 2. This is a clear Release gate failure, not a waiver; the runbook records it as the stop condition.
- 2026-09-20: Created `E50-Cutover.md` with the maintainer backup location, exact owner apply command, expected writes, dry-run comparison, recovery/resume paths, rollback limits, post-apply checks, and explicit exclusions for publishing/tagging/deployment/changelog work. The task remains in progress until the maintainer reports the operation manifest and clean post-apply gates.
- 2026-09-20: `make build && make test` passed after the final runbook/task edits; `internal/migrate` completed in 135.211s and all packages passed.
