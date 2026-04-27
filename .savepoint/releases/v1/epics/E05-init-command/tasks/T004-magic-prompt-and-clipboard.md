---
id: E05-init-command/T004-magic-prompt-and-clipboard
status: planned
objective: "Generate the post-init agent prompt and copy it to the clipboard on a best-effort basis."
depends_on: []
---

# T004: magic-prompt-and-clipboard

## Implementation Plan

- [ ] Create `src/init/magic-prompt.ts` with a deterministic prompt generator that uses template-backed content rather than embedded long strings.
- [ ] Keep the generated prompt short enough for terminal output and explicit enough to route an agent into the new project workflow.
- [ ] Create `src/init/clipboard.ts` with platform-aware best-effort copy behavior and typed success, skipped, and failed outcomes.
- [ ] Ensure clipboard failures never throw through the init success path.
- [ ] Add focused tests for prompt generation, platform command selection or skip behavior, and non-fatal clipboard failure handling.
