---
type: epic-design
status: audited
---

# Epic E02: Data Readers

## Purpose

Read and parse the `.savepoint/` file structure in Go. Parse YAML frontmatter from task files, read router state, config theme, and discover releases/epics/tasks on disk.

## What this epic adds

- Markdown frontmatter parser (`---\n...\n---`) using `gopkg.in/yaml.v3`.
- Task struct with phase support.
- Router state reader (`router.md` yaml block extraction).
- Config reader (`config.yml` theme parsing).
- Release/epic/task file discovery (walk directories).
- Error types for missing files, malformed yaml, schema rejection.

## Definition of Done

- Can read any task `.md` file and parse frontmatter into a Task struct.
- Can read `router.md` and extract the Current state yaml block.
- Can read `config.yml` and parse theme colors.
- Can list all releases, epics within a release, and tasks within an epic.
- All parsing functions have unit tests.
- `go test ./internal/data/...` passes.

## Components and files

## Implemented As

- `internal/data/task.go` defines task metadata, column, progress stage, dependency, and progress structs.
- `internal/data/parser.go` extracts YAML frontmatter from markdown and maps task ID, objective/title, description, status/column, phase/stage, dependency, priority, points, tags, acceptance, notes, release, epic, and progress fields into `Task`.
- `internal/data/router.go` extracts and parses the `## Current state` fenced YAML block from `router.md`.
- `internal/data/config.go` reads `config.yml` theme values, fills missing theme fields from defaults, returns defaults when the file is missing, and surfaces malformed YAML/read errors.
- `internal/data/discover.go` finds the `.savepoint` root and lists release, epic, and task directories in stable order.
- `internal/data/errors.go` defines shared parser/discovery boundary error sentinels.

## Audit Reconciliation

- Fixed parser field coverage for `status`, `phase`, `depends_on`, and the remaining task metadata fields covered by the `Task` struct.
- Fixed config boundary behavior so missing files default while malformed YAML and read errors are returned.
- Replaced discovery tests that depended on the live repository with deterministic temporary fixtures.
- Added parser tests for malformed YAML and task metadata mapping.
- Added the planned `internal/data/errors.go` boundary error file.

| Path | Purpose |
|------|---------|
| `internal/data/task.go` | Task struct, frontmatter parsing |
| `internal/data/parser.go` | Markdown YAML frontmatter extraction |
| `internal/data/router.go` | Router state reader |
| `internal/data/config.go` | Config/theme reader |
| `internal/data/discover.go` | Directory walking, file discovery |
| `internal/data/errors.go` | Error types |
| `internal/data/parser_test.go` | Frontmatter parser tests |
| `internal/data/config_test.go` | Config parsing tests |
| `internal/data/task_test.go` | Task parsing tests |
| `internal/data/router_test.go` | Router parsing tests |
| `internal/data/discover_test.go` | Discovery tests |
