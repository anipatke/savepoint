---
id: T-056
title: Suggest health tools from what's already in the project
objective: O-028
status: planned
depends_on: [{task: T-055, requires: clear}]
owner_validation: {required: false}
planned_by: {role: planner, session: planning-o028-20261001}
---

# Suggest health tools from what's already in the project

## Outcome

A read-only `Discover` service inspects a bounded set of project files and PATH and returns confirmation-ready proposals: one per capability and component, each with a plain-language reason, the exact configuration it would write, default exclusions, and any support gap. It never executes, installs, writes, or uses the network.

## User Check

None beyond the Full Objective Check; the owner sees the proposals through `savepoint health setup` in the setup Task.

## Done When

- `Discover(ctx, root, lookPath)` returns a deterministic `Proposal` list. Each proposal names capability, provider, instance name, scope, exclusions, executable (for executed tools), args, report path, reason, and a gap (`missing_executable`, `unsupported_stack`, `no_report_yet`, or none). Every proposed `CapabilityConfig` passes `Validate()`.
- **Detected inputs:** `go.mod`; `package.json` scripts and dev dependencies (Vitest, `@vitest/coverage-v8`), `vitest.config.*`; `pyproject.toml`, `pytest.ini`, `setup.cfg`, `.coveragerc`; npm/pnpm/Yarn/Bun, Go, and Python lockfiles named in the O-026 catalogue; `.jscpd.json`, `lizard` config; existing reports at the catalogue's conventional paths. Unrecognized stacks produce an explicit `unsupported_stack` gap, never a silent omission.
- **Executed tools** (Lizard, jscpd, OSV-Scanner): proposal args include the report output and the confirmed exclusions translated into that tool's own exclude flags; OSV-Scanner uses its normal online mode, and its proposal reason says the scanner sends package names and versions to OSV.dev. A tool not on PATH is still proposed, with `missing_executable`, and never downloaded.
- **Report-only providers** (tests and coverage): proposals point at the report path the project's gate should produce and include the exact flag the owner would add to their own gate (for example `go test -json ... > <path>`); Savepoint does not edit their scripts.
- **Default exclusions:** one catalogue list (`vendor/**`, `node_modules/**`, `third_party/**`, `dist/**`, `build/**`, `**/*.pb.go`, `**/*_generated.*`, `**/*.min.js`) proposed only when the path exists or the pattern is a generated-file suffix; the list lives in data, not scattered in logic.
- **Monorepos:** manifests are searched to a fixed depth (3) skipping excluded and hidden directories; each component yields its own named, scoped instance. No aggregate instance is proposed.
- **Bounds:** at most a fixed number of files read, each capped in size; a larger or unreadable file becomes a named gap, not a crash. Symlinks outside the project are not followed.
- Tests prove no process is started (a `lookPath` fake is the only outside call) and no file is written (tree bytes and mtimes unchanged).

## Context Files

`.savepoint/objectives/O-028-code-health-provider-execution-and-scaffolding/Objective.md`; `.savepoint/objectives/O-026-code-health-provider-feasibility/Objective.md` (catalogue sections only); `internal/codehealth/model.go`; `internal/codehealth/config.go`; `internal/codehealth/primitives.go`; new `internal/codehealth/discovery.go`, `internal/codehealth/discovery_catalogue.go`, `internal/codehealth/discovery_test.go`, and fixtures under `internal/codehealth/testdata/discovery/`.

## Design References

Design sections 1 and 11; O-026 Selected Provider Catalogue and Module Boundary; O-028 Confirmed Design Decisions (2026-10-01).

## Guardrails

FS-03, FS-05, ARCH-03, ARCH-04, CFG-02, CFG-03, DEP-01, DEP-02, STYLE-09, TEST-01, TEST-02, TEST-04, TEST-08.

## Implementation Plan

1. Confirm the first Task's `Name`/`Required` fields and default-timeout lookup exist; return REPLAN REQUIRED if not.
2. Write the catalogue data (detected files, report conventions, per-tool argument templates and exclusion-flag translation, default exclusions) in `discovery_catalogue.go`.
3. Implement bounded component search and per-stack detection in `discovery.go`; parse only the fields needed (`package.json` via `encoding/json`; TOML/INI files by bounded line scan for known keys; no new dependency).
4. Build proposals, validate each as a `CapabilityConfig`, and sort deterministically (component path, capability order, provider).
5. Fixtures: Go-only, Vitest with V8 coverage, pytest with coverage.py, a polyglot monorepo with two components, a repo with vendored and generated code, an unsupported stack, a missing executable, an oversized manifest, and an escaping symlink.

## Boundaries

No execution, no writes, no network, no install hints beyond naming the missing executable, no CLI, no Design.md edits, no report parsing.

## Technical Verification

Focused runs while iterating; `make build && make test-fast` at handoff. Discovery is filesystem- and platform-sensitive, so also run fresh `make test-full` and note Windows path handling.

## Technical Evidence

Pending execution.

## Drift Notes

Update the AGENTS.md Codebase Map entry for `internal/codehealth`: it currently says the package "performs no collection"; discovery is now in scope. Reconcile through the Full Objective Check.
