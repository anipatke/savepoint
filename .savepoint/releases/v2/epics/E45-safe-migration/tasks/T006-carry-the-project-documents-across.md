---
id: E45-safe-migration/T006-carry-the-project-documents-across
title: Carry the project documents across
status: planned
objective: Convert PRD, router, config, and Health-Check into their V2 places without inventing configuration or rewriting prose.
depends_on:
    - E45-safe-migration/T003-plan-the-conversion-before-touching-anything
complexity_tier: medium
complexity_reason: Several small document mappings with one firm rule: preserve authored content, infer nothing.
---

# T006: Carry the project documents across

## Problem

Records are not the only thing migration moves. `PRD.md` becomes `Idea.md`, `router.md` gains a new state vocabulary, `config.yml` keeps its quality gates and eventually its new schema version, and `Health-Check.md` has nowhere to go because V2 replaced it with config commands and optional procedures.

That last one is where a converter would overreach. `Health-Check.md` is authored Markdown that happens to mention commands; parsing prose into executable configuration means guessing which lines are commands, which are examples, and which are commentary — and then writing the guess into the file that decides what `doctor` runs. The honest alternative is to preserve the document, list the candidate commands in the preview, and let the owner decide.

The router has a smaller version of the same problem: its selected epic or task may be work that migration archives, and the replacement cannot point at a record that no longer exists.

## Context Files

- `internal/migrate/convert_docs.go`
- `internal/migrate/convert_docs_test.go`
- `internal/migrate/plan.go`
- `internal/data/router.go`
- `internal/data/config.go`
- `internal/data/release_doc.go`
- `internal/data/testdata/migration/v1-basic/manifest.yml`
- `internal/data/testdata/migration/v1-history/manifest.yml`
- `.savepoint/releases/v2/v2-Design.md`

## Acceptance Criteria

- [ ] `.savepoint/PRD.md` becomes `.savepoint/Idea.md` with its authored content preserved byte for byte; no section is added, reordered, reworded, or dropped, and the original is archived.
- [ ] `router.md` keeps its `## Current state` YAML anchor and converts to the V2 vocabulary: `pre-implementation`, `epic-design`, and `epic-task-breakdown` map to `design`; `task-building` and `defect-building` map to `task`; `audit-pending` maps to `check`.
- [ ] The converted router names the allocated `O###` and optional `T###` for the previously selected work, and its `next_action` prose is preserved rather than regenerated.
- [ ] When the selected V1 epic or task was archived rather than converted, the router records a named note explaining that and selects no dangling reference; a test covers a router pointing at a completed task.
- [ ] `config.yml` carries its `quality_gates`, theme, and every unknown key across verbatim; the only field migration changes in it is `schema_version`, and that change belongs to the activation step rather than to this task.
- [ ] No `quality_gates` entry is ever inferred from prose: a test proves a project with a `Health-Check.md` full of command-looking lines and no configured gates ends with exactly the gates its V1 `config.yml` declared.
- [ ] `Health-Check.md` is archived and preserved as an optional procedure document, and its candidate commands are reported in the preview as owner decisions rather than written anywhere.
- [ ] `Design.md`, `Guardrails.md`, and `visual-identity.md` keep their names and bytes in place, and a test asserts they are untouched.
- [ ] A project missing any of these optional documents converts without a finding, because absence is the expected V1 shape in `v1-basic`.
- [ ] Every rendered document is deterministic for a given source and injected clock.

## Implementation Plan

- [ ] Add `convert_docs.go` with one small renderer per document and the router state mapping expressed as a table.
- [ ] Implement the PRD-to-Idea move as a byte-preserving relocation, not a re-render.
- [ ] Implement router conversion over `internal/data/router.go`'s existing anchor handling, including the archived-selection note.
- [ ] Implement the `Health-Check.md` preservation path and the candidate-command report, with no write to `config.yml`.
- [ ] Confirm the in-place documents are planned as preserved rather than as targets.
- [ ] Test each mapping, the archived-selection router case, the no-inferred-gates rule, absent optional documents, and deterministic output.
- [ ] Run the focused `internal/migrate` suite, then `make build && make test`.

## Context Log

Pending.
