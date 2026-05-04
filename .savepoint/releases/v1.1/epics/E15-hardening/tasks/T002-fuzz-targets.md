---
id: E15-hardening/T002-fuzz-targets
status: planned
objective: Add fuzz targets for YAML frontmatter parsing round-trip
depends_on: []
---

# T002: Add Fuzz Targets for YAML Frontmatter Parsing

## Context Files

- `internal/data/parser.go` — extractFrontmatter
- `internal/data/write.go` — WriteTaskStatus, SplitFrontmatterBody
- `internal/data/parser_test.go` — existing tests

## Acceptance Criteria

- [ ] Fuzz target added for extractFrontmatter round-trip
- [ ] Fuzz target added for WriteTaskStatus round-trip (write then re-parse)
- [ ] Fuzz targets exercise edge cases (empty content, malformed YAML, unicode)
- [ ] `go test -fuzz=. ./internal/data/` runs without errors

## Implementation Plan

- [ ] Create fuzz_test.go in internal/data with fuzz targets
- [ ] Define corpus seed inputs from known edge cases
- [ ] Run fuzz targets for short duration to verify stability
- [ ] Run `make build && make test`
