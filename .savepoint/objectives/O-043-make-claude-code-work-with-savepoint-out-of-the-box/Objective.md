---
id: O-043
title: Make Claude Code work with Savepoint out of the box
status: planned
depends_on: [O-039]
release: G-002
priority: medium
rank: 4
---

# O-043: Make Claude Code work with Savepoint out of the box

## Outcome

In a project set up or upgraded by Savepoint, a plain Claude Code session knows the Savepoint workflow from its first message, opens the workflow skills on demand, starts already knowing the Next line, is stopped by the harness before it breaks a hard owner-authority rule, and shows the current Savepoint state in its status line. Session-start context cost is no higher than today and falls where skills load on demand.

## Why

`savepoint init` writes `AGENTS.md` and `agent-skills/`, but Claude Code reads neither by itself. It loads only `CLAUDE.md` and skills under `.claude/skills/`. Unless the owner pastes the starter prompt, Claude starts with no Savepoint context, and when it does follow the prompt it spends a tool call reading the guide. Owner-authority rules (only the owner marks a Task done, no router edits in a lane, owner-only CLI commands) are enforced only by prose. Found in a review of the repository on 2026-10-10.

## Success Conditions

- A fresh `savepoint init` project opened in Claude Code mentions Savepoint and its Next line when asked "what should I do next?", with no pasted prompt.
- The four workflow skills appear as Claude Code skills; only their names and descriptions load until one is used.
- Session start adds no more than about 60 tokens beyond the guide (the Next line), and does nothing visible when `savepoint` is unavailable.
- An agent edit that sets a Task to `status: done`, an edit to `.savepoint/router.md` inside a worktree lane, and an agent run of an owner-only `savepoint` command are each blocked with a one-line reason.
- The status line shows the selected record and Next action, and flags owner decisions, at zero context cost.
- `savepoint upgrade-assets` brings existing projects the same files, and never overwrites a user's own `CLAUDE.md` content, skills, or `.claude/settings.json`.
- Session-start token weight before and after is recorded for this repository and a fresh project.

## Architectural Considerations

- Delivery is through `savepoint init` and `savepoint upgrade-assets` (TPL-04), with new templates under `templates/project-v2/`. The embed in `main.go` needs `all:` for the `.claude/` directory.
- `CLAUDE.md` follows the `AGENTS.md` pattern: a Savepoint-managed region (`<!-- SAVEPOINT:BEGIN/END -->`) holding the `@AGENTS.md` import, merged into an existing file without touching text outside the region (FS-01, FS-02).
- Skill entries under `.claude/skills/<skill>/SKILL.md` are thin pointers: the same `name` and `description` as the canonical skill, with a body that says to read and follow `agent-skills/<skill>/SKILL.md`. The canonical skill stays the single home (TPL-01). Pointers are manifest-tracked like skills, so an edited pointer is kept and reported.
- Hooks and the status line are Node scripts under `.claude/hooks/`. Node is already present because Savepoint installs through `npx`. Scripts find `savepoint` on PATH, then in `node_modules/.bin`, and exit silently with no output when it is missing. They never block a session for a missing tool, and they run on Windows (CFG-02, CFG-03).
- `.claude/settings.json` is written only when absent. When it exists, Savepoint never rewrites it. Upgrade reports the exact entries to add (hooks, status line) and leaves the file unchanged (FS-01). A user's own `statusLine` is never replaced.
- Hooks are a Claude Code convenience only. The prose rules in `AGENTS.md` and the skills stay, because Codex and Gemini read them.

## Confirmed Design

Owner-confirmed on 2026-10-10:

- **Goal.** Lives in G-002. G-002's boundaries are widened to allow the scaffold and upgrade changes this Objective needs.
- **Audience.** Ships to every Savepoint project through `init` and `upgrade-assets`, not as a repository-only trial.
- **Scope.** Five Tasks: `CLAUDE.md` import, skills as Claude Code skills, session-start Next line, guard hooks, status line. Extra permission allowlists, a post-compaction reminder, and `/next` or `/check` commands are deferred. The owner is undecided on them.

## Boundaries

**In scope:**
- `templates/project-v2/` additions for `CLAUDE.md`, `.claude/skills/`, `.claude/hooks/`, and `.claude/settings.json`, and the init and upgrade code that installs them.
- This repository's own `CLAUDE.md` and `.claude/` files, regenerated from the templates.
- Content and scaffold tests for the new files.

**Out of scope:**
- Permission allowlists, a post-compaction hook, and slash commands (deferred).
- A Claude Code plugin or marketplace package.
- Removing prose rules from `AGENTS.md` or the skills.
- Integrations for Codex, Gemini, or other agents.
- Any change to lifecycle rules, states, or the `resume` output format.
