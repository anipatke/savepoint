---
id: E15-hardening/T002-fuzz-targets
status: done
objective: Add fuzz targets for YAML frontmatter parsing round-trip
depends_on: []
---

# T002: Add Fuzz Targets for YAML Frontmatter Parsing

## Context Files

- `internal/data/parser.go` — extractFrontmatter
- `internal/data/write.go` — WriteTaskStatus, SplitFrontmatterBody
- `internal/data/parser_test.go` — existing tests

## Acceptance Criteria

- [x] Fuzz target added for extractFrontmatter round-trip
- [x] Fuzz target added for WriteTaskStatus round-trip (write then re-parse)
- [x] Fuzz targets exercise edge cases (empty content, malformed YAML, unicode)
- [x] `go test -fuzz=. ./internal/data/` runs without errors

## Implementation Plan

- [x] Create fuzz_test.go in internal/data with fuzz targets
- [x] Define corpus seed inputs from known edge cases
- [x] Run fuzz targets for short duration to verify stability
- [x] Run `make build && make test`

## Drift Notes

- Fuzzer found bug in `SplitFrontmatterBody`: used `len(TrimSpace(raw))` for bodyStart, causing off-by-N when frontmatter has leading/trailing whitespace. Fixed to compute bodyStart from actual `\n---` offset.
- Fuzzer found `normalizeLineEndings` not idempotent with `\r\r\n` sequences. Fixed to also replace lone `\r` with `\n` (handles legacy Mac CR line endings).
- No new files/modules beyond `internal/data/fuzz_test.go` (in-scope per E15-Detail.md).
