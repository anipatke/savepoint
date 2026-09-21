---
id: O014
title: Give the Next area an Objective word and let the router target Issues
status: planned
depends_on: [O012]
release: R006
---

# O014: Give the Next area an Objective word and let the router target Issues

## Outcome

The board's Next area states a lifecycle word for whichever record `data.Next`
is pointing at — a Task or an Objective alike — instead of only a Task, and
the router's context selection can name a specific Issue the same way it
already names a Release, Objective, or Task.

## Why

Design.md's own Next-area contract (Section 8) says the panel shows "its own
Task or Objective — its lifecycle word ... its identity, and its title." The
actual renderer (`internal/board/v2/next_panel.go`) only builds a word for the
Task branch; an Objective at the ready-for-Check rung prints just its ID and
title, silently breaking that documented parity. This was noticed live while
handing O012 off to its mandatory Full Objective Check, where the Next area
read "O012 — Make Task cards easier to scan" with nothing telling the owner
whether it meant "still building" or "ready for Check."

Separately, the router's selection type (`RouterSelectionV2`) has fields for
Release, Objective, and Task, but nothing for an Issue — there is no way to
point the board or `savepoint resume` at one specific Issue the way the other
three record kinds can already be targeted.

## Success Conditions

- `nextLines`/`renderNext` show one lifecycle word for `next.Objective` the
  same way they already do for `next.Task`: `Planned` (not started), `Build`
  (in progress, an owned Task still needs work), `Check` (every owned Task
  done, ready for the mandatory Full Objective Check), or `Done` — derived
  from the Objective's own status and completion resolution, never borrowing
  a Task-only stage word (`Test`, `Audit`).
- The non-TTY plain rendering states the same word, so piping the board and
  reading it still agree, matching the existing Task-word parity guarantee.
- The router's `## Current state` block can name a single Issue (`issue:
  I###`), decoded and validated with the same strictness and the same "none"
  sentinel convention `release`/`objective`/`task` already use; an existing
  router with no `issue:` key keeps parsing exactly as before (backward
  compatible).
- Selecting an Issue through the router is readable by the board/resume
  without inventing a second selection mechanism alongside
  `data.ResolveSelection`.
- Existing Release/Objective/Task selection behavior, the `ResolveNext`
  ladder's precedence, and non-TTY parity are unchanged for any project that
  never selects an Issue.
- Active documentation (Design.md Section 8's Next-area contract) is
  reconciled to state the implemented word vocabulary precisely.
- Focused board/data tests, `git diff --check`, `make build`, and `make test`
  pass before handoff.

## Architectural Considerations

- `internal/data` remains the sole owner of Objective status/completion
  resolution and Issue lifecycle; this Objective adds presentation (the
  Next-area word) and a router selection field, not a new gate or a second
  status vocabulary.
- `internal/board/v2/badges.go` and `next_panel.go` remain the single
  translation from typed `data` values to glyph/label/accent/word; no second
  word-mapping for Objectives elsewhere in the package.
- Whether a selected Issue can itself become the resolved `data.Next` (a new
  rung) or stays a pure view-context the way Release selection narrows scope
  without being "the" next action is an open interface question — settle it
  during readiness/Task planning for this Objective, not here.
- Keep `router.md`'s YAML shape backward compatible: an existing V2 router
  with no `issue:` key must keep decoding exactly as it does today.

## Boundaries

**In scope:**

- The Objective lifecycle word in the Next area (`internal/board/v2/next_panel.go`,
  and `internal/data/next.go` if the word needs a typed source rather than
  being derived ad hoc in the board).
- Router selection support for a single Issue (`internal/data/router_v2.go`,
  `internal/data/write.go`'s `RouterSelectionV2`): decode, validate, write,
  and the "none" sentinel, plus whatever minimal board surface is needed to
  set it.
- Design.md Section 8 reconciliation for the exact word vocabulary.

**Out of scope:**

- Changing Issue lifecycle states, severity, or type vocabulary.
- Changing the existing Task lifecycle words (`Build`/`Test`/`Check`/`Planned`/`Done`).
- Any Objective, Task, Release, or Issue gate or completion-policy change.
- Renaming Release to Goals (O013's scope) or any other unrelated Next-area
  or router surface not named above.
