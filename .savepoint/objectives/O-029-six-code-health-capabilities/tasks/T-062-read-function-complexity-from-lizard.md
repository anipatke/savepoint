---
id: T-062
title: Read function complexity from Lizard
objective: O-029
status: done
depends_on: [{task: T-059, requires: clear}]
owner_validation:
    required: false
    accepted_check: ""
planned_by: {role: planner, session: planning-o029-20261001}
complexity_tier: low
complexity_reason: One headerless CSV format with fixed columns.
check_waiver:
    task: T-062
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-01T09:16:57Z"
---

# Read function complexity from Lizard

## Outcome

Lizard's CSV report becomes the highest function complexity in the configured scope, with the most complex functions named and their files and lines given.

## User Check

None beyond the Full Objective Check.

## Done When

- The headerless Lizard CSV columns (NLOC, CCN, token, PARAM, length, location, file, function, long name, start, end) are read by position with `encoding/csv`; a row with the wrong column count or a non-numeric CCN is a reader error.
- The value is the maximum CCN. Details: `functions`, `average_ccn`, and `functions_over_10` and `functions_over_20` (the O-027 Good/Watch defaults).
- Affected items are the most complex functions, bounded, with file, start line, and a note naming the function and its CCN.
- A report with no functions is available with value 0 only when the reason says no functions were found. Rows for files outside the root are dropped and make the result partial.
- Lizard reports no version in its CSV, so provenance version is `unknown` and the reason says so.
- Fixtures under `testdata/readers/complexity/`: Go, Python, and TypeScript in one report; empty; malformed row; a row outside the root.

## Context Files

`.savepoint/objectives/O-029-six-code-health-capabilities/Objective.md`; `internal/codehealth/collect.go`; `internal/codehealth/model.go`; `internal/codehealth/snapshot.go`; `internal/codehealth/reader_paths.go`; new `internal/codehealth/reader_lizard.go`, `internal/codehealth/reader_lizard_test.go`, `internal/codehealth/testdata/readers/complexity/` fixtures.

## Design References

O-026 Selected Provider Catalogue "Complexity"; O-027 classification rules.

## Guardrails

FS-05, ARCH-04, DEP-01, DEP-02, STYLE-06, STYLE-07, TEST-01, TEST-02, TEST-04, TEST-07, TEST-08.

## Implementation Plan

1. Confirm `ReportInput.Root` and the path helpers exist; return REPLAN REQUIRED if they do not.
2. Implement the reader and ranking.
3. Fixture-driven table tests.

## Boundaries

No complexity analysis inside Savepoint and no cognitive-complexity claims.

## Technical Verification

Focused tests while iterating; `make build && make test-fast` at handoff.

## Technical Evidence

**Per-criterion outcomes**

1. Positional `encoding/csv` read, wrong column count or non-numeric CCN is an error: met. `FieldsPerRecord = 11`; `TestLizardReaderRejectsMalformedReports` covers wrong columns, bad CCN, and (beyond the criterion) a non-numeric start line.
2. Value is maximum CCN; details `functions`, `average_ccn`, `functions_over_10`, `functions_over_20`: met (`TestLizardReaderValue`, mixed report: 22, 5, 8, 2, 1). The over-counts use `defaultThresholds[CapabilityComplexity]` (Good 10, Watch 20) rather than repeating the numbers.
3. Affected items are the most complex functions, bounded, with file, start line, and a note naming function and CCN: met (`WorstEvidence`, `TestLizardEvidenceIsBounded` checks the 20 cap). Interpretation: only functions above the Good limit (CCN > 10) are listed, so a clean project has no affected items, as `leastCovered` does for fully covered files.
4. No functions gives value 0 with a reason saying none were found; rows outside the root are dropped and make the result partial: met. Decision beyond the text: if every row is outside the root, nothing was measured, so the result is partial with no value rather than 0 (Objective: "never as complexity zero"). Rows inside the root but outside the instance scope are dropped quietly, not partial.
5. Provenance version `unknown`, reason says Lizard's CSV states no version: met on available results. A partial result's reason names the dropped rows first and then carries the same sentence.
6. Fixtures under `testdata/readers/complexity/`: met — `mixed.csv` (Go, Python, TypeScript), `empty.csv`, `wrong-columns.csv`, `bad-ccn.csv`, `bad-line.csv`, `outside.csv`, `all-outside.csv`.

**Commands**: `make test-focused TEST=Lizard` passed; `make build && make test-fast` exited 0.

**Files read**: Task, Objective, `collect.go`, `model.go`, `snapshot.go`, `reader_paths.go`, `AGENTS.md`, savepoint-task skill, `router.md`.
**Extra reads (logged)**: `reader_coverage.go` (first 150 lines) and `reader_coverage_test.go` (first 140 lines) to match reader and test conventions; `classification.go` lines 22-40 and 270-290 and a grep of `internal/codehealth` for shared constants, to reuse `defaultThresholds`, `unknownVersion`, `skippedFilesReason`, `detailMap`, and `equalRefs`.
**Files changed**: new `internal/codehealth/reader_lizard.go`, `reader_lizard_test.go`, `testdata/readers/complexity/*.csv`; this Task file.

**Limitations**: the reader is not registered in `Readers`; registration and the end-to-end containment test belong to T-065. No real Lizard run was made (no network or tool execution in tests); fixtures follow the documented column order. `functions_over_10`/`_20` are named for the default limits and do not follow a project's custom thresholds.

## Drift Notes

None expected.
