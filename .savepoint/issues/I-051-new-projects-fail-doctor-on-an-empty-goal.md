---
id: I-051
title: New and migrated projects fail doctor because their Goal has no Objectives
type: defect
status: resolved
source:
  kind: check
  check: C-923
  actor: {role: checker, session: o022-objective-check-20260925}
  at: '2026-09-25T03:27:22Z'
tasks: [T-038, T-039]
checks: [C-923, C-924]
resolution:
  disposition: verified
  check: C-924
  actor: {role: checker, session: o022-independent-recheck-20260925}
  at: '2026-09-25T04:53:45Z'
severity: medium
history:
  - at: '2026-09-25T03:27:22Z'
    actor: {role: checker, session: o022-objective-check-20260925}
    kind: observed
    check: C-923
    note: A fresh savepoint init project reports PROBLEMS FOUND (exit 1) with v2-release-no-objectives on R-001.
  - at: '2026-09-25T03:58:59Z'
    actor: {role: executor, session: i051-repair-20260925}
    kind: repair_attempted
    note: >-
      Doctor no longer reports the no-objectives completion blocker for an
      in-progress Goal, which is allowed to exist before its first Objective;
      the completion gate still rejects closing an empty Goal, and a done Goal
      without members remains a doctor finding. Added doctor coverage to
      TestV2ScaffoldLoadsCleanThroughRuntimeSchemaAndIndex and
      TestEndToEnd_routerFallbackCasesMigrateToLiveGoal, and updated
      TestCheckReleaseReadiness_ignoresEmptyActiveGoalAndKeepsCanonicalFindings.
      Tests were not run in this repair session; the reported prior full run
      predates these changes. git diff --check passed. Files changed:
      internal/doctor/checks.go, internal/doctor/checks_test.go,
      internal/init/v2_scaffold_test.go, and
      internal/migrate/end_to_end_test.go.
  - at: '2026-09-25T04:29:53Z'
    actor: {role: executor, session: o022-repair-followup-20260925}
    kind: repair_attempted
    note: >-
      The focused checks now pass: a fresh scaffold reports no doctor problem,
      the migrated continuation Goal is accepted while empty and in progress,
      and doctor still reports empty completed Goals. Command:
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
      C-924 independently verified that a fresh project and an in-progress
      continuation Goal with no Objectives pass doctor, while completion
      checks still reject an empty completed Goal.
---

# I-051: New and migrated projects fail doctor because their Goal has no Objectives

## Summary

O-022's Outcome says new and migrated projects start valid. T-038's User Check
runs `savepoint doctor` on a fresh project and expects no Goal problem. Now
that init creates R-001, every fresh project fails doctor. The existing
`v2-release-no-objectives` check flags the empty Goal, and doctor exits 1.
Migrate's generated continuation Goal is empty in the same way. At Idea state,
no Objective exists yet, so the repair doctor offers ("Assign at least one
Objective") cannot be acted on.

## Evidence

- `savepoint init` in an empty folder, then `savepoint doctor`:
  `✗ releases/R-001-first-goal/Release.md: [v2-release-no-objectives] release R-001: release R-001 has no member Objectives`
  and `result: PROBLEMS FOUND (exit code 1)`.
- Before O-022 the scaffold had no Goal record, so this finding did not
  occur on a fresh project.
- The materialized `v1-router-missing` and `v1-router-archived` migrate goldens
  (with `schema_version: 2`) report the same finding for the generated R-002.
- No test runs doctor on a freshly initialized or fallback-migrated project.
  `TestV2ScaffoldLoadsCleanThroughRuntimeSchemaAndIndex` checks strict loading
  only.

## Proof Needed

- A freshly initialized project, and a migrated project whose router selects
  a generated Goal, pass doctor with no problem for the starting Goal. The
  alternative is an explicit planner decision recording why this finding is
  expected and what the owner should do.
- Regression tests run doctor on a fresh scaffold and on a fallback-migrated
  project.
