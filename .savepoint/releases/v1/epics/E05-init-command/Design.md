---
type: epic-design
status: active
---

# Epic E05: init-command

## Purpose

Implement `savepoint init`, the first-run command that creates a Savepoint workflow in an empty or compatible project directory and prints the magic prompt the user gives to an AI agent.

## What this epic adds

- Target directory validation.
- Scaffold writing from template assets.
- Project name and release metadata interpolation.
- Safe handling for existing files.
- Magic prompt output to stdout.
- Best-effort clipboard copy.
- Optional one-shot dev-dependency installation prompt/flag support as defined by the CLI contract.
- Integration tests using temporary directories.

## Components and files

Expected files introduced or extended by this epic:

| Path                          | Purpose                                           |
| ----------------------------- | ------------------------------------------------- |
| `src/commands/init.ts`        | `savepoint init` command implementation.          |
| `src/init/validate-target.ts` | Target directory checks.                          |
| `src/init/write-scaffold.ts`  | Template copy/write orchestration.                |
| `src/init/magic-prompt.ts`    | Prompt text generation.                           |
| `src/init/clipboard.ts`       | Best-effort clipboard integration.                |
| `src/init/dev-deps.ts`        | Optional dependency-install command construction. |
| `test/init/**/*.test.ts`      | Unit and integration tests for init behavior.     |

## Architectural delta

Before this epic, Savepoint has assets and a command shell but cannot create a project. After this epic, `npx savepoint init` can produce the file-only workflow scaffold consumed by agents.

This epic is the first user-visible workflow command, but it should still rely on the existing CLI and template layers rather than embedding content or parsing rules locally.

## Boundaries

In scope:

- Create `.savepoint/` and root `AGENTS.md`.
- Print the initial agent prompt.
- Avoid overwriting user files unless explicitly allowed by flags.
- Support Windows, macOS, and Linux filesystem behavior.

Out of scope:

- Creating epics/tasks interactively.
- Running the board after init.
- Implementing audit or doctor.
- Managing package installation beyond the optional one-shot path.

## Quality gates

- Integration tests should cover empty directories, existing incompatible files, and successful scaffold creation.
- Clipboard failures should not fail init if the scaffold was written successfully.
- File writes should be atomic enough to avoid partially corrupting existing files.

## Design constraints

- Use the template layer as the only source for generated file content.
- Keep user-facing output short and copyable.
- Treat filesystem writes as a boundary with clear errors.
