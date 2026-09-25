---
id: E25-durable-data-writes/T002-atomic-data-write-paths
status: planned
objective: Route all five internal/data write paths through the shared atomic write helper.
depends_on:
  - E25-durable-data-writes/T001-shared-atomic-write
complexity_tier: low
complexity_reason: Mechanical substitution at five known call sites guarded by an existing test suite.
---

# T002: Atomic Data Write Paths

## Problem

`ApplyProposal`, `updateFrontmatterField`, `WriteTaskStatus`, `WriteDefectStatus`, and `WriteRouterState` in `internal/data/write.go` use plain `os.WriteFile`, which truncates before writing. A crash, kill, or full disk mid-write corrupts the task, defect, or router file — the project's source of truth. Audit finding H1.

## Context Files

- `internal/data/write.go`
- `internal/data/write_test.go`
- `internal/data/atomic.go`

## Acceptance Criteria

- [ ] All five write sites in `internal/data/write.go` call `AtomicWrite`; no direct `os.WriteFile` call remains in the file.
- [ ] Mtime conflict detection still works: a successful write produces a new mtime and `ErrMtimeConflict` behavior is unchanged.
- [ ] All existing `internal/data` write tests pass without weakening any assertion.
- [ ] A new test simulates a write failure (e.g. unwritable temp dir or injected failure) and asserts the original file content is intact afterwards.

## Implementation Plan

- [ ] Replace each `os.WriteFile(path, …, 0644)` in `internal/data/write.go` with `AtomicWrite(path, …)`.
- [ ] Verify the board's conflict-retry path (`internal/board/io.go`) still passes its tests, since rename changes mtime semantics slightly on some filesystems.
- [ ] Add the corruption-safety regression test to `internal/data/write_test.go`.
- [ ] Run `go test ./internal/data ./internal/board` and `make build && make test`.

## Context Log

Pending.
