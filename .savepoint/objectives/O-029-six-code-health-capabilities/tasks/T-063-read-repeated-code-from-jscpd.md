---
id: T-063
title: Read repeated code from jscpd
objective: O-029
status: done
depends_on: [{task: T-059, requires: clear}]
owner_validation:
    required: false
    accepted_check: ""
planned_by: {role: planner, session: planning-o029-20261001}
complexity_tier: low
complexity_reason: One JSON format; only the totals and clone list are used.
check_waiver:
    task: T-063
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-01T09:18:40Z"
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

**Per-criterion outcomes**

1. Value is `statistics.total.percentage`, cross-checked against `duplicatedLines`/`lines`; missing statistics is an error; details `duplicated_lines`, `total_lines`, `clones`, `sources`: met (`TestJscpdReaderValue`, `TestJscpdReaderRejectsUnusableReports`). Decision beyond the text: the cross-check allows 0.5 points of difference (jscpd rounds to two decimals), and a total missing `lines`, `duplicatedLines`, or `percentage`, or with `duplicatedLines` > `lines`, is also a reader error.
2. Affected items are the largest clones by lines, bounded, at the first file and start line, with a note naming the second location: met (`WorstEvidence`; `TestJscpdEvidenceIsBounded` checks the 20 cap). Decisions beyond the text: a clone with either file outside the root is left out and makes the reading partial; a clone is kept when at least one of its files is inside the instance scope.
3. Zero lines scanned is a reader error naming the reason, never 0%; health score, complexity estimate, and dead-code fields are ignored: met (`zero-lines.json`; the report struct reads none of those fields).
4. Provenance records the jscpd version when the report includes it, otherwise `unknown`: met (`versioned.json`, others). Assumption: the version would be a top-level `version` field; jscpd's JSON normally has none, so real reports will read `unknown`.
5. Fixtures under `testdata/readers/duplication/`: met. `mixed.json` (Go, TypeScript, Python), `none.json`, `zero-lines.json`, `malformed.json`, `absolute.json`, plus `no-statistics.json`, `inconsistent.json`, `versioned.json`.

**Commands**: `make test-focused TEST=Jscpd` passed; `make build && make test-fast` exited 0.

**Files read**: Task, Objective status line, `model.go`, `reader_paths.go`, `AGENTS.md`, savepoint-task skill, `router.md`.
**Extra reads (logged)**: `reader_lizard.go` and `reader_lizard_test.go` (whole) and `reader_coverage.go` lines 150-316 to match reader and test conventions; a grep of `internal/codehealth` for `percent`, `skippedFilesReason`, `unknownVersion`, `detailMap`, `equalRefs`, `vitestRoot`, and `sanitizeLine` to reuse them; the T-059 status line to confirm the dependency is done; `config.yml` quality gates (none configured; the Makefile gates ran).
**Files changed**: new `internal/codehealth/reader_jscpd.go`, `reader_jscpd_test.go`, `testdata/readers/duplication/*.json`; this Task file.

**Limitations**: the reader is not registered in `Readers`; registration and the end-to-end containment test belong to T-065. No real jscpd run was made; fixtures follow jscpd's documented JSON shape. A real jscpd report may carry its version somewhere other than a top-level `version`, in which case provenance says `unknown`. No Task Check was requested and no waiver was recorded.

## Drift Notes

None expected.
