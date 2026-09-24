---
id: T-028
title: Print one Next line the owner can copy into a new session
objective: O-014
planned_by: {role: planner, session: o014-design-20260924}
status: done
complexity_tier: medium
complexity_reason: "One shared line builder used by the board panel, non-TTY output, and resume, with an Objective word added; presentation only."
depends_on: []
owner_validation:
    required: false
    accepted_check: ""
check_waiver:
    task: T-028
    reason: low complexity
    actor:
        role: owner
        session: owner-chat-20260924
    recorded_at: "2026-09-24T04:31:23Z"
---

# T-028: Print one Next line the owner can copy into a new session

## Outcome

The board's Next area, the non-TTY output, and the first line of
`savepoint resume` print the same plain-text line naming the Objective and
Task with their words, for example
`In Progress O-014 · Build T-028 — Print one Next line the owner can copy`.

## User Check

Open the board with the router on an in-progress Task: the Next line reads
`NEXT: In Progress O-### · Build T-### — <title>`. Run `savepoint resume`:
its first line is the same text. Copy it into a new agent session.

## Done When

- One line builder, fed by `data.Next`, returns:
  `<Objective word> O-### · <Task word> T-### — <Task title>` when a Task is
  on Next; `<Objective word> O-### · Check — <Objective title>` when no Task
  is and every owned Task is done (including CLEAR awaiting owner
  acceptance); `<Objective word> O-### — <Objective title>` otherwise; and
  `Nothing selected` when Next has no Objective.
- Objective word: `Planned`, `In Progress`, `Done` from Objective status.
  Task word: the existing `taskStageWord` (`Planned`/`Build`/`Test`/`Check`/
  `Done`). No glyphs; only ASCII plus `·` and `—` as shown.
- The board prefixes `NEXT:` and colours the words (existing accents); the
  non-TTY output and resume print the same text. A test asserts board plain
  text, non-TTY, and resume first line are byte-equal for each shape.
- The builder lives where board and resume can both use it without breaking
  ARCH-04 package roles (for example `internal/resume` exporting it, or a
  small `data` formatter); record the choice.
- Tests cover every shape, including a Task with a missing Objective record.
- `git diff --check`, `make build && make test-fast` pass.

## Context Files

`internal/board/v2/next_panel.go`, `internal/board/v2/next_panel_test.go`,
`internal/board/v2/plain.go`, `internal/resume/resume.go`,
`internal/resume/resume_test.go`, `internal/data/next.go`,
`main_board_next_parity_test.go`.

## Design References

Design section 8 (Next area).

## Guardrails

DATA-02, ARCH-02, ARCH-04, STYLE-07, STYLE-09, TEST-01, TEST-02, TEST-08.

## Implementation Plan

1. Add the line builder and its table tests.
2. Use it in `nextLines`/`renderNext`, keeping word colouring.
3. Make it resume's first line; keep the rest of resume's report below it.
4. Add the parity test.

## Boundaries

No change to which records Next picks (T-029) or any gate.

## Technical Verification

Focused board and resume tests during iteration; `make build && make
test-fast` at handoff.

## Technical Evidence

Preflight: T-028 has no Task dependencies; O-014 depends on O-012, whose
Objective record is `done` with current Check C-907. The worktree was clean
before this Task started.

Extra reads before execution:

- `.savepoint/Idea.md` and `.savepoint/Design.md` — read because the router
  still said `design`; used to confirm the already owner-confirmed O-014
  design and resolve the user's explicit selection of T-028.
- `.savepoint/objectives/O-012-simplify-task-card-outcomes/Objective.md` —
  read to verify O-014's Objective dependency is satisfied.
- `agent-skills/savepoint-design/SKILL.md` — read to perform the router's
  design-to-task handoff.
- `rg -n '^type (TaskV2|ObjectiveV2) struct|type ColumnType|ColumnPlanned'
  internal/data` — a package-wide declaration/status search to identify the
  types needed for Objective and Task wording; no additional source file was
  opened from its results.
- `git status --short` — confirmed the worktree was clean before execution.
- `internal/board/v2/view_test.go` — opened after `make test-fast` showed its
  empty-board assertion still expected the retired "Nothing selected yet"
  fallback; needed to keep the existing view regression check aligned with
  the new required text.

