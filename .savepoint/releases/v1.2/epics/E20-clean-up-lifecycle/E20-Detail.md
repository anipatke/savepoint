---
type: epic-design
status: audited
---

# E20: Clean-up Lifecycle

## Purpose

Make the task-to-epic workflow contract simple, modular, and hard to drift by centralizing lifecycle rules before adding more one-off compatibility patches.

## What this epic adds

- A single task lifecycle policy surface in `internal/data` for status, stage, legacy aliases, parse compatibility, write validation, and transition behavior.
- Parser and writer behavior that consume the same lifecycle policy instead of duplicating lifecycle decisions.
- Doctor diagnostics that report lifecycle drift without disagreeing with board load behavior.
- Board transitions that use shared lifecycle operations while preserving user-owned task completion boundaries.
- Template, skill, and design updates that document the same lifecycle contract the code enforces.

## Components and files

| Module | Purpose |
|--------|---------|
| `internal/data` | Own the lifecycle contract, parse compatibility, write validation, and transition helpers |
| `internal/doctor` | Report lifecycle drift and repair suggestions using the shared data contract |
| `internal/board` | Use shared lifecycle transitions and keep TUI update code as message handling |
| `AGENTS.md` | Live agent workflow contract |
| `templates/project/AGENTS.md` | Scaffolded agent workflow contract |
| `agent-skills/savepoint-build-task/SKILL.md` | Canonical build workflow guidance |
| `templates/project/agent-skills/savepoint-build-task/SKILL.md` | Scaffolded build workflow guidance |
| `.savepoint/Design.md` | Architecture record for lifecycle ownership and gates |

## Architectural delta

Before this epic, lifecycle rules are mostly consistent but spread across parser validation, write validation, doctor checks, board transitions, and agent-facing guidance. Recent defects show that compatibility behavior can drift when each layer restates the same rule. After this epic, `internal/data` owns the task lifecycle contract and other modules call that contract for validation, compatibility, and transitions.

## Boundaries

**In scope:**
- Task status/stage lifecycle rules, legacy `phase` and stale `stage` compatibility, doctor diagnostics, board transitions, task-to-epic handoff gates, and lifecycle documentation.

**Out of scope:**
- Adding new task statuses.
- Changing the three-column board model.
- Adding a task CRUD CLI.
- Changing defect lifecycle beyond sharing low-level stage helpers where useful.

## Quality gates

- Parser, writer, doctor, and board transition tests cover canonical and legacy lifecycle metadata.
- Doctor remains read-only and reports actionable repair suggestions.
- Board transitions still respect dependency and epic audit gates.
- Live and scaffolded templates remain aligned.
- `make build && make test` passes.

## Open decisions

- Resolved in audit (2026-05-16): `Task.Status` denormalized field removed. `Task.Column` and `Task.Stage` are the sole in-memory lifecycle fields; disk YAML key `status:` is written from `Task.Column`.
