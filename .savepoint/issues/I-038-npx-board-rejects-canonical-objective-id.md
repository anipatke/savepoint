---
id: I-038
title: npx board rejects a canonical hyphenated Objective ID
type: defect
status: resolved
source:
  kind: report
  actor: {role: owner, session: user}
  at: '2026-09-23T20:33:42Z'
resolution:
  disposition: accepted
  actor: {role: owner, session: user}
  at: '2026-09-23T20:53:54Z'
  reason: >-
    Owner reported that the board loads now and explicitly directed resolution.
    This records owner acceptance, not technical clearance from an independent
    Check.
history:
  - at: '2026-09-23T20:33:42Z'
    actor: {role: executor, session: task-report-20260924}
    kind: observed
    note: >-
      Owner reports `npx savepoint board` refuses
      `objectives/O-001-release-validation-cutover/Objective.md` because ID
      `O-001` must match O plus at least three digits.
  - at: '2026-09-23T20:53:54Z'
    actor: {role: owner, session: user}
    kind: owner_decision
    note: >-
      Owner said "loads now, resolve the defect." Issue closed as accepted on
      the owner's report; this does not claim an independent technical Check.
---

# I-038: npx board rejects a canonical hyphenated Objective ID

## Summary

The reported board command rejects Objective ID `O-001`, even though the
current V2 identity grammar requires `O-` plus at least three digits and
accepts `O-001`.

## Evidence

- Owner report: `npx savepoint board` returns `[invalid_v2]` for
  `objectives/O-001-release-validation-cutover/Objective.md` and says
  `objective id "O-001" must match O plus at least three digits`.
- The checked-out V2 parser uses the shared hyphenated identity grammar and
  its identity test includes `O-001` as a valid value.
- The checked-out parser's diagnostic includes the hyphen in `O-`, unlike the
  reported diagnostic. The CLI may be stale or built from a different source
  version; the emitting package/build has not been identified yet.

## Proof Needed

- Identify the package and binary version emitted by `npx savepoint board` and
  compare it with this checkout's identity grammar.
- Verify that the reported project loads with a build containing the current
  parser. If it still fails, inspect the exact Objective frontmatter bytes and
  continue from the first differing validation path.
- Ensure the shipped npm-facing CLI accepts `O-001` and retains rejection of
  genuinely malformed identities.
