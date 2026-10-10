---
id: T-123
title: Bring the skills and guides in line with the new decision rules
objective: O-044
status: done
depends_on: [{task: T-122, requires: clear}]
owner_validation:
    required: false
    accepted_check: ""
planned_by: {role: planner, session: planning-o044-20261010}
check_waiver:
    task: T-123
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-10T06:26:10Z"
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

Executed 2026-10-10. Plan step 1 confirmed: runtime ships `carried_forward` (`check, applies, assessed_by, assessed_at, reason, material_change`), `owner_validation.scope`, Check `unmet`, and the `Assess` Next verb; no REPLAN needed.

Per-criterion outcomes:
1. Pass: `agent-skills/savepoint-check/SKILL.md` Closure Rules drop the exact-Check wording, show `carried_forward`, `scope` and `unmet` shapes, require assessing every prior decision on a re-check, define `Assess` (no new Check), and require a structured `exception` on explicit owner instruction. Trigger also names `Assess`.
2. Pass: `check-method.md` Re-check After Remediation compares each decision's scope with what changed and records the assessment without renewal for unchanged scope.
3. Pass: `Assess` → `savepoint-check` added to the workflow line in root `AGENTS.md` (both occurrences) and `templates/project-v2/AGENTS.md`; one Verification Policy bullet states the carry rule and points to the skill.
4. Pass: README "You decide" bullet describes carry-forward, `Assess`, and ready-to-close-by-exception, matching `internal/resume` wording.
5. Pass: skill and reference copied byte-identical to `templates/project-v2/` (`cmp` clean); existing template-freshness and upgrade tests cover delivery and user-edit preservation.
6. Pass: `TestSavepointCheckSkillClosureRules` asserts the new text and the absence of old wording; new `TestCheckMethodReCheckAssessesPriorDecisions` and `TestAgentsGuideRoutesAssessToCheckSkill`.

Commands run: `make build`; `make test-fast` (exit 0).
Files changed: the files above plus `internal/init/agent_skills_test.go`. Extra reads: `internal/data/next.go`, `internal/resume/*.go`, `Guardrails.md` (to match shipped field names and wording).
Limitations: no Task Check requested and no owner waiver recorded yet; `make test-full` not run (not required for ordinary handoff).

## Drift Notes

None yet.
