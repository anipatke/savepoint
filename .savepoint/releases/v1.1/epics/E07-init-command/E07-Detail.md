---
type: epic-design
status: audited
---

# E07: Init Command

## Purpose

Implement `savepoint init`, the first-run command that creates a Savepoint workflow in an empty or compatible project directory and prints the magic prompt the user gives to an AI agent.

## Interface

```bash
savepoint init [dir]           # scaffold .savepoint/ in target directory
savepoint init [dir] --force   # overwrite existing
savepoint init [dir] --install  # run npm install after scaffold
```

## What this epic adds

- Target directory validation (missing, empty, compatible, already-initialized, conflicting, boundary-error)
- Scaffold writing from `templates/project/`
- Project name interpolation
- Safe handling (refuse conflicts by default, temp-file-plus-rename)
- Magic prompt output to stdout (from `templates/prompts/magic-prompt.prompt.md`)
- Best-effort clipboard copy (doesn't fail if unavailable)
- Optional `--install` flag for dev dependencies
- Integration tests using temporary directories

## Components

| Module | Purpose |
|--------|---------|
| `cmd/init.go` | CLI registration, arg parsing |
| `internal/init/validate.go` | Target directory checks |
| `internal/init/scaffold.go` | Template copy, project name interpolation |
| `internal/init/write.go` | Atomic writes with temp-file-plus-rename |
| `internal/init/prompt.go` | Magic prompt rendering from template |
| `internal/init/clipboard.go` | Platform clipboard detection and copy |

## Implemented As

- `cmd/init.go` owns init argument parsing and delegates execution through `InitRunner`.
- `main.go` embeds `templates/project/` and `templates/prompts/`, then wires validation, scaffold, prompt rendering, clipboard copy, and optional install in sequence.
- `internal/init/` owns validation, scaffold interpolation, atomic writes, prompt rendering, clipboard copy, dependency install, and integration tests.
- Audit reconciled generated workflow templates with the current epic-local audit workflow before closing.

## Boundaries

**In scope:**
- Create `.savepoint/` and root `AGENTS.md`
- Print the initial agent prompt
- Avoid overwriting user files unless explicitly allowed by flags
- Support Windows, macOS, and Linux filesystem behavior

**Out of scope:**
- Creating epics/tasks interactively
- Running the board after init
- Implementing audit or doctor
- Managing package installation beyond the optional one-shot path
