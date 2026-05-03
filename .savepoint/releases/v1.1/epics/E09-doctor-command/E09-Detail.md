---
type: epic-design
status: audited
---

# E09: Doctor Command

## Purpose

Implement `savepoint doctor`, the diagnostic command for corrupted or inconsistent Savepoint projects. Doctor reports problems and suggested repairs, but does not make destructive changes automatically.

## Interface

```bash
savepoint doctor                    # run all checks
savepoint doctor --epic E03          # check specific epic
```

## What this epic adds

- Integrity checks: config, router state, release/epic/task structure
- YAML/frontmatter validation: detect corrupt YAML
- Task acceptance criteria check: verify AC present
- Dependency checks: missing deps, cycles, duplicate IDs
- Audit state check: proposals without audit-pending flag
- Orphan detection: tasks in nonexistent epics
- Ad-hoc quality gate runner
- Human-readable diagnostic output
- Repair suggestions for common failure modes
- Exit codes: clean (0), diagnosed problems (1), internal failure (2)

## Components

| Module | Purpose |
|--------|---------|
| `cmd/doctor.go` | CLI registration, arg parsing |
| `internal/doctor/checks.go` | Individual integrity checks |
| `internal/doctor/report.go` | Diagnostic report formatting |
| `internal/doctor/repairs.go` | Suggested repair text |
| `internal/doctor/gates.go` | Ad-hoc quality gate runner |

## Implemented As

- `cmd/doctor.go` parses `doctor [--epic <epic>]`, reports help, rejects unsupported arguments, and delegates execution through an injected runner.
- `main.go` wires `savepoint doctor` to `internal/doctor.RunAllChecks` and preserves the required 0/1/2 exit-code contract.
- `internal/doctor/checks.go` implements config, router, structure, dependency, duplicate ID, audit-state, and orphan diagnostics.
- `internal/doctor/gates.go`, `report.go`, and `repairs.go` run configured quality gates, format human-readable reports, and attach repair suggestions.
- Tests live in `cmd/doctor_test.go` and `internal/doctor/*_test.go`.

## Boundaries

**In scope:**
- Detect corrupt YAML/frontmatter
- Detect missing config
- Detect missing dependencies and dependency cycles
- Detect duplicate task IDs
- Detect tasks in nonexistent epics
- Detect audit proposals without matching audit-pending state
- Run configured quality gates on demand

**Out of scope:**
- Auto-moving orphaned tasks without user action
- Repairing files destructively
- Launching the TUI
- Calling AI APIs
