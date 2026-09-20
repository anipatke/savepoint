---
type: epic-design
status: audited
---

# Epic E02: Cross-Platform Compatibility

## Purpose

Ensure the project builds, runs, and distributes cleanly across Windows, Linux, and macOS. The codebase itself is platform-agnostic (pure Go with `filepath`), but the build tooling and release workflow need adjustments for multi-platform support.

## Definition of Done

- Makefile works identically across Windows (Git Bash/WSL), Linux, and macOS
- Binary builds succeed for linux-amd64, linux-arm64, darwin-amd64, darwin-arm64
- Smoke tests pass on at least one non-Windows target
- Release artifacts are produced in `dist/` with platform-specific naming

## Components and files

| Path | Purpose |
|------|---------|
| `Makefile` | Build tooling entrypoints that delegate platform-sensitive work to Go |
| `internal/buildtool/` | Go-native build helper for cleanup, cross-compilation, archives, and smoke tests |
| `dist/` | New directory for release artifacts (platform-organized) |

## Architectural notes

- All Go code uses `path/filepath` — no syscalls are platform-specific
- Bubble Tea handles terminal differences across platforms
- No Go build tags or platform-specific files needed
- This is purely a build-tooling and release-workflow epic

## Implemented as

- `Makefile` now has `build-linux`, `build-darwin`, `build-all`, `dist`, and `smoke-test` targets that delegate to `internal/buildtool`.
- `internal/buildtool` uses Go APIs for directory creation, cleanup, tar.gz archive creation, and local smoke-test execution.
- `VERSION` can be passed through `make VERSION=...`; otherwise the helper defaults from `git describe --tags --abbrev=0` with a `v0.0.0` fallback.
- `main.go` handles `--version` before launching the TUI, which gives smoke tests a non-interactive success path.
- `dist/` remains ignored by git and contains both raw cross-compiled binaries and versioned `.tar.gz` archives for Linux and Darwin builds.
- Windows zip packaging was not implemented in this epic; the final implementation covers the Linux and Darwin targets named in the epic's build-all definition of done.

## Definition of Done

- [x] Makefile replaced with Go-native commands (no `rm`, `cp`, `mkdir`)
- [x] `dist/` directory convention established
- [x] `make build-all` produces binaries for linux-amd64, linux-arm64, darwin-amd64, darwin-arm64
- [x] `make dist` creates versioned tar/zip artifacts per platform
- [x] Smoke test script validates binary runs and exits cleanly
