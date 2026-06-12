---
id: E25-durable-data-writes/T001-shared-atomic-write
status: planned
objective: Move the atomic write helper into internal/data so both data and init share one write strategy.
depends_on: []
complexity_tier: medium
complexity_reason: Relocates a cross-package helper and its tests while keeping init behavior byte-identical.
---

# T001: Shared Atomic Write Helper

## Problem

`internal/init/write.go` owns a correct atomic write (temp file + fsync + rename, with a copy fallback), but `internal/data` cannot use it without an import cycle risk, so the data layer writes non-atomically. Audit finding H1.

## Context Files

- `internal/init/write.go`
- `internal/init/write_test.go`
- `internal/data/write.go`
- `internal/data/write_test.go`

## Acceptance Criteria

- [ ] `internal/data` exports an `AtomicWrite(path string, content []byte) error` implementing temp-file + sync + rename with the existing copy fallback for rename failures.
- [ ] `internal/init` delegates to the shared helper; no duplicate atomic-write implementation remains in the repo.
- [ ] Existing `internal/init` write tests pass unchanged (behavior is byte-identical, including 0644 fallback permissions).
- [ ] New tests in `internal/data` cover successful replace, temp-file cleanup on write failure, and the rename-fallback path.
- [ ] Failed writes never leave a stray `.tmp-*` file behind in the target directory.

## Implementation Plan

- [ ] Create `internal/data/atomic.go` with `AtomicWrite` and the unexported `replaceFile` fallback, ported verbatim from `internal/init/write.go`.
- [ ] Port and adapt the relevant tests into `internal/data/atomic_test.go`.
- [ ] Replace the body of `internal/init/write.go` with thin delegation to `data.AtomicWrite` (or update its callers to import `internal/data` directly and delete the file).
- [ ] Run `go test ./internal/data ./internal/init` and `make build`.

## Context Log

Pending.
