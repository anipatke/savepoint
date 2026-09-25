---
id: T-007
title: Present Goals in the V2 board
objective: O-013
planned_by: {role: planner, session: goals-terminology-20260921}
status: done
complexity_tier: medium
complexity_reason: "The V2 board exposes the context in interactive, overlay, reload, and non-TTY paths, all of which must share one vocabulary without changing filtering behavior."
depends_on: [{task: T-006, requires: clear}]
owner_validation:
    required: true
    accepted_check: ""
check_waiver:
    task: T-007
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-09-23T09:13:53Z"
---

# T-007: Present Goals in the V2 board

## Outcome

The V2 interactive and plain board presents the selected delivery context as a
Goal while retaining the existing Objective filtering, selection persistence,
and Release-backed data behavior.

## User Check

Open the V2 board with an existing R### project. Confirm that the selector,
selected-context line, detail view, help, reload/status messages, and plain
output say Goal/Goals; confirm `g` opens the selector and that any retained `r`
alias behaves identically. Confirm Objective and Task membership has not
changed.

## Done When

- Selector, selected-context header, detail overlay, help, and status/error
  messages use the Goal vocabulary.
- `g` is documented and dispatched as the canonical selector shortcut; a
  retained `r` alias has no separate behavior or public wording.
- Interactive and non-TTY output agree on the Goal label and selected record.
- Reload, missing-selection, narrow-terminal, and empty-Goal cases remain
  understandable and deterministic.
- Existing R### data continues to drive the same Objective and Task filtering;
  no rendering path performs new IO or gate evaluation.
- V1 board behavior and historical Release-specific compatibility surfaces are
  left unchanged.

## Context Files

`internal/board/v2/releases.go`, `internal/board/v2/releases_test.go`, `internal/board/v2/model.go`, `internal/board/v2/view.go`, `internal/board/v2/view_test.go`, `internal/board/v2/help.go`, `internal/board/v2/update.go`, `internal/board/v2/objectives.go`, `internal/board/v2/objectives_test.go`, `internal/board/v2/plain.go`, `internal/board/v2/run.go`, `internal/board/v2/run_test.go`, `internal/board/v2/load_test.go`, `internal/board/v2/watch_test.go`.

## Design References

Design sections 2, 3, 6, 8, 11, and 13.

## Guardrails

DATA-02, ARCH-02, ARCH-03, TPL-02, TEST-01, TEST-02, TEST-04, TEST-06,
TEST-08, STYLE-07, STYLE-09, and STYLE-10.

## Implementation Plan

1. Centralize the user-facing Goal labels and shortcut vocabulary at the V2
   board presentation/dispatch boundary rather than scattering new strings.
2. Update selector, header, detail, help, status, reload, and plain-output
   paths while keeping the existing Release-backed model and index calls.
3. Cover the canonical `g` key, optional `r` compatibility alias, selected
   context, empty data, reload, narrow layout, and non-TTY output.
4. Run focused V2 board tests, `git diff --check`, `make build`, and `make test`.

## Boundaries

No data model rename, persisted-file migration, new Goal membership map, gate
policy change, V1 board rewrite, or UI redesign beyond terminology and the
selector shortcut.

## Technical Verification

Focused `go test ./internal/board/v2 ./internal/resume` with interactive,
overlay, reload, narrow-width, shared evidence, and non-TTY assertions;
existing R### fixture load;
`git diff --check`; `make build`; and `make test`.

## Technical Evidence

Start evidence: T-006 is `done` with the explicit owner Task-check waiver
recorded at 2026-09-23T08:50:51Z; O-013 is `in_progress`. The waiver satisfies
T-007's `T-006 requires: clear` dependency, so implementation may start.

Extra reads beyond Context Files, and why:

- `.savepoint/router.md` — confirmed the active route is R-006/O-013/T-007.
- `.savepoint/objectives/O-013-rename-release-to-goals/Objective.md` —
  checked the owning Objective's boundaries and readiness.
- `.savepoint/objectives/O-013-rename-release-to-goals/tasks/T-006-preserve-existing-project-records.md`
  — verified the dependency's completed status and explicit owner waiver.
