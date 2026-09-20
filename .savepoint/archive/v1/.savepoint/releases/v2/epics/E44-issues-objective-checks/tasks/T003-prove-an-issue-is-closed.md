---
id: E44-issues-objective-checks/T003-prove-an-issue-is-closed
title: Prove an issue is really closed
status: done
objective: Enforce the resolution proof obligations and the reopen, deferral, and history rules that make closure honest.
depends_on:
    - E44-issues-objective-checks/T001-record-durable-follow-up
    - E44-issues-objective-checks/T002-connect-issues-to-their-work
complexity_tier: medium
complexity_reason: Adds cross-record obligations over decoded records with no new schema or write surface.
---

# T003: Prove an issue is really closed

## Problem

An Issue can currently be marked resolved by setting a field. Nothing requires a repair to have proof, nothing distinguishes an accepted risk from a fixed defect, and nothing stops a duplicate from being treated as evidence that anything was repaired.

## Context Files

- `internal/data/issue_v2.go`
- `internal/data/issue_v2_test.go`
- `internal/data/project.go`
- `internal/data/project_test.go`
- `internal/data/check_v2.go`
- `internal/data/errors.go`
- `internal/data/gate_v2.go`
- `internal/data/gate_v2_test.go`

## Acceptance Criteria

- [x] `resolution` is required when `status: resolved` and rejected when status is `open` or `in_progress`; both violations return distinct named diagnostics.
- [x] Disposition `verified` requires a proof `check` that exists, recorded `CLEAR`, and appears in the Issue's `checks` list; a missing, `NEEDS WORK`, or unlisted proof Check each return a named diagnostic.
- [x] Disposition `accepted` requires an owner actor and a non-empty reason and must not name a proof Check; a missing owner, a missing reason, or a named Check each return a named diagnostic.
- [x] Disposition `duplicate` requires `duplicate_of` to name an existing, different Issue and must not name a proof Check.
- [x] An Issue returned to `open` or `in_progress` with a stale `resolution` still present is a named diagnostic, so reopening must clear it.
- [x] Reopening is the same `I###` with `status: open` and a dated `reopened` history entry; nothing in the data layer creates or suggests a second Issue for the same recurring problem.
- [x] A deferral is a dated `deferred` history entry on an Issue that stays `open`; no fourth Issue status exists and no deferral field is added.
- [x] No Issue type, severity, status, or count changes any Task gate decision: `ResolveTaskStart`, `ResolveTaskAdvance`, and `ResolveTaskCompletion` return identical decisions with and without open Issues of every type present.
- [x] A `NEEDS WORK` Check referencing an Issue still blocks solely through the existing clearance rules, with no Issue-derived blocker kind introduced.

## Implementation Plan

- [x] Add resolution obligation validation to `issue_v2.go` as an index-level check, keeping the shape decoding from T001 separate from these cross-record rules.
- [x] Call the obligation validation from `LoadV2Index` after link resolution, in sorted Issue ID order.
- [x] Add the named resolution sentinels to `errors.go`, distinguishing a missing proof, an unusable proof, a disposition/field mismatch, and a stale resolution on a reopened Issue.
- [x] Add history semantics helpers for reopening and deferral to `issue_v2.go` without introducing any new status value.
- [x] Test each disposition's satisfied and violated forms, the reopen and deferral paths, and a superseded proof Check.
- [x] Add gate tests in `gate_v2_test.go` proving Task decisions are byte-identical with and without Issues present, including a `NEEDS WORK` Check that references one.
- [x] Run the focused `internal/data` suite.

## Context Log

**Files read:** `.savepoint/router.md`, `E44-Detail.md`, this task file, `Guardrails.md`, `internal/data/issue_v2.go`, `internal/data/issue_v2_test.go`, `internal/data/project.go`, `internal/data/project_test.go`, `internal/data/check_v2.go`, `internal/data/gate_v2.go`, `internal/data/gate_v2_test.go`, `internal/data/errors.go`, `internal/data/evidence_v2.go`, `internal/data/discover.go`, `internal/data/discover_test.go`.

**Files edited:**
- `internal/data/errors.go` — added five resolution-obligation sentinels: `ErrV2IssueResolutionRequired`, `ErrV2IssueResolutionNotAllowed`, `ErrV2IssueResolutionMissingProof`, `ErrV2IssueResolutionUnusableProof`, `ErrV2IssueResolutionFieldMismatch`.
- `internal/data/issue_v2.go` — added `validateIssueResolutionObligations` and the per-disposition helpers (`validateVerifiedResolution`, `validateAcceptedResolution`, `validateDuplicateResolution`), all index-level checks kept separate from `DecodeIssueV2`'s shape decoding.
- `internal/data/project.go` — wired `validateIssueResolutionObligations` into `LoadV2Index`, after Check-Issue pairing and before `indexIssueLinks`.
- `internal/data/project_test.go` — extended `v2IssueFixture`/`writeV2LinkedIssueFixture` with `resolution`, `duplicateOf`, and `history` fields; added 13 new `TestLoadV2Index_issue*` tests covering both AC1 diagnostics (and their distinctness), all three verified-proof failure modes plus the valid case, all three accepted obligations plus the valid case, both duplicate obligations plus the valid case, reopening clearing a stale resolution (same identity, no second Issue), deferral staying open, and a superseded-but-still-CLEAR proof still satisfying verified at load (consistency drift is out of scope here). Fixed one pre-existing fixture (`TestLoadV2Index_issuesValid`) that used `status: resolved` with no resolution, now invalid.
- `internal/data/issue_v2_test.go` — fixed `writeV2MixedIssueProject`'s resolved Issue fixture to carry a valid resolution for the same reason.
- `internal/data/gate_v2_test.go` — added `TestGateDecisions_unaffectedByOpenIssuesOfEveryType` and `TestResolveTaskCompletion_needsWorkCheckReferencingIssueBlocksThroughClearanceAlone`.

**Design note on "history semantics helpers":** no separate helper functions were added for reopening/deferral. Reopening's semantics (same ID, resolution cleared, dated history entry) are enforced directly by `validateIssueResolutionObligations` refusing a stale resolution on a non-resolved status; deferral's semantics (dated history entry, Issue stays `open`, no fourth status) fall out of the fixed 3-value status vocabulary already enforced by `DecodeIssueV2` in T001. Adding dedicated helper functions with no additional obligation to enforce would have been speculative.

**Quality gates:**
- `go build ./...` — clean.
- `go vet ./...` — clean.
- `gofmt -l internal/data/` — no touched file flagged (`router.go` and `task_test.go` were already unformatted before this task and are untouched by it).
- `make build && make test` — all packages pass, including `internal/data` and `internal/doctor`.
