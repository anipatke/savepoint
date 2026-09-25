---
id: v1.2/D015-lifecycle-parse-errors-block-board
release: v1.2
status: resolved
severity: high
title: "Recoverable lifecycle metadata blocks the whole board instead of self-healing"
---

# D015: Lifecycle Parse Errors Block Board Instead of Self-Healing

## Symptom

Loading the board in a user project failed entirely with:

```
parse error for .savepoint/releases/v1.2/defects/D002-incomplete-adaptive-skill-pool.md: stage is required when defect status is in_progress
```

The defect file had `status: in_progress` with no `stage` field. One recoverable
frontmatter omission in one defect file prevented every task and defect in every
release from rendering.

## Expected Behavior

Project principle: **no blocking errors — only notifications and self-healing.**
Recoverable lifecycle metadata (missing stage, stale stage, legacy aliases,
unknown status values) should be healed to a safe canonical value at load time,
with the issue surfaced as a notification (doctor problem, and board status line
once E26/T002 lands) instead of aborting the load.

## Reproduction

1. In any savepoint project, set a defect's frontmatter to `status: in_progress`
   and omit `stage` (an agent saving frontmatter in two steps creates this
   state routinely).
2. Run `savepoint board`.
3. The board exits with a parse error instead of rendering.

The same failure existed for tasks: `status: in_progress` without `stage`
failed `ParseTaskFile` (formerly asserted by
`TestParseTaskFile_rejectsMissingStageForInProgress`, now
`TestParseTaskFile_defaultsMissingStageToBuildForInProgress`).

## Impact

- Total board outage from a single malformed file — the agent/user cannot even
  see the board to find the bad file.
- Recurring defect class: D012 (legacy `implementation` stage), D013 (invalid
  task completion values), and D014 (active `implementation` stage) were all
  the same root cause — strict lifecycle validation on the load path.
- E20/T006 already removed complexity validation from parse for exactly this
  reason; stage/status validation was left behind.

## Fix Plan

Heal at load, stay strict on write, notify via doctor:

- `internal/data/lifecycle.go`: make `ParseTaskLifecycle` heal instead of
  reject — unknown status loads as `planned`, missing/invalid stage on an
  `in_progress` task loads as `build`, stage outside `in_progress` is dropped.
  Add `NormalizeDefectLifecycleForLoad` with the same policy for defects
  (unknown status loads as `open`, task-style aliases `planned`→`open` and
  `done`/`complete`/`completed`→`resolved`), plus `DiagnoseDefectLifecycle`
  returning typed diagnostics mirroring `DiagnoseTaskLifecycle`.
- `internal/data/parser.go`: `ParseDefectFile` normalizes instead of calling
  `validateDefectLifecycle`; lifecycle values can no longer fail a parse.
  Parse errors remain only for structurally broken files (missing/invalid
  YAML frontmatter).
- `internal/data/write.go`: unchanged — `WriteTaskStatus`/`WriteDefectStatus`
  keep strict validation so the app never writes invalid lifecycle values.
  Healed in-memory values persist naturally on the next write.
- `internal/doctor/checks.go`: `checkDefectFile` diagnoses raw frontmatter via
  `DiagnoseDefectLifecycle` (the task-side equivalent already exists at
  `checkTaskLifecycle`), so everything that used to be a blocking parse error
  is reported as a doctor problem instead.
- `internal/doctor/repairs.go`: repair suggestions for the new defect
  lifecycle messages, placed before the generic "release" substring case so
  they are not shadowed by file paths containing `releases/`.
- Out of scope, tracked separately: E26/T002 (tolerant board reload) makes
  structurally broken files skip-and-report with a status-line notification
  instead of aborting, completing the no-blocking-errors principle.

## Acceptance Criteria

- [x] A defect with `status: in_progress` and no `stage` loads as stage `build`; the board renders.
- [x] A task with `status: in_progress` and no `stage` loads as stage `build`; the board renders.
- [x] Unknown defect status loads as `open`; task-style `planned`/`done` aliases map to `open`/`resolved`.
- [x] Unknown task status loads as `planned`; stage outside `in_progress` is dropped on load.
- [x] `savepoint doctor` reports each healed condition as a problem (defect status/stage diagnostics included).
- [x] Write paths still reject invalid lifecycle values (no behavior change in `write.go`).
- [x] `make build && make test` passes.

## Resolution Notes

Load paths now heal instead of reject. `ParseTaskLifecycle` returns a healed
state with no error path; `ParseDefectFile` calls
`NormalizeDefectLifecycleForLoad`. `DiagnoseDefectLifecycle` provides typed
diagnostics consumed by doctor's `checkDefectFile`, which now inspects raw
frontmatter so healed values are still reported, with matching repair
suggestions in `repairs.go`. Write-side validation is unchanged. Parser,
lifecycle, and doctor tests updated from reject-expectations to
heal-expectations; new alias and diagnostic coverage added.

Verified end-to-end against the originally failing project: a copy of the
reporting project with `stage` removed from its in_progress defect renders the
board (exit 0) and `savepoint doctor` reports "defect stage is required when
status is in_progress (loads as build)" with the correct repair suggestion.
`make build && make test` green.
