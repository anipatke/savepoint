---
id: T-030
title: Warn when the router still selects finished work
objective: O-014
planned_by: {role: planner, session: o014-design-20260924}
status: done
complexity_tier: medium
complexity_reason: "One new SelectionDiagnostic kind, rendered by four surfaces that currently show selection diagnostics unevenly (the board shows none)."
depends_on: [{task: T-029, requires: clear}]
owner_validation:
    required: false
    accepted_check: ""
check_waiver:
    task: T-030
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-09-24T06:46:22Z"
---

# T-030: Warn when the router still selects finished work

## Outcome

When the router selects a Task or Objective that is already `done`, the
board's Next area, non-TTY output, `savepoint resume`, and doctor all say the
selection is stale and name the record, alongside whatever Next the records
support.

## User Check

Close a Task on disk while the router still selects it, then open the board:
under the Next line a short warning names the finished Task. `savepoint
resume` and `savepoint doctor` say the same.

## Done When

- `ResolveSelection` returns a `SelectionDone` diagnostic (record kind and ID)
  when the selected Task is done, or when no Task is selected and the
  selected Objective is done. Next still shows the selection as T-029 defines;
  the diagnostic only adds the warning.
- `internal/resume` phrases it in `SelectionPhrase`; the board Next area adds
  one diagnostic line under the Next line for any selection diagnostic, from
  the same phrase source; non-TTY prints the same line; doctor reports it as a
  warning with a repair hint (reselect with `p` or edit router selection).
- Existing selection diagnostics are also visible on the board for the first
  time via the same line; tests confirm their wording is unchanged.
- Tests cover done Task, done Objective, not-done selections (no warning),
  and parity between TTY, non-TTY, resume, and doctor.
- `git diff --check`, `make build && make test-fast` pass.

## Context Files

`internal/data/next.go`, `internal/data/next_test.go`,
`internal/resume/evidence.go`, `internal/resume/resume_test.go`,
`internal/board/v2/next_panel.go`, `internal/board/v2/next_panel_test.go`,
`internal/board/v2/plain.go`, `internal/doctor/v2_runtime.go`,
`internal/doctor/v2_runtime_test.go`, `main_board_next_parity_test.go`.

## Design References

Design sections 1 and 8.

## Guardrails

DATA-03, DATA-04, ARCH-02, ARCH-04, STYLE-07, TEST-01, TEST-02, TEST-08.

## Implementation Plan

1. Add `SelectionDone` to `SelectionDiagnosticKind` and set it in
   `ResolveSelection` without altering the returned `Selection`.
2. Add its phrase to `SelectionPhrase`; confirm the board may import that
   phrase (or move phrase data where both can read it) without breaking the
   ARCH-04 package roles; record the choice.
3. Render the diagnostic line in `nextLines`/`renderNext` and doctor.
4. Add tests.

## Boundaries

No automatic router repair; no change to which Next is chosen.

## Technical Verification

Focused tests per package during iteration; `make build && make test-fast`
at handoff.

## Technical Evidence

Preflight: the router selects T-030 and directs continuing O-014 after T-029.
T-030's `requires: clear` dependency on T-029 is supported by T-029's
`status: done` and its explicit owner waiver for low complexity at
`2026-09-24T05:32:21Z`; the recorded waiver satisfies `requires: clear` but
does not claim technical CLEAR. T-030 started from `planned` and is now
`in_progress` at `stage: audit`.

Extra reads before implementation:

