---
id: E02-data-readers/T003-router-reader
status: done
objective: "Read and parse router.md Current state block"
depends_on: [E02-data-readers/T002-frontmatter-parser]
---

# T003: Router Reader

## Acceptance Criteria

- Reads `router.md` and extracts the `## Current state` yaml block.
- Parses state, release, epic, task, next_action fields.
- Returns error if router.md missing or block not found.
- Tests cover valid and missing router cases.

## Implementation Plan

- [x] Create `internal/data/router.go`.
- [x] Implement `## Current state` block regex extraction.
- [x] Parse yaml block into RouterState struct.
- [x] Write `internal/data/router_test.go`.
- [x] Run `go test`.

## Context Log

Files read: `internal/data/router.go`, `internal/data/router_test.go`
Estimated input tokens: ~700
Notes: Verified during E02 audit closeout. `go test ./internal/data/...`, `go test ./...`, and `go build ./...` passed.
