---
id: E07-init-command/T004-atomic-writes
status: done
objective: "Implement atomic file writes with temp-file-plus-rename"
depends_on: ["E07-init-command/T003-scaffold-writer"]
---

# T004: Atomic Writes

## Acceptance Criteria

- Writes use temp file + rename pattern
- Existing files not corrupted on write failure
- Atomic operation prevents partial writes
- Handles conflicts gracefully
- Works on Windows, Linux, macOS

## Implementation Plan

- [x] Add `internal/init/write.go`
- [x] Implement `AtomicWrite(target, content) error`
- [x] Create temp file in same directory (atomic on POSIX)
- [x] Write content to temp file
- [x] Rename temp to target (atomic on success)
- [x] Clean up temp file on failure
- [x] Test atomic write behavior
- [x] Run `make build && make test`