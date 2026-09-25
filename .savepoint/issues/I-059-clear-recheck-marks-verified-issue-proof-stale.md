---
id: I-059
title: A CLEAR recheck marks every Issue verified by the prior Check as stale
type: defect
status: open
source:
  kind: report
  actor: {role: executor, session: doctor-rca-20260925}
  at: '2026-09-25T11:11:00Z'
severity: low
history:
  - at: '2026-09-25T11:11:00Z'
    actor: {role: executor, session: doctor-rca-20260925}
    kind: observed
    note: >-
      Live doctor reports v2-issue-verified-proof-superseded for I-050..I-053.
      Their proof C-924 was superseded by C-925, a CLEAR recheck that only
      corrected two reviewed file paths and re-verified all four Issues.
  - at: '2026-09-25T11:11:30Z'
    actor: {role: executor, session: doctor-rca-20260925}
    kind: repair_attempted
    note: >-
      InspectIssueConsistency now reports a superseded proof only when the
      latest Check for its scope did not record CLEAR. Added a CLEAR-recheck
      regression test (fails unfixed) and moved two older tests to a NEEDS
      WORK superseder. make build && make test-fast passed. Issue remains
      open for independent verification.
---

# I-059: A CLEAR recheck marks every Issue verified by the prior Check as stale

## Summary

`InspectIssueConsistency` reports a verified Issue's proof as superseded
whenever any later Check exists for the same scope, even one that records
CLEAR and re-confirms the fix. Every recheck then puts every Issue the
previous Check verified into doctor's error list until someone edits each
Issue's `resolution.check`.

Owner decision (2026-09-25): a later CLEAR Check keeps the proof valid; only
a later Check that does not record CLEAR makes it stale.

## Evidence

- `savepoint doctor` at a1b2b9f: 4 `v2-issue-verified-proof-superseded`
  errors (I-050..I-053, each listed twice), all naming C-924 superseded by
  C-925.
- C-925 `result: CLEAR`, `supersedes: C-924`; its verification table marks
  I-050, I-051, I-052, and I-053 `Fixed`.
- `internal/data/issue_v2.go` `InspectIssueConsistency` compared only
  `latest != proof.ID`.

## Repair Attempt

- `internal/data/issue_v2.go`: `InspectIssueConsistency` skips a verified
  Issue when the latest Check for its proof's scope recorded CLEAR.
- New `TestInspectIssueConsistency_clearRecheckKeepsProof` fails against the
  unfixed code. `TestInspectIssueConsistency_verifiedProofSuperseded`,
  `..._sortedOrderReturnsEveryProblem`, and doctor's
  `TestCheckProject_IssueVerifiedProofSuperseded` now use a NEEDS WORK
  superseder, so the stale-proof report is still covered. `TestE44_EpicScenario`
  now expects a CLEAR Objective recheck to keep the Issue's proof.
- `make build && make test-fast` passed (exit 0).
- Live doctor no longer reports I-050..I-053.

## Proof Needed

A test showing a later CLEAR Check on the proof's scope is not reported,
while a later NEEDS WORK Check still is. Live `savepoint doctor` reports no
`v2-issue-verified-proof-superseded` for I-050..I-053.
