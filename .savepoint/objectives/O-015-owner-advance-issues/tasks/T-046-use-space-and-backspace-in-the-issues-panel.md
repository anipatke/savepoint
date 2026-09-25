---
id: T-046
title: Use Space and Backspace in the Issues panel
objective: O-015
status: done
complexity_tier: small
complexity_reason: Wires two keys in the existing Issues key handler to one new tea.Cmd, restores focus after reload, and updates footer/help and guidance wording.
depends_on: [{task: T-045, requires: clear}]
owner_validation:
    required: true
    accepted_check: C-932
    accepted_by: {role: owner, session: user}
planned_by: {role: planner, session: o015-plan-20260925}
check_waiver:
    task: T-046
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-09-25T09:26:33Z"
---

# Use Space and Backspace in the Issues panel

## Outcome

With the Issues panel open, the owner presses Space to move the selected
Issue right and Backspace to move it left, and the selection stays on that
Issue in its new column.

## User Check

Open the board, press `i` for Issues. On an Open Issue press Space: it moves
to In Progress and stays selected. Press Space again: it moves to Resolved
and its detail shows accepted by the owner with the board reason. Press
Backspace twice: it returns to Open, and the history shows every step in
order. Space on a Resolved Issue and Backspace on an Open one do nothing.
The footer and `?` help list both keys.

## Done When

- `handleIssuesKey` (list view only, not detail) maps Space and Backspace to
  a new command in `io.go` that reloads a fresh index, calls the data
  transition from the previous Task for the selected Issue ID with actor
  `{owner, board-owner}` and the current time, and returns `actionMsg` with a
  short message and `reload: true`, or the error via `actionFailure`.
- No command is scheduled for an empty column or missing selection; the
  data-layer refusals show in the status line.
- After reload the Issues cursor follows the moved Issue into its new status
  column.
- The Issues panel footer and help list `space` forward and `backspace` back.
- AGENTS.md, `agent-skills/savepoint-*` and `agent-skills/references/issue-capture.md`
  where they state Issue closure authority, and their
  `templates/project-v2/agent-skills/` copies (byte-identical, TPL-01), say
  the owner may resolve an Issue as `accepted` and reopen it from the board.
  Design section 1/8 notes the two keys.
- Board tests cover both keys moving an Issue, focus following it, the no-op
  cases, a write failure shown in the status line, and the footer/help text.
- `make build && make test-fast` passes.

## Context Files

`internal/board/v2/issues.go`, `internal/board/v2/issues_test.go`,
`internal/board/v2/issues_view.go`, `internal/board/v2/io.go`,
`internal/board/v2/update.go`, `internal/board/v2/help.go`,
`internal/board/v2/footer_test.go`, `internal/board/v2/actions.go`,
`AGENTS.md`, `agent-skills/references/issue-capture.md`,
`templates/project-v2/agent-skills/references/issue-capture.md`,
`.savepoint/Design.md`,
`.savepoint/objectives/O-015-owner-advance-issues/Objective.md`.

## Design References

Design section 1 (V2 follow-up and integration boundary) and section 8 (TUI
layout, keybindings).

## Guardrails

ARCH-02, DATA-02, FS-01, TPL-01, TPL-02, TEST-01..04, TEST-08, STYLE-07.

## Implementation Plan

1. Confirm T-045's transition function exists; return REPLAN
   REQUIRED if not.
2. Add the Issue transition command in `io.go` following
   `writeTaskAdvanceCmd`.
3. Handle `" "` and `"backspace"` in `handleIssuesKey`; make it return the
   command (adjust its caller in `update.go`).
4. Remember the moved Issue ID and restore the cursor to it in
   `refreshIssues` after reload.
5. Update the Issues footer and help copy.
6. Update guidance, templates, and Design wording; keep TPL-01 copies
   identical.
7. Write the tests listed in Done When.

## Boundaries

No prompt, reason entry, new dispositions, detail-overlay actions, or
changes to Task keys.

## Technical Verification

`make test-focused TEST=...` while iterating; `make build && make test-fast`
at handoff.

## Technical Evidence

