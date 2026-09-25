---
id: T-029
title: Make Next show exactly what the router selects
objective: O-014
planned_by: {role: planner, session: o014-design-20260924}
status: done
complexity_tier: medium
complexity_reason: "Mostly deletion of the project- and Release-wide search rungs in ResolveNext, but it changes many resume matrix expectations that must each be updated deliberately."
depends_on: [{task: T-028, requires: clear}]
owner_validation:
    required: false
    accepted_check: ""
check_waiver:
    task: T-029
    reason: low complexity
    actor:
        role: owner
        session: owner-chat-20260924
    recorded_at: "2026-09-24T05:32:21Z"
---

# T-029: Make Next show exactly what the router selects

## Outcome

Next is the router's selected Objective and Task, with gate state from the
existing resolvers. No search picks other work, so the board can no longer
show a Task the router did not choose.

## User Check

Select O-018 (no Tasks) in the router while another Objective has a Task in
progress: Next names O-018, not the other Task. Clear the Objective
selection: Next says nothing is selected.

## Done When

- Selected unfinished Task → today's `resolveTaskRung` (unchanged).
- Selected Objective with no selected Task, or a selected Task that is done
  → the Objective: its integration Check/owner-wait rung when every owned
  Task is done; a plan request when it owns no Tasks; otherwise an Objective
  Next with no Task (the owner or an agent selects one).
- No Objective selected → a "nothing selected" kind, unless a selected
  Release's members are all done, where the existing Release completion
  rungs still apply.
- Removed as sources of Next: `resolveProjectWideIntegrationRung`,
  `resolveReadyRung`, `resolveReleaseActiveTaskRung`,
  `resolveReleaseProjectWideIntegrationRung`, `resolveReleaseReadyRung`
  (delete them if nothing else uses them).
- An unresolved selection (not found, mismatch, wrong Release) yields the
  "nothing selected" kind plus its existing diagnostic, not other work.
- Resume phrasing covers plan-this-Objective, select-a-Task, and
  nothing-selected distinctly. The nothing-selected and select-a-Task
  phrases say how to fix it: press `p` on the board, or ask the agent to
  "set router to O-### T-###"; the phrase text lives in resume's data, not
  inline logic (STYLE-09).
- `next_test.go` and `main_resume_matrix_test.go` expectations change only
  where search behavior was removed; evidence lists each changed case. New
  tests reproduce the two 2026-09-23 incidents.
- `git diff --check`, `make build && make test-fast` pass.

## Context Files

`internal/data/next.go`, `internal/data/next_test.go`,
`internal/resume/resume.go`, `internal/resume/resume_test.go`,
`main_resume_matrix_test.go`, `main_board_next_parity_test.go`.

## Design References

Design sections 4 and 8.

## Guardrails

DATA-02, STYLE-05, STYLE-07, STYLE-10, TEST-01, TEST-02, TEST-05, TEST-06,
TEST-08.

## Implementation Plan

1. Write the two incident tests (they fail today).
2. Rewrite `resolveLadder`/`resolveReleaseLadder` to the rules above; add the
   `NextKind` values needed (for example `NextSelectTask`,
   `NextNothingSelected`) and keep `NextPlanObjective` for a selected
   Task-less Objective.
3. Delete the unused search rungs.
4. Update resume phrasing and matrix expectations.

## Boundaries

No gate or completion-policy change, no router writes (T-032), no Issue or
stale-diagnostic work, no mandatory-Goal behavior (O-022).

## Technical Verification

Focused data and resume tests during iteration; `make build && make
test-fast` at handoff.

## Technical Evidence

Preflight: T-029 depends on T-028 with `requires: clear`. T-028 is `done`
with the owner's explicit low-complexity Task Check waiver; the dependency
resolver treats that waiver as satisfying `requires: clear`, while it remains
distinct from technical `CLEAR`. O-014 depends on O-012, which is `done` with
current Check C-907. The start is allowed. The worktree already contained the
completed T-028 changes before T-029 implementation began.

Extra reads before implementation:

