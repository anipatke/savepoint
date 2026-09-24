---
id: O-019
title: Allocate Task identities without planner selection
status: in_progress
release: R-006
depends_on: [O-018]
---

# O-019: Allocate Task identities without planner selection

## Outcome

Every new V2 Task receives a unique, project-wide, never-reused `T-###` identity from one creation operation. Planners provide an Objective and Task draft without choosing a number.

## Why

I-024 showed that a planner inferred `T-006` from one Objective while another Objective already owned it. Strict loading caught the collision only after the invalid Task reached the live project. The owner chose deterministic allocation and retirement of issued numbers, including deleted Tasks.

## Success Conditions

- One durable project-owned high-water mark issues the next Task number under a lock. Bootstrap from all active Tasks never lowers that mark; failure or deletion never releases an issued number.
- A supported `savepoint create-task` operation takes an ID-free draft, creates matching frontmatter and path, and strict-loads the entire V2 index before reporting success. A failed creation leaves no new Task file and never overwrites existing content.
- Concurrent planners cannot receive the same ID. Failure, interruption, recovery, Linux, and Windows behavior have named evidence.
- The loader keeps its fail-closed duplicate diagnostic naming both paths.
- Canonical and scaffold planning guidance require this operation. `AGENTS.md` grants an exception only for that command, delivered to existing projects through upgrade-assets.
- Cross-Objective workflow scenarios prove planners do not choose or reuse an ID; strict V2 loading follows creation or rename of other identity-bearing records.
- `git diff --check`, `make build && make test`, and the mandatory Full Objective Check verify the integrated result before owner closure.

## Architectural Considerations

- O-018 owns the hyphenated identity grammar and project migration. This Objective follows it and introduces no second ID parser.
- `internal/data` owns strict global V2 indexing and Task validation. A focused new `internal/` package may own allocation and creation if needed. `cmd/` only parses and dispatches.
- Durable high-water state is Savepoint-managed, initialized from the full active Task set for existing projects and safely updated across processes and platforms. A failed reservation may leave a gap, never a reused number.
- Creation validates the project before allocation and after writing. The strict loader remains the structural gate.
- The current agent rule bans all `savepoint` commands. A narrow `create-task` exception and its managed-guide upgrade path are part of this Objective.

## Confirmed Design Decisions

On 2026-09-24, after T-010 returned REPLAN REQUIRED because its planned
platform-write files (`internal/migrate/replace*.go`) had been deleted by
O-021, the owner confirmed a simple, portable allocation design:

- The project lock is a create-exclusive lock file,
  `.savepoint/task-ids.lock`, opened with `O_CREATE|O_EXCL`. The same code
  runs on Linux and Windows; there is no platform-specific lock primitive.
- A caller that finds the lock held retries for a short bounded time, then
  refuses with a named error that gives the lock file's path.
- A lock file left behind by an interrupted process is not recovered
  automatically. Allocation refuses and names the file; the owner removes it.
- The high-water mark is `.savepoint/task-ids.yml` (`last_issued: N`), a
  Savepoint-managed file written only while the lock is held, through the
  existing temp-write, sync, and rename pattern in `internal/data/write.go`.
- Windows evidence comes from the same tests running on the `windows-latest`
  CI job; no Windows-only file replacement is added.
- This replaces Design section 9's "No lockfile" for Task allocation only.
  T-012 reconciles Design once the behavior lands.

## Boundaries

**In scope:** durable concurrent Task-number allocation, safe creation, guidance and upgrade path, cross-Objective and failure verification.

**Out of scope:** allocation for other record kinds; Task lifecycle or Check-policy changes; archived V1 or immutable Check rewrites.

## Originating Issue

I-024 — Task planning can reuse an existing global ID. Escalation transfers repair responsibility; it is not technical clearance.
