## Target File

`.savepoint/Design.md`

## Replace

`last_audited: E01-scaffolding (2026-04-27)`

## With

`last_audited: E02-data-model (2026-04-27)`

## Replace

`## 4. Status model & gates`

## With

```md
## 4. Status model & gates

Five statuses, with explicit gates:

| Status        | Meaning                           | Entry gate                                    |
| ------------- | --------------------------------- | --------------------------------------------- |
| `backlog`     | Task exists, no plan              | created                                       |
| `planned`     | Implementation plan written       | plan section non-empty                        |
| `in_progress` | AI building                       | all `depends_on` are `done`                   |
| `review`      | Build done, awaiting verification | all sub-tasks checked off                     |
| `done`        | Verified, locked                  | assertion + (configurable) quality gates pass |

- `blocked` is a **flag**, not a status — `in_progress` + `blocked: "reason"` is valid.
- Status names and transition validation live in `src/domain/status.ts`.
- `done -> review` is allowed so completed work can be reopened for audit-stale rechecks. The warning itself belongs to future CLI/TUI write layers.
- **Verification mode is configurable per project** (`verify_strict: true|false`). Default: `false` (vibe-coder soft mode).
```

## Replace

`## 5. Dependencies`

## With

```md
## 5. Dependencies

- Declared in YAML frontmatter: `depends_on: [E##-epic/T###-task-id, ...]` (repo-relative IDs).
- Task dependency IDs are parsed with the same domain parser as task IDs (`src/domain/ids.ts`).
- `src/validation/dependencies.ts` detects duplicate task IDs, missing dependencies, and dependency cycles.
- `src/readers/tasks.ts` reads an epic task directory, validates all task markdown, then runs graph validation only after readable task files are parsed.
- Cross-epic deps allowed but warned (signal that epic boundaries may be wrong). Warning behavior belongs to future CLI/TUI command layers.
- TUI shows blocked tasks as visually locked.
- `savepoint doctor` detects cycles.
```

## Replace

`## 12. Distribution & build`

## With

```md
## 12. Distribution & build

- **License:** MIT.
- **Install:** primary `npx savepoint init`, persistent `npm i -g savepoint` → `savepoint`.
- **Runtime:** Node 20.10+ LTS, ESM-only, no native deps. macOS / Linux / Windows-Terminal.
- **Repo:** single package. TypeScript strict. `tsup` build → `dist/`. Bin `dist/cli.js` shebanged.
- **No telemetry.** Ever.
- **Baseline scaffold:** established in epic `E01-scaffolding` (2026-04-27). Package name `savepoint`, version `0.1.0`. Build, typecheck, lint, format, and test gates all pass.
- **Read-only data model:** established in epic `E02-data-model` (2026-04-27). `js-yaml` is the only runtime dependency added for structured YAML/frontmatter parsing.
```

## Replace

`## 13. Testing`

## With

```md
## 13. Testing

| Layer                                                    | Tool                             | Coverage                                                               |
| -------------------------------------------------------- | -------------------------------- | ---------------------------------------------------------------------- |
| Unit: file ops, YAML, frontmatter, snapshot gen          | `vitest`                         | High                                                                   |
| Unit: state transitions, dep resolution, cycle detection | `vitest`                         | High                                                                   |
| Integration: CLI commands in temp dirs                   | `vitest` + `tmp`                 | Medium                                                                 |
| TUI reducers (state, isolated from rendering)            | `vitest` + `ink-testing-library` | Medium                                                                 |
| TUI rendering (snapshot tests)                           | —                                | **None.** Brittle.                                                     |
| End-to-end with real AI agents                           | Manual matrix                    | Pre-release: `[Claude, Cursor, Gemini, Aider]` × `[init, plan, audit]` |

~70% line coverage target; behavior coverage prioritized.

E02 adds focused unit coverage for ID parsing, status validation, task/release/epic/router/config frontmatter validation, markdown read-boundary failures, project root discovery, dependency graph errors, and epic task-set reading. As of the E02 audit, `npm test` reports 16 passing test files and 176 passing tests when run outside the sandbox.
```
