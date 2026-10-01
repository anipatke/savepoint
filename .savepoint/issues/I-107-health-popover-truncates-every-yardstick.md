---
id: I-107
title: Health popover truncates every yardstick
type: defect
status: resolved
source:
  kind: check
  check: C-952
  actor: {role: checker, session: check-o035-20261002-independent}
  at: '2026-10-01T22:30:12Z'
tasks: [T-079]
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
    note: 'Reproduced with real-length copy (TestHealthPopoverKeepsEveryYardstickWithRealCopy). Rows now keep the aim and direction in full; name (short form for Dependency vulnerabilities) then value give way, a value ending in its number keeps the number, and the popover widens to 100 columns. Verified at 80x20, 80x24, 80x40 and 200x50.'
  - at: '2026-10-01T22:48:04Z'
    actor: {role: checker, session: recheck-o035-20261002-independent}
    kind: rechecked
    check: C-953
    note: 'Real-copy aims and full direction are now visible at 80x20/24/40. Remains open: actual outer width is 76 at all three sizes, violating the unchanged T-079 maximum 72. healthPopoverMaxWidth is now 100 (health_view.go:17,105). Frozen Popover normal-size/requirements cell, independent TestCheckO035FrozenDimensions fails. Preserve yardsticks while restoring the specified cap, or obtain an explicit owner design decision before changing it.'
  - at: '2026-10-01T23:05:00Z'
    actor: {role: executor, session: user-request}
    kind: repair_attempted
    note: 'C-953 recheck: width cap restored. Popover is at most 72 cells (healthPopoverMaxWidth=72); real-copy test now asserts the rendered width at 80x20/24/40 and 200x50 alongside aims and direction. Name cap raised to 14 so "Complexity x2" fits.'
  - at: '2026-10-01T23:35:44Z'
    actor: {role: checker, session: recheck-o035-targeted-20261002-independent}
    kind: rechecked
    check: C-954
    note: 'Repair verified by independent CLEAR Check C-954 against C-952 frozen scope; all original cases and admitted adjacent repair cases pass.'
  - at: '2026-10-01T23:40:00Z'
    actor: {role: executor, session: user-request}
    kind: owner_decision
    note: "Owner-directed follow-up: the row's value column is now the number alone (new DashboardRow.Figure; Value keeps the full wording and appears on the selected detail's question line). Tests rows are labelled Tests failing and show the failing count; Dependencies show the total, with ! when blocking; unmeasured shows an em dash. Frees room so name, ten-block spark and aim fit at 72 columns."
---

# I-107: Health popover truncates every yardstick

## Summary

Violates O-035 success conditions 4 and 11; T-079 row content and fit requirements. Found inside C-952's frozen initial Full Check matrix.

## Evidence

Independent scenario `TestCheckO035RealCopy` in a Go overlay (no implementation edits): Use the normal five signal names with values all 3,045 pass; 86%; hardest function scores 46; 8% copy-pasted; 2 unrated, their real default aims, and a three-block spark. Render at 80x20, 80x24, 80x40. All five aims disappear; worse is cut to wor… despite the frame having room geometrically. Even a wider terminal retains the 72-column cap.

Expected: the stated requirement remains true through the supported saved-data-to-board path. Actual: the required information is lost or misrepresented as described above.

internal/board/v2/health_view.go:215-238 pads every name and value to the longest row, then appends the spark and aim and truncates the complete line. Tests use Value=42 and the shorter synthetic aim: 80 or more, so they miss realistic lengths.

## Proof Needed

Reserve visible room for the required yardstick and direction using real projection copy, and prove all five rows at each documented size with thin/early/full history and long values. Truncation may shorten optional text but must not remove the scale.

Direct repair under this Issue is the default. Keep completed Tasks done. A fresh independent Check must verify the frozen C-952 cells and current full gate.

## C-953 Recheck

Real-copy aims and full direction are now visible at 80x20/24/40. Remains open: actual outer width is 76 at all three sizes, violating the unchanged T-079 maximum 72. healthPopoverMaxWidth is now 100 (health_view.go:17,105). Frozen Popover normal-size/requirements cell, independent TestCheckO035FrozenDimensions fails. Preserve yardsticks while restoring the specified cap, or obtain an explicit owner design decision before changing it.
