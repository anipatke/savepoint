---
id: T-008
title: Align active guidance with Goals
objective: O-013
planned_by: {role: planner, session: goals-terminology-20260921}
status: planned
complexity_tier: medium
complexity_reason: "The active repository guide, architecture, public README, workflow skills, and shipped V2 copies must change together while historical evidence remains byte-preserved."
depends_on: [{task: T-007, requires: clear}]
owner_validation:
    required: false
---

# T-008: Align active guidance with Goals

## Outcome

Active repository guidance and shipped V2 workflow documentation describe the
delivery context as Goals, explain the preserved R###/Release storage boundary,
and remain consistent with the implemented V2 board.

## User Check

Read the reconciled README and repository guide alongside the V2 board wording.
Confirm that a new user learns “Goals” as the public concept, while the
compatibility note makes existing `R###` and `release:` files understandable.
No additional product choice is required for this documentation Task.

## Done When

- `README.md`, `AGENTS.md`, `.savepoint/Design.md`, and
  `.savepoint/Guardrails.md` use Goal/Goals for public V2 terminology and
  retain only necessary compatibility references to Release.
- `templates/project-v2/AGENTS.md` and every active V2 skill/reference copy
  describe the same Goal vocabulary, shortcut, optional boundary, and check
  semantics.
- Documentation does not imply that Goals publish, deploy, tag, or create
  changelogs, and does not turn owner waiver/acceptance into technical
  clearance.
- Canonical active skills and their `templates/project-v2/` copies remain
  byte-identical where the template contract requires it.
- Targeted searches distinguish allowed internal, migration, and historical
  Release references from stale public UI instructions.
- Archived V1 material, immutable Checks, Issues, and current historical
  evidence are not rewritten.

## Context Files

`README.md`, `AGENTS.md`, `.savepoint/Design.md`, `.savepoint/Guardrails.md`, `templates/project-v2/AGENTS.md`, `agent-skills/savepoint-idea/SKILL.md`, `agent-skills/savepoint-design/SKILL.md`, `agent-skills/savepoint-task/SKILL.md`, `agent-skills/savepoint-check/SKILL.md`, `agent-skills/references/check-method.md`, `agent-skills/references/issue-capture.md`, `agent-skills/references/commands-and-procedures.md`, `templates/project-v2/agent-skills/savepoint-idea/SKILL.md`, `templates/project-v2/agent-skills/savepoint-design/SKILL.md`, `templates/project-v2/agent-skills/savepoint-task/SKILL.md`, `templates/project-v2/agent-skills/savepoint-check/SKILL.md`, `templates/project-v2/agent-skills/references/check-method.md`, `templates/project-v2/agent-skills/references/issue-capture.md`, `templates/project-v2/agent-skills/references/commands-and-procedures.md`.

## Design References

Design sections 2, 3, 6, 7, 8, 10, and 13.

## Guardrails

TPL-01, TPL-02, TPL-04, ARCH-04, TEST-01, TEST-06, TEST-08, STYLE-07,
STYLE-09, and STYLE-10.

## Implementation Plan

1. Compare each named active document with the completed board vocabulary and
   the compatibility evidence from T-006.
2. Update public Goal terminology, shortcut guidance, hierarchy descriptions,
   and compatibility notes without changing policy semantics.
3. Apply the same edits to the canonical V2 skills and their shipped copies;
   leave V1/archive documentation untouched.
4. Search for stale public Release wording, verify template parity, and run
   `git diff --check`, `make build`, and `make test`.

## Boundaries

No production behavior, persisted schema, migration archive, immutable Check,
Issue, historical V1, or independent Release/Goal policy rewrite.

## Technical Verification

Exact-path active-document terminology searches, template freshness/parity
checks, focused board tests, `git diff --check`, `make build`, and `make test`.
The mandatory Full Objective Check must verify documentation and board agreement
before O-013 can close.

## Technical Evidence

Pending execution: changed/no-change document inventory, search output, template
parity result, command results, files read/changed, and limitations.

## Drift Notes

Any documentation change that requires a persisted schema migration or changes
Goal/Release gate semantics is REPLAN REQUIRED and returns to design.
