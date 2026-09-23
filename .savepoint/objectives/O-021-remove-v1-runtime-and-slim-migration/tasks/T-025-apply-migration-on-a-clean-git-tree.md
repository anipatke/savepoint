---
id: T-025
title: Apply migration on a clean git tree
objective: O-021
planned_by: {role: planner, session: o021-design-20260923}
status: in_progress
stage: audit
complexity_tier: high
complexity_reason: "Replaces the migration write path users rely on, removes ~2k lines of recovery code and its tests, and must keep conversion output byte-identical on both golden fixtures across platforms."
depends_on: [{task: T-024, requires: clear}]
owner_validation:
    required: true
    accepted_check: ""
---

# T-025: Apply migration on a clean git tree

## Outcome

`savepoint migrate --apply` checks that every path it will write or remove
is in a git work tree with no uncommitted or untracked changes, then writes
the converted project directly. If anything fails, it stops, lists what it
wrote, and tells the user how to undo with git. The write-ahead journal,
crash recovery, cutover preflight, and platform replace code are deleted.
Preview output and converted bytes are unchanged.

## User Check

In a temp git repo holding a committed copy of the `v1-history` fixture:
`savepoint migrate` previews without writing. `savepoint migrate --apply`
converts, and `git status` shows exactly the planned changes. Undo with the
printed git command and confirm the tree is back to the commit. Repeat with
an uncommitted edit under `.savepoint/` and outside git; both are refused
before any write, with a message saying what to do.

## Done When

- `--apply` refuses before any write when: `git` is not on PATH; the
  project is not inside a git work tree; or any planned write/remove path
  has uncommitted or untracked changes. Each refusal is a named error with
  the next step (commit or stash; `git init`; install git). Exit codes
  follow the command's existing refusal convention.
- The clean-tree check is injected (`CommandOptions`) so unit tests do not
  need git; one integration test uses a real temp repo and skips with a
  named reason when `git` is absent.
- Apply writes each planned file (create for new paths, replace for
  existing ones, via plain temp-file-and-rename), performs planned removals,
  writes the `v1-to-v2.yml` manifest, and sets `schema_version: 2` last.
  On error it stops and prints the written paths and the git undo command.
- A second `--apply` on an already-migrated project still reports
  "already migrated" and writes nothing.
- Deleted: `operation.go`, the journal and recovery paths in `apply.go`,
  `PendingOperation`, `revalidateSources` (plan and apply already run in one
  process), `cutover.go`, `replace.go`, `replace_unix.go`,
  `replace_windows.go`, their tests, and `data.ResolveReleaseCutover` plus
  `internal/data/release_cutover.go` if nothing else consumes them.
  Recovery-only fields leave the manifest schema; mapping, source hashes,
  archives, and decisions stay.
- Both golden fixtures produce byte-identical output. Plan, convert,
  decisions, inventory, classify, and end-to-end tests pass unchanged apart
  from removed recovery cases.
- README "Safe migration" describes the clean-git rule and undo.
- `git diff --check` and a fresh `make test-full` (migration and
  platform-sensitive, includes the Windows build) pass.

## Context Files

`internal/migrate/command.go`, `internal/migrate/command_test.go`,
`internal/migrate/apply.go`, `internal/migrate/apply_test.go`,
`internal/migrate/operation.go`, `internal/migrate/operation_test.go`,
`internal/migrate/cutover.go`, `internal/migrate/cutover_test.go`,
`internal/migrate/replace.go`, `internal/migrate/replace_unix.go`,
`internal/migrate/replace_windows.go`, `internal/migrate/manifest.go`,
`internal/migrate/manifest_test.go`, `internal/migrate/end_to_end_test.go`,
`internal/migrate/live_boundary_test.go`, `internal/migrate/plan.go`,
`internal/migrate/testdata/golden/v1-basic.yml`,
`internal/migrate/testdata/golden/v1-history.yml`,
`internal/data/release_cutover.go`, `cmd/migrate.go`,
`cmd/migrate_test.go`, `main.go`, `README.md`.

## Design References

Design sections 2, 10 (commands), and the migration section describing
`.savepoint/migrations/v1-to-v2.yml`.

## Guardrails

FS-01, FS-03, FS-05, FS-06, DATA-01, ARCH-01, TEST-02, TEST-03, TEST-04,
TEST-08.

## Implementation Plan

1. Add the injected clean-tree check (`git rev-parse --is-inside-work-tree`
   and `git status --porcelain --untracked-files=all -- <paths>` from the
   project root) and its refusal errors, with unit tests.
2. Rewrite `Apply` as check → write files → removals → manifest → schema,
   with error reporting that lists written paths.
3. Delete the journal, recovery, cutover, and replace code and tests;
   trim the manifest's recovery fields.
4. Confirm golden output is byte-identical and the already-migrated rerun
   is a no-op.
5. Update README "Safe migration"; run a fresh `make test-full`.

## Boundaries

