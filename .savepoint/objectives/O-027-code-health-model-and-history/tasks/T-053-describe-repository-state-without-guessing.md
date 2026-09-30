---
id: T-053
title: Describe repository state without guessing
objective: O-027
status: done
depends_on: [{task: T-051, requires: clear}]
owner_validation:
    required: true
    accepted_check: ""
planned_by: {role: planner, session: planning-o027-20260929}
check_waiver:
    task: T-053
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-09-30T21:07:17Z"
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

Executor session: T-053 build, 2026-10-01. Not a worktree lane (git-dir equals common dir). Audit means ready for a Check, not passed.

**Changed files:** `internal/codehealth/repository.go` (new), `internal/codehealth/repository_test.go` (new), this Task file (status/stage and evidence only).

**Extra reads (beyond Context Files):** `AGENTS.md`, `agent-skills/savepoint-task/SKILL.md` (workflow); `internal/codehealth/errors.go`, `primitives.go`, `config.go`, `snapshot.go`, `identity.go`, `identity_test.go`, `storage.go` (grep only) to reuse `RepositoryIdentity`, digest/path validators, the `.savepoint/health` constants and pattern rules. The `savepoint-task` skill was not registered with the Skill tool, so it was read directly per AGENTS.md. `repository.go` was listed as a Context File but did not exist; the plan creates it, so this was not treated as a material gap.

**Design as built:** `CommandRunner` boundary (argument vectors, no shell; `GitRunner` pins `LC_ALL=C`, `GIT_TERMINAL_PROMPT=0`, `GIT_OPTIONAL_LOCKS=0`, discards stderr, caps output at 32 MiB). `ObserveRepository` returns `Observation{Identity, Shallow}` using `rev-parse`, `ls-files -z`, and `diff --relative --name-only -z HEAD`. Fingerprint is sha256 over sorted, NUL-delimited, length-prefixed entries of working-tree content digests (symlinks hash their target without following; directories and special files are skipped; files Git lists but missing on disk are skipped). `.savepoint/health/**` is always excluded so storing a snapshot cannot change the state it describes. `Relate` uses `cat-file -e`, `merge-base --is-ancestor`, and `rev-list --left-right --count`; counts are set only when not shallow. Wording comes from `Describe`/`DescribeRepositoryError`. Whole observation and comparison use a 30s timeout and honour cancellation.

**Per-criterion outcomes** (all in `internal/codehealth/repository_test.go`):
1. Identity with commit, fingerprint, explicit unborn state: `TestObserveCleanMatchesOracles`, `TestObserveUnborn`, `TestObserveDirtyTracked`; every observation is also checked against `RepositoryIdentity.validate`.
2. Fingerprint stability/sensitivity/exclusion/no retention: `TestFingerprintIgnoresEnumerationOrderAndLocation`, `TestObserveDirtyTracked` (change and revert), `TestObserveUntrackedRelevanceAndScope`, `TestObserveExcludesStoredHealthSnapshots`, `TestObserveLeaksNothingSensitive` (remote credentials, file content, author, message, path absent from output), `TestInputScopeMatching`.
3. Ancestry relations, numeric wording only when proven: `TestRelateAgainstGitOracles` (same, same-with-changed-inputs, behind, ahead, diverged, no-commit, missing object is unknown not diverged), `TestRelateInShallowClone` (ancestry proven but no count; beyond boundary unknown), `TestRelateFailurePaths`, `TestDescriptionsAreDeterministicPlainWords`.
4. Failure and platform cases without panics: `TestObserveBoundaryFailures` (non-repo, missing git, cancelled, invalid scope, unreadable input), `TestObserveRunnerFailuresAreNamed` (runner error, timeout, bad exit code; runner detail not echoed), `TestObserveFromSubdirectoryStaysInsideIt`, `TestObserveLinkedWorktree`, `TestObserveDeletedFile`, `TestObserveSymlinkHashesTargetWithoutFollowing` (incl. dangling), `TestObserveUnusualFilenames` (spaces, unicode, leading dash, quotes, tab, newline off Windows).
5. Structured arguments, no shell, cancellation/timeouts: `TestRunnerBoundaryUsesArgumentVectors`, cancellation and timeout subtests above.
6. Independent oracles: real `git rev-parse`, `git rev-list --count`, and a test-side recomputation of the fingerprint from file bytes (`oracleFingerprint`) over clean, dirty, untracked, ancestor, divergence, unborn, worktree, shallow, and unavailable cases. No test was skipped on this host (Linux).

**Commands run:** `go vet ./internal/codehealth`; `GOOS=windows go vet ./internal/codehealth` (compiles for Windows); `go test -count=1 ./internal/codehealth` (pass); `make build && make test-fast` (pass, all packages). `make test-full` was not run: this Task is neither migration nor evidently platform-sensitive by gate definition, though it does touch path handling (CFG-02/CFG-03); the owner may want a fresh full run or Windows CI before close.

**Limitations:** Fingerprint hashes working-tree bytes, so line-ending conversion (autocrlf) can give different fingerprints for identical commits on different platforms; this is explicit, not normalised. Submodule contents are not fingerprinted. A not-a-repository answer also covers Git refusing an unsafe (foreign-owner) directory, since Git reports both with exit 128. Windows behaviour (path separators, symlink and newline-filename skips) was compiled and vetted but not executed here. Mapping configured per-capability scope/exclusions to one `InputScope` is left to the caller; this Task defines `InputScope` only. Mode-only changes count as dirty but leave the fingerprint unchanged. No Check was written and no owner waiver recorded.

## Drift Notes

Return to planning if truthful ancestry requires network access or if scope fingerprinting cannot avoid sensitive data retention.
