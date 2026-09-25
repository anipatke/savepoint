---
id: E02-data-readers/T004-config-reader
status: done
objective: "Read config.yml and parse theme colors"
depends_on: [E02-data-readers/T003-router-reader]
---

# T004: Config Reader

## Acceptance Criteria

- Reads `.savepoint/config.yml`.
- Parses theme colors (bg, surface, surface_2, border, text, accents).
- Returns default theme if config.yml missing.
- Tests cover valid config and missing file.

## Implementation Plan

- [x] Create `internal/data/config.go`.
- [x] Define Config and Theme structs.
- [x] Parse YAML with defaults.
- [x] Write `internal/data/config_test.go`.
- [x] Run `go test`.

## Context Log

Files read: `internal/data/config.go`, `internal/data/config_test.go`
Estimated input tokens: ~900
Notes: Audit closeout limited default fallback to missing config files, returned malformed YAML/read errors, and added partial-default and malformed YAML coverage. `go test ./internal/data/...`, `go test ./...`, and `go build ./...` passed.
