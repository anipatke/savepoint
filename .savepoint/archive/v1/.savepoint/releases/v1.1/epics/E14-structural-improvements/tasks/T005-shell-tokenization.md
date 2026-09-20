---
id: E14-structural-improvements/T005-shell-tokenization
status: done
objective: Improve splitCommand to handle single quotes and backslash-escapes, or document limitations
depends_on: []
---

# T005: Improve Shell Tokenization in Gates

## Context Files

- `internal/doctor/gates.go:125` — splitCommand function (naïve shell tokenizer)

## Acceptance Criteria

- [x] splitCommand handles single-quote grouping or limitation is documented
- [x] splitCommand handles backslash-escaped characters or limitation is documented
- [x] Existing double-quote handling preserved
- [x] `go test ./...` passes with no regressions

## Implementation Plan

- [x] Review current splitCommand implementation
- [x] Add single-quote parsing
- [x] Add backslash-escape handling
- [x] Update existing tests and add new test cases
- [x] Run `make build && make test`
