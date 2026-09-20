---
id: E08-board-command/T002-non-tty-fallback
status: done
objective: "Implement non-TTY plain text table output"
depends_on: ["E08-board-command/T001-cli-entrypoint"]
---

# T002: Non-TTY Fallback

## Acceptance Criteria

- When not in a TTY, outputs deterministic plain text
- Shows three columns: planned, in_progress, done
- Lists task IDs and objectives in each column
- Includes warning about non-interactive mode
- Shows audit-entry signal if audit proposals exist

## Implementation Plan

- [x] Add `internal/board/plain.go`
- [x] Implement `RenderPlainTable(model) string`
- [x] Check `termenv.IsTerminal()` for TTY detection
- [x] Format tasks in three columns with headers
- [x] Add non-interactive warning banner
- [x] Check for audit proposals, show entry signal
- [x] Test non-TTY output
- [x] Run `make build && make test`