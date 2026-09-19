---
id: E50-release-validation-cutover/T008-verify-and-reconcile-the-live-v2-cutover
status: planned
objective: Verify the owner-run migration, reconcile active guidance, and hand the completed V2 release to independent audit.
depends_on:
    - E50-release-validation-cutover/T007-prepare-the-maintainer-controlled-repository-cutover
complexity_tier: high
complexity_reason: Final verification spans live state, documentation, tests, distribution, and audit handoff.
---

# T008: Verify and reconcile the live V2 cutover

## Problem

An owner-run apply does not complete E50 by itself. The resulting repository must load cleanly, preserve its source history, activate one V2 workflow, pass every gate, and remove transitional claims before independent audit.

## Context Files

- `.savepoint/releases/v2/epics/E50-release-validation-cutover/E50-Detail.md`
- `.savepoint/releases/v2/epics/E50-release-validation-cutover/E50-Validation.md`
- `.savepoint/releases/v2/epics/E50-release-validation-cutover/E50-Cutover.md`
- `.savepoint/router.md`
- `.savepoint/Design.md`
- `.savepoint/Guardrails.md`
- `AGENTS.md`
- `README.md`
- `main.go`
- `internal/data/project.go`
- `internal/data/release_cutover.go`
- `internal/doctor/checks.go`
- `internal/board/v2/load.go`
- `internal/board/v2/run.go`
- `internal/resume/resume.go`
- `internal/init/template_freshness_test.go`
- `internal/migrate/end_to_end_test.go`

## Acceptance Criteria

- [ ] The live repository declares schema V2, has no pending migration operation, loads cleanly, and produces matching board, plain, resume, and doctor interpretations.
- [ ] The migration manifest accounts for every former V1 source and reference at a live or byte-preserved archive destination; the verified backup remains available.
- [ ] Every declared Release passes the canonical cutover decision, with exact current evidence and required owner acceptance visible.
- [ ] Router and AGENTS activate only `idea`, `design`, `task`, and `check`; obsolete V1 skill activation and transitional architecture claims are removed.
- [ ] README, Design, Guardrails, Codebase Map, templates, help, and validation records describe the shipped V2-only runtime consistently.
- [ ] Focused regressions, six-platform distribution checks, `make build`, `make test`, and `git diff --check` pass with named evidence.
- [ ] The router advances to `audit-pending` only after all implementation items and evidence are complete; no task or epic is marked done by the agent.

## Implementation Plan

- [ ] Verify the owner-run operation state, manifest, backup, schema, and canonical Release cutover result.
- [ ] Compare board, plain output, resume, and doctor on the live repository without mutating project state.
- [ ] Reconcile active Design, Guardrails, AGENTS, README, help, templates, and Codebase Map to final reality.
- [ ] Remove remaining transitional live assets or wording while preserving migration code and frozen history.
- [ ] Run focused, full, distribution, and diff-quality gates and record exact outcomes.
- [ ] Complete the validation/cutover records and route E50 to a fresh independent epic audit session.

## Context Log

Pending.
