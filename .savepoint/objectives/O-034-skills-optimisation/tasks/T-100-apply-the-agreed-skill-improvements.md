---
id: T-100
title: Apply the agreed skill improvements
objective: O-034
status: in_progress
stage: audit
depends_on: [{task: T-099, requires: clear}]
complexity_tier: medium
complexity_reason: Text edits across seven parity-checked assets with content tests and upgrade delivery.
owner_validation: {required: true}
planned_by: {role: planner, session: plan-o034-20261003}
---

# Apply the agreed skill improvements

## Outcome

Every **apply** finding from the review Task is made in the canonical skills and references, and their scaffold copies stay byte-identical. Existing projects receive the changes through `savepoint upgrade-assets`, and each **follow-up** finding exists as an Issue.

## User Check

Run `savepoint upgrade-assets --dry-run` on a copy of an existing V2 project and see the revised skills listed as updates. Open `savepoint-task` and `savepoint-check` and find a short, plain rule on how to act on an advisory Code Health report.

## Done When

- Each apply finding is made exactly as recorded, or its deviation is noted in Technical Evidence with the reason. No finding marked no change or follow-up is applied.
- `savepoint-task` and `savepoint-check` carry the Code Health advice rule from O-034's Confirmed Design. It cites guardrails rather than restating them (POL-02).
- Each follow-up finding is recorded as an Issue under `agent-skills/references/issue-capture.md`, including the report next-step wording.
- Live and `templates/project-v2/` copies are byte-identical (TPL-01). Content tests pin each behavioural rule that was added or reworded, with one failure-path assertion per new rule.
- An upgrade test shows a project with unedited package skills receives the revised text, and a project with an edited skill keeps its bytes (FS-02, TEST-03).
- Every scenario walkthrough from the review baseline is re-walked against the revised text, and its result is recorded per scenario: same behaviour, or the intended improvement.
- The full token weight of the seven assets before and after is recorded. An increase is justified by name.
- `make build && make test-fast` passes.

## Context Files

`agent-skills/savepoint-idea/SKILL.md`; `agent-skills/savepoint-design/SKILL.md`; `agent-skills/savepoint-task/SKILL.md`; `agent-skills/savepoint-check/SKILL.md`; `agent-skills/references/check-method.md`; `agent-skills/references/issue-capture.md`; `agent-skills/references/commands-and-procedures.md`, and their copies under `templates/project-v2/agent-skills/`; `internal/init/agent_skills_test.go`; `internal/init/skill_validation_test.go`; `internal/init/template_freshness_test.go`; `internal/init/upgrade_test.go`; the review Task file for its findings. Edit only the files the apply findings name.

## Design References

O-034 Confirmed Design; Design section 1 (bundled Agent Skills, V2 agent workflow assets, upgrade-assets).

## Guardrails

TPL-01, TPL-02, TPL-04, FS-02, POL-01, POL-02, TEST-01, TEST-02, TEST-03, TEST-08, STYLE-07, STYLE-10.

## Implementation Plan

1. Record the before token weight, then apply findings file by file and mirror each change to the scaffold copy.
2. Add or adjust content tests for each new or reworded rule.
3. Create the follow-up Issues.
4. Re-walk the scenarios, record the after weight, and run the gate.

## Boundaries

No Go behaviour change, no new lifecycle state, field, command, or Goal-level Check, and no change to report wording or health classification. The recurring-review guidance belongs to the next Task.

## Technical Verification

Focused `make test-focused TEST=... PKGS=./internal/init` during iteration; `make build && make test-fast` at handoff. The later independent Check applies `agent-skills/references/check-method.md`.

## Technical Evidence

Executor evidence (not a Check, not CLEAR).

### Per-criterion outcomes

