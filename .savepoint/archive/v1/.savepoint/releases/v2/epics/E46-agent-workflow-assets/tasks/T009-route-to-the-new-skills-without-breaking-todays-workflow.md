---
id: E46-agent-workflow-assets/T009-route-to-the-new-skills-without-breaking-todays-workflow
title: Route to the new skills without breaking today's workflow
status: done
objective: Document V2 routing in both guides as explicitly inactive until cutover, and lock the four-role contract with one cross-cutting test.
depends_on:
    - E46-agent-workflow-assets/T002-take-a-rough-idea-without-demanding-a-document
    - E46-agent-workflow-assets/T003-plan-one-objective-at-a-time
    - E46-agent-workflow-assets/T004-give-every-task-a-title-a-person-can-read
    - E46-agent-workflow-assets/T005-execute-a-plan-and-say-what-actually-happened
    - E46-agent-workflow-assets/T006-check-in-a-fresh-session-and-record-the-proof
    - E46-agent-workflow-assets/T007-capture-follow-up-without-a-fourth-phase
    - E46-agent-workflow-assets/T008-carry-project-commands-and-procedures-forward
complexity_tier: medium
complexity_reason: Edits the live guide that governs this repository while it is mid-build, plus the whole-set contract matrix.
---

# T009: Route to the new skills without breaking today's workflow

## Problem

By this point the four V2 skills exist but nothing says how they connect. The routing has to be written down — `state: idea|design|task|check` maps to one skill each, and `REPLAN REQUIRED` derives a planner action rather than becoming a fifth phase — without that routing taking effect. This repository is still building V2 under the V1 lifecycle, and an agent that reads `AGENTS.md` tomorrow must still land on `savepoint-create-task` for `epic-task-breakdown`.

So the V2 routing block ships as documentation marked inactive until cutover, in both the live guide and the shipped copy, while the V1 activation table stays exactly as it is. E47 switches the scaffold default and retires the legacy instructions; E50 removes the transitional readers. Neither happens here.

The other half is the contract nobody has checked yet: that the four roles, read together, actually partition the writes. One skill writing what another owns — a checker that can edit acceptance criteria, an executor that can close an Issue — is the failure that makes the whole model decorative, and it is only visible when all four are asserted as a set.

## Context Files

- `AGENTS.md`
- `templates/project/AGENTS.md`
- `agent-skills/savepoint-idea/SKILL.md`
- `agent-skills/savepoint-design/SKILL.md`
- `agent-skills/savepoint-task/SKILL.md`
- `agent-skills/savepoint-check/SKILL.md`
- `internal/init/agent_skills_test.go`
- `internal/init/template_freshness_test.go`
- `internal/init/skill_validation_test.go`
- `internal/init/stale_reference_test.go`
- `.savepoint/releases/v2/v2-Design.md`

## Acceptance Criteria

- [x] `AGENTS.md` and `templates/project/AGENTS.md` both carry an identical V2 routing section mapping `state: idea` → `savepoint-idea`, `design` → `savepoint-design`, `task` → `savepoint-task`, `check` → `savepoint-check`.
- [x] That section states in its own text that V2 routing is not active, that the V1 activation table governs current work, and that activation belongs to E47's cutover.
- [x] It states that `REPLAN REQUIRED` routes back to `savepoint-design` and is not a fifth router state.
- [x] It names the shared references — `check-method.md`, `issue-capture.md`, `commands-and-procedures.md` — as non-triggerable and loaded by the skills that own them.
- [x] The V1 `## Skill Activation` table and every V1 skill file are unchanged; all nine V1 `savepoint-*` skills are still present in both trees.
- [x] `internal/init/agent_skills_test.go` asserts the full role matrix as a set: for each of the four skills, the records it may write and the records it must not, proving that Check writing, Issue closure, acceptance-criteria edits, and owner acceptance each have exactly one owner.
- [x] The same test asserts that all four V2 skills and all three shared references exist in both trees with byte parity.
- [x] The same test asserts obsolete vocabulary is absent from the V2 assets: no `epic`, no `PRD`, no `audit register`, no `defect-building`, and no `phase` as a lifecycle term.
- [x] `TestProjectGuidanceTemplatesMirrorLiveGuidance`, `TestProjectTemplatesUseCurrentWorkflow`, and the stale-reference tests all pass unchanged.
- [x] Reading `AGENTS.md` top to bottom leaves an agent working the V1 workflow today, with the V2 section readable as future state rather than as competing instructions.

