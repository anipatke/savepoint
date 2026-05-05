---
id: E14-structural-improvements/T007-testutil-package
status: done
objective: Create internal/testutil with shared fixtures to reduce duplication across board/doctor/data/init tests
depends_on: []
---

# T007: Consolidate Test Helpers Into Testutil Package

## Context Files

- All test files across internal/board, internal/doctor, internal/data, internal/init
- Common fixture patterns: temp directory creation, task file writing, config scaffolding

## Acceptance Criteria

- [x] internal/testutil package created
- [x] Common fixture helpers extracted and shared
- [x] All existing tests still pass after refactoring to use shared helpers
- [x] No unnecessary exports in the public API
- [x] `go test ./...` passes with no regressions

## Implementation Plan

- [x] Survey existing test helpers for duplication patterns
- [x] Create internal/testutil with shared fixture functions
- [x] Refactor each test package to use shared helpers
- [x] Run `make build && make test`
