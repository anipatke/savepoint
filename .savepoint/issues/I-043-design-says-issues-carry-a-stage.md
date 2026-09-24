---
id: I-043
title: Design.md says an in-progress Issue needs a stage, but the loader rejects one
type: drift
status: resolved
source:
  kind: report
  actor: {role: planner, session: o014-design-20260924}
  at: '2026-09-24T00:00:00Z'
tasks: []
checks: [C-917]
guardrail_ids: [TPL-02]
severity: low
resolution:
  disposition: verified
  check: C-917
  actor: {role: checker, session: o014-objective-recheck-20260924}
  at: '2026-09-24T09:49:09Z'
history:
  - at: '2026-09-24T00:00:00Z'
    actor: {role: planner, session: o014-design-20260924}
    kind: observed
    note: >-
      Found while planning O-014's Issue Next line. Design.md section 4 and
      internal/data/issue_v2.go disagree.
  - at: '2026-09-24T09:31:25Z'
    actor: {role: executor, session: o014-i043-repair-20260924}
    kind: repair_attempted
    note: >-
      Confirmed Design.md section 4 already says Issues carry no stage, matching
      decodeIssueStatus in internal/data/issue_v2.go. No further code or Design
      edit was needed. The existing repair remains open for an independent
      Check or an explicit owner acceptance decision.
  - at: '2026-09-24T09:49:09Z'
    actor: {role: checker, session: o014-objective-recheck-20260924}
    kind: rechecked
    check: C-917
    note: "Design section 4 matches the Issue decoder: Issues carry no stage."
---

# I-043: Design.md says an in-progress Issue needs a stage, but the loader rejects one

## Summary

`.savepoint/Design.md` section 4 says: "Issues use `open`, `in_progress`, and
`resolved`; `stage` is required only while an Issue is `in_progress`."

`internal/data/issue_v2.go` does the opposite: `decodeIssueStatus` rejects
any Issue that declares a `stage` ("issue status carries no stage"), and the
`IssueStatus` comment says an Issue "never carries a stage".

An agent following Design.md would write an Issue the project cannot load.

## Repair

Correct the Design.md sentence to match the code: Issues carry no stage. No
code change. Small enough to fix directly in O-014's T-034 Design
reconciliation, or in any Task that touches Design.md section 4.
