---
id: E45-safe-migration/T005-carry-unfinished-follow-up-across
title: Carry unfinished follow-up across
status: done
objective: Convert unresolved V1 defects and audit findings into V2 Issues by disposition, archiving the rest without implying a Check.
depends_on:
    - E45-safe-migration/T003-plan-the-conversion-before-touching-anything
complexity_tier: high
complexity_reason: Maps ten finding states and defect lifecycle onto E44 Issue obligations without fabricating proof.
---

# T005: Carry unfinished follow-up across

## Problem

V1 spreads follow-up across three shapes: release defects, ten audit finding states, and a separately maintained register. V2 has one Issue record, and E44 gave it obligations that a careless conversion will violate immediately.

The specific traps: a `fixed` finding is not a verified one — it still awaits independent proof, and converting it to `resolved` would launder an unverified repair into a closed record. A `verified` finding's proof was a V1 audit, not a V2 Check, so it must archive rather than arrive with a fabricated `resolution.check`. A `duplicate` finding can only resolve as `duplicate` if its canonical counterpart also becomes an Issue, because E44 requires `duplicate_of` to name an existing record. And a `deferred` or `owner_decision` finding is still open — E44 made deferral a dated history entry precisely so it could not become a fourth lifecycle state.

Dates are the quiet one. An Issue's seeded history entry needs a timestamp, and inventing one would put a fictional date into an append-only record.

## Context Files

- `internal/migrate/convert_issues.go`
- `internal/migrate/convert_issues_test.go`
- `internal/migrate/convert.go`
- `internal/migrate/plan.go`
- `internal/data/issue_v2.go`
- `internal/data/defect.go`
- `internal/data/audit_finding.go`
- `internal/data/audit_run.go`
- `internal/data/audit_register.go`
- `internal/data/testdata/migration/v1-history/manifest.yml`
- `.savepoint/releases/v2/v2-Design.md`

## Acceptance Criteria

- [x] Rendered Issue content decodes through `DecodeIssueV2` with no diagnostic for every converted record in the `v1-history` fixture.
- [x] An unresolved V1 defect becomes an Issue with `type: defect` and status `open` or `in_progress` mapped from its V1 state; a resolved defect is archived and produces no Issue.
- [x] Findings in `open`, `triaged`, or `mapped` become open Issues preserving their reasons, proof requirements, and relations.
- [x] A `deferred` finding becomes an open Issue carrying a dated `deferred` history entry and no resolution; an `owner_decision` finding that resolved nothing becomes an open Issue carrying a dated `owner_decision` entry.
- [x] An `in_progress` or `fixed` finding becomes an `in_progress` Issue with a `repair_attempted` history entry; neither is ever converted to `resolved`, and a test proves `fixed` still awaits verification.
- [x] A `verified` finding is archived with its original proof and status and produces no Issue and no V2 Check.
- [x] A `waived` finding is archived with its disposition; when active converted work references it, a typed historical decision entry is recorded in the manifest rather than an Issue being invented to hold it.
- [x] A `duplicate` finding resolves as `duplicate` with `duplicate_of` naming the converted canonical Issue only when that canonical also converts; when the canonical is archive-only, the duplicate is archive-only too, and a test covers both directions using the fixture's `F003-duplicate`.
- [x] Every converted Issue's `source` records `kind: migration` with its source-qualified legacy key, and never borrows a Check's authority.
- [x] A seeded history entry uses the date the V1 record recorded; when the source recorded none, it uses the migration time and carries an explicit note saying the original date was absent.
- [x] Converted Issues link to converted repair Tasks through `tasks` when the V1 record named one, and to nothing when it did not; no link is inferred from narrative similarity.
- [x] The V1 audit register, runs, and prompt are archived intact and no register is recreated; Issue listings in V2 are derived, as E44 established.
- [x] A finding whose narrative names no stable identity and cannot be deterministically mapped produces a blocking ambiguity rather than a guessed Issue.

## Implementation Plan

- [x] Add `convert_issues.go` with the disposition mapping table expressed as data, not as a chain of conditionals.
- [x] Implement defect conversion over the existing `internal/data/defect.go` reader, and finding conversion over `audit_finding.go`, treating both as read-only migration inputs.
- [x] Implement the duplicate rule with its dependency on whether the canonical converts, resolving it after allocation so both outcomes are decidable.
- [x] Implement `source` construction and the dated-history seeding rule, including the absent-date note.
- [x] Route unmappable narrative findings to blocking ambiguities rather than to Issues.
- [x] Add round-trip tests decoding every rendered Issue through `DecodeIssueV2`.
- [x] Test each of the ten finding states, both defect states, both duplicate directions, the archived-with-reference waiver path, `source` provenance, and both history-date paths.
- [x] Run the focused `internal/migrate` and `internal/data` suites, then `make build && make test`.

