---
id: I-072
title: Doctor reports every unfinished Objective of an in-progress Goal as an error
type: defect
status: resolved
source:
  kind: report
  actor: {role: owner, session: user}
  at: '2026-09-26T03:34:20Z'
severity: medium
history:
  - at: '2026-09-26T03:34:20Z'
    actor: {role: owner, session: user}
    kind: observed
    note: >-
      galaxy after migration: doctor exits 1 with a
      [v2-release-objective-incomplete] error for each of O-001..O-007 in
      R-002, which is in_progress, including one per open Task of O-001.
  - at: '2026-09-26T04:18:18Z'
    actor: {role: executor, session: fix-i072-doctor-active-goal}
    kind: repair_attempted
    note: >-
      Updated Doctor so active Goals skip empty-membership and unfinished
      member readiness blockers while done Goals still report them. Added
      coverage for an active Goal with an open Task, a done Goal with an
      incomplete Objective, and the migrated active-Goal fixture. The first
      fast-gate run exposed a stale integration-test expectation for an active
      migrated Goal; after updating it, `make build && make test-fast` passed.
      Read the Issue, Task and Issue-capture guidance, router, STYLE policy,
      Doctor implementation and tests, and the affected migration parity test.
      Changed internal/doctor/checks.go, internal/doctor/checks_test.go,
      main_board_next_parity_test.go, and this Issue record. No Check was
      written and the router was not edited. The reported galaxy project was
      unavailable in this worktree, so its direct Doctor command remains
      unverified.
  - at: '2026-09-26T04:41:01Z'
    actor: {role: owner, session: user}
    kind: owner_decision
    note: >-
      Owner instructed the executor to mark this Issue resolved after the
      repair merged to v2 with make build, make test-fast, and make test-full
      passing on the merged branch. No independent Check was run and no
      technical CLEAR is implied.
resolution:
  disposition: accepted
  actor: {role: owner, session: user}
  at: '2026-09-26T04:41:01Z'
  reason: Owner accepted the repair merged to v2 in 7adeff8.
---

# I-072: Doctor reports every unfinished Objective of an in-progress Goal as an error

## Summary

The doctor release check resolves `ResolveReleaseCompletion` for Goals that
are `in_progress` as well as `done`, and reports every blocker as an error
(it already skips `release-no-objectives` for in-progress Goals). An active
Goal's Objectives are unfinished by definition, so an ordinary project in
active work never reports clean, and the real problems are buried. The
check should flag incomplete members only for a Goal recorded `done`.

## Evidence

- Before repair, `internal/doctor/checks.go` around line 83: the loop ran for
  `ColumnInProgress` and `ColumnDone` and appended `releaseBlockerProblem`
  for every blocker except `GateBlockReleaseNoObjectives` on an
  in-progress Goal.
- galaxy: 13 errors, one per unfinished Objective and one per open Task of
  O-001, all for the in-progress Goal R-002.
- Repair attempted 2026-09-26: active Goal readiness no longer reports
  incomplete member blockers; done Goal blockers remain. The updated
  regression coverage and `make build && make test-fast` passed. The migrated
  galaxy project was not available for a direct Doctor run.

## Proof Needed

- An in-progress Goal with unfinished Objectives produces no doctor error;
  a done Goal with an unfinished member still does.
- doctor on the migrated galaxy project reports only real problems.
- Tests cover both Goal states; `make build && make test-fast` pass.
