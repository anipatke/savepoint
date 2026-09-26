---
id: I-082
title: Every skill repeats Goal context, and two carry V1 and keybinding history
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
      The four skills now carry a role-sized Goal paragraph and a one-line pointer to AGENTS.md for Next-line handling; V1 conversion detail and the r keybinding alias are gone from packaged skills. Packaged template copies re-synced; make build and make
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

# I-082: Every skill repeats Goal context, and two carry V1 and keybinding history

## Summary

The same Goal-context paragraph opens all four workflow skills, and
`savepoint-idea` and `savepoint-design` add V1 Release conversion and board
keybinding detail an agent does not need for its role. The Next-line
paragraph also repeats in every skill. This is prompt weight, not a behavior
defect.

## Evidence

- Goal paragraph: `savepoint-idea/SKILL.md:12`, `savepoint-task/SKILL.md:14`,
  `savepoint-check/SKILL.md:14`, and `savepoint-design/SKILL.md` Goal section.
- V1/keybinding detail: `savepoint-design/SKILL.md:97-104` ("`r` is an
  undisplayed compatibility alias") and the matching text in
  `savepoint-idea/SKILL.md:12`.
- Next-line paragraph: `savepoint-idea:20`, `savepoint-design:20`,
  `savepoint-task:22`, `savepoint-check:25`.

## Proof Needed

- Each skill keeps only the Goal facts its role acts on; V1 conversion and
  keybinding detail move to migration or board guidance.
- Router behavior and stale-reference tests still pass; packaged copies stay
  in parity.

## Repair Attempt Evidence

- Goal paragraph keeps every phrase `TestProjectGuidanceRequiresGoalContext` and `TestIdeaGuidanceFillsTheFreshProjectsGoal` pin.
- task and check now tell the agent to report a missing Goal to the owner rather than select one.
- Size: idea 5882→4788 bytes, design 19572→18790, check 13210→12738; task grew slightly (12197→12279) from the I-077/I-079 examples.
- Packaged copies are byte-identical to `agent-skills/`; `make build` and `make test-fast` passed.
