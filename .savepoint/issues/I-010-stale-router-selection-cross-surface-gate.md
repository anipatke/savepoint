---
id: I-010
title: Stale router selection needs one cross-surface real-UX regression gate
type: defect
status: resolved
source:
    kind: migration
    actor:
        role: owner
        session: .savepoint/releases/v2/defects/D011-stale-router-selection-cross-surface-gate.md
    at: "2026-09-20T05:03:29Z"
severity: high
resolution:
    disposition: accepted
    actor: {role: owner, session: user-review-20260922}
    at: "2026-09-22T09:11:47Z"
    reason: >-
        Owner accepts the absence of one combined cross-surface regression fixture.
        The shipped behavior already has focused coverage, so the remaining request is
        treated as a non-blocking test observation rather than a current defect.
history:
    - at: "2026-09-22T09:11:47Z"
      actor: {role: owner, session: user-review-20260922}
      kind: owner_decision
      note: >-
          Reclassified the missing combined regression gate as a non-blocking
          observation and accepted it without a repair or technical CLEAR.
---
## Migrated from V1

Relocated verbatim from the V1 defect body at `.savepoint/releases/v2/defects/D011-stale-router-selection-cross-surface-gate.md` (release `v2`).

## V1 Body (verbatim)

# D011: Stale router selection needs one cross-surface real-UX regression gate

## Symptom

The projection carries a `SelectionDiagnostic` while continuing with a
project-wide `Next` action, and the repository already has focused resume and
board assertions for parts of that behavior. The coverage is split across
constructed projections and separate fixtures; the release gate still needs
one explicit stale-router project exercised through both shipped surfaces so a
future renderer cannot bury the diagnostic beneath the recommendation or
substitute a similarly numbered record.

## Expected Behavior

When the router names a Task or Objective that no longer exists, `savepoint
resume` and the board must both show the stale-selection diagnostic alongside
the project-wide next action, preserve the exact missing identity, and never
pretend that another record is the selected one.

## Reproduction

1. Start with a valid V2 project whose router selects a live Task or
   Objective.
2. Remove or rename that selected record without repairing the router.
3. Invoke the real `resume` command and open/reload the board from the same
   project.
4. Verify that both surfaces show the missing selection and the available
   project-wide action together.

## Impact

A regression in either entry point could make stale context look authoritative
or hide the diagnostic, causing work to be performed against the wrong record
or leaving operators unaware that the router is stale.

## Fix Plan

Confirmed approach: promote the same-disk stale-router fixture to an explicit
cross-surface regression gate covering `resume`, board startup, and board
reload. Assert the diagnostic, the independent next action, and the absence of
a substitute selection.

## Acceptance Criteria

- [ ] One stale-router fixture is exercised through real resume output and the
      board presentation/reload path.
- [ ] Both surfaces retain the missing ID and show the project-wide action.
- [ ] Tests fail if a different record is substituted or the diagnostic is
      hidden.

## Resolution Notes

Resolved by explicit owner acceptance after reclassification as a non-blocking
test observation. This is not a repair or technical `CLEAR`.
