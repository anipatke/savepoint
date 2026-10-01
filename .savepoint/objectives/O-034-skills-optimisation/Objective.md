---
id: O-034
title: Skills optimisation
status: planned
depends_on: [O-032]
release: R-007
priority: low
rank: 1
---

# O-034: Skills optimisation

## Outcome

At the end of each Goal, review Savepoint's skill templates against lessons from the completed work, identify opportunities to improve them, and apply justified tweaks that make the next Goal's workflow clearer and more efficient.

## Why

Planning, execution, and independent Checks reveal repeated instructions, unnecessary context reads, ambiguous handoffs, and guidance that can be improved. A recurring review turns those lessons into better templates rather than leaving them in individual sessions.

## Success Conditions

- The current Goal finishes with a review of the four canonical skills, their shared references, and the shipped template copies, after its feature and hardening Objectives.
- Review findings cite evidence from the Goal's work and assess clarity, duplication, context and token cost, routing, handoffs, and consistency with implemented behavior.
- Useful improvements are applied and verified; opportunities that need separate product or architecture decisions are recorded as follow-up work. A review may conclude that no change is justified, with its reasons recorded.
- Optimisation preserves owner authority, checker independence, acceptance evidence, mandatory Objective Checks, and project-owned verification policy.
- Canonical skills and their scaffold copies remain byte-identical, and existing projects can receive improvements through the supported asset upgrade path.
- Shipped planning guidance establishes the same review at the end of future Goals, with a clear owner and recorded outcome, without adding a Goal-level technical Check or changing Goal completion semantics.
- Representative workflow scenarios and the required gates verify that revised instructions retain correct behavior while reducing demonstrated friction.

## Architectural Considerations

Skill instructions remain the canonical workflow source; shared references hold shared procedures, and AGENTS.md owns routing and repository rules. Changes must respect those boundaries and avoid duplicating policy. Goal membership continues to derive from Objective records.

The exact recurring review mechanism and scenario verification approach will be settled when this Objective is ready for Task planning. Detailed Tasks wait until the preceding work has supplied its lessons.

## Boundaries

**In scope:** an end-of-Goal skill retrospective, evidence-backed template improvements, canonical/scaffold parity, upgrade delivery, representative agent scenarios, and guidance making the review recurring for future Goals.

**Out of scope:** unrelated feature work, speculative rewrites, weakening verification to save tokens, new lifecycle states, automatic Goal completion, a Goal-level Check, and general-purpose agent orchestration.

## Confirmed Direction

The owner requested one final Objective, titled Skills optimisation, and an opportunity to revisit and improve skill templates at the end of every Goal on 2026-10-02. This Objective is the final planned member of R-007 and waits for O-032, which depends on the Goal's preceding Code Health Objectives.
