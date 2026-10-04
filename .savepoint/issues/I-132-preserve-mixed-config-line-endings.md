---
id: I-132
title: Preserve mixed config line endings when saving
type: defect
status: resolved
source:
  kind: check
  check: C-963
  actor: {role: checker, session: check-o037-20261004}
  at: '2026-10-04T02:31:20Z'
tasks: [T-108]
checks: [C-963, C-964]
guardrail_ids: [FS-01, TEST-03]
resolution:
  disposition: verified
  check: C-964
  actor: {role: checker, session: recheck-o037-20261004}
  at: '2026-10-04T02:43:55Z'
  reason: 'Original mixed-ending reproduction and both-separator off/on round-trip matrix preserve unrelated bytes.'
history:
  - at: '2026-10-04T02:31:20Z'
    actor: {role: checker, session: check-o037-20261004}
    kind: observed
    check: C-963
    note: Initial independent Full Objective Check.
  - at: '2026-10-04T02:43:55Z'
    actor: {role: checker, session: recheck-o037-20261004}
    kind: rechecked
    check: C-964
    note: 'CLEAR. Original mixed-ending reproduction and both-separator off/on round-trip matrix preserve unrelated bytes.'
---

# I-132: Preserve mixed config line endings when saving

## Summary

Preserve mixed config line endings when saving before O-037 clearance.

## Evidence

O-037 SC3 and T-108 Done When 2 promise preservation of unrelated authored bytes/comments.

Input bytes: `schema_version: 2\r\n# keep LF\nfeatures:\r\n  parallel_planning: false\r\n`. Load and call WriteParallelPlanning(path, true, config.FeatureSource()). Expected: only false changes to true. Actual: the unrelated comment's LF becomes CRLF. `internal/data/feature_preferences.go:144–150` normalizes all lines then rejoins with one ending.

Independent TestO037IndependentPreservation/mixed_endings fails exact byte comparison; homogeneous LF/CRLF tests miss this supported hand-edited configuration.

## Proof Needed

Preserve each original line separator while changing only the requested value/insertion. Verify mixed endings in both directions, homogeneous LF/CRLF, comments/Unicode, no final newline, repeat no-op, malformed layouts and stale writes in C-963 M2.
