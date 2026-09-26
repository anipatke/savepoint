---
id: I-072
title: Doctor reports every unfinished Objective of an in-progress Goal as an error
type: defect
status: open
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

- `internal/doctor/checks.go` around line 83: the loop runs for
  `ColumnInProgress` and `ColumnDone` and appends `releaseBlockerProblem`
  for every blocker except `GateBlockReleaseNoObjectives` on an
  in-progress Goal.
- galaxy: 13 errors, one per unfinished Objective and one per open Task of
  O-001, all for the in-progress Goal R-002.

## Proof Needed

- An in-progress Goal with unfinished Objectives produces no doctor error;
  a done Goal with an unfinished member still does.
- doctor on the migrated galaxy project reports only real problems.
- Tests cover both Goal states; `make build && make test-fast` pass.
