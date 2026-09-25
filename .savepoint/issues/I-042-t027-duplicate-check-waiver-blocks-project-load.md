---
id: I-042
title: T-027 has two check_waiver keys, so the repository's V2 project no longer loads
type: verification
status: resolved
source:
  kind: check
  check: C-915
  actor: {role: checker, session: o021-objective-recheck-20260924}
  at: '2026-09-23T23:04:15Z'
tasks: [T-027]
checks: [C-915]
severity: medium
resolution:
  disposition: verified
  check: C-915
  actor: {role: checker, session: o021-objective-recheck-20260924}
  at: '2026-09-23T23:31:55Z'
  reason: Owner directed the fix; the duplicate board-written check_waiver was removed, keeping the owner-chat waiver, and the project loads again.
history:
  - at: '2026-09-23T23:04:15Z'
    actor: {role: checker, session: o021-objective-recheck-20260924}
    kind: observed
    check: C-915
    note: >-
      Reproduced with the rebuilt binary against this repository during the
      O-021 re-check.
  - at: '2026-09-23T23:31:55Z'
    actor: {role: checker, session: o021-objective-recheck-20260924}
    kind: rechecked
    check: C-915
    note: Removed the duplicate board-written check_waiver at the owner's direction; savepoint doctor and resume now load the project.
---

# I-042: T-027 has two check_waiver keys, so the repository's V2 project no longer loads

## Summary

`T-027-remove-stale-project-loader-references.md` has two top-level
`check_waiver` keys:

- Lines 10–14: the owner's chat waiver, `owner-chat-20260924`, recorded at 23:01:49Z.
- Lines 15–21: the board's waiver, `board-owner`, recorded at 23:01:55Z.

YAML rejects the duplicate key, so `data.LoadV2Index` fails for the whole
project. Every ordinary command on this repository stops at the same point:

```text
resume: loading project: v2 record could not be decoded: objectives/O-021-remove-v1-runtime-and-slim-migration/tasks/T-027-remove-stale-project-loader-references.md: yaml: unmarshal errors:
  line 15: mapping key "check_waiver" already defined at line 10
```

`savepoint doctor` reports the same problem as `[v2-record-malformed]`. The
board cannot draw.

Because of this, T-027's `done` status and waiver cannot be read by the
`internal/data` resolvers. O-021 cannot show that every owned Task is done,
and cannot close.

The root cause is not established. The board's `WriteTaskEvidenceV2`
replaces an existing key in place (`setMappingNode`). The two timestamps are
six seconds apart, which suggests a concurrent hand edit and board write.
That may relate to I-031. This Issue does not claim a product defect.

## Evidence

- `./savepoint resume` and `./savepoint doctor` in this repository, built
  from the working tree on 2026-09-24.
- The T-027 frontmatter, lines 10–21.

## Proof Needed

The owner keeps exactly one of the two waivers; which one is the owner's
choice. After that, `savepoint doctor` on this repository reports no
`v2-record-malformed` diagnostic, and `savepoint resume` loads the
project. The fix is metadata-only, so C-915's full-gate evidence stays
reusable under TEST-08.
