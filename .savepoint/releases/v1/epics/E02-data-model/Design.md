---
type: epic-design
status: audited
audited: 2026-04-27
---

# Epic E02: data-model

## Context Budget

For implementation tasks in this epic, read only:

- `.savepoint/router.md`
- this epic `Design.md`
- the active task file
- directly touched source/test files

Read `.savepoint/Design.md` only if the task changes architecture. Read `.savepoint/releases/v1/PRD.md`, prior epic docs, audit proposals, or `.savepoint/visual-identity.md` only when the active task explicitly requires them.

## Purpose

Build the core file-backed model for releases, epics, tasks, statuses, dependencies, and project configuration. This epic gives later CLI and TUI work a typed domain layer instead of each command parsing markdown independently.

## What this epic adds

- Typed models for releases, epics, tasks, task statuses, dependency IDs, config, and router state.
- Markdown frontmatter parsing for task and planning files.
- Readers for `.savepoint/releases/{release}/PRD.md`, epic `Design.md`, and task files.
- Status transition validation for `backlog`, `planned`, `in_progress`, `review`, and `done`.
- Dependency resolution, missing dependency detection, and cycle detection.
- Shared filesystem helpers for scoped `.savepoint` reads.
- Unit tests for parsing, validation, dependency logic, and state transitions.

## Implementation strategy

- Use structured parsers for YAML/frontmatter and Markdown boundaries.
- Keep domain types and validation pure; filesystem reads belong in `src/fs/*` and `src/readers/*`.
- Keep writes out of scope except temporary test fixtures.
- Make ID parsing understand prefixed epic and task IDs such as `E02-data-model/T001-readers`.
- Return boundary errors that include the path and reason without forcing callers to parse error text.

## Components and files

Expected files introduced or extended by this epic:

| Path                    | Purpose                                             |
| ----------------------- | --------------------------------------------------- |
| `src/domain/status.ts`  | Status values and transition rules.                 |
| `src/domain/ids.ts`     | Release, epic, and task ID parsing/formatting.      |
| `src/domain/task.ts`    | Task frontmatter and task document types.           |
| `src/domain/release.ts` | Release PRD metadata types.                         |
| `src/domain/config.ts`  | `config.yml` type and defaults.                     |
| `src/fs/markdown.ts`    | Markdown/frontmatter read helpers.                  |
| `src/fs/project.ts`     | Project root and `.savepoint` path helpers.         |
| `src/readers/*.ts`      | Release, epic, task, router, and config readers.    |
| `src/validation/*.ts`   | Status, dependency, duplicate ID, and cycle checks. |
| `test/**/*.test.ts`     | Unit tests for the model and validation behavior.   |

## Architectural delta

Before this epic, Savepoint has only a package scaffold. After this epic, it has a reusable domain and read-only persistence layer for the file-based state machine.

This epic intentionally avoids write operations except where tests create fixtures. Writes belong to later command/TUI/audit work once behavior exists.

## Boundaries

In scope:

- Parse markdown frontmatter with clear boundary errors.
- Model task dependencies using repo-relative task IDs.
- Detect cycles and broken dependencies.
- Expose pure validation functions usable by CLI, TUI, audit, and doctor.

Out of scope:

- Creating or editing tasks.
- Rendering any TUI.
- Implementing command dispatch.
- Running quality gates.
- Writing audit proposals.

## Quality gates

- Branching parser and validator logic must have unit tests.
- Fixture files should be small and local to tests.
- Errors should identify the file path and failure reason.

## Design constraints

- Prefer structured parsers over ad hoc string splitting.
- Keep filesystem access at boundaries; domain validation should be pure.
- Use one source of truth for status transitions and dependency rules.

## Implemented as

Implemented on 2026-04-27 across tasks T001-T007. The read-only data model landed with focused unit tests and passing mechanical gates after formatting the E02 task markdown.

**Implemented components:**

| Path                             | Result                                                                  |
| -------------------------------- | ----------------------------------------------------------------------- |
| `src/domain/ids.ts`              | Release, epic, and full task ID parsing plus formatting helpers.        |
| `src/domain/status.ts`           | Canonical task status values and transition validation.                 |
| `src/domain/task.ts`             | Task frontmatter and task document validation.                          |
| `src/domain/release.ts`          | Release PRD frontmatter validation.                                     |
| `src/domain/epic.ts`             | Epic Design frontmatter validation.                                     |
| `src/domain/router.ts`           | Router state value and Current state validation.                        |
| `src/domain/config.ts`           | Config defaults and merge behavior.                                     |
| `src/fs/markdown.ts`             | Structured YAML frontmatter parsing with path-aware boundary errors.    |
| `src/fs/project.ts`              | `.savepoint` root discovery and scoped path construction.               |
| `src/readers/*.ts`               | Read-only release, epic, router, config, and task-set readers.          |
| `src/validation/dependencies.ts` | Duplicate task ID, missing dependency, and cycle detection.             |
| `test/**/*.test.ts`              | Branch coverage for parser, reader, status, config, and graph behavior. |

**Deviations from original design:**

| Item                          | Deviation            | Reason / follow-up                                                                      |
| ----------------------------- | -------------------- | --------------------------------------------------------------------------------------- |
| `done` status transitions     | Fixed during audit   | `done -> review` is allowed so completed work can be reopened for audit-stale rechecks. |
| Config accent validation      | Fixed during audit   | Non-string custom accent values are ignored before merging defaults.                    |
| Audit snapshot/router flip    | Manual               | Expected until E07 implements the audit CLI pipeline.                                   |
| `.claude/settings.local.json` | Local untracked file | Decide whether to add `.claude/` to `.gitignore` before commit.                         |

**Quality gates:**

| Gate                   | Result                                                                                   |
| ---------------------- | ---------------------------------------------------------------------------------------- |
| `npm run build`        | Passed.                                                                                  |
| `npm run typecheck`    | Passed.                                                                                  |
| `npm run lint`         | Passed.                                                                                  |
| `npm run format:check` | Passed after formatting E02 task markdown.                                               |
| `npm test`             | Passed outside sandbox: 16 test files, 176 tests. Sandbox run hit esbuild `spawn EPERM`. |

## Close Criteria

- All `E02-data-model` tasks are `done`.
- `npm run build`, `npm run typecheck`, `npm run lint`, `npm run format:check`, and `npm test` pass.
- Tests cover parser errors, valid documents, invalid statuses, allowed and disallowed status transitions, missing dependencies, duplicate IDs, and dependency cycles.
- Audit snapshot exists at `.savepoint/audit/E02-data-model/snapshot.md`.
- Audit proposals are accepted, rejected, or explicitly carried forward.
- This epic `Design.md` has `status: audited`.
- Router points to the next epic state.