- `.savepoint/objectives/O-014-next-area-objective-word-and-router-issue-target/tasks/T-028-print-one-copyable-next-line.md` — verified T-029's `requires: clear` dependency and recorded owner waiver.
- `.savepoint/objectives/O-012-simplify-task-card-outcomes/Objective.md` — verified O-014's Objective dependency remains done with current Check C-907.
- `internal/data/gate_v2.go` and `internal/data/dependency.go` — targeted reads/searches to verify the runtime start and Task-dependency treatment of an explicit waiver.
- `rg -n '\b(resolveProjectWideIntegrationRung|resolveReadyRung|resolveReleaseActiveTaskRung|resolveReleaseProjectWideIntegrationRung|resolveReleaseReadyRung|objectiveDependenciesSatisfied|NextReady)\b' --glob '*.go'` — verified the removed search helpers and `NextReady` have no consumers outside the old projection/matrix/resume code being changed.
- `rg -n 'NextPlanObjective' --glob '*.go'` — found two board tests that name the old no-selection kind; opened them to update their empty-project expectation to `NextNothingSelected`.
- `internal/board/v2/load_test.go` and `internal/board/v2/next_panel_test.go` — opened the two out-of-context tests found by that search to update the no-selection assertions.
- `internal/board/v2/next_panel.go` — opened to verify Objective Check-word styling still applies to the new Objective-ready rung used by the shared Next line.
- `internal/board/v2/releases_test.go` — opened after the first `make test-fast` run showed one release-selector expectation still chose a member Task without an Objective selection; updated it to expect `Nothing selected`.
- `git status --short` — distinguished the pre-existing T-028 worktree changes from T-029 changes.
- `.savepoint/router.md` — read the current task route before updating its next action for the user's T-029 selection.

Implementation: `data.ResolveNext` now resolves only the router's Objective
and optional Task. An unfinished selected Task still goes through
`resolveTaskRung`; a selected Objective with no Tasks asks for planning under
it, an Objective with unfinished Tasks asks the owner/agent to select one,
and an Objective whose Tasks are done reports its existing integration gate
or an owner-ready completion state. An empty or unresolved selection returns
`NextNothingSelected`. A selected Release reaches its unchanged completion
rungs only after the Release resolver reports no unfinished member Objective.
The project-wide and Release-wide integration, active-Task, and ready searches
were removed. Resume keeps the plan, select-a-Task, and nothing-selected
instructions as shared copy data; board styling treats Objective-ready as a
Check line.

Acceptance evidence:

- Selected unfinished Tasks still use the existing start/advance/completion
  resolvers: `TestResolveNext_executeWhenTaskPlannedAndReady`,
  `TestResolveNext_executeWhenTaskInProgress`, and the check/owner/dependency
  cases remain green.
- Selected Task-less O-018 with another Objective's T-006 no longer exposes
  the unselected Task: `TestResolveNext_selectedObjectiveDoesNotChooseUnrelatedReadyTask`.
  The active-work form is covered by
  `TestResolveNext_selectedTasklessObjectiveWinsWhenOtherObjectiveHasActiveTask`.
  `TestResolveNext_doneSelectedTaskKeepsObjectiveAheadOfReleaseActiveTask`
  reproduces the selected O-018/T-020 versus active R-006 T-006 incident and
  returns O-018's integration Check. Both incident shapes also run through
  `TestResumeMatrix_everyRungReachedExactlyOnce`, the board/resume parity
  tests, and the built-command parity test.
- An Objective with unfinished owned Tasks returns `NextSelectTask` with the
  Objective and no Task in
  `TestResolveNext_selectedObjectiveWithIncompleteTasksRequestsTaskSelection`.
  An Objective with no Tasks returns `NextPlanObjective`; a selected Task
  already done returns its Objective's integration or completion decision.
  `TestResolveNext_selectedObjectiveReadyUsesCompletionResolver` proves the
  ready state carries the existing allowed gate decision and current Check.
