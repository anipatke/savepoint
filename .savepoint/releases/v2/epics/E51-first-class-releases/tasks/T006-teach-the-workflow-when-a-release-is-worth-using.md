---
id: E51-first-class-releases/T006-teach-the-workflow-when-a-release-is-worth-using
status: in_progress
stage: build
objective: Teach planners and checkers to create and verify optional Releases while fresh projects remain free of mandatory ceremony.
depends_on:
  - E51-first-class-releases/T002-require-release-integration-evidence-and-owner-acceptance
complexity_tier: medium
complexity_reason: Reconciles three public skills and shipped copies with a new optional planning and checking boundary.
---

# T006: Teach the workflow when a Release is worth using

## Problem

A file schema alone does not make Releases usable. The V2 skills must explain when a delivery boundary is justified, how to create/link it without rebuilding hierarchy, and how a fresh checker verifies it. At the same time, init must not synthesize a Release for every weekend project.

## Context Files

- `.savepoint/releases/v2/epics/E51-first-class-releases/E51-Detail.md`
- `agent-skills/savepoint-idea/SKILL.md`
- `agent-skills/savepoint-design/SKILL.md`
- `agent-skills/savepoint-check/SKILL.md`
- `templates/project-v2/agent-skills/savepoint-idea/SKILL.md`
- `templates/project-v2/agent-skills/savepoint-design/SKILL.md`
- `templates/project-v2/agent-skills/savepoint-check/SKILL.md`
- `templates/project-v2/AGENTS.md`
- `internal/init/agent_skills_test.go`
- `internal/init/skill_validation_test.go`
- `internal/init/template_freshness_test.go`
- `internal/init/v2_scaffold_test.go`

## Acceptance Criteria

- [x] Idea guidance treats Release as optional and asks for one only when the owner needs a navigable delivery/package promise across Objectives.
- [x] Design guidance defines the Release outcome/success conditions, allocates a stable `R###`, and links Objectives through one `release` field without nesting files.
- [x] Check guidance supports `scope.kind: release`, reviews cross-Objective integration, reuses Issues, and records no owner acceptance on the owner's behalf.
- [x] Guidance states that Release `done` requires current CLEAR integration evidence plus owner acceptance and does not mean published or deployed.
- [x] Projects without a Release continue through Idea → Design → Task → Check with no missing-record error or extra phase.
- [x] A fresh V2 scaffold contains no default Release record and no placeholder delivery promise.
- [x] Canonical and shipped skill copies are byte-identical after the change and pass existing structural/reference validation.
- [x] No fifth public skill/phase, Release-owned Task list, release PRD helper document, publishing workflow, or duplicated Check method is introduced.
- [x] Scaffold and upgrade tests prove existing user-owned files remain untouched while the updated managed assets install normally.

## Implementation Plan

- [x] Add a concise optional-Release decision boundary to Idea.
- [x] Add Release creation/linking and Design reconciliation rules to Design.
- [x] Add Release-scope verification, Issue routing, and owner-handoff rules to Check.
- [x] Update the shipped copies and managed AGENTS guidance without adding another public phase.
- [x] Assert canonical/template parity and reject stale pre-E51 optional-string wording.
- [x] Prove fresh init creates no Release record and existing no-Release projects remain valid.
- [x] Run focused init asset/skill validation tests.

## Context Log

- Read: `.savepoint/router.md`, `.savepoint/releases/v2/epics/E51-first-class-releases/E51-Detail.md`, this task, `.savepoint/Guardrails.md`, the three canonical V2 skills, their three `templates/project-v2` copies, `templates/project-v2/AGENTS.md`, and the listed init validation/scaffold tests. `.savepoint/Health-Check.md` is absent, so no Quick evidence block was required.
- Targeted verification read/edit: `internal/init/provenance_contract_test.go`, whose typed Check-template fixture was updated for the new `release` scope placeholder.
- Edited: `agent-skills/savepoint-idea/SKILL.md`, `agent-skills/savepoint-design/SKILL.md`, `agent-skills/savepoint-check/SKILL.md`, matching `templates/project-v2/agent-skills/*` copies, `templates/project-v2/AGENTS.md`, `internal/init/agent_skills_test.go`, `internal/init/provenance_contract_test.go`, `internal/init/v2_scaffold_test.go`, and this task file.
- Evidence: `GOCACHE=/tmp/savepoint-gocache go test ./internal/init` passed; `GOCACHE=/tmp/savepoint-gocache make build && GOCACHE=/tmp/savepoint-gocache make test` passed for all packages. Existing parity, structural/reference, user-file-preservation, upgrade, and fresh-scaffold tests passed, along with the new Release guidance and no-Release scaffold checks.
