---
id: I-050
title: Doctor is silent when no live Goal exists and the router names one
type: defect
status: resolved
source:
  kind: check
  check: C-923
  actor: {role: checker, session: o022-objective-check-20260925}
  at: '2026-09-25T03:27:22Z'
tasks: [T-036, T-037, T-040]
checks: [C-923, C-924]
resolution:
  disposition: verified
  check: C-924
  actor: {role: checker, session: o022-independent-recheck-20260925}
  at: '2026-09-25T04:53:45Z'
guardrail_ids: [TPL-02]
severity: medium
history:
  - at: '2026-09-25T03:27:22Z'
    actor: {role: checker, session: o022-objective-check-20260925}
    kind: observed
    check: C-923
    note: With zero Goals and the router still naming R-001, the board says to run doctor and doctor reports ALL CLEAN.
  - at: '2026-09-25T03:55:05Z'
    actor: {role: executor, session: i050-repair-20260925}
    kind: repair_attempted
    note: >-
      Added shared V2Index.HasLiveGoal detection and used it in doctor and the
      board. Doctor now reports missing, unknown, and archived router Goals and
      recommends creating a Goal when none are live. Added regression cases
      TestRunV2ChecksReportsNoLiveGoalForUnknownAndArchivedSelections,
      TestBoardWithZeroGoalsPointsToDoctor, and
      TestBoardWithOnlyArchivedGoalsPointsToDoctor. Tests were not run in this
      repair session; the reported prior full run predates these changes.
      git diff --check passed. Files changed: internal/data/project.go,
      internal/doctor/v2_runtime.go, internal/doctor/v2_runtime_test.go,
      internal/board/v2/next_panel.go, and
      internal/board/v2/releases_test.go.
  - at: '2026-09-25T04:29:53Z'
    actor: {role: executor, session: o022-repair-followup-20260925}
    kind: repair_attempted
    note: >-
      Focused regressions now pass for unknown and archived router selections
      with no live Goal, including doctor and both board renderers. The same
      focused run also passed the I-051 scaffold and migration checks. Command:
      env GOCACHE=/tmp/savepoint-o022-gocache go test ./internal/doctor
      ./internal/board/v2 ./internal/init ./internal/migrate -run
      '^(TestRunV2ChecksReportsNoLiveGoalForUnknownAndArchivedSelections|TestCheckReleaseReadiness_ignoresEmptyActiveGoalAndKeepsCanonicalFindings|TestBoardWithZeroGoalsPointsToDoctor|TestBoardWithOnlyArchivedGoalsPointsToDoctor|TestV2ScaffoldLoadsCleanThroughRuntimeSchemaAndIndex|TestPlan_routerSelectsLiveGoalOrPlansContinuationGoal|TestPlan_auditedReleaseWithActiveEpicShowsDecisionInPreview|TestEndToEnd_routerFallbackCasesMigrateToLiveGoal|TestEndToEnd_goldenRouterFallbackCases|TestEndToEnd_goldenRouterFallbackCasesAreReproducible)$'
      -count=1. git diff --check passed. Issue remains open for independent
      verification.
  - at: '2026-09-25T04:53:45Z'
    actor: {role: checker, session: o022-independent-recheck-20260925}
    kind: rechecked
    check: C-924
    note: >-
      C-924 independently verified that board and doctor share the live-Goal
      predicate and that unknown or archived router selections with no live
      Goal tell the owner to create one. The board pointer reaches that repair.
---

# I-050: Doctor is silent when no live Goal exists and the router names one

## Summary

O-022 promises that with zero Goals the board points to doctor, and T-036
requires doctor's repair to say to create a Goal when none are declared.
AGENTS.md and Design say "If there are no live Goals, `savepoint doctor` says
to create one." Doctor only does this when the router Goal is blank or `none`.
If the router still names a Goal that is missing or archived, doctor says
nothing about Goals. The board sends the owner to doctor, and doctor reports
all clean.

## Evidence

- Run `savepoint init` in an empty folder, then delete
  `.savepoint/releases/R-001-first-goal/`. The router still reads
  `release: R-001`.
- `savepoint board` (non-TTY) prints `No Goals exist; run savepoint doctor.`
- `savepoint doctor` prints `result: ALL CLEAN (exit code 0)` with no Goal
  problem or repair.
- Replace R-001 with an archived (`legacy_completion`) Release and keep the router
  on it. Resume reports the archived selection. Doctor reports no Goal problem.
  The board shows no doctor pointer, because `boardNextLines` checks only
  `len(Index.Releases) == 0`.
- `internal/doctor/v2_runtime.go:62-78` handles only `SelectionReleaseMissing`
  and `SelectionDone`. `SelectionReleaseNotFound` and `SelectionReleaseArchived`
  get no doctor problem. `hasLiveGoal` is used only inside the missing-router
  repair text.
- `internal/board/v2/next_panel.go:38-43` uses "zero Release records". Doctor uses
  "no live Goal". The two surfaces define "no Goals" differently.
- `TestRunV2ChecksReportsMissingGoalsAndConcreteRepairs` covers only the
  blank-router case with zero Goals.

## Proof Needed

- When no live Goal exists, doctor names the problem and says to create a Goal,
  whatever the router names (blank, unknown, or archived).
- The board and doctor use the same definition of "no Goals" so the board's
  pointer never leads to a clean doctor report.
- Regression tests cover zero Goals with an unknown router Goal and
  archived-only Goals with the router on an archived Goal.
