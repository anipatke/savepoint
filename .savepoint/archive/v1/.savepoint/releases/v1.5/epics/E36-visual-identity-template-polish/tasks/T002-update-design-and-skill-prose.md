---
id: E36-visual-identity-template-polish/T002-update-design-and-skill-prose
status: done
objective: Update templates/project/.savepoint/Design.md and agent-skills/bubbletea-tui-design/SKILL.md so they no longer frame visual-identity.md as exclusive to TUI work.
depends_on: ["E36-visual-identity-template-polish/T001-restructure-visual-identity-template"]
complexity_tier: low
complexity_reason: Small prose edits across two files with no code or behaviour changes.
---

# T002: Update Design.md and Skill Prose

## Problem

After T001 makes the visual-identity template project-agnostic, the Design.md template and the bubbletea-tui-design skill still describe it as "loaded only for TUI/theme/visual tasks." This contradiction confuses non-TUI projects and undermines the agnostic framing.

## Context Files

- `templates/project/.savepoint/Design.md`
- `agent-skills/bubbletea-tui-design/SKILL.md`

## Acceptance Criteria

- [ ] `templates/project/.savepoint/Design.md` section header callout (line 11) no longer says "loaded only for TUI/theme/visual tasks."
- [ ] `templates/project/.savepoint/Design.md` directory layout table (line 34) no longer says "loaded conditionally for TUI work."
- [ ] `templates/project/.savepoint/Design.md` TUI section (line 129) references the appendix-adapted VI file appropriately.
- [ ] `agent-skills/bubbletea-tui-design/SKILL.md` lines 8 and 23 frame visual-identity as a general design-system reference, not TUI-only.
- [ ] `make build && make test` passes.

## Implementation Plan

- [x] Edit `templates/project/.savepoint/Design.md`:
  - Line 11: Change to `> **Visual identity** lives separately in `.savepoint/visual-identity.md` and covers your project's design system.`
  - Line 34: Change to `├── visual-identity.md          ← design system (palette, type, patterns)`
  - Line 129: Adjust to reference the general sections and note the TUI appendix is for terminal projects.
- [x] Edit `agent-skills/bubbletea-tui-design/SKILL.md`:
  - Line 8: Generalize the description of visual-identity.md.
  - Line 23: Change "Read `.savepoint/visual-identity.md` only when the task touches rendering..." to "Read `.savepoint/visual-identity.md` for design-system context; focus on the TUI appendix when the task touches rendering..."
- [x] Verify consistency between the two files.

## Context Log

Implementation present in the working tree, verified 2026-07-26 against all five AC.

- `templates/project/.savepoint/Design.md:11` now reads "defines the project's design system (palette, typography, patterns)" — the "loaded only for TUI/theme/visual tasks" framing is gone.
- `templates/project/.savepoint/Design.md:34` directory row now reads "design system (palette, type, patterns)" without the conditional-loading note.
- `templates/project/.savepoint/Design.md:129` now references the file generally and points terminal projects at the TUI adaptation appendix.
- `agent-skills/bubbletea-tui-design/SKILL.md:8` reframes visual-identity as the project's design system with Atari-Noir as the example theme.
- `agent-skills/bubbletea-tui-design/SKILL.md:23` now reads "for design-system context; focus on the TUI adaptation appendix when...".
- Gate: `make build` exit 0, `make test` exit 0 (all packages ok).

All acceptance criteria met. Awaiting user sign-off for `status: done`.
