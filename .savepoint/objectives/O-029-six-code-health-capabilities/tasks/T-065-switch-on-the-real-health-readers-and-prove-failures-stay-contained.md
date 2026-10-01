---
id: T-065
title: Switch on the real health readers and prove failures stay contained
objective: O-029
status: planned
depends_on: [{task: T-060, requires: clear}, {task: T-061, requires: clear}, {task: T-062, requires: clear}, {task: T-063, requires: clear}, {task: T-064, requires: clear}]
owner_validation: {required: false}
planned_by: {role: planner, session: planning-o029-20261001}
complexity_tier: medium
complexity_reason: Cross-reader integration over real fixtures plus Design and Codebase Map reconciliation.
---

# Switch on the real health readers and prove failures stay contained

## Outcome

Code Health offers one registry of the real readers for every approved provider, and a full collection over a mixed project shows each measure reported truthfully even when other tools fail, time out, or produce bad reports.

## User Check

None beyond the Full Objective Check.

## Done When

- `DefaultReaders()` returns a reader for every `ProviderKey` in the catalogue; a test fails if a catalogue key has no reader.
- An end-to-end `Collect` test over a temporary Git project with fixture reports and a fake tool runner configures all nine providers, two instances of one capability, and in the same run one malformed report, one timed-out tool, one unavailable executable, and one absent report. Every valid instance keeps its value and details, and every failed instance keeps its own outcome without a value and is never classified as bad code.
- A saved snapshot from that run round-trips through the store and passes validation.
- The Codebase Map entry for `internal/codehealth` in `AGENTS.md` and the Code Health part of `Design.md` say that production readers exist and that the collection command is still to come.

## Context Files

`.savepoint/objectives/O-029-six-code-health-capabilities/Objective.md`; `internal/codehealth/collect.go`, `internal/codehealth/collect_test.go`, `internal/codehealth/model.go`, `internal/codehealth/storage.go`, `internal/codehealth/reader_tests.go`, `internal/codehealth/reader_junit.go`, `internal/codehealth/reader_coverage.go`, `internal/codehealth/reader_lizard.go`, `internal/codehealth/reader_jscpd.go`, `internal/codehealth/reader_osv.go`; `AGENTS.md`; `.savepoint/Design.md`; new `internal/codehealth/readers.go`, `internal/codehealth/readers_integration_test.go`.

## Design References

Design section 1 and the Codebase Map; O-026 Module Boundary "Result and error states"; O-028 Confirmed Design Decisions.

## Guardrails

FS-05, ARCH-04, TPL-02, CFG-02, CFG-03, STYLE-07, TEST-01, TEST-02, TEST-04, TEST-06, TEST-08.

## Implementation Plan

1. Confirm all five reader Tasks landed; return REPLAN REQUIRED if any reader is missing.
2. Add `DefaultReaders()` and the completeness test.
3. Build the integration test from the reader fixtures and the existing fake runner and Git helpers.
4. Update the Codebase Map and Design wording.

## Boundaries

No collection command, no Check or TUI wiring, and no new reader logic. Shared helper cleanup is limited to removing duplication the readers introduced.

## Technical Verification

Process and path handling make this platform-sensitive: fresh `make test-full` at handoff.

## Technical Evidence

Pending execution.

## Drift Notes

None expected.
