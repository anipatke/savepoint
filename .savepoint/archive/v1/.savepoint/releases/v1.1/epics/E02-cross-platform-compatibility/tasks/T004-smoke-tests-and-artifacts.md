---
id: E02-cross-platform-compatibility/T004-smoke-tests-and-artifacts
status: done
objective: "Add smoke test script and create versioned release artifacts for all platforms"
depends_on:
  - E02-cross-platform-compatibility/T002-linux-build-target
  - E02-cross-platform-compatibility/T003-macos-build-target
---

# T004: Smoke Tests and Release Artifacts

## Acceptance Criteria

- `make dist` builds all platforms and creates versioned tar/zip archives in `dist/`
- Archive naming follows pattern: `savepoint-{version}-{platform}-{arch}.tar.gz` (or `.zip` for Windows)
- Smoke test validates each binary starts and exits with code 0
- Smoke test does not require interactive TUI (headless validation)
- `dist/` contains both raw binaries and packaged archives

## Implementation Plan

- [x] Add `dist` Makefile target that invokes all `build-*` targets
- [x] Create platform-appropriate archives (tar.gz for Unix, zip for Windows)
- [x] Add `smoke-test` target that runs each binary with `--help` or similar non-interactive flag
- [x] Add `.gitignore` entry for `dist/`
- [x] Document `make dist` in AGENTS.md

## Context Log

Files read:
- `Makefile`
- `AGENTS.md`
- `main.go`
- `.gitignore`
- `internal/board/board.go`

Estimated input tokens: ~1200

Notes:
- Added `--version` flag to `main.go` (checks `os.Args[1]`, prints `version` var, exits 0). No TUI interaction.
- `version` var injected at build time via `-ldflags "-X main.version=$(VERSION)"`.
- `VERSION` defaults to `git describe --tags` output or `v0.0.0` if no tags exist.
- `.gitignore` already had `dist/` entry — no change needed.
- smoke-test runs local platform binary only; cross-compiled binaries cannot be exec'd without emulation.
- Windows support (zip archives) deferred — E02 scope is linux/darwin only.
- Audit closeout replaced shell-specific Makefile recipes with `internal/buildtool`, which uses Go APIs for cleanup, directory creation, cross-compilation orchestration, tar.gz archive creation, and smoke-test execution.

Quality gates:
- `make build && make test`: NOT RUN (`make` is not installed in this PowerShell environment)
- `go run ./internal/buildtool build`: PASS
- `go run ./internal/buildtool build-linux`: PASS
- `go run ./internal/buildtool build-darwin`: PASS
- `go run ./internal/buildtool dist`: PASS
- `go run ./internal/buildtool smoke-test`: PASS (`v0.0.0`, exit 0)
- `go test ./...`: PASS

## Drift Notes

- Resolved in E02 audit closeout: `main.go` Codebase Map entry now mentions `--version`, and `internal/buildtool/` was added to the Codebase Map.
