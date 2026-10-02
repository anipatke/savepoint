---
id: T-060
title: Read test results from Go, Vitest, and pytest reports
objective: O-029
status: done
depends_on: [{task: T-059, requires: clear}]
owner_validation:
    required: false
    accepted_check: ""
planned_by: {role: planner, session: planning-o029-20261001}
complexity_tier: medium
complexity_reason: Two report formats (event stream and JUnit XML) with build-error and empty-suite edge cases.
check_waiver:
    task: T-060
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-01T09:11:30Z"
---

# Read test results from Go, Vitest, and pytest reports

## Outcome

A project's existing Go test JSON, Vitest JUnit, or pytest JUnit report becomes a failed-test count with totals and the names of failing tests, and an unusable report is never mistaken for a passing suite.

## User Check

None beyond the Full Objective Check.

## Done When

- **Go test JSON:** the event stream is read line by line. The value is the count of failed tests plus one per package that failed to build (a package `fail` with no test failures and build output), each named in a note. Details: `total_tests`, `passed_tests`, `skipped_tests`, `build_failures`. Affected items point at the package directory (via the Go module helper) with the test name in the note; raw output is never stored. A stream cut off mid-run (tests started without a final package action) is partial with a reason.
- **JUnit (Vitest and pytest share one parser):** the value is failures plus errors. Details: `total_tests`, `skipped_tests`, `errors`. Affected items use the testcase `file` attribute or classname-derived path where present, with the test name in the note. Hostname and system-out content is ignored.
- An empty suite (zero tests) is available with value 0 and `total_tests` 0, and the reason says no tests ran. Malformed JSON or XML is a reader error.
- Provenance records the runner version where the report has it, otherwise `unknown`.
- Fixtures under `testdata/readers/tests/`: passing, failing, build failure, skips, empty suite, truncated stream, malformed, and a monorepo Go report spanning two modules, for each provider as applicable.

## Context Files

`.savepoint/objectives/O-029-six-code-health-capabilities/Objective.md`; `internal/codehealth/collect.go`; `internal/codehealth/model.go`; `internal/codehealth/snapshot.go`; `internal/codehealth/reader_paths.go`; new `internal/codehealth/reader_tests.go`, `internal/codehealth/reader_junit.go`, `internal/codehealth/reader_tests_test.go`, `internal/codehealth/testdata/readers/tests/` fixtures.

## Design References

O-026 Selected Provider Catalogue "Tests"; O-027 classification rules; O-029 Confirmed Design Decisions.

## Guardrails

FS-05, ARCH-04, DEP-01, DEP-02, STYLE-06, STYLE-07, TEST-01, TEST-02, TEST-04, TEST-07, TEST-08.

## Implementation Plan

1. Confirm `ReportInput.Root` and the path helpers exist; return REPLAN REQUIRED if they do not.
2. Implement the Go test JSON reader with `encoding/json` decoding per line.
3. Implement one JUnit reader with `encoding/xml`, keyed for Vitest and pytest provenance.
4. Fixture-driven table tests for every case above.

## Boundaries

No running tests, no gate-status fallback, no new dependency.

## Technical Verification

Focused tests while iterating; `make build && make test-fast` at handoff.

## Technical Evidence

Commands run: `go test ./internal/codehealth -run 'TestGoTest|TestJUnit'` (pass); `make build && make test-fast` (both pass, no failures).

Per-criterion outcomes:
- Go test JSON: `GoTestReader` decodes per line; value = failed tests + one per package that failed with no failing test (build failure); details `total_tests`, `passed_tests`, `skipped_tests`, `build_failures`; evidence is the package directory via `LoadGoModules`/`Resolve` with the test name in the note, no raw output read. A stream with started tests and no package result is partial with a reason. `TestGoTestReader` (pass, fail+skip, build failure, empty, truncated), `TestGoTestReaderMonorepoResolvesEachModule`, `TestGoTestReaderKeepsFailureCountWhenPackageUnresolved`.
- JUnit: one `JUnitReader` for Vitest and pytest; value = failures + errors; details `total_tests`, `skipped_tests`, `errors`; evidence from the `file` attribute, else the class name (pytest dotted module path, trailing test classes dropped); hostname and system-out ignored (`TestJUnitReader` checks no leak). `TestJUnitReaderHonoursScope`.
- Empty suite: value 0, `total_tests` 0, reason "no tests ran" (both readers). Malformed, empty, or non-JUnit input is an error: `TestGoTestReaderRejectsUnusableReports`, `TestJUnitReaderRejectsUnusableReports`.
- Provenance: neither format carries a runner version, so `provider_version` is `unknown`.
- Fixtures: `testdata/readers/tests/` has passing, failing, build failure, skips, empty, truncated, malformed, and monorepo Go reports, plus Vitest/pytest failing, passing, empty, and malformed JUnit.

Files changed: new `internal/codehealth/reader_tests.go`, `reader_junit.go`, `reader_tests_test.go`, fixtures under `testdata/readers/tests/`; status/stage edits to this Task and `.savepoint/router.md`.
Files read: only the Context Files, plus `AGENTS.md` and `agent-skills/savepoint-task/SKILL.md` (extra reads: workflow), `snapshot.go` was a listed Context File.

Limitations: the readers are not yet registered in `Readers` (T-065). Go subtests count as tests in their own right, so a failing subtest and its parent both count. A Go package that failed to build is detected as a package `fail` with no failing test; the newer `build-fail` event is not read. No real Go/Vitest/pytest output was generated; fixtures are hand-written to the documented formats. `gofmt -l` lists three files outside this Task that were already unformatted.

## Drift Notes

None expected.
