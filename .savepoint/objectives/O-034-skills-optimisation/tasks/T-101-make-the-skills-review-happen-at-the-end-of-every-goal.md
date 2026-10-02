---
id: T-101
title: Make the skills review happen at the end of every Goal
objective: O-034
status: done
depends_on: [{task: T-100, requires: clear}]
complexity_tier: low
complexity_reason: One planning-guidance rule with parity, upgrade, and content tests.
owner_validation:
    required: true
    accepted_check: ""
planned_by: {role: planner, session: plan-o034-20261003}
check_waiver:
    task: T-101
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-02T20:18:06Z"
---

# Make the skills review happen at the end of every Goal

## Outcome

The shipped `savepoint-design` skill tells the planner that, once a Goal's other Objectives are planned, it adds a final workflow-retrospective Objective. That Objective has a clear owner and a recorded outcome, and it says what the review may change in a project that receives the skills from the package.

## User Check

Read the new rule in `savepoint-design`. It says when the retrospective Objective is added, who owns it, and where its outcome is recorded. It also says a downstream project tunes its own Guardrails, AGENTS.md project rules, and gates, and records suggestions for the packaged skills as Issues. Nothing in it adds a Goal Check or changes when a Goal is complete.

## Done When

- `savepoint-design` states the rule in its Goal guidance: when to add the retrospective, planner ownership, the outcome recorded in that Objective (including a "no change, because…" conclusion), and the downstream rule. These match O-034's Confirmed Design.
- The rule adds no field, state, command, Goal-owned Task list, or Goal-level Check, and it leaves "a Goal is complete when every member Objective is complete" unchanged. A content test asserts the rule and the absence of a Goal Check.
- Live and scaffold copies are byte-identical (TPL-01). An upgrade test shows an existing project receives the rule through `upgrade-assets` (TPL-04, FS-02).
- The scenario "planner finishes planning a Goal's Objectives" is walked against the revised text and recorded.
- `make build && make test-fast` passes.

## Context Files

`agent-skills/savepoint-design/SKILL.md`; `templates/project-v2/agent-skills/savepoint-design/SKILL.md`; `internal/init/agent_skills_test.go`; `internal/init/upgrade_test.go`; `.savepoint/objectives/O-034-skills-optimisation/Objective.md`.

## Design References

O-034 Confirmed Design (Recurrence, Downstream behaviour); AGENTS.md Required Goal Context.

## Guardrails

TPL-01, TPL-02, TPL-04, FS-02, POL-02, TEST-01, TEST-02, TEST-08.

## Implementation Plan

1. Add the rule to the Goal guidance in `savepoint-design` and mirror it to the scaffold.
2. Add the content test and the upgrade assertion.
3. Walk the scenario and run the gate.

## Boundaries

Guidance only. No Go behaviour, Goal template, or completion-semantics change.

## Technical Verification

Focused `make test-focused TEST=... PKGS=./internal/init` during iteration; `make build && make test-fast` at handoff. The later independent Check applies `agent-skills/references/check-method.md`.

## Technical Evidence

Executor evidence (not a Check, not CLEAR).

### Per-criterion outcomes

- **Rule stated.** `savepoint-design` gains "### Goal Workflow Retrospective" under Required Goal Context: add one final retrospective Objective once the Goal's other Objectives are planned; the planner owns it and records the outcome in it, including a "no change, because…" conclusion; downstream projects tune their own Guardrails, AGENTS.md project rules and gates and record packaged-skill suggestions as Issues. Matches O-034's Confirmed Design.
- **No new mechanism.** The text adds no field, state, command, Goal-owned Task list or Goal-level Check, and the existing "a Goal is complete when every member Objective is complete" lines are untouched. `TestSkillReviewRulesArePinned` ("design adds a retrospective Objective to each Goal") pins the rule and fails if "Goal Check" or "Goal-level Check" appears.
- **Parity (TPL-01).** The scaffold copy is a byte copy; `TestSavepointDesignSkillLiveAndTemplateMatch` passes.
- **Upgrade (TPL-04, FS-02).** `TestUpgradeDeliversGoalRetrospectiveRuleToDesignSkill`: a manifest-tracked old design skill is updated from the real templates and contains the new section.
- **Scenario "planner finishes planning a Goal's Objectives".** Before, the text said nothing, so no review happened. Now the planner adds the final Objective, owns it, and records its outcome there. In this repo the review edits canonical skills; downstream it tunes project-owned files and files Issues. Completion semantics are unchanged.
- **Gate.** `make build && make test-fast` passed (exit 0).

### Limitations

- The scenario walk is written, not a live agent run. No `make test-full`: no code changed.
- Extra reads, logged: `internal/init/upgrade_test.go` helpers and `agent_skills_test.go` case table, to match T-100's test pattern.

## Drift Notes

None expected.
