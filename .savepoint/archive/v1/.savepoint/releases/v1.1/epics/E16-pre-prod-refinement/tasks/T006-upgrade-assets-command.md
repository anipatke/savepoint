---
id: E16-pre-prod-refinement/T006-upgrade-assets-command
status: done
objective: Add a safe project asset updater so npm package updates can refresh Savepoint templates and agent skills without overwriting user-authored project state
depends_on: []
---

# T006: Upgrade Project Assets Command

## Problem

Savepoint's scaffolded templates and agent skills are embedded in the npm-distributed binary and will continue to improve as part of package releases. Today, `savepoint init --force` can refresh scaffolded files, but it is too broad for routine package updates because it can rewrite project state files that contain user-authored planning and release work.

Users need a safe, explicit command they can run after `npm install` or `npm update` to refresh package-owned assets in an existing Savepoint project while preserving PRDs, design docs, releases, epics, tasks, audits, and other project state.

## Context Files

- `main.go` - CLI dispatch, embedded template filesystem, init runner pattern
- `cmd/init.go` - command parsing style to mirror for upgrade-assets
- `internal/init/scaffold.go` - current template walking and AGENTS.md merge behavior
- `internal/init/agents.go` - managed AGENTS.md block merge helpers
- `internal/init/validate.go` - target validation patterns and conflict behavior
- `internal/init/write.go` - atomic write helper for safe file updates
- `internal/init/scaffold_test.go` - scaffold overwrite and merge test patterns
- `internal/init/integration_test.go` - end-to-end init pipeline patterns
- `package.json` - npm package scripts and shipped files
- `README.md` - install/update documentation

## Acceptance Criteria

- [ ] A new command `savepoint upgrade-assets [dir] [--dry-run] [--force]` is available
- [ ] The command refuses to run against a directory that is not an existing Savepoint project
- [ ] The command refreshes package-owned `agent-skills/**/SKILL.md` files from the embedded `templates/project` assets
- [ ] The command refreshes the Savepoint-managed block in an existing agent guide through the existing AGENTS.md merge behavior
- [ ] The command does not overwrite `.savepoint/PRD.md`, `.savepoint/Design.md`, `.savepoint/releases/**`, task files, epic files, or audit files
- [ ] `--dry-run` reports planned updates, merges, skips, and unchanged files without writing to disk
- [ ] Re-running the command is idempotent and reports unchanged assets when no content differs
- [ ] `--force` is reserved for explicitly package-owned asset refresh behavior and does not delete unknown local files by default
- [ ] npm install/update behavior includes a lightweight postinstall notice that tells users to run `savepoint upgrade-assets` in each project, without mutating any project automatically
- [ ] README documents the update flow: `npm update -g savepoint` followed by `savepoint upgrade-assets`
- [ ] `make build && make test` passes

## Implementation Plan

- [ ] Add command parsing for `upgrade-assets [dir] [--dry-run] [--force]` in `cmd/`
- [ ] Wire `upgrade-assets` into `main.go` using the embedded `templates/project` filesystem
- [ ] Add an internal upgrade-assets runner that validates the target contains `.savepoint`
- [ ] Implement an allowlist so only `agent-skills/**/SKILL.md` and the managed `AGENTS.md` block are refreshed initially
- [ ] Add a structured report for updated, merged, unchanged, skipped, and refused paths
- [ ] Implement dry-run mode using the same diff decisions as write mode
- [ ] Preserve unknown local files and all project state files regardless of `--force`
- [ ] Add focused unit tests for parsing, validation, writes, dry-run, idempotency, and skipped project-state files
- [ ] Add a package `postinstall` script that prints the explicit upgrade command only
- [ ] Update README install/update documentation
- [ ] Run `make build && make test`

## Context Log

- Files read: main.go, cmd/init.go, cmd/board.go, cmd/doctor.go, internal/init/scaffold.go, internal/init/agents.go, internal/init/validate.go, internal/init/write.go, internal/data/discover.go, templates/project, Makefile, package.json, README.md, related test files
- Files edited: main.go, cmd/upgrade-assets.go (new), cmd/upgrade-assets_test.go (new), internal/init/upgrade.go (new), internal/init/upgrade_test.go (new), package.json, README.md, this file
- Token estimate: ~12K
- Quality gates: make build (pass), make test (pass — all 9 packages green)

## Drift Notes

Added new source files not yet listed in AGENTS.md Codebase Map:
- `cmd/upgrade-assets.go` — CLI arg parsing for upgrade-assets command
- `internal/init/upgrade.go` — upgrade runner: project validation, allowlist walks, structured report

Corresponding tests were also added at `cmd/upgrade-assets_test.go` and `internal/init/upgrade_test.go`.

These fall under the existing `cmd/` and `internal/init/` modules in the Codebase Map, but the specific upgrade-assets functionality should be documented there.
