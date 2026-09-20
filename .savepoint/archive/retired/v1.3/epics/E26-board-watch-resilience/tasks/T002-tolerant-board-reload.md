---
id: E26-board-watch-resilience/T002-tolerant-board-reload
status: planned
objective: Make board reloads skip malformed files with a visible status note instead of aborting entirely.
depends_on:
  - E26-board-watch-resilience/T001-rearm-watch-after-failed-reload
complexity_tier: medium
complexity_reason: Changes loadBoardData error policy across tasks, defects, and epics with status reporting and tests.
---

# T002: Tolerant Board Reload

## Problem

`loadBoardData` propagates the first parse error, so one malformed task or defect file makes the whole board fail to load or reload — the exact trigger for the watcher-death bug, and on startup it prevents the board from opening at all. An agent saving a file in two steps routinely creates this transient state. Root cause behind audit finding H2.

## Context Files

- `internal/board/board.go`
- `internal/board/board_test.go`
- `internal/board/watch.go`
- `internal/board/update.go`

## Acceptance Criteria

- [ ] A release with one malformed task file still loads: valid tasks render normally and the board opens.
- [ ] The status line reports the number of skipped files and the first failing path (e.g. `2 files skipped: tasks/T003-foo.md …`).
- [ ] Defect and epic-status parse failures are skipped-and-reported the same way as task failures.
- [ ] Directory-level errors (releases dir missing/unreadable) still fail the load outright — only per-file parse errors are tolerated.
- [ ] A reload that previously failed wholesale now succeeds with skips, and the next file change is still observed (works with T001's re-arm).
- [ ] `go test ./internal/board` passes; existing tests for fully-valid projects are unchanged.

## Implementation Plan

- [ ] Change `loadEpicTasks`, `loadReleaseDefects`, and the epic-status read in `loadBoardData` to collect per-file errors instead of returning on the first one.
- [ ] Thread a skipped-files summary through `reloadMsg` (and initial load) into `StatusMessage`.
- [ ] Keep directory-level discovery errors fatal to avoid masking a broken project root.
- [ ] Add board tests: one malformed task among valid ones (startup and reload paths), malformed defect, and unreadable releases dir still fatal.
- [ ] Run `go test ./internal/board` and `make build && make test`.

## Context Log

Pending.
