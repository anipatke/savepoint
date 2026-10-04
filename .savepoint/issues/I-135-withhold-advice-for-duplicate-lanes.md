---
id: I-135
title: Withhold advice for duplicate lanes
type: defect
status: resolved
source:
  kind: check
  check: C-965
  actor: {role: checker, session: check-o033-20261004}
  at: '2026-10-04T03:26:27Z'
tasks: [T-102, T-103]
checks: [C-965, C-966]
resolution:
  disposition: verified
  check: C-966
  actor: {role: checker, session: check-o033-recheck-20261004}
  at: '2026-10-04T03:46:04Z'
  reason: Duplicate-core and missing-title probes now withhold affected advice while ordinary start and unrelated candidates remain available.
history:
  - at: '2026-10-04T03:26:27Z'
    actor: {role: checker, session: check-o033-20261004}
    kind: observed
    check: C-965
    note: Initial independent Full Objective Check.
  - at: '2026-10-04T03:46:04Z'
    actor: {role: checker, session: check-o033-recheck-20261004}
    kind: rechecked
    check: C-966
    note: Duplicate-core and missing-title probes now withhold affected advice while ordinary start and unrelated candidates remain available.
---

# I-135: Withhold advice for duplicate lanes

## Summary

Withhold advice for duplicate lanes before O-033 clearance.

## Evidence

O-033 SC5 and Confirmed Design Optional plan records, T-102 DW1/2, and T-103 DW1 require malformed/ambiguous advice to diagnose and suppress affected recommendations while ordinary actions remain available.

Decode O-001 with lanes [{key: core, title: First}, {key: core, title: Second}, {key: board, title: Board}]. Give T-001/core reads [] writes [a.go] and T-002/board reads [] writes [b.go]. Validate/index and enable advice. Expected: duplicate-core diagnostic and no recommendation involving T-001. Actual: diagnostic exists, but Groups is [[T-001,T-002]]. Ordinary start still succeeds (correct).

`internal/data/concurrency_plan_v2.go:304–308` keeps the first duplicate declaration; validateV2Planning only marks Task diagnostics/unresolved references unusable, so core remains usable. Independent TestO033MalformedLaneSuppressesAdvice/0 fails; missing-title sibling correctly withholds. Existing duplicate-lane decoder test checks diagnostics, not recommendation suppression.

## Proof Needed

Suppress recommendations involving the ambiguous lane without making the project fatal or adding a lifecycle gate. Replay duplicate-core/board and missing-title cells, plus unique lanes, unrelated valid lanes, old/mixed records, malformed Task advice and off-mode membership. Preserve diagnostics and ordinary start decisions. Frozen M1/M3.

## Independent Recheck — C-966

Duplicate-core and missing-title probes now withhold affected advice while ordinary start and unrelated candidates remain available.

Verified within the original C-965 scope by CLEAR Check C-966 on 19a6807. Original evidence remains above.
