## Target File

`.savepoint/Design.md`

## Replace

last_audited: E03-cli-foundation (2026-04-27)

## With

last_audited: E04-templates-and-prompts (2026-04-27)

## Target File

`.savepoint/Design.md`

## Replace

- **Token-efficiency principle.**
  - Cold session bootstrap: ~5–7K tokens (one-time per conversation).
  - Per-task incremental: <2KB.
  - Audit: 5–15KB.
  - Anything that breaks these bounds violates the wedge.

## With

- **Token-efficiency principle.**
  - Cold session bootstrap: ~5–7K tokens (one-time per conversation).
  - Per-task incremental: <2KB.
  - Audit: 5–15KB.
  - Anything that breaks these bounds violates the wedge.
- **Template asset boundary:** established in epic `E04-templates-and-prompts` (2026-04-27). Project scaffolding and workflow prompts now live as versioned markdown/YAML assets under `templates/`, with small TypeScript helpers in `src/templates/` for manifest lookup, path resolution, loading, and interpolation.

## Target File

`.savepoint/Design.md`

## Replace

    ├── audit/
    │   └── {E##-epic}/
    │       ├── snapshot.md
    │       └── proposals/
    └── releases/

## With

    ├── audit/
    │   └── {E##-epic}/
    │       ├── snapshot.md
    │       └── proposals.md
    └── releases/

## Target File

`.savepoint/Design.md`

## Replace

- `AGENTS.md` is at root (uppercase, cross-vendor spec).
- `Design.md` lives in `.savepoint/` (working doc, not public-facing).
- `visual-identity.md` is conditional — only loaded by router for TUI/theme/visual tasks.
- **Subtasks are inline checklists** inside the task `.md` — never separate files.
- Epic folders are prefixed for navigation (`E01-scaffolding/`). Task files and IDs are prefixed the same way (`T001-package-baseline.md`).

## With

- `AGENTS.md` is at root (uppercase, cross-vendor spec).
- `Design.md` lives in `.savepoint/` (working doc, not public-facing).
- `visual-identity.md` is conditional — only loaded by router for TUI/theme/visual tasks.
- **Subtasks are inline checklists** inside the task `.md` — never separate files.
- Epic folders are prefixed for navigation (`E01-scaffolding/`). Task files and IDs are prefixed the same way (`T001-package-baseline.md`).
- The Savepoint package itself owns scaffold assets under `templates/project/`, `templates/release/v1/`, and `templates/prompts/`; generated projects receive rendered copies, not hardcoded strings from command handlers.

## Target File

`.savepoint/Design.md`

## Replace

- **Proposal bundles** use delta-shaped edits by default: `Insert After`, `Replace`, or `Delete` blocks anchored to exact text. Agents avoid rewriting whole sections unless a section is genuinely being replaced.
- **Quality review** is a section inside the proposal bundle. A standalone `quality-review.md` is only needed when the review UI explicitly asks for separate files.
- **Codebase Map** is generated mechanically in `AGENTS.md` between markers from changed-module metadata. Agents provide or refine short purpose text for new modules; they do not hand-rewrite the whole map unless the generator is unavailable.

## With

- **Proposal bundles** use delta-shaped edits by default: `Insert After`, `Replace`, or `Delete` blocks anchored to exact text. Agents avoid rewriting whole sections unless a section is genuinely being replaced.
- **Quality review** is a section inside the proposal bundle. A standalone `quality-review.md` is only needed when the review UI explicitly asks for separate files.
- **Snapshot availability is an audit precondition.** The router should enter `audit-pending` only after `.savepoint/audit/{E##-epic}/snapshot.md` exists. While the audit CLI is still a stub, the closeout task should create one manual snapshot instead of making the next agent probe for a missing file at audit startup.
- **Codebase Map** is generated mechanically in `AGENTS.md` between markers from changed-module metadata. Agents provide or refine short purpose text for new modules; they do not hand-rewrite the whole map unless the generator is unavailable.

## Target File

`.savepoint/Design.md`

## Replace

- **CLI foundation:** established in epic `E03-cli-foundation` (2026-04-27). The binary now has a testable parser, help/version handling, command dispatch, terminal capability detection, and deterministic stubs for `init`, `board`, `audit`, and `doctor`.

## With

- **CLI foundation:** established in epic `E03-cli-foundation` (2026-04-27). The binary now has a testable parser, help/version handling, command dispatch, terminal capability detection, and deterministic stubs for `init`, `board`, `audit`, and `doctor`.
- **Template and prompt assets:** established in epic `E04-templates-and-prompts` (2026-04-27). Default project, release, config, visual identity, and workflow prompt content now lives under `templates/`, with `src/templates/` exposing typed manifest, path, load, and render helpers.

## Target File

`.savepoint/Design.md`

## Replace

E03 adds focused CLI coverage for bare invocation, global help/version flags, command-level help, command dispatch, unknown top-level commands and flags, unknown command-level flags, environment detection, and command stubs. As of E03 closeout, the focused CLI runner suite reports 35 passing tests.

## With

E03 adds focused CLI coverage for bare invocation, global help/version flags, command-level help, command dispatch, unknown top-level commands and flags, unknown command-level flags, environment detection, and command stubs. As of E03 closeout, the focused CLI runner suite reports 35 passing tests.

E04 adds focused template coverage for required project/release/prompt assets, router states, embedded agent markers, audit prompt shape, template path resolution, missing-template boundary errors, and placeholder interpolation. As of the E04 audit snapshot, `npm test` reports 26 passing test files and 299 passing tests.

## Target File

`AGENTS.md`

## Replace

| Module                           | Epic                                                                            | Purpose                                                     |
| -------------------------------- | ------------------------------------------------------------------------------- | ----------------------------------------------------------- |
| `src/cli.ts`                     | [E03-cli-foundation](.savepoint/releases/v1/epics/E03-cli-foundation/Design.md) | Process entrypoint; delegates process globals to `runCli()` |
| `src/version.ts`                 | [E01-scaffolding](.savepoint/releases/v1/epics/E01-scaffolding/Design.md)       | Single source for package version string                    |
| `src/cli/args.ts`                | [E03-cli-foundation](.savepoint/releases/v1/epics/E03-cli-foundation/Design.md) | CLI argument parsing and normalization                      |
| `src/cli/environment.ts`         | [E03-cli-foundation](.savepoint/releases/v1/epics/E03-cli-foundation/Design.md) | TTY, color, and platform capability detection               |
| `src/cli/exit-codes.ts`          | [E03-cli-foundation](.savepoint/releases/v1/epics/E03-cli-foundation/Design.md) | Shared CLI exit code constants                              |
| `src/cli/help.ts`                | [E03-cli-foundation](.savepoint/releases/v1/epics/E03-cli-foundation/Design.md) | Top-level and command-level help text generation            |
| `src/cli/run.ts`                 | [E03-cli-foundation](.savepoint/releases/v1/epics/E03-cli-foundation/Design.md) | Testable CLI runner and command dispatcher                  |
| `src/commands/*.ts`              | [E03-cli-foundation](.savepoint/releases/v1/epics/E03-cli-foundation/Design.md) | Deterministic command stub handlers                         |
| `src/domain/config.ts`           | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/Design.md)         | Config defaults and typed project config model              |
| `src/domain/epic.ts`             | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/Design.md)         | Epic Design frontmatter model and validation                |
| `src/domain/ids.ts`              | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/Design.md)         | Release, epic, and task ID parsing and formatting           |
| `src/domain/release.ts`          | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/Design.md)         | Release PRD frontmatter model and validation                |
| `src/domain/router.ts`           | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/Design.md)         | Router state values and validation                          |
| `src/domain/status.ts`           | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/Design.md)         | Task status values and transition validation                |
| `src/domain/task.ts`             | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/Design.md)         | Task frontmatter/document model and validation              |
| `src/fs/markdown.ts`             | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/Design.md)         | Markdown/frontmatter reader with path-aware boundary errors |
| `src/fs/project.ts`              | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/Design.md)         | Project root discovery and scoped `.savepoint` path helpers |
| `src/readers/config.ts`          | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/Design.md)         | Read-only `config.yml` reader with defaults                 |
| `src/readers/epic.ts`            | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/Design.md)         | Read-only epic Design reader                                |
| `src/readers/release.ts`         | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/Design.md)         | Read-only release PRD reader                                |
| `src/readers/router.ts`          | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/Design.md)         | Read-only router Current state reader                       |
| `src/readers/tasks.ts`           | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/Design.md)         | Epic task-set reader with graph validation                  |
| `src/validation/dependencies.ts` | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/Design.md)         | Duplicate, missing dependency, and cycle detection          |
| `test/smoke.test.ts`             | [E01-scaffolding](.savepoint/releases/v1/epics/E01-scaffolding/Design.md)       | Baseline Vitest smoke test proving test runner works        |
| `test/cli/*.test.ts`             | [E03-cli-foundation](.savepoint/releases/v1/epics/E03-cli-foundation/Design.md) | CLI parser, help, environment, stub, and runner tests       |
| `test/domain/*.test.ts`          | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/Design.md)         | Unit tests for pure domain models and validation            |
| `test/fs/*.test.ts`              | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/Design.md)         | Unit tests for filesystem boundary helpers                  |
| `test/readers/*.test.ts`         | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/Design.md)         | Unit tests for read-only Savepoint document readers         |
| `test/validation/*.test.ts`      | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/Design.md)         | Unit tests for dependency graph validation                  |

