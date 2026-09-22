---
id: T011
title: Create Tasks through the allocator
objective: O019
planned_by: {role: planner, session: i024-allocation-design-20260923}
status: planned
complexity_tier: high
complexity_reason: Creation must preserve authored content and leave the project loadable across failed writes and concurrent calls.
depends_on: [{task: T010, requires: clear}]
owner_validation: {required: false}
---

# T011: Create Tasks through the allocator

## Outcome

One `savepoint create-task` operation accepts an Objective and ID-free Task draft, reserves an ID, writes matching frontmatter and filename, and validates the full V2 index before success.

## Done When

- Refuse an authored ID and occupied file; preserve other supported draft fields and body content.
- Hold the project lock through reservation, file creation, and strict validation. Failure leaves no new Task file; its number stays retired.
- Refuse missing owners and already-invalid projects without modifying existing content.
- Concurrent creation in different Objectives yields distinct IDs and a valid index on Linux and Windows.
- Command dispatch stays thin; filesystem behavior lives in `internal/`.

## Context Files

`main.go`, `main_test.go`, `cmd/init.go`, `cmd/init_test.go`, `internal/data/project.go`, `internal/data/task_v2.go`, `internal/data/write.go`, `internal/data/project_test.go`, `internal/data/task_v2_test.go`, `internal/data/write_test.go`, `AGENTS.md`.

## Design References

Design sections 1, 2, 5, and 6; O019 Architectural Considerations.

## Guardrails

FS-01, FS-05, FS-06, DATA-01, DATA-03, CFG-01, CFG-02, ARCH-01, ARCH-03, ARCH-04, TEST-01, TEST-02, TEST-03, TEST-04, TEST-05, TEST-08.

## Implementation Plan

1. Specify the ID-free draft contract and success/failure output.
2. Add thin parsing and dispatch to an internal creation operation using T010's reservation and lock.
3. Create without overwriting content, then strict-load. On failure remove only the new file and retain the reservation.
4. Test malformed input, occupied paths, write/validation failure, invalid projects, and simultaneous requests.

## Boundaries

No general record editor, Task status transition, or planner-selected ID override.

## Technical Verification

Focused command and internal tests, `git diff --check`, and `make build && make test`.

## Technical Evidence

Pending execution: per-criterion outcomes, command results, files read/changed, and limitations.

## Drift Notes

If the current decoder cannot preserve the ID-free draft's required data, return REPLAN REQUIRED before changing the record schema.
