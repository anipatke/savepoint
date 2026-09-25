---
id: E05-init-command/T004-magic-prompt-and-clipboard
status: done
objective: "Generate the post-init agent prompt and copy it to the clipboard on a best-effort basis."
depends_on: []
---

# T004: magic-prompt-and-clipboard

## Implementation Plan

- [x] Create `src/init/magic-prompt.ts` with a deterministic prompt generator that uses template-backed content rather than embedded long strings.
- [x] Keep the generated prompt short enough for terminal output and explicit enough to route an agent into the new project workflow.
- [x] Create `src/init/clipboard.ts` with platform-aware best-effort copy behavior and typed success, skipped, and failed outcomes.
- [x] Ensure clipboard failures never throw through the init success path.
- [x] Add focused tests for prompt generation, platform command selection or skip behavior, and non-fatal clipboard failure handling.
- [x] Update `.savepoint/metrics/context-bench.md` with this task's measured context before setting status to review.

## Context Log

- Files read: 15
- Estimated input tokens: ~7,625
- Notes: Read T004, AGENTS.md, E05 Design, validate-target, write-scaffold, task-building prompt, templates layer (manifest, load, render, paths), project template files (AGENTS.md, router.md), environment.ts, and both existing init test files for patterns.
