---
id: T-031
title: Let the router select an Issue to work on
objective: O-014
planned_by: {role: planner, session: o014-design-20260924}
status: planned
complexity_tier: medium
complexity_reason: "New optional router key through decode, write, and selection resolution with backward compatibility, plus context display on two surfaces."
depends_on: [{task: T-030, requires: clear}]
owner_validation: {required: false}
---

# T-031: Let the router select an Issue to work on

## Outcome

The router can select an Issue, alone or alongside an Objective/Task. An
Issue selected alone is the Next line, for example
`In Progress I-042 — <title>`; alongside a Task it is shown as
context.

## User Check

Ask the agent to "set router to I-042": the board, non-TTY output, and the
first line of `savepoint resume` all read `<Issue word> I-042 — <title>`.
Then set the router to O-014 T-028 with `issue: I-042`: Next is the T-028
line, and I-042 appears as context.

## Done When

- `ReadStateV2` decodes an optional `issue` key with the same shape check
  and `none` sentinel as `objective`/`task`; an Issue alone (no Objective)
  is valid. A router without the key decodes exactly as before.
- `RouterSelectionV2` gains `Issue`; `validate` and `WriteRouterStateV2`
  write it without disturbing other keys, adding the key only when a value
  is set or the key already exists, so untouched routers stay byte-identical.
- `ResolveSelection` resolves it into `Selection.Issue`; an unknown Issue
  yields `SelectionNotFound` with record kind `issue`; a resolved Issue
  yields T-030's `SelectionDone`.
- With only an Issue selected, `ResolveNext` returns an Issue Next kind
  carrying the Issue; with an Objective/Task also selected, T-029's rules
  apply and the Issue rides along on `Next` as context.
- T-028's line builder adds the Issue shape `<Open|In Progress|Resolved>
  I-### — <title>`; Issues carry no stage, so there is no `·` part. Board, non-TTY, and resume first line stay byte-equal.
- Resume adds an Issue context line when the Issue rides alongside a Task.
- Tests: decode (present, absent, `none`, malformed, Issue alone), write byte
  preservation, not-found and resolved diagnostics, each line shape, and
  Next unchanged for routers without `issue:`.
- `git diff --check`, `make build && make test-fast` pass.

## Context Files

`internal/data/router_v2.go`, `internal/data/router_v2_test.go`,
`internal/data/write.go`, `internal/data/write_test.go`,
`internal/data/next.go`, `internal/data/next_test.go`,
`internal/resume/resume.go`, `internal/resume/resume_test.go`,
`internal/data/issue_v2.go`, `internal/board/v2/next_panel.go`,
`internal/board/v2/next_panel_test.go`, `internal/board/v2/plain.go`,
`internal/resume/evidence.go`, `main_board_next_parity_test.go`.

## Design References

Design sections 1 and 8.

## Guardrails

DATA-01, DATA-03, FS-01, STYLE-07, TEST-02, TEST-03, TEST-08.

## Implementation Plan

1. Extend the router frontmatter and state types; validate shape.
2. Extend the writer and its tests, including byte preservation.
3. Resolve the Issue through `Selection`; add the Issue-only Next kind.
4. Extend the line builder and resume; add tests.

## Boundaries

No board key to select or advance an Issue (O-015); selecting one is the
owner's request to an agent or a file edit. No Issue lifecycle change.

## Technical Verification

Focused data/resume/board tests during iteration; `make build && make
test-fast` at handoff.

## Technical Evidence

Pending execution.

## Drift Notes

None expected.
