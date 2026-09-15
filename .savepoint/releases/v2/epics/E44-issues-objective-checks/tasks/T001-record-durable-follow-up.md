---
id: E44-issues-objective-checks/T001-record-durable-follow-up
title: Record durable follow-up
status: planned
objective: Load Issue records into the project index with strict decoding, confined discovery, and a validated duplicate graph.
depends_on: []
complexity_tier: high
complexity_reason: Adds a fourth record family across decoding, confined discovery, and index graph validation.
---

# T001: Record durable follow-up

## Problem

Follow-up that outlives a single evaluation has nowhere to live. A Check can name `I###` references, but no Issue record exists for them to resolve against, so a repeated observation cannot keep its identity and a duplicate cannot point anywhere.

## Context Files

- `internal/data/issue_v2.go`
- `internal/data/issue_v2_test.go`
- `internal/data/check_v2.go`
- `internal/data/errors.go`
- `internal/data/parser.go`
- `internal/data/evidence_v2.go`
- `internal/data/discover.go`
- `internal/data/discover_test.go`
- `internal/data/project.go`
- `internal/data/project_test.go`
- `internal/data/dependency.go`

## Acceptance Criteria

- [ ] An Issue record requires a valid global `I` ID, a non-empty `title`, a `type` of `defect`, `drift`, `guardrail`, `verification`, or `other`, a `status` of `open`, `in_progress`, or `resolved`, and a `source` naming what produced it with its actor and time.
- [ ] Issue `status` is its own vocabulary: `planned`, `done`, `in_progress` with a stage, or any other Task lifecycle value is rejected with a named diagnostic and never healed into a valid-looking default.
- [ ] `tasks`, `checks`, and `guardrail_ids` decode as lists; `tasks` and `checks` are validated as `T###` and `C###` identity references for shape only, and `guardrail_ids` are opaque policy strings that are never resolved against a Guardrails file.
- [ ] Optional `severity`, `resolution`, and `duplicate_of` decode when present; an absent optional field stays absent rather than becoming a permissive default.
- [ ] `resolution` decodes its disposition (`verified`, `accepted`, `duplicate`), optional proof `check`, actor, time, and reason as a shape; cross-record proof obligations are out of scope for this task.
- [ ] `history` decodes as an ordered list of `{at, actor, kind, note, check}` entries whose `kind` is one of `observed`, `repair_attempted`, `rechecked`, `deferred`, `reopened`, `owner_decision`; a malformed or unrecognized entry is a named diagnostic.
- [ ] Issues are discovered only from `.savepoint/issues/`, through the same path confinement as Checks, rejecting traversal, symlink escape, case aliasing, duplicate IDs, and a filename that does not start with its declared ID.
- [ ] `V2Index` exposes Issues keyed by `I###`.
- [ ] `duplicate_of` resolves to an existing, different Issue; a missing target, a self-reference, or a cycle each return a distinct named diagnostic.
- [ ] A project with no `issues/` directory loads with no Issues and no error.
- [ ] Malformed Issue records fail the load closed, in deterministic ID order, exactly as E42 and E43 structural diagnostics do.

## Implementation Plan

- [ ] Add `IssueV2`, its type/status/disposition/history-kind types, and `DecodeIssueV2` in a new `issue_v2.go`, reusing `ParseV2Document` and the existing anchored ID patterns.
- [ ] Add the named Issue sentinels to `errors.go` for malformed records, invalid identity, invalid lifecycle vocabulary, and duplicate-graph conflicts.
- [ ] Add `v2IssuesDirName` and confined `.savepoint/issues/` discovery in `discover.go`, reusing `v2PathConfiner` and the existing filename/ID mismatch rule.
- [ ] Extend `V2Index` and `LoadV2Index` in `project.go` with the Issues map and duplicate-graph validation, reusing `findV2Cycle`/`describeV2Cycle` from `dependency.go`.
- [ ] Test valid records, every malformed field, Task vocabulary rejection, duplicate IDs, path mismatch, unsafe paths, history entry shapes, valid and missing and self and cyclic `duplicate_of`, absent `issues/`, and deterministic diagnostic ordering.
- [ ] Run the focused `internal/data` suite.

## Context Log

Pending.
