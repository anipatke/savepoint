---
id: E05-init-command/T006-init-command-integration
status: planned
objective: "Wire `savepoint init` end to end from CLI entrypoint to scaffold, prompt output, clipboard, and optional install behavior."
depends_on:
  - E05-init-command/T001-init-cli-contract
  - E05-init-command/T002-target-validation
  - E05-init-command/T003-scaffold-writer
  - E05-init-command/T004-magic-prompt-and-clipboard
  - E05-init-command/T005-dev-deps-install-option
---

# T006: init-command-integration

## Implementation Plan

- [ ] Replace the init stub with orchestration that validates the target, writes the scaffold, prints the magic prompt, attempts clipboard copy, and handles optional install.
- [ ] Return success, usage, and boundary-error exit codes consistently with the existing CLI runner.
- [ ] Keep stdout short and copyable while sending actionable failure details to stderr.
- [ ] Add end-to-end temporary-directory tests for empty directory success, incompatible existing files, already initialized project handling, clipboard failure success, and install failure reporting.
- [ ] Run focused init tests first, then the full quality gate suite for the epic closeout path.