No change to what the conversion produces, the preview format, the
decisions file, or archive layout. No backup-copy fallback for non-git
projects (Confirmed Decision 2). No change to V2 runtime loading.

## Technical Verification

Focused `internal/migrate` tests during iteration; fresh `make test-full`
at handoff. Owner validation per the User Check above.

## Technical Evidence

**Extra reads beyond the Context Files:**

- `.savepoint/objectives/O-021-remove-v1-runtime-and-slim-migration/tasks/T-024-check-the-schema-without-migrate.md` — read to verify T-025's `requires: clear` dependency. It contains T-024's explicit owner waiver (recorded at `2026-09-23T11:11:35Z`); the owner reaffirmed that waiver in this conversation at `2026-09-23T11:15:05Z` and directed its evidence to the mandatory Full Objective Check. This waiver satisfies the downstream dependency and does not claim technical `CLEAR`.
- Targeted `rg` search results in `internal/data/release_gate_v2_test.go` — needed to identify the remaining tests of `ResolveReleaseCutover` before removing the resolver and its implementation file.
- Repository-wide Go reference search for `PendingOperation`, `PreflightCutover`, `ResolveReleaseCutover`, `ReplaceFile`, and recovery-only manifest fields — needed to find every compile-time consumer before deleting those APIs and files.
- The same search identified out-of-context consumers to reconcile before deletion: `main_test.go` recovery dispatch cases, `main_board_next_parity_test.go` and `internal/data/release_gate_v2_test.go` Release-cutover tests, `internal/data/next_test.go`'s stale replace-primitive comment, `internal/board/v2/watch.go`'s stale pending-operation comment, and the three platform-specific replace test files. These references must be checked so removed APIs do not leave compile failures or stale claims.
- `main_test.go`, `main_board_next_parity_test.go`, `internal/data/release_gate_v2_test.go`, `internal/board/v2/watch.go`, and `internal/data/next_test.go` — inspected only the matching tests/comments to remove obsolete recovery/cutover claims and the import-graph test comment after deleting those APIs.
- `internal/board/v2/watch_test.go` — targeted search found `.migration` in the excluded-path cases; inspect those exact cases before dropping the obsolete watcher directory.
- `internal/board/v2/watch_test.go` — read and updated the relevant watch-set cases so `.migration` is explicitly excluded now that it no longer carries live operation state.
- `internal/init/upgrade_test.go` — inspect the matching recovery-guidance assertion found by the symbol scan to determine whether it describes a live caller or retired journal behavior.
- `internal/init/upgrade_test.go` — read the matching test body; it creates a journal only to verify upgrade-assets ignores leftovers, so replace that deleted-API setup with an inert historical directory fixture and remove the obsolete `--recover` claim.
- `internal/migrate/inventory.go` and `internal/migrate/inventory_test.go` — inspect the remaining `.migration` exclusion so its comment and test describe inert journals left by earlier versions rather than active operation state.
- `main_board_test.go` and `main_resume_test.go` — the affected-package compile surfaced two remaining fixtures calling deleted `migrate.CreateOperation`; inspect those tests and replace them with inert legacy journal files so they continue proving runtime callers ignore leftovers.
- `internal/data`'s `ReadSchemaVersion` implementation — targeted symbol lookup and definition read at `2026-09-23T12:01:44Z`; confirmed that a missing `config.yml` is a supported V1 shape and the schema reader is read-only. This explained the repository-copy E2E case and led to creating `config.yml` during final schema activation.

**Execution evidence:**

- Clean-tree refusal coverage: `go test ./internal/migrate ./cmd -count=1` passed. Unit cases cover missing Git, a non-worktree, modified and untracked planned paths, and ignored planned paths; refusals happen before the writeability probe. The real-Git integration test previews read-only, applies from a committed fixture, checks every planned path, runs the matching path-scoped restore/clean operations, confirms a clean tree, then proves a modified `.savepoint/config.yml` is refused without project-file changes.
- Apply ordering and retry coverage: the migration tests pass for direct creates/replacements, archive copies and removals, manifest before schema activation, written-path error reporting, missing-config creation as the final write, and already-migrated no-op behavior.
- Conversion preservation: both golden fixtures and repository-copy migration passed as part of the full suite; recovery-only manifest fields and journal/cutover/replace code are removed while source hashes, identities, archives, and decisions remain.
- `git diff --check` passed after implementation.
- Latest `make test-full` passed on Go `go1.26.2 linux/amd64`, including the full Go suite and Linux, Darwin, and Windows builds; completion evidence recorded at `2026-09-23T12:15:58Z`. Two earlier full-gate attempts reported a root-package setup failure. `go test . -count=1` and direct `go test -json -count=1 ./...` both passed, and the complete `make test-full` rerun passed.
- The independent Task Check is explicitly waived by the owner; its evidence routes to the mandatory Full Objective Check. Owner acceptance and marking T-025 `done` remain pending.

## Drift Notes

Design and AGENTS.md still describe recoverable apply until T-026.
