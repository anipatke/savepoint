---
type: epic-design
status: planned
---

# E34: Audit Register Doctor and Quality Gates

## Purpose

Make audit-register state diagnosable through `savepoint doctor` and prove the release with focused regression coverage across scaffold, data, board, and documentation workflows.

## What this epic adds

- Doctor checks for malformed audit-register files.
- Report output and repair suggestions for invalid statuses, missing proof, broken links, and stale duplicates.
- Regression coverage that proves projects without an audit register remain valid.
- Release quality gate coverage for the new templates, data model, and TUI review section.

## Components and files

| Module | Purpose |
|--------|---------|
| `internal/doctor` | Validate and report audit-register consistency |
| `cmd/doctor.go` | Preserve CLI doctor behavior while adding audit diagnostics |
| `internal/data` | Provide audit validation results consumed by doctor |
| `internal/board` | Provide TUI regression coverage for the read-only overlay |
| `internal/init` | Provide scaffold and upgrade regression coverage |

## Architectural delta

Doctor gains optional audit-register diagnostics. Absence of `.savepoint/audit/` is not an error; malformed present files are reported with actionable messages.

## Boundaries

**In scope:**
- Doctor checks and report output
- Typed repair suggestions
- Regression matrix across affected packages
- Release quality gate documentation

**Out of scope:**
- Doctor auto-repair of audit files
- Running audits
- Verifying source code fixes for individual findings
- Blocking projects that have not adopted the register

## Quality gates

- `go test ./internal/doctor ./internal/data ./internal/board ./internal/init` passes.
- `make build && make test` passes before release audit.
- Doctor output stays stable for projects without audit-register files.

## Open decisions

None.
