---
id: O-034
title: Skills optimisation
status: in_progress
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

The recurring review mechanism and scenario verification approach are settled in Confirmed Design below.

## Confirmed Design

Owner-confirmed on 2026-10-03 (planning session plan-o034-20261003), after O-032 closed CLEAR (C-960):

- **Review scope.** The four workflow skills (`savepoint-idea`, `savepoint-design`, `savepoint-task`, `savepoint-check`), the three shared references, AGENTS.md routing guidance, and the scaffolded project documents (template `Design.md`, `Guardrails.md`, `Idea.md`, `config.yml`, `router.md`, Goal `Release.md`, `AGENTS.md`; extended by the owner on 2026-10-03). `bubbletea-tui-design` is out of scope. Evidence comes from R-007's records: REPLAN REQUIRED Tasks, NEEDS WORK Checks and their rechecks, Issues I-084..I-121, and Lessons Carried In.
- **Recurrence.** Shipped `savepoint-design` guidance tells the planner to add a final workflow-retrospective Objective to each Goal once its other Objectives are planned. The planner owns it, and its outcome is recorded in that Objective. There are no new fields, states, commands, or Goal-level Check, and Goal completion semantics are unchanged.
- **Downstream behaviour.** In a project that receives the skills from the package, the retrospective tunes project-owned files (Guardrails, project rules in AGENTS.md, configured gates) and records suggestions for the packaged skills as Issues, rather than editing package-owned skills. In this repository the review edits the canonical skills directly.
- **Scenario verification.** Written scenario walkthroughs are checked against the revised text (at least: Next → skill routing, REPLAN REQUIRED re-entry, a waived Task Check reaching the Full Objective Check, an issue-only repair, and acting on an advisory Code Health report). Go content tests pin the key rules and live/scaffold byte parity. Live agent runs are not required.
- **Code Health advice.** `savepoint-task` and `savepoint-check` gain a short rule for acting on an advisory Code Health report: fix to the watch line, not the aim; prefer production code that is risky or often changed; leave flat dispatch tables alone; report what remains; treat narrowing the measured scope as an owner decision. A change to the generated report's next-step wording is product work, recorded as a follow-up Issue rather than applied here.

## Lessons Carried In

Evidence from O-032's Code Health complexity work (2026-10-03), for the review to weigh when it runs:

- The complexity signal has an aim (10 or less) and a watch line (20 or less). Bringing the hardest function from 42 to 20 took about 25 functions, mostly long scenario tests; reaching the aim would touch about 280 more, over 120 of them production code. The report's "investigate each signal and apply a fix" wording invites chasing the aim without saying where to stop.
- Costs seen: indirection from splitting, scenario tests that no longer read top to bottom, refactor regressions (tests changed alongside the code that they guard), review churn, and metric gaming, since moving branches elsewhere lowers the score without removing them.
- Judgment worth stating in skills: fix to the watch line, prefer production code that is risky or often changed, leave flat dispatch tables alone, report what remains rather than silently stopping or silently continuing, and treat narrowing the measured scope (for example excluding test files) as an owner decision.
- Candidate review questions: whether `savepoint-task` and `savepoint-check` should say how an agent acts on an advisory Code Health report, and whether the generated report's next-step text should name the watch line as the stopping point. Any change to report wording is product work and would be recorded as follow-up, not applied by this review.

## Boundaries

**In scope:** an end-of-Goal skill retrospective, evidence-backed template improvements, canonical/scaffold parity, upgrade delivery, representative agent scenarios, and guidance making the review recurring for future Goals.

**Out of scope:** unrelated feature work, speculative rewrites, weakening verification to save tokens, new lifecycle states, automatic Goal completion, a Goal-level Check, and general-purpose agent orchestration.

## Confirmed Direction

The owner requested one final Objective, titled Skills optimisation, and an opportunity to revisit and improve skill templates at the end of every Goal on 2026-10-02. This Objective is the final planned member of R-007 and waits for O-032, which depends on the Goal's preceding Code Health Objectives.
