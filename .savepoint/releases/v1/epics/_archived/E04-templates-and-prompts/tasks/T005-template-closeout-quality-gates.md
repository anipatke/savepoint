---
id: E04-templates-and-prompts/T005-template-closeout-quality-gates
status: done
objective: "Run the E04 closeout checks and prepare the epic for audit after template assets and helpers are complete."
depends_on: ["E04-templates-and-prompts/T004-template-integrity-tests"]
---

# T005: template-closeout-quality-gates

## Implementation Plan

- [x] Review the E04 Design against implemented template assets, prompt assets, helpers, and tests for any missing scope.
- [x] Run the focused template test suite.
- [x] Run `npm run typecheck`, `npm run lint`, `npm run format:check`, and `npm test` for epic closeout.
- [x] Update E04 task or design notes only for factual deltas discovered during closeout.
- [x] Set this task to review when all checks pass or document any blocked quality gate with the exact command and failure.
