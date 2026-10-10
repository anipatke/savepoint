---
id: T-125
title: Check records after editing and name trigger words
objective: O-040
status: done
depends_on: [{task: T-124, requires: clear}]
complexity_tier: medium
complexity_reason: One managed-block rule, pointer lines and descriptions in four skills, scaffold copies, and matching content tests.
owner_validation:
    required: true
    accepted_check: ""
planned_by: {role: planner, session: plan-o040-20261010}
check_waiver:
    task: T-125
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-10T07:15:21Z"
---

# Check records after editing and name trigger words

## Outcome

The managed AGENTS.md block says once that after an agent edits any `.savepoint/` record, it runs `savepoint resume`. If the record fails to load, the agent fixes it and runs `savepoint resume` again until it loads. Each skill that writes records points to that rule. Each workflow skill's `description` names the `Next` words that select it.

## User Check

Read CLI Rules in `AGENTS.md` and find one write → resume → fix rule. Open each skill and find a one-line pointer to it, not a restatement. Read the four skill descriptions and see `Plan`/`Replan` on design, `Start`/`Build`/`Test`/`Pick a Task in`/`Fix` on task, and `Check`/`Assess` on check. Idea keeps "router state is idea".

## Done When

- CLI Rules in `templates/project-v2/AGENTS.md` widens the existing "After creating or renaming any other identity-bearing record…" sentence into one rule. It covers any `.savepoint/` record write, says to run `savepoint resume`, and says to fix the record and re-run until it strict-loads. It says not to substitute other work. This repository's managed block is refreshed to match through `./savepoint upgrade-assets`.
- Idea, design, task, and check each carry a one-line pointer to that rule by section name. Check's workflow step 4 drops its own restatement in favour of the pointer. A content test asserts the rule at its home and that no skill restates it.
- Skill `description` fields name their `Next` words exactly as AGENTS.md's Workflow routes them, stay under 1024 characters, and keep their current role wording. A content test cross-checks each word against the Workflow list, failing on a missing or extra word.
- What each rule requires is otherwise unchanged. Live and scaffold copies are byte-identical (TPL-01). An upgrade test shows that unedited projects receive the revised skills and block (FS-02, TEST-03).
- O-034 scenarios re-walked with a result each; token weight before and after recorded, with the increase justified by name.
- `make build && make test-fast` passes.

## Context Files

`templates/project-v2/AGENTS.md`; `AGENTS.md`; `agent-skills/savepoint-idea/SKILL.md`; `agent-skills/savepoint-design/SKILL.md`; `agent-skills/savepoint-task/SKILL.md`; `agent-skills/savepoint-check/SKILL.md`; their copies under `templates/project-v2/agent-skills/`; `internal/init/agent_skills_test.go`; `internal/init/skill_validation_test.go`; `internal/init/template_freshness_test.go`; `internal/init/upgrade_test.go`; `internal/init/agents_test.go`; `.savepoint/objectives/O-034-skills-optimisation/tasks/T-099-review-the-skills-against-what-the-code-health-goal-taught-us.md` for the scenarios.

## Design References

O-040 Confirmed Design; Design section 1 (bundled Agent Skills, V2 agent workflow assets, upgrade-assets).

## Guardrails

TPL-01, TPL-02, TPL-04, FS-02, POL-01, POL-02, TEST-01, TEST-02, TEST-03, TEST-08, STYLE-07, STYLE-10.

## Implementation Plan

1. Confirm T-124's reference changes landed; return REPLAN REQUIRED if they did not. Record the before token weight.
2. Widen the CLI Rules sentence in the template and refresh this repository's managed block.
3. Add the pointer line to each skill and replace Check step 4's restatement.
4. Add the `Next` words to the four descriptions (Idea unchanged apart from any pointer).
5. Mirror the changes to the scaffold copies and add the content tests (rule home, not restated, description words against Workflow, parity).
6. Re-walk the scenarios, record the after weight, and run the gate.

## Boundaries

Do not add new commands, validators, or `Next` words, and do not change routing logic in `internal/data`. Do not restructure the references (previous Task). Do not change what any other rule requires.

## Technical Verification

Focused tests during iteration; handoff gate from AGENTS.md's Verification Policy. Later Check evaluation follows `agent-skills/references/check-method.md`.

## Technical Evidence

### Criteria outcomes

- Rule: `templates/project-v2/AGENTS.md` CLI Rules now holds one paragraph: after writing or editing any `.savepoint/` record (including creating or renaming), run `savepoint resume`; if it fails to load, fix and re-run until it loads; do not substitute other work. It replaces the old identity-bearing sentence. This repository's block was refreshed with `./savepoint upgrade-assets` (AGENTS.md diff: +3/-1).
- Pointers: idea, design, task, check each end with one line, "After writing any `.savepoint/` record, follow AGENTS.md's CLI Rules (write → resume → fix)." Design's and Check step 4's own restatements were removed. `TestWriteResumeFixRuleHasOneHomeAndSkillsPointToIt` asserts the rule once at each home and no restatement in any skill (live and scaffold).
- Descriptions: design names `Plan`/`Replan`; task names `Start`/`Build`/`Test`/`Pick a Task in`/`Fix`; check names `Check`/`Assess`; idea keeps "router state is idea". All under 1024 characters (250, 356, 346, 230). `TestSkillDescriptionsNameTheirNextWords` checks each word is still routed by the Workflow line and fails on a missing or extra word.
- Unchanged otherwise: only descriptions, pointers and the two removed restatements changed. Live and scaffold skills are byte-identical (copied; parity tests pass). `TestUpgradeDeliversRevisedSkillsAndRuleBlock` shows unedited projects receive the revised skills and block.
- Three older tests that asserted the old sentence were updated to the new rule (`agent_skills_test.go`, `template_freshness_test.go`, `v2_scaffold_test.go`).
- Gate: `make build && make test-fast` exit 0; `savepoint doctor` all clean.

### Token weight (words / bytes)

| Part | Before | After |
|------|--------|-------|
| AGENTS.md managed block | 2253 / 14904 | 2273 / 15016 |
| Four skills + three references | 11274 / 76340 | 11326 / 76677 |

Increase (+72 words, +449 bytes): the single rule paragraph, four one-line pointers, and the Next words in three descriptions, partly offset by the two restatements removed.

### O-034 scenarios re-walked

1. Next to skill routing: improved; descriptions now name the Next words, so a pasted `Next: Check O-0xx` matches the check skill without reconciling against router state (A5).
2. REPLAN REQUIRED re-entry: unchanged; `Replan` is now in design's description alongside the existing REPLAN REQUIRED wording.
3. Waived Task Check to Full Check: unchanged; no waiver text touched.
4. Issue-only repair: unchanged; `Fix` is now named in task's description, routing and handoff text untouched.
5. Advisory Code Health report: unchanged; not touched.

### Limitations

Re-walk is against written text only, no live agent run. `make test-full` not run. T-124 had landed (references carry contents lists), confirmed before starting. `gofmt -l` flags `internal/init/manifest_test.go`, which this Task did not touch. Extra reads: `internal/init/v2_scaffold_test.go`.

## Drift Notes

None yet.
