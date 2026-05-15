---
id: v1.2/D012-legacy-implementation-stage-parse-error
release: v1.2
status: resolved
severity: high
title: "Legacy implementation stage blocks task board operations"
---

# D012: Legacy Implementation Stage Blocks Task Board Operations

## Symptom

In another project, moving a completed task from Planned to Done fails with a task parse error:

`invalid stage "implementation": use build, test, or audit`

Observed file:

`C:\Users\User\Projects\Kids Tutor\.savepoint\releases\v1\epics\E02-retention-and-observability\tasks\T006-stale-heartbeat-alert.md`

## Expected Behavior

The board should not be blocked by older task files that contain a legacy `stage: implementation` value when the task is not actively in the canonical `in_progress` workflow. Users should be able to complete the task or receive an actionable repair path that does not prevent the intended board action.

## Reproduction

1. Open a Savepoint project containing a task file with `stage: implementation`.
2. Attempt to move the task from Planned to Done in the board.
3. The operation fails while parsing the task file with `invalid stage "implementation": use build, test, or audit`.

## Impact

Projects with legacy task frontmatter can become stuck after the stricter stage validation introduced by D011. A user cannot complete affected tasks from the board until the frontmatter is manually repaired.

## Fix Plan

1. Inspect task frontmatter parsing and lifecycle validation for legacy or out-of-workflow `stage` values.
2. Preserve canonical `stage: build|test|audit` requirements for `status: in_progress`.
3. Allow planned/done task operations to recover from legacy `stage: implementation` where it is not part of the active lifecycle, or report it through doctor with a concrete repair.
4. Add regression coverage for a planned/done task with legacy `stage: implementation` and for invalid in-progress stages.

## Acceptance Criteria

- [x] Board/data operations can parse or repair legacy non-in-progress tasks with `stage: implementation`.
- [x] `status: in_progress` still rejects invalid stages such as `implementation`.
- [x] Doctor or parser messaging remains actionable and uses canonical `stage` vocabulary.
- [x] `make build && make test` passes.

## Resolution Notes

Resolved. Task parsing now clears stale `stage` values for non-in-progress tasks, including legacy `stage: implementation`, so the board can load and rewrite affected task files during normal transitions. In-progress tasks still require canonical `stage: build`, `stage: test`, or `stage: audit`; invalid active stages still fail parsing. Doctor continues to report stale non-in-progress task stage fields with an actionable remove-stage repair. Verification: `go test ./internal/data ./internal/doctor ./internal/board`, `make build`, and `make test` passed.
