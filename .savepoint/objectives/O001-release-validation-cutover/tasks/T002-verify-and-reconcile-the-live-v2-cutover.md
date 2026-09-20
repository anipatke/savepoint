---
id: T002
title: Verify the owner-run migration, reconcile active guidance, and hand the completed V2 release to independent audit.
objective: O001
planned_by:
    role: planner
    session: migration
status: in_progress
stage: audit
depends_on:
    - task: T001
release: v2
---
## Migrated from V1

Relocated verbatim from the V1 task body at `.savepoint/releases/v2/epics/E50-release-validation-cutover/tasks/T008-verify-and-reconcile-the-live-v2-cutover.md` (release `v2`).

## Context Files (post-migration)

- `.savepoint/archive/v1/.savepoint/releases/v2/epics/E50-release-validation-cutover/E50-Detail.md`
- `.savepoint/archive/v1/.savepoint/releases/v2/epics/E50-release-validation-cutover/E50-Validation.md`
- `.savepoint/archive/v1/.savepoint/releases/v2/epics/E50-release-validation-cutover/E50-Cutover.md`
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

## V1 Body (verbatim)

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

- [x] Verify the owner-run operation state, manifest, backup, schema, and canonical Release cutover result.
- [x] Compare board, plain output, resume, and doctor on the live repository without mutating project state.
- [x] Reconcile active Design, Guardrails, AGENTS, README, help, templates, and Codebase Map to final reality.
- [x] Remove remaining transitional live assets or wording while preserving migration code and frozen history.
- [x] Run focused, full, distribution, and diff-quality gates and record exact outcomes.
- [ ] Complete the validation/cutover records and route E50 to a fresh independent epic audit session.

## Context Log

- 2026-09-20: T001 is recorded `done` under the owner's explicit remote waiver; its migration operation, backup, schema, and canonical cutover evidence are the dependency for this Task.
- 2026-09-20: Started T002 at `stage: build` to reconcile the post-migration V2 guidance and legacy test assumptions before the independent check.
- 2026-09-20: REPLAN REQUIRED. The migration correctly moved the three declared E50 source context files out of their V1 paths; all three `.savepoint/releases/v2/epics/E50-release-validation-cutover/*` paths are missing, with byte-preserved equivalents under `.savepoint/archive/v1/.savepoint/releases/v2/epics/E50-release-validation-cutover/`. The task plan must be updated to name post-migration context paths before implementation continues; no substitute path was silently adopted.
- 2026-09-20: The live V2 parity probe loaded schema 2 with 6 Releases, 1 Objective, 2 Tasks, 0 Checks, and 11 Issues. `data.ResolveNext` selected `execute` for T002; non-TTY board, resume, and doctor agreed on the same task. Doctor reported the expected two current blockers: T001 is done without technical clearance, and T002 is still in progress, so R006 is not clear. Historical R001-R005 remain allowed or retired according to their migrated records.
- 2026-09-20: `go test ./internal/board ./internal/init ./internal/migrate` passed after redirecting preserved V1 fixture reads to `.savepoint/archive/v1/` and reconstructing V1 repository boundaries inside migration tests. `go test ./...` and `make build && make test` both passed; `internal/migrate` completed in about 123 seconds on the full run. `git diff --check` passed.
- 2026-09-20: Distribution gates passed: `make build-all`, `make dist`, `make verify-dist`, and the package test. The repository `make package-check` target reached the npm test but its default npm cache/log directory was not writable in this environment; the equivalent `npm --cache /tmp/savepoint-npm-cache pack --dry-run --ignore-scripts` passed and listed all six platform binaries.
- 2026-09-20: Re-ran the exact `make package-check` target with `NPM_CONFIG_CACHE=/tmp/savepoint-npm-cache`; npm test and six-platform dry-run packaging both passed.
- 2026-09-20: Active V2 guidance now names the four-state router, V2 Design boundary, archived V1 source, and legacy-template boundary. `templates/project/AGENTS.md` is explicitly labeled as a V1 migration fixture; `templates/project-v2/` remains the active scaffold. V1 skills remain available for that compatibility path but are not activated by the live router.
- 2026-09-20: REPLAN REQUIRED. The owner clarified a global V2 gate policy: an individual Task Check is optional and may be explicitly waived for any implementation, while the Full Objective Check (the V2 equivalent of an epic check) is mandatory and Release Checks remain mandatory when a Release exists. Do not create a T001 Task Check under the superseded policy; reconcile the shared gate contract before resuming T002.
- 2026-09-20: Reconciled `internal/data`'s runtime gate resolvers and tests with the documented policy in `a62f7c7`. Added a `check_waiver` Evidence sub-block (`CheckWaiver`, decoded by `decodeCheckWaiverV2` in `evidence_v2.go`) requiring `task`, `reason`, `actor.role: owner`, and `recorded_at`; it decodes only on Task evidence (`ErrV2EvidenceMalformed` on Objective/Release) and only when it names its own Task. `ResolveTaskCompletion` in `gate_v2.go` now grants completion `AllowedByWaiver` under owner authority only when clearance is `ClearanceMissing` (no Check was ever requested) and an applicable waiver is recorded — never as a `CLEAR` result, and never once any Check exists (`needs_work`, `stale`, and `unknown` clearance stay unwaivable). `ResolveObjectiveCompletion`, `ResolveTaskDependencyV2` (so a `requires: clear` dependency still ignores waivers), and Release resolution are unchanged, keeping Objective and Release Checks mandatory. `InspectTaskConsistency` now excludes `AllowedByWaiver` from `evidence_contradicts_status`, matching how it already excludes `AllowedByException`. Added decode tests in `evidence_v2_test.go` and gate/consistency tests in `gate_v2_test.go` (including that a waiver does not survive once a Check exists, does not apply across Tasks, and does not loosen Objective completion). `go vet ./...`, `go test ./...` (all packages, `internal/migrate` ~124s), `make build`, and `git diff --check` all passed. This was a data/gate-resolver change only; `internal/board/v2` and `internal/resume` UI surfacing of `AllowedByWaiver` (an owner-visible board action analogous to the existing exception-close key) was not touched and remains open follow-up if the owner wants board-driven waiver closure rather than a direct file edit.

## Replan Required

Resolved: the post-migration Context Files section above names the archived, byte-preserved E50 sources and the active V2 implementation files. The V1 body remains unchanged for provenance; implementation may resume within those reconciled boundaries.

Resolved: the runtime gate resolvers and their tests now enforce the optional-Task / mandatory-Objective-and-Release split via the `CheckWaiver` evidence block and `GateDecision.AllowedByWaiver`, matching the policy already recorded in `.savepoint/Design.md`, `.savepoint/Guardrails.md`, and the skill files. Implementation may resume within T002's remaining plan item.
