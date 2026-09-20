---
id: E45-safe-migration/T003-plan-the-conversion-before-touching-anything
title: Plan the conversion before touching anything
status: done
objective: Produce a deterministic write-free conversion plan with global identity allocation and the source-qualified reference map.
depends_on:
    - E45-safe-migration/T002-take-stock-of-what-exists
complexity_tier: high
complexity_reason: Deterministic target planning, global identity allocation, and the reference map every later task depends on.
---

# T003: Plan the conversion before touching anything

## Problem

A user has to be able to read what migration intends to do, on a project that is still untouched, and decide whether to allow it. That requires the plan to be a value rather than a side effect: computed from the inventory, reviewable in full, and identical the next time it runs.

Identity is where the determinism gets hard. A V1 `T001` is not a V2 `T001`, and the `v1-history` fixture proves why: the same short ID `T001-shared` exists under `E01-example` in both release `v1` and release `v1.1`. Allocation has to be keyed by a source-qualified identity, walk sources in a fixed order, and never reuse an identity — including never reserving one for work that is only archived.

The reference map is the other half. Some V1 facts have no V2 record to live in, and the manifest is where they land instead of being lost or faked.

## Context Files

- `internal/migrate/plan.go`
- `internal/migrate/plan_test.go`
- `internal/migrate/manifest.go`
- `internal/migrate/manifest_test.go`
- `internal/migrate/inventory.go`
- `internal/migrate/classify.go`
- `internal/data/config.go`
- `internal/data/objective_v2.go`
- `internal/data/task_v2.go`
- `internal/data/issue_v2.go`
- `internal/init/migrate_audit_skill.go`
- `internal/data/testdata/migration/v1-history/manifest.yml`
- `.savepoint/releases/v2/v2-Design.md`

## Acceptance Criteria

- [x] `Plan` is a pure function of the inventory, the owner decisions, an injected clock, and an injected operation-ID source; it opens no file for writing, creates no directory, and a byte-and-mtime snapshot across a full plan of both fixtures proves it.
- [x] Planning the same fixture twice with the same injected clock and operation ID produces an identical plan: identical targets, identical allocated IDs, identical ordering.
- [x] Legacy identities are source-qualified by release, epic, path, and original ID; `v1-history`'s two `E01-example/T001-shared` records are allocated two distinct global Task IDs, and the map records which is which.
- [x] Global IDs are allocated as the next unused `O`, `T`, or `I` value in a fixed source order, and no identity is ever reused.
- [x] Completed Tasks and audited epics are planned as archive entries with a mapping record and no V2 identity; they reserve nothing, because they receive no ID at all.
- [x] A dependency from an active Task onto an archived completed Task is planned as a typed `legacy_prerequisite` manifest entry keyed by the new `T###`, naming the archive path and the original recorded completion or waiver evidence; it never appears in the converted Task's `depends_on`.
- [x] The plan enumerates, in full, every target record, every archive entry, the complete identity map, every conflict, and every ambiguity, so the preview a user reads is the plan itself rather than a summary of it.
- [x] The manifest model serializes the source hashes, the identity map, archive references, typed legacy prerequisites, and recorded owner decisions to `.savepoint/migrations/v1-to-v2.yml`.
- [x] The manifest is create-only and coexists with the directory's existing occupants: `internal/init`'s archived legacy skill copies and its `README.md` are inventoried, planned for preservation, and never rewritten.
- [x] An existing `v1-to-v2.yml` is refused with a named diagnostic distinguishing an already-migrated project from a conflicting one, rather than being overwritten.
- [x] A project already at `schema_version: 2` plans no work and reports that; an unknown explicit schema version fails through `ReadSchemaVersion`'s existing named diagnostic rather than a second one.
- [x] Every planned target path is confined to the project and built with `filepath` joins.

## Implementation Plan

- [x] Define the plan value type in `plan.go`: planned records, archive entries, identity map, legacy prerequisites, conflicts, and ambiguities, all as data rather than as rendered text.
- [x] Implement source-qualified legacy keys and fixed-order identity allocation, with the archived-record rule that allocates nothing.
- [x] Implement the active-to-archived dependency rule that produces a typed legacy prerequisite instead of a `depends_on` entry.
- [x] Add the version gate over `data.ReadSchemaVersion`, reusing its diagnostics for unknown versions.
- [x] Add `manifest.go`: the `v1-to-v2.yml` model, its serialization, and the create-only and coexistence rules against `internal/init`'s existing `.savepoint/migrations/` contents.
- [x] Inject the clock and operation-ID source through the plan's inputs so tests can fix both and later tasks can render goldens.
- [x] Test determinism across repeated runs, source-qualified allocation on `v1-history`, archived-record handling, legacy prerequisites, manifest coexistence and refusal, both version-gate outcomes, and the write-free snapshot.
- [x] Run the focused `internal/migrate` suite, then `make build && make test`.

## Context Log

**Read:** `.savepoint/router.md`, `E45-Detail.md`, this task file, all Context Files listed above, plus (to ground the design against real V1 shapes) `internal/data/task.go`, `defect.go`, `dependency.go`, `discover.go`, `parser.go`, `lifecycle.go`, `audit_finding.go`, `audit_run.go`, `audit_register.go`, `evidence_v2.go`, `project.go`, `errors.go`, the frozen `v1-basic`/`v1-history` fixture files and their `manifest.yml`s, `internal/migrate/fixture_test.go` and `inventory_test.go` (for reusable test helpers), and `.savepoint/Guardrails.md`.

**Wrote:** `internal/migrate/plan.go` (new), `internal/migrate/manifest.go` (new), `internal/migrate/plan_test.go` (new), `internal/migrate/manifest_test.go` (new).

**Design notes / scope boundary:** `Plan` decides identity and destination only — which V1 source becomes which V2 Objective/Task/Issue, or an archive entry, plus the legacy reference map. It deliberately does not yet assign a V2 fate to `config.yml`, `router.md`, `PRD.md`, `Design.md`, `Health-Check.md`, the managed guide, or skill files — those require content-rendering decisions owned by T004-T006. Immutable byte-preserve-only history (epic audits, the audit prompt/register/runs, resolved/verified/waived dispositions, and unclassified files) is archived here since archiving never requires a content decision. `AmbiguityDuplicateSourceIdentity` and `AmbiguityUnresolvedNarrativeFind` are defined in the shared vocabulary per the epic but not yet detected — that full decision-file-driven contract is T007's scope; this task detects and blocks on unrecognized lifecycle values and missing dependency targets, since both fall directly out of identity allocation and dependency resolution.

**Quality gates:**
- `go build ./...` — clean.
- `go vet ./...` — clean.
- `gofmt -l` over changed files — clean.
- `make build && make test` — all packages pass, including the new `internal/migrate` plan/manifest tests (`go test ./internal/migrate/...`).
- No `.savepoint/Health-Check.md` exists at the project root; the Quick check step is skipped per `savepoint-build-task` (absence is not a finding).
