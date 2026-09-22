---
id: T015
title: Run independent tests together
objective: O016
status: planned
depends_on: [{task: T013, requires: clear}, {task: T014, requires: clear}]
owner_validation: {required: false}
planned_by: {role: planner, session: o016-design-20260923}
complexity_tier: high
complexity_reason: Package-wide hooks and filesystem mutations make parallel safety a test-specific decision.
---

# T015: Run independent tests together

## Outcome

Independent migration test groups execute concurrently, while tests using shared hooks or other process-wide state stay serial and deterministic.

## User Check

If a Task Check is requested, review the concurrency inventory and repeated stress results, especially interruption and recovery cases.

## Done When

- Each parallelized group has a documented isolation basis from T013/T014; global hook and environment users remain serial unless dependencies are made per-run.
- Temporary directories and prepared fixture data are not shared mutably across parallel tests.
- Repeated runs and race-enabled runs show no nondeterminism; named fault-injection and recovery cases still pass.
- Before/after timings and any groups deliberately left serial are recorded.
- Configured quality gates and per-criterion evidence pass.

## Context Files

`internal/migrate/end_to_end_test.go`, `internal/migrate/apply_test.go`, `internal/migrate/operation_test.go`, `internal/migrate/command_test.go`, `internal/migrate/cutover_test.go`, `.savepoint/objectives/O016-fast-test-feedback/tasks/T013-measure-and-design-test-gates.md`, `.savepoint/objectives/O016-fast-test-feedback/tasks/T014-share-read-only-migration-results.md`.

## Design References

Design sections 1 and 13; O016 Architectural Considerations.

## Guardrails

TEST-01..05, TEST-07..08, CFG-02, STYLE-03, STYLE-05.

## Implementation Plan

1. Apply T013's shared-state inventory to each proposed concurrent group.
2. Add parallel execution only to isolated subtests and fixtures.
3. Keep fault-injection and process-global controls serial.
4. Run repeated normal and race-enabled suites; investigate any inconsistent result before handoff.
5. Compare timing and run the configured quality gates.

## Boundaries

No global test hook races, shared mutable fixture data, or assertions that depend on completion order.

## Technical Verification

Repeated focused migration runs, race-enabled migration tests, recovery and interruption cases, `make build && make test`.

## Technical Evidence

Pending execution: isolation matrix, repeated-run results, race results, timing, and command results.
