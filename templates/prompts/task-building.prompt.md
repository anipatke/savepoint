# Prompt: Task Building

<!-- AGENT: Execute the task plan. Tick checkboxes as you complete them. Run focused tests. Stop at review. -->

You are implementing one Savepoint task that is `in_progress` with all `depends_on` done.

## Input

- `.savepoint/router.md` — confirms you are on the correct task.
- `.savepoint/releases/v{N}/epics/{E##-slug}/Design.md` — epic context.
- The active task file: `.savepoint/releases/v{N}/epics/{E##-slug}/tasks/TNNN-slug.md` — your checklist.
- Any directly touched source or test files.

## Output

- Code changes that satisfy every checkbox in the task's `## Implementation Plan`.
- Updated task file with all checkboxes ticked: `- [x]`.
- Task frontmatter `status: review`.

## Rules

1. **Follow the checklist in order.** Do not skip steps.
2. **Run focused tests first.** If the task touches new behavior, write or run tests for that behavior before moving to the next checkbox.
3. **Reserve full quality-gate suite for closeout.** Run `npm run build`, `npm run typecheck`, `npm run lint`, and `npm test` after all checkboxes are ticked, or only if the task explicitly says so.
4. **Make minimal changes.** Touch as few files as possible. Small diffs are easier to audit.
5. **Stop at review.** Do not start the next task, do not update the router, and do not commit.
6. **If a checkbox is blocked,** leave it unchecked, add a `<!-- BLOCKED: reason -->` comment under it, and stop.
