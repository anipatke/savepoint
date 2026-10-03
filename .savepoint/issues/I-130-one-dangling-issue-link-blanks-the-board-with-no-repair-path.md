---
id: I-130
title: One dangling issue link blanks the board with no repair path
type: defect
status: open
source:
  kind: report
  actor: {role: owner, session: owner-report-20261003}
  at: '2026-10-03T00:00:00Z'
severity: medium
history:
  - at: '2026-10-03T00:00:00Z'
    actor: {role: owner, session: owner-report-20261003}
    kind: observed
    note: Owner saw INVALID PROJECT DATA on a project whose issue I-017 listed missing check C-005; the board offered only q:quit.
---

# I-130: One dangling issue link blanks the board with no repair path

## Summary

A single issue→check or issue→task reference to a record that does not exist
makes the whole project fail to load, so the board draws only the diagnostic
screen. Strict loading is intentional, but the failure gives the owner little
help recovering.

## Evidence

- `internal/data/project.go:322` returns `ErrV2IssueMissingLinkTarget` for the
  first dangling issue check link; `:317` does the same for a task link.
- `internal/board/v2/view.go` `renderDiagnostic` draws only the heading, the
  first diagnostic and `q:quit`. Only the first error is reported, so fixing one
  link can reveal the next.
- `internal/doctor/checks.go:383` names the error
  (`v2-issue-missing-link-target`) but has no repair for it.
- Interim: the diagnostic screen now points to `savepoint doctor`.

## Proof Needed

Add a doctor repair for a dangling issue check/task link, ship it with the next
version, and confirm the diagnostic screen's hint leads to it. Consider whether
the loader should report every dangling link rather than only the first.
