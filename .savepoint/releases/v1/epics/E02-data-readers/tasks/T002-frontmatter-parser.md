---
id: E02-data-readers/T002-frontmatter-parser
status: done
objective: "Parse YAML frontmatter from markdown task files"
depends_on: [E02-data-readers/T001-task-struct]
---

# T002: Frontmatter Parser

## Acceptance Criteria

- Can parse `---\nyaml\n---\nbody` format from any `.md` file.
- Extracts id, status, phase, objective, depends_on from frontmatter.
- Returns structured Task or error for malformed files.
- Tests cover valid, missing frontmatter, malformed yaml cases.

## Implementation Plan

- [x] Create `internal/data/parser.go`.
- [x] Implement frontmatter regex extraction.
- [x] Parse YAML with `gopkg.in/yaml.v3`.
- [x] Map parsed fields to Task struct.
- [x] Write `internal/data/parser_test.go`.
- [x] Run `go test`.

## Context Log

Files read: `internal/data/parser.go`, `internal/data/parser_test.go`, `internal/data/task.go`
Estimated input tokens: ~1,500
Notes: Audit closeout replaced regex extraction with delimiter-based extraction, mapped task metadata fields including `status`, `phase`, and `depends_on`, and added malformed YAML coverage. `go test ./internal/data/...`, `go test ./...`, and `go build ./...` passed.
