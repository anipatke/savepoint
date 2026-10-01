---
id: T-064
title: Read dependency vulnerabilities and their severity from OSV-Scanner
objective: O-029
status: done
depends_on: [{task: T-059, requires: clear}]
owner_validation:
    required: false
    accepted_check: ""
planned_by: {role: planner, session: planning-o029-20261001}
complexity_tier: medium
complexity_reason: Severity bucketing plus a small change to the severity rules in the model and classification.
check_waiver:
    task: T-064
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-01T09:21:49Z"
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

**Per-criterion outcomes**

1. Value is the number of groups across `results[].packages[].groups`; `max_severity` buckets 9.0+ critical, 7.0+ high, 4.0+ medium, above 0 low, missing/unparsable unknown; five detail counts sum to the total: met (`TestOSVScannerReaderValue`, `all-severities.json` covers each boundary and checks the sum). Decisions beyond the text: a score of 0 or above 10 is unknown, since "above 0 low" leaves 0 without a bucket; a source outside the instance scope is excluded from the count, while a source outside the root is counted but has no affected item.
2. Model admits `medium`, `low`, `unknown`; `validateSeverity` checks all five; `unknown` above zero is a hard blocker; medium and low never block; the "total above zero with severity missing" rule is kept: met (`TestVulnerabilitySeverityCountsAreValidated`, `TestAssessCurrentValueThresholds`, `TestOSVScannerUnknownSeverityBlocks`). Existing `high`/`critical` constant names kept; added `vulnerabilitySeverityKeys`.
3. Affected items point at each lockfile or manifest, repository-relative, with a note naming package, version, and up to three advisory IDs (`+N more`), worst severity first, bounded: met (`TestOSVScannerEvidenceIsWorstFirstAndBounded`). Decision: an unknown group ranks with high (7.0), since it blocks the same way.
4. Empty `results` is value 0; unresolved versions or a scanner error makes the result partial with a reason; provenance records the scanner version when present; the reason states the database was queried at collection time with snapshot `unknown` unless the report gives one: met (`none.json`, `unresolved.json`, `scanner-error.json`). Assumptions: a top-level `version` and `database_snapshot`, a `results[].error` or `packages[].error` string, and an empty package `version` are how the report says these; the plan did not name the fields. Beyond the text: packages with `vulnerabilities` but no `groups` also make the reading partial, so they are not silently undercounted. The database statement is in the reason of every reading, partial or not, within the 200-character bound.
5. Fixtures under `testdata/readers/vulnerabilities/`: met. `multi-ecosystem.json` (Go, npm, Python), `none.json`, `all-severities.json`, `missing-severity.json`, `unresolved.json`, `malformed.json`, `absolute.json`, plus `scanner-error.json`, `ungrouped.json`, `no-results.json`.

**Commands**: `make test-focused TEST=OSVScanner` passed; `make build && make test-fast` exited 0.

**Files read**: Task, `model.go`, `snapshot.go`, `reader_paths.go`, `classification.go`, `collect.go` (lines 1-150, 275-300), `AGENTS.md`, savepoint-task skill, `router.md`.
**Extra reads (logged)**: `reader_jscpd.go` and `reader_jscpd_test.go` to match reader and test conventions; a grep of `internal/codehealth` for `skippedFilesReason`, `unknownVersion`, `sanitizeLine`, `boundText`, `detailMap`, `equalRefs`, `vitestRoot`; `classification_test.go` and `model_test.go` severity cases (also edited).
**Files changed**: `classification.go`, `snapshot.go`, `classification_test.go`, `model_test.go`; new `reader_osv.go`, `reader_osv_test.go`, `testdata/readers/vulnerabilities/*.json`; this Task file.

**Limitations**: the reader is not registered in `Readers`; registration and the end-to-end containment test belong to T-065. No real OSV-Scanner run was made; fixtures follow its documented JSON shape (`results[].source.path`, `packages[].groups[].max_severity`). The version, database-snapshot, and error fields above are assumptions that real reports probably lack, so real reports will read `unknown` and rarely show a scanner error. `max_severity` is read as a numeric CVSS score only; a vector string is unknown. No Task Check was requested and no waiver was recorded.

## Drift Notes

None expected.
