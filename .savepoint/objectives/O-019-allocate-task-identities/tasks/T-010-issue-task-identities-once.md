---
id: T-010
title: Issue each Task number once
objective: O-019
planned_by: {role: planner, session: i024-allocation-design-20260923}
status: in_progress
stage: build
complexity_tier: high
complexity_reason: Durable allocation and concurrent writers must behave consistently across supported platforms.
depends_on: []
owner_validation: {required: false}
---

# T-010: Issue each Task number once

## Outcome

A project-owned allocation operation issues monotonically increasing Task IDs under a project lock, including when Tasks are deleted or creation fails.

## Done When

- Strict-load before allocation; refuse invalid projects without changing records.
- Initialize a durable high-water mark from all active Tasks; never lower it after deletion or failure.
- Concurrent callers cannot receive the same number. A held lock is retried for a bounded time, then refused naming `.savepoint/task-ids.lock`; a leftover lock from an interrupted process is refused the same way and never removed automatically. The same tests pass on the Linux and Windows CI jobs.
- Existing duplicate-ID validation still fails closed and names both paths.

## Context Files

`internal/data/project.go`, `internal/data/discover.go`, `internal/data/task_v2.go`, `internal/data/write.go`, `internal/data/project_test.go`, `internal/data/discover_test.go`, `internal/data/task_v2_test.go`, `internal/data/write_test.go`, `internal/data/write_linux_test.go`, `.github/workflows/ci.yml`, `AGENTS.md`.

## Design References

Design sections 1, 2, 5, 6, and 9; O-019 Architectural Considerations and Confirmed Design Decisions.

## Guardrails

FS-01, FS-04, FS-05, FS-06, DATA-03, CFG-02, ARCH-03, ARCH-04, TEST-01, TEST-02, TEST-03, TEST-04, TEST-05, TEST-08.

## Implementation Plan

1. Confirm O-018's grammar and the temp-write, sync, and rename pattern in `internal/data/write.go` (`replaceV2File`); return REPLAN REQUIRED if interfaces contradict this plan.
2. Add the create-exclusive `.savepoint/task-ids.lock` and the `.savepoint/task-ids.yml` high-water mark per O-019's Confirmed Design Decisions; initialize existing projects from their complete active Task set.
3. Persist the next reservation while holding the lock. Never roll it back after failure.
4. Test bootstrap, deletion, failure, concurrency, a held lock, a leftover lock, invalid projects, and the unchanged two-path duplicate diagnostic. Keep the tests platform-neutral so the Windows CI job runs them.

## Boundaries

No allocator for other record kinds, Task lifecycle change, or weakening of strict loading.

## Technical Verification

Focused `internal/data` tests, `git diff --check`, and `make test-full` (platform-sensitive), with the Windows CI job as Windows evidence.

## Technical Evidence

REPLAN REQUIRED — the Context Files name `internal/migrate/replace.go`,
`internal/migrate/replace_unix.go`, and `internal/migrate/replace_windows.go`,
but none of those paths exists. Implementation Plan step 1 depends on confirming
the existing platform-safe write behavior, so that assumption cannot be checked
within this plan's declared context. No implementation code was changed. Task
remains `in_progress` at `build` pending planner review.

Files read: `agent-skills/savepoint-task/SKILL.md`, `.savepoint/router.md`, this
Task, `.savepoint/objectives/O-019-allocate-task-identities/Objective.md`,
user-provided `AGENTS.md` instructions, targeted section/rule searches in
`.savepoint/Design.md` and `.savepoint/Guardrails.md`, and targeted declaration
searches in `internal/data/project.go`, `internal/data/discover.go`,
`internal/data/task_v2.go`, and `internal/data/write.go`. The named
`internal/migrate/replace*.go` paths were checked by the targeted search and
reported missing. Remaining Context Files were not read because the missing
files make the planned platform review impossible.

Extra reads: `.savepoint/Design.md` section headings, because the Task cites
Design sections 1, 2, 5, and 6; `.savepoint/Guardrails.md` rules FS-01, FS-04,
FS-05, FS-06, DATA-03, CFG-02, ARCH-03, ARCH-04, TEST-01 through TEST-05, and
TEST-08, because the Task names them. These reads did not resolve the missing
platform-write context.

Files changed: this Task record (start lifecycle and replan evidence), the
owning `Objective.md` (start lifecycle), and `.savepoint/router.md` (replan
handoff to design). No acceptance criteria changed. Limitations: no acceptance
criterion was implemented or verified; no build or test gate was run.

Commands: initial `savepoint resume` exited 0 and selected Start T-010; after
the replan handoff, `savepoint resume` exited 0 but selected Build T-010 while
`.savepoint/router.md` records `state: design`. This routing discrepancy is
included for planner review.

## Drift Notes

If cross-platform locking or crash recovery requires a different storage contract, return REPLAN REQUIRED.

Replan 2026-09-24 (planner, owner-chat-20260924): the executor's REPLAN REQUIRED above was resolved by replacing the deleted `internal/migrate/replace*.go` Context Files with `internal/data/write.go`'s existing safe write and by recording the owner-confirmed lock design in O-019. The executor's routing observation was not a defect: `savepoint resume` shows `Replan` only when the Task carries a `replan:` frontmatter block, and none was recorded.
