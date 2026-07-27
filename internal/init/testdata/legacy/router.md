<!-- FROZEN FIXTURE: a pre-v1.5 project router, kept byte-for-byte. Do not
     re-sync with templates/project/.savepoint/router.md. A failure here is a
     compatibility break in the reader. -->

# Agent State Machine

## Read order

1. This file (router.md)
2. Current state → next action
3. Active epic E##-Detail.md
4. Active task file

## Current state

```yaml
state: task-building
release: v1
epic: E03-board-tui
task: E03-board-tui/T002-columns
next_action: "Build E03-board-tui/T002-columns."
updated: 2024-11-02
```

## State → action

### task-building

Task `in_progress`, depends satisfied.

**Next:** Execute plan, tick checkboxes, run quality gates, update router to next task.