## With

| Module                           | Epic                                                                                          | Purpose                                                     |
| -------------------------------- | --------------------------------------------------------------------------------------------- | ----------------------------------------------------------- |
| `src/cli.ts`                     | [E03-cli-foundation](.savepoint/releases/v1/epics/E03-cli-foundation/Design.md)               | Process entrypoint; delegates process globals to `runCli()` |
| `src/version.ts`                 | [E01-scaffolding](.savepoint/releases/v1/epics/E01-scaffolding/Design.md)                     | Single source for package version string                    |
| `src/cli/args.ts`                | [E03-cli-foundation](.savepoint/releases/v1/epics/E03-cli-foundation/Design.md)               | CLI argument parsing and normalization                      |
| `src/cli/environment.ts`         | [E03-cli-foundation](.savepoint/releases/v1/epics/E03-cli-foundation/Design.md)               | TTY, color, and platform capability detection               |
| `src/cli/exit-codes.ts`          | [E03-cli-foundation](.savepoint/releases/v1/epics/E03-cli-foundation/Design.md)               | Shared CLI exit code constants                              |
| `src/cli/help.ts`                | [E03-cli-foundation](.savepoint/releases/v1/epics/E03-cli-foundation/Design.md)               | Top-level and command-level help text generation            |
| `src/cli/run.ts`                 | [E03-cli-foundation](.savepoint/releases/v1/epics/E03-cli-foundation/Design.md)               | Testable CLI runner and command dispatcher                  |
| `src/commands/*.ts`              | [E03-cli-foundation](.savepoint/releases/v1/epics/E03-cli-foundation/Design.md)               | Deterministic command stub handlers                         |
| `src/domain/config.ts`           | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/Design.md)                       | Config defaults and typed project config model              |
| `src/domain/epic.ts`             | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/Design.md)                       | Epic Design frontmatter model and validation                |
| `src/domain/ids.ts`              | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/Design.md)                       | Release, epic, and task ID parsing and formatting           |
| `src/domain/release.ts`          | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/Design.md)                       | Release PRD frontmatter model and validation                |
| `src/domain/router.ts`           | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/Design.md)                       | Router state values and validation                          |
| `src/domain/status.ts`           | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/Design.md)                       | Task status values and transition validation                |
| `src/domain/task.ts`             | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/Design.md)                       | Task frontmatter/document model and validation              |
| `src/fs/markdown.ts`             | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/Design.md)                       | Markdown/frontmatter reader with path-aware boundary errors |
| `src/fs/project.ts`              | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/Design.md)                       | Project root discovery and scoped `.savepoint` path helpers |
| `src/readers/config.ts`          | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/Design.md)                       | Read-only `config.yml` reader with defaults                 |
| `src/readers/epic.ts`            | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/Design.md)                       | Read-only epic Design reader                                |
| `src/readers/release.ts`         | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/Design.md)                       | Read-only release PRD reader                                |
| `src/readers/router.ts`          | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/Design.md)                       | Read-only router Current state reader                       |
| `src/readers/tasks.ts`           | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/Design.md)                       | Epic task-set reader with graph validation                  |
| `src/templates/*.ts`             | [E04-templates-and-prompts](.savepoint/releases/v1/epics/E04-templates-and-prompts/Design.md) | Template manifest, path, load, and render helpers           |
| `src/validation/dependencies.ts` | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/Design.md)                       | Duplicate, missing dependency, and cycle detection          |
| `templates/project/**`           | [E04-templates-and-prompts](.savepoint/releases/v1/epics/E04-templates-and-prompts/Design.md) | Default project scaffold markdown and YAML assets           |
| `templates/release/**`           | [E04-templates-and-prompts](.savepoint/releases/v1/epics/E04-templates-and-prompts/Design.md) | Release starter PRD templates                               |
| `templates/prompts/*.prompt.md`  | [E04-templates-and-prompts](.savepoint/releases/v1/epics/E04-templates-and-prompts/Design.md) | Agent workflow prompt templates                             |
| `test/smoke.test.ts`             | [E01-scaffolding](.savepoint/releases/v1/epics/E01-scaffolding/Design.md)                     | Baseline Vitest smoke test proving test runner works        |
| `test/cli/*.test.ts`             | [E03-cli-foundation](.savepoint/releases/v1/epics/E03-cli-foundation/Design.md)               | CLI parser, help, environment, stub, and runner tests       |
| `test/domain/*.test.ts`          | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/Design.md)                       | Unit tests for pure domain models and validation            |
| `test/fs/*.test.ts`              | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/Design.md)                       | Unit tests for filesystem boundary helpers                  |
| `test/readers/*.test.ts`         | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/Design.md)                       | Unit tests for read-only Savepoint document readers         |
| `test/templates/*.test.ts`       | [E04-templates-and-prompts](.savepoint/releases/v1/epics/E04-templates-and-prompts/Design.md) | Template asset, prompt, router, and render integrity tests  |
| `test/validation/*.test.ts`      | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/Design.md)                       | Unit tests for dependency graph validation                  |

