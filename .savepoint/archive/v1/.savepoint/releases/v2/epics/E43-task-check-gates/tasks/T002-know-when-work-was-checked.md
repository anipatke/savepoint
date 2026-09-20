---
id: E43-task-check-gates/T002-know-when-work-was-checked
title: Know when work was last checked
status: done
objective: Decode the shared freshness, acceptance, exception, and replan evidence on Task and Objective records.
depends_on:
    - E43-task-check-gates/T001-record-check-results
complexity_tier: medium
complexity_reason: Adds one shared evidence contract to two record decoders with strict reference validation.
---

# T002: Know when work was last checked

## Problem

Tasks carry no recorded evidence, so nothing distinguishes work that was independently checked from work that merely claims a status, and nothing links a record to the Check that cleared it.

## Context Files

- `internal/data/evidence_v2.go`
- `internal/data/evidence_v2_test.go`
- `internal/data/task_v2.go`
- `internal/data/task_v2_test.go`
- `internal/data/objective_v2.go`
- `internal/data/objective_v2_test.go`
- `internal/data/check_v2.go`
- `internal/data/project.go`
- `internal/data/project_test.go`
- `internal/data/errors.go`

## Acceptance Criteria

- [x] Task and Objective records decode one shared evidence block: `last_check`, `freshness`, `owner_validation`, optional `exception`, and optional `replan`.
- [x] `freshness` requires `state: current|stale|unknown`, a `check` reference, `assessed_by` actor provenance, `assessed_at`, and a `basis`; any other state value is a named diagnostic.
- [x] `owner_validation` decodes `required` and an optional `accepted_check`; an absent block means owner validation is not required.
- [x] `exception` requires a non-empty `requirements` list, `reason`, `owner`, `recorded_at`, and the `check` it applies to.
- [x] `replan` requires `reason`, `recorded_by`, and `recorded_at`, and never changes the record's status or stage.
- [x] Every Check reference in the evidence block — `last_check`, `freshness.check`, `owner_validation.accepted_check`, `exception.check` — must name a Check that exists in the index, or the load fails with a named diagnostic.
- [x] All evidence fields are optional; a record carrying none decodes cleanly and carries no evidence rather than a defaulted one.
- [x] No evidence value is healed: an unknown state, a malformed timestamp, or a partially filled block is a diagnostic, never a value that could support completion.
- [x] There is no `checked: true` field, and existing V1 Task parsing remains behaviorally unchanged.

## Implementation Plan

- [x] Define the actor reference, freshness, owner-validation, exception, and replan types with one strict decoder in a new `evidence_v2.go`.
- [x] Decode the shared block from both `DecodeTaskV2` and `DecodeObjectiveV2` through that single decoder.
- [x] Add evidence-reference resolution to `LoadV2Index` so every named Check must exist, reporting path and identity context.
- [x] Add the evidence diagnostics to `errors.go`.
- [x] Test each required and optional field, each malformed variant, absent evidence, every dangling Check reference, and that a partially filled block is rejected.
- [x] Add a V1 regression test proving V1 Task parsing is untouched, then run the focused `internal/data` suite.

## Context Log

**Files read:** `internal/data/task_v2.go`, `internal/data/task_v2_test.go`, `internal/data/objective_v2.go`, `internal/data/objective_v2_test.go`, `internal/data/check_v2.go`, `internal/data/project.go`, `internal/data/project_test.go`, `internal/data/errors.go`, `internal/data/discover_test.go` (targeted read for `writeV2ObjectiveFixture`/`writeV2TaskFixture`/`writeV2CheckFixture` fixture helpers), `internal/data/parser.go` (targeted read to confirm `ParseV2Document`/`V2SourceDocument`), `.savepoint/Guardrails.md` (STYLE rules), `AGENTS.md`, `.savepoint/router.md`, `agent-skills/savepoint-build-task/SKILL.md`, `E43-Detail.md`, `T001-record-check-results.md` (targeted read to confirm dependency status and prior design decisions, e.g. `checked_at`/`recorded_at` RFC 3339 parsing, four-role actor vocabulary).

