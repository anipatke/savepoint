# Prompt: Task Planning (Repair or Late Add)

<!-- AGENT: A single task needs a plan. Write its Implementation Plan as checkboxes and set status: planned. Do not implement. -->

You are repairing or late-adding one task to an epic that is already in progress.

## Input

- The task file stub (may have an empty or incomplete `## Implementation Plan`).
- `.savepoint/releases/v{N}/epics/{E##-slug}/Design.md` — epic context.
- `.savepoint/router.md` — current state.

## Output

The updated task file with:

```yaml
---
id: {E##-slug}/TNNN-slug
status: planned
objective: "<one sentence>"
depends_on: ["..."]
---

# TNNN: slug

## Acceptance Criteria

- [ ] Criterion one (observable outcome)
- [ ] Criterion two (observable outcome)

## Implementation Plan

- [ ] Step one
- [ ] Step two

## Context Log

- Files read:
- Estimated input tokens:
- Notes:
```

## Rules

1. **One task only.** Do not touch other tasks, the epic Design, or the release PRD.
2. **Respect existing depends_on.** Do not remove dependencies; you may add them if the task truly requires prior work.
3. **Acceptance criteria come first.** Write `## Acceptance Criteria` (observable outcomes) before `## Implementation Plan` (build checklist).
4. **Plan before code.** Every task must have checkboxes before an agent implements it.
5. **Use only canonical task status.** Planned task frontmatter must use `status: planned`. Do not write `todo`, `doing`, `blocked`, `review`, `audit`, or phase names into `status`.
6. **Use phase separately.** Only implementation flow may add `phase: build`, `phase: test`, or `phase: audit`, and only when `status: in_progress`.
7. **Do not implement.** Set `status: planned` and stop.
8. **If the task is a repair,** reference the defect or gap that triggered it in a `## Repair Context` section above the plan.
