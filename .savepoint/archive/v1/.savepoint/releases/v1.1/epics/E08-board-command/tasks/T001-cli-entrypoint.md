---
id: E08-board-command/T001-cli-entrypoint
status: done
objective: "Create CLI entrypoint for savepoint board command"
depends_on: []
---

# T001: CLI Entrypoint

## Acceptance Criteria

- `savepoint board --help` shows usage: `board [--release <release>] [--epic <epic>]`
- `savepoint` (no args) defaults to board
- `savepoint board` launches TUI
- `savepoint board --release v1` filters by release
- `savepoint board --epic E03` filters by epic

## Implementation Plan

- [x] Update main.go dispatch to recognize subcommands
- [x] Add `cmd/board.go` with CLI arg parsing
- [x] Wire board command into main.go
- [x] Implement arg validation (--release, --epic)
- [x] Default to board when no subcommand given
- [x] Test `savepoint board --help` output
- [x] Run `make build && make test`