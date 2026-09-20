---
id: E05-project-cleanup/T003-delete-assets
status: planned
objective: "Delete templates/, agent-skills/, CLAUDE.md, GEMINI.md, ink-cli-ui-design.zip, vitest.config.ts"
depends_on: []
---

# T003: Delete Assets

## Acceptance Criteria

- `templates/` directory does not exist.
- `agent-skills/` directory does not exist.
- `CLAUDE.md`, `GEMINI.md`, `ink-cli-ui-design.zip` do not exist.
- `vitest.config.ts` does not exist (keep `vitest.config.js`).
- No references to deleted files remain.

## Implementation Plan

- [ ] Search codebase for references to templates/, agent-skills/, CLAUDE.md, GEMINI.md.
- [ ] Delete `templates/` and all files within.
- [ ] Delete `agent-skills/` and all files within.
- [ ] Delete `CLAUDE.md`, `GEMINI.md`, `ink-cli-ui-design.zip`.
- [ ] Delete `vitest.config.ts`.
- [ ] Run `npm run build`.
