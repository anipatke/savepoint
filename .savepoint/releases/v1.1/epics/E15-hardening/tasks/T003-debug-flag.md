---
id: E15-hardening/T003-debug-flag
status: planned
objective: Add --debug flag and SAVEPOINT_DEBUG env var for structured debug logging
depends_on: []
---

# T003: Add --debug Flag and SAVEPOINT_DEBUG Env Var

## Context Files

- `main.go` — CLI entrypoint, --version flag
- `cmd/` — command dispatch
- `internal/board/` — board init, file watchers, command dispatch

## Acceptance Criteria

- [ ] --debug flag accepted at CLI level
- [ ] SAVEPOINT_DEBUG env var recognized
- [ ] Debug output includes board init, file watcher events, and command dispatch
- [ ] Debug output is off by default (no performance impact)
- [ ] `go test ./...` passes with no regressions

## Implementation Plan

- [ ] Define debug logging helper in a shared location
- [ ] Add --debug flag parsing in main.go
- [ ] Instrument key points in board init, file watcher, and update dispatch
- [ ] Add tests for debug flag behavior
- [ ] Run `make build && make test`
