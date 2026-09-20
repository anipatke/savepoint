---
id: E02-data-readers/T005-discovery
status: done
objective: "Walk .savepoint directory to discover releases, epics, and tasks"
depends_on: [E02-data-readers/T004-config-reader]
---

# T005: Discovery

## Acceptance Criteria

- Can list all releases under `.savepoint/releases/`.
- Can list all epics within a release.
- Can list all tasks within an epic.
- Can find `.savepoint` root by walking up from cwd.
- Tests cover discovery logic.

## Implementation Plan

- [x] Create `internal/data/discover.go`.
- [x] Implement `FindSavepointRoot(start string)`.
- [x] Implement `ListReleases(root)`, `ListEpics(root, release)`, `ListTasks(root, release, epic)`.
- [x] Write `internal/data/discover_test.go`.
- [x] Run `go test`.

## Context Log

Files read: `internal/data/discover.go`, `internal/data/discover_test.go`, `internal/data/errors.go`
Estimated input tokens: ~800
Notes: discover_test.go had unused `os`/`path/filepath` imports causing build failure — removed. Audit closeout moved discovery coverage to deterministic temp-directory fixtures and added a shared missing-root error sentinel. All 4 discovery tests pass. `go test ./internal/data/...`, `go test ./...`, and `go build ./...` passed.
