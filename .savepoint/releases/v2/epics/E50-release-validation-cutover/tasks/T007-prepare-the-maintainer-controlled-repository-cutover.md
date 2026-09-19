---
id: E50-release-validation-cutover/T007-prepare-the-maintainer-controlled-repository-cutover
status: planned
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

- [ ] The latest independent E51 re-audit is clear before this task can request live apply.
- [ ] The runbook names the exact source revision, inventory hashes, preflight result, Release decisions, backup location, expected writes, verification commands, recovery paths, and rollback limits.
- [ ] A fresh dry-run against the live repository is write-free and matches the previously validated repository-copy plan or explains every intentional difference.
- [ ] Any ambiguity, pending operation, user edit, invalid V2 candidate, unclear Release, material Issue, stale evidence, or missing exact acceptance blocks the handoff.
- [ ] The agent does not run the Savepoint CLI; the maintainer receives one explicit apply command and a stop point for informed approval.
- [ ] No publishing, tagging, deployment, or changelog action is bundled with migration.

## Implementation Plan

- [ ] Confirm E51 re-audit and all T001-T006 evidence are current and internally consistent.
- [ ] Capture the live repository inventory, revision, and read-only dry-run evidence.
- [ ] Compare the live plan to the proven temporary-copy result and reconcile intentional drift.
- [ ] Record backups, expected file mapping, recovery procedure, rollback limits, and post-apply checks.
- [ ] Present the exact owner-run apply command and stop for explicit maintainer action.
- [ ] After the owner reports the result, capture the operation manifest and proceed only if no conflict or recovery remains.

## Context Log

Pending.
