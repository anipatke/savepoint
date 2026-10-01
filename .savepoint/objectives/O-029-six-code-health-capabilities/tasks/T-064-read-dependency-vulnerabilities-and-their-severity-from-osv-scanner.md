---
id: T-064
title: Read dependency vulnerabilities and their severity from OSV-Scanner
objective: O-029
status: planned
depends_on: [{task: T-059, requires: clear}]
owner_validation: {required: false}
planned_by: {role: planner, session: planning-o029-20261001}
complexity_tier: medium
complexity_reason: Severity bucketing plus a small change to the severity rules in the model and classification.
---

# Read dependency vulnerabilities and their severity from OSV-Scanner

## Outcome

OSV-Scanner's JSON report becomes a vulnerability count split into critical, high, medium, low, and unknown severity, with the affected packages named, and unknown severity blocks just as high and critical do.

## User Check

None beyond the Full Objective Check.

## Done When

- The value is the number of distinct vulnerability groups across `results[].packages[].groups`, one per package. Each group's `max_severity` CVSS score sets its bucket: 9.0+ critical, 7.0+ high, 4.0+ medium, above 0 low, missing or unparsable unknown. Details carry all five counts, and they sum to the total.
- The model admits `medium`, `low`, and `unknown` detail keys beside `high` and `critical`; `validateSeverity` checks all five. Classification treats an `unknown` count above zero as a hard blocker; medium and low never block; the existing "total above zero with severity missing" rule still blocks.
- Affected items point at each lockfile or manifest (`source.path` made repository-relative) with a note naming the package, version, and advisory IDs, bounded with the worst severity first.
- An empty `results` list is available with value 0. A source scanned without resolved versions, or a scanner error entry, makes the result partial with a reason. Provenance records the scanner version when present; the reason states that the vulnerability database was queried at collection time and its snapshot is `unknown` unless the report gives one.
- Fixtures under `testdata/readers/vulnerabilities/`: Go, npm, and Python lockfiles in one report; no findings; every severity bucket; missing `max_severity`; manifest without resolved versions; malformed; absolute source paths.

## Context Files

`.savepoint/objectives/O-029-six-code-health-capabilities/Objective.md`; `internal/codehealth/collect.go`; `internal/codehealth/model.go`; `internal/codehealth/snapshot.go`; `internal/codehealth/reader_paths.go`; `internal/codehealth/classification.go`, `internal/codehealth/classification_test.go`, `internal/codehealth/model_test.go`; new `internal/codehealth/reader_osv.go`, `internal/codehealth/reader_osv_test.go`, `internal/codehealth/testdata/readers/vulnerabilities/` fixtures.

## Design References

O-026 Selected Provider Catalogue "Dependency vulnerabilities"; O-027 classification rules; O-029 Confirmed Design Decisions.

## Guardrails

FS-05, ARCH-04, DEP-01, DEP-02, STYLE-06, STYLE-07, TEST-01, TEST-02, TEST-04, TEST-07, TEST-08. DATA-03.

## Implementation Plan

1. Confirm `ReportInput.Root` and the path helpers exist; return REPLAN REQUIRED if they do not.
2. Add the three detail keys, extend `validateSeverity`, and add the unknown-blocks rule in classification with tests.
3. Implement the reader and bucketing.
4. Fixture-driven table tests.

## Boundaries

No vulnerability matching inside Savepoint, no network in tests, and no change to the snapshot schema version.

## Technical Verification

Focused tests while iterating; `make build && make test-fast` at handoff.

## Technical Evidence

Pending execution.

## Drift Notes

None expected.
