---
id: T-124
title: Make the long references easy to navigate
objective: O-040
status: done
depends_on: []
complexity_tier: medium
complexity_reason: Structural edits to two references plus scaffold copies, a pointer audit across four skills, and matching content tests.
owner_validation:
    required: true
    accepted_check: ""
planned_by: {role: planner, session: plan-o040-20261010}
check_waiver:
    task: T-124
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-10T07:13:17Z"
---

# Make the long references easy to navigate

## Outcome

`check-method.md` and `issue-capture.md` open with a contents list. `check-method.md` has a copyable Full Objective Check progress checklist beside the Quick Check Procedure. A recorded audit shows that every reference a workflow skill must apply is one step from that skill or from the always-loaded AGENTS.md managed block. Any remaining chain is fixed with a direct pointer.

## User Check

Open `agent-skills/references/check-method.md`. The first screen lists every section, and the Full checklist can be copied into a Check record and ticked off in order. Open `issue-capture.md` and see its contents list. Read the audit table in this Task's evidence and confirm that no rule needs two hops.

## Done When

- Both references open, after the intro paragraph, with a contents list naming every `##` section in file order. A content test fails if a `##` heading is missing from the list or is out of order.
- `check-method.md` has a `### Full Check Progress Checklist` under Quick And Full Evidence Modes. It is a fenced Markdown checkbox list that names each section Full mode applies, in method order, and the Full bullet points to it. A content test asserts that every listed step names an existing heading.
- Quick Check Procedure wording and every section's rule content are unchanged (diff shows additions only, plus any pointer fixes from the audit).
- Audit table in evidence: for each of the four skills, every reference it applies, every further file that reference sends the reader to, and the result (direct, always-loaded, or fixed). Any fix makes the skill name the target directly.
- Live and `templates/project-v2/agent-skills/` copies are byte-identical (TPL-01). An upgrade test shows that unedited projects receive the revised references (FS-02, TEST-03).
- The O-034 scenarios (T-099 baselines) are re-walked against the revised text, with a result per scenario. Token weight before and after is recorded, and the increase is justified by name (contents lists and checklist).
- `make build && make test-fast` passes.

## Context Files

`agent-skills/references/check-method.md`; `agent-skills/references/issue-capture.md`; `agent-skills/references/commands-and-procedures.md`; `agent-skills/savepoint-idea/SKILL.md`; `agent-skills/savepoint-design/SKILL.md`; `agent-skills/savepoint-task/SKILL.md`; `agent-skills/savepoint-check/SKILL.md`; their copies under `templates/project-v2/agent-skills/`; `templates/project-v2/AGENTS.md`; `internal/init/agent_skills_test.go`; `internal/init/skill_validation_test.go`; `internal/init/template_freshness_test.go`; `internal/init/upgrade_test.go`; `.savepoint/objectives/O-034-skills-optimisation/tasks/T-099-review-the-skills-against-what-the-code-health-goal-taught-us.md` for the scenarios.

## Design References

O-040 Confirmed Design; Design section 1 (bundled Agent Skills, V2 agent workflow assets, upgrade-assets).

## Guardrails

TPL-01, TPL-02, TPL-04, FS-02, POL-01, POL-02, TEST-01, TEST-02, TEST-03, TEST-08, STYLE-07, STYLE-10.

## Implementation Plan

1. Record the before token weight (managed block + four skills + three references, words and bytes with `wc`), measured as in T-113.
2. Run the pointer audit and record the table. Fix any remaining chain by naming the target directly in the skill.
3. Add the contents lists to both references.
4. Add the Full Check Progress Checklist and point the Full bullet to it.
5. Mirror every change to the scaffold copies and add the content tests (list completeness and order, checklist steps name real headings, parity).
6. Re-walk the scenarios, record the after weight, and run the gate.

## Boundaries

Do not change what any rule requires. Do not add or edit AGENTS.md rules, skill descriptions, or the write → resume → fix rule (next Task). Do not add a contents list to a skill or to `commands-and-procedures.md` (under 100 lines). Do not add new commands or validators.

## Technical Verification

Focused tests during iteration; handoff gate from AGENTS.md's Verification Policy. Later Check evaluation follows `agent-skills/references/check-method.md`.

## Technical Evidence

### Criteria outcomes

- Contents lists: `check-method.md` (16 `##` sections) and `issue-capture.md` (9) open with a `**Contents**` list after the intro. `TestLongReferencesOpenWithContentsListInFileOrder` compares each list to the file's `##` headings, skipping fenced code, in order; it runs on the live and scaffold copies.
- Checklist: `### Full Check Progress Checklist` sits under Quick And Full Evidence Modes as a fenced markdown checkbox list of 16 steps in method order; the Full bullet now points to it. `TestCheckMethodFullChecklistNamesRealHeadingsInOrder` asserts the placement, the pointer, that each step names an existing heading, that steps keep file order, and that no Full section is omitted.
- Unchanged rules: the only removed line is the old Full bullet, which was extended with the pointer sentence. Every other change is an addition. Quick Check Procedure text is untouched.
- Parity: reference files copied to `templates/project-v2/agent-skills/references/`; existing live/template parity tests pass. `TestUpgradeDeliversRevisedReferencesToUneditedProjects` shows an unedited project (previous text, no list) is refreshed to the revised references by upgrade-assets.
- Gate: `make build && make test-fast` exit 0.

### Pointer audit (no fixes needed)

| Skill | References it applies | Further files those references name | Result |
|-------|-----------------------|--------------------------------------|--------|
| savepoint-idea | none | none | direct |
| savepoint-design | `check-method.md`, `issue-capture.md`, `commands-and-procedures.md`, each named in the skill | `commands-and-procedures.md` names `check-method.md`, which Design also names directly; `AGENTS.md` | direct / always-loaded |
| savepoint-task | `issue-capture.md`, `check-method.md`, each named in the skill | `AGENTS.md` (Verification Policy) | direct / always-loaded |
| savepoint-check | `check-method.md` (in full), `issue-capture.md`, each named in the skill | `AGENTS.md`; `.savepoint/Guardrails.md` (a project file, not a reference) | direct / always-loaded |

### Token weight (words / bytes, as in T-113)

| Part | Before | After |
|------|--------|-------|
| AGENTS.md managed block | 2253 / 14904 | 2253 / 14904 |
| Four skills + three references | 11016 / 74879 | 11274 / 76340 |

The +258 words / +1461 bytes are the two contents lists and the checklist (with its intro line and the Full-bullet pointer).

### O-034 scenarios re-walked against the revised text

1. Next to skill routing: unchanged; no skill or AGENTS.md text touched.
2. REPLAN REQUIRED re-entry: unchanged; no `savepoint-task` or `savepoint-design` text touched.
3. Waived Task Check to Full Check: unchanged; Full still reviews every owned Task, and the checklist adds no waiver rule.
4. Issue-only repair: unchanged; `issue-capture.md` gained only the contents list.
5. Advisory Code Health report: unchanged; the checklist names Collect Code Health Evidence without restating it.

### Limitations

Scenarios were re-walked against the written text only, with no live agent run. `make test-full` was not run (not required for an ordinary Task). Extra reads: `internal/init/skill_validation_test.go` helpers and `internal/init/upgrade.go`/`upgrade_test.go` signatures, all listed Context Files except `upgrade.go`.

## Drift Notes

None yet.
