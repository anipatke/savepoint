---
id: E15-hardening/T006-abbreviation-splitting
status: planned
objective: Fix splitChecklistSentences to skip periods preceded by known abbreviations
depends_on: []
---

# T006: Fix Abbreviation Handling in Checklist Sentence Splitting

## Context Files

- `internal/board/detail.go:88` — splitChecklistSentences function

## Acceptance Criteria

- [ ] Known abbreviations (e.g., "e.g.", "i.e.") do not trigger sentence breaks
- [ ] Existing sentence splitting behavior preserved for non-abbreviation periods
- [ ] Abbreviation list is configurable/extensible
- [ ] `go test ./...` passes with no regressions

## Implementation Plan

- [ ] Define known abbreviation set
- [ ] Add abbreviation check before period-based sentence split
- [ ] Add test cases for abbreviations in sentences
- [ ] Run `make build && make test`
