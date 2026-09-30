---
id: T-052
title: Keep health history safe and compact
objective: O-027
status: done
depends_on: [{task: T-051, requires: clear}]
owner_validation:
    required: false
    accepted_check: ""
planned_by: {role: planner, session: planning-o027-20260929}
check_waiver:
    task: T-052
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-09-30T20:58:01Z"
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

**Commands run:** `go test -count=1 -race ./internal/codehealth` (pass); `GOOS=windows go vet ./internal/codehealth` (pass, cross-vet only, not a native Windows run); `make build && make test-fast` (exit 0).

**Files changed:** `internal/codehealth/storage.go` (new), `internal/codehealth/storage_test.go` (new), this Task's `status`/`stage`.

**Files read:** the Task's Context Files (Objective.md, Design.md sections 2 and 8-11, Guardrails.md, model.go). `storage.go` and `storage_test.go` did not exist; both were created.

**Extra reads (outside Context Files), needed to reuse the model's validators and fixtures:** `internal/codehealth/config.go`, `snapshot.go`, `identity.go`, `primitives.go`, `errors.go` (validation, `DecodeConfig`/`DecodeSnapshot`, ID derivation, sentinels); the fixture helpers in `model_test.go` (`validSnapshot`, `validConfig`, `reseal`). `AGENTS.md` and `agent-skills/savepoint-task/SKILL.md` were read as workflow.

**Design:** `Store` (project path only) keeps config at `.savepoint/health/config.json` and snapshots at `.savepoint/health/snapshots/<64-hex>.json`. The file name is the snapshot identity. Every directory step is `Lstat`ed and a symlink or non-directory is refused (`ErrUnsafePath`). Only `health/` and `snapshots/` are ever created, never `.savepoint` (`ErrNotProject`). Snapshot writes go to a temp file and `os.Link` into place, so they are atomic and create-only. Config replace is temp + rename. Reads are bounded to 1 MiB. `Prune` only removes manual snapshots beyond the newest ten and never runs from save or load.

**Per-criterion outcomes:**
1. Confinement and rejection: met. `TestSymlinkEscapesAreRefused` (health dir, snapshots dir, config file, snapshot file), `TestSaveSnapshotRefusesInvalid` (traversing ID `sha256:../../escape`), `TestLoadConfigRejectsBadFiles` (malformed, unknown field, unsupported version, trailing data), `TestLoadSnapshotsRejectsBadFiles` (malformed, unsupported version, name/content mismatch, non-snapshot name, uppercase digest), `TestLoadSnapshotsRejectsDuplicateContentUnderAnotherName`, `TestLoadSnapshotsRejectsDirectoryEntry`. Paths are built with `filepath.Join`. Duplicate identities cannot occur under a valid name (identity and file name are 1:1), so a duplicate surfaces as a name mismatch or the explicit duplicate error.
2. Immutable, collision-safe, atomic, retry-safe: met. `TestSaveSnapshotIsIdempotent` (repeat reports unchanged, mtime unchanged), `TestSaveSnapshotConflictKeepsExistingFile` (corrupt, different snapshot, unsupported version: refused, file byte-identical), `TestInterruptedWriteLeftoverIsIgnoredAndRetrySucceeds`, `TestFailedSnapshotWriteLeavesStoreValid` (skips on Windows/root).
3. Deterministic order; series kept apart: met. `TestLoadSnapshotsOrderIsDeterministic` (time, then identity on a tie), `TestIncompatibleSeriesStayApart`.
4. Official never pruned; newest ten manual kept; exact reporting: met. `TestPruneRetentionBoundaries` (0/9/10/11/14; plan equals what `Prune` removed, oldest removed), `TestPruneNeverTouchesOfficialSnapshots` (12 official plus 13 manual), `TestPruneTiesUseIdentity`.
5. No pruning as a side effect; failures keep valid records: met. `TestSavingNeverPrunes`, `TestPruneRefusesWhenHistoryIsInvalid`, `TestPrunePartialFailureKeepsRemainingValid` and `TestConfigFailedReplaceKeepsOldFile` (skip on Windows/root), `TestPruneIsRepeatable`.
6. Test coverage list: met. Clean and dirty directories (`TestLoadSnapshotsAbsentDirectories`, `TestSaveRefusesNonProjectWithoutScaffolding`, stray-file and leftover tests), absent directories, malformed/unsupported files, conflicting and interrupted writes, symlinks, 9/10/11 boundaries, ties, mixed origins, repeat execution. Also config round trip, unchanged repeat, atomic replace with no leftovers, and refused invalid save.

**Limitations:**
- Not run natively on Windows (CFG-03). Symlink and permission tests skip there or when run as root; the `os.Link` create-only write assumes a filesystem with hard links (NTFS and ext4 yes; FAT-like mounts would fail the write clearly rather than overwrite).
- A single invalid file fails the whole load (and blocks pruning) by design; there is no partial-load or repair mode.
- No cross-process locking. Two concurrent writers of the same identity are safe (link is exclusive); a concurrent `Prune` and save are not coordinated.
- The Task was not independently checked. No owner Task-check waiver has been recorded.
- `.claude/settings.local.json` was already modified before this Task and is unrelated.

## Drift Notes

Return to planning if safe immutable writes require changing shared project-file writers or introducing deletion during collection.
