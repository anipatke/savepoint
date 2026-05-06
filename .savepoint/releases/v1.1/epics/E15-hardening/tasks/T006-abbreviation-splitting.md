---
id: E15-hardening/T006-abbreviation-splitting
status: done
objective: Fix splitChecklistSentences to skip periods preceded by known abbreviations
depends_on: []
---

# T006: Fix Abbreviation Handling in Checklist Sentence Splitting

## Context Files

- `internal/board/detail.go:88` — splitChecklistSentences function

## Acceptance Criteria

- [x] Known abbreviations (e.g., "e.g.", "i.e.") do not trigger sentence breaks
- [x] Existing sentence splitting behavior preserved for non-abbreviation periods
- [x] Abbreviation list is configurable/extensible
- [x] `go test ./...` passes with no regressions

## Implementation Plan

- [x] Define known abbreviation set
- [x] Add abbreviation check before period-based sentence split
- [x] Add test cases for abbreviations in sentences
- [x] Run `make build && make test`
