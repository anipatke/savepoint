---
id: E04-templates-and-prompts/T001-project-template-assets
status: done
objective: "Create the default project-level Savepoint template assets copied by future init workflows."
depends_on: []
---

# T001: project-template-assets

## Implementation Plan

- [x] Create `templates/project/AGENTS.md` with root workflow instructions, build/test placeholders, code style rules, and the generated Codebase Map markers.
- [x] Create `templates/project/.savepoint/router.md` with the initial state machine, router states, conditional read guidance, and embedded `<!-- AGENT: ... -->` instructions.
- [x] Create starter `templates/project/.savepoint/PRD.md` and `templates/project/.savepoint/Design.md` files with frontmatter, headings, and placeholders safe for a new project.
- [x] Create `templates/project/.savepoint/config.yml` with default Savepoint configuration values.
- [x] Create `templates/project/.savepoint/visual-identity.md` as a concise optional visual design starter.
- [x] Run a focused format check over the new template markdown/YAML files.
