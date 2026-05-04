---
id: E14-structural-improvements/T007-testutil-package
status: planned
objective: Create internal/testutil with shared fixtures to reduce duplication across board/doctor/data/init tests
depends_on: []
---

# T007: Consolidate Test Helpers Into Testutil Package

## Context Files

- All test files across internal/board, internal/doctor, internal/data, internal/init
- Common fixture patterns: temp directory creation, task file writing, config scaffolding

## Acceptance Criteria

- [ ] internal/testutil package created
- [ ] Common fixture helpers extracted and shared
- [ ] All existing tests still pass after refactoring to use shared helpers
- [ ] No unnecessary exports in the public API
- [ ] `go test ./...` passes with no regressions

## Implementation Plan

- [ ] Survey existing test helpers for duplication patterns
- [ ] Create internal/testutil with shared fixture functions
- [ ] Refactor each test package to use shared helpers
- [ ] Run `make build && make test`
