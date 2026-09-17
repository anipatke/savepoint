---
id: E45-safe-migration/T003-plan-the-conversion-before-touching-anything
title: Plan the conversion before touching anything
status: planned
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

- [ ] `Plan` is a pure function of the inventory, the owner decisions, an injected clock, and an injected operation-ID source; it opens no file for writing, creates no directory, and a byte-and-mtime snapshot across a full plan of both fixtures proves it.
- [ ] Planning the same fixture twice with the same injected clock and operation ID produces an identical plan: identical targets, identical allocated IDs, identical ordering.
- [ ] Legacy identities are source-qualified by release, epic, path, and original ID; `v1-history`'s two `E01-example/T001-shared` records are allocated two distinct global Task IDs, and the map records which is which.
- [ ] Global IDs are allocated as the next unused `O`, `T`, or `I` value in a fixed source order, and no identity is ever reused.
- [ ] Completed Tasks and audited epics are planned as archive entries with a mapping record and no V2 identity; they reserve nothing, because they receive no ID at all.
- [ ] A dependency from an active Task onto an archived completed Task is planned as a typed `legacy_prerequisite` manifest entry keyed by the new `T###`, naming the archive path and the original recorded completion or waiver evidence; it never appears in the converted Task's `depends_on`.
- [ ] The plan enumerates, in full, every target record, every archive entry, the complete identity map, every conflict, and every ambiguity, so the preview a user reads is the plan itself rather than a summary of it.
- [ ] The manifest model serializes the source hashes, the identity map, archive references, typed legacy prerequisites, and recorded owner decisions to `.savepoint/migrations/v1-to-v2.yml`.
- [ ] The manifest is create-only and coexists with the directory's existing occupants: `internal/init`'s archived legacy skill copies and its `README.md` are inventoried, planned for preservation, and never rewritten.
- [ ] An existing `v1-to-v2.yml` is refused with a named diagnostic distinguishing an already-migrated project from a conflicting one, rather than being overwritten.
- [ ] A project already at `schema_version: 2` plans no work and reports that; an unknown explicit schema version fails through `ReadSchemaVersion`'s existing named diagnostic rather than a second one.
- [ ] Every planned target path is confined to the project and built with `filepath` joins.

## Implementation Plan

- [ ] Define the plan value type in `plan.go`: planned records, archive entries, identity map, legacy prerequisites, conflicts, and ambiguities, all as data rather than as rendered text.
- [ ] Implement source-qualified legacy keys and fixed-order identity allocation, with the archived-record rule that allocates nothing.
- [ ] Implement the active-to-archived dependency rule that produces a typed legacy prerequisite instead of a `depends_on` entry.
- [ ] Add the version gate over `data.ReadSchemaVersion`, reusing its diagnostics for unknown versions.
- [ ] Add `manifest.go`: the `v1-to-v2.yml` model, its serialization, and the create-only and coexistence rules against `internal/init`'s existing `.savepoint/migrations/` contents.
- [ ] Inject the clock and operation-ID source through the plan's inputs so tests can fix both and later tasks can render goldens.
- [ ] Test determinism across repeated runs, source-qualified allocation on `v1-history`, archived-record handling, legacy prerequisites, manifest coexistence and refusal, both version-gate outcomes, and the write-free snapshot.
- [ ] Run the focused `internal/migrate` suite, then `make build && make test`.

## Context Log

Pending.
