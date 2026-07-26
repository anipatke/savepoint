---
type: audit-findings
audited: 2026-07-26
---

# Audit Findings: E36 Visual-identity Template Polish

## Main Findings

E36 is audited and closed. The task implementation plans now reflect the completed work, the epic is marked `audited`, `.savepoint/Design.md` records E36 as the latest audit, and the router advances to E37's planned T001.

The user explicitly accepted one residual deviation: the original standalone scanline and glow descriptions remain omitted rather than being restored under the TUI appendix. No template change was applied for that finding.

All other reviewed outcomes stand. T001's general sections are TUI-neutral, Atari-Noir is clearly a replaceable example, the terminal table and guardrails are in the appendix, and the `savepoint init` note and required frontmatter remain intact. T002 meets all five acceptance criteria: the Design template and Bubble Tea skill consistently describe visual identity as general design-system context and direct terminal work to the appendix.

Verification completed on 2026-07-26: `make build && make test` passed across all Go packages. No Go source or test changes were in scope.

## Code Style Review

- [x] One job per file — each changed document retains one clear purpose.
- [x] One job per function — not applicable; E36 changes prose only.
- [x] Test branches — not applicable to the prose changes; the full build/test gate passes.
- [x] Types document intent — not applicable; no code or data-model types changed.
- [x] Build only what is needed — changes stay within the three scoped deliverables.
- [x] Handle errors at boundaries — not applicable; no runtime boundary changed.
- [x] One source of truth — the visual-identity template remains the shipped design-system source.
- [x] Comments explain WHY — no implementation comments were added.
- [x] Content in data files — design-system guidance remains in Markdown templates and skills.
- [x] Small diffs — the two supporting prose edits are focused; the larger template rewrite is required by T001.

## Proposed Changes

### Target File
.savepoint/releases/v1.5/epics/E36-visual-identity-template-polish/tasks/T001-restructure-visual-identity-template.md

### Replace
```md
- [ ] Write new top-level structure: generic overview → palette template → typography → spacing → visual patterns → interaction principles → replication brief → flex & constraints → TUI appendix → savepoint init note.
- [ ] Populate each section with project-agnostic guidance; use Atari-Noir values in example table cells with `e.g.` or `Example:` labels.
- [ ] Copy TUI-specific content into the appendix without rewording.
- [ ] Verify no TUI-specific framing leaks into the general sections.
```

### With
```md
- [x] Write new top-level structure: generic overview → palette template → typography → spacing → visual patterns → interaction principles → replication brief → flex & constraints → TUI appendix → savepoint init note.
- [x] Populate each section with project-agnostic guidance; use Atari-Noir values in example table cells with `e.g.` or `Example:` labels.
- [x] Move TUI-specific content into the appendix, generalizing implementation-specific language where needed.
- [x] Verify no TUI-specific framing leaks into the general sections.
```

### Target File
.savepoint/releases/v1.5/epics/E36-visual-identity-template-polish/tasks/T002-update-design-and-skill-prose.md

### Replace
```md
- [ ] Edit `templates/project/.savepoint/Design.md`:
  - Line 11: Change to `> **Visual identity** lives separately in `.savepoint/visual-identity.md` and covers your project's design system.`
  - Line 34: Change to `├── visual-identity.md          ← design system (palette, type, patterns)`
  - Line 129: Adjust to reference the general sections and note the TUI appendix is for terminal projects.
- [ ] Edit `agent-skills/bubbletea-tui-design/SKILL.md`:
  - Line 8: Generalize the description of visual-identity.md.
  - Line 23: Change "Read `.savepoint/visual-identity.md` only when the task touches rendering..." to "Read `.savepoint/visual-identity.md` for design-system context; focus on the TUI appendix when the task touches rendering..."
- [ ] Verify consistency between the two files.
```

### With
```md
- [x] Edit `templates/project/.savepoint/Design.md`:
  - Line 11: Change to `> **Visual identity** lives separately in `.savepoint/visual-identity.md` and covers your project's design system.`
  - Line 34: Change to `├── visual-identity.md          ← design system (palette, type, patterns)`
  - Line 129: Adjust to reference the general sections and note the TUI appendix is for terminal projects.
- [x] Edit `agent-skills/bubbletea-tui-design/SKILL.md`:
  - Line 8: Generalize the description of visual-identity.md.
  - Line 23: Change "Read `.savepoint/visual-identity.md` only when the task touches rendering..." to "Read `.savepoint/visual-identity.md` for design-system context; focus on the TUI appendix when the task touches rendering..."
- [x] Verify consistency between the two files.
```
