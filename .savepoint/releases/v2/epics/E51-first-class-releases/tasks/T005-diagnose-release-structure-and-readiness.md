---
id: E51-first-class-releases/T005-diagnose-release-structure-and-readiness
status: planned
objective: Report actionable Release integrity and readiness diagnostics from the same indexed records and gate decisions.
depends_on:
  - E51-first-class-releases/T002-require-release-integration-evidence-and-owner-acceptance
complexity_tier: medium
complexity_reason: Adds one consumer of canonical Release diagnostics and gates across doctor checks, reports, and repair guidance.
---

# T005: Diagnose Release structure and readiness

## Problem

Direct file edits can leave Release identity, membership, evidence, acceptance, or historical mappings inconsistent. Doctor must explain those states without creating Issues, repairing files, or inventing a second interpretation of Release readiness.

## Context Files

- `.savepoint/releases/v2/epics/E51-first-class-releases/E51-Detail.md`
- `internal/data/project.go`
- `internal/data/release_v2.go`
- `internal/data/release_gate_v2.go`
- `internal/doctor/interfaces.go`
- `internal/doctor/interfaces_test.go`
- `internal/doctor/checks.go`
- `internal/doctor/checks_test.go`
- `internal/doctor/gates.go`
- `internal/doctor/gates_test.go`
- `internal/doctor/report.go`
- `internal/doctor/report_test.go`
- `internal/doctor/repairs.go`
- `internal/doctor/repairs_test.go`

## Acceptance Criteria

- [ ] Doctor reports malformed or duplicate Release identity, unsafe/mismatched paths, and dangling Objective Release references with file and ID.
- [ ] Doctor distinguishes unassigned Objectives, which are valid, from an Objective that explicitly names a missing Release.
- [ ] Doctor reports empty `in_progress`/`done` Releases, incomplete member Objectives, missing/stale/unknown/NEEDS WORK evidence, unresolved blockers, and missing/stale owner acceptance.
- [ ] A valid historical-completion reference is described as historical evidence, not current CLEAR; malformed or dangling legacy references are diagnosed.
- [ ] Doctor consumes canonical Release/index/gate results and does not reparse Release frontmatter or recalculate completion in its reporting layer.
- [ ] Open advisory Issues do not make an otherwise actionable project generically unhealthy; material blockers are named from the gate result.
- [ ] Diagnostics include manual, non-destructive repair guidance and never create a Release, Issue, Check, or acceptance record.
- [ ] Projects with no Releases produce no Release warning and retain byte-identical existing doctor output where Release context is irrelevant.
- [ ] Human and structured test assertions cover ordering and stable wording for multiple simultaneous Release findings.

## Implementation Plan

- [ ] Expose the indexed Release and canonical readiness information through doctor's consumer interface.
- [ ] Add structural diagnostics for Release records and Objective references.
- [ ] Add gate/evidence/acceptance diagnostics by formatting canonical blocker kinds.
- [ ] Add historical-completion diagnostics without treating archives as current evidence.
- [ ] Add targeted repair suggestions that point to the authoritative record and requirement.
- [ ] Freeze no-Release output and multi-finding order in report tests.
- [ ] Verify all doctor paths remain read-only on success and failure.

## Context Log

Pending.
