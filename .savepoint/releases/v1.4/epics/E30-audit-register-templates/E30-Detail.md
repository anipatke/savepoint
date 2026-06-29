---
type: epic-design
status: audited
---

# E30: Audit Register Templates

## Purpose

Add the scaffolded markdown assets that define the Audit Register's prompt, current register, finding records, and run history without changing runtime behavior.

## What this epic adds

- A canonical audit prompt template with a changelog and required output fields.
- A repo-wide register template that summarizes current finding state.
- Folder guidance for durable finding records and immutable run records.
- Init and upgrade-assets handling for audit-register documentation assets.
- Template freshness and scaffold coverage for the new files.

## Components and files

| Module | Purpose |
|--------|---------|
| `templates/project/.savepoint/audit/` | Scaffold the audit prompt, register, findings guidance, and run history guidance |
| `internal/init` | Copy and refresh audit-register documentation assets safely |
| `templates/project/AGENTS.md` | Introduce the audit-register workflow entry point for generated projects |

## Architectural delta

Savepoint gains a new optional project documentation area under `.savepoint/audit/`. These files are durable project truth when present, but existing projects without the directory remain valid.

## Boundaries

**In scope:**
- Markdown templates and placeholder guidance
- Safe creation during init
- Safe refresh during upgrade-assets
- Template freshness tests

**Out of scope:**
- Audit data parsing
- Board rendering
- Doctor diagnostics
- Actual register population for this repository

## Quality gates

- Init tests prove new projects contain the audit-register assets.
- Upgrade tests prove user-edited audit files are not overwritten.
- Template freshness tests include the new audit assets.
- `go test ./internal/init` passes.

## Open decisions

None.
