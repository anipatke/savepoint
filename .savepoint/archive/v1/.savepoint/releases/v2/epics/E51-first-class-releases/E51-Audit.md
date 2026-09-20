---
type: audit-findings
audited: 2026-09-20
---

# Audit Findings: E51 Make releases first-class in V2

## Main Findings

### Verdict

CLEAR. The two issues from the initial audit are closed, no in-scope finding remains, and no owner waiver or additional owner-run evidence is required for E51. The repository is **CLEAR TO COMMIT/PUSH** for this epic.

### Prior Finding Closure Map

| Prior finding | Result | Evidence |
|---|---|---|
| Release switch persisted a rejected cross-Release selection | Closed | The normal `r` → Enter workflow now clears stale Objective/Task context, reloads with no selection diagnostic, retains R002-scoped Next, and produces matching board/resume facts. Direct cross-Release writes fail without changing `router.md`. `TestReleaseSelectionFiltersIndexedObjectivesAndPersistsOnlyRouterContext` and `TestSelectionWriteRejectsCrossReleaseObjective` pass uncached. |
| Invalid free-text Release metadata was accepted in Release-free projects | Closed | `indexReleaseObjectives` now rejects every non-empty non-`R###` reference regardless of Release count, naming the file, Objective, and value. `TestLoadV2Index_rejectsLegacyPackagingTextWithoutReleases` and the valid unassigned no-Release regression pass uncached. |

### What Needs Attention

No product or repository issue remains in the frozen E51 audit scope.

### Materiality Summary

No materiality actions are required.

### What Is Proven / Not Proven

Proven: optional Release identity and confined discovery; typed Objective membership and derived reverse links; Release-scoped immutable Checks; completion refusal for incomplete Objectives, missing/NEEDS WORK/stale/unknown or superseded evidence, material unresolved Issues, and missing or stale exact owner acceptance; historical-completion separation; source-preserving and conflict-safe writes; strict Release-aware selection and Next behavior; board, plain-output, and resume parity; doctor diagnostics; optional no-Release scaffolds; canonical/template skill parity; migration inventory, preview, archive mapping, apply, interruption recovery, edit-conflict refusal, retry, and unchanged second apply; temporary repository-copy accountability; and E50's composition of the canonical Release cutover decision without a second readiness rule.

No E51 acceptance criterion remains unverified. Publishing, deployment, tagging, changelog generation, network services, authentication, billing, and migration of the live repository remain intentionally outside E51.

### Audit Evidence

- Scope lock: reused the initial audit's immutable T001–T009 acceptance scope, E51 quality gates, public data/doctor/board/resume/migration/init surfaces, Release and selection state classes, evidence/Issue/acceptance states, TTY and non-TTY representations, and migration operation sequence. No new blocking axis or dependency layer was admitted.
- Coverage and workflow result: absent/one/multiple Releases; unassigned, valid, dangling, malformed, and historical membership; planned/in-progress/done lifecycle; current/stale/unknown/NEEDS WORK/superseded evidence; owner wait/acceptance; unresolved/resolved/excepted Issues; valid/missing/mismatched selection; selector open/navigation/cancel/switch/reload/conflict; and migration preview/apply/interruption/conflict/retry/no-op cells all pass. The two original failed cells pass after remediation.
- Side-effecting workflows: selection validates fresh indexed ownership before its guarded router write, clears stale cross-Release context, reloads, and preserves bytes on refusal. Migration still orders source inventory, decisions, targets/archives, preview, backup, staged writes, archive verification, protected publication, schema activation, manifest recovery, conflict refusal, retry, and unchanged repeat without silent overwrite.
- File reality and drift: every task Context File exists; no unexplained phantom or discarded scratch file was found. Canonical and scaffolded Idea, Design, and Check skills are byte-identical. Release architecture remains reconciled in Design, v2-Design, README, E50, and the AGENTS Codebase Map.
- Gates: focused remediation and Release tests passed uncached; `go test ./internal/data ./internal/doctor ./internal/board/v2 ./internal/resume ./internal/init . -count=1` passed; `go test ./internal/migrate -count=1` passed in 128.950s; `make build` passed; `make test` passed, including migration in 122.880s; `git diff --check 61cf510..HEAD` and the working-tree `git diff --check` passed.

### Guardrails Verification

- Rule IDs checked: FS-01, FS-03–06; DATA-01–04; TPL-01–04; ARCH-01–04; CFG-01–02; DEP-01–02; TEST-01–08; POL-01–02. REL-01–03 remain not applicable because E51 changes no distribution target, archive, checksum, or version path.
- Health check mode: Full. `.savepoint/Health-Check.md` is absent, so no additional project procedure was applicable.
- Evidence: dry-run/no-write, managed-write preservation and idempotence, named malformed-data diagnostics, canonical gate ownership, renderer I/O boundaries, migration recovery/conflict safety, temporary-fixture isolation, skill/template parity, focused failures, and the required full build/test gates all pass.
- File reality evidence: all scoped paths and remediation files exist; current source lines and uncached tests independently confirm both repairs.
- Waivers or unresolved findings: none.

## Code Style Review

- [x] STYLE-01 **One job per file** — split files when responsibilities mix.
- [x] STYLE-02 **One job per function** — small, named, testable units.
- [x] STYLE-03 **Test branches** — cover meaningful conditionals and edge cases.
- [x] STYLE-04 **Types document intent** — prefer explicit types over comments.
- [x] STYLE-05 **Build only what is needed** — no speculative abstractions.
- [x] STYLE-06 **Handle errors at boundaries** — validate inputs, APIs, IO, and external data.
- [x] STYLE-07 **One source of truth** — no duplicated rules, constants, state, or config.
- [x] STYLE-08 **Comments explain why** — not what the code already says.
- [x] STYLE-09 **Content lives in data** — keep copy/config out of logic.
- [x] STYLE-10 **Small diffs** — minimal, reviewable, behaviour-preserving changes.
