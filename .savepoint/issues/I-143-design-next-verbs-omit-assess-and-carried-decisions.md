---
id: I-143
title: Design's Next-verb and acceptance rules omit Assess and carried decisions
type: drift
status: resolved
source:
  kind: check
  check: C-968
  actor: {role: checker, session: check-o044-20261010}
  at: '2026-10-10T06:29:46Z'
tasks: [T-122, T-123]
checks: [C-968, C-969, C-970]
resolution:
  disposition: verified
  check: C-970
  actor: {role: checker, session: check-o044-final-20261010}
  at: '2026-10-10T06:38:20Z'
  reason: 'Design lists Assess, the widened Accept and carried owner decisions; proven in C-969 and unchanged since.'
history:
  - at: '2026-10-10T06:29:46Z'
    actor: {role: checker, session: check-o044-20261010}
    kind: observed
    check: C-968
    note: Found during O-044 Design reconciliation.
  - at: '2026-10-10T06:34:45Z'
    actor: {role: checker, session: check-o044-recheck-20261010}
    kind: rechecked
    check: C-969
    note: 'Re-check — Design lists Assess, widened Accept and carried decisions; repair proven. Awaits a CLEAR Check to record verified.'
  - at: '2026-10-10T06:38:20Z'
    actor: {role: checker, session: check-o044-final-20261010}
    kind: rechecked
    check: C-970
    note: 'Verified by C-970 (CLEAR).'
---

# I-143: Design's Next-verb and acceptance rules omit Assess and carried decisions

## Summary

O-044 added the `Assess` Next verb, made `Accept` also mean "renew a decision
a material change ended", and let owner acceptance and exceptions carry
forward to later Checks through `carried_forward`. `.savepoint/Design.md` was
not reconciled: it still lists the verbs without `Assess`, defines `Accept`
only as "current Check, owner acceptance the only blocker", and nowhere
records the carry-forward fields, the `unmet` Check field or the new decision
blockers. O-044 Success Condition 6 requires the runtime, resume, board,
skills and documentation to state the same rule; the Full Objective Check
includes Design reconciliation.

## Evidence

- `.savepoint/Design.md:198`: "Task: `Start` …, then at audit `Check`,
  `Accept` (current Check, owner acceptance the only blocker), or `Close` …
  Objective: `Check`, `Accept`, `Close` …" and "Check, Accept, and Close in
  the Check-phase colour" — no `Assess`, although `internal/resume/resume.go:125`
  and `:143` return it and `internal/board/v2/next_panel.go:84` colours it.
- `.savepoint/Design.md:31` owns "acceptance, exception" decisions in
  `internal/data` but records no carried-forward applicability, `owner_validation.scope`
  or Check `unmet`.
- `grep -n "carried_forward\|Assess\|unmet" .savepoint/Design.md` returns
  nothing.

## Proof Needed

- Design lists `Assess` for Task and Objective with its meaning
  (applicability only, not a new Check) and the widened meaning of `Accept`.
- Design's evidence boundary records `carried_forward`, `owner_validation.scope`
  and Check `unmet`, and the routing for unassessed, changed and uncovered
  decisions, consistent with the fix for I-142.
- Planner-owned edit (`savepoint-design`); verified by the next O-044 Check.
