---
id: E25-durable-data-writes/T003-router-write-search-hardening
status: planned
objective: Make WriteRouterState locate the state block without ToLower index math that breaks on non-ASCII text.
depends_on: []
complexity_tier: low
complexity_reason: Single-function fix with a targeted regression test and no interface changes.
---

# T003: Router Write Search Hardening

## Problem

`WriteRouterState` finds the "Current state" block by lowercasing the whole document and using the resulting index into the original string (`internal/data/write.go:263`). `strings.ToLower` can change byte length for non-ASCII runes, so router files containing such characters before the state block would be corrupted on write. Audit finding L5.

## Context Files

- `internal/data/write.go`
- `internal/data/write_test.go`

## Acceptance Criteria

- [ ] The state block is located case-insensitively without indexing the original string via a lowercased copy (e.g. line-wise `strings.EqualFold` scan).
- [ ] A regression test writes router state to a file containing multibyte characters (e.g. `İ`, emoji) before the state block and asserts the file round-trips uncorrupted.
- [ ] Existing router write tests pass unchanged.

## Implementation Plan

- [ ] Rework the block search in `WriteRouterState` to scan lines with `strings.EqualFold`/case-preserving matching and compute offsets from the original string only.
- [ ] Add the multibyte regression test to `internal/data/write_test.go`.
- [ ] Run `go test ./internal/data`.

## Context Log

Pending.
