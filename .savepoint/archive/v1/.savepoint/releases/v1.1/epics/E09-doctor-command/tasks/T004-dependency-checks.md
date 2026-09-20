---
id: E09-doctor-command/T004-dependency-checks
status: done
objective: "Validate dependencies for missing deps, cycles, duplicates"
depends_on: ["E09-doctor-command/T003-structure-checks"]
---

# T004: Dependency Checks

## Acceptance Criteria

- Missing dependencies detected (depends_on points to non-existent task)
- Dependency cycles detected (A→B→C→A)
- Duplicate task IDs detected
- Clear error messages for each failure

## Implementation Plan

- [x] Add to `internal/doctor/checks.go`
- [x] Implement `CheckDependencies(root, epicFilter) error`
- [x] Collect all task IDs from release
- [x] Check each depends_on references existing task
- [x] Build dependency graph, detect cycles (DFS)
- [x] Detect duplicate task IDs (same ID in multiple files)
- [x] Return clear error messages for each issue
- [x] Test dependency validation on valid and invalid projects
- [x] Run `make build && make test`