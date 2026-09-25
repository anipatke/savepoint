---
id: T-048
title: Stop describing a Goal Check in guidance
objective: O-025
status: planned
depends_on: [{task: T-047, requires: clear}]
owner_validation: {required: false}
planned_by: {role: planner, session: o025-plan-20260925}
---

# Stop describing a Goal Check in guidance

## Outcome

Every active document and the new-project scaffold describe two Checks: the
mandatory Full Objective Check and the optional Task Check. A Goal is
complete when all its Objectives are complete.

## User Check

Search the active guidance for "Goal Check", "Release Check", and
`scope.kind: release`: only the compatibility note that old
`scope.kind: release` records still load remains. Read AGENTS.md's
Verification Policy and Check sections: they name only the Objective and Task
Checks.

## Done When

- `.savepoint/Design.md` describes Goal completion as all member Objectives
  complete, removes the Goal Check step from the verification order, and
  notes that `scope.kind: release` records still load but no decision reads
  them.
- `.savepoint/Guardrails.md` TEST-08, TEST-09, and the Check summary no
  longer mention a Goal Check.
- `AGENTS.md`, `.savepoint/router.md`, and `README.md` drop the Goal Check.
- `agent-skills/savepoint-check`, `savepoint-design`, `savepoint-task`, and
  the three shared references drop the Goal Check, including savepoint-design's
  Goal `done` rule.
- The `templates/project-v2` copies of every file above match the live
  wording where they are meant to.
- `internal/init/template_freshness_test.go`,
  `internal/init/agent_skills_test.go`, and migration goldens are updated
  only where their expected text changes.
- `make build && make test-fast` pass.

## Context Files

`.savepoint/Design.md`, `.savepoint/Guardrails.md`, `.savepoint/router.md`,
`AGENTS.md`, `README.md`,
`agent-skills/savepoint-check/SKILL.md`,
`agent-skills/savepoint-design/SKILL.md`,
`agent-skills/savepoint-task/SKILL.md`,
`agent-skills/references/check-method.md`,
`agent-skills/references/commands-and-procedures.md`,
`agent-skills/references/issue-capture.md`,
`templates/project-v2/.savepoint/Design.md`,
`templates/project-v2/.savepoint/Guardrails.md`,
`templates/project-v2/.savepoint/router.md`,
`templates/project-v2/AGENTS.md`,
`templates/project-v2/agent-skills/savepoint-check/SKILL.md`,
`templates/project-v2/agent-skills/savepoint-design/SKILL.md`,
`templates/project-v2/agent-skills/savepoint-task/SKILL.md`,
`templates/project-v2/agent-skills/references/check-method.md`,
`templates/project-v2/agent-skills/references/commands-and-procedures.md`,
`templates/project-v2/agent-skills/references/issue-capture.md`,
`internal/init/template_freshness_test.go`,
`internal/init/agent_skills_test.go`,
`.savepoint/objectives/O-025-remove-the-goal-check/Objective.md`.

## Design References

Design Section 1 (Goal boundary, Check workflow), the Check table, and
Section 10 verification order.

## Guardrails

TEST-08, TEST-09, TPL-02, TEST-01..04.

## Implementation Plan

1. Confirm the runtime Task landed; return REPLAN REQUIRED if Goal
   completion still reads a Goal Check.
2. Edit the live guidance files, then mirror the scaffold copies.
3. Update template and skill tests and any migration goldens whose text
   changed.
4. Search for leftover Goal Check wording in active guidance.
5. Run `make build && make test-fast`.

## Boundaries

No runtime code changes, and no edits to historical Checks, Issues,
Objectives, or Tasks.

## Technical Verification

Focused template and skill tests during iteration; `make build &&
make test-fast` at handoff; the Full Objective Check requires
`make test-full`.

## Technical Evidence

Pending execution.

## Drift Notes

None expected.
