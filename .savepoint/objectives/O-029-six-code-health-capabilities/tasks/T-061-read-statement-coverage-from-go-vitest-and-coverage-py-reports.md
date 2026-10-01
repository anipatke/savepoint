---
id: T-061
title: Read statement coverage from Go, Vitest, and coverage.py reports
objective: O-029
status: planned
depends_on: [{task: T-059, requires: clear}]
owner_validation: {required: false}
planned_by: {role: planner, session: planning-o029-20261001}
complexity_tier: medium
complexity_reason: Three formats; Go profiles need import-path mapping and duplicate-block merging.
---

# Read statement coverage from Go, Vitest, and coverage.py reports

## Outcome

A project's existing coverage report becomes one statement-coverage percentage that means the same thing on every stack, with the least-covered files named.

## User Check

None beyond the Full Objective Check.

## Done When

- **Go cover profile:** `mode:` is required; rows are merged by block so the same block from several packages is counted once; covered and total statements give the value. Import paths become repository paths through the Go module helper; rows that cannot be resolved make the result partial with a reason.
- **Vitest V8 `coverage-final.json`:** statements come from each file's `s` counters; details add `covered_functions`/`total_functions` from `f` and `covered_branches`/`total_branches` from `b`. Absolute file keys become repository-relative.
- **coverage.py JSON:** the value comes from `totals.covered_lines` and `num_statements` (statement-based); details add branch counts when `meta.branch_coverage` is true; provenance records `meta.version`.
- Every reader adds `covered_statements` and `total_statements` details. Zero total statements is a reader error with a named reason, never 0% or 100%. Affected items are the least-covered files, bounded.
- Fixtures under `testdata/readers/coverage/`: populated, fully covered, zero statements, malformed, Go profile with an unresolvable import path, and a monorepo Go profile across two modules.

## Context Files

`.savepoint/objectives/O-029-six-code-health-capabilities/Objective.md`; `internal/codehealth/collect.go`; `internal/codehealth/model.go`; `internal/codehealth/snapshot.go`; `internal/codehealth/reader_paths.go`; new `internal/codehealth/reader_coverage.go`, `internal/codehealth/reader_coverage_test.go`, `internal/codehealth/testdata/readers/coverage/` fixtures.

## Design References

O-026 Selected Provider Catalogue "Coverage"; O-027 classification rules; O-029 Confirmed Design Decisions.

## Guardrails

FS-05, ARCH-04, DEP-01, DEP-02, STYLE-06, STYLE-07, TEST-01, TEST-02, TEST-04, TEST-07, TEST-08.

## Implementation Plan

1. Confirm `ReportInput.Root` and the path helpers exist; return REPLAN REQUIRED if they do not.
2. Implement the three readers with a shared statement total and least-covered ranking.
3. Fixture-driven table tests.

## Boundaries

No line-coverage headline, no instrumentation, no running tests.

## Technical Verification

Focused tests while iterating; `make build && make test-fast` at handoff.

## Technical Evidence

Pending execution.

## Drift Notes

None expected.
