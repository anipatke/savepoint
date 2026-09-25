---
id: T-047
title: Complete a Goal when its Objectives are complete
objective: O-025
status: done
owner_validation:
    required: false
    accepted_check: ""
planned_by: {role: planner, session: o025-plan-20260925}
check_waiver:
    task: T-047
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-09-25T20:40:15Z"
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

Extra read: package-wide Go reference search for `checkIssueIDs`,
`mapsKeysString`, and `releaseExceptionCoversIssue`; needed to verify which
planned helpers could be removed safely. `checkIssueIDs` is also used by
`objective_gate_v2.go`, so it must remain; `releaseExceptionCoversIssue` has
no other caller. Extra read: repository-wide Go reference search for the
Goal-specific Next kinds and Release acceptance phrase helpers, needed to
confirm their callers before removing the obsolete Goal Check path; this
found callers only in the Context Files listed above.
Extra read: `internal/board/v2/releases_test.go`, required after its package
test exposed an assertion for the removed Goal owner-validation section.

Implementation: Goal completion now returns allowed after all member
Objectives pass, while retaining the empty-members, incomplete-member, and
legacy branches. The release Check, linked Issue, exception, and owner
acceptance paths were removed from that decision. Next no longer has Goal
Check or Goal owner-acceptance rungs; resume and board omit their clearance
and acceptance output; doctor reports membership blockers only.

Acceptance evidence:

- Completion uses member Objective completion only. No-members and incomplete
  member blockers remain; legacy completion remains; absent, stale,
  NEEDS WORK, and linked-Issue Goal Check evidence does not change the result.
  Covered by `TestResolveReleaseCompletion_requiresMemberObjectivesOnly`,
  `TestResolveReleaseCompletion_ignoresEveryGoalCheckClearanceState`,
  `TestResolveReleaseCompletion_ignoresGoalCheckAndItsIssues`,
  `TestResolveReleaseCompletion_allowsCompleteMembersWithoutGoalCheck`, and
  `TestResolveReleaseCompletion_legacyCompletionWithLiveMembersIsUnchanged`.
- The Goal Check and owner-acceptance Next kinds and their resume phrases and
  evidence lines are removed. A completed Goal projects `NextReleaseReady`
  without carrying Goal clearance. Covered by
  `TestResolveNext_selectedReleaseWithCompleteMembersIsReadyRegardlessOfGoalCheck`,
  `TestRender_goalReadyUsesObjectiveCompletionWithoutGoalCheckEvidence`, and
  the on-disk resume matrix cases for no Goal Check and a NEEDS WORK Goal
  Check.
- Doctor now maps only no-members and incomplete-member blockers. Its tests
  load a Goal-scoped Check (including NEEDS WORK), an unresolved linked Issue,
  and a Goal with no Check without reporting Goal Check problems.
- Goal detail omits clearance and owner-validation sections while retaining
  membership readiness and Check history. Covered by
  `TestGoalDetailShowsLatestChecksCodeStyleReview` and
  `TestReleaseDetailShowsEvidenceHistoryIssuesAndMembershipReadiness`.
- Release-scoped Checks still decode and load. Covered by
  `TestDecodeCheckV2_releaseScope`, the on-disk resume matrix, and doctor
  loading tests; Goal completion ignores their result and linked Issue.
- `go test ./internal/data ./internal/resume ./internal/board/v2 ./internal/doctor .`
  passed after the stale Goal Check assertions were updated.
- `make build && make test-fast` passed. `git diff --check` passed.
- `savepoint resume` passed after the stage moved to `audit`; Next is
  `Check T-047` and asks the owner to request the optional Task Check or record
  an explicit waiver. The router remains selected on T-047.

Files read: `.savepoint/router.md`; this Task; its owning
`.savepoint/objectives/O-025-remove-the-goal-check/Objective.md`;
`.savepoint/issues/I-060-remove-the-goal-check.md`;
`.savepoint/Guardrails.md` for ARCH-01, DATA-02, STYLE-07, TEST-01..04,
TEST-08; and the Task Context Files `internal/data/release_gate_v2.go`,
`internal/data/release_gate_v2_test.go`, `internal/data/next.go`,
`internal/data/next_test.go`, `internal/data/gate_v2.go`,
`internal/resume/resume.go`, `internal/resume/evidence.go`,
`internal/resume/resume_test.go`, `internal/board/v2/detail.go`,
`internal/board/v2/detail_view.go`, `internal/board/v2/detail_test.go`,
`internal/doctor/checks.go`, `internal/doctor/checks_test.go`, and
`main_resume_matrix_test.go`. Extra reads: `internal/data/objective_gate_v2.go`
was identified by a package-wide helper-reference search and confirms
`checkIssueIDs` must remain; `internal/board/v2/releases_test.go` was read
after its test exposed an obsolete Goal owner-validation assertion. A
repository-wide reference search confirmed the removed Next kinds and Goal
acceptance phrase helpers had no other callers.

Files changed: this Task and its owning Objective lifecycle records;
`internal/data/release_gate_v2.go`, `internal/data/release_gate_v2_test.go`,
`internal/data/next.go`, `internal/data/next_test.go`,
`internal/resume/resume.go`, `internal/resume/evidence.go`,
`internal/resume/resume_test.go`, `internal/board/v2/detail.go`,
`internal/board/v2/detail_view.go`, `internal/board/v2/detail_test.go`,
`internal/board/v2/releases_test.go`, `internal/doctor/checks.go`,
`internal/doctor/checks_test.go`, and `main_resume_matrix_test.go`.

Limitations: the direct `savepoint doctor` command was not run because the
repository CLI rules restrict agent commands to `savepoint resume`; doctor
behavior was verified through `RunV2Checks` tests. The live router remained
selected on T-047, so R-006 with no Objective selection was verified through
the on-disk resume matrix fixtures instead of changing the active selection.

## Drift Notes

None expected; Design is reconciled by the guidance Task.
