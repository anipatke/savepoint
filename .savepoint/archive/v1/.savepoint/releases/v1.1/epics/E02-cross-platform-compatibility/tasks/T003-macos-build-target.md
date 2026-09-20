---
id: E02-cross-platform-compatibility/T003-macos-build-target
status: done
objective: "Add macOS build targets for amd64 and arm64 architectures"
depends_on:
  - E02-cross-platform-compatibility/T001-fix-makefile
---

# T003: Add macOS Build Target

## Acceptance Criteria

- `make build-darwin` produces darwin-amd64 and darwin-arm64 binaries in `dist/`
- Binaries have correct Mach-O format
- `dist/` directory is created automatically if it doesn't exist
- Binaries are named `dist/darwin-amd64/savepoint` and `dist/darwin-arm64/savepoint`

## Implementation Plan

- [x] Add `build-darwin` Makefile target with `GOOS=darwin GOARCH=amd64` and `GOOS=darwin GOARCH=arm64`
- [x] Ensure output goes to `dist/darwin-amd64/` and `dist/darwin-arm64/`
- [x] Verify binaries build correctly with `make build-darwin`

## Context Log

Files read:
- `Makefile`
- `.savepoint/releases/v1.1/epics/E02-cross-platform-compatibility/E02-Detail.md`
- `.savepoint/releases/v1.1/epics/E02-cross-platform-compatibility/tasks/T003-macos-build-target.md`

Estimated input tokens: 400

Notes:
- Mirrors T002 pattern for darwin
- Binaries verified: dist/darwin-amd64/savepoint Mach-O 64-bit x86_64, dist/darwin-arm64/savepoint Mach-O 64-bit arm64
- Quality gates: `go build` PASS, `go test ./...` PASS (4 packages)
