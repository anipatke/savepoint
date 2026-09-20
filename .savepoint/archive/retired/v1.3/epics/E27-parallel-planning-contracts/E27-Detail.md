---
type: epic-design
status: planned
---

# E27: Parallel Planning Contracts

## Purpose

Add the optional configuration and task metadata needed to prove that planned tasks can run safely in separate worktrees before any worktree is created.

## What this epic adds

- Disabled-by-default advanced parallel-worktree configuration nested under the launcher.
- Optional task execution metadata for a parallel group and exact repo-relative write scope.
- Canonical validation for dependency readiness, duplicate paths, and overlapping task ownership.
- Planning-skill guidance that creates conservative parallel groups and falls back to sequential execution when ownership is unclear.

## Components and files

| Module | Purpose |
|--------|---------|
| `internal/data/config.go` | Parse and validate advanced parallel-worktree configuration |
| `internal/data/task.go` | Parse optional task execution metadata |
| `internal/data/write.go` | Preserve execution metadata through lifecycle writes |
| `internal/data/parallel.go` | Own parallel eligibility and exact write-scope conflict rules |
| `agent-skills/savepoint-create-task/SKILL.md` | Define planning rules for parallel groups and write scopes |
| `templates/project/agent-skills/savepoint-create-task/SKILL.md` | Ship the same planning contract to generated projects |

## Architectural delta

Task planning gains an optional execution contract without changing the task lifecycle. `internal/data` remains the source of truth for parsing and validation, while advanced execution consumers receive a typed eligibility result instead of reimplementing dependency and scope rules.

## Boundaries

**In scope:**
- Optional configuration and task metadata
- Exact path normalization and overlap detection
- Dependency-aware parallel eligibility
- Planning skill and template guidance

**Out of scope:**
- Git commands or worktree creation
- Board actions
- Runtime manifests
- Automatic grouping inferred from source code

## Quality gates

- Existing task/config fixtures parse unchanged when advanced fields are absent.
- Table tests cover valid groups, missing metadata, unsatisfied dependencies, duplicate paths, and parent/child collisions.
- `go test ./internal/data` passes.
- `make build && make test` passes.

## Open decisions

None. Write scopes use exact repo-relative paths; ambiguous or generated shared ownership remains sequential.
