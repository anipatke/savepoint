---
id: F003
title: Duplicate of shared-epic dependency resolution finding
status: duplicate
severity: medium
confidence: high
source_auditor: agent
work_item: E01-example/T001-shared
releases:
  - v1.1
tasks:
  - E01-example/T001-shared
duplicate_of: F001
first_seen: 2026-06-25
last_seen: 2026-07-01
proof_needed: "n/a; canonical finding is F001"
---

## Summary

Independently filed report of the same shared-epic dependency resolution risk already
tracked as `F001`, found again before reconciliation.

## Evidence

Same symptom and location as `F001`: `E01-example/T002-follow-up`'s short dependency
reference and the same-short-ID collision with the done v1 Task.

## Proof

No independent proof is tracked here; resolution follows `F001`'s proof, not this
record's.

## History

- 2026-06-25: Filed independently before the 2026-07-01 reconciliation pass.
- 2026-07-01: Reconciled as `duplicate`, pointing at the canonical `F001`. The stable ID
  `F003` is retained rather than deleted.
