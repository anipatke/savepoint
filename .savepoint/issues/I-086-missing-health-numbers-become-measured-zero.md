---
id: I-086
title: Missing health numbers become measured zero
type: defect
status: resolved
source:
  kind: check
  check: C-939
  actor: {role: checker, session: o027-full-check-20261001}
  at: '2026-09-30T21:34:30Z'
tasks: [T-051]
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

# I-086: Missing health numbers become measured zero

## Summary

O-027's explicit rule that missing measurements never become zero/passing/healthy; T-051 criteria 2, 3, 5 and 6. Frozen M1.

Create a valid tests snapshot containing value={"number":0,"unit":"count"} and Good summary. Remove the number member (or set it to null), preserving the original snapshot ID. DecodeSnapshot accepts the record, restores measured zero, and leaves the Good summary intact.

Evidence: internal/codehealth/snapshot.go:46 defines Value.Number as float64; config.go:153 decodes into structs without checking numeric field presence; snapshot.go:238 validates only the resulting zero. Embedded harness M1_missing_numbers reports "missing/null number accepted as measured 0, summary good, original ID unchanged" for both variants.

Expected: missing/null measured numbers yield a named diagnostic, distinct from an explicit zero. Actual: the decoder silently invents zero and the content identity does not detect the difference. TestMeasuredOutcomeRules checks a nil Value pointer, but not a present value object whose number is missing or null.

## Proof Needed

Add presence-aware decoding/validation for required measured numbers; test omitted/null/explicit-zero through DecodeSnapshot and storage loads, preserving valid zero measurements. Apply the same principle to required numeric detail and threshold fields where applicable within M1.

Repair directly under this Issue and retain completed Task statuses. A fresh independent Check verifies the repair; this record grants no clearance.

