---
id: E15-hardening/T007-root-test-allowlist
status: planned
objective: Move agent_skills_test.go to a test package and extract audit allowlist to a named constant
depends_on: []
---

# T007: Move Root-Level Test and Extract Audit Allowlist

## Context Files

- `agent_skills_test.go` — root-level package main test
- `internal/board/epic_panel.go:116-119` — allowedSections map

## Acceptance Criteria

- [ ] agent_skills_test.go moved to cmd_test package (or appropriate location)
- [ ] allowedSections extracted to a named constant with documentation
- [ ] All existing tests still pass after refactoring
- [ ] `go test ./...` passes with no regressions

## Implementation Plan

- [ ] Move agent_skills_test.go to an appropriate internal test package
- [ ] Update imports and paths in moved test
- [ ] Extract allowedSections map to a named constant
- [ ] Update any references in epic_panel.go
- [ ] Run `make build && make test`
