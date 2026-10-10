---
id: I-142
title: An exception that misses an unmet requirement routes to another Check instead of the owner
type: defect
status: resolved
source:
  kind: check
  check: C-968
  actor: {role: checker, session: check-o044-20261010}
  at: '2026-10-10T06:29:46Z'
tasks: [T-121, T-122]
checks: [C-968, C-969, C-970]
guardrail_ids: [TPL-02]
resolution:
  disposition: verified
  check: C-970
  actor: {role: checker, session: check-o044-final-20261010}
  at: '2026-10-10T06:38:20Z'
  reason: 'Uncovered unmet requirement routes to Accept for Objective and Task, names the IDs, and asks the owner to widen or renew the exception or send the work back; full gate passed.'
history:
  - at: '2026-10-10T06:29:46Z'
    actor: {role: checker, session: check-o044-20261010}
    kind: observed
    check: C-968
    note: Found by the O-044 Full Objective Check probe P4; exception_scope blocker has no Next rung.
  - at: '2026-10-10T06:34:45Z'
    actor: {role: checker, session: check-o044-recheck-20261010}
    kind: rechecked
    check: C-969
    note: 'Re-check — route now Accept for Objective and Task and evidence names uncovered IDs; action line still says accept the current Check (NEEDS WORK). Still open.'
  - at: '2026-10-10T06:38:20Z'
    actor: {role: checker, session: check-o044-final-20261010}
    kind: rechecked
    check: C-970
    note: 'Verified by C-970 (CLEAR).'
---

# I-142: An exception that misses an unmet requirement routes to another Check instead of the owner

## Summary

When the latest Check lists an `unmet` requirement that an otherwise
applicable exception does not cover, the gate correctly refuses completion
with a `GateBlockExceptionScope` blocker naming the uncovered IDs. But Next
has no rung for that blocker, so resume and the board fall back to
`Check O-…` / "Record the Objective … integration Check." and never show the
uncovered IDs. This is the loop O-044 set out to remove, and it contradicts
the Confirmed Design ("an uncovered ID is a named blocker"; materially
changed → `Accept`, for that decision only, naming the change), O-044 Success
Condition 2, T-122 Done When 1–2 ("Accept lines name … the material change or
uncovered requirement IDs") and the new README text ("you are asked again
only for the decision that changed").

## Evidence

Independent probe (C-968, P4), TheShed-shaped temporary project: Objective
O-002 with two done Tasks, owner acceptance and an exception
(`requirements: [TEST-08, CFG-03]`) recorded at C-005, both carried to C-006
by a checker entry `applies: true`; C-006 is NEEDS WORK with
`unmet: [TEST-08, DESIGN-02]`.

- Expected: `Accept O-002`, naming the exception and the uncovered
  `DESIGN-02` (owner widens/renews the exception, or the work is fixed).
- Actual:

  ```
  Check O-002 — Shed
  Technical clearance: Check C-006 recorded NEEDS WORK.
  Next action: Record the Objective O-002 integration Check.
  ```

Code: `internal/data/decision_gate_v2.go:173-183` emits
`GateBlockExceptionScope`; `internal/data/next.go:581-593` (`decisionRung`)
maps only `GateBlockDecisionUnassessed` and `GateBlockDecisionChanged`, so
`resolveObjectiveIntegrationRung` (`next.go:651`) falls through to
`NextObjectiveIntegration`, and Task routing (`next.go:560`) falls through to
the clearance rung. `internal/resume/resume.go:288-300` has no evidence line
for the blocker. The same probe with `unmet: [TEST-08]` correctly reads
`Close O-002 … Ready to close by exception`.

Missing test evidence: `TestResolveObjectiveCompletion_exceptionMustCoverEveryUnmetRequirement`
proves only the gate blocker; no Next, resume or board test covers the
uncovered-requirement route.

## Proof Needed

- The P4 shape (Objective and Task variants) routes to `Accept` (or another
  owner-facing verb the planner confirms) — never `Check` — and its evidence
  line names the exception, the latest Check and every uncovered ID, with
  authored text sanitised.
- Board and resume read the same projection for it; the board refusal text
  for `exception_scope` stays accurate.
- `unmet` fully covered still reads `Close … Ready to close by exception`;
  unassessed still reads `Assess`.
- `make build && make test-fast` pass; the next Full Objective Check runs
  `make test-full`.
