---
id: E46-agent-workflow-assets/T002-take-a-rough-idea-without-demanding-a-document
title: Take a rough idea without demanding a document
status: done
objective: Add the savepoint-idea planner skill and its shipped copy, with a write boundary that stops at Idea.md and the routing handoff.
depends_on: []
complexity_tier: medium
complexity_reason: Two mirrored skill files plus a contract test, but the read/write boundary and Idea sections are real contract decisions.
---

# T002: Take a rough idea without demanding a document

## Problem

V1's entry point is `savepoint-draft-prd`, and a PRD is a document the user is expected to arrive with or sit through. V2's entry point accepts one rough sentence and refines it by asking, which means the skill has to be explicit about two things V1 never had to be: what it is allowed to write, and when it must stop and ask the owner instead of deciding.

The write boundary is the whole point of the skill. An idea session that quietly starts naming components, or writes an Objective because the shape seemed obvious, has skipped the design conversation and produced a backlog nobody agreed to. Architecture, Objectives, Tasks, and code are all out of bounds here; the only outputs are `.savepoint/Idea.md` and the routing handoff that sends the work on to design.

Existing-project evidence is allowed but targeted: reading a repository to understand what already exists is legitimate, reading it to design a solution is not.

## Context Files

- `agent-skills/savepoint-idea/SKILL.md`
- `templates/project/agent-skills/savepoint-idea/SKILL.md`
- `agent-skills/savepoint-draft-prd/SKILL.md`
- `internal/init/skill_validation_test.go`
- `internal/init/agent_skills_test.go`
- `internal/init/template_freshness_test.go`
- `.savepoint/releases/v2/v2-Design.md`

## Acceptance Criteria

- [x] `agent-skills/savepoint-idea/SKILL.md` exists with frontmatter `name: savepoint-idea` and a non-empty `description`, and `templates/project/agent-skills/savepoint-idea/SKILL.md` is byte-identical to it.
- [x] It has non-empty `## Purpose`, `## Trigger`, `## Read`, `## Workflow`, and `## Rules` sections, so it passes the existing structure validation in both trees.
- [x] Its trigger is router `state: idea`, stated as a V2 routing state and not as a V1 phase name.
- [x] Its `## Read` names router, `.savepoint/Idea.md`, stated user intent, and targeted existing-project evidence — and nothing else.
- [x] Its write boundary names exactly two outputs: `.savepoint/Idea.md` and the routing handoff. Design, Guardrails, Objectives, Tasks, Checks, Issues, and production code are named as forbidden.
- [x] It carries the Idea artifact template with the sections Intent, User, Core Experience, Scope, Out of Scope, Success Criteria.
- [x] It states that a single rough sentence is a valid starting input and that no prepared PRD, research document, or completed template is required before the conversation starts.
- [x] It states the escalation rule: material product uncertainty goes to the owner as a question; the skill does not resolve product choices by inference.
- [x] It names the handoff target `savepoint-design` and states that the idea session does not detail any Objective itself.
- [x] `agent-skills/savepoint-draft-prd/SKILL.md` and every other V1 skill are unchanged and still present in both trees.
- [x] `internal/init/agent_skills_test.go` asserts the skill's write boundary, the forbidden-output list, the Idea template sections, and live/template parity.
- [x] `TestProjectGuidanceTemplatesMirrorLiveGuidance` passes with the new skill counted in both trees.

## Implementation Plan

- [x] Read design section 7 for the Idea template sections and section 8 for the planner read/write/forbidden row.
- [x] Write `agent-skills/savepoint-idea/SKILL.md` with Purpose, Trigger, Read, Workflow, Rules, and the Idea artifact template.
- [x] State the write boundary and forbidden outputs as explicit rules, not as prose implication, so the contract test can assert them.
- [x] Copy the file verbatim to `templates/project/agent-skills/savepoint-idea/SKILL.md`.
- [x] Extend `internal/init/agent_skills_test.go` with the savepoint-idea contract case.
- [x] Run `go test ./internal/init/...`, then `make build && make test`.

## Context Log

**Files read:** `agent-skills/savepoint-idea/SKILL.md` (new), `templates/project/agent-skills/savepoint-idea/SKILL.md` (new), `agent-skills/savepoint-draft-prd/SKILL.md`, `agent-skills/savepoint-system-design/SKILL.md`, `agent-skills/savepoint-audit-task/SKILL.md`, `internal/init/skill_validation_test.go`, `internal/init/agent_skills_test.go`, `internal/init/template_freshness_test.go`, `.savepoint/releases/v2/v2-Design.md` (sections 7-8), `E46-Detail.md`.

**Files written:** `agent-skills/savepoint-idea/SKILL.md`, `templates/project/agent-skills/savepoint-idea/SKILL.md` (byte-identical copy), `internal/init/agent_skills_test.go` (added `TestSavepointIdeaSkillReadBoundary`, `TestSavepointIdeaSkillWriteBoundary`, `TestSavepointIdeaSkillHasArtifactTemplate`, `TestSavepointIdeaSkillAcceptsRoughInputAndEscalates`, `TestSavepointIdeaSkillLiveAndTemplateMatch`).

**Quality gates:** `go test ./internal/init/...` — pass. `make build && make test` — pass (all packages).

**Health check:** No `.savepoint/Health-Check.md` in this project; step skipped per skill instructions.

No V1 skill files were touched; `git status` shows no diff under `agent-skills/savepoint-draft-prd` or `templates/project/agent-skills/savepoint-draft-prd`.
