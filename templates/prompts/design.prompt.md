# Prompt: Project Design Creation

<!-- AGENT: You have an approved project PRD. Produce a project-level Design.md. Stop when the design is complete; do not proceed to epics or tasks. -->

You are architecting a Savepoint-managed project based on an approved PRD.

## Input

- `.savepoint/PRD.md` — the project vision and constraints.
- Any additional context the user provides about tech stack or patterns.

## Output

A single markdown file matching the structure of `.savepoint/Design.md`:

```
---
type: project-design
status: active
last_audited: never
---

# {Project Name} — System Architecture

## 1. Architecture model
## 2. Directory layout
## 3. Hierarchy semantics
## 4. Status model & gates
## 5. Dependencies
## 6. CLI surface
## 7. Agent audit workflow
## 8. Testing strategy
## 9. Release versioning
```

## Rules

1. **Design for the PRD, not for the future.** Include only patterns and components needed to satisfy the approved scope.
2. **Directory layout should be specific.** List expected top-level directories and key files; avoid vague bullets.
3. **Status model should name every state** and describe what gates move work between them.
4. **CLI surface is optional.** If the project has no CLI, state that explicitly.
5. **Do not write epics or tasks.** The Design is architecture only.
6. **Keep visual identity separate.** If the project has UI/TUI/theme concerns, note that `.savepoint/visual-identity.md` will carry them.
7. **Audit is agent-led.** Do not design a `savepoint audit` CLI pipeline. The audit gate is performed by a fresh agent using the audit skill, producing one epic-local `E##-Audit.md` file.
8. **Audit file structure is fixed.** `E##-Audit.md` user-facing sections are `## Main Findings` and `## Code Style Review`; file-specific replacement blocks belong under `## Proposed Changes` as admin/apply metadata and should not be rendered in the Epic Detail Audit tab.
