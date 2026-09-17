---
id: E45-safe-migration/T007-ask-instead-of-guessing
title: Ask instead of guessing
status: planned
objective: Give every ambiguity a stable ID and an owner decision input, and make unresolved blocking ambiguities refuse application.
depends_on:
    - E45-safe-migration/T003-plan-the-conversion-before-touching-anything
complexity_tier: medium
complexity_reason: Adds stable ambiguity identity and a decision input over an existing plan contract.
---

# T007: Ask instead of guessing

## Problem

The earlier conversion tasks all route the same way when a source is unclear: raise an ambiguity. None of them can say what happens next, because there is nowhere to answer.

An ambiguity needs three things to be answerable. A stable identity, so the same project previews the same IDs today and tomorrow and the user can refer to one. A place to put the answer, which cannot be the preview itself, because the preview writes nothing. And a consequence, so that an unanswered blocking ambiguity stops the migration instead of being quietly resolved by whichever branch the code happened to take.

The distinction that matters is blocking versus advisory. An unrecognized status on an active Task changes what the migrated project means, so it blocks. An unclassified file that will be archived intact is worth reporting and does not.

## Context Files

- `internal/migrate/decisions.go`
- `internal/migrate/decisions_test.go`
- `internal/migrate/plan.go`
- `internal/migrate/plan_test.go`
- `internal/migrate/manifest.go`
- `internal/data/errors.go`
- `internal/data/testdata/migration/v1-history/manifest.yml`
- `.savepoint/Guardrails.md`

## Acceptance Criteria

- [ ] Every ambiguity carries a stable ID derived deterministically from its source path and kind, so repeated previews of an unchanged project produce identical IDs.
- [ ] Each ambiguity states its source path, its kind, why it is ambiguous, and the concrete choices that would resolve it.
- [ ] Blocking ambiguities cover at least: an unrecognized lifecycle value on an active record, a dependency naming a missing target, a duplicate source identity, and a narrative finding that names no stable identity.
- [ ] Advisory ambiguities — including an unclassified file that is archived intact — are reported and do not block.
- [ ] A decisions input file maps ambiguity IDs to concrete outcomes, and the plan reports whether it is appliable: a plan with an unresolved blocking ambiguity is not, and names every unresolved ID.
- [ ] A decisions file naming an ID the plan did not raise is an error, not a silently ignored entry, so a stale decisions file cannot be mistaken for a current one.
- [ ] A decisions file supplying a value outside an ambiguity's stated choices is rejected with a named diagnostic.
- [ ] Resolved decisions are recorded into the manifest model with their provenance: the decision, the source file, and the time.
- [ ] Reading a decisions file writes nothing, and the preview remains write-free with decisions supplied.
- [ ] No ambiguity is ever resolved by a default, a heuristic, or a similarity comparison; a test proves two narrative findings describing the same symptom are not merged.

## Implementation Plan

- [ ] Add `decisions.go` with the ambiguity type, the stable ID derivation, and the blocking-versus-advisory classification.
- [ ] Add the decisions file schema and its reader, with named diagnostics for unknown IDs and out-of-range values.
- [ ] Extend the plan value with its appliable state and the list of unresolved blocking IDs, keeping `Plan` pure.
- [ ] Record accepted decisions and their provenance into the manifest model.
- [ ] Test ID stability across repeated previews, each blocking kind, the advisory path, unknown and invalid decision entries, provenance recording, the write-free property with decisions supplied, and the no-merge proof.
- [ ] Run the focused `internal/migrate` suite, then `make build && make test`.

## Context Log

Pending.
