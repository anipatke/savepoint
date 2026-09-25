---
type: epic-design
status: audited
---

# E21: Document Template Optimisation

## Purpose

Audit, clean up, and optimise the four core project document templates — `AGENTS.md`, `PRD.md`, `Design.md`, and a new `Concept.md` — so scaffolded projects start with tighter, more actionable guidance and a richer pre-PRD ideation surface.

## What this epic adds

- A revised `templates/project/AGENTS.md` that is leaner, removes any guidance drift introduced since E18, and aligns with the E20 lifecycle contract.
- A revised `templates/project/.savepoint/PRD.md` with sharper section prompts that help vibe coders write concrete, scoped product visions faster.
- A revised `templates/project/.savepoint/Design.md` with section headings and prompts that match the current directory layout and architecture model.
- A new `templates/project/.savepoint/Concept.md` template that gives projects a lightweight pre-PRD ideation document for capturing raw ideas, target feelings, and open questions before committing to a product direction.
- Template freshness tests updated to cover all four documents.

## Components and files

| File | Purpose |
|------|---------|
| `templates/project/AGENTS.md` | Scaffolded agent entry point |
| `AGENTS.md` | Live agent entry point (kept in sync where guidance applies equally) |
| `templates/project/.savepoint/PRD.md` | Scaffolded product vision template |
| `templates/project/.savepoint/Design.md` | Scaffolded architecture template |
| `templates/project/.savepoint/Concept.md` | New pre-PRD ideation template (scaffolded only) |
| `internal/init/template_freshness_test.go` | Freshness / alignment guards |

## Architectural delta

Before this epic, `AGENTS.md` and the three `.savepoint/` document templates have accrued small inconsistencies since E18's template optimisation pass. `Design.md` section numbering does not yet reflect the current directory layout. There is no scaffolded ideation document between raw idea and formal PRD. After this epic, the four document templates form a coherent pre-build funnel: Concept → PRD → Design → AGENTS.

## Boundaries

**In scope:**
- `templates/project/AGENTS.md` and live `AGENTS.md` alignment (guidance only — no skill changes).
- `templates/project/.savepoint/PRD.md` and `Design.md` section/prompt revisions.
- Authoring `templates/project/.savepoint/Concept.md` from scratch.
- Freshness test coverage for all four templates.

**Out of scope:**
- Changes to agent skills (`agent-skills/` or `templates/project/agent-skills/`).
- Changes to `router.md`, `config.yml`, or `visual-identity.md` templates.
- Live project `PRD.md` and `Design.md` (project-state files, not templates).
- Adding `Concept.md` to `upgrade-assets` for existing projects (deferred because `.savepoint/` project-state files are intentionally skipped).

## Open decisions

- **Existing-project upgrade inclusion:** `savepoint init` now scaffolds `Concept.md` automatically because it walks embedded `templates/project`. Existing Savepoint projects do not receive `.savepoint/Concept.md` through `upgrade-assets`; decide after real use whether that command should create the file or continue treating it as project state.

## Quality gates

- `make build && make test` passes.
- Template freshness tests cover `AGENTS.md`, `PRD.md`, `Design.md`, and `Concept.md`.
- Live and scaffolded `AGENTS.md` are consistent on lifecycle guidance.
