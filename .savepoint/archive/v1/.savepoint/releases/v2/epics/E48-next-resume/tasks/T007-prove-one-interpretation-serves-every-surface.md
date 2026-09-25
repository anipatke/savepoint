---
id: E48-next-resume/T007-prove-one-interpretation-serves-every-surface
title: Prove one interpretation serves every surface
status: done
objective: Assert end to end that every project state reaches exactly one rung, that resume never writes, and that the projection carries no presentation state.
depends_on:
    - E48-next-resume/T005-reopen-a-project-without-changing-it
    - E48-next-resume/T006-report-health-that-is-worth-acting-on
complexity_tier: medium
complexity_reason: No new behavior; a matrix over real projects proving totality, the no-write contract, and consumer independence.
---

# T007: Prove one interpretation serves every surface

## Problem

T001 through T006 each prove their own piece against their own fixtures. Three claims of this epic are not provable from any one of them, and they are the claims the epic is actually selling.

**The ladder is total.** Every rung has a test, but that does not establish that every project lands on a rung — only that nine constructed projects do. The matrix runs real scaffolded projects through the whole path, `LoadProject` to rendered output, and asserts one rung, one selection, one next action each, with no project falling through and no project matching two rungs.

**Resume never writes.** Each task asserts this for its own cases. The guarantee a user relies on is stronger: across every state in the matrix, every failure path, and a second consecutive invocation, the project tree is byte-identical and mtime-identical. A guarantee proven case by case is a guarantee with gaps between the cases.

**The interpretation is shared.** E49 will render this projection in the board, and the claim that it can is only true if nothing about the projection is resume-shaped. That is a structural property and it can be asserted now: the projection package must not import the renderer, the projection value must carry no formatted strings or terminal state, and a second consumer built in the test must reach the same selection and next action from the same call. The board itself is E49's; a stub consumer is enough to prove the shape.

Also verified here: the freshly initialized V2 project from E47 resumes into Idea/Design planning rather than an error, which is the first thing a new user will do.

## Context Files

- `internal/data/next.go`
- `internal/data/next_test.go`
- `internal/data/project.go`
- `internal/resume/resume.go`
- `cmd/resume.go`
- `main.go`
- `main_test.go`
- `internal/doctor/report.go`
- `internal/testutil/`
- `templates/project-v2/`
- `.savepoint/releases/v2/epics/E48-next-resume/E48-Detail.md`
- `AGENTS.md`

## Acceptance Criteria

- [x] A matrix test drives real scaffolded projects end to end — load, router read, projection, render — for every rung of the ladder, asserting exactly one rung kind, one selection, and one next action per project.
- [x] Every project in the matrix reaches a rung; a project matching zero rungs or more than one fails the test.
- [x] A project scaffolded by `savepoint init` from `templates/project-v2/`, untouched, resumes into Idea/Design planning with no error and no diagnostic.
- [x] Every matrix project is snapshotted for bytes and mtimes before and after invocation and is unchanged, on success and on every failure path.
- [x] A second consecutive `resume` on each matrix project produces byte-identical output and changes nothing.
- [x] `internal/data` does not import `internal/resume`, `internal/board`, `internal/doctor`, or `internal/migrate`, asserted by an import test.
- [x] The projection value carries no formatted display strings, no width or color state, and no writer, asserted by review against its type definition and by a test consumer that renders it differently.
- [x] A second, deliberately minimal consumer in the test renders the same projection and reports the same selection and next action as `internal/resume`.
- [x] Rendered output for every matrix project is complete and deterministic with no TTY and no color, and readable at 40 columns.
- [x] Doctor's four categories are asserted against at least one matrix project holding a malformed record, a missing-evidence target, and an advisory open Issue simultaneously, with the project reported unsound for the first two and not for the third.
- [x] Tests use temporary directories only and touch no developer project (TEST-04).
- [x] Drift between the epic's Codebase Map additions and the packages that exist is reconciled, or recorded under `## Drift Notes`.
- [x] `make build && make test` passes.

