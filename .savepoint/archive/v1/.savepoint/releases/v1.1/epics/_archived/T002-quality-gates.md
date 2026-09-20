---
id: E06-audit-command/T002-quality-gates
status: done
objective: "Implement quality gate runner that blocks on failure"
depends_on: ["E06-audit-command/T001-cli-entrypoint"]
phase: build
---

# T002: Quality Gates

## Acceptance Criteria

- Reads `quality_gates` from `.savepoint/config.yml`
- Runs each configured command (lint, typecheck, test)
- Blocks audit if any gate fails and `block_on_failure: true`
- Continues if `block_on_failure: false` but logs failure
- Shows clear output: which gates passed/failed

## Implementation Plan

- [x] Add `internal/audit/gates.go`
- [x] Read config `quality_gates` struct
- [x] Implement `RunGates(ctx) error` that executes each command
- [x] Handle null gates (skip if not configured)
- [x] Add exit code propagation for blocking behavior
- [x] Test with configured gates in .savepoint/config.yml
- [x] Run `make build && make test`