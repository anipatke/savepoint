---
id: F001
title: Shared-epic dependency resolution lacked regression coverage
status: fixed
severity: medium
confidence: high
source_auditor: agent
work_item: E01-example/T001-shared
releases:
  - v1.1
epics:
  - E01-example
tasks:
  - E01-example/T001-shared
guardrail_ids:
  - TEST-05
locations:
  - internal/data/dependency.go
first_seen: 2026-06-20
last_seen: 2026-07-01
proof_needed: A regression test proving the same-epic short reference resolves within the active release, not a prior release's completed Task of the same short ID
reviewer_note: "kept from v1 audit tooling; not part of the current finding schema"
---

## Summary

The shared-epic short dependency reference risked resolving to the completed v1
baseline Task instead of the in-progress v1.1 Task sharing the same short ID.

## Evidence

`E01-example/T002-follow-up` (v1.1) depends on the short reference `T001-shared`, and
the same short ID also names a `done` Task in the v1 release, `E01-example/T001-shared`.

## Proof

Not yet verified: the regression test named in `proof_needed` has not landed. Raw
`fixed` status is not proof of verification.

## History

- 2026-06-20: Opened during the v1.1 baseline review.
- 2026-07-01: Marked `fixed` pending the named regression test; still unverified.