- `agent-skills/savepoint-task/SKILL.md` — applied the canonical task lifecycle,
  evidence, and write boundaries.
- `.savepoint/Design.md` sections 2, 3, 6, 8, 11, and 13, and the named rule
  rows in `.savepoint/Guardrails.md` — checked the Task's explicit design and
  policy references.
- Targeted `internal/board/v2` symbol searches for `newReleaseDetail` and the
  detail renderer — the Task requires renaming Release detail copy, but its
  Context Files list only the selector call site, not those implementations.
- Targeted `internal/board/v2` search for Release selector status/error
  strings and their assertions — the Task requires Goal wording in status and
  error messages, while its Context Files omit the I/O helper and some focused
  tests.
- `internal/board/v2/detail_view.go`, `io.go`, and `footer_test.go` — their
  detail labels, selection messages, and footer assertions are explicit Task
  acceptance surfaces omitted from Context Files.
- `internal/resume/evidence.go`, `resume.go`, and `resume_test.go` — inspect
  the shared user-facing Release evidence phrases called by the detail overlay
  and keep Goal terminology at one presentation source rather than duplicating
  it; a targeted search found no separate `evidence_test.go`.
- `main_resume_matrix_test.go` — the configured repository test gate identified
  three cross-surface resume assertions still expecting the old Release output;
  inspect those exact expectations because they cover the shared wording used
  by the V2 board.
- Targeted board Next-area search for Release-specific presentation and tests
  — the Task requires interactive/non-TTY vocabulary parity, while the Next
  renderer is not listed in Context Files.
- `internal/board/v2/actions.go` targeted key-table check — `handleKey`
  dispatches record actions before the Goal selector, so `g` must not collide.

Implementation evidence:

- Selector and shortcut: `TestReleaseSelectorOpensOverTheBoardAndStartsOnCurrentRelease`
  proves `g` opens the Goal selector; `TestLegacyReleaseKeyRemainsAnUndisclosedGoalSelectorAlias`
  proves `r` remains an alias without public wording. Help and footer tests
  require only `g` / “Goals” to be advertised.
- Selected context and detail: `TestReleaseSelectionFiltersIndexedObjectivesAndPersistsOnlyRouterContext`
  checks the `GOAL:` line and unchanged indexed Objective/Task membership;
  detail tests cover Goal headings/readiness and preserve read-only behavior.
- Reload, errors, narrow, and empty states: reload tests require the missing
  selection message to say Goal; cross-Goal selection refusal remains explicit;
  `TestGoalSelectorFitsNarrowBoardWidths` and detail width checks cover narrow
  layouts; the empty selector says “no Goals in this project”.
- Non-TTY: `TestRunWithoutTTYLabelsSelectedReleaseAsGoal` checks the same
  selected Goal identity in plain output.
- The shared `internal/resume` evidence/action phrases now use Goal wording so
  board Next/detail and resume keep one phrase source. Its unit tests and the
  root resume matrix assert the new wording. The `R###` fields, Release-backed
  index/filtering, gate decisions, persisted records, and V1 surfaces were not
  renamed or changed; rendering still resolves no new IO or gate decision.

Files changed for T-007: `internal/board/v2/{releases.go,releases_test.go,
view.go,help.go,update.go,plain.go,run_test.go,footer_test.go,detail_view.go,
io.go}`, `internal/resume/{evidence.go,resume.go,resume_test.go}`, and
`main_resume_matrix_test.go`. No model or index source was modified.

Verification on 2026-09-23:

- `go test ./internal/board/v2 ./internal/resume` — PASS.
- `make build` — PASS.
- `make test` — initial run found three stale Release-worded expectations in
  `main_resume_matrix_test.go`; after updating those assertions, rerun PASS
  across `./...` (the slowest package, `internal/migrate`, took about 1m51s).
- `git diff --check` — PASS after the final edits.

Task is at `stage: audit`, still `status: in_progress`. No Check record was
written; the owner decides whether to request the optional independent Task
Check or provide an explicit waiver, and only the owner may mark this Task done.

## Drift Notes

If the board cannot expose Goals without changing persisted identity or gate
ownership, stop and return REPLAN REQUIRED rather than introducing a second
domain model.
