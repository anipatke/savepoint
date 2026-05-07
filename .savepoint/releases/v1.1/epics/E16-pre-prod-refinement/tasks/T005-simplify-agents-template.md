---
id: E16-pre-prod-refinement/T005-simplify-agents-template
status: in_progress
stage: test
objective: Reduce scaffolded `AGENTS.md` line count and context overhead while preserving Savepoint operating rules
depends_on: []
---

# T005: Simplify Scaffolded Agents Guide

## Problem

The scaffolded `templates/project/AGENTS.md` guide is correct but verbose. It repeats router-read instructions, separates closely related task and context-budget rules, and spends extra lines on audit structure that can be stated more compactly without changing behavior.

This increases context cost for every generated Savepoint project because agents must read the guide at session start.

## Context Files

- `templates/project/AGENTS.md` - scaffolded agent guide to simplify
- `.savepoint/router.md` - current state wording that overlaps with the template
- `agent-skills/savepoint-build-task/SKILL.md` - implementation-loop rules that must remain compatible with the template

## Acceptance Criteria

- [x] `templates/project/AGENTS.md` is shorter and removes redundant wording from workflow, implementation, audit, build, and context-budget sections
- [x] The template still preserves required skill activation, task status, `stage`, implementation, drift, audit, code style, codebase map, build gate, and CLI rules
- [x] The implementation/context-budget rules still prohibit broad exploration for task builds unless explicitly required
- [x] The audit instructions still require a fresh audit agent, one audit file, visible `## Main Findings` and `## Code Style Review` sections, and file-specific proposed-change blocks under `## Proposed Changes`
- [x] The resulting guide remains clear enough for generated projects without relying on repository-specific Epic 16 context
- [x] `make build && make test` passes

## Implementation Plan

- [x] Collapse `## Workflow` and `## Skill Activation` into a shorter startup section
- [x] Tighten `## Task Status` into the parse-critical status and stage rules
- [x] Merge context-budget constraints into the implementation read step and remove the separate context-budget section
- [x] Compress audit handoff and audit file structure into one audit section
- [x] Replace the fenced build command with a one-line build gate while preserving `make build && make test`
- [x] Keep code style and CLI rules concise without weakening meaning
- [x] Run `make build && make test`

## Context Log

- Files read: `.savepoint/router.md`, `agent-skills/savepoint-build-task/SKILL.md`, `.savepoint/releases/v1.1/epics/E16-pre-prod-refinement/E16-Detail.md`, `.savepoint/releases/v1.1/epics/E16-pre-prod-refinement/tasks/T005-simplify-agents-template.md`, `templates/project/AGENTS.md`.
- Files edited: `.savepoint/releases/v1.1/epics/E16-pre-prod-refinement/tasks/T005-simplify-agents-template.md`, `templates/project/AGENTS.md`.
- Token estimate: low.
- Quality gates: `make build && make test` passed.
