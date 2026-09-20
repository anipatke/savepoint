---
id: E17-defect-workflow-tui/T001-defect-data-model
status: done
objective: Add the release-level defect data model, parser, discovery, and validation primitives
depends_on: []
---

# T001: Defect Data Model And Discovery

## Problem

Defects need to be represented as first-class release repair artifacts without overloading task files. The rest of the system needs a typed data boundary before board, doctor, or router behavior can safely consume defect files.

## Context Files

- `internal/data/task.go` - existing task model and parsed markdown fields
- `internal/data/frontmatter.go` - shared frontmatter parsing and splitting helpers
- `internal/data/discover.go` - release, epic, and task discovery patterns
- `internal/data/status.go` - canonical status and stage constants
- `internal/data/write.go` - canonical write validation patterns
- `internal/data/task_test.go` - parser and lifecycle validation test style
- `internal/data/discover_test.go` - discovery fixture patterns

## Acceptance Criteria

- [ ] A typed `Defect` model represents id, release, status, severity, optional introduced/reference fields, stage, title, and body sections
- [ ] Defects are discovered from `.savepoint/releases/{release}/defects/*.md`
- [ ] Missing `defects/` directories are treated as zero defects, not project corruption
- [ ] Defect status is limited to `planned`, `in_progress`, and `done`
- [ ] `stage` is required when defect status is `in_progress`
- [ ] Defect parsing preserves the sections needed by the TUI detail view
- [ ] Unit tests cover valid defects, missing directories, malformed frontmatter, invalid statuses, and missing in-progress stage
- [ ] `make build && make test` passes

## Implementation Plan

- [x] Add `Defect` and any small supporting types in `internal/data`
- [x] Reuse existing frontmatter/body parsing helpers instead of duplicating markdown splitting
- [x] Implement release-scoped defect discovery with deterministic filename ordering
- [x] Add validation for allowed statuses and in-progress stage requirements
- [x] Add focused parser and discovery tests
- [x] Run `make build && make test`

## Context Log

- Files read: task.go, discover.go, write.go, lifecycle.go, parser.go, errors.go, task_test.go, discover_test.go, dependency.go
- Files edited: parser.go, discover.go, lifecycle.go
- Files created: defect.go, defect_test.go
- Quality gates: make build && make test PASS

