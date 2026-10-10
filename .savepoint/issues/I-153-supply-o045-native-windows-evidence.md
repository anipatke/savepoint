---
id: I-153
title: Native Windows evidence is missing for O-045
type: verification
status: open
source:
  kind: check
  check: C-978
  actor: {role: checker, session: check-o045-20261011}
  at: '2026-10-10T21:10:26Z'
tasks: [T-127, T-128, T-129, T-130]
checks: [C-978]
guardrail_ids: [CFG-03, CFG-02]
history:
  - at: '2026-10-10T21:10:26Z'
    actor: {role: checker, session: check-o045-20261011}
    kind: observed
    check: C-978
    note: Initial Full Objective Check of O-045.
---

# I-153: Native Windows evidence is missing for O-045

## Summary

Current full Go suite must pass natively on Windows CI and be named against the reviewed code. No covering CI result exists for this working-tree implementation.

Violated requirements: O-045 SC7; T-129 DW5; T-130 DW7.

## Evidence

- Location: `T-129 and T-130 Technical Evidence / Limitations; internal/run/proc_windows.go`.
- Smallest reproduction: The reviewed changes are uncommitted on 87ee065; both platform-sensitive Task records explicitly state the branch was not pushed and native Windows tests were not run. Local full tests cross-build Windows but cannot execute its process-tree and liveness implementation.
- Independent harness and full output are embedded in C-978. Probe PASS means the observation ran successfully, not that the requirement passed.

## Proof Needed

Supply the native Windows full-suite CI job URL, reviewed commit, and successful result, including run stop/tree/liveness tests. Cross-compilation or historical CI for other code is insufficient.
