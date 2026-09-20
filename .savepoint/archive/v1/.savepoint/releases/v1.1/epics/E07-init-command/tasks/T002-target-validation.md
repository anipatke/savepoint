---
id: E07-init-command/T002-target-validation
status: done
objective: "Implement target directory validation"
depends_on: ["E07-init-command/T001-cli-entrypoint"]
---

# T002: Target Validation

## Acceptance Criteria

- Target missing → error with helpful message
- Target empty → proceed with init
- Target has compatible files (package.json, .git, etc.) → proceed
- Target already has .savepoint → error unless --force
- Target has conflicting files → error unless --force
- Boundary errors (permission denied) → clear error message

## Implementation Plan

- [x] Add `internal/init/validate.go`
- [x] Implement `ValidateTarget(path, force) error`
- [x] Check directory exists, is writable
- [x] Check for existing .savepoint (conflict)
- [x] Check for compatible project files
- [x] Handle permission errors with clear messages
- [x] Test all validation scenarios
- [x] Run `make build && make test`