## Target File

`templates/prompts/audit-reconciliation.prompt.md`

## Replace

- `.savepoint/audit/{E##-slug}/snapshot.md` — what changed during the epic.
- `.savepoint/releases/v{N}/epics/{E##-slug}/Design.md` — the epic design (may have deltas from the original plan).
- `.savepoint/Design.md` — project architecture.
- `.savepoint/AGENTS.md` — agent guide and codebase map.
- Only the source and test files listed in the snapshot.

## With

- `.savepoint/audit/{E##-slug}/snapshot.md` — what changed during the epic. If this file is missing while the audit CLI is still unavailable, create one manual snapshot from the known epic scope once; do not search broadly for replacement inputs.
- `.savepoint/releases/v{N}/epics/{E##-slug}/Design.md` — the epic design (may have deltas from the original plan).
- `.savepoint/Design.md` — project architecture.
- `AGENTS.md` — agent guide and codebase map.
- Only the source and test files listed in the snapshot.

## Target File

`.savepoint/releases/v1/epics/E04-templates-and-prompts/Design.md`

## Replace

## Architectural delta

Before this epic, project scaffolding would need to invent file contents in code. After this epic, generated Savepoint projects come from data files that agents and humans can inspect.

Templates become a stable internal asset boundary consumed by `init-command`, audit prompts, and future workflow docs.

