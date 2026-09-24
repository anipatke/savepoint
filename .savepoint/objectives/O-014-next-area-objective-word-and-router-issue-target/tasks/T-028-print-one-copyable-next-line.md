---
id: T-028
title: Print one Next line the owner can copy into a new session
objective: O-014
planned_by: {role: planner, session: o014-design-20260924}
status: planned
complexity_tier: medium
complexity_reason: "One shared line builder used by the board panel, non-TTY output, and resume, with an Objective word added; presentation only."
depends_on: []
owner_validation: {required: false}
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

Pending execution.

## Drift Notes

Design section 8 is reconciled in T-034.
