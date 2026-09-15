---
id: E44-issues-objective-checks/T003-prove-an-issue-is-closed
title: Prove an issue is really closed
status: planned
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

- [ ] `resolution` is required when `status: resolved` and rejected when status is `open` or `in_progress`; both violations return distinct named diagnostics.
- [ ] Disposition `verified` requires a proof `check` that exists, recorded `CLEAR`, and appears in the Issue's `checks` list; a missing, `NEEDS WORK`, or unlisted proof Check each return a named diagnostic.
- [ ] Disposition `accepted` requires an owner actor and a non-empty reason and must not name a proof Check; a missing owner, a missing reason, or a named Check each return a named diagnostic.
- [ ] Disposition `duplicate` requires `duplicate_of` to name an existing, different Issue and must not name a proof Check.
- [ ] An Issue returned to `open` or `in_progress` with a stale `resolution` still present is a named diagnostic, so reopening must clear it.
- [ ] Reopening is the same `I###` with `status: open` and a dated `reopened` history entry; nothing in the data layer creates or suggests a second Issue for the same recurring problem.
- [ ] A deferral is a dated `deferred` history entry on an Issue that stays `open`; no fourth Issue status exists and no deferral field is added.
- [ ] No Issue type, severity, status, or count changes any Task gate decision: `ResolveTaskStart`, `ResolveTaskAdvance`, and `ResolveTaskCompletion` return identical decisions with and without open Issues of every type present.
- [ ] A `NEEDS WORK` Check referencing an Issue still blocks solely through the existing clearance rules, with no Issue-derived blocker kind introduced.

## Implementation Plan

- [ ] Add resolution obligation validation to `issue_v2.go` as an index-level check, keeping the shape decoding from T001 separate from these cross-record rules.
- [ ] Call the obligation validation from `LoadV2Index` after link resolution, in sorted Issue ID order.
- [ ] Add the named resolution sentinels to `errors.go`, distinguishing a missing proof, an unusable proof, a disposition/field mismatch, and a stale resolution on a reopened Issue.
- [ ] Add history semantics helpers for reopening and deferral to `issue_v2.go` without introducing any new status value.
- [ ] Test each disposition's satisfied and violated forms, the reopen and deferral paths, and a superseded proof Check.
- [ ] Add gate tests in `gate_v2_test.go` proving Task decisions are byte-identical with and without Issues present, including a `NEEDS WORK` Check that references one.
- [ ] Run the focused `internal/data` suite.

## Context Log

Pending.
