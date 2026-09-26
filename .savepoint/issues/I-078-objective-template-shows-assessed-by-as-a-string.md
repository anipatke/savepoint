---
id: I-078
title: Objective template shows freshness.assessed_by as a string
type: defect
status: open
source:
  kind: report
  actor: {role: owner, session: user}
  at: '2026-09-26T05:18:49Z'
severity: low
history:
  - at: '2026-09-26T05:18:49Z'
    actor: {role: owner, session: user}
    kind: observed
    note: >-
      Raised by an independent review of the packaged Savepoint skills; claim verified against the code in a follow-up review session before capture.
---

# I-078: Objective template shows freshness.assessed_by as a string

## Summary

The Objective template in `savepoint-design` shows
`assessed_by: role/session`, a single string. The runtime decodes
`assessed_by` as a `{role, session}` map, so an Objective that copies the
template's shape fails to load as malformed evidence.

## Evidence

- `agent-skills/savepoint-design/SKILL.md:133`: `assessed_by: role/session`.
- `internal/data/evidence_v2.go:98` types it as `evidenceActorFrontmatter`
  (`role`, `session`); `decodeV2Actor` (`:374`) requires both.
- Every other actor in the skills is already written as
  `{role: ..., session: ...}`.

## Proof Needed

- The template shows `assessed_by: {role: checker, session: ...}`.
- An Objective with a freshness block copied from the template loads.
- Packaged template copies stay byte-identical.
