---
id: T013
title: Find the slow test work
objective: O016
status: planned
depends_on: []
owner_validation: {required: true}
planned_by: {role: planner, session: o016-design-20260923}
complexity_tier: spike
complexity_reason: Gate selection and safe reuse depend on measured test cost and shared-state analysis.
---

# T013: Find the slow test work

## Outcome

A measured baseline and a written gate-selection decision identify the costly tests, the exact fast/full test sets, and the isolation limits for optimization.

## User Check

Review the proposed gate membership and confirm that the full gate still runs every existing test. Confirm the recorded machine and timing method are useful for tracking the stated targets.

## Done When

- Cold and warm package and individual-test timings identify the main cost within internal/migrate, with command, Go version, OS, cache state, and repeated sample results recorded.
- An inventory maps every test to full coverage and names the proposed fast subset. No test is omitted from full coverage.
- The decision names repeated fixture work that can share read-only preparation, mutating scenarios that need fresh directories, and process-global hooks that prevent parallel execution.
- Fast/full command design, timing output, failure behavior, and the baseline against which improvements will be compared are documented for T014-T016.
- The current configured quality gates pass, and each criterion has evidence.

## Context Files

`Makefile`, `.github/workflows/ci.yml`, `internal/migrate/apply_test.go`, `internal/migrate/end_to_end_test.go`, `internal/migrate/fixture_test.go`, `internal/migrate/operation_test.go`, `internal/migrate/cutover_test.go`, `.savepoint/Guardrails.md`, `.savepoint/objectives/O016-fast-test-feedback/Objective.md`.

## Design References

Design sections 1, 12, and 13; O016 Confirmed Verification Policy and Architectural Considerations.

## Guardrails

TEST-01..04, TEST-07..08, CFG-02, STYLE-05, STYLE-07.

## Implementation Plan

1. Record the host/toolchain and run repeatable cold and warm timing samples without editing production code.
2. Attribute slow package time to named tests and fixture operations.
3. Inspect shared hooks and mutation boundaries for safe reuse or parallelism.
4. Produce the gate-selection and optimization decision as Task evidence; flag any target that cannot be met without changing O016's boundaries.
5. Run the currently configured quality gates before handoff.

## Boundaries

Measurement and a decision deliverable only. Do not remove tests or change the verification policy in this Task.

## Technical Verification

Repeat timing commands; verify that every existing test is included in the proposed full selection; `make build && make test`.

## Technical Evidence

Pending execution: timing table, test inventory, decision, command results, and limitations.
