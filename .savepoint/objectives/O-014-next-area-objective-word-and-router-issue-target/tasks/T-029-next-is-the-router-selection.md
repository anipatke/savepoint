---
id: T-029
title: Make Next show exactly what the router selects
objective: O-014
planned_by: {role: planner, session: o014-design-20260924}
status: planned
complexity_tier: medium
complexity_reason: "Mostly deletion of the project- and Release-wide search rungs in ResolveNext, but it changes many resume matrix expectations that must each be updated deliberately."
depends_on: [{task: T-028, requires: clear}]
owner_validation: {required: false}
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

Pending execution.

## Drift Notes

Design section 8 is reconciled in T-034.
