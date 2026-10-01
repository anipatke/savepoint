---
id: T-061
title: Read statement coverage from Go, Vitest, and coverage.py reports
objective: O-029
status: done
depends_on: [{task: T-059, requires: clear}]
owner_validation:
    required: false
    accepted_check: ""
planned_by: {role: planner, session: planning-o029-20261001}
complexity_tier: medium
complexity_reason: Three formats; Go profiles need import-path mapping and duplicate-block merging.
check_waiver:
    task: T-061
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-01T09:14:40Z"
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

### Per-criterion outcomes (2026-10-01, executor session)

- **Go cover profile:** met. `mode:` required (`go missing mode` rejected); rows merged by block (`go-merged.out`: same block from two runs counts once, covered if any run covered it); import paths resolved through `LoadGoModules`/`Resolve` (monorepo test maps two modules); unresolvable rows make the result partial with a reason (`go-unresolvable.out`). Files outside the instance scope are dropped quietly, not partial (resolution runs scope-free, scope applied after).
- **Vitest V8:** met. Statements from `s`; `covered/total_functions` from `f`; `covered/total_branches` from `b` (each arm counted). Absolute keys become repo-relative via `RelPath`; a key outside the root makes the result partial.
- **coverage.py:** met. Value from `totals.covered_lines`/`num_statements`; branch details only when `meta.branch_coverage` is true; `meta.version` recorded as provider version.
- **Shared:** met. All readers add `covered_statements`/`total_statements`; zero statements is an error for all three; affected items are the least-covered files via `WorstEvidence` (bounded); every reading passes `validateBounded`.
- **Fixtures:** met. `testdata/readers/coverage/` holds populated, fully covered, zero, malformed, Go unresolvable, and Go monorepo cases (plus merged, branches, outside-root, no-totals).

### Commands

- `go test ./internal/codehealth -run 'Coverage|GoCover'`: pass.
- `make build && make test-fast`: both exit 0.

### Files

- Read: Objective.md, collect.go, model.go, reader_paths.go, plus reader_tests.go and reader_tests_test.go as the pattern.
- Extra reads (pattern only, not in Context Files): `reader_tests.go`, `reader_tests_test.go`; `snapshot.go` (Detail/validateBounded) was a listed Context File.
- Changed: new `reader_coverage.go`, `reader_coverage_test.go`, fixtures under `testdata/readers/coverage/`; `.gitignore` (extra edit outside Context Files: added `!internal/codehealth/testdata/readers/coverage/` because the existing `coverage/` rule ignored the mandated fixture directory).

### Limitations

- coverage.py headline uses the report's own `totals`, so instance scope/exclusions do not change the value (they only filter affected items).
- Readers are not registered; registration is T-065.
- No `make test-full` run; not required for ordinary handoff.

## Drift Notes

None expected.
