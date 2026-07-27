---
id: E36-visual-identity-template-polish/T001-restructure-visual-identity-template
status: done
objective: Rewrite templates/project/.savepoint/visual-identity.md so general design-system sections are TUI-neutral and the Atari-Noir theme is presented as an illustrative example, with terminal/TUI content demoted to an optional appendix.
depends_on: []
complexity_tier: medium
complexity_reason: Restructures one template with no code changes, but rewrites every general design-system section.
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

- [x] Write new top-level structure: generic overview → palette template → typography → spacing → visual patterns → interaction principles → replication brief → flex & constraints → TUI appendix → savepoint init note.
- [x] Populate each section with project-agnostic guidance; use Atari-Noir values in example table cells with `e.g.` or `Example:` labels.
- [x] Move TUI-specific content into the appendix, generalizing implementation-specific language where needed.
- [x] Verify no TUI-specific framing leaks into the general sections.

## Context Log

Implementation present in the working tree, verified 2026-07-26. Four of six AC met; one partially met.

Met:
- Structure is generic-first with the appendix last: Palette (13), Typography (34), Spacing & Layout Rhythm (44), Visual Patterns (50), Interaction Principles (57), Replication Brief (66), Flex & Constraints (82), then `## Appendix: TUI / Terminal Adaptation` (92).
- The appendix carries the demoted terminal content: "What survives in the terminal" table (96) including the scanlines and glow rows, and "Terminal UI guardrails" (111).
- The `## When \`savepoint init\` ships` section is preserved (121).
- Frontmatter `type: visual-identity` and the filename are unchanged.
- Atari-Noir is labelled as illustration — "Hex (example)", "Role (example)", "*Example (Atari-Noir theme): ...*" — and the header states the shipped values are an example theme to replace.
- Gate: `make build` exit 0, `make test` exit 0 (all packages ok).

Previously partially met — AC "general sections are TUI-neutral" had two leaks above the appendix boundary. Both fixed 2026-07-26:
- Line 17, Palette table: "Page / terminal background" → "Base background".
- Line 80, Replication Brief: "uniform terminal black backgrounds" → "uniform black backgrounds".

Re-verified: a scan of lines 1–91 for terminal/TUI vocabulary returns no TUI-specific framing. Remaining matches are incidental — "underglow" (line 63), "intuitive" (line 76), and the monospace note at line 39, which is medium-neutral typography advice. Gate re-run after the fix: `make build` exit 0, `make test` exit 0 (all packages ok).

All six acceptance criteria met. Awaiting user sign-off for `status: done`.

Scope extension, authorised by the user 2026-07-26: pre-existing project-specific audience copy (present in `HEAD`, not introduced here) was removed rather than deferred to a follow-up task.

- Line 76: "Young Explorer Baseline: ... understandable by a 7-year-old ... discoverable through play" → "Audience baseline: ... understandable by your least-expert intended user ... discoverable through use".
- Line 88: "physical/visual analogies a 7-year-old can grasp" → "concrete analogies your audience already understands".

Both retain the original principle while dropping the fixed audience. Remaining `savepoint` references in the file are confined to the `## When savepoint init ships` section, which an acceptance criterion requires preserved. Gate after the change: `make build` exit 0, `make test` exit 0.
