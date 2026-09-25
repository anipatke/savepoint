---
id: E16-pre-prod-refinement/T003-version-command
status: done
objective: Add a simple CLI version-reporting entry point so users can confirm the installed Savepoint build
depends_on: []
---

# T003: Add CLI Version Command

## Context Files

- `main.go` — CLI entrypoint and version reporting surface
- `cmd/` — CLI command parsing and dispatch if the version surface needs command-level handling
- `internal/buildtool/` — build metadata or release helpers if version values are already centralized there
- `README.md` — user-facing CLI usage documentation if the version command is documented elsewhere

## Acceptance Criteria

- [x] Running `savepoint --version` prints the current Savepoint version and exits successfully
- [x] Running the chosen equivalent form, if implemented, behaves consistently with `savepoint --version`
- [x] Version output is concise, script-friendly, and does not start the TUI or require a project directory
- [x] Unknown command and normal command behavior remain unchanged
- [x] Tests cover the version-reporting path and at least one normal CLI path that should not be affected
- [x] User-facing docs mention the version command if the repo already documents CLI commands
- [x] `make build && make test` passes

## Implementation Plan

- [x] Inspect the existing CLI entrypoint and command dispatch flow
- [x] Choose the smallest version surface, preferring `savepoint --version` unless existing parsing favors `savepoint version`
- [x] Wire version output before project validation, TUI startup, or command-specific work
- [x] Reuse the existing version constant or build metadata source instead of duplicating version strings
- [x] Add focused tests for version output and unchanged normal dispatch
- [x] Update CLI documentation only where command usage is already documented
- [x] Run `make build && make test`

## Context Log

Files read:
- `main.go`
- `cmd/init.go`
- `cmd/init_test.go`
- `cmd/board.go`
- `cmd/board_test.go`
- `cmd/doctor.go`
- `cmd/doctor_test.go`
- `internal/buildtool/main.go`
- `internal/buildtool/main_test.go`
- `README.md`

Files edited:
- `main_test.go`
- `README.md`
- `.savepoint/releases/v1.1/epics/E16-pre-prod-refinement/tasks/T003-version-command.md`

Estimated input tokens: ~18k

Notes:
- `savepoint --version` was already wired in `main.go` before this task; implementation added subprocess coverage around the real `main()` exit path.
- No equivalent `savepoint version` command was added; `--version` remains the single version-reporting surface.
- Quality gates passed with workspace-local Go cache: `make build` and `make test`.
