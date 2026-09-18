---
id: E46-agent-workflow-assets/T003-plan-one-objective-at-a-time
title: Plan one objective at a time
status: done
objective: Add the savepoint-design planner skill, its Objective template, and its readiness rules, so only the next Objective is ever detailed.
depends_on:
    - E46-agent-workflow-assets/T002-take-a-rough-idea-without-demanding-a-document
complexity_tier: medium
complexity_reason: One mirrored skill pair, but it carries the Objective template, readiness gate, and research-task rule the rest of planning depends
---

# T003: Plan one objective at a time

## Problem

V1 planning produces a whole release of epics up front, then detail for each in turn. V2 collapses that: the planner owns Design, Guardrails, and one Objective, and it details Tasks only for the Objective that is next. Everything further out stays a named outcome with boundaries, because detailed plans written ahead of the decisions that shape them are the thing that gets thrown away.

That makes readiness a real gate rather than a feeling. The planner may only detail an Objective's Tasks when interfaces and data ownership are settled, constraints are scoped, dependency outcomes are known, and there is a verification approach. When the implementation approach itself is unknown, the honest output is a bounded research Task with a decision deliverable — not a confident plan that the executor discovers is fiction.

The skill also has to say what the planner is not: it does not write production code, and it does not treat owner product choices as technical decisions it can settle. Technical readiness does not require the owner to review code.

## Context Files

- `agent-skills/savepoint-design/SKILL.md`
- `templates/project/agent-skills/savepoint-design/SKILL.md`
- `agent-skills/savepoint-system-design/SKILL.md`
- `agent-skills/savepoint-create-plan/SKILL.md`
- `agent-skills/savepoint-create-task/SKILL.md`
- `internal/init/skill_validation_test.go`
- `internal/init/agent_skills_test.go`
- `.savepoint/releases/v2/v2-Design.md`
- `.savepoint/Guardrails.md`

## Acceptance Criteria

- [x] `agent-skills/savepoint-design/SKILL.md` exists with frontmatter `name: savepoint-design`, and the shipped copy under `templates/project/agent-skills/` is byte-identical.
- [x] It passes the existing structure validation: non-empty `## Purpose`, `## Trigger`, `## Read`, `## Workflow`, `## Rules`.
- [x] Its trigger is router `state: design`, and it states that `REPLAN REQUIRED` from an executor routes back into this skill rather than into a separate phase.
- [x] Its write boundary names Design, Guardrails, the current Objective, the next Objective's detailed Tasks, and routing — and forbids production code and detailed backlog beyond the next Objective.
- [x] It carries the Objective artifact template: frontmatter `id`, `title`, `status: planned|in_progress|done`, `depends_on: [O###]`, optional `release`, `last_check`, and evidence/freshness fields; body sections Outcome, Why, Success Conditions, Architectural Considerations, Boundaries.
- [x] The Objective template states that Task membership is derived from Task ownership and must not be maintained as a second list in the Objective body.
- [x] It carries the Design template sections — Architecture, Components/Codebase Map, Interfaces and Data Flow, Boundaries, Decisions, Current Technical State — and states that Design describes implemented reality while Objective deltas hold planned change until reconciliation.
- [x] It states the Guardrails rule: durable project constraints only, roughly 10–20 substantive rules with stable category IDs, severity and exception authority owned by policy, and no duplicated prose in tasks or checks.
- [x] It states the readiness gate for detailing an Objective's Tasks: settled interfaces and data ownership, scoped constraints, known dependency outcomes, and a verification approach.
- [x] It states that an unknown implementation approach becomes a bounded research Task with a named decision deliverable, and that a Task with multiple unrelated outcomes or an unresolved architectural decision is split.
- [x] It states that product choices go to the owner and that technical readiness does not require owner code review.
- [x] `internal/init/agent_skills_test.go` asserts the write boundary, the Objective template fields and body sections, the readiness gate, the one-Objective rule, and live/template parity.
- [x] The V1 planning skills `savepoint-create-plan`, `savepoint-system-design`, and `savepoint-create-task` are unchanged and still present in both trees.

## Implementation Plan

- [x] Read design sections 6, 7, and 8 for the Objective contract, planning readiness, and the planner role row.
- [x] Write `agent-skills/savepoint-design/SKILL.md` with Purpose, Trigger, Read, Workflow, Rules, the Objective artifact template, and the Design/Guardrails section contracts.
- [x] State the one-Objective rule, the readiness gate, the research-Task rule, and the split rule as numbered rules the contract test can assert.
- [x] Reference `agent-skills/references/check-method.md` for how planned verification is later evaluated rather than restating it.
- [x] Copy the file verbatim to the template tree.
- [x] Extend `internal/init/agent_skills_test.go` with the savepoint-design contract case.
- [x] Run `go test ./internal/init/...`, then `make build && make test`.

## Context Log

**Read:** `agent-skills/savepoint-system-design/SKILL.md`, `agent-skills/savepoint-create-plan/SKILL.md`, `agent-skills/savepoint-create-task/SKILL.md`, `internal/init/skill_validation_test.go`, `internal/init/agent_skills_test.go`, `.savepoint/Guardrails.md`, `.savepoint/releases/v2/v2-Design.md` (sections 5-8), `agent-skills/savepoint-idea/SKILL.md` (structural precedent for the mirrored skill pair), `agent-skills/references/check-method.md` (to reference, not restate).

**Edited:**
- `agent-skills/savepoint-design/SKILL.md` — new skill: Purpose, Trigger (`state: design`, `REPLAN REQUIRED` re-entry), Read, Workflow, Objective artifact template, Design template, Guardrails rule, readiness gate, Rules (one-Objective, research-Task, split, product-choice/technical-readiness).
- `templates/project/agent-skills/savepoint-design/SKILL.md` — byte-identical copy.
- `internal/init/agent_skills_test.go` — added `TestSavepointDesignSkillPassesStructureValidation`, `TestSavepointDesignSkillTrigger`, `TestSavepointDesignSkillWriteBoundary`, `TestSavepointDesignSkillObjectiveTemplate`, `TestSavepointDesignSkillDesignTemplate`, `TestSavepointDesignSkillGuardrailsRule`, `TestSavepointDesignSkillReadinessGate`, `TestSavepointDesignSkillOneObjectiveResearchAndSplitRules`, `TestSavepointDesignSkillLiveAndTemplateMatch`.

**Quality gates:**
- `go test ./internal/init/... -run 'Skill|SavepointDesign' -v` — all new and existing skill/init tests pass.
- `make build && make test` — build succeeds, full suite passes across all packages.
- `git status --short` on `agent-skills/savepoint-create-plan`, `agent-skills/savepoint-system-design`, `agent-skills/savepoint-create-task` and their template copies — no output, confirming the V1 planning skills are untouched.

No `.savepoint/Health-Check.md` in this project; Quick-check step skipped per its documented absence-is-not-a-finding rule.
