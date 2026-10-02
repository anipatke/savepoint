---
id: I-111
title: Health refresh exposes provider name
type: defect
status: resolved
source:
  kind: check
  check: C-952
  actor: {role: checker, session: check-o035-20261002-independent}
  at: '2026-10-01T22:30:12Z'
tasks: [T-080]
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
    note: 'Reproduced (TestHealthRefreshNeverNamesAProvider). Refresh progress now reads "Refreshing N of M: <Signal>" with the Esc hint; checked for all five signals in the popover and history mode. The older test that asserted the provider was corrected.'
  - at: '2026-10-01T22:48:04Z'
    actor: {role: checker, session: recheck-o035-20261002-independent}
    kind: rechecked
    check: C-953
    note: 'Repair reproduced and proven: progress omits provider names for all five signals, in signals/history modes, while retaining signal and cancellation hint. Original progress output cells pass. Verified resolution awaits a CLEAR Check under issue-capture.md.'
  - at: '2026-10-01T23:35:44Z'
    actor: {role: checker, session: recheck-o035-targeted-20261002-independent}
    kind: rechecked
    check: C-954
    note: 'Repair verified by independent CLEAR Check C-954 against C-952 frozen scope; all original cases and admitted adjacent repair cases pass.'
---

# I-111: Health refresh exposes provider name

## Summary

Violates O-035 success condition 9: provider names absent from default views; T-079 popover excludes providers. Found inside C-952's frozen initial Full Check matrix.

## Evidence

Independent scenario `TestCheckO035ProgressProvider` in a Go overlay (no implementation edits): Start refresh and deliver Coverage progress with provider go_cover, position 2 of 5. The default popover visibly says Refreshing 2 of 5: Coverage (go_cover).

Expected: the stated requirement remains true through the supported saved-data-to-board path. Actual: the required information is lost or misrepresented as described above.

internal/board/v2/health_view.go:31 and :84-95 include p.Provider in the progress format. The refresh test asserts Coverage and Esc but never bans the provider.

## Proof Needed

Keep the signal progress and cancel hint while omitting raw provider identifiers. Prove each configured signal during refresh in the popover and history mode; retain the no-provider assertion already used for measured rows.

Direct repair under this Issue is the default. Keep completed Tasks done. A fresh independent Check must verify the frozen C-952 cells and current full gate.

## C-953 Recheck

Repair reproduced and proven: progress omits provider names for all five signals, in signals/history modes, while retaining signal and cancellation hint. Original progress output cells pass. Verified resolution awaits a CLEAR Check under issue-capture.md.
