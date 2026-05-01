---
type: epic-design
status: planned
---

# Epic E01: Go Project Setup

## Purpose

Initialize the savepoint CLI as a Go module. Set up the project directory structure, add Bubble Tea + Lip Gloss dependencies, and establish a working build/test loop.

## What this epic adds

- `go.mod` with module path `github.com/opencode/savepoint`.
- `main.go` entrypoint that launches the Bubble Tea program.
- Directory structure: `cmd/`, `internal/board/`, `internal/data/`, `internal/styles/`.
- `go.sum` with resolved Bubble Tea, Lip Gloss, reflow, and YAML dependencies.
- `Makefile` or script for `build`, `test`, `run`.

## Definition of Done

- `go mod init` succeeds and `go mod tidy` produces clean `go.mod`/`go.sum`.
- `go build` produces a working binary.
- `go test ./...` runs (even if no tests yet, it must not error).
- Bubble Tea `tea.NewProgram` launches and renders a blank screen.
- No TypeScript files, `package.json`, or `node_modules` remain at project root.

## Components and files

| Path | Purpose |
|------|---------|
| `go.mod` | Go module definition |
| `go.sum` | Dependency checksums |
| `main.go` | Entrypoint: tea.NewProgram |
| `Makefile` | Build/test/run shortcuts |
| `cmd/` | CLI command setup |
| `internal/board/` | Board TUI package |
| `internal/data/` | Data reading/parsing package |
| `internal/styles/` | Lip Gloss styles package |
