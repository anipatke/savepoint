---
id: E02-cross-platform-compatibility/T002-linux-build-target
status: done
objective: "Add Linux build targets for amd64 and arm64 architectures"
depends_on:
  - E02-cross-platform-compatibility/T001-fix-makefile
---

# T002: Add Linux Build Target

## Acceptance Criteria

- `make build-linux` produces linux-amd64 and linux-arm64 binaries in `dist/`
- Binaries have correct ELF format (can be verified via `file` command on Linux)
- `dist/` directory is created automatically if it doesn't exist
- Binaries are named `dist/linux-amd64/savepoint` and `dist/linux-arm64/savepoint`

## Implementation Plan

- [x] Add `build-linux` Makefile target with `GOOS=linux GOARCH=amd64` and `GOOS=linux GOARCH=arm64`
- [x] Ensure `dist/` directory creation is part of the build target
- [x] Use `go build` with output flags to place binaries in `dist/` subdirectories
- [x] Verify binaries build correctly with `make build-linux`

## Context Log

Files read:
- `Makefile`
- `.savepoint/releases/v1.1/epics/E02-cross-platform-compatibility/E02-Detail.md`
- `.savepoint/releases/v1.1/epics/E02-cross-platform-compatibility/tasks/T001-fix-makefile.md`

Estimated input tokens: 500

Notes:
- Go cross-compilation built-in — no external tooling required
- `mkdir -p` used for dist subdirs; portable across Linux, macOS, Git Bash on Windows
- Binaries verified: dist/linux-amd64/savepoint ELF 64-bit x86-64, dist/linux-arm64/savepoint ELF 64-bit ARM aarch64
- Quality gates: `go build` PASS, `go test ./...` PASS (4 packages)
