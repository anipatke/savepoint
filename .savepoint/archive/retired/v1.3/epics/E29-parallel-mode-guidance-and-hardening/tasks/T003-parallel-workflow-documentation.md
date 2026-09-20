---
id: E29-parallel-mode-guidance-and-hardening/T003-parallel-workflow-documentation
status: planned
objective: Document the complete advanced planning, launch, integration, toggle, and cleanup workflow.
depends_on:
  - E29-parallel-mode-guidance-and-hardening/T001-doctor-parallel-diagnostics
  - E29-parallel-mode-guidance-and-hardening/T002-scaffold-parallel-mode
complexity_tier: medium
complexity_reason: Guidance spans user decisions, agent boundaries, Git recovery, and compatibility guarantees across several artifacts.
---

# T003: Parallel Workflow Documentation

## Problem

The advanced workflow has Git and ownership constraints that must be explicit before users trust parallel launches.

## Context Files

- `README.md`
- `templates/project/AGENTS.md`
- `agent-skills/savepoint-build-task/SKILL.md`
- `templates/project/agent-skills/savepoint-build-task/SKILL.md`
- `.savepoint/releases/v1.3/v1.3-PRD.md`

## Acceptance Criteria

- [ ] Documentation explains opt-in planning metadata, eligibility, clean-checkout preflight, and one-worktree-per-task launch.
- [ ] Worker guidance names the launch-scope override, router ownership, task lifecycle, and exact write boundary.
- [ ] Integration remains user-controlled and merge conflicts are described as a planning-boundary failure.
- [ ] Toggle-off behavior and explicit safe cleanup are documented without implying deletion or monitoring.
- [ ] Default sequential instructions remain first and sufficient for users who never enable advanced mode.

## Implementation Plan

- [ ] Add a concise advanced-mode guide to the README after the sequential launcher workflow.
- [ ] Update build-task skills with the scoped worktree launch contract.
- [ ] Keep normal router-first behavior unchanged outside a generated launch scope.
- [ ] Document integration and cleanup commands conceptually without automating them.
- [ ] Verify source and template skill parity.

## Context Log

Pending.
