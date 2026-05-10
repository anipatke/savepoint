---
type: epic-design
status: audited
---

# E16: Pre-Production Refinement

## Purpose

Refine the pre-production workflow before release by tightening init safety and reducing unnecessary agent context reads. This epic focuses on small, high-leverage changes that protect user files and better enforce Savepoint's token-efficiency goals.

## What this epic adds

- Existing `AGENTS.md` files no longer block init validation.
- Savepoint instructions are written as a managed block inside existing agent guide files.
- Re-running init refreshes the managed block without duplicating it.
- User-owned content outside the managed block is preserved.
- Common filename casing variants such as `Agents.MD` are detected and reused.
- Savepoint phase skills define explicit read budgets and phase-specific allowed reads.
- Scaffolded skill templates carry the same context-discipline rules as bundled skills.
- "Create epic/task only" workflows explicitly prohibit source-code reads, broad searches, tests, and status checks.
- The CLI exposes a simple version-reporting entry point such as `savepoint --version`.

## Components

| Module | Purpose |
|--------|---------|
| `internal/init/validate.go` | Allow existing agent guide files while keeping real scaffold conflicts protected |
| `internal/init/scaffold.go` | Route `AGENTS.md` template writes through merge behavior |
| `internal/init/agents.go` | Add focused helper for detecting, inserting, and refreshing managed agent-guide blocks |
| `internal/init/validate_test.go` | Update validation expectations for existing agent guides |
| `internal/init/scaffold_test.go` | Cover insert, idempotency, casing, and force refresh behavior |
| `internal/init/integration_test.go` | Verify init succeeds in projects with an existing agent guide |
| `agent-skills/` | Update live phase skill guides with read-budget and tool-discipline rules |
| `templates/project/agent-skills/` | Update scaffolded phase skill templates to match live skills |
| `templates/project/AGENTS.md` | Add concise context-budget rules for generated projects |
| `main.go` | Expose and test the version-reporting CLI surface |
| `cmd/` | Keep command parsing behavior consistent with the version surface if parsing changes are needed |

## Implemented as

- `internal/init/agents.go` owns agent-guide casing detection plus managed block insertion/replacement.
- `internal/init/scaffold.go` routes `AGENTS.md` scaffold writes through the managed merge path and now scaffolds the release skeleton from template assets.
- `cmd/upgrade-assets.go` adds the CLI parser for `upgrade-assets [dir] [--dry-run] [--force]`.
- `internal/init/upgrade.go` owns existing-project validation, package-owned asset allowlisting, dry-run reporting, idempotent skill refresh, and managed agent-guide block refresh.
- `templates/project/.savepoint/releases/v1/v1-PRD.md` seeds the release PRD referenced by the scaffolded router.
- `package.json` postinstall prints a notice only; project mutation remains explicit via `savepoint upgrade-assets`.

## Boundaries

**In scope:**
- `savepoint init` behavior for existing agent guide files
- Managed block insertion and replacement
- Case-variant detection for root agent guide filenames
- Unit and integration tests for init behavior
- Skill and skill-template guidance that enforces minimal file reads by phase
- Context-budget guidance for planning-only and create-only workflows
- A narrow version-reporting CLI surface such as `savepoint --version`, `savepoint version`, or equivalent

**Out of scope:**
- New CLI flags or prompts other than the version-reporting surface
- Changing the magic prompt contract unless required by casing behavior
- Modifying generated Savepoint instruction content beyond managed markers
- Changing `agent-skills` conflict behavior
- Large process redesigns beyond pre-production guidance refinements
