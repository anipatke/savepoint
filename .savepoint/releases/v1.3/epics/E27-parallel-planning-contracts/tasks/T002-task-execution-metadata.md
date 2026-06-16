---
id: E27-parallel-planning-contracts/T002-task-execution-metadata
status: planned
objective: Parse and preserve optional parallel-group and exact write-scope task metadata.
depends_on:
  - E25-durable-data-writes/T002-atomic-data-write-paths
complexity_tier: medium
complexity_reason: Task parsing and canonical writes must add nested metadata without lifecycle regressions.
---

# T002: Task Execution Metadata

## Problem

Tasks cannot declare pre-planned parallel membership or the exact files a worker owns.

## Context Files

- `internal/data/task.go`
- `internal/data/task_test.go`
- `internal/data/parser.go`
- `internal/data/write.go`
- `internal/data/write_test.go`

## Acceptance Criteria

- [ ] Tasks may optionally define `execution.parallel_group` and `execution.write_scope`.
- [ ] Write-scope entries are normalized repo-relative paths and reject empty, absolute, escaping, or duplicate values.
- [ ] A partial execution block is invalid only when evaluated for advanced parallel launch; ordinary task loading remains backward compatible.
- [ ] `WriteTaskStatus` preserves a valid execution block byte-semantically where practical.
- [ ] Tasks without execution metadata retain identical lifecycle and rendering behavior.

## Implementation Plan

- [ ] Add typed optional execution metadata to the task model.
- [ ] Add path normalization and metadata validation helpers separate from lifecycle validation.
- [ ] Preserve the execution block through canonical status writes.
- [ ] Add parser tests for absent, valid, partial, malformed, and traversal paths.
- [ ] Add write tests proving lifecycle changes do not remove or rewrite valid execution metadata.

## Context Log

Pending.
