---
id: T-063
title: Read repeated code from jscpd
objective: O-029
status: planned
depends_on: [{task: T-059, requires: clear}]
owner_validation: {required: false}
planned_by: {role: planner, session: planning-o029-20261001}
complexity_tier: low
complexity_reason: One JSON format; only the totals and clone list are used.
---

# Read repeated code from jscpd

## Outcome

jscpd's JSON report becomes the share of duplicated lines in the configured scope, with the largest repeated blocks named.

## User Check

None beyond the Full Objective Check.

## Done When

- The value is `statistics.total.percentage`, cross-checked against `duplicatedLines`/`lines`; a missing statistics block is a reader error. Details: `duplicated_lines`, `total_lines`, `clones`, `sources`.
- Affected items are the largest clones by lines, bounded, each pointing at the first file and start line with a note naming the second location.
- Zero lines scanned is a reader error with a named reason, never 0%. jscpd's health score, complexity estimate, and dead-code fields are ignored.
- Provenance records the jscpd version when the report includes it, otherwise `unknown`.
- Fixtures under `testdata/readers/duplication/`: mixed Go/TypeScript/Python clones, no clones, zero lines, malformed, and absolute file paths.

## Context Files

`.savepoint/objectives/O-029-six-code-health-capabilities/Objective.md`; `internal/codehealth/collect.go`; `internal/codehealth/model.go`; `internal/codehealth/snapshot.go`; `internal/codehealth/reader_paths.go`; new `internal/codehealth/reader_jscpd.go`, `internal/codehealth/reader_jscpd_test.go`, `internal/codehealth/testdata/readers/duplication/` fixtures.

## Design References

O-026 Selected Provider Catalogue "Duplication"; O-027 classification rules.

## Guardrails

FS-05, ARCH-04, DEP-01, DEP-02, STYLE-06, STYLE-07, TEST-01, TEST-02, TEST-04, TEST-07, TEST-08.

## Implementation Plan

1. Confirm `ReportInput.Root` and the path helpers exist; return REPLAN REQUIRED if they do not.
2. Implement the reader and clone ranking.
3. Fixture-driven table tests.

## Boundaries

No duplication analysis inside Savepoint and no Node or Rust runtime dependency.

## Technical Verification

Focused tests while iterating; `make build && make test-fast` at handoff.

## Technical Evidence

Pending execution.

## Drift Notes

None expected.
