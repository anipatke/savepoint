---
id: E15-hardening/T005-windows-targets
status: done
objective: Add Windows amd64 and arm64 build targets to buildtool
depends_on: []
---

# T005: Add Windows Build Targets

## Context Files

- `internal/buildtool/main.go` — targets list, build logic, localExecutable

## Acceptance Criteria

- [x] Windows amd64 target added
- [x] Windows arm64 target added
- [x] .exe suffix handled for Windows binaries
- [x] Existing Linux and Darwin targets preserved
- [x] `go test ./...` passes

## Implementation Plan

- [x] Add windows-amd64 and windows-arm64 to targets list
- [x] Handle .exe suffix in build output path via `executableName(goos)`
- [x] Update localExecutable to detect Windows (already done in prior task)
- [x] Update tests for new targets
- [x] Run `make build && make test`
