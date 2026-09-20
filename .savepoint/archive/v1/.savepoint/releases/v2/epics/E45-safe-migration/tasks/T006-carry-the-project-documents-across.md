---
id: E45-safe-migration/T006-carry-the-project-documents-across
title: Carry the project documents across
status: done
objective: Convert PRD, router, config, and Health-Check into their V2 places without inventing configuration or rewriting prose.
depends_on:
    - E45-safe-migration/T003-plan-the-conversion-before-touching-anything
complexity_tier: medium
complexity_reason: "Several small document mappings with one firm rule: preserve authored content, infer nothing."
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

- [x] `.savepoint/PRD.md` becomes `.savepoint/Idea.md` with its authored content preserved byte for byte; no section is added, reordered, reworded, or dropped, and the original is archived.
- [x] `router.md` keeps its `## Current state` YAML anchor and converts to the V2 vocabulary: `pre-implementation`, `epic-design`, and `epic-task-breakdown` map to `design`; `task-building` and `defect-building` map to `task`; `audit-pending` maps to `check`.
- [x] The converted router names the allocated `O###` and optional `T###` for the previously selected work, and its `next_action` prose is preserved rather than regenerated.
- [x] When the selected V1 epic or task was archived rather than converted, the router records a named note explaining that and selects no dangling reference; a test covers a router pointing at a completed task.
- [x] `config.yml` carries its `quality_gates`, theme, and every unknown key across verbatim; the only field migration changes in it is `schema_version`, and that change belongs to the activation step rather than to this task.
- [x] No `quality_gates` entry is ever inferred from prose: a test proves a project with a `Health-Check.md` full of command-looking lines and no configured gates ends with exactly the gates its V1 `config.yml` declared.
- [x] `Health-Check.md` is archived and preserved as an optional procedure document, and its candidate commands are reported in the preview as owner decisions rather than written anywhere.
- [x] `Design.md`, `Guardrails.md`, and `visual-identity.md` keep their names and bytes in place, and a test asserts they are untouched.
- [x] A project missing any of these optional documents converts without a finding, because absence is the expected V1 shape in `v1-basic`.
- [x] Every rendered document is deterministic for a given source and injected clock.

## Implementation Plan

- [x] Add `convert_docs.go` with one small renderer per document and the router state mapping expressed as a table.
- [x] Implement the PRD-to-Idea move as a byte-preserving relocation, not a re-render.
- [x] Implement router conversion over `internal/data/router.go`'s existing anchor handling, including the archived-selection note.
- [x] Implement the `Health-Check.md` preservation path and the candidate-command report, with no write to `config.yml`.
- [x] Confirm the in-place documents are planned as preserved rather than as targets.
- [x] Test each mapping, the archived-selection router case, the no-inferred-gates rule, absent optional documents, and deterministic output.
- [x] Run the focused `internal/migrate` suite, then `make build && make test`.

## Context Log

**Files read:** `.savepoint/router.md` (router state), AGENTS.md, `agent-skills/savepoint-build-task/SKILL.md`, E45-Detail.md, T006 task file, `internal/migrate/plan.go`, `internal/migrate/classify.go`, `internal/migrate/convert.go`, `internal/migrate/manifest.go`, `internal/migrate/fixture_test.go`, `internal/migrate/plan_test.go`, `internal/migrate/convert_test.go`, `internal/data/router.go`, `internal/data/config.go`, `internal/data/release_doc.go`, both fixture `manifest.yml` and `project/` trees under `internal/data/testdata/migration/`, `.savepoint/releases/v2/v2-Design.md` (router vocabulary, retained documents).

**Files edited:**
- `internal/migrate/plan.go` — added `DocumentKind`/`PlannedDocument`, `ConversionPlan.Documents`, `ArchiveEntry.CandidateCommands`, `planBuilder.planDocuments` (new Pass 4: PRD→Idea archive+document, router document, Health-Check archive+candidate commands), and `preservedInPlaceDocs` (name-based exemption for `Guardrails.md`/`visual-identity.md` from the unclassified-archive sweep, since the frozen `Role` vocabulary in `classify.go` — pinned by `TestClassify_roleVocabularyIsComplete` — deliberately gets no new roles for them).
- `internal/data/router.go` — added `ReplaceStateBlock`, the writer half of `ReadState`'s anchor-finding, reused to rewrite only the `## Current state` YAML block.
- `internal/migrate/convert_docs.go` (new) — `ConvertIdea` (byte-preserving relocation), `ConvertRouter` (state-vocabulary mapping + O###/T### resolution + archived-selection note), `CandidateHealthCheckCommands` (fenced-code-block extraction for preview only).
- `internal/migrate/convert_docs_test.go` (new) — fixture-based and ad hoc coverage for every AC above.

**Quality gates:** `go build ./...`, `go test ./internal/migrate/...`, `go test ./...`, and `make build && make test` all pass. No `.savepoint/Health-Check.md` exists in this project, so its Quick-check step is skipped per the skill's own instruction.

**Note:** `classify.go`, though not in this task's `## Context Files`, required a read (and `preservedInPlaceDocs` in `plan.go`, not a `classify.go` edit) to satisfy AC8 without widening the frozen `Role` vocabulary `TestClassify_roleVocabularyIsComplete` pins.
