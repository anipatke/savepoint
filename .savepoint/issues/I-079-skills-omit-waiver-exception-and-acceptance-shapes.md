---
id: I-079
title: Skills describe waiver, exception, and owner acceptance without their frontmatter shapes
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
  - at: '2026-09-26T05:25:10Z'
    actor: {role: executor, session: skill-review-fixes}
    kind: repair_attempted
    note: >-
      savepoint-task shows the check_waiver shape; savepoint-check shows owner_validation acceptance and exception shapes, all from evidence_v2.go field names. Packaged template copies re-synced; make build and make
      test-fast passed.
  - at: '2026-09-26T05:26:38Z'
    actor: {role: owner, session: user}
    kind: owner_decision
    note: >-
      Owner waived the independent Check and instructed the executor to
      resolve this Issue. No Check was run and no technical CLEAR is implied.
resolution:
  disposition: accepted
  actor: {role: owner, session: user}
  at: '2026-09-26T05:26:38Z'
  reason: Owner accepted the skill repair committed to v2 in 0ddf034 without a Check.
---

# I-079: Skills describe waiver, exception, and owner acceptance without their frontmatter shapes

## Summary

The skills say a Task-check waiver must name "the Task, reason, actor, and
time", and describe owner exceptions and owner acceptance, but never show the
frontmatter the runtime reads. An agent recording one of these by hand has to
guess the field names; prose in the body is ignored by the gates.

## Evidence

- Waiver described in prose only: `savepoint-task/SKILL.md:43,61`,
  `savepoint-check/SKILL.md:103`, `references/check-method.md`,
  scaffold `AGENTS.md:46-49`.
- Runtime shapes in `internal/data/evidence_v2.go`: `check_waiver`
  (`task`, `reason`, `actor`, `recorded_at`), `exception` (`requirements`,
  `reason`, `owner`, `recorded_at`, `check`), `owner_validation`
  (`required`, `accepted_check`, `accepted_by`).
- `savepoint-design/SKILL.md:179` shows only `owner_validation: {required: true}`.
- Lower risk than I-077: the board writes waivers itself
  (`internal/data/write.go:842`), and T-010's board-written waiver is
  well-formed. The gap bites when an agent records an owner decision by hand.

## Proof Needed

- The skill that records each field shows its exact shape once, from the
  runtime field names.
- Records written from those examples load and satisfy the matching gates.
- Packaged template copies stay byte-identical.

## Repair Attempt Evidence

- Shapes match the decoders: waiver actor role owner and `task` equal to its own ID; exception requires `check`; acceptance requires `accepted_check` and an owner `accepted_by`.
- An `exception:` block written from the example, added to a copy of T-010 alongside its board-written waiver, loads under `savepoint resume`.
- Packaged copies are byte-identical to `agent-skills/`; `make build` and `make test-fast` passed.
