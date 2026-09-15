---
id: E44-issues-objective-checks/T002-connect-issues-to-their-work
title: Connect issues to the work they came from
status: planned
objective: Resolve Issue links to Checks and Tasks in both directions and derive Issue listings from the index.
depends_on:
    - E44-issues-objective-checks/T001-record-durable-follow-up
complexity_tier: medium
complexity_reason: Resolves one relation across two record families with a single authoritative direction.
---

# T002: Connect issues to the work they came from

## Problem

A Check's `issues` field is still an unresolved list of strings, and an Issue's `tasks` and `checks` fields name records nothing verifies. Nothing says which direction of the Check-to-Issue relation is authoritative, so the two sides can disagree silently and every consumer would have to maintain its own summary.

## Context Files

- `internal/data/issue_v2.go`
- `internal/data/issue_v2_test.go`
- `internal/data/check_v2.go`
- `internal/data/project.go`
- `internal/data/project_test.go`
- `internal/data/errors.go`
- `internal/data/task_v2.go`

## Acceptance Criteria

- [ ] Every `I###` reference in a Check's `issues` resolves to an existing Issue; a missing target returns a named diagnostic identifying the Check and the reference.
- [ ] Every `T###` in an Issue's `tasks` and every `C###` in its `checks` resolves to an existing record; a missing target returns a named diagnostic identifying the Issue and the reference.
- [ ] The immutable Check is the authoritative direction: a Check naming an Issue that does not appear in that Issue's `checks` list returns a named diagnostic naming both records, and the load fails closed rather than reconciling either side.
- [ ] An Issue may list a Check that does not name it — a recheck or proof Check that observed the Issue without recording it — and that asymmetry is legal.
- [ ] `V2Index` exposes a Task ID to Issue IDs map and a Check ID to Issue IDs map, both in deterministic ID order.
- [ ] Derived listings answer open, in_progress, and resolved Issue sets and counts, and per-type counts, from the index alone with no stored summary, register, or cached total anywhere in the project files.
- [ ] Listings are computed from the loaded index on each call and never persisted to disk.
- [ ] A project whose Issues reference no Tasks or Checks loads with empty link maps and no error.

## Implementation Plan

- [ ] Add link resolution to `LoadV2Index` in `project.go` after Issue discovery: resolve Check `issues`, Issue `tasks`, and Issue `checks`, then assert the Check-authoritative pairing rule.
- [ ] Add the named link sentinels to `errors.go`, distinguishing a missing target from an unpaired Check-to-Issue link.
- [ ] Build the Task-to-Issue and Check-to-Issue maps during index construction, sorting IDs for deterministic output.
- [ ] Add derived listing and count helpers to `issue_v2.go` operating on `*V2Index`, returning sorted IDs.
- [ ] Update the doc comment on `CheckV2.Issues` in `check_v2.go` so it no longer says the references are never loaded or validated.
- [ ] Test resolved links, missing Issue/Task/Check targets, the unpaired-link diagnostic in both orderings, the legal asymmetric case, map ordering, and listing and count results across mixed statuses and types.
- [ ] Run the focused `internal/data` suite.

## Context Log

Pending.
