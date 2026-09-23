---
id: T-012
title: Route planners through Task creation
objective: O-019
planned_by: {role: planner, session: i024-allocation-design-20260923}
status: planned
complexity_tier: medium
complexity_reason: Canonical and scaffold guidance, upgrades, and agent scenarios must agree on one required path.
depends_on: [{task: T-011, requires: clear}]
owner_validation: {required: true}
---

# T-012: Route planners through Task creation

## Outcome

Planners use the creation operation without selecting an ID. Guidance permits that one command while retaining the ban on other agent-run `savepoint` commands.

## User Check

Review the planning instructions and a two-Objective example. Confirm the planner supplies an ID-free draft, never chooses a number, and may invoke only the narrow Task creation command.

## Done When

- `AGENTS.md` permits `savepoint create-task` only for new Tasks and continues to prohibit other agent-run `savepoint` commands.
- Canonical design skill and byte-identical V2 scaffold copy require the operation; no example directs manual Task ID allocation.
- Fresh-init and upgrade-assets guidance deliver the rule. Strict V2 loading is required after creating or renaming other identity-bearing records.
- A cross-Objective workflow scenario creates Tasks concurrently and fails on ID reuse or an invalid index.
- Command documentation and Design describe implemented behavior and durable allocation state.

## Context Files

`AGENTS.md`, `agent-skills/savepoint-design/SKILL.md`, `templates/project-v2/agent-skills/savepoint-design/SKILL.md`, `agent-skills/savepoint-check/SKILL.md`, `templates/project-v2/agent-skills/savepoint-check/SKILL.md`, `agent-skills/references/issue-capture.md`, `templates/project-v2/agent-skills/references/issue-capture.md`, `internal/init/agent_skills_test.go`, `internal/init/template_freshness_test.go`, `internal/init/v2_scaffold_test.go`, `.savepoint/Design.md`, `README.md`.

## Design References

Design sections 1, 2, 5, and 6; O-019 Architectural Considerations.

## Guardrails

TPL-01, TPL-02, TPL-04, ARCH-03, TEST-01, TEST-02, TEST-04, TEST-05, TEST-08.

## Implementation Plan

1. Replace manual Task-ID guidance with ID-free draft and command instructions. Document the narrow agent exception and its upgrade path.
2. Align canonical and scaffold skills; reconcile Check/Issue guidance only where identity-bearing validation needs it.
3. Add cross-Objective agent scenarios and scaffold/upgrade assertions.
4. Reconcile Design and docs, then run focused tests, `git diff --check`, and `make build && make test`.

## Boundaries

No permission for unrelated CLI commands or duplicated allocation policy across skills.

## Technical Verification

Focused `internal/init` and workflow scenarios, skill byte comparison, `git diff --check`, and `make build && make test`.

## Technical Evidence

Pending execution: per-criterion outcomes, command results, files read/changed, and limitations.

## Drift Notes

If the managed `AGENTS.md` merge cannot safely deliver the exception to existing projects, return REPLAN REQUIRED.
