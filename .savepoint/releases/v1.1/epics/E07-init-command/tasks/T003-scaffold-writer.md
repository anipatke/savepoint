---
id: E07-init-command/T003-scaffold-writer
status: done
objective: "Implement scaffold writing from templates"
depends_on: ["E07-init-command/T002-target-validation"]
---

# T003: Scaffold Writer

## Acceptance Criteria

- Reads templates from `templates/project/` directory
- Interpolates project name into template files (AGENTS.md, config.yml)
- Creates parent directories as needed
- Excludes conflicting files unless --force
- Copies all .savepoint/ contents
- Creates root AGENTS.md from template

## Implementation Plan

- [x] Add `internal/init/scaffold.go`
- [x] Implement `Scaffold(targetDir, projectName, force) error`
- [x] Load templates from embedded or filesystem path
- [x] Implement template variable interpolation (project name)
- [x] Create .savepoint/ directory structure
- [x] Copy all template files to target
- [x] Create root AGENTS.md from template
- [x] Test scaffold output in temp directory
- [x] Run `make build && make test`

## Context Log

- Files read: router.md, E07-Detail.md, T002-target-validation.md, T001-cli-entrypoint.md, cmd/init.go, cmd/init_test.go, main.go, internal/init/validate.go, internal/init/validate_test.go, Makefile, templates/project/**/*
- Files created: internal/init/scaffold.go, internal/init/scaffold_test.go
- Files modified: main.go, T003-scaffold-writer.md
- Estimated input tokens: ~14k
- Key decisions:
  - Used Go `embed.FS` with `fs.Sub` to embed `templates/project/` recursively.
  - Explicit `//go:embed templates/project/.savepoint` needed because Go's embed excludes dotfile-prefixed directories from recursive directory embedding.
  - Scaffold accepts `fs.FS` for testability; tested with `testing/fstest.MapFS`.
  - Init runner in `main.go` validates first (reuses T002), then scaffolds.
  - Template variables: `{{PROJECT_NAME}}` and `{{RELEASE_NUMBER}}` replaced via `strings.ReplaceAll`.

## Drift Notes

- `cmd/` and `internal/init/` modules were not in AGENTS.md Codebase Map. Updated AGENTS.md to include both.