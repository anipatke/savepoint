---
id: I-100
title: OSV reader double counts distinct vulnerability groups
type: defect
status: resolved
source:
  kind: check
  check: C-944
  actor: {role: checker, session: o029-full-check-20261001}
  at: '2026-10-01T09:35:25Z'
tasks: [T-064]
checks: [C-944, C-945]
guardrail_ids: [TEST-01, TEST-02]
severity: medium
history:
  - at: '2026-10-01T09:35:25Z'
    actor: {role: checker, session: o029-full-check-20261001}
    kind: observed
    check: C-944
    note: Initial Full O-029 Check; see C-944 frozen matrix and embedded independent harness.
  - at: '2026-10-01T10:30:00Z'
    actor: {role: executor, session: o029-issue-repair-20261001}
    kind: repair_attempted
    note: 'OSV reader counts each distinct advisory-ID group once per package (order-insensitive); distinct groups and the same group in different packages stay separate. Tests cover identical, reordered, different, and cross-package cases, with severity totals matching. make build && make test-fast pass. Awaiting independent recheck.'
  - at: '2026-10-01T10:42:32Z'
    actor: {role: checker, session: o029-independent-recheck-20261001}
    kind: rechecked
    check: C-945
    note: Original reproduction and admitted adjacent frozen cells pass; see C-945 for independent evidence and successful fresh full gate.
resolution:
  disposition: verified
  check: C-945
  actor: {role: checker, session: o029-independent-recheck-20261001}
  at: '2026-10-01T10:42:32Z'
  reason: Independent Full Objective recheck proved the repair within C-944 frozen scope.
---

# I-100: OSV reader double counts distinct vulnerability groups

## Summary

T-064 Done When 1 requires the number of distinct vulnerability groups, one per package; O-029 requires trustworthy vulnerability totals. Frozen M5 duplicate identical group cell.

## Evidence

C-944 `M5 OSV duplicate group`: `{"results":[{"source":{"path":"go.mod"},"packages":[{"package":{"name":"x","version":"1"},"groups":[{"ids":["X"],"max_severity":"7"},{"ids":["X"],"max_severity":"7"}]}]}]}`.

Expected one distinct group for package x@1 (or a named invalid-duplicate report diagnostic). Actual value=2, high=2, two identical affected items, complete reading. reader_osv.go:110-113 counts each array element with no distinct-group identity check. TestOSVScannerReaderValue checks totals against its fixture array length and does not test duplicate groups. Inflation distorts totals/severity/trends; it does not bypass the high-severity blocker.

## Proof Needed

Deduplicate equivalent advisory-ID groups within each package, or reject contradictory duplicates with a named error; keep different packages distinct. Prove identical/reordered alias group variants and distinct packages using the original M5 duplicate axis; severity totals must agree with the distinct total.

Repair under this Issue by default and preserve every completed Task status. This record grants no clearance; a later independent Check verifies remediation.
