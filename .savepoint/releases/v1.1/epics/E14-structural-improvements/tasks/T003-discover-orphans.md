---
id: E14-structural-improvements/T003-discover-orphans
status: done
objective: Add ListRootDirs to Discover and use it in CheckOrphans instead of direct os.ReadDir
depends_on: []
---

# T003: Route CheckOrphans Through Discover

## Context Files

- `internal/doctor/checks.go:488-514` — CheckOrphans uses os.ReadDir directly
- `internal/data/discover.go` — existing ListReleases, ListEpics, ListTasks methods

## Acceptance Criteria

- [x] ListRootDirs(releasesPath string) ([]string, error) added to Discover
- [x] CheckOrphans uses data.Discover for filesystem traversal
- [x] `go test ./...` passes with no regressions

## Implementation Plan

- [x] Add ListRootDirs method to Discover in discover.go
- [x] Add tests for ListRootDirs in discover_test.go
- [x] Refactor CheckOrphans to use Discover instead of os.ReadDir
- [x] Run `make build && make test`

## Context Log

- Files read: `internal/doctor/checks.go`, `internal/data/discover.go`, `internal/data/discover_test.go`, `internal/doctor/interfaces.go`.
- Files edited: `internal/data/discover.go`, `internal/data/discover_test.go`, `internal/doctor/interfaces.go`, `internal/doctor/interfaces_test.go`, `internal/doctor/checks.go`.
- Token estimate: ~8k.
- Quality gates: `go test ./internal/data ./internal/doctor` passed; `make build && make test` passed (`go test ./...` all packages).
