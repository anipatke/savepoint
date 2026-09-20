---
id: T002
title: Verify the owner-run migration, reconcile active guidance, and hand the completed V2 release to independent audit.
objective: O001
planned_by:
    role: planner
    session: migration
status: planned
depends_on:
    - task: T001
release: v2
check_waiver:
    task: T002
    reason: Owner waived T002's optional local Task Check; evidence routes to O001's mandatory Full Objective Check instead.
    actor:
        role: owner
        session: owner-remote-control
    recorded_at: "2026-09-20T07:48:58Z"
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

- [x] The live repository declares schema V2, has no pending migration operation, loads cleanly, and produces matching board, plain, resume, and doctor interpretations.
- [x] The migration manifest accounts for every former V1 source and reference at a live or byte-preserved archive destination; the verified backup remains available.
- [ ] Every declared Release passes the canonical cutover decision, with exact current evidence and required owner acceptance visible. **Not yet true by design**: `ResolveReleaseCutover` blocks R006 on `release_objective_incomplete` (member Objective O001: "task T002 is not done") — the one remaining, expected blocker until this Task closes and O001's Full Objective Check runs. This criterion closes after that Check and owner action, not by this Task alone.
- [x] Router and AGENTS activate only `idea`, `design`, `task`, and `check`; obsolete V1 skill activation and transitional architecture claims are removed.
- [x] README, Design, Guardrails, Codebase Map, templates, help, and validation records describe the shipped V2-only runtime consistently.
- [x] Focused regressions, six-platform distribution checks, `make build`, `make test`, and `git diff --check` pass with named evidence.
- [x] The router advances to `check` (V2's `audit-pending` equivalent) only after all implementation items and evidence are complete; no task or epic is marked done by the agent.

## Implementation Plan

- [x] Verify the owner-run operation state, manifest, backup, schema, and canonical Release cutover result.
- [x] Compare board, plain output, resume, and doctor on the live repository without mutating project state.
- [x] Reconcile active Design, Guardrails, AGENTS, README, help, templates, and Codebase Map to final reality.
- [x] Remove remaining transitional live assets or wording while preserving migration code and frozen history.
- [x] Run focused, full, distribution, and diff-quality gates and record exact outcomes.
- [x] Complete the validation/cutover records and route E50 to a fresh independent epic audit session.

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
- 2026-09-20: Re-verified the live repository with a throwaway Go probe (`data.LoadV2Index`, `migrate.PreflightCutover`, `data.ResolveReleaseCutover`, `data.NewRouterReader().ReadStateV2`, `data.ResolveNext` against this repo root; deleted after use, never committed) rather than running the `savepoint` CLI. Findings: no pending migration operation; schema already V2; index loads 6 Releases / 1 Objective / 2 Tasks / 0 Checks / 11 Issues; R001-R005 are `done`, R006 is `in_progress`; `ResolveReleaseCutover` reports exactly one blocker, `release_objective_incomplete` on R006 because T002 itself is not yet done (expected — this Task cannot close itself). **Defect found and fixed**: `data.NewRouterReader().ReadStateV2` failed on the live `.savepoint/router.md` with `yaml: line 5: mapping values are not allowed in this context` — the `next_action` value written into this file in the prior gate-resolver-reconciliation edit contained an unquoted `:`, which YAML parses as a nested mapping key. The same unquoted-colon pattern was already present in the *previous* `next_action` value (visible via `git log -p`), so this parse failure predates this session and was never actually exercised end-to-end before this probe — meaning the "board, plain, resume, and doctor interpretations agree" claims in earlier Context Log entries were evaluated against router.md content that had not yet re-triggered this specific bug, and this entry is the first to actually parse the file with today's content. Fixed by quoting `next_action` as a YAML string; re-ran the probe and `ReadStateV2`/`ResolveNext` succeeded, resolving `Next.Kind = check_needed` for T002 — consistent with `ResolveTaskCompletion`'s `ClearanceMissing` result (no Check, no waiver recorded on T002) and with routing this Task to Check next. `go test ./...` still passes after the router.md fix (data package tests don't touch this file; the parse only runs against the live repo through `resume`/`board`/`doctor`, which the probe exercised directly).
- 2026-09-20: Acceptance criteria updated against the probe's concrete findings: schema/no-pending-op/clean-load/resolver-parity, manifest (`​.savepoint/migrations/v1-to-v2.yml`, byte-preserved sources under `.savepoint/archive/v1/`), router/AGENTS four-state-only, and doc/gate/build evidence are all satisfied and checked. The Release-cutover criterion is left unchecked on purpose — it is blocked on the one expected, correctly-reported condition (T002 itself is not done, so R006 cannot be clear yet), not a defect; it closes after O001's Full Objective Check and owner action, outside this Task's own authority. `.savepoint/router.md` is now `state: check`, `objective: O001`, `task: T002` — this is V2's equivalent of routing E50 to independent audit. Whether the fresh Check session runs an optional local Task Check on T002 first, or the owner instead records an explicit Task-check waiver (as with T001) so the evidence routes straight to O001's mandatory Full Objective Check, is the owner's call per the Write Boundary — not invented here.
- 2026-09-20: The owner chose to waive T002's optional Task Check. Recorded a `check_waiver` evidence block in this Task's frontmatter (`task: T002`, `actor: {role: owner, session: owner-remote-control}`, `recorded_at: 2026-09-20T07:48:58Z`) — the first real-world use of the `CheckWaiver` mechanism built earlier today. Verified by decode: `index.Tasks["T002"].Evidence.CheckWaiver` is populated and `data.ResolveTaskCompletion(index, "T002")` now returns `Allowed=true, AllowedByWaiver=true, Actor=owner`, i.e. T002's evidence is ready for the owner to close; it is not technical `CLEAR` and does not by itself close the Task or satisfy O001's own mandatory Full Objective Check. Per the Write Boundary, this executor session does not set `status: done` on T002 itself — that remains the owner's own action.

## Replan Required

Resolved: the post-migration Context Files section above names the archived, byte-preserved E50 sources and the active V2 implementation files. The V1 body remains unchanged for provenance; implementation may resume within those reconciled boundaries.

Resolved: the runtime gate resolvers and their tests now enforce the optional-Task / mandatory-Objective-and-Release split via the `CheckWaiver` evidence block and `GateDecision.AllowedByWaiver`, matching the policy already recorded in `.savepoint/Design.md`, `.savepoint/Guardrails.md`, and the skill files. Implementation may resume within T002's remaining plan item.
