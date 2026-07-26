---
id: E40-upgrade-safety/T004-backward-compat-fixtures
status: planned
objective: Pin old file shapes and the template/skill lifecycle matrix with regression tests and explicit one-command upgrade guidance.
depends_on: ["E40-upgrade-safety/T002-skill-conflict-policy", "E40-upgrade-safety/T003-agents-marker-conflict"]
complexity_tier: medium
complexity_reason: Fixture authoring, reader-contract assertions, full install/upgrade ownership coverage, and user-facing upgrade guidance span multiple boundaries.
---

# T004: Backward-Compatibility Fixtures

## Problem

`upgrade.go:118-121` skips the whole `.savepoint/` subtree, so `router.md` is never migrated and a project's router keeps its original shape indefinitely. Compatibility therefore has to hold in the reader, not the template — and nothing tests that.

`ReadState` (`internal/data/router.go:28`) currently tolerates unknown and missing YAML keys, because `yaml.Unmarshal` is called without `KnownFields`. That tolerance is accidental rather than asserted: a future change enabling strict decoding, or adding a required field, would break every existing router with no failing test.

It also hard-fails without the literal `## Current state` heading (`router.go:10`) and a fenced YAML block inside it. Those two anchors are a format contract, and nothing records that.

`template_freshness_test.go` asserts templates match *current* guidance. Nothing asserts that *old* files still work, that every installed asset has an explicit upgrade ownership policy, or that user-facing instructions distinguish the required upgrade command from its optional dry-run preview.

## Context Files

- `internal/init/testdata/legacy/` (will create)
- `internal/data/router_test.go`
- `internal/init/upgrade_test.go`
- `internal/init/integration_test.go`
- `internal/data/router.go` (contract under test)
- `README.md`

## Acceptance Criteria

- [ ] `internal/init/testdata/legacy/` holds pinned fixtures: a pre-v1.5 `router.md`, an `AGENTS.md` with no markers, an `AGENTS.md` with markers, and a customized `SKILL.md`.
- [ ] Fixtures are byte-frozen and commented as such — they represent files already in the wild and must never be "fixed up" to match current templates. A test that fails against a fixture indicates a compatibility break, not a stale fixture.
- [ ] A router test asserts the legacy fixture parses, and that every `RouterState` field it carries reads back correctly.
- [ ] A router test asserts unknown YAML keys are ignored rather than erroring, pinning the non-strict decode.
- [ ] A router test asserts a missing optional key (`defect`) yields the zero value rather than an error.
- [ ] A router test asserts the two structural anchors: content without `## Current state`, and content without a fenced YAML block, each produce a clear error. The shipped `templates/project/.savepoint/router.md` is asserted to contain both.
- [ ] An upgrade test runs `UpgradeProjectAssets` over a project built from the legacy fixtures and asserts: the marker-less `AGENTS.md` is reported `conflict` and left byte-identical, the marked one is merged in place, and the customized `SKILL.md` survives per the T002 migration rule with its `.bak` written.
- [ ] A test asserts a task file using the legacy `phase` frontmatter field still parses, covering the existing `phase` → `stage` compatibility path (`parser.go:82`).
- [ ] Install and upgrade tests enforce the E40 Template and skill lifecycle matrix: all shipped assets exist after `init`; only package-owned skills/references and the marked AGENTS.md block refresh in place; Guardrails.md, Health-Check.md, and audit assets install only when missing; every other `.savepoint/` file remains untouched on upgrade.
- [ ] The permanent project-owned boundary is explicit and tested: `--force` does not overwrite router, config, PRD, Design, Concept, visual identity, Guardrails, Health-Check, audit state, release PRDs, epics, tasks, defects, or audit records.
- [ ] README documents `savepoint upgrade-assets` as the single required per-project update command, describes `--dry-run` as an optional read-only preview, and makes clear that installing or updating the Savepoint binary does not mutate existing projects automatically.
- [ ] `make build && make test` passes.

## Implementation Plan

- [ ] Author the fixture files, each with a header comment stating its provenance and that it is frozen.
- [ ] Add the router reader-contract tests to `internal/data/router_test.go`.
- [ ] Add the legacy-project upgrade test to `internal/init/upgrade_test.go`, assembling a temp project from the fixtures.
- [ ] Confirm the legacy `phase` task case is covered; add it if the existing parser tests only cover `stage`.
- [ ] Add table-driven install/upgrade ownership coverage for every lifecycle-matrix row, including `--force` assertions at the project-owned boundary.
- [ ] Update README with the one-command upgrade flow, optional dry run, and conflict-resolution behavior.

## Notes

This is the task that makes the rest of the epic durable. T001–T003 fix three specific defects; this one is the guard that catches the fourth before it ships.

The fixtures encode a policy as much as a test: Savepoint may change its templates freely, but may not change how it *reads* files it has already written without a deliberate migration.
