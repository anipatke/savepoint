---
id: C-907
scope: {kind: objective, id: O-012}
result: CLEAR
checked_by: {role: checker, session: o012-recheck-20260922}
executed_session: i018-remediation-20260922
checked_at: '2026-09-22T10:00:27Z'
reviewed:
  base_commit: 84e9e502985aeb9c9602f3e08c8dbddd21233e4c
  head_commit: 84e9e502985aeb9c9602f3e08c8dbddd21233e4c
  files:
    - .savepoint/objectives/O-900-test-objective-for-ui-checks/tasks/T-903.md
    - .savepoint/objectives/O-900-test-objective-for-ui-checks/tasks/T-909.md
    - internal/board/v2/card_test.go
    - internal/board/v2/fixture_test.go
    - internal/board/v2/run_test.go
    - .savepoint/objectives/O-012-simplify-task-card-outcomes/tasks/T-005-refresh-the-o900-visual-test-spread.md
    - .savepoint/objectives/O-012-simplify-task-card-outcomes/tasks/T-009-restore-the-o900-pending-check-example.md
  dependencies: []
issues: []
supersedes: C-906
---

# C-907: O-012 Full Objective Re-check

## Verdict

`CLEAR`. I-018 is repaired. O-900 now renders all six retained Task-card review
outcomes through interactive and non-TTY paths while preserving its 4/4/4
column spread and all wait, replan, and owner blockers. Production Task-card
behavior is unchanged.

## Closure Map

| Prior Issue | Result | Evidence |
| --- | --- | --- |
| I-018 | Closed — verified | T-903 is the Done/no-Check `[ ] CHECK` example; T-909 carries the in-progress stale `[!] REVIEW` role; durable card and non-TTY matrix tests pass |

## Frozen-Scope Re-check

The C-906 scope lock was reused without expansion. Every original outcome,
blocker, column-spread, interactive/plain parity, retired-label, policy, and
documentation cell remains proven. The failed O-900 pending-outcome cell now
passes on both rendering paths.

## Evidence

- `TestO900OutcomeSpreadRendersEveryOutcomeAndBlockerOnCards` — pass.
- `TestRunWithoutTTYRendersTheCompleteO900OutcomeSpread` — pass.
- Focused fresh recheck command — pass in 0.013s.
- `git diff --check` — pass.
- `make build` — pass.
- Executor's post-repair `make test` — pass, including `internal/migrate` in
  118.614s. The owner explicitly directed the checker not to repeat this
  unchanged full gate; an attempted duplicate run was interrupted and is not
  used as evidence.
- Strict live-project loading after the T-009 rename and after the fixture
  repair — pass.

I-019-I-024 are separate follow-up Issues outside the frozen I-018/O-012 repair
scope and do not change this verdict.

## Materiality

No material actions remain inside the C-906 scope lock. I-018's original Medium
materiality is resolved by the live fixture repair and durable regression
coverage.

## Owner Action

O-012 has current mandatory Full Objective clearance. The checker does not mark
the Objective `done`; owner closure remains a separate decision.
