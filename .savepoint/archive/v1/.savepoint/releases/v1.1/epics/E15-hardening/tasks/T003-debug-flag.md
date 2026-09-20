---
id: E15-hardening/T003-debug-flag
status: done
objective: Add --debug flag and SAVEPOINT_DEBUG env var for structured debug logging
depends_on: []
---

# T003: Add --debug Flag and SAVEPOINT_DEBUG Env Var

## Context Files

- `main.go` — CLI entrypoint, --version flag
- `cmd/` — command dispatch
- `internal/board/` — board init, file watchers, command dispatch

## Acceptance Criteria

- [x] --debug flag accepted at CLI level
- [x] SAVEPOINT_DEBUG env var recognized
- [x] Debug output includes board init, file watcher events, and command dispatch
- [x] Debug output is off by default (no performance impact)
- [x] `go test ./...` passes with no regressions

## Implementation Plan

- [x] Define debug logging helper in a shared location
- [x] Add --debug flag parsing in main.go
- [x] Instrument key points in board init, file watcher, and update dispatch
- [x] Add tests for debug flag behavior
- [x] Run `make build && make test`
