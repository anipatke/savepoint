---
type: epic-design
status: active
---

# Epic E04: templates-and-prompts

## Purpose

Create the template and prompt assets that `savepoint init` will copy into new projects. This epic makes the Router Pattern concrete as versioned markdown/config assets instead of hardcoded strings.

## What this epic adds

- Default `AGENTS.md`, `.savepoint/router.md`, project PRD, project Design, release PRD, config, and visual identity templates.
- Prompt templates for PRD creation, Design creation, epic design, task breakdown, task planning, task building, and audit reconciliation.
- Embedded `<!-- AGENT: ... -->` instructions in prompt/template markdown where agent behavior matters.
- Template loading helpers that future commands can call.
- Tests proving template files exist, render, and contain required markers.

## Components and files

Expected files introduced or extended by this epic:

| Path                                              | Purpose                                |
| ------------------------------------------------- | -------------------------------------- |
| `templates/project/AGENTS.md`                     | Root agent entrypoint template.        |
| `templates/project/.savepoint/router.md`          | Default router template.               |
| `templates/project/.savepoint/PRD.md`             | Project PRD starter.                   |
| `templates/project/.savepoint/Design.md`          | Project Design starter.                |
| `templates/project/.savepoint/config.yml`         | Default config.                        |
| `templates/project/.savepoint/visual-identity.md` | Default visual identity.               |
| `templates/release/v1/PRD.md`                     | Release PRD starter.                   |
| `templates/prompts/*.prompt.md`                   | Agent-facing workflow prompts.         |
| `src/templates/*.ts`                              | Template lookup and rendering helpers. |
| `test/templates/**/*.test.ts`                     | Template integrity tests.              |

## Architectural delta

Before this epic, project scaffolding would need to invent file contents in code. After this epic, generated Savepoint projects come from data files that agents and humans can inspect.

Templates become a stable internal asset boundary consumed by `init-command`, audit prompts, and future workflow docs.

## Boundaries

In scope:

- Template file structure and required content.
- Simple variable interpolation where needed, such as project name and release number.
- Integrity tests for required files and markers.

Out of scope:

- Writing templates into a target directory.
- Clipboard integration.
- TUI consumption of templates beyond later audit-review display needs.
- AI semantic-review content deferred to v0.2.0.

## Quality gates

- Template tests should fail if required files, frontmatter, router states, or agent instruction markers are missing.
- Keep prompts concise enough to preserve the token-efficiency goals.

## Design constraints

- Content lives in markdown/YAML files, not large TypeScript string literals.
- Prompt text should be agent-agnostic.
- Avoid references to implementation files that do not exist in a newly initialized project.
