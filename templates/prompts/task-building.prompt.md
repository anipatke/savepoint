# Prompt: Task Building

<!-- AGENT: Execute the task plan. Tick checkboxes as you complete them. Run focused tests. Stop at done. -->

You are implementing one Savepoint task that is `in_progress` with all `depends_on` done.

## Input

- `.savepoint/router.md` — confirms you are on the correct task.
- `.savepoint/releases/v{N}/epics/{E##-slug}/E##-Detail.md` — epic context.
- The active task file: `.savepoint/releases/v{N}/epics/{E##-slug}/tasks/TNNN-slug.md` — your checklist.
- Any directly touched source or test files.

## Output

- Code changes that satisfy every line in the task's `## Acceptance Criteria`.
- Updated task file with all checkboxes ticked: `- [x]`.
- Updated task file `## Context Log` with files read, estimated input tokens, notes, and quality-gate results.
- Task frontmatter `status: done`.

## Status Rules

- Task frontmatter `status` must only be `planned`, `in_progress`, or `done`.
- Active phases are separate: `phase: build`, `phase: test`, or `phase: audit`.
- `phase` is only valid with `status: in_progress`; remove `phase` when setting `status: planned` or `status: done`.
- Never write `todo`, `doing`, `blocked`, `review`, `audit`, or phase names into `status`.

## Rules

1. **Follow the checklist in order.** Do not skip steps.
2. **Satisfy acceptance criteria first.** Every `## Acceptance Criteria` line must have a passing test or a verified manual outcome before the task can be `done`. The implementation checklist exists to produce those outcomes; every checked box should map to an observable result.
3. **Run focused tests first.** If the task touches new behavior, write or run tests for that behavior before moving to the next checkbox.
4. **Run full quality gates at closeout.** After all checkboxes are ticked, run the project's full quality-gate suite (build, typecheck, lint, format:check, test). Record pass/fail for each gate in the Context Log.
5. **If a gate fails,** fix it or document the blocker in the task file before setting `status: done`.
6. **Make minimal changes.** Touch as few files as possible. Small diffs are easier to audit.
7. **Stop at done. Update the router. Prompt the user.** Do not start the next task, do not commit, and do not advance without user acknowledgment.
8. **If a checkbox is blocked,** leave it unchecked, add a `<!-- BLOCKED: reason -->` comment under it, and stop.
9. **Drift check before handoff.** Ask yourself: did this task add new files/modules not in the Codebase Map? Did it change the architecture from Design.md? If yes, append `## Drift Notes` to the task file.
