---
id: E43-task-check-gates/T006-explain-broken-evidence
title: Explain broken check evidence
status: done
objective: Report every new Check and evidence diagnostic through doctor and prove E43 behavior under full regression gates.
depends_on:
    - E43-task-check-gates/T003-save-checks-safely
    - E43-task-check-gates/T005-decide-when-tasks-close
complexity_tier: medium
complexity_reason: Integrates existing data diagnostics with doctor and validates combined epic behavior.
---

# T006: Explain broken check evidence

## Problem

The new Check and evidence contracts are only useful if a user can see what is wrong and what to do about it. Doctor currently reports one structural load failure and has no path for evaluation-level inconsistencies that do not fail the load.

## Context Files

- `internal/doctor/checks.go`
- `internal/doctor/checks_test.go`
- `internal/doctor/repairs.go`
- `internal/doctor/repairs_test.go`
- `internal/doctor/report.go`
- `internal/doctor/report_test.go`
- `internal/data/gate_v2.go`
- `internal/data/errors.go`

## Acceptance Criteria

- [x] Every new Check and evidence sentinel has a stable diagnostic name; the name for a given problem does not change between runs.
- [x] Each new diagnostic name maps to a manual repair suggestion; doctor never repairs or rewrites a Check, an evidence field, or a record status.
- [x] Doctor reports evaluation-level inconsistencies that do not fail the load, alongside structural load failures, with record path and identity context.
- [x] Multiple inconsistencies in one project are all reported, in deterministic order.
- [x] Doctor reports no V2 Check or evidence problems for a V1 project and none for a clean V2 project.
- [x] Doctor writes no project files: a full run leaves every file's bytes and modified time unchanged.
- [x] Doctor consumes the data-layer decisions and diagnostics without restating clearance, authority, or dependency rules.
- [x] Epic-level behavior is proven end to end on a temporary project: records load, clearance resolves, a technical Task closes, an owner-validated Task waits, a rerun supersedes, and evidence writes preserve authored content.
- [x] The frozen V1 fixtures from E41 load and behave exactly as they do today.
- [x] `make build && make test` passes, with named cases and outcomes recorded in this task.

## Implementation Plan

- [x] Extend the diagnostic-name mapping in `checks.go` to cover every new sentinel.
- [x] Add a doctor check that reports the data-layer evidence inconsistencies for a loaded V2 project.
- [x] Add repair suggestions for each new name in `repairs.go` and confirm report rendering handles multiple problems.
- [x] Test each diagnostic name, each repair suggestion, multiple simultaneous problems, the V1 and clean-V2 no-problem cases, and read-only behavior by snapshotting bytes and modified times.
- [x] Add the end-to-end epic scenario test on a temporary project.
- [x] Re-run the E41 V1 fixture characterization tests, then run `go vet ./...` and `make build && make test`.

## Context Log

**Files read:** `internal/doctor/checks.go`, `internal/doctor/repairs.go`, `internal/doctor/report.go`, `internal/doctor/checks_test.go`, `internal/doctor/repairs_test.go`, `internal/doctor/report_test.go`, `internal/data/errors.go`, `internal/data/check_v2.go`, `internal/data/evidence_v2.go`, `internal/data/gate_v2.go`, `internal/data/project.go`, `internal/data/task_v2.go`, `internal/data/discover.go`, `internal/data/write.go`, `internal/data/gate_v2_test.go`, `internal/data/project_test.go`, `internal/data/discover_test.go`, `internal/data/write_test.go`, `.savepoint/Guardrails.md` (no `.savepoint/Health-Check.md` present — Quick check step skipped, absence is not a finding).

**Files edited:**
- `internal/doctor/checks.go` — extended `v2DiagnosticName` with the 7 new Check/evidence sentinels (`v2-check-malformed`, `v2-check-missing-scope-target`, `v2-check-missing-reference`, `v2-check-supersedes-conflict`, `v2-evidence-malformed`, `v2-evidence-missing-reference`, `v2-check-immutable`); `CheckProject` now also runs `data.InspectTaskConsistency` on a successfully loaded V2 index and reports each result as a Problem via new `v2ConsistencyProblems`/`v2ConsistencyDiagnosticName`/`v2TaskSourcePath` helpers (`v2-done-without-clearance`, `v2-acceptance-superseded`, `v2-evidence-contradicts-status`).
- `internal/doctor/repairs.go` — added the 7 new cases to `V2ProblemRepair` and a new `V2ConsistencyRepair(name string) string` for the 3 consistency names, all manual-only (no auto-heal).
- `internal/doctor/checks_test.go` — added `TestCheckProject_CheckMalformed`, `TestCheckProject_CheckMissingScopeTarget`, `TestCheckProject_CheckMissingReference`, `TestCheckProject_CheckSupersedesConflict`, `TestCheckProject_EvidenceMalformed`, `TestCheckProject_EvidenceMissingReference`, `TestCheckProject_v2ValidWithChecksNoProblems`, `TestCheckProject_ConsistencyDiagnostics` (3 simultaneous problems, deterministic task-ID order), `TestCheckProject_ConsistencyReadOnly` (bytes + mtime unchanged).
- `internal/doctor/repairs_test.go` — added `TestV2ProblemRepair_checkAndEvidenceNames`, `TestV2ConsistencyRepair`.
- `internal/data/e2e_v2_test.go` (new) — `TestE43_EpicScenario`: builds a real temp project, proves missing→current clearance, a technical Task closing under checker authority, an owner-validated Task blocked on `GateBlockOwnerAcceptance` then allowed after acceptance, a superseding rerun re-blocking completion (`ClearanceStale`), and evidence writes preserving an unknown frontmatter field and the authored Markdown body throughout.

**Quality gates:**
- `go vet ./...` — clean.
- `make build && make test` — all packages pass (`internal/data` 0.287s, `internal/doctor` cached/clean, plus `internal/board`, `internal/init`, `cmd`, root, `internal/buildtool`, `internal/styles`).
- E41 V1 fixture characterization tests re-run explicitly: `TestMigrationHistory*`, `TestMigrationSourceBasic*`, `TestRouterReader*` — all PASS, unchanged.
- Health-Check.md Quick check: skipped (file absent from this project).