## Implementation Plan

- [x] Draft the V2 routing section: the four-state table, the inactive-until-cutover statement, the replan routing rule, and the shared reference list.
- [x] Add it to `AGENTS.md` below the existing V1 sections, then mirror it verbatim into `templates/project/AGENTS.md`.
- [x] Verify the existing mirroring and stale-reference tests still pass with the new section present, and adjust the new text rather than the existing tests if they fail.
- [x] Add the cross-cutting role matrix test to `internal/init/agent_skills_test.go`, driven by a table of skill name, permitted writes, and forbidden writes.
- [x] Add the set-completeness, parity, and obsolete-vocabulary assertions over the four skills and three references.
- [x] Re-read the four skills as a set and fix any overlap the matrix exposes in the skill text, not in the test.
- [x] Run `go test ./internal/init/...`, then `make build && make test`.

## Context Log

**Files read:** all `## Context Files` listed above, plus `.savepoint/releases/v2/epics/E46-agent-workflow-assets/tasks/T002` through `T008` frontmatter (extra read, to confirm dependency status before starting — all `status: done`).

**Files edited:**
- `AGENTS.md`, `templates/project/AGENTS.md` — appended identical `## V2 Routing (Inactive Until Cutover)` section: four-state table, inactive-until-cutover statement naming E47, the `REPLAN REQUIRED` routing rule, and the three shared references named non-triggerable with their owning skill.
- `agent-skills/savepoint-idea/SKILL.md`, `templates/project/agent-skills/savepoint-idea/SKILL.md` — reworded "a prepared PRD" to "a prepared requirements document" so the obsolete-vocabulary test's ban on `PRD` holds; no behavior change, both trees kept byte-identical.
- `internal/init/agent_skills_test.go` — added `TestV1SkillSetUnchangedAlongsideV2Routing`, `TestV2SkillSetIsCompleteWithByteParity`, `TestV2SkillSetHasNoObsoleteVocabulary`, `TestV2SkillRoleMatrixPartitionsSensitiveWrites` (table-driven ownership check for Check-record writing, Issue closure, acceptance-criteria edits, and owner-acceptance claims across the four V2 skills), and `TestAgentsGuidesCarryV2RoutingSection`.

**Per-criterion outcome:** all ten acceptance criteria above verified directly — the routing section text, the role-matrix test, and a full local test run (see commands below). No pre-existing skill text needed changes beyond the single `PRD` reword; the four skills' existing Write Boundary / Rules prose already partitioned Check writing, Issue closure, acceptance-criteria edits, and owner-acceptance claims correctly, confirmed by grepping each phrase before writing the test.

**Commands run:**
- `go test ./internal/init/...` — pass.
- `go test ./internal/init/... -run 'TestProjectGuidanceTemplatesMirrorLiveGuidance|TestProjectTemplatesUseCurrentWorkflow|TestLiveSourcesDoNotReferenceGenericAuditSkill|TestScaffoldedProjectDoesNotReferenceGenericAuditSkill|TestStaleReferenceCheckExcludesHistoricalReleaseRecords' -v` — all pass.
- `make build && make test` — pass, all packages.

**Limitations:** the acceptance-criteria-edit ownership check anchors on `## Done When` (V2's Task template heading) rather than the literal phrase "acceptance criteria", since V2 renamed that field; this is a naming judgment call, not a verified-by-execution fact. No owner-facing UI or CLI surface exists yet to exercise this routing; verification is text-level only, as the task scope requires (routing stays inactive).
