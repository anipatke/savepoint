---
id: E46-agent-workflow-assets/T005-execute-a-plan-and-say-what-actually-happened
title: Execute a plan and say what actually happened
status: done
objective: Add the savepoint-task executor skill so one Task is built within its boundaries, with logged extra reads, recorded evidence, and no self-granted clearance.
depends_on:
    - E46-agent-workflow-assets/T004-give-every-task-a-title-a-person-can-read
    - E46-agent-workflow-assets/T001-write-the-shared-checking-method-once
complexity_tier: medium
complexity_reason: One mirrored skill pair carrying the lifecycle, replan handoff, and evidence contract the Check later depends on.
---

# T005: Execute a plan and say what actually happened

## Problem

The executor is the role with the most opportunity to quietly lie. It can widen scope and call it necessary, redesign around an inconvenient plan and not mention it, tick acceptance criteria it did not verify, or mark its own work cleared. V1's `savepoint-build-task` already blocks the last one by handing completion to the user; V2 needs the rest written down, because completion now depends on recorded evidence that a fresh checker will read as fact.

Two mechanics carry most of the weight. Extra reads are allowed but logged: the plan's Context Files are the budget, and anything read beyond them is recorded so replanning frequency and context growth are measurable instead of anecdotal. And a materially invalid plan produces `REPLAN REQUIRED` with preserved partial work — not an improvised redesign. Stopping honestly is the success case there.

The lifecycle is the other half: planned → in_progress at stage build, then test, then audit, where audit means ready for Check and never means passed.

## Context Files

- `agent-skills/savepoint-task/SKILL.md`
- `templates/project/agent-skills/savepoint-task/SKILL.md`
- `agent-skills/savepoint-build-task/SKILL.md`
- `agent-skills/references/check-method.md`
- `internal/init/agent_skills_test.go`
- `internal/data/task_v2.go`
- `internal/data/gate_v2.go`
- `.savepoint/releases/v2/v2-Design.md`
- `.savepoint/Guardrails.md`

## Acceptance Criteria

- [x] `agent-skills/savepoint-task/SKILL.md` exists with frontmatter `name: savepoint-task`, passes the existing structure validation, and has a byte-identical shipped copy.
- [x] Its trigger is router `state: task`, and its `## Read` names router, the Task, the owning Objective's boundaries, the Task's Context Files, and applicable policy — and states that this is the read budget.
- [x] Its write boundary names scoped implementation, recorded evidence, lifecycle progress, and the replan handoff; it forbids editing acceptance criteria, writing a Check, closing an Issue, or claiming clearance or owner acceptance.
- [x] It states the lifecycle transitions it may make: `planned` → `in_progress` with `stage: build`, then `build` → `test` → `audit`, with `audit` defined as ready for Check and explicitly not as passed.
- [x] It states that starting requires satisfied Task dependencies and a ready owning Objective, and that a blocked start is reported rather than worked around.
- [x] It requires reads beyond the Task's Context Files to be recorded with what was read and why.
- [x] It defines `REPLAN REQUIRED`: the executor sets the replan reason with handoff evidence, keeps current status and stage, preserves partial work, and stops for the planner — and never silently redesigns.
- [x] It requires recorded technical evidence at handoff: per-criterion outcome, the named commands run including `make build && make test`, files read and changed, and stated limitations.
- [x] It requires the handoff to go to a fresh `savepoint-check` session and states that the executor's own session can never be that Check.
- [x] It treats `STYLE` guardrail rules as advisory and references guardrail IDs rather than restating rule prose.
- [x] `internal/init/agent_skills_test.go` asserts the read budget, the write boundary and forbidden claims, the stage sequence, the replan contract, the evidence requirement, and live/template parity.
- [x] `agent-skills/savepoint-build-task/SKILL.md` is unchanged and still present in both trees.

## Implementation Plan

- [x] Read design sections 4 and 8 for the lifecycle table and the executor role row, and `internal/data/gate_v2.go` for the start and completion gates the wording must match.
- [x] Write `agent-skills/savepoint-task/SKILL.md` with Purpose, Trigger, Read, Workflow, Rules.
- [x] Write the lifecycle steps and the `audit` meaning explicitly, matching the recorded stage vocabulary rather than inventing UI words.
- [x] Write the extra-read logging rule and the `REPLAN REQUIRED` contract, including what must be preserved.
- [x] Write the evidence and handoff requirements, pointing at `agent-skills/references/check-method.md` for what the checker will do with them.
- [x] Mirror the file to the template tree.
- [x] Extend `internal/init/agent_skills_test.go` with the savepoint-task contract case.
- [x] Run `go test ./internal/init/...`, then `make build && make test`.

