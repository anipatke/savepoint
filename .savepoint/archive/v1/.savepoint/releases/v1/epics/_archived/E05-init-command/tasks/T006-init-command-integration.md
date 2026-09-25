---
id: E05-init-command/T006-init-command-integration
status: done
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

- [x] Replace the init stub with orchestration that validates the target, writes the scaffold, prints the magic prompt, attempts clipboard copy, and handles optional install.
- [x] Return success, usage, and boundary-error exit codes consistently with the existing CLI runner.
- [x] Keep stdout short and copyable while sending actionable failure details to stderr.
- [x] Add end-to-end temporary-directory tests for empty directory success, incompatible existing files, already initialized project handling, clipboard failure success, and install failure reporting.
- [x] Run focused init tests first, then the full quality gate suite for the epic closeout path.
- [x] Update `.savepoint/metrics/context-bench.md` with this task's measured context before setting status to review.

## Context Log

- Files read: 21 (AGENTS.md, E05 Design, T006, T001, validate-target.ts, write-scaffold.ts, magic-prompt.ts, clipboard.ts, dev-deps.ts, args.ts, run.ts, cli.ts, exit-codes.ts, result.ts, run.test.ts, validate-target.test.ts, write-scaffold.test.ts, manifest.ts, render.ts, paths.ts, context-bench.md)
- Estimated input tokens: ~10,925
- Notes: runCli made async; InitResult replaces CommandResult for init; EXIT_BOUNDARY_ERROR added; InitDeps injection for clipboard/install test overrides; stale commands.test.ts runInit block removed.
