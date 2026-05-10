---
type: epic-design
status: planned
---

# E17: Defect Workflow TUI

## Purpose

Introduce first-class defect documentation and TUI visibility for observed correctness fixes without forcing defects into the planned epic/task model. Defects remain release-level repair work, while the board stays focused on the existing task flow.

## What this epic adds

- Release-level defect files under `.savepoint/releases/{release}/defects/`.
- A `Defect` data model and discovery path parallel to tasks.
- Router support for `defect-building` priority.
- A header or Next Activity defect signal showing open defect count.
- A keyboard-driven Defects overlay for browsing open, in-progress, and resolved defects.
- A defect detail overlay showing symptom, expected behavior, reproduction, impact, fix plan, acceptance criteria, and resolution notes.
- Optional related-defect markers on task cards when defects reference the selected epic/task.
- Doctor diagnostics for defect frontmatter, status/stage rules, and broken references.

## Components

| Module | Purpose |
|--------|---------|
| `internal/data` | Add defect model, frontmatter parsing, discovery, and validation |
| `internal/board` | Add defect state, overlay navigation, summary rendering, detail rendering, and router priority support |
| `internal/doctor` | Validate defect structure and references |
| `templates/project` | Add defect documentation/template guidance for scaffolded projects |
| `.savepoint/releases/{release}/defects` | Store release-level defect files |

## Implemented As

- Defect model, parser, release-scoped discovery, and lifecycle validation live in `internal/data`.
- Board loading reads release defects alongside tasks and renders an open-defect header signal without adding a fourth column.
- The `d` key opens a defect overlay; enter opens a defect detail overlay that renders the evidence sections from the defect markdown body.
- Related task cards can show compact defect markers when a defect `reference` matches a visible task.
- Doctor validates defect files for parse errors, required fields, invalid status/stage, and broken task-like references.
- Scaffolded AGENTS/router templates and README now document the defect lane, but `defect-building` state guidance needs the audit proposal updates before close.

## Boundaries

**In scope:**
- Release-level defect markdown files
- TUI defect count, overlay, detail view, and router priority display
- Defect validation in doctor
- Defect documentation and scaffold templates
- Focused tests for data parsing, board rendering, and diagnostics

**Out of scope:**
- A fourth board column
- Treating defects as epics or task files
- Automatic defect creation from runtime failures
- New CLI commands for creating, editing, or closing defects
- Changing task status ownership rules
