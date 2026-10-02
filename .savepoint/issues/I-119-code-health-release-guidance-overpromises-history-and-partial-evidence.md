---
id: I-119
title: Code Health release guidance overpromises history and partial evidence
type: drift
status: resolved
source:
  kind: check
  check: C-957
  actor: {role: checker, session: check-o032-20261002-independent}
  at: '2026-10-02T05:23:00Z'
tasks: [T-096, T-098]
checks: [C-957, C-958, C-960]
guardrail_ids: [TPL-02]
severity: medium
resolution:
  disposition: verified
  check: C-960
  actor: {role: checker, session: recheck-o032-20261003-independent}
  at: '2026-10-02T19:20:00Z'
  reason: "Release copy qualifies sparse-series trend and partial values; Design.md now agrees with code and README."
history:
  - at: '2026-10-02T05:23:00Z'
    actor: {role: checker, session: check-o032-20261002-independent}
    kind: observed
    check: C-957
    note: Initial independent Full O-032 Check; frozen scope and reproducible harness in C-957.
  - at: '2026-10-02T05:27:50Z'
    actor: {role: executor, session: repair-o032-20261002}
    kind: repair_attempted
    note: 'README and CHANGELOG now qualify trend equivalence for series added or removed within older history, and describe partial evidence as incomplete (may carry a value, never Good). Classification unchanged. Design.md not edited.'
  - at: '2026-10-02T06:19:18Z'
    actor: {role: checker, session: recheck-o032-20261002-independent}
    kind: rechecked
    check: C-958
    note: "Release copy now qualifies sparse-series trend and partial measured values; original sparse fixture matches approved limit. Technically proven, open pending CLEAR Check proof."
  - at: '2026-10-02T19:20:00Z'
    actor: {role: checker, session: recheck-o032-20261003-independent}
    kind: rechecked
    check: C-960
    note: "CLEAR. Release copy qualifies sparse-series trend and partial values; Design.md now agrees with code and README."
---

# I-119: Code Health release guidance overpromises history and partial evidence

## Summary

T-096 Done When 1-2 and 5 require guidance to match implemented history and partial evidence. README promises unchanged trends unconditionally, while the approved window deliberately permits thinner trends for sparse series. README and CHANGELOG also say partial evidence is not measured, although supported partial reports can carry an honest measured value.

## Evidence

`README.md:284-287`: beyond the bounded window, "The trend, labels and sign-off are unchanged." T-097 explicitly limits trend equivalence to signals present in every window snapshot and allows thinner history for recently added/removed instances. Independent `TestO032SparseHistoryDocumentation`: 16 official snapshots, coverage present at indices 0,1,2,15 and absent at 3-14. Full history shows `Steady over 4 official checks, 90% to 92%` with four spark points; bounded dashboard shows `No trend yet` and no spark. This is the approved implementation limitation, so repair the copy, not the window's sparse-series policy. `README.md:339-341` and `CHANGELOG.md:73-74` describe partial evidence as not measured. `LizardReader` at `reader_lizard.go:75-85` can return `Partial:true` with a measured maximum CCN and reason; `TestPolyglotReadersKeepEachStackIndependent`, `TestDashboardPartialSupportStaysVisibleAndNotGood`, and direct empty/malformed-reader probes preserve this distinction. It is not healthy zero; it is an incomplete measurement.

## Proof Needed

Qualify trend preservation with the approved comparable-series/window limitation. Explain partial measurements as incomplete and never automatically Good, distinguishing them from failed/unavailable reports with no value. Reconcile README, CHANGELOG and relevant Design wording without changing classification or broadening the provider contract.

Repair directly under this Issue by default; preserve all owner-completed Task statuses. This record grants no clearance or owner acceptance.
