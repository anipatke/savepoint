---
id: E48-next-resume/T005-reopen-a-project-without-changing-it
title: Reopen a project without changing it
status: done
objective: Ship savepoint resume as a thin read-only command that loads, projects, renders, and writes nothing on any path.
depends_on:
    - E48-next-resume/T004-say-what-the-evidence-actually-shows
complexity_tier: medium
complexity_reason: Thin command plus main wiring, but the no-write guarantee must hold on every failure path.
---

# T005: Reopen a project without changing it

## Problem

Everything above this task is a value transformation. This is where it becomes a command a user runs on a real project at the exact moment they are least sure what state that project is in — which is precisely when a tool must not touch anything.

The no-write guarantee is the product, not a side effect. It covers the success path and every failure path equally: a missing directory, a project that is not a Savepoint project, a malformed router, an index that fails to load, an unresolvable selection, and a writer that fails halfway through rendering all leave the project byte-identical and mtime-identical. No lock file, no cache, no self-healed router line, no "while I was here" repair. A user must be able to run `resume` reflexively, including mid-migration, without wondering whether it changed their state.

This is also where migration state gets injected. `internal/migrate` imports `internal/data`, so the projection cannot call `migrate.PendingOperation` itself; `main.go` resolves it at the same point it resolves the project root and passes it in. That call is a read of the operation record on disk — it stays a read.

A V1 project gets a message, not a projection. `LoadProject` dispatches on `schema_version`, and building a second interpretation over releases and epics — for a lifecycle E50 deletes — would recreate the divergence this epic exists to close. Resume names the schema it found, names `savepoint migrate` as the route, writes nothing, and exits nonzero. The live Savepoint repository is V1 until E50, so its own maintainers will meet this message; that is the accurate state of a mid-release cutover.

The command itself stays thin (ARCH-01): argument parsing, `--help`, an optional project directory, runner injection matching `cmd/board.go`. No domain parsing, no gate reading, no rendering decisions.

## Context Files

- `cmd/board.go`
- `cmd/board_test.go`
- `cmd/doctor.go`
- `cmd/doctor_test.go`
- `main.go`
- `main_test.go`
- `internal/resume/resume.go`
- `internal/data/next.go`
- `internal/data/project.go`
- `internal/data/router_v2.go`
- `internal/migrate/operation.go`
- `internal/testutil/`
- `AGENTS.md`

## Acceptance Criteria

- [x] `cmd/resume.go` parses `--help` and an optional project directory, rejects unknown flags and extra positional arguments with named errors, and dispatches through an injected runner exactly as `cmd/board.go` does.
- [x] `cmd/resume.go` contains no record parsing, no gate reading, and no rendering (ARCH-01).
- [x] `main.go` wires `resume` to load the project, read the V2 router, obtain migration state from `migrate.PendingOperation`, resolve the projection, and render it to stdout.
- [x] With no directory argument, the project is resolved the same way the existing commands resolve it, and correctness does not depend on the working directory beyond that resolution (ARCH-03).
- [x] A V2 project renders the projection and exits zero.
- [x] A V1 project prints a message naming the schema version it read and `savepoint migrate` as the route to V2, writes nothing, and exits nonzero.
- [x] A missing directory, a directory that is not a Savepoint project, an unreadable `schema_version`, a malformed router, and an index that fails to load each print a named diagnostic to stderr and exit nonzero (FS-06, CFG-01, DATA-03).
- [x] An unresolvable router selection renders the diagnostic and the available next action and exits zero, because the project is intact and the router line is not.
- [x] A writer failure mid-render propagates with context and exits nonzero.
- [x] Every path above leaves the project's file bytes and mtimes unchanged, verified by a snapshot helper comparing the whole project tree before and after invocation.
- [x] Resume starts no subprocess, no board process, and no network connection; no configured quality gate is executed.
- [x] A project with a pending migration operation renders the migration rung, and the operation's own files are unchanged by the invocation.
- [x] `AGENTS.md` records the `resume` dispatch in the `main.go` and `cmd/` Codebase Map rows.
- [x] `make build && make test` passes.

