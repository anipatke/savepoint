---
id: E09-doctor-command/T006-quality-gates-report
status: done
objective: "Run quality gates on demand and format diagnostic report"
depends_on: ["E09-doctor-command/T005-audit-orphan-checks"]
---

# T006: Quality Gates + Report

## Acceptance Criteria

- Run configured lint, typecheck, test commands from config
- Report pass/fail status for each gate
- Format all diagnostics into human-readable report
- Include file paths and actionable fix suggestions
- Return appropriate exit code

## Implementation Plan

- [x] Add `internal/doctor/gates.go`
- [x] Implement `RunQualityGates(root) []GateResult`
- [x] Read config quality_gates
- [x] Execute each command, capture output
- [x] Report pass/fail for each gate
- [x] Add `internal/doctor/report.go`
- [x] Collect all check results
- [x] Format into human-readable report
- [x] Add repair suggestions from internal/doctor/repairs.go
- [x] Return exit code: 0 if all clean, 1 if problems found
- [x] Test doctor on valid and invalid projects
- [x] Run `make build && make test`