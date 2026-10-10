---
id: T-123
title: Bring the skills and guides in line with the new decision rules
objective: O-044
status: planned
depends_on: [{task: T-122, requires: clear}]
owner_validation: {required: false}
planned_by: {role: planner, session: planning-o044-20261010}
---

# Bring the skills and guides in line with the new decision rules

## Outcome

The Check skill, its shared method, the agent guide and the README describe the same decision rules the runtime now applies: decisions keep their originating Check, a later Check records whether each still applies and why, the owner renews only a materially changed decision, `Assess` is an applicability step and not a new Check, waived evidence is never CLEAR, and owner decisions are recorded in structured form rather than only as Issue history. Existing projects receive the guidance through `savepoint upgrade-assets`.

## User Check

Read the savepoint-check Closure Rules and the AGENTS.md Next-verb list, then walk the TheShed O-002 story against them: a checker following them after C-006 would record carry entries and report "ready to close by exception", not ask for the waiver and acceptance again.

## Done When

1. `agent-skills/savepoint-check/SKILL.md` Closure Rules replace "acceptance naming a superseded Check does not count" and "applies only to the Check it names" with the carry-forward rule, show the `carried_forward`, `owner_validation.scope` and Check `unmet` shapes, say when a re-check must assess every prior owner decision, define the `Assess` procedure (applicability only, no new Check record), and require an owner's waiver to be recorded as a structured `exception` on the owner's explicit instruction.
2. `agent-skills/references/check-method.md` Re-check After Remediation tells a checker to compare each prior decision's scope with what changed and record the assessment, without demanding renewal for unchanged scope.
3. Both copies of the AGENTS.md workflow block add `Assess` → `savepoint-check` (applicability only) and the Verification Policy states the rule once, pointing to the skill rather than restating it (STYLE-07).
4. README's owner-facing description of acceptance and exceptions matches the runtime (TPL-02).
5. Canonical skills and references stay byte-identical to their `templates/project-v2/` copies (TPL-01); `upgrade-assets` delivers them to existing projects (TPL-04) with user-edited files preserved (FS-02, TEST-03).
6. Content tests assert the new rule text and the absence of the old exact-Check wording; `make build && make test-fast` passes (TEST-08). Per-criterion evidence follows AGENTS.md's Verification Policy, including an explicit owner waiver if the optional Task Check is skipped.

## Context Files

`agent-skills/savepoint-check/SKILL.md`, `templates/project-v2/agent-skills/savepoint-check/SKILL.md`, `agent-skills/references/check-method.md`, `templates/project-v2/agent-skills/references/check-method.md`, `AGENTS.md`, `templates/project-v2/AGENTS.md`, `README.md`, `agent_skills_test.go`, `internal/init/agent_skills_test.go`, `internal/init/agents_test.go`, `internal/init/template_freshness_test.go`, `.savepoint/objectives/O-044-keep-owner-decisions-across-rechecks/Objective.md`.

## Design References

O-044 Confirmed Design. I-141 Root cause 3–4 and Planned Fix, Guidance.

## Guardrails

TPL-01, TPL-02, TPL-04, FS-02, STYLE-07, TEST-01, TEST-03, TEST-08.

## Implementation Plan

1. Confirm the field names, blocker wording and `Assess` verb shipped by the earlier O-044 Tasks; return REPLAN REQUIRED if they differ from this plan.
2. Rewrite the Closure Rules and re-check guidance; keep each rule in one home.
3. Update both AGENTS.md copies and README.
4. Copy canonical files to the scaffold and update content tests.
5. Run `make build && make test-fast`.

## Boundaries

Guidance only: no Go runtime change. Coordinate wording with O-039 and O-040 if they have landed in the same files; do not reintroduce duplicated policy they removed. Do not edit TheShed or any project's records.

## Technical Verification

Focused `make test-focused TEST=...` during iteration; `make build && make test-fast` for handoff. The Full Objective Check needs fresh `make test-full` under `agent-skills/references/check-method.md`.

## Technical Evidence

Pending execution.

## Drift Notes

None yet.