## Implementation Plan

- [x] Add `cmd/resume.go` mirroring `cmd/board.go`'s option/runner shape, with its own options type and usage string.
- [x] Add `cmd/resume_test.go` covering help, the directory argument, unknown flags, and extra arguments.
- [x] Wire the `resume` case in `main.go`, keeping the migration-state lookup and project resolution in the runner closure.
- [x] Add the V1-project message path, exiting nonzero without loading a V2 index.
- [x] Add or extend a tree-snapshot helper in `internal/testutil` capturing relative path, bytes, and mtime for a whole project.
- [x] Add the no-write assertions to every success and failure case, using the snapshot helper on both sides of the invocation.
- [x] Add the writer-failure test with a writer that fails after N bytes.
- [x] Add the pending-migration end-to-end case over a real staged operation directory.
- [x] Update the `main.go` and `cmd/` Codebase Map rows.
- [x] Run `go test ./cmd/... ./internal/resume/... ./internal/data/...`, then `make build && make test`.

## Context Log

**Files read:** this task file, `E48-Detail.md`, `AGENTS.md` Codebase Map, completed `T004` task file, `cmd/board.go`, `cmd/board_test.go`, `cmd/doctor.go`, `cmd/migrate.go`, `cmd/migrate_test.go`, `cmd/init.go`, `cmd/upgrade-assets.go`, `main.go`, `main_test.go`, `internal/resume/resume.go`, `internal/data/next.go`, `internal/data/project.go`, `internal/data/router_v2.go`, `internal/data/router.go`, `internal/data/config.go` (`ReadSchemaVersion`/`SchemaVersion`), `internal/data/discover.go` (V2 directory/file-name constants for building fixtures without disk-copying a template), `internal/migrate/command.go` (`ResolveTarget`'s exact-directory, no-upward-walk resolution and its two named diagnostics), `internal/migrate/operation.go` (`PendingOperation`/`PendingOperationReport`, and that `CreateOperation` writes only inside `.savepoint/.migration/<opID>/`), `internal/testutil/fs.go`, `internal/testutil/fixture.go`, `internal/data/project_test.go` and `internal/data/discover_test.go` (V2 fixture file shapes reused as literal content in the new fixtures below).

**Files edited:**
- `cmd/resume.go` (new) — `ResumeOptions{Dir}`, `ResumeRunner`, `RunResume`, `ParseResumeArgs`, mirroring `cmd/migrate.go`'s exact shape (optional single positional directory defaulting to `"."`, `--help`, named errors for an unknown flag or a second directory) rather than `cmd/board.go`'s flag-only shape, since resume needs the directory argument migrate already has and board does not; the runner-injection and `--help`-bypasses-runner contract is identical to both. Contains no import beyond `context`/`fmt`/`io` (ARCH-01), asserted structurally rather than just by convention.
- `cmd/resume_test.go` (new) — help, default directory, explicit directory, unknown-flag rejection, multiple-directory rejection, and runner code/error passthrough, mirroring `cmd/migrate_test.go`'s test shapes. Added `TestPackage_staysThinNoDomainImports`, a `go/build.ImportDir` check over the whole `cmd` package (not just resume.go) proving every file in it — old and new — imports only `context`/`errors`/`fmt`/`io`, since ARCH-01 is a property of the package, not one file.
- `main.go` — added the `resume` dispatch case (mirroring `migrate`'s: `cmd.RunResume` then `os.Exit(code)`, printing `err` to stderr first), `resumeRunner` (thin production wiring calling `runResume(opts.Dir, os.Stdout)`, mirroring `migrateRunner`/`upgradeAssetsRunner`), `runResume(dir string, stdout io.Writer) (int, error)`, and `schemaVersionLabel`. `runResume` resolves `dir` via `migrate.ResolveTarget` (exact-match, no upward walk — the same resolution `migrate` already uses, satisfying ARCH-03 without inventing a second one), then reads `migrate.PendingOperation` **before** the schema-version check: a migration's own recovery guidance already documents that every path can be verified without `schema_version: 2` yet being activated (`RecoveryGuidance` in `operation.go`), so a project mid-conversion can show `schema_version` absent while a pending operation exists — checking pending first is what makes `data.ResolveNext`'s "migration outranks everything" rule true at the command layer too, not just inside the ladder. On a pending operation, `ResolveNext` is called with only `Migration` set (`Index`/`Router` left zero) because `ResolveNext`'s first line returns before touching either. Otherwise: `data.ReadSchemaVersion` decides the V1-message-to-stdout-exit-1-no-error path (mirroring `upgradeAssetsRunner`'s existing V1 migrate-route note, which also reports through stdout rather than as an error) from the five-diagnostics-to-stderr-via-a-returned-error path (`ErrMalformedSchemaVersion`/`ErrUnsupportedSchemaVersion` are real errors from `ReadSchemaVersion`, not V1 fallbacks); `data.LoadProject`, `os.ReadFile` on `router.md`, and `ReadStateV2` each wrap their error with a `resume:`-prefixed context string on failure. No step here writes anything — every one is `os.Stat`/`os.ReadFile`/pure computation — so the no-write guarantee holds structurally, not by discipline.
- `main_resume_test.go` (new) — `writeResumeV2Project` and `resumeRouterV2Content` build a minimal valid V2 project's raw file content directly (no template copy), reusing the exact directory/file-name shapes `internal/data/discover_test.go`'s own fixture helpers use, since those helpers are unexported to `internal/data` and unavailable here. Tests: V2 project renders and exits zero; V1 project prints the schema-version-and-migrate-route message to stdout and exits nonzero; missing directory and not-a-Savepoint-project each name `migrate.ResolveTarget`'s own diagnostic on stderr; a malformed router (`state: bogus`) names the router diagnostic on stderr; an index that fails to load (a Task missing its required title) names the decode diagnostic on stderr; an unresolvable selection (router names a missing Task) renders the diagnostic and the next action together and exits zero; a real pending operation (via `migrate.CreateOperation`, the same helper `main_test.go`'s own `--recover` tests already use) renders the migration rung naming the operation ID and exits zero. Every case above wraps the `runMainForTest` subprocess invocation in the existing `snapshotDir`/`assertSameSnapshot` pair `main_test.go` already defines for the migrate tests, rather than introducing a second copy in `internal/testutil` — the Implementation Plan's suggestion to add one there was superseded once an equivalent helper turned out to already exist exactly where the assertion needs to run (`package main`, alongside `runMainForTest`); reusing it satisfies the acceptance criterion's actual requirement (a snapshot-verified before/after comparison) without duplicating it. The writer-failure case is the one test that calls `runResume` directly rather than through the subprocess harness — a `resumeFailingWriter` that always errors — because the subprocess tests capture real OS stdout, which cannot be made to fail on demand; `runResume`'s injectable `io.Writer` parameter exists specifically so this path is unit-testable at all.
- `AGENTS.md` — extended the `main.go` and `cmd/` Codebase Map rows to name the `resume` dispatch and `runResume`'s responsibilities, and to state ARCH-01 explicitly for `cmd/`.

**Quality gates:** `go build ./...` (clean), `gofmt -l` over every changed/new file (clean), `go vet ./...` (clean), `go test ./cmd/... ./internal/resume/... ./internal/data/... .` (all pass), `go test ./...` and `make build && make test` (pass, all packages). Manually ran `go run . resume --help` and `go run . resume <fresh-init-project>` (the latter correctly resumed into Idea/Design planning) as an end-to-end sanity check beyond the automated suite. No `.savepoint/Health-Check.md` in this project, so the Quick check step is skipped per the build-task skill.

No drift: `cmd/resume.go` and the `main.go` wiring are exactly what `E48-Detail.md`'s Components and files table names for this task (`cmd/resume.go`, and `main.go`'s "`resume` dispatch, and the wiring that supplies migration state"); `main_resume_test.go` is a new test file rather than a new production package, and the `AGENTS.md` rows it explicitly requires are updated.
