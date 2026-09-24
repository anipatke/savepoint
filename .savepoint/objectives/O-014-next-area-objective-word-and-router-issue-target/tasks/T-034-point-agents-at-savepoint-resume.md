---
id: T-034
title: Point agents at savepoint resume for the next step
objective: O-014
planned_by: {role: planner, session: o014-design-20260924}
status: planned
complexity_tier: medium
complexity_reason: "Guidance-only, but spans AGENTS.md, the router template, four live skills with byte-identical scaffold copies, and Design reconciliation for the whole Objective."
depends_on: [{task: T-028, requires: clear}, {task: T-033, requires: clear}]
owner_validation: {required: false}
---

# T-034: Point agents at savepoint resume for the next step

## Outcome

Agent guidance names `savepoint resume` as the one CLI command agents may
run and the source of the next action, says who updates the router selection
after a closure, and Design.md describes what O-014 implemented.

## User Check

Paste a Next line such as `In Progress O-014 · Build T-029 — ...` into a
fresh session: the agent opens T-029 with `savepoint-task` without asking
what the line means.

## Done When

- AGENTS.md (live and `templates/project-v2/AGENTS.md`), the router
  template, and the phase skills replace "act on the router's next_action"
  with "run `savepoint resume` and act on its Next". CLI Rules allow exactly
  `savepoint resume` (read-only) and keep every other command human-only.
- AGENTS.md's Workflow starts with the Next line: a pasted or `resume`-printed
  line names the Objective and Task; its word picks the skill (Task
  `Planned`/`Build`/`Test` → `savepoint-task`, Task or Objective `Check` →
  `savepoint-check`, Objective with no Task → `savepoint-design`; an
  `Open` or `In Progress` Issue → repair under `savepoint-task` per
  `issue-capture.md`, leaving closure to a checker or the owner).
- AGENTS.md documents the owner request "set router to O-### [T-###]
  [I-###]": confirm each record exists and the Task belongs to the
  Objective (refuse and say why otherwise), edit only those selection keys,
  leave `release:` unchanged unless the owner names a Release, and finish by
  running `savepoint resume` and showing its Next line.
- When the binary is unavailable, guidance says to read the router
  selection, report the missing tool, and not guess the next step.
- Router-advance ownership is stated once in AGENTS.md and referenced by
  skills: the board advances the router when the owner closes a Task; `savepoint-design`
  selects the next Objective; `savepoint-task` selects the Task it starts;
  nobody blanks or changes `release:` except an explicit Goal choice.
- Skills stop writing `next_action` (`savepoint-design` step 11,
  `savepoint-idea`). Live and scaffold skills stay byte-identical.
- Design.md section 8 states the Next line format, that Next is the router
  selection, the board's closure advance, the diagnostic line, and Issue
  context; sections 1 and 4 drop `next_action` and outdated `p` hand-off
  wording. Section 4's Issue sentence says Issues carry no stage (repairs
  I-043; record `repair_attempted` in the Issue, leave closure to a checker).
- Template freshness and agent-skill tests pass; `git diff --check`,
  `make build && make test-fast` pass.

## Context Files

`AGENTS.md`, `templates/project-v2/AGENTS.md`,
`templates/project-v2/.savepoint/router.md`, `.savepoint/router.md`,
`agent-skills/savepoint-idea/SKILL.md`,
`agent-skills/savepoint-design/SKILL.md`,
`agent-skills/savepoint-task/SKILL.md`,
`agent-skills/savepoint-check/SKILL.md`,
`templates/project-v2/agent-skills/savepoint-idea/SKILL.md`,
`templates/project-v2/agent-skills/savepoint-design/SKILL.md`,
`templates/project-v2/agent-skills/savepoint-task/SKILL.md`,
`templates/project-v2/agent-skills/savepoint-check/SKILL.md`,
`.savepoint/Design.md`, `internal/init/template_freshness_test.go`,
`internal/init/agent_skills_test.go`, `agent_skills_test.go`.

## Design References

Design sections 1, 4, 6, and 8.

## Guardrails

TPL-01, TPL-02, TPL-04, POL-02, TEST-06, TEST-08.

## Implementation Plan

1. Update AGENTS.md workflow, CLI Rules, and router-advance ownership.
2. Update skills in both trees together; diff to confirm byte identity.
3. Reconcile Design.md sections 1, 4, 6, and 8 against T-028..T-033.
4. Run template and skill tests.

## Boundaries

No production code. No Goal-is-mandatory wording (O-022).

## Technical Verification

`make build && make test-fast` at handoff; the Full Objective Check then
runs a fresh `make test-full`.

## Technical Evidence

Pending execution.

## Drift Notes

None expected.
