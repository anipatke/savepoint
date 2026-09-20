---
id: E29-parallel-mode-guidance-and-hardening/T001-doctor-parallel-diagnostics
status: planned
objective: Diagnose advanced parallel configuration, planning metadata, and runtime Git drift without mutation.
depends_on:
  - E28-parallel-worktree-launch/T004-worktree-status-and-cleanup
complexity_tier: high
complexity_reason: Diagnostics reconcile config, tasks, manifests, filesystem paths, branches, and live Git worktrees.
---

# T001: Doctor Parallel Diagnostics

## Problem

Advanced users need precise explanations when planned groups or runtime worktrees are invalid, stale, or unsafe to clean.

## Context Files

- `internal/doctor/checks.go`
- `internal/doctor/checks_test.go`
- `internal/doctor/interfaces.go`
- `internal/doctor/repairs.go`
- `internal/doctor/report.go`
- `internal/data/parallel.go`
- `internal/worktree/manifest.go`

## Acceptance Criteria

- [ ] Doctor reports invalid advanced enablement, limits, execution metadata, and scope overlaps only when relevant.
- [ ] Runtime checks report stale manifests, missing paths, branch mismatches, dirty worktrees, and orphaned live worktrees.
- [ ] Disabled mode does not make old tasks or leftover runtime state a blocking project error.
- [ ] Suggestions are explicit and non-destructive; doctor never removes worktrees or branches.
- [ ] Git inspection is injected and tests require no global repository state.

## Implementation Plan

- [ ] Extend doctor interfaces with read-only parallel eligibility and Git runtime inspection.
- [ ] Add typed problem codes and repair suggestions for planning and runtime failures.
- [ ] Gate config/task findings by advanced-mode relevance while still surfacing stale runtime information as advisory.
- [ ] Add temporary-repository and fake-interface tests for every mismatch class.
- [ ] Verify report ordering and existing doctor output remain stable when the feature is absent.

## Context Log

Pending.