- Start: dependency gate reported satisfied by `./savepoint resume` at 2026-09-25; router selected O-015/T-046 and Goal R-006.
- Extra reads: `.savepoint/router.md` for routing; `.savepoint/Guardrails.md` for ARCH-02, DATA-02, FS-01, TPL-01/02, TEST-01..04, TEST-08, STYLE-07; targeted `internal/data` transition API search and `internal/data/write.go` excerpt for implementation plan step 1 (confirmed `AdvanceIssueV2` and `RetreatIssueV2`); targeted board symbol search for reload selection; `internal/board/v2/view.go` for footer ownership; `internal/board/v2/columns_view_test.go` helper signatures for board command tests; `agent-skills/savepoint-design/SKILL.md`, `savepoint-task/SKILL.md`, `savepoint-check/SKILL.md` and matching `templates/project-v2/agent-skills/` copies plus shared Issue-capture references for closure-authority wording required by Done When; targeted `internal/init/agent_skills_test.go` read and assertion updates after `make test-fast` exposed obsolete exact-wording expectations.
- Acceptance evidence:
  1. Added an Issues transition command that loads a fresh index, invokes `AdvanceIssueV2` or `RetreatIssueV2` as `{owner, board-owner}` with the current UTC time, and returns `actionMsg` with status text, reload, and the selected Issue ID. Space/Backspace dispatch only in list mode.
  2. Empty columns and missing selections schedule no command. Data refusals for Space on Resolved and Backspace on Open are shown in the status line without file changes; tested both.
  3. Successful actions carry the moved ID through reload; the board test moves I-001 through all four transitions and confirms focus follows its status column and cursor.
  4. The Issues footer and context help list Space and Backspace; footer/help assertions pass.
  5. Updated `AGENTS.md`, the active Design/Task/Check skills and Issue-capture reference, and their scaffold copies to describe owner resolution as `accepted` and reopening from the board. Updated `.savepoint/Design.md` sections 1 and 8 with the keys. Live/scaffold skill/reference pairs compare byte-identically.
  6. Board tests cover both keys, append-only owner history through a full cycle, focus restoration, boundary refusals, empty/missing selection, detail-mode exclusion, write failure in the status line, and footer/help text.
  7. `make build && make test-fast` passed on the final run. Two earlier runs exposed stale exact-wording expectations in `internal/init/agent_skills_test.go`; those expectations were updated to assert the new authority wording.
- Files read: Task Context Files (`internal/board/v2/issues.go`, `issues_test.go`, `issues_view.go`, `io.go`, `update.go`, `help.go`, `footer_test.go`, `actions.go`, `AGENTS.md`, active/template `issue-capture.md`, `.savepoint/Design.md`, and the owning Objective); required route/guardrail reads and all extra reads are listed above.
- Files changed: this Task record; `AGENTS.md`; `.savepoint/Design.md`; active and scaffolded Issue-capture, Design, Task, and Check guidance; `internal/board/v2/issues.go`, `io.go`, `update.go`, `help.go`, `view.go`, `issues_test.go`, `footer_test.go`; `internal/init/agent_skills_test.go`.
- Limitations: no independent Task Check or owner waiver was requested or recorded. The Task is ready for a fresh optional Task Check; the mandatory Full Objective Check remains required before O-015 can close.

## Post-C-930 Repair Evidence

- Direct repair for I-056: updated the four Issue-authority passages in
  `templates/project-v2/AGENTS.md` to say the owner may resolve as `accepted`
  with Space and reopen with Backspace, matching `AGENTS.md`.
- Extra reads for this repair: `.savepoint/router.md`, C-930, I-056, and
  `templates/project-v2/AGENTS.md`. These were needed to preserve the Check's
  scope, locate the stale scaffold passages, and compare them with the live
  guidance.
- `make build && make test-fast` passed. Task status remains `done`; I-056
  remains open for independent verification.

## I-056 Follow-up Repair Evidence

- Restored the scaffold's rule that an agent may record an owner decision
  only when directly instructed. Rewrote its Issue Capture bullet with the
  repository's checker/owner/planner wording, removing the checker-only
  closure contradiction, and aligned the four passages on “Issues panel.”
- Extra read: `templates/project-v2/AGENTS.md` was re-read to apply the
  review's exact wording corrections; `AGENTS.md` was re-read for the
  matching canonical passages.
- `make build && make test-fast` passed after the wording correction. Task
  status remains `done`; I-056 remains open for independent verification.

## Owner Acceptance

- `2026-09-25T10:20:38Z`: after C-932 recorded CLEAR, the owner said in
  conversation "mark O015 as done, commit and push all changes". C-932 is
  CLEAR and supersedes C-931. On that instruction the agent recorded
  `owner_validation.accepted_check: C-932` and set O-015 `status: done`.

## Drift Notes

None expected.
