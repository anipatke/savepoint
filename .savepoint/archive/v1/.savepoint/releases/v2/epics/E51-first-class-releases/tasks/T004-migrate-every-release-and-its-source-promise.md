---
id: E51-first-class-releases/T004-migrate-every-release-and-its-source-promise
status: done
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

- [x] Migration inventory classifies every V1 release directory and release PRD before conversion planning.
- [x] Deterministic `R###` allocation is source-order stable, collision-aware, reserved in the manifest, and unchanged on a repeated preview/apply.
- [x] Every converted Objective references the Release produced from its source-qualified V1 release; duplicate legacy Task/Epic IDs do not affect that mapping.
- [x] Release PRD intent has an accountable live destination in the Release body, while the original file remains exact-byte archived with its hash recorded.
- [x] Unmapped authored sections remain visibly preserved as legacy source and are never silently discarded or duplicated into unrelated Objectives.
- [x] Settled historical Releases use typed legacy-completion references and never receive fabricated V2 Checks or CLEAR evidence.
- [x] Active Releases enter `in_progress` with unknown Release clearance; unresolved defects/findings remain live Issues associated through preserved source/Check links.
- [x] Ambiguous release lifecycle, source identity, or completion evidence appears in preview as a blocking owner decision before apply writes anything.
- [x] Migration manifest mappings resolve source release path/name and every converted Objective/reference to the same `R###`.
- [x] Existing backup, recovery, conflict, schema-activation-last, and second-run no-op guarantees remain intact.

## Implementation Plan

- [x] Add V1 release and release-PRD roles to the frozen migration inventory/classification vocabulary.
- [x] Reserve and allocate Release identities from the source inventory before Objective conversion.
- [x] Generate Release records and Objective references in the conversion plan.
- [x] Map release PRD sections into the Release body and preserve all remaining authored bytes/accountability.
- [x] Add deterministic historical-completion versus active/ambiguous decisions without treating V1 audit state as V2 clearance.
- [x] Extend manifest/reference mappings and preview output with Release destinations and owner decisions.
- [x] Expand basic/history golden fixtures for multiple releases, duplicate scoped IDs, active work, settled history, and unresolved follow-up.
- [x] Prove existing recovery and no-op invariants still cover the enlarged plan.

## Context Log

**Files read:** `.savepoint/router.md`, `.savepoint/releases/v2/epics/E51-first-class-releases/E51-Detail.md`, `.savepoint/Guardrails.md`, this task, `internal/data/release_v2.go`, `internal/data/release_gate_v2.go`, the migration inventory/classification/planning/conversion/decision/manifest/preview sources and tests listed above, and the frozen `v1-basic`/`v1-history` migration fixtures and goldens.

**Files edited:**
- `internal/migrate/classify.go`, `classify_test.go` — release-directory and release-PRD classification vocabulary.
- `internal/migrate/plan.go`, `decisions.go` — deterministic R### reservation/allocation, release-source planning, Objective→Release mapping, lifecycle/completion/source ambiguities, duplicate-source owner decisions, and release archive hashes.
- `internal/migrate/convert_releases.go`, `convert_releases_test.go` — active and historical Release rendering, legacy-source preservation/fencing, typed historical completion, round-trip and write-free tests.
- `internal/data/release_v2.go` — fenced legacy-source headings no longer masquerade as duplicate required Release sections.
- `internal/migrate/convert.go`, `convert_docs.go`, `apply.go`, `manifest.go`, `preview.go` — Release references, router selection, target rendering, manifest identity/archive evidence, and preview output.
- `internal/migrate/convert_test.go`, `testdata/golden/v1-basic.yml`, `testdata/golden/v1-history.yml` — expected R### references and live/archive Release outputs.
- This task file — canonical build lifecycle, verified criteria/checklist, and handoff log; it remains `status: in_progress` for user review.

**Verification:**
- Focused release planning/conversion tests — PASS.
- `go test ./internal/migrate ./internal/data` — PASS.
- `make build && make test` — PASS across all packages.
- End-to-end migration goldens and existing recovery/second-run no-op coverage — PASS.
- No `.savepoint/Health-Check.md` exists at the project root, so the Quick health check is skipped per the build-task workflow.

No Drift Notes: all new code and tests remain within the documented `internal/migrate` and `internal/data` modules; no architecture or Codebase Map change was required.
