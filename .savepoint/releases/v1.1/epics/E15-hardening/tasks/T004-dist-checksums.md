---
id: E15-hardening/T004-dist-checksums
status: done
objective: Generate checksums.txt during make dist using Go crypto APIs
depends_on: []
---

# T004: Add Distribution Checksums

## Context Files

- `internal/buildtool/main.go` — dist() function that creates tar.gz archives
- `Makefile` — dist target

## Acceptance Criteria

- [x] dist() generates checksums.txt file with SHA256 hashes
- [x] checksums.txt is included in the dist directory
- [x] Existing archive creation behavior preserved
- [x] `go test ./...` passes

## Implementation Plan

- [x] Add SHA256 checksum computation in dist() using crypto/sha256
- [x] Write checksums.txt to dist directory
- [x] Add tests for checksum generation
- [x] Run `make build && make test`