## Context Log

**Files read (Context Files budget):** `agent-skills/savepoint-build-task/SKILL.md`, `agent-skills/references/check-method.md`, `internal/init/agent_skills_test.go`, `internal/data/task_v2.go`, `internal/data/gate_v2.go`, `.savepoint/releases/v2/v2-Design.md`, `.savepoint/Guardrails.md`. `agent-skills/savepoint-task/SKILL.md` and `templates/project/agent-skills/savepoint-task/SKILL.md` did not exist yet — this task creates them.

**Extra reads beyond the Context Files budget, logged with reason:**
- `.savepoint/router.md`, `.savepoint/releases/v2/epics/E46-agent-workflow-assets/E46-Detail.md` — required read order per AGENTS.md/router before touching any task.
- `AGENTS.md` — confirmed workflow/terminology rules and the `savepoint-build-task` skill governing this session (router `state: task-building` is V1; `savepoint-task` is the V2 asset this task ships).
- `agent-skills/savepoint-audit-task/SKILL.md` and `agent-skills/savepoint-design/SKILL.md` — the nearest sibling V2 skills (design already built under this epic's T003/T004), read to match established section structure, tone, and how they cite `agent-skills/references/check-method.md` and Guardrail IDs instead of restating them.
- `internal/init/skill_validation_test.go` — needed the exact `skillRoots()`, `frontmatterField`, and `sectionBody` test helpers before extending `agent_skills_test.go`.
- `internal/init/template_freshness_test.go` (`TestProjectGuidanceTemplatesMirrorLiveGuidance`) — confirmed the live/template skill-count and parity check is dynamic (discovers `savepoint-*` directories), so adding `savepoint-task` needed no separate count update and would not silently break that test.

**Files written:**
- `agent-skills/savepoint-task/SKILL.md` (new) and `templates/project/agent-skills/savepoint-task/SKILL.md` (new, byte-identical mirror).
- `internal/init/agent_skills_test.go` — added `TestSavepointTaskSkillPassesStructureValidation`, `TestSavepointTaskSkillTrigger`, `TestSavepointTaskSkillReadBudget`, `TestSavepointTaskSkillWriteBoundaryAndForbiddenClaims`, `TestSavepointTaskSkillLifecycleStages`, `TestSavepointTaskSkillBlockedStartReported`, `TestSavepointTaskSkillExtraReadsLogged`, `TestSavepointTaskSkillReplanRequiredContract`, `TestSavepointTaskSkillEvidenceRequirement`, `TestSavepointTaskSkillFreshCheckHandoff`, `TestSavepointTaskSkillStyleAdvisoryRule`, `TestSavepointTaskSkillLiveAndTemplateMatch`, `TestSavepointBuildTaskSkillPresentInBothTrees`.

**Per-criterion outcome:** every acceptance criterion above is backed by a passing assertion in `internal/init/agent_skills_test.go` (see test names above) except the structure/frontmatter/byte-parity criterion, which is covered by both the new `TestSavepointTaskSkillPassesStructureValidation`/`TestSavepointTaskSkillLiveAndTemplateMatch` and the pre-existing generic `TestSavepointSkillsHaveValidFrontmatter`, `TestSavepointSkillsHaveNonEmptyTriggerAndWorkflow`, and `TestProjectGuidanceTemplatesMirrorLiveGuidance` (auto-discovers `savepoint-*` dirs, so it now also covers `savepoint-task`). `savepoint-build-task/SKILL.md` unchanged: `git diff --stat` for both `agent-skills/savepoint-build-task` and `templates/project/agent-skills/savepoint-build-task` is empty.

**Commands run:**
- `go test ./internal/init/... -run 'SavepointTask|SavepointBuildTask' -v` — all 13 new/targeted tests pass.
- `go test ./internal/init/...` — full package passes (no regressions from the new directory or count-based parity test).
- `make build && make test` — build succeeds; full suite passes across all packages.

**Limitations:** No independent Check has been run against this Task; per the skill just written, this session cannot be that Check. `internal/data/gate_v2.go` and `internal/data/task_v2.go` are read-only context here — this task does not wire the V2 skill into any V2 runtime code (none exists yet on this branch), only ships the skill guidance and its contract tests, matching the epic's stated scope (skill contracts, matching scaffold copies, shared method).
