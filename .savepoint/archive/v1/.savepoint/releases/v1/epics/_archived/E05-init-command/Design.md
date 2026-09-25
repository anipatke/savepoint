---
type: epic-design
status: audited
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

## Implemented As

- `src/cli/args.ts`, `src/cli/help.ts`, and `src/cli/run.ts` now recognize `savepoint init [dir]`, `--force`/`-f`, `--install`/`-y`, and init command help.
- `src/commands/init.ts` orchestrates target validation, scaffold writing, magic prompt generation, best-effort clipboard copy, and optional install execution.
- `src/init/validate-target.ts` checks missing, empty, compatible, already-initialized, conflicting, and boundary-error target states.
- `src/init/write-scaffold.ts` loads project templates through the template layer, renders project metadata, creates parent directories, refuses conflicts by default, and writes via temp-file-plus-rename.
- `src/init/magic-prompt.ts` renders `templates/prompts/magic-prompt.prompt.md` rather than embedding prompt copy in code.
- `src/init/clipboard.ts` selects platform clipboard commands and returns success, skipped, or failed results without throwing through the init success path.
- `src/init/dev-deps.ts` detects package managers from lockfiles and builds or runs the matching dev-dependency install command.
- Project scaffold templates now include `.savepoint/metrics/context-bench.md` so generated projects can track context load from the start.
- Init tests cover parser behavior, validation, scaffold writes, prompt generation, clipboard behavior, optional install handling, and end-to-end temporary-directory init flows.
- The audit snapshot was generated manually because `savepoint audit` is still a stub.

Design delta notes:

- `runCli()` and `src/cli.ts` became async so the init command can perform filesystem, clipboard, and subprocess work.
- The optional install path is explicit flag-driven behavior in this implementation; interactive prompting is not implemented.
- Clipboard failures and install failures do not fail the scaffold path; install failures are reported on stderr after successful scaffold and prompt output.

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
