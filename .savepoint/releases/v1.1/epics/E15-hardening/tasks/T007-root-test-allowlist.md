---
id: E15-hardening/T007-root-test-allowlist
status: done
objective: Move agent_skills_test.go to a test package and extract audit allowlist to a named constant
depends_on: []
---

# T007: Move Root-Level Test and Extract Audit Allowlist

## Context Files

- `agent_skills_test.go` — root-level package main test
- `internal/board/epic_panel.go:116-119` — allowedSections map

## Acceptance Criteria

- [x] agent_skills_test.go moved to cmd_test package (or appropriate location)
- [x] allowedSections extracted to a named constant with documentation
- [x] All existing tests still pass after refactoring
- [x] `go test ./...` passes with no regressions

## Implementation Plan

- [x] Move agent_skills_test.go to an appropriate internal test package
- [x] Update imports and paths in moved test
- [x] Extract allowedSections map to a named constant
- [x] Update any references in epic_panel.go
- [x] Run `make build && make test`

## Notes

- `agent_skills_test.go` changed from `package main` → `package main_test` (external test package at root). Paths remain valid since Go test working dir is the package dir.
- `epicAuditHiddenSectionHeadings` already extracted as package-level var; added doc comment explaining suppression rationale.
