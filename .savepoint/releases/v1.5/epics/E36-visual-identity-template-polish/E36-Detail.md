---
type: epic-design
status: planned
---

# E36: Visual-identity Template Polish

## Purpose

Make the shipped `visual-identity.md` template project-agnostic so any project type (web, CLI, desktop, API) gets a useful design-system starter. Atari-Noir stays as a filled-in illustrative example. Terminal/TUI-specific content is demoted to a clearly delimited optional appendix.

## What this epic adds

- A restructured `templates/project/.savepoint/visual-identity.md` with TUI-neutral general sections first and an optional TUI adaptation appendix at the bottom.
- Updated prose in `templates/project/.savepoint/Design.md` and `agent-skills/bubbletea-tui-design/SKILL.md` so visual-identity is no longer framed as "loaded only for TUI tasks."

## Components and files

| Module | Purpose |
|--------|---------|
| `templates/project/.savepoint/visual-identity.md` | Main template to restructure — keep filename, change body |
| `templates/project/.savepoint/Design.md` | Update prose (lines 11, 34, 129) that frame VI as TUI-only |
| `agent-skills/bubbletea-tui-design/SKILL.md` | Update prose (lines 8, 23) that frame VI as TUI-only |

## Architectural delta

No Go code changes. The template remains opaque to the binary: embedded at `main.go:18`, scaffolded verbatim at `internal/init/scaffold.go:32-51`, and skipped by `upgrade-assets` at `internal/init/upgrade.go:118-121`. The file's existence (not its contents) is what the tool cares about.

## Boundaries

**In scope:**

- Restructuring the visual-identity template body
- Updating Design.md template and bubbletea-tui-design skill prose to match the agnostic framing
- Keeping the filename `visual-identity.md` and `type: visual-identity` frontmatter intact

**Out of scope:**

- The repo's own `.savepoint/visual-identity.md` (not a template, not shipped to users)
- Go code changes of any kind
- Guardrails.md or Health-Check.md (separate epics E37, E38)

## Quality gates

- `make build && make test` passes (no Go tests assert VI content).
- Fresh scaffold produces a `visual-identity.md` whose general sections contain no TUI-specific framing.
- Atari-Noir is presented as an illustrative example, not a mandatory default.
- Terminal-specific content (feasibility table, terminal UI guardrails, scanlines/glow material, What Survives table) lives in a clearly delimited optional appendix.

## Open decisions

None. TUI adaptation appendix remains in the same file so it stays visible to TUI projects while being removable by non-TUI projects.
