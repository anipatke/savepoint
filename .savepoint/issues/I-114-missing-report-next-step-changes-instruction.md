---
id: I-114
title: Missing-report next step changes the required agent instruction
type: defect
status: resolved
source:
  kind: check
  check: C-955
  actor: {role: checker, session: check-o036-20261002-independent}
  at: '2026-10-02T03:01:24Z'
tasks: [T-084]
checks: [C-955, C-956]
resolution:
  disposition: verified
  check: C-956
  actor: {role: checker, session: recheck-o036-20261002-independent}
  at: '2026-10-02T03:12:42Z'
  reason: Repair independently verified against the frozen C-955 scope.
history:
  - at: '2026-10-02T03:12:42Z'
    actor: {role: checker, session: recheck-o036-20261002-independent}
    kind: rechecked
    note: 'Persisted missing-report next step matches the required investigate-it text; full prompt and original signal rows are visible at 80x20/24/40.'
    check: C-956
---

# I-114: Missing-report next step changes the required agent instruction

## Summary

T-084 Done When 1 explicitly requires the missing-report next step to be "Run savepoint health report, then ask your agent to investigate it." The implementation says "fix it." This alters the requested instruction and makes the referent ambiguous: the report is derived, and the agent should investigate its findings rather than fix the report itself. This is a small copy correction, not a change in report behavior.

## Evidence

`internal/codehealth/dashboard_copy.go:104` defines "Run savepoint health report, then ask your agent to fix it." LoadDashboard uses that literal for a Needs Attention row when no report exists. Independent persisted fixture: save three official snapshots with complexity 46 and no report, load the dashboard, call PopoverNextStep on complexity. Actual is the fix-it literal; expected is the investigate-it literal above.

TestRedSignalPointsAtTheReport (`internal/codehealth/dashboard_copy_test.go:388`) and TestHealthPopoverNextStepPointsAtTheReportAtEverySize (`internal/board/v2/health_realcopy_test.go:188`) explicitly assert the changed wording, so green tests do not prove the named acceptance text. Present-report wording and frame size checks pass.

## Proof Needed

Restore the required missing-report instruction, update its assertions, and prove the full text remains visible at 80x20, 80x24 and 80x40. Recheck the frozen red/report-missing cells and full gate. Alternatively the owner/planner must explicitly reconcile the acceptance wording; a checker does not rewrite it.

## Recheck Proof

C-956: Persisted missing-report next step matches the required investigate-it text; full prompt and original signal rows are visible at 80x20/24/40. Fresh full and focused gates pass.
