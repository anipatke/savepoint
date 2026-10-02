---
id: I-112
title: Health popover hides a capability with multiple instances
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
    note: 'Reproduced (TestHealthPopoverReachesEveryInstanceOfEverySignal). The five-row window now follows the cursor over every row with up/down cues, and movement is no longer capped at five. Proved with two Tests instances, two instances in the first and middle capabilities, and a blocking last capability.'
  - at: '2026-10-01T22:48:04Z'
    actor: {role: checker, session: recheck-o035-20261002-independent}
    kind: rechecked
    check: C-953
    note: 'Final capability is reachable now, but remains open: the five-row window scrolls on down selection with six configured instances, dropping an earlier instance and moving other rows. healthWindowTop (health_view.go:279) and uncapped moveHealth (health.go:167) introduce scrolling contrary to unchanged O-035 success condition 2 and T-079 no-key-scroll criteria. Frozen Popover multiple-instance/all-selection and never-scroll requirement cells; independent TestCheckO035MultipleInstances fails the no-scroll condition. Preserve accessible multi-instance information without scrolling, or obtain an explicit owner design decision.'
  - at: '2026-10-01T23:05:00Z'
    actor: {role: executor, session: user-request}
    kind: repair_attempted
    note: 'C-953 recheck: scrolling removed. Several instances of one signal now share one row for the worst instance (blocking first, then label severity) named "Tests x2", with "(worst of N: name)" in the selected detail. The window, cues and extra navigation are gone; TestHealthPopoverShowsOneRowPerSignalWithoutScrolling proves rows never move, the blocking final signal is reachable with its sign-off, at 80x20/24/40.'
  - at: '2026-10-01T23:35:44Z'
    actor: {role: checker, session: recheck-o035-targeted-20261002-independent}
    kind: rechecked
    check: C-954
    note: 'Repair verified by independent CLEAR Check C-954 against C-952 frozen scope; all original cases and admitted adjacent repair cases pass.'
---

# I-112: Health popover hides a capability with multiple instances

## Summary

Violates O-035 outcome and success conditions 3/4: five signals presented; existing supported multi-instance configuration retained. Found inside C-952's frozen initial Full Check matrix.

## Evidence

Independent scenario `TestCheckO035MultipleInstances` in a Go overlay (no implementation edits): Configure two instances for an earlier capability and one for the other four, giving six dashboard rows. For example prepend a second Tests row to the normal five. The renderer breaks after row five and movement clamps to index four. Dependency vulnerabilities is never visible or selectable, although its result can block sign-off. A saved multi-instance complexity config is supported by TestDashboardMultipleInstancesOfOneCapabilityEachGetARow.

Expected: the stated requirement remains true through the supported saved-data-to-board path. Actual: the required information is lost or misrepresented as described above.

internal/board/v2/health_view.go:220-223 drops rows beyond five; internal/board/v2/health.go:167-169 caps navigation. DashboardRows emits each configured instance, not one aggregate per signal. No board test covers that integration.

## Proof Needed

Ensure all five capabilities remain represented and any blocking configured instance has accessible meaning/sign-off information within the fixed frame. Choose a bounded aggregate or other explicit representation without silently dropping supported instances. Prove two instances in the first and middle capabilities, mixed health and a blocking final capability.

Direct repair under this Issue is the default. Keep completed Tasks done. A fresh independent Check must verify the frozen C-952 cells and current full gate.

## C-953 Recheck

Final capability is reachable now, but remains open: the five-row window scrolls on down selection with six configured instances, dropping an earlier instance and moving other rows. healthWindowTop (health_view.go:279) and uncapped moveHealth (health.go:167) introduce scrolling contrary to unchanged O-035 success condition 2 and T-079 no-key-scroll criteria. Frozen Popover multiple-instance/all-selection and never-scroll requirement cells; independent TestCheckO035MultipleInstances fails the no-scroll condition. Preserve accessible multi-instance information without scrolling, or obtain an explicit owner design decision.
