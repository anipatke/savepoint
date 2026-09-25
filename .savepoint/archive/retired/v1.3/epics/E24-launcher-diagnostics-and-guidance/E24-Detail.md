---
type: epic-design
status: planned
---

# E24: Launcher Diagnostics and Guidance

## Purpose

Make the optional launcher understandable and supportable through doctor diagnostics, scaffolded configuration, workflow guidance, and release-level regression verification.

## What this epic adds

- Doctor validation for enabled launcher profiles, placeholders, and terminal configuration.
- Disabled-by-default launcher examples in newly scaffolded projects.
- User documentation for configuring builder/auditor roles and supported terminal behavior.
- Agent guidance explaining bounded item launches and fresh-session epic audits.
- End-to-end tests proving disabled compatibility and enabled action dispatch.

## Components and files

| Module | Purpose |
|--------|---------|
| `internal/doctor` | Diagnose invalid enabled launcher configuration |
| `templates/project/.savepoint/config.yml` | Scaffold the opt-in configuration surface |
| `templates/project/AGENTS.md` | Explain launch entrypoints without replacing phase skills |
| `README.md` | Document setup, limitations, and platform behavior |
| `internal/board/integration_test.go` | Verify launcher behavior through board boundaries |

## Architectural delta

No new runtime boundary is added. This epic exposes and validates the contracts introduced by E22 and E23 while preserving configuration ownership and upgrade-assets rules.

## Boundaries

**In scope:**
- Read-only diagnostics
- New-project templates and user documentation
- Workflow guidance and regression tests
- Manual cross-platform verification checklist

**Out of scope:**
- Auto-editing existing user config during upgrade
- Installing agent CLIs or terminal applications
- Authentication diagnostics
- End-to-end tests that consume paid agent API calls

## Quality gates

- `savepoint doctor` remains read-only and ignores absent/disabled launcher config.
- Template tests prove newly initialized projects remain disabled by default.
- Full `make build && make test` passes.

## Open decisions

None.
