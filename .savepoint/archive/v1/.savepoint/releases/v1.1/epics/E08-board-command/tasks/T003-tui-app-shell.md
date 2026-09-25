---
id: E08-board-command/T003-tui-app-shell
status: done
objective: "Implement Bubble Tea TUI program shell"
depends_on: ["E08-board-command/T001-cli-entrypoint"]
---

# T003: TUI App Shell

## Acceptance Criteria

- Launches Bubble Tea program with alt screen
- Loads theme from config.yml
- Initializes model with release/epic filters
- Handles quit (q, ctrl+c) gracefully
- Handles window resize

## Implementation Plan

- [x] Add `internal/board/tui.go`
- [x] Implement `RunTUI(release, epic) error`
- [x] Create tea.Program with board model
- [x] Enable alt screen mode
- [x] Load and apply theme from config
- [x] Handle quit key (q, ctrl+c)
- [x] Handle window resize events
- [x] Test TUI launch in terminal
- [x] Run `make build && make test`
