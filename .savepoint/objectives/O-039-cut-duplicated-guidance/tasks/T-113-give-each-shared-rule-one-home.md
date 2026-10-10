---
id: T-113
title: Give each shared rule one home
objective: O-039
status: planned
depends_on: [{task: T-112, requires: clear}]
complexity_tier: medium
complexity_reason: Text moves across the managed block, four skills, three references, scaffold copies, and content tests.
owner_validation: {required: true}
planned_by: {role: planner, session: plan-o039-20261006}
---

# Give each shared rule one home

## Outcome

Goal Context and the Verification Policy rules (Task-check waiver, `clear`/`accepted` dependencies, Objective Check coverage) are stated once, in the managed AGENTS.md block. The skills and references keep only role-specific consequences and cite the section by name. The `savepoint health report` human-only rule is in the managed block template.

## User Check

Search the four skills and three references for the Goal Context paragraph and the waiver rule, and find only a one-line pointer to the AGENTS.md section. Run `savepoint upgrade-assets --dry-run` on a copy of a V2 project and see the revised assets listed.

## Done When

- No shared rule paragraph appears in more than one of: the managed block, the four skills, the three references.
- Each skill still states what its own role does with the rule (for example, Check's authority to write closure evidence).
- `health report` is human-only in the managed block template, and this repository's managed block matches the template.
- Live and `templates/project-v2/` copies are byte-identical (TPL-01). Content tests assert each moved rule at its single home, with a failure-path assertion that it is not restated.
- An upgrade test shows unedited package skills receive the revised text, and an edited skill keeps its bytes (FS-02, TEST-03).
- O-034's scenario walkthroughs are re-walked against the revised text, with the result recorded per scenario.
- Token weight of the managed block, four skills, and three references before and after is recorded.
- `make build && make test-fast` passes.

## Context Files

`templates/project-v2/AGENTS.md`; `AGENTS.md`; `agent-skills/savepoint-idea/SKILL.md`; `agent-skills/savepoint-design/SKILL.md`; `agent-skills/savepoint-task/SKILL.md`; `agent-skills/savepoint-check/SKILL.md`; `agent-skills/references/check-method.md`; `agent-skills/references/issue-capture.md`; `agent-skills/references/commands-and-procedures.md`, and their copies under `templates/project-v2/agent-skills/`; `internal/init/agent_skills_test.go`; `internal/init/skill_validation_test.go`; `internal/init/template_freshness_test.go`; `internal/init/upgrade_test.go`; `internal/init/agents_test.go`; `.savepoint/objectives/O-034-skills-optimisation/tasks/T-099-review-the-skills-against-what-the-code-health-goal-taught-us.md` for the scenarios.

## Design References

O-039 Confirmed Design; Design section 1 (bundled Agent Skills, V2 agent workflow assets, upgrade-assets).

## Guardrails

TPL-01, TPL-02, TPL-04, FS-02, POL-01, POL-02, TEST-01, TEST-02, TEST-03, TEST-08, STYLE-07, STYLE-10.

## Implementation Plan

1. Record the before token weight.
2. For each shared rule, make the managed block the full statement and replace restatements with a pointer plus any role-specific consequence. Mirror each change to the scaffold copy.
3. Move the `health report` rule into the template and refresh this repository's managed block.
4. Move content-test assertions to the single home and add the not-restated checks.
5. Re-walk the scenarios, record the after weight, and run the gate.

## Boundaries

No change to what any rule requires; no reference restructuring or checklists (O-040); no V1/V2 history or terminology edits beyond the moved text (next Task).

## Technical Verification

Focused tests during iteration; handoff gate from AGENTS.md's Verification Policy.

## Technical Evidence

Pending execution.

## Drift Notes

None yet.
