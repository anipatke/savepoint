---
id: E09-doctor-command/T002-config-router-validation
status: done
objective: "Validate config.yml and router state"
depends_on: ["E09-doctor-command/T001-cli-entrypoint"]
---

# T002: Config + Router Validation

## Acceptance Criteria

- Config.yml exists and is valid YAML
- Config has required fields: quality_gates, theme
- Router state is valid (valid state name, release, epic)
- Router state matches actual release/epic directories

## Implementation Plan

- [x] Add `internal/doctor/checks.go`
- [x] Implement `CheckConfig(root) error`
- [x] Validate config.yml exists and parses
- [x] Validate required fields present
- [x] Implement `CheckRouter(root, epicFilter) error`
- [x] Validate router state is valid YAML
- [x] Validate state name is one of: pre-implementation, epic-design, epic-task-breakdown, task-building, audit-pending
- [x] Validate release/epic directories exist
- [x] Test validation on valid and invalid projects
- [x] Run `make build && make test`