---
id: E29-parallel-mode-guidance-and-hardening/T002-scaffold-parallel-mode
status: planned
objective: Scaffold disabled advanced configuration, ignored runtime paths, and upgrade-safe managed guidance.
depends_on:
  - E27-parallel-planning-contracts/T004-parallel-planning-skill
  - E28-parallel-worktree-launch/T002-runtime-manifests-and-scope
complexity_tier: medium
complexity_reason: Init and upgrade assets must add discoverability without enabling features or overwriting user-owned files.
---

# T002: Scaffold Parallel Mode

## Problem

Generated projects need discoverable advanced settings and safe ignore rules while preserving current defaults and user content.

## Context Files

- `templates/project/config.yml`
- `templates/project/AGENTS.md`
- `templates/project/.gitignore`
- `internal/init/scaffold_test.go`
- `internal/init/upgrade_test.go`
- `main.go`

## Acceptance Criteria

- [ ] New projects include a documented parallel-worktree example with `enabled: false`.
- [ ] Runtime manifests, launch scopes, and managed worktree paths are ignored by Git.
- [ ] Init never creates a worktree, runtime manifest, or enabled advanced configuration.
- [ ] Upgrade refreshes package-owned guidance without changing project state or enabling advanced mode.
- [ ] Existing user-owned ignore and agent-guide content remains preserved.

## Implementation Plan

- [ ] Add the disabled advanced block to the scaffold config template.
- [ ] Add narrow runtime/worktree ignore entries to the project template.
- [ ] Update managed AGENTS guidance for scoped worktree launches.
- [ ] Extend embedded asset wiring only where required by existing init/upgrade patterns.
- [ ] Add init and upgrade regression tests for disabled defaults and content preservation.

## Context Log

Pending.
