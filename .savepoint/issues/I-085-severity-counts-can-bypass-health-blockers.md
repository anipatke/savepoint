---
id: I-085
title: Severity counts can bypass health blockers
type: defect
status: resolved
source:
  kind: check
  check: C-939
  actor: {role: checker, session: o027-full-check-20261001}
  at: '2026-09-30T21:34:30Z'
tasks: [T-051, T-054]
checks: [C-939, C-940]
guardrail_ids: [DATA-03, TEST-01, TEST-02]
severity: high
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

# I-085: Severity counts can bypass health blockers

## Summary

O-027 classification success condition and confirmed high/critical rule; T-054 criterion 1; T-051 criteria 3 and 5. Frozen cells M1/M7.

A dependency-vulnerability result with total 0, high 1 and critical 0 is accepted by Snapshot.Validate and classified Good after two prior zero official observations. The critical=1 variant also passes as Good. Negative and fractional severity counts also pass snapshot validation.

Evidence: internal/codehealth/classification.go:244 returns before inspecting severity when total is zero; internal/codehealth/snapshot.go:292 checks detail numbers only for finiteness. C-939's embedded harness M7_severity_validation_integration reproduces all four cases. Output for zero_total_high: validation=<nil>, classification=good, explanation="Vulnerabilities 0 meets the good threshold. Unchanged at 0 over 3 official checks."

Expected: a positive high/critical count cannot produce Good; contradictory, negative or fractional severity counts receive a named validation failure. Actual: contradictory evidence is accepted and can bypass a hard blocker. This is reachable through the pure assessment API and a valid persisted snapshot, without bypassing validation. Existing TestAssessCurrentValueThresholds tests consistent totals only.

## Proof Needed

Validate severity counts and their relationship to total, and make the assessment boundary conservative for contradictory evidence. Reproduce total=0/high=1 and total=0/critical=1, negative/fractional counts and the existing missing-severity cases; retain non-overridable blockers under configured thresholds.

Repair directly under this Issue and retain completed Task statuses. A fresh independent Check verifies the repair; this record grants no clearance.