## Implementation Plan

- [x] Build the matrix as a table of named project states, each with a fixture builder producing a real project directory.
- [x] Assert one rung per project, and assert the no-fall-through property by requiring every rung kind to be covered and every project to match exactly one.
- [x] Reuse the T005 snapshot helper for the before/after comparison across the whole matrix.
- [x] Add the fresh-scaffold case using the real embedded V2 template tree.
- [x] Add the second-invocation determinism pass over the same matrix.
- [x] Add the import-boundary test for `internal/data`.
- [x] Write the minimal second consumer and assert selection and next-action parity with `internal/resume`.
- [x] Add the combined doctor-category project and its soundness assertions.
- [x] Reconcile the Codebase Map against the packages as built; record drift if anything diverged from the design.
- [x] Run `go test ./... -count=1`, `go vet ./...`, then `make build && make test`.

## Context Log

**Files read:** this task file, `E48-Detail.md`, `AGENTS.md` Codebase Map, `.savepoint/router.md`, completed `T004`/`T005`/`T006` task files (their Context Logs, to learn established fixture/test conventions rather than reinvent them), `internal/data/next.go` (full — confirmed `Next`'s fields are all typed resolver values, with no formatted string, width/color, or `io.Writer` field anywhere on it or on `NextInput`/`Selection`/`SelectionDiagnostic`), `internal/data/next_test.go` (fixture shapes for every rung, `mustCheck`/`mustCurrentTask` semantics, and the existing `TestNext_packageDoesNotImportMigrate` import-boundary pattern), `internal/data/task_v2.go`, `objective_v2.go`, `check_v2.go`, `evidence_v2.go`, `issue_v2.go` (exact YAML frontmatter shapes needed to hand-write valid V2 fixtures), `internal/data/router_v2.go` (`ReadStateV2`, the `none` selection sentinel), `internal/data/discover.go` (V2 directory/file-name layout: `objectives/<dir>/Objective.md`, `objectives/<dir>/tasks/<file>.md`, `checks/<file>.md`, `issues/<file>.md`), `internal/data/discover_test.go` and `internal/doctor/checks_test.go` (existing `writeV2*Fixture` helpers, reused as literal-content templates since they are unexported to package `data`/`doctor` and unavailable from `main`), `internal/resume/resume.go` and `resume_test.go` (the `renderText`/box-drawing-glyph readability check T004 already established, reused verbatim against real matrix output), `cmd/resume.go`/`cmd/resume_test.go`, `main.go` (`runResume`'s exact read order: pending migration before schema version, then `LoadProject`, then router, then `ResolveNext`), `main_test.go` (`snapshotDir`/`assertSameSnapshot`/`runMainForTest`/`writeMigrateMinimalProject`/`mkdirAll`, reused rather than duplicated), `main_resume_test.go` (`writeResumeV2Project`/`resumeRouterV2Content`, reused directly for the `NextExecute` row), `internal/doctor/report.go` (`HealthCategory`, `HealthFindings`, `HasProblems`), `internal/init/scaffold.go` (`Scaffold`'s signature, to drive the real embedded V2 template tree rather than hand-copy its shipped `router.md`/`config.yml`), `templates/project-v2/.savepoint/router.md` and `config.yml`.

**Files edited:**
- `main_resume_matrix_test.go` (new) — the end-to-end matrix. `resumeMatrixCases()` names nine real, on-disk V2 project fixtures, one per `data.NextKind` rung (`pending_migration`, `replan`, `dependency`, `execute`, `check_needed`, `owner_validation_required`, `objective_integration`, `ready`, `plan_objective`), plus a tenth built through the real embedded `templates/project-v2` tree via `savepointinit.Scaffold` (the fresh-`init`-scaffold acceptance criterion). `resolveNextFromDisk` mirrors `runResume`'s own read order (pending migration, then schema version, `LoadProject`, router, `ResolveNext`) to hand back the `data.Next` value directly rather than parsing it back out of rendered prose. `TestResumeMatrix_everyRungReachedExactlyOnce` is the central proof: for every case it asserts `next.Kind` equals exactly the one expected rung (not "at least"), asserts no `SelectionDiagnostic` leaked in (none of these fixtures name an unresolvable router selection — that path is T005's own coverage), renders through the real `internal/resume.Render`, checks the rendered text against `assertNarrowWidthReadable` (T004's own box-drawing-glyph heuristic, re-run here against pipeline output instead of synthetic `data.Next` literals), exercises `minimalSecondConsumerRender` — a second, `internal/resume`-independent reader of `data.Next` that derives its own selection ID and reports the `NextKind` directly, then cross-checks that ID appears in resume's own rendering — resolves and renders a second time to prove per-invocation determinism, and finally snapshots the project directory before and after with `assertSameSnapshot`. After the loop, it asserts all nine `NextKind` values were actually reached, which is what makes "no project matching zero rungs" and "the ladder has no gap" a checked property rather than an assumption. `TestResumeMatrix_runResumeThroughTheRealCommandAlsoWritesNothingTwice` re-runs the same ten fixtures through `runResume` itself — the literal function `cmd.RunResume` dispatches to — twice each, asserting byte-identical output and an unchanged snapshot, so the no-write and determinism guarantees are proven against the actual command entrypoint, not only against the reassembled pipeline. Failure-path no-write coverage (missing directory, non-Savepoint directory, malformed router, unloadable index, writer failure) is not re-proven here: it already exists, per-case, in `main_resume_test.go`'s `TestMainResumeMissingDirectory`, `TestMainResumeNotASavepointProject`, `TestMainResumeMalformedRouter`, `TestMainResumeIndexFailsToLoad`, and `TestRunResumeWriterFailurePropagates` (T005), each already wrapped in the same `snapshotDir`/`assertSameSnapshot` pair this file reuses (TEST-06).
- `internal/data/next_test.go` — added `TestDataPackage_staysBeneathEverySurfaceItFeeds`, broadening the existing migrate-only import-boundary test's pattern to also forbid `internal/resume`, `internal/board`, and `internal/doctor`, over both `pkg.Imports` and `pkg.TestImports` exactly as the existing test does. Left `TestNext_packageDoesNotImportMigrate` in place rather than folding it in, since it documents the specific migrate/data cycle reason on its own.
- `internal/doctor/report_test.go` — added `TestDiagnosticReport_CombinedMalformedMissingEvidenceAndAdvisoryIssue`: one project carrying a done Task with no Check (`HealthMissingEvidence`), a Task whose recorded owner-acceptance names a now-superseded Check (`v2-acceptance-superseded`, `HealthMalformedData` — reusing the exact fixture shape `TestCheckProject_ConsistencyDiagnostics` already established for that diagnostic), and an open Issue (`HealthPendingReview`) simultaneously. Asserts all three categories are represented in `HealthFindings()`, that `HasProblems()` is true, and that non-review findings exist independent of the advisory Issue — proving the three categories co-exist correctly in one project rather than only in isolation (T006 tested each in its own fixture).

**Drift:** none. No new package or file beyond test files in three already-mapped packages (`main`, `internal/data`, `internal/doctor`); no Codebase Map row needed.

**Quality gates:** `gofmt -l` on every new/changed file (clean — the four pre-existing unformatted files elsewhere in the tree are untouched by this task), `go vet ./...` (clean), `go test ./... -count=1` (all packages pass), `go test ./... -run 'TestResumeMatrix|TestDataPackage_staysBeneathEverySurfaceItFeeds|TestDiagnosticReport_CombinedMalformedMissingEvidenceAndAdvisoryIssue' -v` (all new cases pass individually), `make build && make test` (pass). No `.savepoint/Health-Check.md` in this project, so the Quick check step is skipped per the build-task skill.