- **Findings applied as recorded.** A1, A2, A3, A4, A5, A6, A7, A8, T1, T2, T4 applied. F1 became I-122; N1..N6, T3, T5 untouched. Deviations: (a) A7 also trimmed the key names from the Task Write Boundary prose and the `savepoint-check` Purpose, and kept one mention in `issue-capture.md` and one in each AGENTS.md (the Check section), as recorded. (b) A8 removed the whole "Legacy V1 compatibility" block, including its closing Task `stage` / `status: done` lines; both are already in Terminology, so nothing is lost. (c) T2 put CODE-01 under a new "Code Health" table in the template Guardrails.md. (d) T4 needed only the A7 key-name trim; the template AGENTS.md already lacked the repo-only sections.
- **Code Health rule.** `savepoint-task` gains "Acting On A Code Health Report"; `savepoint-check` gains one sentence in Code Health Evidence. Both name `.savepoint/health/config.json` and cite a Guardrail by ID without restating it (POL-02); pinned by `TestSkillReviewRulesCiteGuardrailsInsteadOfRestatingThem`.
- **Follow-up Issue.** I-122 records the report next-step wording (F1). It is the only follow-up finding.
- **Parity (TPL-01).** `diff -rq agent-skills templates/project-v2/agent-skills` differs only by `bubbletea-tui-design`, which is out of scope. Existing byte-parity tests pass.
- **Content tests.** `TestSkillReviewRulesArePinned` (table: three triggers, free reads, Code Health rule, platform evidence, repair re-run, key-name removal, with stale-phrase failure assertions), `TestSkillReviewRulesCiteGuardrailsInsteadOfRestatingThem`, `TestScaffoldedPolicyDocumentsPointInsteadOfRestating`, `TestRepoAgentsGuideDropsDuplicateRoutingAndLegacySections`. Two existing phrase pins were updated for the reworded text.
- **Upgrade.** `TestUpgradeDeliversRevisedSkillsToUneditedProjects` (manifest-tracked old text receives the revised skill) and `TestUpgradeKeepsAnEditedSkillAndOffersTheRevision` (edited skill keeps its bytes, revision offered as `.new`). A manual `savepoint upgrade-assets --dry-run` on a scaffolded copy with the old task, check and design skills reported all three as updated.
- **Scenario re-walk.** 1 Next to skill routing: improved; each Trigger now names the `Next` line first and says a pasted line is the selection, so the router-state mismatch no longer needs reconciling. 2 REPLAN re-entry: same behaviour; the Replan trigger now names `Replan`, the frontmatter-block rule is untouched. 3 Waived Task Check reaching the Full Check: same behaviour; the five waiver statements are unchanged (T1 removed the sixth, in template Design.md). 4 Issue-only repair: improved; the executor now re-runs Proof Needed and neighbouring cases before `repair_attempted` and marks unrun platforms unverified; the router-advance steps are unchanged. 5 Advisory Code Health report: improved; the executor stops at the watch line and reports the rest, the checker records an aim-to-watch signal as an observation.
- **Token weight (bytes, ~4 per token).** Seven assets 82,210 B before, 83,819 B after (+1,609 B, about +400 tokens, +2.0%). The increase is the named Code Health section and Check sentence (A1, A2), the repair re-run rule (A4), the platform-evidence sentence (A6), the free-reads sentence (A3) and the three Next-line triggers (A5), partly offset by the key-name trim (A7). Outside the seven: repo `AGENTS.md` 23,138 to 20,857 B (-2,281, loaded every session in this repo); template `AGENTS.md` 17,642 to 17,584 B; template `Design.md` 2,828 to 1,711 B; template `Guardrails.md` 3,729 to 3,901 B (CODE-01).
- **Gate.** `make build && make test-fast` passed (exit 0) after the final edit, including `internal/init`.

### Limitations

- Existing projects receive the skill, reference and AGENTS.md changes through `upgrade-assets`, but not the template `Design.md` and `Guardrails.md` edits (T1, T2): the upgrade skips files that already exist, so only new projects get them. Delivering CODE-01 to existing projects is a project decision for the owner.
- Scenario re-walks are written, not live agent runs (as the Objective's Confirmed Design allows). No `make test-full` was run; this Task changes no code.
- Extra reads, logged: `README.md` and the `internal/init` upgrade and manifest code (to answer the README review and the upgrade criterion), `savepoint init` and `upgrade-assets --help` output.
- README: at the owner's request the project README was reviewed against the code and updated: `init` creates G-001 (it said R-001), the `.savepoint/` tree (router comment, `health/`, `releases/G-001`), `upgrade-assets` flags and edited-skill behaviour, the Goal filter wording, and a sentence on the watch line in the Code Health section.

Files read: the Context Files, `AGENTS.md`, `templates/project-v2/.savepoint/Design.md`, `Guardrails.md`, `README.md`. Files changed: the seven live and scaffold skill and reference pairs that changed (task, check, design, issue-capture), both AGENTS.md files, template Design.md and Guardrails.md, `README.md`, `internal/init/agent_skills_test.go`, `template_freshness_test.go`, `upgrade_test.go`, Issue I-122, and this Task.

## Drift Notes

If an apply finding turns out to need a product decision, return it as a follow-up Issue rather than widening this Task.
