---
id: I018
title: O900 no longer renders the pending Check outcome
type: drift
status: resolved
source:
  kind: check
  check: C906
  actor: {role: checker, session: o012-full-check-20260922}
  at: '2026-09-22T09:28:27Z'
tasks: [T005, T009]
checks: [C906, C907]
guardrail_ids: [TEST-01]
severity: medium
resolution:
  disposition: verified
  check: C907
  actor: {role: checker, session: o012-recheck-20260922}
  at: '2026-09-22T10:00:27Z'
  reason: O900 now renders all six retained outcomes through interactive and non-TTY paths with durable regression coverage.
history:
  - at: '2026-09-22T09:28:27Z'
    actor: {role: checker, session: o012-full-check-20260922}
    kind: observed
    note: Live O900 card and plain-output probes render five retained outcomes but never the pending [ ] CHECK outcome required by O012 and T005.
    check: C906
  - at: '2026-09-22T09:36:11Z'
    actor: {role: owner, session: user}
    kind: owner_decision
    note: Keep T005 done. Remediate this Objective-level finding through linked new work under O012; do not require the owner to reopen a completed Task.
    check: C906
  - at: '2026-09-22T09:41:35Z'
    actor: {role: planner, session: i018-remediation-20260922}
    kind: observed
    note: Linked new remediation Task T009 under O012 while preserving T003-T005 as completed history.
  - at: '2026-09-22T09:41:35Z'
    actor: {role: executor, session: i018-remediation-20260922}
    kind: repair_attempted
    note: T009 swapped disposable O900 roles so T903 is Done with no Check and renders pending CHECK while T909 remains the stale REVIEW example in progress; added temporary-project regression coverage for all six outcomes and wait/replan/owner blockers through interactive and non-TTY paths. Strict loading, focused tests, git diff check, make build, and make test pass. Ready for independent recheck; Issue remains open.
  - at: '2026-09-22T10:00:27Z'
    actor: {role: checker, session: o012-recheck-20260922}
    kind: rechecked
    note: C907 verified the live O900 role swap and durable interactive/non-TTY regression matrix; I018 is resolved.
    check: C907
---

# I018: O900 no longer renders the pending Check outcome

## Summary

O012 requires the disposable O900 Objective to visibly exercise every retained
Task-card review outcome. T005 assigns the pending `[ ] CHECK` example to T903,
but T003's final renderer deliberately suppresses missing and current clearance
on every open Task. T903 is still in progress at build, so neither the
interactive card path nor the non-TTY path renders the pending outcome.

## Evidence

- `.savepoint/objectives/O012-simplify-task-card-outcomes/Objective.md:46`
  requires O900 to visibly exercise every retained outcome.
- `.savepoint/objectives/O012-simplify-task-card-outcomes/tasks/T005-refresh-the-o900-visual-test-spread.md:30`
  requires visible pending and checked `CHECK` examples; line 162 records the
  now-false claim `T903 → [ ] CHECK`.
- `internal/board/v2/card.go:152` permits an open outcome only when
  `reviewOutcomeIsActionable` is true; missing clearance is intentionally
  false, so T903 cannot render `[ ] CHECK`.
- An independent C906 harness loaded the live O900 records through the real
  card projection and plain renderer. Both lacked `[ ] CHECK`; the card matrix
  otherwise contained `[✓] CHECK`, `[!] NEEDS WORK`, `[!] REVIEW`,
  `[✓] WAIVED`, `[✓] OWNER ACCEPTED`, and all three required blockers.
- Existing package tests pass because they test the badge function's pending
  text synthetically but do not assert that the live O900 fixture reaches it.

Expected: O900 provides a supported, unambiguous card that renders `[ ] CHECK`
in both interactive and plain output.

Actual: no O900 card renders `[ ] CHECK`; T903 renders only `▣ BUILD`.

## Proof Needed

- Through new remediation work under O012, reconcile T005's recorded evidence
  and the O900 fixture with the final open-card rule so one valid fixture card
  visibly renders `[ ] CHECK` without losing the promised column spread or
  another retained outcome. T005 remains `done`.
- Add a durable fixture-level regression test that projects an O900-equivalent
  record set through both card and non-TTY paths and asserts all six retained
  outcomes plus wait, replan, and owner blockers.
- Correct the stale audit-stage test comment and update T005's technical
  evidence to the behavior actually verified.
- Pass focused board/data tests, `git diff --check`, `make build`, and
  `make test`, then obtain a fresh Full Objective Check that supersedes C906.
