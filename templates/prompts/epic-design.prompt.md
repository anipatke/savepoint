# Prompt: Epic Design Creation

<!-- AGENT: You have a release PRD and a project Design. Write one epic Design.md. Stop when that single epic design is complete. -->

You are designing one epic for a Savepoint-managed release.

## Input

- `.savepoint/releases/v{N}/PRD.md` — release scope and epic list.
- `.savepoint/Design.md` — project architecture and constraints.
- The epic name and one-sentence purpose from the release PRD table.

## Output

A single markdown file at `.savepoint/releases/v{N}/epics/{E##-slug}/Design.md`:

```
---
type: epic-design
status: active
---

# Epic {E##}: {slug}

## Purpose
## What this epic adds
## Components and files
## Architectural delta
## Boundaries
## Quality gates
## Design constraints
## Open decisions
```

## Rules

1. **One epic only.** Do not write tasks, do not design other epics, and do not modify the release PRD.
2. **Purpose is one paragraph.** Explain what the epic adds to the system and why it is ordered where it is.
3. **Components and files must be a table.** Each row needs a path and a one-sentence purpose.
4. **Architectural delta explains before/after.** What does the system lack before this epic? What will it have after?
5. **Boundaries split in/out of scope explicitly.** Use "In scope" and "Out of scope" sub-lists.
6. **Quality gates are testable.** List commands, file checks, or behavioral criteria that prove the epic is done.
7. **Open decisions are honest.** If something is unclear, state it as an open decision rather than guessing.
