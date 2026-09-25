---
id: E44-issues-objective-checks/T002-connect-issues-to-their-work
title: Connect issues to the work they came from
status: done
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

- [x] Every `I###` reference in a Check's `issues` resolves to an existing Issue; a missing target returns a named diagnostic identifying the Check and the reference.
- [x] Every `T###` in an Issue's `tasks` and every `C###` in its `checks` resolves to an existing record; a missing target returns a named diagnostic identifying the Issue and the reference.
- [x] The immutable Check is the authoritative direction: a Check naming an Issue that does not appear in that Issue's `checks` list returns a named diagnostic naming both records, and the load fails closed rather than reconciling either side.
- [x] An Issue may list a Check that does not name it — a recheck or proof Check that observed the Issue without recording it — and that asymmetry is legal.
- [x] `V2Index` exposes a Task ID to Issue IDs map and a Check ID to Issue IDs map, both in deterministic ID order.
- [x] Derived listings answer open, in_progress, and resolved Issue sets and counts, and per-type counts, from the index alone with no stored summary, register, or cached total anywhere in the project files.
- [x] Listings are computed from the loaded index on each call and never persisted to disk.
- [x] A project whose Issues reference no Tasks or Checks loads with empty link maps and no error.

## Implementation Plan

- [x] Add link resolution to `LoadV2Index` in `project.go` after Issue discovery: resolve Check `issues`, Issue `tasks`, and Issue `checks`, then assert the Check-authoritative pairing rule.
- [x] Add the named link sentinels to `errors.go`, distinguishing a missing target from an unpaired Check-to-Issue link.
- [x] Build the Task-to-Issue and Check-to-Issue maps during index construction, sorting IDs for deterministic output.
- [x] Add derived listing and count helpers to `issue_v2.go` operating on `*V2Index`, returning sorted IDs.
- [x] Update the doc comment on `CheckV2.Issues` in `check_v2.go` so it no longer says the references are never loaded or validated.
- [x] Test resolved links, missing Issue/Task/Check targets, the unpaired-link diagnostic in both orderings, the legal asymmetric case, map ordering, and listing and count results across mixed statuses and types.
- [x] Run the focused `internal/data` suite.

## Context Log

**Read:** `.savepoint/router.md`, `AGENTS.md`, `E44-Detail.md`, this task, `T001-record-durable-follow-up.md` (what it deliberately left to T002), `.savepoint/Guardrails.md` STYLE rules, and the listed context files in `internal/data`.

**Edited:** `internal/data/project.go`, `internal/data/project_test.go`, `internal/data/issue_v2.go`, `internal/data/issue_v2_test.go`, `internal/data/errors.go`, `internal/data/check_v2.go`.

**Decisions made during build.**

- Two sentinels, not four. `ErrV2IssueMissingLinkTarget` covers all three dangling directions (Check→Issue, Issue→Task, Issue→Check) because a caller's response to any of them is the same — repair the reference — while `ErrV2IssueUnpairedCheckLink` is separate because its repair is a judgement about which side is wrong. Each message names the owning record and the reference, so the direction is readable without a distinct sentinel per field.
- Both link maps are built from the Issue side only. The pairing rule guarantees every Check-side link also appears in that Issue's `checks`, so walking Issues alone captures every link plus the legal asymmetric ones, and the map has one derivation rather than two that could disagree (STYLE-07).
- Repeated references are listed once. A record may legally write `tasks: [T001, T001]`; `appendUniqueIssueID` drops the repeat, and a tail check suffices because one Issue's references are walked together and Issues are walked in ascending ID order.
- Listings are methods on `*V2Index` declared in `issue_v2.go`, returning fresh slices and maps on each call. Counts pre-seed every status and type at zero so a caller never reads an absent key as none.
- `issueStatuses`/`issueTypes` now hold the two vocabularies once; `decodeIssueStatus` and `decodeIssueType` read them instead of restating them in a switch, so the counts cannot drift from what decoding accepts (STYLE-07, STYLE-09). Diagnostics are unchanged.
- Issue IDs sort lexically via `slices.Sorted`, matching the ordering T001 established for Issue diagnostics. The numeric `sortedV2CheckIDs` ordering stays Check-only; Check lists are walked in that recorded order so link diagnostics follow the same order as the existing Check diagnostics.

**Acceptance verification.** 10 new tests (14 cases with subtests), all passing:

- Resolution — `TestLoadV2Index_issueLinksResolve` (a fully paired link lands in both maps), `TestLoadV2Index_checkNamesMissingIssue` (names `C001` and `I999`), `TestLoadV2Index_issueNamesMissingRecord` (missing task, missing check; each names the Issue and the reference).
- Authority direction — `TestLoadV2Index_checkIssueLinkUnpaired` runs both orderings (`C001`→`I002` and `C002`→`I001`), asserting the sentinel, both records named, and a nil index alongside the diagnostic.
- Legal asymmetry — `TestLoadV2Index_issueMayNameCheckThatDoesNotNameIt` proves a recheck Check that never recorded the Issue loads, and still appears in `CheckIssues`.
- Maps — `TestLoadV2Index_issueLinkMapsAreOrdered` (three Issues written out of order plus a repeated `T001` reference yield `[I001 I002 I003]` once each in both maps), `TestLoadV2Index_issueLinkMapsEmptyWithoutReferences`, `TestLoadV2Index_issueLinkDiagnosticOrder` (identical error across runs, lowest ID first).
- Listings and counts — `TestV2Index_issueListingsAndCounts` (four Issues across three statuses and three types; open/in_progress/resolved sets, status counts, per-type counts), `TestV2Index_issueCountsCoverEveryVocabularyValue` (empty index reports every status and type at zero), `TestV2Index_issueListingsAreDerivedNotPersisted` (project file set and mtimes unchanged after three listing calls, and a status changed in the index changes the next answer).

**Quality gates.** `go test ./internal/data/` ok. `make build` ok. `make test` ok — all packages pass. `go vet ./...` clean. `gofmt -l` reports only `router.go` and `task_test.go`, both already unformatted on `v2` before this change and untouched here.

**Not done here, by design.** Resolution proof obligations (a `verified` Issue's proof Check existing, being CLEAR, and appearing in `checks`), reopening, and deferral rules are T003 — `checks` membership is validated for existence here, not for what a resolution requires of it. Issue writes and append-only history are T004. Objective gates are T005–T006. Doctor reporting of the new diagnostics is T007. `.savepoint/Health-Check.md` is absent from this project, so the Quick check step was skipped.
