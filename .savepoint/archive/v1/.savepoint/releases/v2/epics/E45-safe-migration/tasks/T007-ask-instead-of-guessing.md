---
id: E45-safe-migration/T007-ask-instead-of-guessing
title: Ask instead of guessing
status: done
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

- [x] Every ambiguity carries a stable ID derived deterministically from its source path and kind, so repeated previews of an unchanged project produce identical IDs.
- [x] Each ambiguity states its source path, its kind, why it is ambiguous, and the concrete choices that would resolve it.
- [x] Blocking ambiguities cover at least: an unrecognized lifecycle value on an active record, a dependency naming a missing target, a duplicate source identity, and a narrative finding that names no stable identity.
- [x] Advisory ambiguities — including an unclassified file that is archived intact — are reported and do not block.
- [x] A decisions input file maps ambiguity IDs to concrete outcomes, and the plan reports whether it is appliable: a plan with an unresolved blocking ambiguity is not, and names every unresolved ID.
- [x] A decisions file naming an ID the plan did not raise is an error, not a silently ignored entry, so a stale decisions file cannot be mistaken for a current one.
- [x] A decisions file supplying a value outside an ambiguity's stated choices is rejected with a named diagnostic.
- [x] Resolved decisions are recorded into the manifest model with their provenance: the decision, the source file, and the time.
- [x] Reading a decisions file writes nothing, and the preview remains write-free with decisions supplied.
- [x] No ambiguity is ever resolved by a default, a heuristic, or a similarity comparison; a test proves two narrative findings describing the same symptom are not merged.

## Implementation Plan

- [x] Add `decisions.go` with the ambiguity type, the stable ID derivation, and the blocking-versus-advisory classification.
- [x] Add the decisions file schema and its reader, with named diagnostics for unknown IDs and out-of-range values.
- [x] Extend the plan value with its appliable state and the list of unresolved blocking IDs, keeping `Plan` pure.
- [x] Record accepted decisions and their provenance into the manifest model.
- [x] Test ID stability across repeated previews, each blocking kind, the advisory path, unknown and invalid decision entries, provenance recording, the write-free property with decisions supplied, and the no-merge proof.
- [x] Run the focused `internal/migrate` suite, then `make build && make test`.

## Context Log

**Files read:** `.savepoint/router.md`, `E45-Detail.md`, this task file, `internal/migrate/plan.go`, `internal/migrate/manifest.go`, `internal/migrate/plan_test.go`, `internal/migrate/manifest_test.go`, `internal/migrate/classify.go`, `internal/migrate/fixture_test.go`, `internal/migrate/inventory_test.go`, `internal/data/errors.go`, `internal/data/lifecycle.go`, `internal/data/audit_finding.go` (+ its test for the finding frontmatter shape), `.savepoint/Guardrails.md`.

**Files edited:**
- `internal/migrate/decisions.go` (new) — `AmbiguityKind`/`Ambiguity`/`Decisions`/`DecisionValue` moved here from `plan.go` and extended with `Blocking`, `Choices`, `DecisionSourceFile`/`DecisionAt`; `isBlockingAmbiguityKind`; `ReadDecisionsFile` (YAML list schema, provenance-stamped, pure read); `validateDecisions` (unknown-ID and out-of-choices diagnostics via `ErrUnknownAmbiguityID`/`ErrDecisionValueNotAllowed`); `addAmbiguity`/`addAmbiguityDiscriminated` (the latter folds a discriminator, e.g. a missing dependency's reference text or the first duplicate path, into the ID so a kind that can recur against the same path never collides).
- `internal/migrate/plan.go` — removed the superseded `Ambiguity`/`AmbiguityKind`/`Decisions`/`addAmbiguity` definitions; added `ConversionPlan.Appliable`/`UnresolvedBlockingIDs`; `Plan` now calls `validateDecisions` after `build()` and computes the unresolved-blocking list before returning; added `AmbiguityDuplicateSourceIdentity` detection in `planTasksImpl` (a `taskIdentitySeen` map keyed by release/epic/original task ID); added the `AmbiguityUnclassifiedFile` advisory in `archiveRemainingRoles`; every existing `addAmbiguity` call site now supplies `Choices` (`epicStatusChoices`, `taskStatusChoices`, `findingStatusChoices`, `missingDependencyChoices`, `unresolvedNarrativeFindingChoices`); the two missing-dependency call sites in `resolveTaskDependencies` now use `addAmbiguityDiscriminated` keyed by the unresolved reference.
- `internal/migrate/manifest.go` — `ManifestDecision` gained `SourceFile`/`DecidedAt`; `BuildManifest` now carries `Ambiguity.DecisionSourceFile`/`DecisionAt` into it.
- `internal/migrate/decisions_test.go` (new) — covers every AC: ID stability (`TestAmbiguity_idStableAcrossRepeatedPreviews`), self-description (`TestAmbiguity_statesPathKindDetailAndChoices`), each new/updated blocking kind (`TestPlan_duplicateSourceIdentity_isBlockingAmbiguity`, `TestPlan_unresolvedNarrativeFinding_isBlockingAmbiguity`), the advisory path (`TestPlan_unclassifiedFile_isAdvisoryNotBlocking`), decision resolution and appliability (`TestPlan_decisionResolvesBlockingAmbiguity_andMakesPlanAppliable`), the two decisions-file diagnostics (`TestPlan_decisionsFile_unknownID_isError`, `TestPlan_decisionValueOutsideChoices_isError`), manifest provenance (`TestBuildManifest_recordsDecisionProvenance`), the write-free property for both reading a decisions file and previewing with one supplied (`TestReadDecisionsFile_readsWithoutWriting`, `TestPlan_writeFreeWithDecisionsSupplied`), and the no-auto-merge proof (`TestPlan_similarFindingsAreNeverAutoMerged`).
- `.savepoint/releases/v2/epics/E45-safe-migration/tasks/T007-ask-instead-of-guessing.md` — this file.

**Design notes not obvious from the diff:**
- `Choices` is deliberately about resolving *the ambiguity Plan raised*, not about re-deriving a correct value from the V1 source: `missing_dependency_target`'s only choice is `drop_dependency` (migration never invents a replacement target — a real fix is a V1 source edit followed by a rerun), and `duplicate_source_identity`'s choices are `keep_first`/`keep_second` naming which of the two conflicting declarations is canonical. Applying a decision's effect on downstream Targets/Archives is out of this task's scope (a later apply task); T007 only records that the ambiguity is resolved and gates `Appliable`.
- `AmbiguityUnresolvedNarrativeFind` (a `duplicate` finding whose `duplicate_of` names no known finding) was already partially implemented by an earlier task; this task added its `Choices`, confirmed it is `Blocking`, and added dedicated coverage since none existed before.
- Kept the existing `kind:path` ID format for single-occurrence kinds exactly as `plan_test.go`'s pre-existing `TestPlan_unrecognizedTaskStatus_isBlockingAmbiguityNotHealed` already asserted; added a `#<discriminator>` suffix only for kinds that can recur against the same path, so no existing ID changed.

**Quality gates:**
- `go build ./...` — pass.
- `go vet ./...` — pass.
- `go test ./internal/migrate/...` — pass (full package, including all pre-existing tests).
- `make build && make test` — pass.
- `.savepoint/Health-Check.md` — absent from this project; Quick-check step skipped per the skill (absence is not a finding).
