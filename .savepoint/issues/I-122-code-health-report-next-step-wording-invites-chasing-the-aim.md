---
id: I-122
title: Code Health report next-step wording invites chasing the aim
type: drift
status: open
source:
  kind: report
  actor: {role: executor, session: build-t100-20261003}
  at: '2026-10-03T12:00:00Z'
tasks: [T-099, T-100]
guardrail_ids: [STYLE-05]
severity: low
history:
  - at: '2026-10-03T12:00:00Z'
    actor: {role: executor, session: build-t100-20261003}
    kind: observed
    note: Follow-up F1 from the T-099 skills review. Not applied in T-100 because report wording is product work.
---

# I-122: Code Health report next-step wording invites chasing the aim

## Summary

The generated Code Health report tells an agent to "investigate each signal ... propose a fix, then apply it", and the complexity line says "aim for 10 or less". An agent following the report chases the aim rather than stopping at the watch line (20 or less). `savepoint-task` and `savepoint-check` now state the stopping point, but the report an agent reads first still does not.

## Evidence

- `internal/codehealth/report.go` and `internal/codehealth/dashboard_copy.go` carry the next-step and complexity wording.
- O-032 Lessons Carried In: bringing the hardest function from 42 to 20 took about 25 functions; reaching the aim would touch about 280 more, over 120 of them production code.
- C-957 and C-960 observations: most complexity improvement is production refactoring, and the Watch band stayed populated at close.

## Proof Needed

Decide the report wording (name the watch line as the stopping point, say what to report when signals remain), change the generated text and its tests, and show a regenerated `.savepoint/health/report.md` carrying it. Report wording and health classification are product behaviour, so the owner decides the text.
