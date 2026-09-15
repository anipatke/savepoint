---
id: E44-issues-objective-checks/T004-save-issues-without-losing-history
title: Save issues without losing history
status: planned
objective: Create Issue records once, patch their managed fields safely, and append history entries without rewriting the past.
depends_on:
    - E44-issues-objective-checks/T001-record-durable-follow-up
    - E44-issues-objective-checks/T003-prove-an-issue-is-closed
complexity_tier: medium
complexity_reason: Extends the existing preserving write path with append-only list semantics and a new create path.
---

# T004: Save issues without losing history

## Problem

An Issue is mutable where a Check is immutable, and its history is the part that must never be rewritten. The current write path patches scalars and nested blocks, but nothing can append to a list while refusing to shorten, reorder, or edit what is already recorded. Objective evidence also has no writer.

## Context Files

- `internal/data/write.go`
- `internal/data/write_test.go`
- `internal/data/issue_v2.go`
- `internal/data/issue_v2_test.go`
- `internal/data/project.go`
- `internal/data/errors.go`
- `internal/data/objective_v2.go`
- `internal/data/evidence_v2.go`

## Acceptance Criteria

- [ ] `CreateIssueV2` allocates the next unused `I###` over active records, writes create-only to `.savepoint/issues/{id}-{slug}.md`, validates the marshalled content decodes as an Issue before creating the file, and refuses an existing path rather than overwriting it.
- [ ] A rejected Issue creation leaves no file behind.
- [ ] Managed Issue writes patch only `status`, `resolution`, `duplicate_of`, `tasks`, `checks`, and `severity`, preserving every other YAML key, unknown fields, and the authored Markdown body unchanged.
- [ ] A history write appends entries only: a write whose entry list is shorter than, reorders, or edits any already-recorded entry is refused with a named diagnostic and leaves the file untouched.
- [ ] Appending a history entry preserves every earlier entry's fields byte-for-byte in the rewritten frontmatter.
- [ ] `WriteObjectiveEvidenceV2` patches an Objective's `last_check`, `freshness`, `owner_validation`, `exception`, and `replan` through the same patch builders `WriteTaskEvidenceV2` uses, with no second evidence contract.
- [ ] Every write validates the patched content decodes before any file is replaced, so a malformed value is refused and the file is left untouched.
- [ ] Writing values a record already has is a no-op: the file's bytes and modification time are unchanged.
- [ ] Preserving writes retain the supported line-ending form, matching existing V2 write behavior.

## Implementation Plan

- [ ] Add `NewIssueV2`, `CreateIssueV2`, and `nextV2IssueID` to `write.go`, reusing the create-only file sequence and next-unused allocation pattern `CreateCheckV2` established.
- [ ] Add `WriteIssueV2` for managed field patches through `writeV2Record`, validating with `DecodeIssueV2`.
- [ ] Add an append-only history patch that compares the incoming entry list against the recorded one and refuses any non-append change before writing.
- [ ] Add the named append-only violation sentinel to `errors.go`.
- [ ] Generalize the evidence patch builder so `WriteObjectiveEvidenceV2` reuses it, validating with `DecodeObjectiveV2`.
- [ ] Test create-only collision, validation failure leaving no file, managed patches preserving unknown fields and bodies, append-only acceptance and every refusal form, Objective evidence patches, no-op stability, and line-ending preservation.
- [ ] Run the focused `internal/data` suite.

## Context Log

Pending.
