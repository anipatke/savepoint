---
id: I-048
title: The board cannot close an Objective whose Check is CLEAR
type: defect
status: open
source:
  kind: report
  actor: {role: owner, session: owner-chat-20260924}
  at: '2026-09-24T10:30:00Z'
checks: [C-917, C-918]
guardrail_ids: [TPL-02]
severity: medium
history:
  - at: '2026-09-24T10:30:00Z'
    actor: {role: owner, session: owner-chat-20260924}
    kind: observed
    note: >-
      After C-918 cleared O-023 and Next read "Close O-023", the owner asked
      "i cant mark an obj done??". The board offers no ordinary Objective
      close; O-014 and O-023 were closed by hand-editing their status and the
      router.
---

# I-048: The board cannot close an Objective whose Check is CLEAR

## Summary

When every owned Task is done and the Objective's latest Check is CLEAR,
`savepoint resume` and the board's Next both say `Close O-###`, but nothing
on the board performs that closure. The only Objective lifecycle write the
board has is closure by recorded exception (`x`). Design section 8 already
promises that "Closing the selected Objective clears `objective` and `task`",
which no board path does for an ordinary, cleared Objective.

## Evidence

- `internal/board/v2/actions.go:48` `actionsForRecord` offers only `x`
  (requires `AllowedByException`) and `a` (owner acceptance block). A cleared
  Objective with `decision.Allowed` and no exception gets no action.
- `internal/board/v2/io.go:100` `writeExceptionCompletionCmd` is the only
  code that sets `objective.Status = done`; it refuses unless
  `AllowedByException`.
- `internal/board/v2/update.go:119` Space and Backspace act only on a
  focused Task card, never on a focused sidebar Objective.
- `.savepoint/Design.md` section 8 (line 183) describes an Objective closure
  that clears the router's `objective` and `task`.
- Repro: O-023 with T-035 done and C-918 CLEAR. Next reads
  `Close O-023 — ...`; focusing O-023 in the sidebar shows only `p` in the
  help's OWNER ACTIONS, and Space does nothing. Commits `b13955b` (O-014) and
  `95bfa4e` (O-023) closed each Objective by hand-editing.

## Proposed Repair

Mirror Task completion. Space on a focused sidebar Objective closes it:

1. Re-resolve `data.ResolveObjectiveCompletion` against a fresh index.
2. If `Allowed` with owner authority, write `status: done` through
   `data.WriteObjectiveV2`, then call `completedRecordAction` with
   `closed = {Objective: id}` so `RouterSelectionAfterClosureV2` clears
   `objective`/`task` and keeps `release:` byte-for-byte.
3. Otherwise refuse with `decisionRefusal` (unfinished Task, missing,
   NEEDS WORK, stale or unknown Check, owner acceptance pending, open linked
   Issue). An already-done Objective reports "already done".
4. Leave `x` as the exception path; a gate that is allowed only by exception
   keeps requiring `x`.
5. List the close in the help's OWNER ACTIONS when the gate allows it, so
   the key is visible only when it will work.

No `internal/data` policy changes; the gate already decides.

## Proof Needed

- Board tests: a cleared Objective closes on Space, becomes done, and the
  router clears `objective` and `task` while `release:` and `issue` are
  unchanged.
- Refusal tests: an unfinished Task, no Check, NEEDS WORK, owner acceptance
  pending, an open linked Issue, and an exception-only gate each leave the
  Objective and router untouched and show the refusal.
- A router write conflict after the Objective write reports the stale
  selection, as Task closure does.
- Help lists the Objective close only when it is allowed.
- `make build && make test-fast` pass.
