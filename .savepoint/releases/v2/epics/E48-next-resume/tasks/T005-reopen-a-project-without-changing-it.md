---
id: E48-next-resume/T005-reopen-a-project-without-changing-it
title: Reopen a project without changing it
status: planned
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

- [ ] `cmd/resume.go` parses `--help` and an optional project directory, rejects unknown flags and extra positional arguments with named errors, and dispatches through an injected runner exactly as `cmd/board.go` does.
- [ ] `cmd/resume.go` contains no record parsing, no gate reading, and no rendering (ARCH-01).
- [ ] `main.go` wires `resume` to load the project, read the V2 router, obtain migration state from `migrate.PendingOperation`, resolve the projection, and render it to stdout.
- [ ] With no directory argument, the project is resolved the same way the existing commands resolve it, and correctness does not depend on the working directory beyond that resolution (ARCH-03).
- [ ] A V2 project renders the projection and exits zero.
- [ ] A V1 project prints a message naming the schema version it read and `savepoint migrate` as the route to V2, writes nothing, and exits nonzero.
- [ ] A missing directory, a directory that is not a Savepoint project, an unreadable `schema_version`, a malformed router, and an index that fails to load each print a named diagnostic to stderr and exit nonzero (FS-06, CFG-01, DATA-03).
- [ ] An unresolvable router selection renders the diagnostic and the available next action and exits zero, because the project is intact and the router line is not.
- [ ] A writer failure mid-render propagates with context and exits nonzero.
- [ ] Every path above leaves the project's file bytes and mtimes unchanged, verified by a snapshot helper comparing the whole project tree before and after invocation.
- [ ] Resume starts no subprocess, no board process, and no network connection; no configured quality gate is executed.
- [ ] A project with a pending migration operation renders the migration rung, and the operation's own files are unchanged by the invocation.
- [ ] `AGENTS.md` records the `resume` dispatch in the `main.go` and `cmd/` Codebase Map rows.
- [ ] `make build && make test` passes.

## Implementation Plan

- [ ] Add `cmd/resume.go` mirroring `cmd/board.go`'s option/runner shape, with its own options type and usage string.
- [ ] Add `cmd/resume_test.go` covering help, the directory argument, unknown flags, and extra arguments.
- [ ] Wire the `resume` case in `main.go`, keeping the migration-state lookup and project resolution in the runner closure.
- [ ] Add the V1-project message path, exiting nonzero without loading a V2 index.
- [ ] Add or extend a tree-snapshot helper in `internal/testutil` capturing relative path, bytes, and mtime for a whole project.
- [ ] Add the no-write assertions to every success and failure case, using the snapshot helper on both sides of the invocation.
- [ ] Add the writer-failure test with a writer that fails after N bytes.
- [ ] Add the pending-migration end-to-end case over a real staged operation directory.
- [ ] Update the `main.go` and `cmd/` Codebase Map rows.
- [ ] Run `go test ./cmd/... ./internal/resume/... ./internal/data/...`, then `make build && make test`.

## Context Log

Pending.
