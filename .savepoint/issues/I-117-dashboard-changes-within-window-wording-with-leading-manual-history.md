---
id: I-117
title: Dashboard changes within-window wording with leading manual history
type: defect
status: open
source:
  kind: check
  check: C-957
  actor: {role: checker, session: check-o032-20261002-independent}
  at: '2026-10-02T05:23:00Z'
tasks: [T-098]
checks: [C-957, C-958]
guardrail_ids: [TEST-01, TEST-02]
severity: low
history:
  - at: '2026-10-02T05:23:00Z'
    actor: {role: checker, session: check-o032-20261002-independent}
    kind: observed
    check: C-957
    note: Initial independent Full O-032 Check; frozen scope and reproducible harness in C-957.
  - at: '2026-10-02T05:27:50Z'
    actor: {role: executor, session: repair-o032-20261002}
    kind: repair_attempted
    note: 'windowStart now loads the whole history when no official snapshot precedes the tenth newest. Added TestLoadDashboardIsUnchangedWithLeadingManualHistory (leading, manual-only, trailing, interleaved, ten-official cases).'
  - at: '2026-10-02T06:19:18Z'
    actor: {role: checker, session: recheck-o032-20261002-independent}
    kind: rechecked
    check: C-958
    note: "All 36 original official/manual count scenarios pass after windowStart repair. Technically proven, open pending CLEAR Check proof."
---

# I-117: Dashboard changes within-window wording with leading manual history

## Summary

T-098 Done When 2 and the approved T-097 decision promise byte-identical dashboard output for histories with at most ten official snapshots. Exactly ten official snapshots preceded by manual refreshes violate that promise.

## Evidence

C-957 independent `TestO032WindowIndependentMatrix/official=10/leading`: save two manual snapshots first, then ten official snapshots with increasing creation times, using valid existing fixture builders. `LoadSnapshots` returns 12; `LoadWindow` returns 10 with `cut=true`. The dashboard's basis changes from `Compared with 9 earlier comparable official checks. 2 manual results shown, not counted.` to `Compared with the 9 most recent comparable official checks. Older history was not read.` Nothing was beyond the promised <=10-official boundary. `internal/codehealth/window.go:88-97` returns immediately at the tenth official without determining whether there are any older official snapshots. Existing `TestLoadDashboardIsUnchangedWithinTheWindow` uses histories beginning with an official snapshot and misses leading manuals. Independent matrix covered official counts 0-11 crossed with leading/trailing/interleaved manuals; this is its only failed cell.

## Proof Needed

Prove whole-Dashboard and Chip equivalence at 0-10 official snapshots including leading manuals, manual-only, interleaved and trailing cases. Keep the approved bounded behavior for >10 official snapshots, retention and no-write rules. If the owner wants a different promise, obtain a specific acceptance change through planning rather than editing criteria to match this Check.

Repair directly under this Issue by default; preserve all owner-completed Task statuses. This record grants no clearance or owner acceptance.
