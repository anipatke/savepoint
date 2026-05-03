---
id: E13-audit-remediation/T007-test-coverage
status: done
objective: Add tests for untested packages (buildtool, styles), convert repair matching to typed errors, fix minor test issues
depends_on:
    - E13-audit-remediation/T001-safe-cleanup
---

# T007: Test Coverage — Buildtool, Styles, Typed Repairs, Minor Fixes

## Context Files

- `internal/buildtool/main.go` — zero test files
- `internal/styles/palette.go` — zero test files
- `internal/styles/styles.go` — zero test files
- `internal/doctor/repairs.go` — substring-based error matching fragile
- `internal/data/errors.go` — existing sentinel errors
- `agent_skills_test.go` — hardcodes expected skill count of 6

## Acceptance Criteria

- [ ] `internal/buildtool/main.go` has test coverage for `version()`, `splitCommand()` (if exists), and `localExecutable()`
- [ ] `internal/styles/` has a test file verifying palette constants and color completeness
- [ ] `SuggestRepair` in `repairs.go` accepts `error` and uses `errors.Is()` with sentinel errors from `data/errors.go` (and any new ones needed)
- [ ] New sentinel errors added to `data/errors.go` for repair-matching: at minimum `ErrInvalidStatus`, `ErrMissingFrontmatter`, `ErrConfigNotFound`, `ErrStructureProblem`
- [ ] `agent_skills_test.go` no longer hardcodes the skill count (derives it from directory listing)
- [ ] `go test ./...` passes

## Implementation Plan

- [ ] Create `internal/buildtool/main_test.go` with tests for `version()` (env var, flag override, git tag fallback) and `localExecutable()`
- [ ] Create `internal/styles/styles_test.go` verifying all 3 color tiers exist for each named constant, and `color()` produces correct CompleteColor
- [ ] Add new sentinel errors to `data/errors.go` for doctor repair matching
- [ ] Refactor `SuggestRepair` to accept `error` and use `errors.Is()` instead of substring matching
- [ ] Update `repairs_test.go` to test typed error matching
- [ ] Fix `agent_skills_test.go` to derive expected skill count from directory listing
- [ ] Run `make build && make test`