---
id: I-090
title: Repository identity ignores POSIX colon filenames
type: defect
status: resolved
source:
  kind: check
  check: C-939
  actor: {role: checker, session: o027-full-check-20261001}
  at: '2026-09-30T21:34:30Z'
tasks: [T-053]
checks: [C-939, C-940]
guardrail_ids: [CFG-02, TEST-01, TEST-02]
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

# I-090: Repository identity ignores POSIX colon filenames

## Summary

T-053 criteria 1/2/4 promise fingerprints that change with relevant content and deterministic unusual-filename handling or named errors. Frozen M5.

On Linux, commit a regular source:part.go file, observe with empty InputScope (all inputs), edit its content and observe again. Both observations say Dirty:false and have identical fingerprints despite Git's tracked change.

Evidence: internal/codehealth/repository.go:321 rejects every path containing a colon; paths at line 309 silently omits it from both input and changed-file listings. Embedded M5_colon_filename reproduces this on a temporary repository. Existing TestObserveUnusualFilenames covers spaces, Unicode, quotes, tabs and newline but not colon.

Expected: supported POSIX filenames are fingerprinted, or an explicit named error prevents claiming clean/unchanged state. Actual: relevant source content disappears from repository identity. Windows drive/stream path protection should remain intact.

## Proof Needed

Use platform-correct confinement checks or report rejected Git paths explicitly. Recheck tracked and relevant untracked POSIX colon paths and ordinary drive/traversal rejection, without silently losing inputs.

Repair directly under this Issue and retain completed Task statuses. A fresh independent Check verifies the repair; this record grants no clearance.

