---
id: E14-structural-improvements/T005-shell-tokenization
status: planned
objective: Improve splitCommand to handle single quotes and backslash-escapes, or document limitations
depends_on: []
---

# T005: Improve Shell Tokenization in Gates

## Context Files

- `internal/doctor/gates.go:125` — splitCommand function (naïve shell tokenizer)

## Acceptance Criteria

- [ ] splitCommand handles single-quote grouping or limitation is documented
- [ ] splitCommand handles backslash-escaped characters or limitation is documented
- [ ] Existing double-quote handling preserved
- [ ] `go test ./...` passes with no regressions

## Implementation Plan

- [ ] Review current splitCommand implementation
- [ ] Add single-quote parsing
- [ ] Add backslash-escape handling
- [ ] Update existing tests and add new test cases
- [ ] Run `make build && make test`