## Context Log

**Files read:** `internal/migrate/convert.go`, `plan.go`, `manifest.go`, `classify.go`, `fixture_test.go`, `convert_test.go`, `plan_test.go`, `manifest_test.go`; `internal/data/issue_v2.go`, `issue_v2_test.go`, `defect.go`, `audit_finding.go`, `evidence_v2.go`, `parser.go`; the `v1-history` fixture project (defects, findings, tasks, epics) and `templates/project/.savepoint/audit/findings/README.md` (to confirm `deferral_reason` covers both `deferred` and `owner_decision`).

**Files edited:**
- `internal/migrate/convert_issues.go` (new) — `ConvertIssue`, dispatching to `convertDefectIssue`/`convertFindingIssue` by `Classify`; the `findingDispositionRules` data table; source/history/duplicate rendering; `findingMetadataBody` (renders `proof_needed`, dates, and relation fields that have no dedicated Issue field, so AC3's "proof requirements and relations" survive rather than being silently dropped as "known but unmapped").
- `internal/migrate/convert_issues_test.go` (new) — round-trip/determinism over both frozen fixtures, defect dispositions, all finding dispositions, waived-reference recording, both duplicate directions, source provenance, history date seeding (both the parsed and absent-date paths), task linking, `legacy_fields` preservation, and the dangling-`duplicate_of` ambiguity.
- `internal/migrate/plan.go` — added `PlannedTarget.DuplicateOfGlobalID` (plan.go previously computed but discarded this via `_ = canonicalID`); added `WaivedReference` type, `ConversionPlan.WaivedRefs`, and `planBuilder.recordWaivedReferences`; fixed `planFinding`'s duplicate case so a `duplicate_of` naming no known finding raises `AmbiguityUnresolvedNarrativeFind` instead of silently archiving; added the shared `resolveTaskReference` helper (matches a V1 task reference — already shaped like `LegacyKey.OriginalID` — against active Task targets, scoped by release).
- `internal/migrate/manifest.go` — added `ManifestWaivedReference` and `ManifestV1ToV2.WaivedReferences`, populated and sorted in `BuildManifest`.
- `internal/migrate/convert_test.go` — replaced the `TargetIssue: // belongs to a later task` stub in the shared round-trip/determinism tests with real `ConvertIssue` calls, now that this task implements it.

**Design decisions:**
- Migration actor convention: every migrated Issue's `source`/`resolution`/`history` actor is `{role: owner, session: <V1 source path>}`. Session carries the source-qualified legacy key directly in the frontmatter (LegacyKey's own "actual uniqueness qualifier"), rather than only implying it via the relocated body.
- `source.at` and a duplicate's `resolution.at` (when no date parses) use `plan.GeneratedAt` — the true moment migration ran — never a V1-authored date presented as if it were.
- Finding `IssueType` is `guardrail` when the finding names `guardrail_ids`, else `other`; this is purely descriptive (no gate reads it) and derived from an authored field rather than invented.
- A `duplicate` finding converts as a **resolved** Issue with `resolution.disposition: duplicate` (not merely an open Issue naming `duplicate_of`), matching `IssueDispositionDuplicate`'s purpose in `internal/data`: closed, provable-of-nothing, pointing at the canonical.
- The archive-only-canonical duplicate direction is tested with a synthetic project rather than a modified frozen fixture (the frozen `F003-duplicate`'s canonical, `F001`, is `fixed`, which always converts) — the epic's own fixtures are frozen and must not be edited.

**Quality gates:** `go build ./...`, `go vet ./...`, `go test ./internal/migrate/...` (all new tests plus the full existing suite pass), `make build && make test` (all packages pass). No `.savepoint/Health-Check.md` exists at this project's own root, so the Quick check step is skipped per the build-task skill's own guidance (absence is not a finding).

**Drift:** None. No new files or modules outside the epic's documented `internal/migrate` package; `internal/data` was read-only. `PlannedTarget.DuplicateOfGlobalID`, `WaivedReference`, and `ManifestWaivedReference` extend `plan.go`/`manifest.go` (both in this task's Context Files) rather than introducing a new module.