## With

## Architectural delta

Before this epic, project scaffolding would need to invent file contents in code. After this epic, generated Savepoint projects come from data files that agents and humans can inspect.

Templates become a stable internal asset boundary consumed by `init-command`, audit prompts, and future workflow docs.

## Implemented As

- Project scaffold assets were added under `templates/project/`, including `AGENTS.md`, `.savepoint/router.md`, `.savepoint/PRD.md`, `.savepoint/Design.md`, `.savepoint/config.yml`, and `.savepoint/visual-identity.md`.
- Release starter content was added at `templates/release/v1/PRD.md`.
- Prompt assets were added under `templates/prompts/` for PRD creation, project design, epic design, task breakdown, task planning, task building, and audit reconciliation.
- `src/templates/manifest.ts` defines typed template sets for project, release, and prompt assets.
- `src/templates/paths.ts` resolves template roots and family-specific paths.
- `src/templates/load.ts` loads named templates from disk and returns path-aware `not_found` errors.
- `src/templates/render.ts` interpolates `PROJECT_NAME`, `RELEASE_NUMBER`, and `RELEASE_NAME` placeholders while leaving unresolved placeholders intact as a visible failure signal.
- `src/templates/index.ts` exposes a narrow public helper surface.
- Template tests were added under `test/templates/` for required files, frontmatter, router states, agent instruction markers, audit prompt shape, path lookup, loading failures, and interpolation behavior.
- The audit snapshot was generated manually because `savepoint audit` is still a stub.

Design delta notes:

- The loader reports all read failures as `not_found`; richer filesystem error typing is deferred until command handlers need to distinguish missing files from permission or IO failures.
- The audit prompt correctly requires one proposal bundle, but audit startup exposed a workflow gap: the router can point agents at a snapshot before one exists. The follow-up improvement is to make snapshot creation an explicit precondition before entering `audit-pending`, or to instruct the agent to create one manual snapshot once without searching broadly.

## Quality Review

## Must Fix Before Close

None.

## Carry Forward

- Audit initiation currently assumes `.savepoint/audit/{E##-epic}/snapshot.md` exists. In this audit, that caused the agent to probe a missing file before reconstructing the snapshot manually. Future audit closeout should create the snapshot before setting `state: audit-pending`, or the router/audit prompt should say to create one manual snapshot once when the CLI is unavailable.
- `loadTemplate()` currently maps every `readFile` failure to `not_found`. This is acceptable for E04's missing-template boundary, but future command code may need distinct `permission_denied` or `read_failed` errors for better diagnostics.

## Already Fixed

- E04 T005 is marked `done`.
- A manual E04 audit snapshot now exists at `.savepoint/audit/E04-templates-and-prompts/snapshot.md`.
- The audit proposal bundle is a single `.savepoint/audit/E04-templates-and-prompts/proposals.md` file, matching the E04 audit prompt.
- Quality gates passed on 2026-04-27: `npm run build`, `npm run typecheck`, `npm run lint`, `npm run format:check`, and `npm test` with 26 passing test files and 299 passing tests.
