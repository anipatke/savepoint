---
id: E51-first-class-releases/T004-migrate-every-release-and-its-source-promise
status: planned
objective: Convert every V1 release into a deterministic V2 Release with live source intent, archived bytes, and accountable references.
depends_on:
  - E51-first-class-releases/T002-require-release-integration-evidence-and-owner-acceptance
complexity_tier: high
complexity_reason: Extends recoverable migration across records, documents, identities, evidence, and ambiguous legacy state.
---

# T004: Migrate every Release and its source promise

## Problem

E45 currently maps V1 release names to optional metadata and release PRD material toward Objective scope. That loses the delivery boundary E51 makes authoritative. Migration must assign stable Release identity, preserve the release PRD as live intent plus exact-byte history, link converted Objectives, and remain recoverable and truthful about legacy evidence.

## Context Files

- `.savepoint/releases/v2/epics/E51-first-class-releases/E51-Detail.md`
- `internal/data/release_v2.go`
- `internal/data/release_gate_v2.go`
- `internal/migrate/inventory.go`
- `internal/migrate/inventory_test.go`
- `internal/migrate/classify.go`
- `internal/migrate/classify_test.go`
- `internal/migrate/plan.go`
- `internal/migrate/plan_test.go`
- `internal/migrate/convert.go`
- `internal/migrate/convert_test.go`
- `internal/migrate/convert_docs.go`
- `internal/migrate/convert_docs_test.go`
- `internal/migrate/decisions.go`
- `internal/migrate/decisions_test.go`
- `internal/migrate/manifest.go`
- `internal/migrate/manifest_test.go`
- `internal/migrate/preview.go`
- `internal/migrate/testdata/golden/v1-basic.yml`
- `internal/migrate/testdata/golden/v1-history.yml`

## Acceptance Criteria

- [ ] Migration inventory classifies every V1 release directory and release PRD before conversion planning.
- [ ] Deterministic `R###` allocation is source-order stable, collision-aware, reserved in the manifest, and unchanged on a repeated preview/apply.
- [ ] Every converted Objective references the Release produced from its source-qualified V1 release; duplicate legacy Task/Epic IDs do not affect that mapping.
- [ ] Release PRD intent has an accountable live destination in the Release body, while the original file remains exact-byte archived with its hash recorded.
- [ ] Unmapped authored sections remain visibly preserved as legacy source and are never silently discarded or duplicated into unrelated Objectives.
- [ ] Settled historical Releases use typed legacy-completion references and never receive fabricated V2 Checks or CLEAR evidence.
- [ ] Active Releases enter `in_progress` with unknown Release clearance; unresolved defects/findings remain live Issues associated through preserved source/Check links.
- [ ] Ambiguous release lifecycle, source identity, or completion evidence appears in preview as a blocking owner decision before apply writes anything.
- [ ] Migration manifest mappings resolve source release path/name and every converted Objective/reference to the same `R###`.
- [ ] Existing backup, recovery, conflict, schema-activation-last, and second-run no-op guarantees remain intact.

## Implementation Plan

- [ ] Add V1 release and release-PRD roles to the frozen migration inventory/classification vocabulary.
- [ ] Reserve and allocate Release identities from the source inventory before Objective conversion.
- [ ] Generate Release records and Objective references in the conversion plan.
- [ ] Map release PRD sections into the Release body and preserve all remaining authored bytes/accountability.
- [ ] Add deterministic historical-completion versus active/ambiguous decisions without treating V1 audit state as V2 clearance.
- [ ] Extend manifest/reference mappings and preview output with Release destinations and owner decisions.
- [ ] Expand basic/history golden fixtures for multiple releases, duplicate scoped IDs, active work, settled history, and unresolved follow-up.
- [ ] Prove existing recovery and no-op invariants still cover the enlarged plan.

## Context Log

Pending.
