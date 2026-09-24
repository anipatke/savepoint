---
id: T-030
title: Warn when the router still selects finished work
objective: O-014
planned_by: {role: planner, session: o014-design-20260924}
status: planned
complexity_tier: medium
complexity_reason: "One new SelectionDiagnostic kind, rendered by four surfaces that currently show selection diagnostics unevenly (the board shows none)."
depends_on: [{task: T-029, requires: clear}]
owner_validation: {required: false}
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

Pending execution.

## Drift Notes

Design section 8's one-line contract changes; reconciled in T-034.
