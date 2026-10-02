---
id: I-115
title: Native Windows release evidence is missing
type: verification
status: resolved
source:
  kind: check
  check: C-957
  actor: {role: checker, session: check-o032-20261002-independent}
  at: '2026-10-02T05:23:00Z'
tasks: [T-091, T-095]
checks: [C-957, C-958, C-960]
guardrail_ids: [CFG-03, CFG-02, TEST-08]
severity: medium
resolution:
  disposition: verified
  check: C-960
  actor: {role: checker, session: recheck-o032-20261003-independent}
  at: '2026-10-02T19:20:00Z'
  reason: "Current-head native evidence: run 36975366853 on exact 37dffd3 has ci and windows-tests successful; fresh Linux make ci passed."
history:
  - at: '2026-10-02T05:23:00Z'
    actor: {role: checker, session: check-o032-20261002-independent}
    kind: observed
    check: C-957
    note: Initial independent Full O-032 Check; frozen scope and reproducible harness in C-957.
  - at: '2026-10-02T06:19:18Z'
    actor: {role: checker, session: recheck-o032-20261002-independent}
    kind: rechecked
    check: C-958
    note: "Native current-head Windows full suite and Linux make ci success verified from https://github.com/anipatke/savepoint/actions/runs/36971878006 on d00543d5252fc5557aeba067995427236da2823f, go1.26.8. Evidence gap filled; remains open pending CLEAR Check proof."
  - at: '2026-10-02T19:20:00Z'
    actor: {role: checker, session: recheck-o032-20261003-independent}
    kind: rechecked
    check: C-960
    note: "CLEAR. Current-head native evidence: run 36975366853 on exact 37dffd3 has ci and windows-tests successful; fresh Linux make ci passed."
---

# I-115: Native Windows release evidence is missing

## Summary

O-032's confirmed observable outcome 3 requires native Windows CI to pass. T-091 Done When 5 and T-095 Done When 4 explicitly carry this evidence to the Full Objective Check. Current Linux full tests and cross-compilation are successful, but native runtime verification of this work is **Unverified**, not a claimed Windows product failure.

## Evidence

Read-only GitHub queries on 2026-10-02: `gh run list --branch v2.1 --limit 5 --json databaseId,headSha,status,conclusion,workflowName,url` returns `[]`. The recent green CI is https://github.com/anipatke/savepoint/actions/runs/36296646883 on master head 9beffa5b01aa4e92425498154f5055a8ca642776, which predates this work. `.github/workflows/ci.yml:5-8` triggers only push/PR to master and v2; its `windows-tests` job at :31 runs the full native suite. Reviewed head is 136fe69 plus the uncommitted history-window and documentation work listed in C-957. T-095 accurately records this missing evidence. No push, PR, workflow dispatch or CI change was authorized or performed.

## Proof Needed

Provide a successful native `windows-tests` full-suite result for the final committed tree including O-032's uncommitted work. Record its commit/run URL and toolchain. Cross-builds and older green runs do not substitute. Keep completed Tasks done; this Issue records an evidence gap.

Repair directly under this Issue by default; preserve all owner-completed Task statuses. This record grants no clearance or owner acceptance.
