# Prompt: Project PRD Creation

<!-- AGENT: Read the user's intent, then interview them until you can write a concise project PRD. Stop when the PRD is complete; do not proceed to design or implementation. -->

You are helping a user define a new Savepoint-managed project. Your goal is to produce a project-level PRD file that an agent can use later to design and build the system.

## Input

The user may provide:

- A rough idea, feature list, or problem statement.
- Existing code, sketches, or references.
- Constraints (tech stack, timeline, team size).

## Output

A single markdown file matching the structure of `.savepoint/PRD.md`:

```
---
type: project-prd
status: active
---

# {Project Name} — Product Vision

## What it is
## Why
## Target user
## Headline differentiator
## Success metrics
## Constraints
## Out of scope (forever or for now)
```

## Rules

1. **One paragraph per section** unless the user explicitly provides more detail.
2. **Out of scope is mandatory.** If the user does not specify it, ask what they are explicitly not building.
3. **Do not invent implementation details.** The PRD is about problems and outcomes, not libraries or file structures.
4. **Stop after the PRD.** Do not suggest next steps, design, or task breakdown.
5. **Use concise language.** Token discipline starts here.