- `.savepoint/objectives/O-014-next-area-objective-word-and-router-issue-target/tasks/T-029-next-is-the-router-selection.md` — confirmed the dependency record is done and carries the owner waiver named by the current router route.
- `.savepoint/Guardrails.md` (T-030's DATA-03, DATA-04, ARCH-02, ARCH-04, STYLE-07, TEST-01, TEST-02, TEST-08 entries) — checked applicable parsing, rendering, package-boundary, single-source, and verification constraints.
- `.savepoint/Design.md` sections 1 and 8 — checked the data/rendering ownership and the existing Next-area presentation contract referenced by T-030.
- `rg -n 'type DiagnosticReport|type Problem|HealthFindings|V2ProblemRepair' internal/doctor`, then targeted reads of `internal/doctor/report.go`, `internal/doctor/checks.go`, and `internal/doctor/repairs.go` — locate the existing report categories and repair-hint path needed to add the stale-selection warning without inventing a parallel doctor surface.
- `rg -n 'func writeCompleteV2Project' internal/doctor` and targeted read of `internal/doctor/report_test.go`'s fixture — locate the existing complete-project fixture so the doctor warning test can change only the selected Task's lifecycle and router input.
- `rg -n 'type TaskV2 struct' internal/data` and targeted read of `internal/data/task_v2.go` — check whether the resolved Task exposes its source path, avoiding extra filesystem discovery in the cross-surface fixture.
- `main_resume_matrix_test.go` targeted reads of the matrix case definitions and `SelectionDiagnostic` assertion — the first handoff gate exposed three completed-Task matrix cases whose former no-diagnostic assertion must now verify the additive `SelectionDone` diagnostic while preserving their Next rung.
- `git status --short` — final worktree check after the build gate, to verify the scoped changed files and confirm no tracked build artifacts were generated.

Implementation: `ResolveSelection` returns `SelectionDone` with the original
selection intact for a done Task or a done Objective with no selected Task.
`ResolveNext` continues through the existing selected-record ladder for that
diagnostic, preserving T-029's Next. `resume.SelectionPhrase` owns the warning
wording; the board appends it under Next for both TUI and non-TTY output, and
doctor reports the same phrase as a non-blocking warning with a router repair
hint. Existing unresolved-selection phrases remain unchanged.

Acceptance evidence:

- Done Task → `SelectionDone` names Task T-001, retains Objective O-001, and
  leaves the existing Objective Next rung in place:
  `TestResolveSelection_doneTaskKeepsSelectionAndAddsDiagnostic`.
- Done Objective with no selected Task → `SelectionDone` names Objective
  O-002, retains its existing Next rung; planned Objective and Task selections
  produce no done diagnostic:
  `TestResolveSelection_doneObjectiveWithoutTaskAddsDiagnostic` and
  `TestResolveSelection_notDoneSelectionsHaveNoDoneDiagnostic`.
- Resume names finished Task and Objective through `SelectionPhrase`, renders
  the warning alongside Next, and keeps the existing not-found wording:
  `TestSelectionPhrase_doneSelectionsNameFinishedRecord`,
  `TestRender_doneSelectionWarningAlongsideCurrentNext`, and
  `TestNextPanelAddsSelectionDiagnosticUnderTheExistingNextLine`.
- TUI, non-TTY, resume, and doctor parity for a router-selected done Task is
  covered by `TestDoneTaskSelectionWarningIsSharedAcrossSurfaces`; doctor
  warning category, router path, and repair hint are covered by
  `TestRunV2ChecksWarnsWhenRouterSelectsDoneTask`.

Iteration commands passed:

- `go test ./internal/data -run '^(TestResolveSelection_doneTaskKeepsSelectionAndAddsDiagnostic|TestResolveSelection_doneObjectiveWithoutTaskAddsDiagnostic|TestResolveSelection_notDoneSelectionsHaveNoDoneDiagnostic)$' -count=1`
- `go test ./internal/resume -run '^(TestSelectionPhrase_doneSelectionsNameFinishedRecord|TestRender_doneSelectionWarningAlongsideCurrentNext)$' -count=1`
- `go test ./internal/board/v2 -run '^(TestNextPanelAddsSelectionDiagnosticUnderTheExistingNextLine|TestRenderNextIncludesTheSelectionDiagnostic)$' -count=1`
- `go test ./internal/doctor -run '^TestRunV2ChecksWarnsWhenRouterSelectsDoneTask$' -count=1`
- `go test . -run '^TestDoneTaskSelectionWarningIsSharedAcrossSurfaces$' -count=1`

Handoff gate:

- The first `make build && make test-fast` built successfully but failed the
  resume matrix's previous nil-diagnostic assertion for three completed-Task
  cases. Updated those cases to expect `SelectionDone` while preserving their
  existing Next kind/action.
- `go test . -run '^TestResumeMatrix_everyRungReachedExactlyOnce$' -count=1`
  — passed after updating the matrix expectations.
- Final `make build && make test-fast` — passed.
- `git diff --check` — passed.

Files read: `.savepoint/router.md`, this Task, its owning `Objective.md`,
`agent-skills/savepoint-task/SKILL.md`, the Context Files listed above,
`.savepoint/Guardrails.md` (named rules only), `.savepoint/Design.md`
(sections 1 and 8), T-029's dependency evidence, and the targeted extra files
and searches recorded above.

Files changed: this Task, `internal/data/next.go`,
`internal/data/next_test.go`, `internal/resume/evidence.go`,
`internal/resume/resume_test.go`, `internal/board/v2/next_panel.go`,
`internal/board/v2/next_panel_test.go`, `internal/doctor/v2_runtime.go`,
`internal/doctor/v2_runtime_test.go`, `internal/doctor/report.go`, and
`main_board_next_parity_test.go`.

Also changed `main_resume_matrix_test.go` to record the new additive
`SelectionDone` diagnostic in its three done-Task cases.

Limitations: no optional Task Check has been requested or waived. This
executor has not written a Check or marked T-030 done.

## Drift Notes

Design section 8's one-line contract changes; reconciled in T-034.
