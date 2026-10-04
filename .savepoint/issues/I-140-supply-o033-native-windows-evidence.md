---
id: I-140
title: Supply native Windows evidence for O-033
type: verification
status: open
source:
  kind: check
  check: C-965
  actor: {role: checker, session: check-o033-20261004}
  at: '2026-10-04T03:26:27Z'
tasks: [T-102, T-103, T-104, T-105, T-106, T-107]
checks: [C-965]
guardrail_ids: [CFG-03]
history:
  - at: '2026-10-04T03:26:27Z'
    actor: {role: checker, session: check-o033-20261004}
    kind: observed
    check: C-965
    note: Initial independent Full Objective Check.
---

# I-140: Supply native Windows evidence for O-033

## Summary

Supply native Windows evidence for O-033 before O-033 clearance.

## Evidence

O-033 Skills, verification and delivery and T-107 Technical Verification explicitly require native windows-tests CI evidence for the final implementation. CFG-03 cannot be satisfied by cross-compilation.

Reviewed HEAD is 53fbe53405cdfc99636f57e830c0ed08a6f38d9a plus uncommitted T-105–107 guidance/docs/integration tests. Fresh local make build && make test-full passed all host tests/cross-builds. Independent gh run list --branch v2.20 shows only run 37171504754 on older 4520c93 (O-037), before O-033 changes. No current native run has been supplied. The owner asked why the evidence was needed; no waiver or walkthrough acceptance was inferred.

This is missing platform proof, not an observed Windows failure. Existing I-134 was for O-037 and is already resolved; it does not cover new O-033 source/test inputs.

## Proof Needed

After the narrow repairs, supply successful repository native windows-tests CI evidence tied to the final O-033 code and test inputs (URL/revision/job result). Fresh full host gate remains required after repair. Recheck frozen M8, no extra CI job or orchestration requirement.