- No Objective selection leaves unrelated ready or active work unselected:
  `TestResolveNext_noObjectiveSelectionDoesNotPickReadyTask`,
  `TestResolveNext_unselectedObjectiveIsNotChosen`,
  `TestResolveNext_nothingSelectedWhenNothingSelected`, and
  `TestResolveNext_planObjectiveForEmptyProject`. Missing/mismatched selections
  retain their diagnostic and return no substituted record in
  `TestResolveNext_unresolvedSelectionDoesNotSubstituteAvailableWork` and the
  `TestResolveSelection_*` cases.
- Release completion remains available for a selected Release with no
  Objective only when its member Objective work is complete:
  `TestResolveNext_selectedReleaseProjectsCheckOwnerAndReadyRungs`. Release
  selection with no Objective and incomplete member work now returns
  `NextNothingSelected`; the board behavior is covered by
  `TestReleaseSelectionFiltersIndexedObjectivesAndPersistsOnlyRouterContext`.
- `TestActionPhraseSelectionGuidance` and
  `TestNextLineFormatsEverySelectionShape` cover the distinct selected-plan,
  select-a-Task, nothing-selected, integration, and owner-ready text. The
  no-selection and select-a-Task phrases both explain `p` and
  `set router to O-### T-###`.
- The package-wide helper search recorded above found no remaining callers for
  the removed search rungs. The resume matrix reaches each new Next kind and
  verifies determinism and no writes; board and resume assert the same line.

Verification:

- Regression-first run of the three incident tests failed on the pre-change
  resolver as expected: it returned `ready` for selected Task-less O-018 and
  `execute/T-006` for the done T-020 Release-selection case. The fixture's
  Release maps were initialized, then those tests passed with the new
  selection behavior.
- `go test . ./internal/data ./internal/resume ./internal/board/v2 -run '^(TestResolveNext|TestNextLine|TestRender_|TestActionPhraseSelectionGuidance|TestResumeMatrix_|TestBoardNextAndResumeReportTheSameAnswer|TestBuiltBoardAndResumeReportTheSameAnswer|TestLoadProjectEmptyProjectFromTemplateIsNormal|TestNextPanel)' -count=1` — passed.
- `go test ./internal/board/v2 -run '^TestReleaseSelectionFiltersIndexedObjectivesAndPersistsOnlyRouterContext$' -count=1` — passed after updating its stale unselected-member-Task expectation.
- The first `make build && make test-fast` run found that same stale Release-selector assertion. After updating the test to expect `Nothing selected` while retaining the T-002 card in the Release view, the final `make build && make test-fast` run passed.
- `git diff --check` — passed after the implementation and evidence edits.

Files read: `.savepoint/router.md`, `.savepoint/Guardrails.md`,
`.savepoint/objectives/O-014-next-area-objective-word-and-router-issue-target/Objective.md`,
this Task, `agent-skills/savepoint-task/SKILL.md`, the six Context Files
(`internal/data/next.go`, `internal/data/next_test.go`,
`internal/resume/resume.go`, `internal/resume/resume_test.go`,
`main_resume_matrix_test.go`, `main_board_next_parity_test.go`), and the extra
files listed above (`T-028`, O-012, `internal/data/gate_v2.go`,
`internal/data/dependency.go`, and the four board Next/release test/source
files). The search commands and their reasons are recorded in Extra reads.

Files changed: `.savepoint/router.md`, T-028's waiver evidence, this Task,
`internal/data/next.go`, `internal/data/next_test.go`,
`internal/resume/resume.go`, `internal/resume/resume_test.go`,
`main_resume_matrix_test.go`, `internal/board/v2/next_panel.go`,
`internal/board/v2/next_panel_test.go`, `internal/board/v2/load_test.go`, and
`internal/board/v2/releases_test.go`. The existing
`main_board_next_parity_test.go` consumes the expanded matrix without code
changes.

Handoff: T-029 is `done`; the owner waived its optional Task Check for low
complexity at `2026-09-24T05:32:21Z`. This waiver is not technical `CLEAR`.
The mandatory Full Objective Check remains required before O-014 can close.

## Drift Notes

Design section 8 is reconciled in T-034.
