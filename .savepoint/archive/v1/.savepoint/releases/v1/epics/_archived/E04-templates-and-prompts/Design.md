---
type: epic-design
status: audited
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

## Implemented As

- Project scaffold assets were added under `templates/project/`, including `AGENTS.md`, `.savepoint/router.md`, `.savepoint/PRD.md`, `.savepoint/Design.md`, `.savepoint/config.yml`, and `.savepoint/visual-identity.md`.
- Release starter content was added at `templates/release/v1/PRD.md`.
- Prompt assets were added under `templates/prompts/` for PRD creation, project design, epic design, task breakdown, task planning, task building, and audit reconciliation.
- `src/templates/manifest.ts` defines typed template sets for project, release, and prompt assets.
- `src/templates/paths.ts` resolves template roots and family-specific paths.
- `src/templates/load.ts` loads named templates from disk and returns path-aware `not_found` errors.
- `src/templates/render.ts` interpolates `PROJECT_NAME`, `RELEASE_NUMBER`, and `RELEASE_NAME` placeholders while leaving unresolved placeholders intact as a visible failure signal.
- `src/templates/index.ts` exposes a narrow public helper surface.
- Template tests were added under `test/templates/` for required files, frontmatter, router states, agent instruction markers, audit prompt shape, path lookup, loading failures, and interpolation behavior.
- The audit snapshot was generated manually because `savepoint audit` is still a stub.

Design delta notes:

- The loader reports all read failures as `not_found`; richer filesystem error typing is deferred until command handlers need to distinguish missing files from permission or IO failures.
- The audit prompt correctly requires one proposal bundle, but audit startup exposed a workflow gap: the router can point agents at a snapshot before one exists. The follow-up improvement is to make snapshot creation an explicit precondition before entering `audit-pending`, or to instruct the agent to create one manual snapshot once without searching broadly.

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
