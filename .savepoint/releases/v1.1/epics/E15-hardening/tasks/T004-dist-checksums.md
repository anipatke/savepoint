---
id: E15-hardening/T004-dist-checksums
status: planned
objective: Generate checksums.txt during make dist using Go crypto APIs
depends_on: []
---

# T004: Add Distribution Checksums

## Context Files

- `internal/buildtool/main.go` — dist() function that creates tar.gz archives
- `Makefile` — dist target

## Acceptance Criteria

- [ ] dist() generates checksums.txt file with SHA256 hashes
- [ ] checksums.txt is included in the dist directory
- [ ] Existing archive creation behavior preserved
- [ ] `go test ./...` passes

## Implementation Plan

- [ ] Add SHA256 checksum computation in dist() using crypto/sha256
- [ ] Write checksums.txt to dist directory
- [ ] Add tests for checksum generation
- [ ] Run `make build && make test`
