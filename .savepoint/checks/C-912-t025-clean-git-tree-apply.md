---
id: C-912
scope: {kind: task, id: T-025}
result: NEEDS WORK
checked_by: {role: checker, session: t025-task-check-20260924}
executed_session: unrecorded
checked_at: '2026-09-23T21:00:00Z'
reviewed:
  base_commit: 89ed584
  head_commit: 4f3b2df02eb76c3f0e1669d09faa17655c19cfc6
  files:
    - internal/migrate/command.go
    - internal/migrate/command_test.go
    - internal/migrate/apply.go
    - internal/migrate/apply_test.go
    - internal/migrate/manifest.go
    - internal/migrate/manifest_test.go
    - internal/migrate/end_to_end_test.go
    - internal/migrate/plan.go
    - internal/migrate/testdata/golden/v1-basic.yml
    - internal/migrate/testdata/golden/v1-history.yml
    - cmd/migrate.go
    - cmd/migrate_test.go
    - main.go
    - README.md
  dependencies: []
issues: [I-039]
supersedes: null
---

# C-912: T-025 Task Check — apply migration on a clean git tree

Quick evidence, requested by the owner ("Check T025"). This session did not
build T-025. The Task's earlier owner waiver of the optional Task Check is
replaced by this explicit request.

## Scope Lock

1. Criteria: T-025's eight Done When items, the User Check, and the Boundaries
   (no change to conversion output, preview, decisions, or archive layout).
   Guardrails as named by the Task. Gate: `git diff --check` plus a fresh
   `make test-full`.
2. Changed files: listed under `reviewed.files`. Public entry point:
   `savepoint migrate [--apply] [--dry-run] [--decisions F] [dir]` through
   `cmd.RunMigrate` → `migrate.RunCommand` → `migrate.Apply`.
3. External boundary relied on: the `git` subprocess (`rev-parse`,
   `status --porcelain --untracked-files=all --ignored=matching`), and the
   printed `git restore` / `git clean` undo string.
4. Matrix below. Unsupported and not applicable: text-width classes (no
   rendering is in scope), no-colour or TTY sinks (plain output only), and
   network or timeout (git runs locally and synchronously).
5. Materiality: an Issue must break a Done When item or the Outcome through
   `savepoint migrate` on a V1 project that `ResolveTarget` accepts.

## Coverage Matrix

Independent probe: a CLI built from HEAD and run against scratch git repos
copied from `internal/data/testdata/migration/v1-history/project`, plus a
minimal V1 project. Tree hashes and `git status` are the oracle.

| # | Surface / state | Input | Expected | Actual | Result |
|---|---|---|---|---|---|
| 1 | preview, clean repo | v1-history | exit 0, no writes | exit 0, 0 status lines | pass |
| 2 | apply, clean repo | v1-history | converts; only planned paths change | exit 0; 38 changed paths, which is exactly the planned set (also asserted in the integration test) | pass |
| 3 | printed undo after success | v1-history | tree back to the commit | exit 0; `git status --ignored` empty | pass (see observation 1) |
| 4 | retry apply after undo | v1-history plus leftover empty dirs | converts again | exit 0 | pass |
| 5 | second apply | migrated project | "already migrated", no writes | as expected; exit 0; no new changes | pass |
| 6 | dirty planned path | `router.md` edited | refused before any write, next step named | exit 1, `ErrDirtyGitPaths` naming `M .savepoint/router.md`; tree hash unchanged | pass |
| 7 | dirty non-planned path | untracked `notes.txt` | allowed | exit 0 | pass |
| 8 | outside git | no `.git` | refused, `git init` hint | exit 1, `ErrNotGitWorkTree`; files unchanged | pass |
| 9 | git missing | `PATH=/nonexistent` | refused, install hint | exit 1, `ErrGitUnavailable`; no changes | pass |
| 10 | failure mid-apply | `releases/v1` chmod 555 (removal fails) | stops, lists written paths, prints undo, schema not activated | exit 1; 23 written paths listed; no `schema_version` line; printed undo leaves the tree clean; no temp files left | pass |
| 11 | printed undo, no tracked planned paths | `.savepoint/Design.md` only, committed | undo returns the tree to the commit | `fatal: you must specify path(s) to restore` (exit 128); `config.yml` and the manifest are left behind | **Issue I-039** |
| 12 | ignored planned path | unit test | refused | covered by `TestRunCommand_applyGitRefusalsHappenBeforeAnyWrite` | pass |
| 13 | refusal order | unit tests | git refusals come before the writeability probe | as tested | pass |
| 14 | golden fixtures | v1-basic, v1-history | byte-identical apart from the recovery fields that were removed | golden diff drops only `generated_at` and `operation_id` from the manifest | pass |

