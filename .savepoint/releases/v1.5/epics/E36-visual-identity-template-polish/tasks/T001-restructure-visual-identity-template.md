---
id: E36-visual-identity-template-polish/T001-restructure-visual-identity-template
status: planned
objective: Rewrite templates/project/.savepoint/visual-identity.md so general design-system sections are TUI-neutral and the Atari-Noir theme is presented as an illustrative example, with terminal/TUI content demoted to an optional appendix.
depends_on: []
complexity_tier: medium
---

# T001: Restructure visual-identity.md Template

## Problem

The shipped visual-identity template is written entirely from a TUI perspective — palette, typography, and patterns assume a terminal UI project. Projects of other types (web, API, desktop) receive a template that is confusing and requires heavy editing.

## Context Files

- `templates/project/.savepoint/visual-identity.md`
- `.savepoint/visual-identity.md` (reference for the repo's own copy — not to be changed)

## Acceptance Criteria

- [ ] General design-system sections (palette, typography, spacing, visual patterns, interaction principles, replication brief, flex & constraints) are written in TUI-neutral language.
- [ ] The Atari-Noir theme (colors, fonts, patterns) is present as a filled-in illustrative example — not as required default content. Template instructions make it clear users replace these with their own design system.
- [ ] The terminal feasibility table ("What survives in the terminal"), "Terminal UI guardrails," scanlines/glow material, and the "What Survives in the Terminal" table live in a clearly delimited, removable appendix titled `## Appendix: TUI / Terminal Adaptation`.
- [ ] The section explaining the file's role at `savepoint init` time is preserved.
- [ ] The filename `visual-identity.md` and frontmatter `type: visual-identity` remain unchanged.
- [ ] `make build && make test` passes.

## Implementation Plan

- [ ] Write new top-level structure: generic overview → palette template → typography → spacing → visual patterns → interaction principles → replication brief → flex & constraints → TUI appendix → savepoint init note.
- [ ] Populate each section with project-agnostic guidance; use Atari-Noir values in example table cells with `e.g.` or `Example:` labels.
- [ ] Copy TUI-specific content into the appendix without rewording.
- [ ] Verify no TUI-specific framing leaks into the general sections.

## Context Log

Pending.
