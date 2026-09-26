---
id: I-081
title: Quick Task Checks carry the full Objective check method
type: drift
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

# I-081: Quick Task Checks carry the full Objective check method

## Summary

`check-method.md` says Quick and Full modes "apply the method below at
different reach", so an optional Task Check on a small change still owes a
scope lock, a multi-axis coverage matrix, and the other Full-method records.
The per-axis rules are conditional ("when applicable"), but nothing gives
Quick mode a short procedure of its own.

## Evidence

- `agent-skills/references/check-method.md:39` ("Both modes apply the method
  below at different reach"), `:65-78` (scope lock), `:106-128` (coverage
  matrix axes), `:242-246` (admission ledger).
- The reference is 351 lines, the largest packaged skill file.

## Proof Needed

- Quick mode has its own short procedure (diff, acceptance criteria, edge
  probes on changed code, configured gates, relevant guardrails).
- The full method stays required for the mandatory Objective Check.
- Packaged template copies stay byte-identical.