**Files edited:**
- `internal/data/evidence_v2.go` (new) — `FreshnessState`, `Freshness`, `OwnerValidation`, `Exception`, `Replan`, `Evidence` types; raw frontmatter counterparts (`freshnessV2Frontmatter`, `ownerValidationV2Frontmatter`, `exceptionV2Frontmatter`, `replanV2Frontmatter`, `evidenceActorFrontmatter`, `evidenceV2Frontmatter`); one shared `decodeEvidenceV2` plus per-sub-block helpers (`decodeFreshnessV2`, `decodeOwnerValidationV2`, `decodeExceptionV2`, `decodeReplanV2`, `decodeEvidenceActor`, `decodeEvidenceTimestamp`), reusing `checkIDPatternV2` (check_v2.go) and `Actor`/`ActorRole` (check_v2.go) rather than redefining them.
- `internal/data/task_v2.go` — embedded `evidenceV2Frontmatter` inline into `taskV2Frontmatter` (confirmed gopkg.in/yaml.v3 supports `yaml:",inline"` on an unexported embedded struct type with a throwaway snippet before relying on it); added `Evidence *Evidence` to `TaskV2`; `DecodeTaskV2` now calls `decodeEvidenceV2(path, "task", fields.ID, fields.evidenceV2Frontmatter)`.
- `internal/data/objective_v2.go` — same treatment: inline-embedded `evidenceV2Frontmatter`, added `Evidence *Evidence` to `ObjectiveV2`, wired `decodeEvidenceV2(path, "objective", fields.ID, ...)` into `DecodeObjectiveV2`.
- `internal/data/project.go` — added `validateEvidenceReferences` and `checkEvidenceReferences`, called from `LoadV2Index` after `indexChecks`; walks Task IDs then Objective IDs in sorted order and checks `last_check`, `freshness.check`, `owner_validation.accepted_check`, `exception.check` against `index.Checks` in fixed field order per record, for deterministic diagnostics.
- `internal/data/errors.go` — added `ErrV2EvidenceMalformed` (unknown enum values, unparseable timestamps, malformed list entries) and `ErrV2EvidenceMissingReference` (dangling Check reference), mirroring the existing `ErrV2Check*` sentinel split.
- `internal/data/evidence_v2_test.go` (new) — direct `decodeEvidenceV2` coverage: absent evidence is nil, `last_check`-only, each freshness state value, every required freshness/exception/replan field individually cleared (partial-block rejection), unknown enum values, malformed/missing timestamps, owner_validation absent/required-only/with accepted_check/malformed accepted_check, and all five sub-blocks decoding together on one record.
- `internal/data/task_v2_test.go` — added `TestDecodeTaskV2_evidenceValid`, `TestDecodeTaskV2_noEvidenceIsNil`, `TestDecodeTaskV2_malformedEvidencePropagatesDiagnostic` (full-document decode through `DecodeTaskV2`, not just the unit-level decoder).
- `internal/data/objective_v2_test.go` — added `TestDecodeObjectiveV2_evidenceValid`, `TestDecodeObjectiveV2_noEvidenceIsNil`.
- `internal/data/project_test.go` — added `TestLoadV2Index_evidenceReferencesResolve`, `TestLoadV2Index_evidenceMissingReference` (table over all four reference fields), `TestLoadV2Index_evidenceMissingReferenceOnObjective`.

**Named cases and results:**
- `TestDecodeEvidenceV2_*` (18 top-level, several table-driven covering every required/optional field and partial-block rejection) — PASS
- `TestDecodeTaskV2_evidence*`, `TestDecodeTaskV2_noEvidenceIsNil`, `TestDecodeTaskV2_malformedEvidencePropagatesDiagnostic` — PASS
- `TestDecodeObjectiveV2_evidence*`, `TestDecodeObjectiveV2_noEvidenceIsNil` — PASS
- `TestLoadV2Index_evidence*` (3, one table-driven over 4 reference fields) — PASS
- `TestParseTaskFile_v1BehaviorUnaffectedByV2Types` (pre-existing V1 regression test) — PASS, confirms V1 `Task` parsing (legacy `todo` healing, objective-as-title fallback) is untouched
- Full `go test ./internal/data/...` — PASS, no regressions in E42/T001 suites
- `go vet ./...` — clean
- `gofmt -l` on all edited/added files — clean

**Quality gates:** `make build && make test` — PASS across all packages (`.`, `cmd`, `internal/board`, `internal/buildtool`, `internal/data`, `internal/doctor`, `internal/init`, `internal/styles`).

No `.savepoint/Health-Check.md` in this project; Quick check step skipped per `savepoint-build-task`.

**Design decisions not fully spelled out in the task/epic:**
- `evidenceV2Frontmatter` is embedded inline (`yaml:",inline"`) into both `taskV2Frontmatter` and `objectiveV2Frontmatter` so the five raw evidence fields are declared exactly once, reusable as a single value (`fields.evidenceV2Frontmatter`) passed to the shared decoder — verified this works with an unexported embedded struct type against gopkg.in/yaml.v3 before committing to the approach, since duplicate-declaring the five fields on both frontmatter structs was the fallback if it hadn't.
- `owner_validation.required` has no separate presence check: YAML's zero-value-on-omission for `bool` (`false`) is indistinguishable from, and semantically consistent with, an explicit `required: false` — both mean "not required" — so no diagnostic is needed for an omitted `required` field within a present `owner_validation` block.
- `exception.requirements` entries are validated only as non-empty trimmed strings, not against an ID pattern (e.g. `I###`): the epic's Exception schema names them "requirement IDs" generically (illustrated in tests as Guardrails rule IDs like `TEST-08`), with no fixed shape stated in E43-Detail.md, unlike the `I###`-anchored `issues` list on Check records.
- `checked_at`-style timestamps (`freshness.assessed_at`, `exception.recorded_at`, `replan.recorded_at`) are parsed strictly as RFC 3339, matching T001's established convention for `Check.checked_at`.
- Missing-reference diagnostics report all four Check-reference fields (`last_check`, `freshness.check`, `owner_validation.accepted_check`, `exception.check`) through one sentinel, `ErrV2EvidenceMissingReference`, with the specific field named in the message text — mirrors `ErrV2CheckMissingReference`'s single-sentinel-per-diagnostic-class pattern from T001 rather than one sentinel per field.

No Drift Notes: `evidence_v2.go` is the exact new file E43-Detail.md's Components table already names (`internal/data/evidence_v2.go (new) | Define and decode the evidence blocks shared by Task and Objective records`), and every other edited file (`task_v2.go`, `objective_v2.go`, `project.go`, `errors.go`, and their tests) is already named in that same table for this epic.
