---
id: v1.2/D014-active-implementation-stage-parse-error
release: v1.2
status: resolved
severity: high
title: "Active legacy implementation stage blocks task parsing"
---

# D014: Active Legacy Implementation Stage Blocks Task Parsing

## Symptom

A Savepoint project can fail to load or update an active task when the task frontmatter contains:

```yaml
status: in_progress
stage: implementation
```

Observed file:

`C:\Users\User\Projects\Kids Tutor\.savepoint\releases\v1\epics\E08-report-card-persona-fixtures\tasks\T002-item-level-add-delete-fixture-workflow.md`

The parser reports:

`invalid stage "implementation": use build, test, or audit. Add 'stage: build' to task frontmatter`

## Expected Behavior

Savepoint should recover from legacy `stage: implementation` task files by loading them as canonical `stage: build`, then writing canonical lifecycle frontmatter during the next board update.

## Reproduction

1. Open a Savepoint project containing an active task with `status: in_progress` and `stage: implementation`.
2. Load or update the task through the board.
3. The task parse fails before Savepoint can rewrite canonical task lifecycle metadata.

## Impact

Projects created or edited with older workflow language can become blocked until users manually repair task frontmatter.

## Fix Plan

1. Map legacy `implementation` stage metadata to canonical `build` during task load.
2. Keep write validation strict so Savepoint only writes `build`, `test`, or `audit`.
3. Keep unknown active stages rejected.
4. Add regression coverage for active `stage: implementation` and `phase: implementation`.

## Acceptance Criteria

- [x] Active tasks with legacy `stage: implementation` parse as `stage: build`.
- [x] Active tasks with legacy `phase: implementation` parse as `stage: build`.
- [x] Unknown active stages still fail parsing.
- [x] `go test ./internal/data` passes.

## Resolution Notes

Resolved. Task load compatibility now maps legacy active `implementation` stage metadata to canonical `build`, including both `stage: implementation` and legacy `phase: implementation`. Write validation remains canonical-only, so subsequent board writes persist `stage: build`, `stage: test`, or `stage: audit`. Verification: `go test ./internal/data`, `make build`, and `make test` passed.
