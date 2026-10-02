---
id: T-101
title: Make the skills review happen at the end of every Goal
objective: O-034
status: planned
depends_on: [{task: T-100, requires: clear}]
complexity_tier: low
complexity_reason: One planning-guidance rule with parity, upgrade, and content tests.
owner_validation: {required: true}
planned_by: {role: planner, session: plan-o034-20261003}
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

Pending execution.

## Drift Notes

None expected.
