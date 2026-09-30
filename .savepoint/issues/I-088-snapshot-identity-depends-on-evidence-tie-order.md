---
id: I-088
title: Snapshot identity depends on evidence tie order
type: defect
status: resolved
source:
  kind: check
  check: C-939
  actor: {role: checker, session: o027-full-check-20261001}
  at: '2026-09-30T21:34:30Z'
tasks: [T-051]
checks: [C-939, C-940]
guardrail_ids: [TEST-01, TEST-02]
severity: medium
resolution:
  disposition: verified
  check: C-940
  actor: {role: checker, session: o027-recheck-20261001}
  at: '2026-09-30T22:21:36Z'
  reason: Repair reproduced and proven by CLEAR Full Objective re-check C-940.
history:
  - at: '2026-09-30T21:34:30Z'
    actor: {role: checker, session: o027-full-check-20261001}
    kind: observed
    check: C-939
    note: Reproduced during the initial Full O-027 Check; see frozen matrix and embedded harness in C-939.
  - at: '2026-09-30T22:21:36Z'
    actor: {role: checker, session: o027-recheck-20261001}
    kind: rechecked
    check: C-940
    note: C-939 reproduction and independent re-check probes pass; see C-940 closure map.
---

# I-088: Snapshot identity depends on evidence tie order

## Summary

T-051 criterion 4 promises stable canonical snapshot identity. Frozen M2.

Use two valid EvidenceRef entries with the same path report.json and line 1, but notes first and second. Seal and validate the snapshot. Reverse only those entries: ComputeID changes from sha256:fae282b44beaeec064dff3f192cde7f034ad89a52fe9dbb86a230ba115f70485 to sha256:8785ebbce8f94c0aa065da4f5a2925c4ecb7f5a9ada5ccfc7d9f7ece7e235cb4.

Evidence: internal/codehealth/identity.go:51 compares evidence Path and Line only, ignoring Note when those keys tie. Snapshot validation permits these distinct references. Embedded harness M2_evidence_tie_permutation fails.

Expected: reordering an order-insensitive list preserves identity and an identical retry. Actual: one logical observation gets different identities or fails its retained ID validation after reordering. TestSnapshotIDIgnoresListOrder covers results, summaries and details, not evidence ties.

## Proof Needed

Define a total canonical evidence ordering (including Note), or explicitly reject ambiguous duplicate path/line entries; verify permutations, valid distinct notes, same-ID validation and save retry.

Repair directly under this Issue and retain completed Task statuses. A fresh independent Check verifies the repair; this record grants no clearance.

