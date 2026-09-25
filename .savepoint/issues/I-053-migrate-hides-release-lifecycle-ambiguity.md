---
id: I-053
title: Migrate fails before showing a release lifecycle decision the owner must make
type: defect
status: resolved
source:
  kind: check
  check: C-923
  actor: {role: checker, session: o022-objective-check-20260925}
  at: '2026-09-25T03:27:22Z'
tasks: [T-039]
checks: [C-923, C-924]
resolution:
  disposition: verified
  check: C-924
  actor: {role: checker, session: o022-independent-recheck-20260925}
  at: '2026-09-25T04:53:45Z'
guardrail_ids: [DATA-03]
severity: low
history:
  - at: '2026-09-25T03:27:22Z'
    actor: {role: checker, session: o022-objective-check-20260925}
    kind: observed
    check: C-923
    note: An audited or unrecognized V1 release status with an active epic makes Plan return "has no resolvable V2 Goal" instead of listing the blocking ambiguity.
  - at: '2026-09-25T04:29:53Z'
    actor: {role: executor, session: o022-repair-followup-20260925}
    kind: repair_attempted
    note: >-
      Planning an Objective without a mapped Goal now defers when that
      release has a blocking completion or lifecycle ambiguity, allowing the
      migration preview to report the owner decision instead of failing with
      "no resolvable V2 Goal". Added the audited release with an active epic
      regression; it confirms a blocked preview with the lifecycle choices and
      no misleading Goal error. The focused migration test and golden checks
      passed; git diff --check passed. Issue remains open for independent
      verification.
  - at: '2026-09-25T04:53:45Z'
    actor: {role: checker, session: o022-independent-recheck-20260925}
    kind: rechecked
    check: C-924
    note: >-
      C-924 independently verified that audited and unrecognized lifecycle
      statuses with active work appear as owner decisions in preview, block
      Apply before writes, and plan when the decision is supplied.
---

# I-053: Migrate fails before showing a release lifecycle decision the owner must make

## Summary

T-039 fails the preview with a named error for an epic whose release has "no
resolvable" Goal. Its Boundaries also rule out changes to legacy completion
handling. A V1 release with status `audited`, or an unrecognized status,
produces a blocking ambiguity. The preview lists that ambiguity so the owner
can resolve it with `--decisions`. If the release also holds an active epic,
`planEpic` now fails the whole plan first. The preview never shows the
ambiguity, and the error hides the real cause and the decision that fixes it.

## Evidence

- A temporary `internal/migrate` probe (removed after the run) used release
  `v1` with `status: audited` and one `in_progress` epic with a planned Task.
  `Plan` returned
  `plan Objective for V1 epic v1/E02-x: V1 release "v1" has no resolvable V2 Goal`.
  `status: weird` returned the same error.
- The same fixture with the `release_completion_evidence:.savepoint/releases/v1/v1-PRD.md`
  decision set to `in_progress` planned successfully (`Appliable: true`).
  The owner could resolve it, but the preview never says which decision is needed.
- `internal/migrate/plan.go` `planEpic` returns the new error when
  `releaseIDs[release]` is empty. `resolveReleaseStatus` leaves it empty while
  the release's lifecycle ambiguity is unresolved.
- `TestPlan_ambiguousReleaseDispositionBlocksCutover` covers only a release
  with no epics.

## Proof Needed

- A release whose lifecycle is waiting on an owner decision produces its
  ambiguity in the preview, even when it holds active epics. The preview is not
  appliable, and it no longer fails with "no resolvable V2 Goal".
- A regression test covers an audited release that holds an active epic.
