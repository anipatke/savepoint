---
id: O-019
title: Allocate Task identities without planner selection
status: planned
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

## Boundaries

**In scope:** durable concurrent Task-number allocation, safe creation, guidance and upgrade path, cross-Objective and failure verification.

**Out of scope:** allocation for other record kinds; Task lifecycle or Check-policy changes; archived V1 or immutable Check rewrites.

## Originating Issue

I-024 — Task planning can reuse an existing global ID. Escalation transfers repair responsibility; it is not technical clearance.
