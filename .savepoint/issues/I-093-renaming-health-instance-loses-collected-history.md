---
id: I-093
title: Renaming a health instance loses collected history
type: defect
status: open
source:
  kind: check
  check: C-941
  actor: {role: checker, session: o028-check-20261001}
  at: '2026-10-01T08:38:18Z'
tasks: [T-055, T-057]
checks: [C-941]
guardrail_ids: [TEST-01, TEST-02]
history:
  - at: '2026-10-01T08:50:00Z'
    actor: {role: executor, session: o028-repair-20261001}
    kind: repair_attempted
    note: "Collect groups history by capability and provider, not name; Assess's series selection separates scopes and configurations. Tests: TestCollectKeepsHistoryAcrossARename (unnamed to named, named to named, changed scope), TestCollectKeepsTwoScopedInstancesSeparate. Gates: make build, make test-fast, make test-full (go1.26.2 linux/amd64) pass; awaiting a checker."
    check: C-941
---

# I-093: Renaming a health instance loses collected history

## Summary

T-055 requires rename to preserve a comparison series and history to keep instances separate by series identity. Collect partitions history by a key including name before Assess can select comparable observations. A rename therefore loses the old history despite an unchanged SeriesID and Digest.

## Evidence

C-941 frozen M2: collect three official observations with complexity 5 for name `api`, scope `api/**`, then collect a fourth with name `backend` and otherwise identical configuration. Expected: four comparable official observations and a steady range. Actual saved summary: `Not enough comparable official history for a recent range yet. Good needs three comparable official checks.`

`internal/codehealth/collect.go:184` groups history by `r.key()`; `collect.go:151` gives Assess only `history[r.key()]`. The key includes Name (`config.go:83`, `snapshot.go:201`), although SeriesID intentionally excludes it (`identity.go:114`). Existing `TestRenameKeepsSeriesButScopeStartsOne` checks only the hash; `TestInstancesAreAssessedSeparately` supplies Assess a complete history directly, bypassing collection's filtering.

## Proof Needed

An integrated stored-history regression must preserve history across a rename while keeping different scoped instances separate. Check unnamed→named, named→named, unchanged name with changed scope/exclusions, and two independent instances within frozen M2. Re-run full gate.
