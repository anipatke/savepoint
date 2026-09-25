---
id: T-045
title: Move an Issue forward or back in its file
objective: O-015
status: done
complexity_tier: small
complexity_reason: Adds one Issue writer beside the existing Task/Objective writers, reusing writeV2Record, plus one transition function over the existing Issue vocabulary.
depends_on: []
owner_validation:
    required: false
    accepted_check: ""
planned_by: {role: planner, session: o015-plan-20260925}
check_waiver:
    task: T-045
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-09-25T09:11:53Z"
---

# Move an Issue forward or back in its file

## Outcome

`internal/data` can move an Issue one status forward or back, recording the
owner's resolution or reopening and one history entry, without disturbing
anything else in the file.

## User Check

None beyond tests; the owner sees this through the board in the next Task.

## Done When

- A data function (for example `AdvanceIssueV2(index, id, actor, at)` and
  `RetreatIssueV2(index, id, actor, at)`, or one function with a direction)
  applies exactly the four transitions in O-015's Success Conditions and
  returns a named error for Space on resolved, Backspace on open, or an
  unknown ID, writing nothing.
- Resolve writes `resolution` with `disposition: accepted`, owner actor, time,
  and reason `Resolved by the owner from the board.`, and no `check`.
- Reopen removes `resolution`, `duplicate_of`, and `escalated_to`, and its
  `reopened` history note names the removed disposition.
- Each transition appends exactly one history entry (`owner_decision` or
  `reopened`) after the existing entries, which stay byte-identical.
- The write goes through `writeV2Record`: body and unknown frontmatter are
  preserved, the patched content decodes as `IssueV2` and passes the
  resolution rules before replacement, and a file changed on disk is refused
  with the existing source-conflict error.
- Tests in `internal/data/issue_v2_test.go` or a new
  `internal/data/issue_write_test.go` cover all four transitions, both
  refusals, history append-only, body/unknown-field preservation, a
  verified-then-reopened Issue, and a stale-file conflict.

## Context Files

`internal/data/issue_v2.go`, `internal/data/issue_v2_test.go`,
`internal/data/write.go`, `internal/data/write_test.go`,
`internal/data/errors.go`, `internal/data/evidence_v2.go` (Actor encoding),
`.savepoint/objectives/O-015-owner-advance-issues/Objective.md`.

## Design References

Design section 1 (V2 follow-up and integration boundary) and the Issue
lifecycle rules.

## Guardrails

FS-01, FS-04, DATA-01, DATA-02, DATA-03, ARCH-02, TEST-01..04, TEST-08,
STYLE-07.

## Implementation Plan

1. Add an Issue record writer in `write.go` (or a new `issue_write.go`) that
   builds `v2FieldPatch` entries for `status`, `resolution`, `duplicate_of`,
   `escalated_to`, and `history`, using `encodeV2Node` for the structured
   blocks and `Remove` for cleared keys.
2. Validate with `DecodeIssueV2` plus the per-record resolution rule
   (resolution present exactly when resolved; accepted obligations), since
   the index-wide validator is not available per file.
3. Add the transition function that looks up the Issue, picks the next or
   previous status, builds the resolution/history change, and calls the
   writer. Keep the reason string and board session in one place.
4. Update the stale `WriteIssueV2`/`WriteIssueHistoryV2` mention in the
   `NewIssueV2` comment to name the real function.
5. Write the tests listed in Done When.

## Boundaries

No board changes, new dispositions, new history kinds, or status vocabulary
changes.

## Technical Verification

`make test-focused TEST=...` while iterating; `make build && make test-fast`
at handoff.

## Technical Evidence

Per-criterion evidence:

- All four one-step transitions and terminal/unknown-ID refusals are covered
  by `TestAdvanceAndRetreatIssueV2_transitions` and
  `TestAdvanceAndRetreatIssueV2_refusalsAndUnknownIDDoNotWrite`.
- Resolve records accepted, owner, action time, the fixed reason, and no Check;
  the `in_progress_to_resolved` subtest checks those fields.
- Reopen clears resolution and disposition targets, and names the removed
  disposition. The verified, duplicate, and escalated reopen subtests cover
  those cases.
- Each transition appends one attributed entry and preserves the decoded
  existing history entries. The transition table checks append-only order,
  body preservation, unknown frontmatter, and an unknown history field.
- `writeIssueTransitionV2` uses `writeV2Record`; its validation decodes the
  patched Issue, enforces resolution/status consistency, and checks accepted
  resolution obligations. `TestAdvanceIssueV2_refusesStaleIssueFile` proves a
  disk conflict leaves the external bytes intact.

Commands and results:

- `go test ./internal/data -run 'TestAdvanceAndRetreatIssueV2|TestAdvanceIssueV2' -count=1` — passed after correcting the encoded YAML node shape. The first focused run exposed that `encodeV2Node` returns a mapping node rather than a document node.
- `make build && make test-fast` — passed on the final run. The initial run built successfully but `test-fast` failed on the same node-shape assumption; the failure was fixed before this passing run.
- `git diff --check` — passed.

Files read: `agent-skills/savepoint-task/SKILL.md`, `.savepoint/router.md`,
the T-045 Task and O-015 Objective, all T-045 Context Files, and
`.savepoint/Guardrails.md` rules named by the Task. Extra reads: a targeted
symbol search across `internal/data`, then definitions in
`internal/data/project.go`, `internal/data/parser.go`, and
`internal/data/check_v2.go` for `V2Index`, source freshness metadata, and
owner actor types.

Files changed: the T-045 Task evidence/lifecycle, O-015 status, and
`internal/data/errors.go`, `internal/data/write.go`, and
`internal/data/issue_v2_test.go`.

Limitations: board key handling, reload/focus behavior, and footer/help text
remain for the other O-015 Task; they were not verified here. No Task Check
waiver has been recorded. Stage `audit` means this Task is ready for a Check;
it does not claim `CLEAR`.

Extra reads: a targeted symbol search across `internal/data` located the
definitions of `V2Index`, `V2SourceDocument`, and `Actor`; the relevant
definitions were read in `internal/data/project.go`, `internal/data/parser.go`,
and `internal/data/check_v2.go`. These were needed to implement the
index-based Issue transition using the existing source freshness metadata
and actor types.

## Post-C-930 Repair Evidence

- Direct repair for I-057: `writeV2Record` now uses indentation inferred from
  retained source YAML node locations, so a managed write preserves the
  source document's nested block indentation. Added
  `TestAdvanceIssueV2_preservesExistingHistoryBytes`, which checks the exact
  original 2-space history block after an Issue transition.
- Extra reads for this repair: `.savepoint/router.md`, C-930, I-057, the
  sibling T-046 record, `agent-skills/savepoint-task/SKILL.md`,
  `agent-skills/references/issue-capture.md`, `.savepoint/Guardrails.md`, and
  `internal/data/parser.go`. These were needed to follow repair routing,
  preserve the original Check's scope, confirm the promised byte-level
  behavior, and use source-node line/column values for indentation detection.
- Focused command
  `go test ./internal/data -run 'TestAdvanceIssueV2_preservesExistingHistoryBytes|TestAdvanceAndRetreatIssueV2_transitions|TestWrite.*Preserves' -count=1`
  passed after rerunning with Go cache access. `make build && make test-fast`
  passed.
- Repair files: `internal/data/write.go`,
  `internal/data/issue_v2_test.go`, and I-057. Task status remains `done`;
  I-057 remains open for independent verification.

## Drift Notes

None expected.
