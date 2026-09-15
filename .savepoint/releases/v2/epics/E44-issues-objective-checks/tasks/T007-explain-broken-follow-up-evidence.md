---
id: E44-issues-objective-checks/T007-explain-broken-follow-up-evidence
title: Explain broken follow-up evidence
status: planned
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

- [ ] Consistency inspection reports, without rewriting any record: an Objective `done` without current integration clearance, an Objective `done` while an owned Task is not `done`, and an Issue resolved as `verified` whose proof Check is no longer `CLEAR` or has been superseded.
- [ ] Inspection walks records in sorted ID order and returns every problem found, not only the first.
- [ ] Every new `internal/data` sentinel added in T001 through T006 maps to a stable doctor diagnostic name; no diagnostic falls through to the generic fallback.
- [ ] Every new diagnostic name has a manual repair suggestion that names the field or record to fix.
- [ ] Doctor reports Issue posture from derived counts computed at report time, with no stored summary, register file, or cached total read or written.
- [ ] An open advisory Issue does not make a project unhealthy; structural and gate evidence decide health, and the distinction is asserted by test.
- [ ] Doctor writes no project files on any path, including every new diagnostic and consistency path.
- [ ] Doctor re-derives its own schema-free view: it consumes `data` diagnostics and inspection results and recreates no Issue, Objective, or gate rule.
- [ ] An end-to-end V2 project fixture exercising Issues, an Objective Check, and a dependent Objective loads, gates, and reports as specified.
- [ ] `make build && make test` passes across all packages, and `go vet ./...` is clean.

## Implementation Plan

- [ ] Add `InspectObjectiveConsistency` and Issue proof consistency to `objective_gate_v2.go` and `issue_v2.go`, following `InspectTaskConsistency`'s read-only, sorted, return-everything shape.
- [ ] Extend `v2DiagnosticName` and the consistency diagnostic mapping in `internal/doctor/checks.go` with the new sentinels and kinds.
- [ ] Extend `V2ProblemRepair` and `V2ConsistencyRepair` in `internal/doctor/repairs.go` with a repair line per new name.
- [ ] Add derived Issue counts to the doctor report, reusing the T002 listing helpers.
- [ ] Add an end-to-end fixture in `e2e_v2_test.go` covering Issues, an Objective integration Check, and a dependent Objective.
- [ ] Test every new diagnostic name and repair, the health distinction for advisory Issues, no-write behavior across all paths, and deterministic ordering.
- [ ] Run the focused `internal/data` and `internal/doctor` suites, then `go vet ./...`, `git diff --check`, and `make build && make test`.

## Context Log

Pending.
