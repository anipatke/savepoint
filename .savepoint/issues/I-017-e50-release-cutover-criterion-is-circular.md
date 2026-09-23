---
id: I-017
title: E50 cannot satisfy its Release-cutover criterion before its own Check closes
type: verification
status: resolved
source:
  kind: check
  check: C-905
  actor: {role: checker, session: e50-objective-check-20260921}
  at: '2026-09-21T08:06:07Z'
tasks: [T-002]
checks: [C-905]
severity: high
resolution:
  disposition: accepted
  actor: {role: owner, session: owner-accept-20260921}
  at: '2026-09-21T19:15:00Z'
  reason: >-
    Owner accepts the reworded T-002 criterion rather than gating cutover on a
    fresh Check re-verifying it against the gate resolvers: low residual
    risk, not worth another verification loop. Not a technical CLEAR — the
    reworded criterion (T-002 Context Log, 2026-09-21) and the executor's own
    resolver read (ResolveObjectiveCompletion does not depend on Release
    state) are not independently re-proven by this decision.
history:
  - at: '2026-09-21T08:06:07Z'
    actor: {role: checker, session: e50-objective-check-20260921}
    kind: observed
    note: E50 requires R-006 to pass before E50 can clear, while R-006 requires E50 to clear and close first.
  - at: '2026-09-21T18:47:00Z'
    actor: {role: executor, session: remediation-20260921}
    kind: repair_attempted
    note: >-
      Confirmed against the resolvers, not just the report: `ResolveObjectiveCompletion`
      (`internal/data/objective_gate_v2.go`) never reads Release state — O-001's own
      completion depends only on its owned Tasks being done and its own Check clearance
      being current. The circularity was entirely in how T-002's acceptance criterion was
      worded, asking this Check to verify a fact (`R-006 passes ResolveReleaseCutover`) that
      structurally cannot be true until after O-001 closes, which cannot happen until after
      this same Check. Reworded the criterion in T-002 (Context Log entry 2026-09-21) to
      state what the Check can actually prove now — no O-001-owned blocker to R-006's
      cutover — and to explicitly defer the Release-level pass and owner acceptance to
      R-006's own mandatory, strictly-downstream Release Check, which remains required and
      is not shortcut by this change. Not closing this Issue — a fresh Check should
      independently verify the reworded criterion is no longer circular and that no gate
      resolver was weakened.
  - at: '2026-09-21T19:15:00Z'
    actor: {role: owner, session: owner-accept-20260921}
    kind: owner_decision
    note: >-
      Owner accepted the risk instead of commissioning a fresh Check to
      re-verify the reworded criterion. See resolution.reason.
---

# I-017: E50 cannot satisfy its Release-cutover criterion before its own Check closes

## Summary

E50/T-002 requires every declared Release to pass `ResolveReleaseCutover`. R-006 owns O-001 (E50), so its canonical completion decision cannot pass until O-001 has a current CLEAR Objective Check and is owner-closed. R-006 then independently requires a current Release Check and exact owner acceptance. The Objective Check is therefore asked to prove a condition that can only become true after that same Check and later owner/Release actions.

## Evidence

- Violated criterion: T-002 acceptance criterion at line 76 and O-001 architectural delta at line 46.
- `ResolveObjectiveCompletion` requires current Objective clearance (`internal/data/objective_gate_v2.go:42-65`).
- `ResolveReleaseCompletion` first composes every member Objective's completion (`internal/data/release_gate_v2.go:74-98`), then requires a current Release Check and exact owner acceptance (`internal/data/release_gate_v2.go:100-151`).
- Smallest reproduction: run O-001's mandatory Full Objective Check while R-006 is `in_progress`, O-001 is not yet owner-closed, and no R-006 Release Check can yet be current. The criterion remains false regardless of the Objective implementation result.
- Expected: every O-001 acceptance criterion is independently provable during its mandatory Objective Check.
- Actual: the Release-cutover criterion is only satisfiable after O-001's Check and closure, so this Check cannot return CLEAR without waiving or redefining its own acceptance boundary.
- Missing evidence: no staged criterion or non-circular ordering contract defines Objective clearance first and Release cutover acceptance afterward.

## Proof Needed

Replan the criterion and gate ordering so the Objective Check can prove an O-001-owned outcome without requiring its downstream Release gate to have already completed. Preserve the mandatory later Release Check and exact owner acceptance. A fresh Full Objective Check must verify the revised, non-circular contract.
