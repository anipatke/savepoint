---
id: I-089
title: Invalid health globs silently empty the scope
type: defect
status: resolved
source:
  kind: check
  check: C-939
  actor: {role: checker, session: o027-full-check-20261001}
  at: '2026-09-30T21:34:30Z'
tasks: [T-051, T-053]
checks: [C-939, C-940]
guardrail_ids: [CFG-01, DATA-03, TEST-02]
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

# I-089: Invalid health globs silently empty the scope

## Summary

T-051 configuration validation, T-053 criteria 1/2/4 and CFG-01. Frozen M1/M5.

InputScope{Include: []string{"["}} and configuration scope ["src/[abc"] validate successfully. ObserveRepository with an unterminated-bracket include returns a clean observation with no error. Changing a.go leaves the observation and fingerprint unchanged because nothing matches.

Evidence: internal/codehealth/primitives.go:120 validates relative path spelling but not glob syntax; internal/codehealth/repository.go:163 discards path.Match errors as a false match. Embedded M5_malformed_pattern records nil errors, clean=true, unchanged=true. Existing invalid-scope tests cover traversal, not malformed glob syntax.

Expected: a malformed include/exclusion returns a named configuration error. Actual: a typo silently suppresses relevant inputs and can misrepresent stale evidence as current.

## Proof Needed

Validate the supported glob grammar at configuration/InputScope boundaries and propagate invalid-pattern errors instead of treating them as nonmatches. Cover unterminated bracket patterns in includes and exclusions plus valid directory, *, **, and bracket patterns.

Repair directly under this Issue and retain completed Task statuses. A fresh independent Check verifies the repair; this record grants no clearance.

