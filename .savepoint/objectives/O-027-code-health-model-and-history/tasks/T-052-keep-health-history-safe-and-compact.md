---
id: T-052
title: Keep health history safe and compact
objective: O-027
status: planned
depends_on: [{task: T-051, requires: clear}]
owner_validation: {required: false}
planned_by: {role: planner, session: planning-o027-20260929}
---

# Keep health history safe and compact

## Outcome

Code Health safely loads and stores versioned configuration and immutable snapshots under `.savepoint/health/`, preserves permanent Full Check evidence, and prunes manual history only through an explicit, deterministic maintenance operation.

## User Check

No separate owner validation is required. An optional Task Check may inspect persistence safety and retention behavior; if skipped, the owner records the required waiver when completing the Task.

## Done When

- Configuration and snapshot paths are confined beneath `.savepoint/health/`, use cross-platform path handling, and reject traversal, symlink escape, unknown schema versions, malformed files, duplicate snapshot identities, and content/filename identity mismatch.
- Snapshot creation is immutable, collision-safe, atomic, and safe to retry; an existing identical snapshot is unchanged, while differing content at the same identity is refused.
- Loading returns snapshots in deterministic order and preserves old compatible and incompatible series without silently joining them.
- Full Objective Check snapshots are never selected for pruning; explicit manual pruning keeps the newest ten manual snapshots and reports exactly what would be or was removed.
- Collection never prunes as a side effect, and failed writes or pruning leave previously valid records intact.
- Tests cover clean and dirty directories, absent directories, malformed and unsupported files, interrupted or conflicting writes, symlink/path escapes where supported, retention boundaries at 9/10/11 items, ties, mixed origins, and repeat execution.

## Context Files

`.savepoint/objectives/O-027-code-health-model-and-history/Objective.md`; `.savepoint/Design.md`; `.savepoint/Guardrails.md`; `internal/codehealth/model.go`; `internal/codehealth/storage.go`; `internal/codehealth/storage_test.go`.

## Design References

Design sections 1, 2, 8, 9, and 11.

## Guardrails

FS-01, FS-04, FS-05, FS-06, DATA-03, ARCH-03, CFG-01, CFG-02, TEST-01, TEST-02, TEST-04, TEST-08.

## Implementation Plan

1. Add confined config and snapshot path resolution rooted at the project path.
2. Implement strict versioned decoding and deterministic loading without repairing or rewriting malformed input.
3. Implement atomic create-only snapshot writes with canonical identity verification and idempotent identical-content handling.
4. Implement an explicit manual-retention planning/apply API that never includes Full Check snapshots and has no collection side effect.
5. Exercise failure timing, collision, retry, Windows path behavior, and retention boundaries in temporary-project tests.

## Boundaries

No provider execution, repository fingerprint calculation, classification or trend policy, CLI command, Check wiring, TUI rendering, background pruning, or raw report storage.

## Technical Verification

Run focused storage tests during iteration, then `make build && make test-fast`. A later Check follows `agent-skills/references/check-method.md`.

## Technical Evidence

Pending execution: record failure-path cases, retention matrices, commands, changed files, and limitations.

## Drift Notes

Return to planning if safe immutable writes require changing shared project-file writers or introducing deletion during collection.
