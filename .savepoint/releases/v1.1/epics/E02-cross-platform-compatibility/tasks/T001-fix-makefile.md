---
id: E02-cross-platform-compatibility/T001-fix-makefile
status: done
objective: "Replace Unix-only Makefile commands with Go-native cross-platform equivalents"
depends_on: []
---

# T001: Fix Makefile for Cross-Platform Use

## Acceptance Criteria

- `make build` produces a working binary on Windows (via Git Bash), Linux, and macOS
- `make test` runs all tests on all three platforms
- `make clean` removes the binary without using `rm -f` or other Unix-only commands
- No shell commands in Makefile are platform-dependent

## Implementation Plan

- [x] Read current `Makefile` and identify all Unix-specific commands
- [x] Replace `rm -f savepoint` with `go clean` or a Go-native removal approach
- [x] Use `go build -o savepoint main.go` directly (already portable)
- [x] Test `make build && make test && make clean` works on current platform
- [x] Update `AGENTS.md` Build/Test/Run section if commands changed

## Context Log

Files read:
- `Makefile`
- `AGENTS.md`
- `.savepoint/releases/v1.1/epics/E02-cross-platform-compatibility/E02-Detail.md`

Estimated input tokens: 400

Notes:
- `go clean` removes the binary `go build` would produce — cross-platform, no-fail if missing
- `go build` and `go test` were already portable; only `rm -f` needed replacement
- Quality gates: `go build -o savepoint main.go` PASS, `go test ./...` PASS (4 packages), `go clean` PASS
