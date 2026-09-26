---
id: I-081
title: Quick Task Checks carry the full Objective check method
type: drift
status: resolved
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
  - at: '2026-09-26T05:29:40Z'
    actor: {role: executor, session: skill-review-fixes}
    kind: repair_attempted
    note: >-
      check-method.md now gives Quick mode its own six-step Quick Check Procedure with a short scope lock, and marks the coverage matrix, workflow and side-effect lock, and adversarial pass as Full mode only. Packaged template copies re-synced; make build and make
      test-fast passed.
  - at: '2026-09-26T05:30:33Z'
    actor: {role: owner, session: user}
    kind: owner_decision
    note: >-
      Owner waived the independent Check and instructed the executor to
      resolve this Issue. No Check was run and no technical CLEAR is implied.
resolution:
  disposition: accepted
  actor: {role: owner, session: user}
  at: '2026-09-26T05:30:33Z'
  reason: Owner accepted the skill repair committed to v2 in 03d4629 without a Check.
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

## Repair Attempt Evidence

- `references/check-method.md` adds `### Quick Check Procedure` under Quick And Full Evidence Modes; the three heavy sections open with a Full-mode-only line; re-check covers a Quick scope lock.
- All section headings the tests pin are unchanged; `savepoint-check` still loads the method in full.
- Packaged copies are byte-identical to `agent-skills/`; `make build` and `make test-fast` passed.
