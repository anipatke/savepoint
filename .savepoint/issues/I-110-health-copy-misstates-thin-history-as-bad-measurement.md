---
id: I-110
title: Health copy misstates thin history as bad measurement
type: defect
status: resolved
source:
  kind: check
  check: C-952
  actor: {role: checker, session: check-o035-20261002-independent}
  at: '2026-10-01T22:30:12Z'
tasks: [T-075]
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
    note: 'Reproduced (TestMeaningNamesTheCauseOfTheLabel). New whyLabel names the cause when the label is worse than the value earns: worse than recent checks, partial, out of date, or fewer than three comparable checks. Stored classification and thresholds are untouched; good readings with thin history no longer say to repair anything.'
  - at: '2026-10-01T22:48:04Z'
    actor: {role: checker, session: recheck-o035-20261002-independent}
    kind: rechecked
    check: C-953
    note: 'Repair reproduced and proven: first fresh coverage 90% now says Meets the aim; it stays Watch until three comparable checks. Passing tests, low complexity/duplication at 1/2/3, partial/stale and baseline-decline adjacent cases pass. Classifier and thresholds unchanged. Verified resolution awaits a CLEAR Check under issue-capture.md.'
  - at: '2026-10-01T23:35:44Z'
    actor: {role: checker, session: recheck-o035-targeted-20261002-independent}
    kind: rechecked
    check: C-954
    note: 'Repair verified by independent CLEAR Check C-954 against C-952 frozen scope; all original cases and admitted adjacent repair cases pass.'
---

# I-110: Health copy misstates thin history as bad measurement

## Summary

Violates O-035 success condition 5: explanation says what the result means; T-075 Meaning and NextStep. Found inside C-952's frozen initial Full Check matrix.

## Evidence

Independent scenario `TestCheckO035Meaning` in a Go overlay (no implementation edits): Save one complete fresh official check with coverage 90% (Good yardstick 80%). Assess correctly retains Watch because three comparable checks are needed. The new Meaning says the tests skip a fair amount of code and recommends adding tests, omitting the actual thin-history reason. The same label table calls passing fresh tests incomplete/out of date and complexity 5 getting hard to follow.

Expected: the stated requirement remains true through the supported saved-data-to-board path. Actual: the required information is lost or misrepresented as described above.

internal/codehealth/dashboard_copy.go:49-71 maps a composite classification directly to a measurement-level assertion; classification.go:118-122 explicitly assigns Watch to good readings with fewer than three official checks. Existing copy tests check all label cells for presence, not fidelity to the actual assessment cause.

## Proof Needed

Keep stored classification and thresholds unchanged. Make the plain explanation accurate for thin-history, confidence/stale/partial and baseline-worsening reasons as well as threshold-derived labels. Verify fresh high coverage, passing tests, low complexity/duplication at 1/2/3 checks and measured partial/stale outcomes; do not tell users to repair a good value solely because history is thin.

Direct repair under this Issue is the default. Keep completed Tasks done. A fresh independent Check must verify the frozen C-952 cells and current full gate.

## C-953 Recheck

Repair reproduced and proven: first fresh coverage 90% now says Meets the aim; it stays Watch until three comparable checks. Passing tests, low complexity/duplication at 1/2/3, partial/stale and baseline-decline adjacent cases pass. Classifier and thresholds unchanged. Verified resolution awaits a CLEAR Check under issue-capture.md.
