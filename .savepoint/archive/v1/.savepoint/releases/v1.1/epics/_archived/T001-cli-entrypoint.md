---
id: E06-audit-command/T001-cli-entrypoint
status: in_progress
objective: "Create CLI entrypoint for savepoint audit command"
depends_on: []
---

# T001: CLI Entrypoint

## Acceptance Criteria

- `savepoint audit --help` shows usage: `audit <epic-id|release> [--skip --reason]`
- `savepoint audit <epic-id>` runs audit for that epic
- `savepoint audit <release>` runs audit for all unaudited epics in that release
- `savepoint audit --skip --reason "..."` skips the audit and logs reason
- Unknown arguments show error and exit non-zero

## Implementation Plan

- [x] Add `cmd/audit.go` with CLI argument parsing
- [x] Wire audit command into main.go
- [x] Implement arg validation for epic-id (E##) and release (v#.#) formats
- [x] Add subcommand dispatch to audit workflow
- [x] Test `savepoint audit --help` output
- [x] Run `make build && make test`