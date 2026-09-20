---
id: E07-init-command/T001-cli-entrypoint
status: done
objective: "Create CLI entrypoint for savepoint init command"
depends_on: []
---

# T001: CLI Entrypoint

## Acceptance Criteria

- `savepoint init --help` shows usage: `init [dir] [--force] [--install]`
- `savepoint init` uses current directory
- `savepoint init <dir>` uses specified directory
- `--force` flag allows overwriting existing .savepoint
- `--install` flag triggers dev dependency install after scaffold
- Unknown flags show error and exit non-zero

## Implementation Plan

- [x] Add `cmd/init.go` with CLI arg parsing
- [x] Wire init command into main.go dispatch
- [x] Implement arg validation (dir, --force, --install)
- [x] Default to current directory if no dir specified
- [x] Test `savepoint init --help` output
- [x] Run quality gates (`make` unavailable; used direct Go build/test)

## Context Log

- Files read: `.savepoint/router.md`, `E07-Detail.md`, this task file, E07 task files T002-T007, `main.go`, `go.mod`, `Makefile`, `internal/board/board.go`, archived E05 init task.
- Files edited: `cmd/init.go`, `cmd/init_test.go`, `main.go`, this task file, `.savepoint/router.md`.
- Estimated input tokens: ~18k.
- Quality gates: `go test ./cmd` passed; `make build` could not run because `make` is not installed; `go build ./...` passed with a Go module cache warning; `go test ./...` passed; `go run . init --help` printed `Usage: init [dir] [--force] [--install]`; `go run . init --bogus` exited non-zero with an unknown flag error.

## Drift Notes

- Added new `cmd/` package for CLI command parsing and dispatch. This module is listed in the E07 epic design but is not yet listed in the AGENTS.md Codebase Map.
