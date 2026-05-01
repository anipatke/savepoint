---
id: E05-project-cleanup/T005-rewrite-agents-md
status: planned
objective: "Rewrite AGENTS.md as a lean agent guide"
depends_on: [E04-board-phase-integration/T008-board-tests]
---

# T005: Rewrite AGENTS.md

## Acceptance Criteria

- `AGENTS.md` is under 100 lines.
- Covers: router read order, phase workflow, build commands, code style.
- No Codebase Map table.
- No audit workflow description.
- No 6-state planning machine.
- No capability check warning.
- Valid markdown, readable by any agent.

## Implementation Plan

- [ ] Read existing `AGENTS.md`.
- [ ] Write new lean version: router -> epic -> task workflow.
- [ ] Document phase model: build -> test -> audit -> done.
- [ ] List build/test/lint commands.
- [ ] Keep code style rules (the 10 rules).
- [ ] Remove everything else.

