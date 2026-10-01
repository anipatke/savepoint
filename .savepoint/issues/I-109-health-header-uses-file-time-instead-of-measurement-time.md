---
id: I-109
title: Health header uses file time instead of measurement time
type: defect
status: resolved
source:
  kind: check
  check: C-952
  actor: {role: checker, session: check-o035-20261002-independent}
  at: '2026-10-01T22:30:12Z'
tasks: [T-076]
checks: [C-952, C-953, C-954]
guardrail_ids: [TEST-01, TEST-02]
resolution:
  disposition: verified
  check: C-954
  actor: {role: checker, session: recheck-o035-targeted-20261002-independent}
  at: '2026-10-01T23:35:44Z'
  reason: 'Frozen original and remediation scenarios pass in CLEAR Check C-954, with fresh full gate and official health evidence.'
history:
  - at: '2026-10-01T22:40:00Z'
    actor: {role: executor, session: user-request}
    kind: repair_attempted
    note: 'Reproduced (TestLatestSnapshotFollowsStoredRecencyNotFileTime). LatestSnapshot now orders by stored CreatedAt then ID, as LoadSnapshots does, reading each file only for its time and fully decoding only the newest. Proved with reversed, equal and tied times, and the header chip equals Dashboard.Chip.'
  - at: '2026-10-01T22:48:04Z'
    actor: {role: checker, session: recheck-o035-20261002-independent}
    kind: rechecked
    check: C-953
    note: 'Repair reproduced and proven: stored CreatedAt and identity decide recency; reversed mtimes, equal mtimes and timestamp ties match dashboard ordering. Original Chip recency cells pass independently. Verified resolution awaits a CLEAR Check under issue-capture.md.'
  - at: '2026-10-01T23:35:44Z'
    actor: {role: checker, session: recheck-o035-targeted-20261002-independent}
    kind: rechecked
    check: C-954
    note: 'Repair verified by independent CLEAR Check C-954 against C-952 frozen scope; all original cases and admitted adjacent repair cases pass.'
---

# I-109: Health header uses file time instead of measurement time

## Summary

Violates O-035 success condition 1; T-076 newest-snapshot lookup and refresh consistency. Found inside C-952's frozen initial Full Check matrix.

## Evidence

Independent scenario `TestCheckO035ChipRecency` in a Go overlay (no implementation edits): Save a later official snapshot with Watch overall and an older manual snapshot with failing tests and Needs Attention. Set the later file mtime earlier than the older file (as after restore/copy). LoadChip shows Needs Attention; LoadDashboard and Dashboard.Chip show Watch. Opening H changes the header even though no measurement changed.

Expected: the stated requirement remains true through the supported saved-data-to-board path. Actual: the required information is lost or misrepresented as described above.

internal/codehealth/chip.go:97-129 selects by filesystem mtime; internal/codehealth/storage.go:205-255 orders by stored CreatedAt then ID. The chip test deliberately aligns mtime and measurement ordering. T-076 records the divergence as a limitation, but no owner exception permits it.

## Proof Needed

Keep the cheap newest lookup while using the same persisted recency semantics as the dashboard. Prove reversed save order/mtime, equal mtimes, restored histories, tie ordering, and matching header before/after H.

Direct repair under this Issue is the default. Keep completed Tasks done. A fresh independent Check must verify the frozen C-952 cells and current full gate.

## C-953 Recheck

Repair reproduced and proven: stored CreatedAt and identity decide recency; reversed mtimes, equal mtimes and timestamp ties match dashboard ordering. Original Chip recency cells pass independently. Verified resolution awaits a CLEAR Check under issue-capture.md.
