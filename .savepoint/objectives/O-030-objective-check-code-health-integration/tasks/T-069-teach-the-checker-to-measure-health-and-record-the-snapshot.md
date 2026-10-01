---
id: T-069
title: Teach the checker to measure health and record the snapshot
objective: O-030
status: planned
depends_on: [{task: T-067, requires: clear}, {task: T-068, requires: clear}]
owner_validation: {required: true}
planned_by: {role: planner, session: planning-o030-20261001}
complexity_tier: medium
complexity_reason: Workflow wording across mirrored skill, method, and guide files that must stay byte-identical and match shipped behavior.
---

# Teach the checker to measure health and record the snapshot

## Outcome

A fresh checker running a Full Objective Check knows to run the full gate, then `savepoint health check`, record the snapshot ID in the Check, and treat the verdict as supporting evidence. A blocking verdict prevents CLEAR. Warnings lead to Issues only by the checker's own judgment. Task Checks never collect.

## User Check

Read the updated Full Objective Check steps in `agent-skills/savepoint-check/SKILL.md` and the CLI Rules exception in `AGENTS.md`. Confirm they say, in plain words: collection happens only in a Full Objective Check; the full gate runs first, so tests aren't run twice; what blocks clearance; that a manual snapshot never counts; and that no Issue is created automatically.

## Done When

- `savepoint-check` SKILL.md and `references/check-method.md`, Full evidence only, say:
  - run the full gate first, then `savepoint health check O-### [dir]`;
  - record `health_snapshot` in the Check frontmatter;
  - a verdict that blocks clearance prevents `CLEAR`;
  - a project with no health configuration records "Code Health not configured" and is not a finding;
  - optional failures and warnings go through the checker's existing Issue-capture judgment;
  - Task Checks and other activity never run collection.
- The Check record template in the skill shows the optional `health_snapshot` field.
- `AGENTS.md` and `templates/project-v2/AGENTS.md` CLI Rules add the narrow exception for `savepoint health check`, limited to Full Objective Checks. `health setup` stays human-only.
- Design.md section 6 lists the command. Section 7 describes Code Health as supporting evidence in the Full Objective Check, and section 1 no longer says the collection command is still to come. The Codebase Map entries for `internal/codehealth` and `internal/doctor` reflect the gate and the snapshot-reference diagnostic.
- Live and template skill and reference copies are byte-identical (TPL-01). `agent_skills_test.go` and the template tests pass.

## Context Files

`.savepoint/objectives/O-030-objective-check-code-health-integration/Objective.md`; `agent-skills/savepoint-check/SKILL.md`, `agent-skills/references/check-method.md`; `templates/project-v2/agent-skills/savepoint-check/SKILL.md`, `templates/project-v2/agent-skills/references/check-method.md`; `AGENTS.md`, `templates/project-v2/AGENTS.md`; `.savepoint/Design.md`; `agent_skills_test.go`.

## Design References

Design sections 1, 6, and 7; O-030 Confirmed Design Decisions.

## Guardrails

TPL-01, TPL-02, TPL-04, POL-01, POL-02, TEST-01, TEST-06, TEST-08.

## Implementation Plan

1. Confirm the command and the `health_snapshot` field behave as their Tasks recorded; return REPLAN REQUIRED if the shipped behavior differs.
2. Edit the live skill and method, then copy them byte-for-byte to the template tree.
3. Update both AGENTS.md CLI Rules sections and Design.md.
4. Run the skill-mirror tests and the fast gate.

## Boundaries

No code changes beyond test expectations that assert guidance text. No Task Check changes, and no new blocking rule beyond those T-066 implements.

## Technical Verification

`make build && make test-fast` at handoff. The mandatory Full Objective Check uses fresh `make test-full`, following `agent-skills/references/check-method.md`.

## Technical Evidence

Pending execution.

## Drift Notes

Existing projects receive the refreshed skills through `upgrade-assets` (TPL-04); no new scaffold file is added.
