---
id: T-053
title: Describe repository state without guessing
objective: O-027
status: planned
depends_on: [{task: T-051, requires: clear}]
owner_validation: {required: true}
planned_by: {role: planner, session: planning-o027-20260929}
---

# Describe repository state without guessing

## Outcome

Code Health derives a deterministic, privacy-safe repository identity for clean, dirty, divergent, shallow, unborn, and non-Git worktrees, including honest ancestor-based wording for how observations relate.

## User Check

Review the plain-language clean, dirty, commits-behind, ahead, diverged, unborn, and unavailable descriptions against small temporary repositories. Confirm that relevant untracked changes affect identity without exposing file content or private remote information.

## Done When

- Repository identity records commit when useful, a deterministic fingerprint of configured-scope tracked changes and relevant untracked inputs, and an explicit state when no useful commit exists.
- Fingerprints are stable across enumeration order and path separators, change when relevant content changes, exclude configured irrelevant paths, and never retain file contents, usernames, hostnames, remote URLs, or author/commit message text.
- Relationship calculation uses actual ancestry to distinguish same, behind, ahead, diverged, unavailable, shallow/unknown, and no-useful-commit cases; numeric “commits behind” wording appears only when proven.
- Git failures, missing Git, non-repository directories, subdirectories, worktrees, deleted files, symlinks, unusual filenames, and relevant untracked files produce deterministic results or named errors without panics.
- The implementation uses structured executable arguments without a shell and supports cancellation/timeouts at its boundary.
- Temporary-repository tests provide independent Git-command or file-hash oracles for clean, dirty tracked, relevant/irrelevant untracked, ancestor, divergence, unborn, worktree, and unavailable cases.

## Context Files

`.savepoint/objectives/O-027-code-health-model-and-history/Objective.md`; `.savepoint/Design.md`; `.savepoint/Guardrails.md`; `internal/codehealth/model.go`; `internal/codehealth/repository.go`; `internal/codehealth/repository_test.go`.

## Design References

Design sections 1, 8, 9, 11, and 12.

## Guardrails

ARCH-03, ARCH-04, CFG-02, CFG-03, DEP-01, TEST-01, TEST-02, TEST-04, TEST-08.

## Implementation Plan

1. Define a narrow command-runner boundary and derive Git state with explicit argument vectors from the project root.
2. Build the configured-scope fingerprint from normalized path, tracked state, and content digests for relevant tracked and untracked inputs.
3. Calculate observation relationships from merge-base/ancestry evidence and keep unknown states explicit.
4. Render deterministic owner-facing descriptions from typed repository relations without leaking Git metadata.
5. Add temporary-repository matrices, boundary failures, cancellation, and platform-aware filename/path cases.

## Boundaries

No snapshot filesystem persistence, provider execution, remote fetch, network access, branch mutation, Check wiring, TUI rendering, or generic VCS abstraction.

## Technical Verification

Run focused repository-state tests during iteration, then `make build && make test-fast`. A later Check follows `agent-skills/references/check-method.md`.

## Technical Evidence

Pending execution: record repository matrices, independent oracles, commands, changed files, and platform limitations.

## Drift Notes

Return to planning if truthful ancestry requires network access or if scope fingerprinting cannot avoid sensitive data retention.
