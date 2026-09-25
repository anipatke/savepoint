---
id: E28-parallel-worktree-launch/T002-runtime-manifests-and-scope
status: planned
objective: Persist ignored coordinator manifests and per-worktree launch scope without changing router ownership.
depends_on:
  - E28-parallel-worktree-launch/T001-worktree-service
complexity_tier: high
complexity_reason: Runtime state must survive restarts, reconcile with Git, and never leak into branch integration.
---

# T002: Runtime Manifests and Scope

## Problem

Parallel sessions need durable local coordination and task routing without modifying the shared router in every branch.

## Context Files

- `internal/worktree/manifest.go`
- `internal/worktree/manifest_test.go`
- `internal/worktree/service.go`
- `internal/launcher/prompt.go`
- `internal/launcher/prompt_test.go`
- `templates/project/AGENTS.md`

## Acceptance Criteria

- [ ] The primary checkout records group, base commit, task, branch, path, scope, and launch result in an ignored manifest.
- [ ] Each worktree receives an ignored `launch-scope.yml` naming the selected task and declaring router and write boundaries.
- [ ] Runtime writes are atomic and tolerate an absent runtime directory.
- [ ] Manifest loading distinguishes active, stale, missing, and mismatched Git state without mutating it.
- [ ] Parallel builder prompts direct the agent to the launch scope, its own task lifecycle, and no router edits.
- [ ] Runtime files are excluded from commits and worktree integration by documented policy and tests.

## Implementation Plan

- [ ] Define versioned coordinator and launch-scope schemas.
- [ ] Add atomic runtime read/write helpers inside the worktree boundary.
- [ ] Write launch scope after worktree creation and before agent dispatch.
- [ ] Add manifest reconciliation against `git worktree list` and branch identity.
- [ ] Extend deterministic launcher prompts and generated AGENTS routing guidance for scoped worktree sessions.

## Context Log

Pending.
