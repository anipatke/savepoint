---
id: T-047
title: Complete a Goal when its Objectives are complete
objective: O-025
status: planned
owner_validation: {required: false}
planned_by: {role: planner, session: o025-plan-20260925}
---

# Complete a Goal when its Objectives are complete

## Outcome

The runtime treats a Goal as complete when every member Objective is
complete. Resume, the board, and doctor never ask for a Goal Check.

## User Check

On this repository, run `savepoint resume`: with R-006 selected and no
Objective selected, Next no longer says `Record a fresh Goal Check`; it gives
the Goal-ready action. Run `savepoint doctor`: no `v2-release-clearance-*`,
`v2-release-checker-authority`, `v2-release-issue-unresolved`, or Goal
owner-acceptance problem appears. Open R-006's detail on the board: no Goal
Check clearance or owner-acceptance section is shown.

## Done When

- `ResolveReleaseCompletion` allows completion when the Goal has member
  Objectives and each passes `ResolveObjectiveCompletion`. It no longer reads
  Goal-scoped clearance, Issues linked to a Goal Check, a Goal exception, or
  Goal owner acceptance. Legacy completion and the no-Objectives and
  incomplete-member blockers behave as before.
- `NextReleaseCheckNeeded`, `NextReleaseOwnerValidationRequired`, and
  `NextReleaseIntegration` are removed along with their resume phrases and
  evidence lines, unless a remaining blocker still needs one. Next for a
  complete Goal is `NextReleaseReady`.
- Doctor's Goal problems are limited to the membership blockers the gate can
  still return.
- The board Goal detail drops the Goal Check clearance and
  `releaseOwnerValidationLines` output.
- A `scope.kind: release` Check still decodes and loads, and changes no
  decision (tested).
- Tests cover: all members complete → allowed with no Goal Check; a member
  incomplete → blocked; a leftover NEEDS WORK Goal-scoped Check → still
  allowed; legacy completion unchanged; resume, board, and doctor output for
  each.
- `make build && make test-fast` pass.

## Context Files

`internal/data/release_gate_v2.go`, `internal/data/release_gate_v2_test.go`,
`internal/data/next.go`, `internal/data/next_test.go`,
`internal/data/gate_v2.go`, `internal/resume/resume.go`,
`internal/resume/evidence.go`, `internal/resume/resume_test.go`,
`internal/board/v2/detail.go`, `internal/board/v2/detail_view.go`,
`internal/board/v2/detail_test.go`, `internal/doctor/checks.go`,
`internal/doctor/checks_test.go`, `main_resume_matrix_test.go`,
`.savepoint/objectives/O-025-remove-the-goal-check/Objective.md`,
`.savepoint/issues/I-060-remove-the-goal-check.md`.

## Design References

Design Goal boundary (Section 1), Check workflow, and Section 10 verification
order; reconciled by the guidance Task.

## Guardrails

ARCH-01, DATA-02, STYLE-07, TEST-01..04, TEST-08.

## Implementation Plan

1. Reduce `resolveReleaseCompletionForRecord` to the legacy, no-Objectives,
   and member-Objective checks; return allowed after members pass. Delete
   helpers left without callers (`releaseExceptionCoversIssue`, and
   `checkIssueIDs` only if nothing else uses it).
2. Simplify `resolveReleaseCompletionRung` and remove Next kinds that can no
   longer occur; update `resolveReleaseLadder`.
3. Remove the Goal Check phrases and evidence lines in `internal/resume`.
4. Remove Goal Check clearance and owner-acceptance lines from the board Goal
   detail.
5. Drop the unreachable Goal clearance cases from `releaseBlockerProblem`.
6. Update and add tests; run focused tests, then
   `make build && make test-fast`.

## Boundaries

No change to Objective or Task gates, storage names, Check decoding, or
guidance prose. No board action to close a Goal.

## Technical Verification

Focused tests during iteration; `make build && make test-fast` at handoff;
the Full Objective Check requires `make test-full`.

## Technical Evidence

Pending execution.

## Drift Notes

None expected; Design is reconciled by the guidance Task.
