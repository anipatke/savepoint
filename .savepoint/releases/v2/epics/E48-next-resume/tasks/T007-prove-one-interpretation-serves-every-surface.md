---
id: E48-next-resume/T007-prove-one-interpretation-serves-every-surface
title: Prove one interpretation serves every surface
status: planned
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

- [ ] A matrix test drives real scaffolded projects end to end — load, router read, projection, render — for every rung of the ladder, asserting exactly one rung kind, one selection, and one next action per project.
- [ ] Every project in the matrix reaches a rung; a project matching zero rungs or more than one fails the test.
- [ ] A project scaffolded by `savepoint init` from `templates/project-v2/`, untouched, resumes into Idea/Design planning with no error and no diagnostic.
- [ ] Every matrix project is snapshotted for bytes and mtimes before and after invocation and is unchanged, on success and on every failure path.
- [ ] A second consecutive `resume` on each matrix project produces byte-identical output and changes nothing.
- [ ] `internal/data` does not import `internal/resume`, `internal/board`, `internal/doctor`, or `internal/migrate`, asserted by an import test.
- [ ] The projection value carries no formatted display strings, no width or color state, and no writer, asserted by review against its type definition and by a test consumer that renders it differently.
- [ ] A second, deliberately minimal consumer in the test renders the same projection and reports the same selection and next action as `internal/resume`.
- [ ] Rendered output for every matrix project is complete and deterministic with no TTY and no color, and readable at 40 columns.
- [ ] Doctor's four categories are asserted against at least one matrix project holding a malformed record, a missing-evidence target, and an advisory open Issue simultaneously, with the project reported unsound for the first two and not for the third.
- [ ] Tests use temporary directories only and touch no developer project (TEST-04).
- [ ] Drift between the epic's Codebase Map additions and the packages that exist is reconciled, or recorded under `## Drift Notes`.
- [ ] `make build && make test` passes.

## Implementation Plan

- [ ] Build the matrix as a table of named project states, each with a fixture builder producing a real project directory.
- [ ] Assert one rung per project, and assert the no-fall-through property by requiring every rung kind to be covered and every project to match exactly one.
- [ ] Reuse the T005 snapshot helper for the before/after comparison across the whole matrix.
- [ ] Add the fresh-scaffold case using the real embedded V2 template tree.
- [ ] Add the second-invocation determinism pass over the same matrix.
- [ ] Add the import-boundary test for `internal/data`.
- [ ] Write the minimal second consumer and assert selection and next-action parity with `internal/resume`.
- [ ] Add the combined doctor-category project and its soundness assertions.
- [ ] Reconcile the Codebase Map against the packages as built; record drift if anything diverged from the design.
- [ ] Run `go test ./... -count=1`, `go vet ./...`, then `make build && make test`.

## Context Log

Pending.
