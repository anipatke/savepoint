---
id: E09-doctor-command/T001-cli-entrypoint
status: done
objective: "Create CLI entrypoint for savepoint doctor command"
depends_on: []
---

# T001: CLI Entrypoint

## Acceptance Criteria

- `savepoint doctor --help` shows usage: `doctor [--epic <epic>]`
- `savepoint doctor` runs all checks
- `savepoint doctor --epic E03` runs checks for specific epic
- Exit codes: 0 = clean, 1 = diagnosed problems, 2 = internal error

## Implementation Plan

- [x] Add `cmd/doctor.go` with CLI arg parsing
- [x] Wire doctor command into main.go dispatch
- [x] Implement arg validation (--epic)
- [x] Implement exit code logic (0/1/2)
- [x] Test `savepoint doctor --help` output
- [x] Run `make build && make test`