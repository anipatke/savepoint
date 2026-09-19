---
id: E48-next-resume/T006-report-health-that-is-worth-acting-on
title: Report health that is worth acting on
status: planned
objective: Split doctor's health reporting into malformed data, missing evidence, failing gates, and pending review, so an advisory backlog stops reading as a broken project.
depends_on: []
complexity_tier: medium
complexity_reason: Reshapes an existing report's categories and its healthy/unhealthy verdict without changing any check's logic.
---

# T006: Report health that is worth acting on

## Problem

`DiagnosticReport.HasProblems` collapses everything doctor finds into one bit. That was survivable when findings were mostly structural, but a V2 project accumulates open Issues as a normal condition of being worked on — E44 made Issues the place durable follow-up lives, and having some open is the healthy state, not a fault. A report that calls such a project unhealthy teaches its owner that the report is noise, and the next time it reports a genuinely malformed record they will skip past it.

So health becomes four categories that mean different things and demand different responses:

- **Malformed data** — a record does not decode, an ID is duplicated, a reference dangles. The project cannot be interpreted correctly until this is fixed.
- **Missing evidence** — records are valid but a target has no Check, or no freshness assessment. Nothing is broken; something has not been done yet.
- **Failing configured gates** — a configured quality gate ran and failed. This is a real failure, and it is the project's own definition of one.
- **Pending semantic review** — work that needs a human or a checker to look at it. Doctor is reporting a queue, not a defect.

Only the first three bear on whether the project is structurally sound. Advisory open Issues are listed and never count against it.

What does not change: doctor creates no Issues, repairs no files, and runs no check it does not run today. The inspectors from E43 and E44 — `InspectTaskConsistency`, `InspectIssueConsistency`, `InspectObjectiveConsistency` — already produce typed diagnostics; this task routes their output into the right category rather than inventing new detection. Keep archive content out of live diagnostics except for reference integrity.

## Context Files

- `internal/doctor/report.go`
- `internal/doctor/report_test.go`
- `internal/doctor/checks.go`
- `internal/doctor/checks_test.go`
- `internal/doctor/gates.go`
- `internal/doctor/interfaces.go`
- `internal/data/gate_v2.go`
- `internal/data/issue_v2.go`
- `internal/data/objective_gate_v2.go`
- `cmd/doctor.go`
- `.savepoint/releases/v2/v2-Design.md`
- `.savepoint/releases/v2/epics/E48-next-resume/E48-Detail.md`

## Acceptance Criteria

- [ ] The diagnostic report carries the four named categories, and every finding doctor produces is assigned to exactly one of them.
- [ ] Structural soundness is derived from malformed data, missing evidence, and failing gates only; pending semantic review never affects it.
- [ ] A V2 project whose only findings are advisory open Issues is not reported as structurally unsound, and its exit status reflects that.
- [ ] Open Issues are still listed, with their type and status, under the pending-review category.
- [ ] Malformed records, duplicate identities, and dangling references report under malformed data, sourced from the existing inspectors rather than new detection logic.
- [ ] A target with no recorded Check and a target with a `stale` or `unknown` freshness assessment report under missing evidence, distinguished from one another in the output.
- [ ] A failing configured quality gate reports under failing gates, retaining its existing timeout and failure output behavior.
- [ ] The formatted output labels all four categories and omits an empty category rather than printing an empty heading.
- [ ] Doctor creates no Issue file, repairs no file, and writes nothing during a report run, verified by a project-tree snapshot before and after (FS-03).
- [ ] Existing V1 doctor behavior and its tests are unchanged.
- [ ] Exit-status behavior for a structurally unsound project is unchanged from today's failing case.
- [ ] `make build && make test` passes.

## Implementation Plan

- [ ] Read `report.go` and enumerate every finding it can emit today, mapping each to one of the four categories before changing any code.
- [ ] Introduce the categories on the report type and route existing findings into them, leaving each check's detection logic untouched.
- [ ] Replace the single `HasProblems` verdict with a soundness derivation over the three blocking categories, keeping the existing exported behavior for V1 callers.
- [ ] Route the E43/E44 inspector output into malformed data and missing evidence.
- [ ] Update `Format()` to label categories and skip empty ones.
- [ ] Add the advisory-backlog-only fixture asserting the project is sound and the Issues are still listed.
- [ ] Add the missing-Check versus stale-assessment distinction test.
- [ ] Add the no-write snapshot assertion around a full report run.
- [ ] Run `go test ./internal/doctor/... ./cmd/...`, then `make build && make test`.

## Context Log

Pending.