Implementation choice: `internal/resume` will export the shared Next-line,
Objective-word, and Task-word builders. `data.ResolveNext` will include the
Task's owning Objective when that record exists, without changing selection
or gate decisions. If the Objective record is absent, the line will report
the Task alone rather than invent an Objective state.

Acceptance evidence:

- `internal/resume.NextLine` produces the requested Objective/Task, Objective
  integration, and empty-selection shapes. `TestNextLineFormatsEverySelectionShape`
  covers Objective and Task statuses, Build/Test/Check stages, a CLEAR
  Objective awaiting owner acceptance, and a Task whose Objective record is
  missing. That fallback keeps the Task visible without inventing an
  Objective state.
- `data.ResolveNext` adds the selected Task's owning Objective record when it
  exists; it does not change the selected rung or gate decisions.
  `TestBoardNextAndResumeReportTheSameAnswer` verifies that relationship for
  resolved Task cases and checks the TUI, non-TTY board, and resume first line
  against one shared line. `TestBuiltBoardAndResumeReportTheSameAnswer`
  checks the command paths as well.
- The board adds `NEXT:` and uses the existing phase accents for Objective,
  Task, and Check words. `TestRenderNextUsesExistingLifecycleAccents` covers
  those style choices; the parity tests cover the visible `NEXT:` prefix.
  The empty-board expectation now matches `Nothing selected` in
  `TestViewEmptyTemplateProjectOpensWithThreeEmptyColumns`.

Verification:

- `go test . ./internal/board/v2 ./internal/resume -run 'TestNextPanel|TestNextArea|TestViewEmptyTemplateProjectOpensWithThreeEmptyColumns|TestBoardNextAndResumeReportTheSameAnswer|TestBuiltBoardAndResumeReportTheSameAnswer|TestNextLineFormatsEverySelectionShape|TestRender_selectionDiagnosticAndNextActionTogether' -count=1` — passed.
- `go test ./internal/board/v2 -run '^TestRenderNextUsesExistingLifecycleAccents$' -count=1` — passed.
- `go test -count=1 -skip '^(TestEndToEnd_temporaryRepositoryCopyMigratesWithReleaseAccountability|TestApply_recoversAtEveryPublishBoundaryWithoutOverwritingUserEdits|TestEndToEnd_goldenIsReproducible)$' ./...` — passed during iteration, before the final accent assertion was added.
- `make build` — passed.
- `make test-fast` — passed on the final run. One intermediate invocation reported a root-package setup failure; the direct full Go gate passed and a repeat of `make test-fast` passed.
- `git diff --check` — passed.

Files read: `.savepoint/router.md`, `.savepoint/objectives/O-014-next-area-objective-word-and-router-issue-target/Objective.md`, `.savepoint/objectives/O-014-next-area-objective-word-and-router-issue-target/tasks/T-028-print-one-copyable-next-line.md`, `.savepoint/Guardrails.md`, `.savepoint/Idea.md`, `.savepoint/Design.md`, `.savepoint/objectives/O-012-simplify-task-card-outcomes/Objective.md`, `agent-skills/savepoint-design/SKILL.md`, `agent-skills/savepoint-task/SKILL.md`, `internal/board/v2/next_panel.go`, `internal/board/v2/next_panel_test.go`, `internal/board/v2/plain.go`, `internal/resume/resume.go`, `internal/resume/resume_test.go`, `internal/data/next.go`, `main_board_next_parity_test.go`, and `internal/board/v2/view_test.go` (extra read noted above). The package-wide `internal/data` search is also recorded above.

Files changed: `.savepoint/router.md`, this Task, `internal/board/v2/next_panel.go`, `internal/board/v2/next_panel_test.go`, `internal/board/v2/view_test.go`, `internal/data/next.go`, `internal/resume/resume.go`, `internal/resume/resume_test.go`, and `main_board_next_parity_test.go`.

Handoff: The owner explicitly waived T-028's optional Task Check for low complexity at `2026-09-24T04:31:23Z` (actor `owner`, session `owner-chat-20260924`). This waiver is not technical `CLEAR`; the mandatory Full Objective Check remains required. The owner has set T-028 to `done`.

## Drift Notes

Design section 8 is reconciled in T-034.
