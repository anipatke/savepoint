## Target File

`.savepoint/releases/v1/epics/E02-data-model/Design.md`

## Replace

```yaml
---
type: epic-design
status: active
---
```

## With

```yaml
---
type: epic-design
status: audited
audited: 2026-04-27
---
```

## Replace

`## Close Criteria`

## With

```md
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
```
