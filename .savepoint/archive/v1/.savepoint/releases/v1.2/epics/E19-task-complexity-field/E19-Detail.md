---
type: epic-design
status: audited
---

# E19: Task Complexity Field

## Purpose

Add an explicit complexity signal to task planning so the board and task-writing workflow can show how hard a task is and which implementation path it is likely to need.

## What this epic adds

- Task frontmatter fields for a complexity tier and a short reason.
- Validation and persistence that keep complexity stable across parse and write flows.
- Task card and task detail rendering that expose complexity in the TUI.
- An updated `create-task` skill that assigns complexity with the shared rubric.
- Backfill of the existing v1.2 task set so release planning uses the new field consistently.

## Components and files

| Module | Purpose |
|--------|---------|
| `internal/data` | Task model, parsing, validation, write helpers, and parser tests |
| `internal/doctor` | Validation for malformed or missing complexity metadata |
| `internal/board` | Card/detail rendering and tests for the complexity display |
| `agent-skills/savepoint-create-task/SKILL.md` | Canonical root task-planning guidance |
| `templates/project/agent-skills/savepoint-create-task/SKILL.md` | Scaffolded copy of the task-planning guidance |
| `internal/init/template_freshness_test.go` | Protect the live/scaffolded skill sync surface |
| `.savepoint/releases/v1.2/epics/*/tasks/*.md` | Backfill the existing v1.2 task files with complexity metadata |

## Architectural delta

Before this epic, task files carry status and phase but not implementation complexity, so the board and planner cannot signal how hard a task is or what sort of agent it should attract. After this epic, complexity becomes part of the task contract, the TUI can show it at a glance, and the planning workflow records the rationale in a consistent format.

## Boundaries

**In scope:**
- Data model, validation, TUI rendering, create-task skill updates, v1.2 backfill, and tests.

**Out of scope:**
- Changing router behavior.
- Inventing per-agent assignment logic.
- Adding new task lifecycle states.

## Quality gates

- Complexity round-trips through task parse and write paths.
- Task cards and the task detail view display the new field without breaking existing layout.
- Root and scaffolded create-task skills match and encode the same rubric.
- Existing v1.2 task files are updated and `make build && make test` passes.

## Open decisions

None.
