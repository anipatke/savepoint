---
id: E15-hardening/T005-windows-targets
status: planned
objective: Add Windows amd64 and arm64 build targets to buildtool
depends_on: []
---

# T005: Add Windows Build Targets

## Context Files

- `internal/buildtool/main.go` — targets list, build logic, localExecutable

## Acceptance Criteria

- [ ] Windows amd64 target added
- [ ] Windows arm64 target added
- [ ] .exe suffix handled for Windows binaries
- [ ] Existing Linux and Darwin targets preserved
- [ ] `go test ./...` passes

## Implementation Plan

- [ ] Add windows-amd64 and windows-arm64 to targets list
- [ ] Handle .exe suffix in build output path
- [ ] Update localExecutable to detect Windows
- [ ] Update tests for new targets
- [ ] Run `make build && make test`
