---
id: E45-safe-migration/T011-run-migration-from-the-command-line
title: Run migration from the command line
status: planned
objective: Add a thin migrate command that previews by default, applies only when asked, and reports named errors.
depends_on:
    - E45-safe-migration/T009-publish-the-conversion-once
complexity_tier: low
complexity_reason: Argument parsing and dispatch over an existing package, following the established command pattern.
---

# T011: Run migration from the command line

## Problem

The whole operation is unreachable without a command, and the command is where the safest default has to be chosen once.

`savepoint migrate` should preview. A user who types a command they have not used before, on a project they care about, should get a report rather than a converted project — and the flag they have to add to change that is the moment they confirm intent. `v2-Design.md` spells the preview as `--dry-run`, which stays supported as an explicit way to say the default out loud, but the default itself is what protects the case that matters.

Everything else here is the existing pattern: `cmd/` parses arguments and dispatches, behavior lives in `internal/`, and `main.go` gains one case.

## Context Files

- `cmd/migrate.go`
- `cmd/migrate_test.go`
- `cmd/doctor.go`
- `cmd/upgrade-assets.go`
- `main.go`
- `main_test.go`
- `internal/migrate/plan.go`
- `internal/migrate/apply.go`
- `internal/migrate/operation.go`
- `README.md`

## Acceptance Criteria

- [ ] `savepoint migrate [dir]` previews and writes nothing; a test snapshots bytes and modification times across the command and asserts nothing changed.
- [ ] `--apply` is required to write; `--dry-run` is accepted as an explicit synonym of the default; passing both is accepted and previews; the help text states that preview is the default.
- [ ] `--decisions FILE` supplies the owner decision input, and a plan with unresolved blocking ambiguities exits nonzero naming every unresolved ID.
- [ ] `--recover` reports an incomplete operation and its recovery guidance, and resumes it when `--apply` is also given.
- [ ] A missing directory, an unwritable directory, and a directory that is not a Savepoint project each produce a distinct named error and a nonzero exit, with no partial write, satisfying FS-06.
- [ ] An unknown flag produces a named error and the usage text rather than being ignored.
- [ ] The preview output enumerates the planned records, archives, identity mappings, conflicts, and decisions required — the plan itself, not a summary that omits entries.
- [ ] Output is deterministic and readable without a TTY, matching the existing non-interactive command behavior.
- [ ] `cmd/migrate.go` contains argument parsing and dispatch only, with no frontmatter parsing, no conversion rule, and no filesystem walking, keeping ARCH-01 intact.
- [ ] `main.go` gains one `migrate` case following the existing dispatch shape, and `--version` and every existing command behave unchanged.
- [ ] `README.md` documents the command, its default, and the fact that preview writes nothing.

## Implementation Plan

- [ ] Add `cmd/migrate.go` following the argument-parsing and runner-injection shape used by `cmd/doctor.go` and `cmd/upgrade-assets.go`.
- [ ] Implement flag handling for `--apply`, `--dry-run`, `--decisions`, and `--recover`, with the preview default and named errors for unknown flags.
- [ ] Render the preview from the plan value, keeping the rendering in `internal/migrate` if it needs any domain knowledge.
- [ ] Wire the `migrate` case into `main.go`.
- [ ] Document the command in `README.md`.
- [ ] Test the write-free default, each flag, each target-directory failure, unknown flags, deterministic non-TTY output, and unchanged behavior of the existing commands.
- [ ] Run the focused `cmd` and `internal/migrate` suites, then `make build && make test`.

## Context Log

Pending.
