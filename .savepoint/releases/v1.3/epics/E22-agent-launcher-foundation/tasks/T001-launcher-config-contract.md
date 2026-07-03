---
id: E22-agent-launcher-foundation/T001-launcher-config-contract
status: in_progress
stage: build
objective: Add a disabled-by-default, vendor-neutral launcher configuration contract.
depends_on: []
complexity_tier: medium
complexity_reason: The config adds nested defaults and validation while preserving all existing config behavior.
---

# T001: Launcher Config Contract

## Problem

Savepoint has no configuration for enabling agent launches, assigning builder and auditor roles, or selecting terminal behavior.

## Context Files

- `internal/data/config.go`
- `internal/data/config_test.go`
- `.savepoint/config.yml`

## Acceptance Criteria

- [x] `agent_launcher.enabled` defaults to `false` when absent.
- [x] Builder and auditor profiles support an executable plus an ordered argument list with documented placeholders.
- [x] Auditor configuration is optional; builder configuration is required only when the launcher is enabled.
- [x] Terminal mode supports `auto` plus a structured override without requiring shell command strings.
- [x] Existing config files parse with unchanged theme and quality-gate behavior.
- [x] Validation errors name the exact missing or invalid launcher field.

## Implementation Plan

- [x] Add typed launcher, agent-profile, and terminal configuration fields to `internal/data`.
- [x] Define defaults and validation helpers separately from YAML parsing.
- [x] Keep command arguments as structured values and define supported placeholders centrally.
- [x] Add table-driven tests for absent, disabled, enabled, partial, and malformed configurations.
- [x] Add a disabled example block to the repository-local config for development.

## Context Log

**Files read:** `internal/data/config.go`, `internal/data/config_test.go`, `.savepoint/config.yml`, `E22-Detail.md`.

**Files edited/added:**
- `internal/data/launcher_config.go` (new) — `AgentLauncher`, `CommandSpec`, `TerminalConfig`, `TerminalMode` (`auto`/`override`), centrally-defined placeholders (`{{prompt}}`, `{{project_root}}`) with `SupportedPlaceholders`, plus `fillLauncherDefaults` and `Validate` kept separate from YAML parsing.
- `internal/data/config.go` — added `AgentLauncher` field to `Config`, baked the auto terminal default into `defaultConfig`, and called fill + `Validate` in `Read` (returns a field-naming error on invalid launcher config).
- `internal/data/launcher_config_test.go` (new) — table-driven coverage for absent/missing-file, disabled, enabled-full, partial (optional auditor omitted), and malformed/invalid (missing builder, missing command, invalid mode, override without command), plus existing theme/quality-gate behavior and placeholder set.
- `.savepoint/config.yml` — added a commented, disabled `agent_launcher` example block.

**Acceptance criteria:** each verified by a passing test (defaults, optional auditor, terminal auto/override, unchanged theme/gates, and field-named validation errors).

**Quality gates:**
- `go test ./internal/data` → ok
- `make build` → ok
- `make test` (`go test ./...`) → ok (all packages)
- `gofmt`/`go vet` clean for the changed files. The epic gate `go test ./internal/launcher` is deferred: that package is created in T004.

## Drift Notes

- The epic's Components table named `internal/data/config.go` for launcher parsing. Per the AGENTS.md "one job per file" rule, the launcher contract lives in a new `internal/data/launcher_config.go` (same `data` package), keeping `config.go` focused on theme/quality-gate parsing. No new module or architectural change versus the Codebase Map.