### Workflow / side-effect lock (Apply)

| Order | Operation | Failure effect | Oracle |
|---|---|---|---|
| 0 | target resolve, plan, conflicts, appliable | no writes; exit 1 or 2 | rows 1, 5 |
| 1 | git work tree and path-status check | no writes; exit 1 | rows 6, 8, 9 |
| 2 | writeability probe (temp file, then removed) | no project writes | existing unit test |
| 3 | render the whole batch before any write | nothing written, undo printed | code read (`buildApplyBatch`) |
| 4 | creates, archive copies, router replace (temp file then rename) | stop; written paths and undo printed | row 10 |
| 5 | archived source removals | stop; same report | row 10 |
| 6 | manifest | stop; same report | unit test |
| 7 | `schema_version: 2` last (created when `config.yml` is missing) | stop; project stays V1 | rows 2, 10 |

## Acceptance Classification

| Done When | Result | Evidence |
|---|---|---|
| Refusals (no git, not a work tree, dirty or untracked paths), named errors and next steps, exit code 1 | Proven | rows 6, 8, 9, 12 |
| Injected clean-tree check; real-git test that skips with a named reason | Proven | `CommandOptions.RunGit`; `t.Skip("git is not available; ...")` |
| Apply order, temp file and rename, manifest, schema last, error lists written paths **and the git undo command** | **Issue** | row 10 passes; row 11 fails (I-039) |
| Second apply is a no-op | Proven | row 5 |
| Deletions | Proven | `operation*`, `cutover*`, `replace*`, and `release_cutover.go` are gone; no Go references remain to `PendingOperation`, `revalidateSources`, `ResolveReleaseCutover`, `PreflightCutover`, or `ReplaceFile` |
| Golden byte-identity | Proven | row 14 |
| README "Safe migration" | Proven | README.md:201-223 |
| `git diff --check` and fresh `make test-full` | Proven | `git diff --check 89ed584 4f3b2df` clean; fresh `make test-full` at about 2026-09-23T20:54Z, go1.26.2 linux/amd64: all packages ok, Linux, Darwin, and Windows builds ok, EXIT 0 |

Focused tests: `go test ./internal/migrate ./cmd -count=1` passed.

## Materiality

| Issue | Likelihood | Impact | Materiality | Recommendation |
|---|---|---|---|---|
| I-039 undo string fails with an empty tracked group | Low: needs a V1 project with no `config.yml`, no router, and nothing archived | Low: nothing is destroyed; two untracked files remain and git states the cause | Low | Fix now; about a one-line change plus a test that runs the printed string. The owner may accept it instead. |

## Observations (non-blocking)

1. Undo leaves empty directories behind (`objectives/…`, `archive/v1/…`,
   `migrations/`). Git ignores them and retry works (row 4).
2. If apply is killed mid-write (a signal, not an error return), a
   `.savepoint-migrate-apply-*` temp file can survive. The printed undo does
   not name it. This is outside the Done When items, which cover the error
   path only.
3. AGENTS.md still names `migrate.PendingOperation`,
   `internal/data/release_cutover.go`, and `ResolveReleaseCutover`. This is
   declared drift and belongs to T-026.
4. README says the undo command is "scoped to those paths" (the changed
   paths). It actually covers every planned path. That is harmless, because
   the clean-tree check proved every planned path was clean before apply.
5. At check time the working tree held uncommitted edits to T-025 evidence,
   I-031, and I-038. These are metadata only; no code or test inputs changed.

## Owner Validation Still Needed

T-025 has `owner_validation.required: true`. The owner can accept only a
current `CLEAR` Check. After I-039 is repaired, a re-check that supersedes
C-912 is needed. The Full Objective Check for O-021 remains mandatory.
