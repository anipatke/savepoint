---
id: I-108
title: Health sparklines drop early and restart notes
type: defect
status: resolved
source:
  kind: check
  check: C-952
  actor: {role: checker, session: check-o035-20261002-independent}
  at: '2026-10-01T22:30:12Z'
tasks: [T-078, T-079]
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
    note: 'Reproduced (TestHealthPopoverShowsSparkNotesWithAndWithoutBlocks). The selected signal now shows a History line with the early/restart/thin note; rows without a sparkline say no trend yet. Tested with and without blocks at 80x20, 80x24, 120x40. Data-level point counts and restarted series are covered by the existing codehealth spark tests.'
  - at: '2026-10-01T22:48:04Z'
    actor: {role: checker, session: recheck-o035-20261002-independent}
    kind: rechecked
    check: C-953
    note: 'Repair reproduced and proven: selected detail now shows early, restart and thin history notes with and without blocks. Original Spark output/all-selection cells and focused regressions pass. No technical verified resolution is recorded yet because the superseding Objective Check is NEEDS WORK; issue-capture.md reserves verified for proof from a CLEAR Check.'
  - at: '2026-10-01T23:35:44Z'
    actor: {role: checker, session: recheck-o035-targeted-20261002-independent}
    kind: rechecked
    check: C-954
    note: 'Repair verified by independent CLEAR Check C-954 against C-952 frozen scope; all original cases and admitted adjacent repair cases pass.'
---

# I-108: Health sparklines drop early and restart notes

## Summary

Violates O-035 success condition 6. Found inside C-952's frozen initial Full Check matrix.

## Evidence

Independent scenario `TestCheckO035SparkNotes` in a Go overlay (no implementation edits): Give a row Spark=▁▄█, SparkWord=better and SparkNote=early, restarted because settings changed, or both. The renderer returns only ▁▄█ better in all three cases. No selected-detail line renders SparkNote either.

Expected: the stated requirement remains true through the supported saved-data-to-board path. Actual: the required information is lost or misrepresented as described above.

internal/board/v2/health_view.go:254-264 only reads SparkNote when Spark is empty. Data tests prove the note exists, but the renderer test only checks blocks and direction.

## Proof Needed

Render early and restart information whenever applicable. Prove 2/3/4/5/10 points and changed settings through the final popover, including a restarted series that already has three comparable points.

Direct repair under this Issue is the default. Keep completed Tasks done. A fresh independent Check must verify the frozen C-952 cells and current full gate.

## C-953 Recheck

Repair reproduced and proven: selected detail now shows early, restart and thin history notes with and without blocks. Original Spark output/all-selection cells and focused regressions pass. No technical verified resolution is recorded yet because the superseding Objective Check is NEEDS WORK; issue-capture.md reserves verified for proof from a CLEAR Check.
