# Audit Proposals: E02-data-readers

## Target File

`.savepoint/Design.md`

## Replace

```md
- **Template asset boundary:** established in epic `E04-templates-and-prompts` (2026-04-27). Project scaffolding and workflow prompts now live as versioned markdown/YAML assets under `templates/`, with small TypeScript helpers in `src/templates/` for manifest lookup, path resolution, loading, and interpolation.
```

## With

```md
- **Go data-reader boundary:** established in epic `E02-data-readers` (2026-05-01). `internal/data` owns Savepoint file parsing and discovery for the Go implementation: task frontmatter models, markdown YAML extraction, router state parsing, config theme defaults, and release/epic/task directory listing.
- **Template asset boundary:** established in epic `E04-templates-and-prompts` (2026-04-27). Project scaffolding and workflow prompts now live as versioned markdown/YAML assets under `templates/`, with small TypeScript helpers in `src/templates/` for manifest lookup, path resolution, loading, and interpolation.
```

## Target File

`.savepoint/Design.md`

## Replace

```md
## 12. Distribution & build
```

## With

```md
## 12. Distribution & build

> Audit note: the live repository is transitioning from the documented TypeScript/Node implementation to a Go module (`github.com/opencode/savepoint`). The architecture document still contains substantial TypeScript-era implementation detail and should be reconciled as Go epics are audited.
```

## Target File

`AGENTS.md`

## Replace

```md
| `src/domain/config.ts`               | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/E02-Detail.md)                       | Config defaults and typed project config model                                                       |
| `src/domain/epic.ts`                 | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/E02-Detail.md)                       | Epic Design frontmatter model and validation                                                         |
| `src/domain/ids.ts`                  | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/E02-Detail.md)                       | Release, epic, and task ID parsing and formatting                                                    |
| `src/domain/release.ts`              | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/E02-Detail.md)                       | Release PRD frontmatter model and validation                                                         |
| `src/domain/router.ts`               | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/E02-Detail.md)                       | Router state values and validation                                                                   |
| `src/domain/status.ts`               | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/E02-Detail.md)                       | Task status values and transition validation                                                         |
| `src/domain/task.ts`                 | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/E02-Detail.md)                       | Task frontmatter/document model and validation                                                       |
| `src/fs/markdown.ts`                 | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/E02-Detail.md)                       | Markdown/frontmatter reader with path-aware boundary errors                                          |
| `src/fs/project.ts`                  | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/E02-Detail.md)                       | Project root discovery and scoped `.savepoint` path helpers                                          |
| `src/readers/config.ts`              | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/E02-Detail.md)                       | Read-only `config.yml` reader with defaults                                                          |
| `src/readers/epic.ts`                | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/E02-Detail.md)                       | Read-only epic Design reader                                                                         |
| `src/readers/release.ts`             | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/E02-Detail.md)                       | Read-only release PRD reader                                                                         |
| `src/readers/router.ts`              | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/E02-Detail.md)                       | Read-only router Current state reader                                                                |
| `src/readers/tasks.ts`               | [E02-data-model](.savepoint/releases/v1/epics/E02-data-model/E02-Detail.md)                       | Epic task-set reader with graph validation                                                           |
```

## With

```md
| `internal/data/task.go`              | [E02-data-readers](.savepoint/releases/v1/epics/E02-data-readers/E02-Detail.md)                   | Task, column, progress stage, dependency, and metadata structs for Go data readers                   |
| `internal/data/parser.go`            | [E02-data-readers](.savepoint/releases/v1/epics/E02-data-readers/E02-Detail.md)                   | Markdown YAML frontmatter extraction and task frontmatter parsing                                    |
| `internal/data/router.go`            | [E02-data-readers](.savepoint/releases/v1/epics/E02-data-readers/E02-Detail.md)                   | Router current-state YAML block extraction and parsing                                               |
| `internal/data/config.go`            | [E02-data-readers](.savepoint/releases/v1/epics/E02-data-readers/E02-Detail.md)                   | Savepoint config/theme parsing with default theme fallback                                           |
| `internal/data/discover.go`          | [E02-data-readers](.savepoint/releases/v1/epics/E02-data-readers/E02-Detail.md)                   | `.savepoint` root discovery and release, epic, and task directory listing                            |
```

## Target File

`.savepoint/releases/v1/epics/E02-data-readers/E02-Detail.md`

## Insert After

```md
## Components and files
```

## With

```md
## Implemented As

- `internal/data/task.go` defines task metadata, column, progress stage, dependency, and progress structs.
- `internal/data/parser.go` extracts YAML frontmatter from markdown and currently maps the task ID and objective into `Task`.
- `internal/data/router.go` extracts and parses the `## Current state` fenced YAML block from `router.md`.
- `internal/data/config.go` reads `config.yml` theme values and falls back to a default theme when the config is missing.
- `internal/data/discover.go` finds the `.savepoint` root and lists release, epic, and task directories in stable order.

## Deltas From Plan

- The planned `internal/data/errors.go` was not created; errors currently use formatted `error` values.
- `ParseTaskFile` does not yet populate `Column`, `Stage`, `DependsOn`, `Priority`, `Points`, `Tags`, `Acceptance`, or `Notes` from YAML frontmatter.
- The parser uses a regular expression rather than a line-based frontmatter scanner.
- Config malformed YAML currently falls back to defaults instead of surfacing a parse error.
```

## Quality Review

## Must Fix Before Close

All close blockers resolved in audit closeout.

## Must Fix Before Next Epic

All next-epic blockers resolved in audit closeout.

## Carry Forward

- `.savepoint/Design.md` and `AGENTS.md` still mostly describe the previous TypeScript implementation. The Go architecture should be reconciled incrementally as each Go epic is audited.

## Already Fixed

- `go test ./internal/data/...` passed during audit.
- `go test ./...` passed during audit.
- `ParseTaskFile` now maps task frontmatter fields including `status`, `phase`, `depends_on`, priority, points, tags, acceptance, notes, release, epic, and progress.
- `ConfigReader.Read` now returns defaults only for missing config files and returns malformed YAML/read errors.
- Discovery tests now use temporary fixtures instead of the live repository tree.
- Parser tests now assert task metadata mapping and malformed YAML behavior.
- `internal/data/errors.go` now provides data-reader boundary error sentinels.
- E02 task files now include context logs.
