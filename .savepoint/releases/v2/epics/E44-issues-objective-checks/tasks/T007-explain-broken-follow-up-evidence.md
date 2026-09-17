---
id: E44-issues-objective-checks/T007-explain-broken-follow-up-evidence
title: Explain broken follow-up evidence
status: done
objective: Report every new Issue and Objective diagnostic through doctor and prove E44 behavior under full regression gates.
depends_on:
    - E44-issues-objective-checks/T004-save-issues-without-losing-history
    - E44-issues-objective-checks/T006-wait-for-objectives-that-are-not-ready
complexity_tier: medium
complexity_reason: Wires existing diagnostics into doctor and runs the epic's full regression evidence.
---

# T007: Explain broken follow-up evidence

## Problem

The new Issue and Objective rules fail closed inside `internal/data`, but a user running doctor sees none of them by name and gets no repair guidance. Inconsistencies that a load cannot refuse — a done Objective without integration clearance, a verified Issue whose proof was superseded — are not detected at all.

## Context Files

- `internal/doctor/checks.go`
- `internal/doctor/checks_test.go`
- `internal/doctor/repairs.go`
- `internal/doctor/repairs_test.go`
- `internal/doctor/report.go`
- `internal/doctor/report_test.go`
- `internal/data/objective_gate_v2.go`
- `internal/data/objective_gate_v2_test.go`
- `internal/data/issue_v2.go`
- `internal/data/gate_v2.go`
- `internal/data/errors.go`
- `internal/data/e2e_v2_test.go`

## Acceptance Criteria

- [x] Consistency inspection reports, without rewriting any record: an Objective `done` without current integration clearance, an Objective `done` while an owned Task is not `done`, and an Issue resolved as `verified` whose proof Check is no longer `CLEAR` or has been superseded.
- [x] Inspection walks records in sorted ID order and returns every problem found, not only the first.
- [x] Every new `internal/data` sentinel added in T001 through T006 maps to a stable doctor diagnostic name; no diagnostic falls through to the generic fallback.
- [x] Every new diagnostic name has a manual repair suggestion that names the field or record to fix.
- [x] Doctor reports Issue posture from derived counts computed at report time, with no stored summary, register file, or cached total read or written.
- [x] An open advisory Issue does not make a project unhealthy; structural and gate evidence decide health, and the distinction is asserted by test.
- [x] Doctor writes no project files on any path, including every new diagnostic and consistency path.
- [x] Doctor re-derives its own schema-free view: it consumes `data` diagnostics and inspection results and recreates no Issue, Objective, or gate rule.
- [x] An end-to-end V2 project fixture exercising Issues, an Objective Check, and a dependent Objective loads, gates, and reports as specified.
- [x] `make build && make test` passes across all packages, and `go vet ./...` is clean.

## Implementation Plan

- [x] Add `InspectObjectiveConsistency` and Issue proof consistency to `objective_gate_v2.go` and `issue_v2.go`, following `InspectTaskConsistency`'s read-only, sorted, return-everything shape.
- [x] Extend `v2DiagnosticName` and the consistency diagnostic mapping in `internal/doctor/checks.go` with the new sentinels and kinds.
- [x] Extend `V2ProblemRepair` and `V2ConsistencyRepair` in `internal/doctor/repairs.go` with a repair line per new name.
- [x] Add derived Issue counts to the doctor report, reusing the T002 listing helpers.
- [x] Add an end-to-end fixture in `e2e_v2_test.go` covering Issues, an Objective integration Check, and a dependent Objective.
- [x] Test every new diagnostic name and repair, the health distinction for advisory Issues, no-write behavior across all paths, and deterministic ordering.
- [x] Run the focused `internal/data` and `internal/doctor` suites, then `go vet ./...`, `git diff --check`, and `make build && make test`.

## Context Log

**Files read:** `internal/doctor/checks.go`, `internal/doctor/checks_test.go`, `internal/doctor/repairs.go`, `internal/doctor/repairs_test.go`, `internal/doctor/report.go`, `internal/doctor/report_test.go`, `internal/data/objective_gate_v2.go`, `internal/data/objective_gate_v2_test.go`, `internal/data/issue_v2.go`, `internal/data/issue_v2_test.go`, `internal/data/gate_v2.go`, `internal/data/check_v2.go`, `internal/data/project.go`, `internal/data/errors.go`, `internal/data/e2e_v2_test.go`, `internal/data/discover_test.go`, `internal/data/project_test.go`, `internal/data/write.go`, `internal/data/dependency_test.go`, `.savepoint/Guardrails.md`.

**Files edited:**
- `internal/data/objective_gate_v2.go`: added `ObjectiveConsistencyDiagnosticKind`/`ObjectiveConsistencyDiagnostic` and `InspectObjectiveConsistency` (done-without-clearance, done-with-incomplete-task), mirroring `InspectTaskConsistency`.
- `internal/data/issue_v2.go`: added `IssueConsistencyDiagnosticKind`/`IssueConsistencyDiagnostic` and `InspectIssueConsistency` (verified proof superseded by a later Check on the same scope).
- `internal/doctor/checks.go`: extended `v2DiagnosticName` with the 13 new E44 Issue sentinels (T001–T006); extended `v2ConsistencyProblems` to also report Objective and Issue consistency diagnostics via new `v2ObjectiveConsistencyDiagnosticName`/`v2IssueConsistencyDiagnosticName`/`v2ObjectiveSourcePath`/`v2IssueSourcePath` helpers; added `IssuePosture`/`IssuePostureReport` (derived Issue counts, no stored summary).
- `internal/doctor/repairs.go`: added repair suggestions for every new `V2ProblemRepair` and `V2ConsistencyRepair` name.
- `internal/doctor/report.go`: added `Issues *IssuePosture` field to `DiagnosticReport`, wired into `RunAllChecks`, and a new "Issue Posture" `Format()` section (advisory only — not counted in `HasProblems`).
- `internal/data/objective_gate_v2_test.go`, `internal/data/issue_v2_test.go`: unit tests for the new Inspect* functions (ignore-not-applicable, each diagnostic kind, clean case, sorted/return-everything order).
- `internal/data/e2e_v2_test.go`: added `TestE44_EpicScenario` — a dependent Objective blocked on readiness, an Objective closing under checker authority, a verified Issue against the Objective Check, then a later Check making both the Objective's clearance and the Issue's proof stale, detected by the new Inspect functions.
- `internal/doctor/checks_test.go`: per-sentinel `CheckProject` wiring tests, an exhaustive `v2DiagnosticName`/`V2ProblemRepair` fallback test over all 13 new sentinels, Objective/Issue consistency diagnostic tests, the open-advisory-Issue health test, `IssuePostureReport` tests, and two read-only (no-write) tests covering the new paths.
- `internal/doctor/repairs_test.go`, `internal/doctor/report_test.go`: repair-suggestion and Issue Posture section/health tests.

**Quality gates:**
- `go build ./...` — clean.
- `go test ./internal/data/... ./internal/doctor/...` — pass (also `go test ./...` via `make test`).
- `go vet ./...` — clean.
- `git diff --check` — clean.
- `make build && make test` — pass across all packages.
