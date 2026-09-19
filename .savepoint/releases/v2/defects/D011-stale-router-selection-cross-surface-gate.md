---
id: v2/D011-stale-router-selection-cross-surface-gate
release: v2
status: open
severity: high
title: "Stale router selection needs one cross-surface real-UX regression gate"
---

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

Promote the same-disk stale-router fixture to an explicit cross-surface
regression gate covering `resume`, board startup, and board reload. Assert the
diagnostic, the independent next action, and the absence of a substitute
selection.

## Acceptance Criteria

- [ ] One stale-router fixture is exercised through real resume output and the
      board presentation/reload path.
- [ ] Both surfaces retain the missing ID and show the project-wide action.
- [ ] Tests fail if a different record is substituted or the diagnostic is
      hidden.

## Resolution Notes

Pending.
