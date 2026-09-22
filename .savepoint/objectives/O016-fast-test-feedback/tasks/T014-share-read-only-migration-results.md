---
id: T014
title: Reuse migration fixture results
objective: O016
status: planned
depends_on: [{task: T013, requires: clear}]
owner_validation: {required: false}
planned_by: {role: planner, session: o016-design-20260923}
complexity_tier: high
complexity_reason: Migration assertions share expensive setup while mutation and recovery evidence must stay isolated.
---

# T014: Reuse migration fixture results

## Outcome

Read-only end-to-end migration assertions reuse independently prepared converted fixture results where safe, reducing repeated copies and full applies without weakening their assertions.

## User Check

If a Task Check is requested, compare the before/after test inventory and a representative converted fixture result; confirm mutating recovery scenarios still use fresh directories.

## Done When

- The T013 decision identifies every assertion moved to shared read-only preparation.
- Shared preparation cannot be mutated by a consumer, and each mutating, interruption, recovery, no-op, and user-file-safety scenario retains isolated input.
- Failure messages still identify the fixture and assertion; test ordering cannot change results.
- Before/after timing and the retained test inventory are recorded.
- Named happy and failure paths, configured quality gates, and per-criterion evidence pass.

## Context Files

`internal/migrate/end_to_end_test.go`, `internal/migrate/apply_test.go`, `internal/migrate/fixture_test.go`, `internal/migrate/operation_test.go`, `.savepoint/objectives/O016-fast-test-feedback/tasks/T013-measure-and-design-test-gates.md`.

## Design References

Design sections 1 and 13; O016 Architectural Considerations.

## Guardrails

FS-01, TEST-01..05, TEST-07..08, STYLE-03, STYLE-07.

## Implementation Plan

1. Use T013's assertion map to group read-only checks by fixture and one prepared result.
2. Keep every write or fault-injection scenario on an independent temporary project.
3. Preserve named subtests and fixture-specific diagnostics.
4. Compare test inventory and timing before and after; rerun recovery and file-preservation cases.
5. Run the currently configured quality gates before handoff.

## Boundaries

No production migration behavior change, skipped assertion, shared mutable fixture, or test-order dependency.

## Technical Verification

Focused migration test cases, full migration package, filesystem preservation and recovery cases, `make build && make test`.

## Technical Evidence

Pending execution: assertion mapping, named cases, isolation proof, timing, and command results.
