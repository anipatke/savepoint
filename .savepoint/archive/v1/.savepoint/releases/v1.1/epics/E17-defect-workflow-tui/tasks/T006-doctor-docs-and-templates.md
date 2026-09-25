---
id: E17-defect-workflow-tui/T006-doctor-docs-and-templates
status: done
objective: Document the defect workflow and add doctor diagnostics for malformed defect files
depends_on: [E17-defect-workflow-tui/T001-defect-data-model, E17-defect-workflow-tui/T002-defect-router-priority]
---

# T006: Doctor, Docs, And Templates

## Problem

Defects become part of the methodology once the TUI can display them. New projects and diagnostics need to explain and enforce the defect file shape so repair work stays evidence-driven instead of becoming ad hoc notes.

## Context Files

- `internal/doctor/checks.go` - project integrity diagnostics
- `internal/doctor/report.go` - diagnostic report formatting
- `internal/doctor/checks_test.go` - doctor diagnostic test patterns
- `templates/project/AGENTS.md` - scaffolded agent methodology guidance
- `templates/project/.savepoint/router.md` - scaffolded router state guidance
- `templates/project/.savepoint/Design.md` - scaffolded architecture documentation
- `README.md` - user-facing workflow documentation
- `AGENTS.md` - live methodology and Codebase Map guidance

## Acceptance Criteria

- [x] Doctor reports malformed defect frontmatter
- [x] Doctor reports invalid defect status and missing stage on in-progress defects
- [x] Doctor reports broken introduced-by or related task references when present
- [x] Scaffolded methodology docs describe when to use defects instead of epics or tasks
- [x] README documents the defect file location and TUI `d` overlay
- [x] Live AGENTS.md Codebase Map is updated only if new modules or responsibilities are added
- [x] Tests cover doctor diagnostics for valid defects and each malformed defect case
- [x] `make build && make test` passes

## Implementation Plan

- [x] Add defect validation checks through doctor using the data-layer parser
- [x] Add report messages and repair suggestions for defect-specific errors
- [x] Update scaffolded methodology documentation with the defect lane
- [x] Update README with concise defect workflow usage
- [x] Update live Codebase Map if the implementation adds new module responsibilities
- [x] Run `make build && make test`

## Context Log

- Files read: checks.go, checks_test.go, interfaces.go, interfaces_test.go, report.go, repairs.go, defect.go, defect_test.go, discover.go, lifecycle.go, parser.go, templates/project/AGENTS.md, templates/project/.savepoint/router.md, README.md, AGENTS.md
- Files edited: interfaces.go, interfaces_test.go, checks.go, checks_test.go, report.go, repairs.go, templates/project/AGENTS.md, README.md
- Token estimate: ~18k
- Quality gates: make build && make test — all pass
