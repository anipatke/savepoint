---
id: T-126
title: Write the end-of-Goal review as a note, not an Objective
objective: O-042
status: done
depends_on: []
complexity_tier: low
complexity_reason: One skill section and its scaffold copy, plus the two tests that pin its wording.
owner_validation:
    required: true
    accepted_check: ""
planned_by: {role: planner, session: plan-o042-20261010}
check_waiver:
    task: T-126
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-10T08:16:17Z"
---

# Write the end-of-Goal review as a note, not an Objective

## Outcome

`savepoint-design` tells the planner to write the end-of-Goal workflow review as a `## Workflow Review` section in the Goal's `Release.md`, with no Objective, Task or Check, and to turn any real change it finds into an ordinary Objective or Issue.

## User Check

Open `agent-skills/savepoint-design/SKILL.md` and read the review section: it should ask for a note on the Goal, say a "no change" review needs no Objective, Task or Check, and no longer say "add one final workflow-retrospective Objective". Compare it with the G-002 note in `.savepoint/releases/G-002-skills-optimisation/Release.md`.

## Done When

- The `### Goal Workflow Retrospective` section is replaced by a `### Goal Workflow Review` section that: says when the planner writes it (the first planning session after every member Objective of the Goal is done, or when the owner asks); names the `## Workflow Review` section in the Goal's `Release.md`; keeps the existing review scope and evidence list; keeps the "no change, because…" conclusion; says it is not an Objective, Task or Check and that a real change becomes an Objective or an Issue; keeps the package-receiving rule (tune own Guardrails, AGENTS.md project rules and gates; record packaged-skill suggestions as Issues); and keeps "a Goal is still complete when every member Objective is complete".
- The scaffold copy in `templates/project-v2/agent-skills/savepoint-design/SKILL.md` is byte-identical (TPL-01).
- `internal/init/agent_skills_test.go` pins the new wording and lists "retrospective Objective" among forbidden phrases; `internal/init/upgrade_test.go`'s upgrade test asserts the new heading and is renamed to match.
- No other skill, reference or AGENTS.md text mentions a required retrospective Objective (confirm with a search of `agent-skills/`, `templates/project-v2/` and `AGENTS.md`).
- `make build && make test-fast` pass.

## Context Files

`agent-skills/savepoint-design/SKILL.md`, `templates/project-v2/agent-skills/savepoint-design/SKILL.md`, `internal/init/agent_skills_test.go`, `internal/init/upgrade_test.go`, `.savepoint/releases/G-002-skills-optimisation/Release.md`, `.savepoint/objectives/O-042-skills-optimisation-workflow-retrospective/Objective.md`.

## Design References

O-042 Confirmed Decisions.

## Guardrails

TPL-01, TPL-02, TPL-04.

## Implementation Plan

1. Rewrite the review section in the canonical design skill per Done When; keep it about the same length.
2. Copy it byte-for-byte to the scaffold.
3. Update the content-test row and the upgrade test.
4. Search for leftover "retrospective Objective" wording in the skills, references, scaffold and AGENTS.md; fix any that describe the old rule.
5. Run the handoff gates.

## Boundaries

Guidance and its tests only. No change to Objective closing, the Goal loader, fields, states or commands. Leave past records such as O-038 as they are.

## Technical Verification

Focused `make test-focused TEST=...` on the init package while iterating; handoff `make build && make test-fast`; the Full Objective Check under `agent-skills/references/check-method.md` with fresh `make test-full`. No platform-sensitive code, so the normal CI run is enough.

## Technical Evidence

Executor session, 2026-10-10. No extra reads beyond the Context Files, plus a search of `agent-skills/`, `templates/project-v2/`, `AGENTS.md` and `internal/init/*_test.go` for "retrospective" (a Done When step).

Per criterion:
- Section replaced: `### Goal Workflow Review` in `agent-skills/savepoint-design/SKILL.md` states when it is written (first planning session after every member Objective is done, or on owner request), names `## Workflow Review` in the Goal's `Release.md`, keeps scope and evidence list, keeps "no change, because…", says it is not an Objective, Task or Check and that a real change becomes an Objective or Issue, keeps the package-receiving rule and the "Goal is still complete when every member Objective is complete" clause. Met.
- Scaffold copy byte-identical: `cmp` of the two files reported identical. Met (TPL-01).
- Tests: `agent_skills_test.go` row renamed and pins the new wording, forbids "retrospective Objective" and "workflow-retrospective Objective"; `upgrade_test.go` test renamed `TestUpgradeDeliversGoalWorkflowReviewRuleToDesignSkill` and asserts `### Goal Workflow Review`. Met.
- No other mention: search found "retrospective" only in the forbidden-phrase row of `agent_skills_test.go`. Met.
- `make build && make test-fast`: both exit 0. Met.

Files changed: both design skill copies, `internal/init/agent_skills_test.go`, `internal/init/upgrade_test.go`; lifecycle fields in this Task and O-042.

Limitations: `make test-full` not run (Full Objective Check will); no Task Check requested or waived.

## Drift Notes

None yet.
