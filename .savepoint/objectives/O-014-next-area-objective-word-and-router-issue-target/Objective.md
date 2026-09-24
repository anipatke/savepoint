---
id: O-014
title: Give the Next area an Objective word and let the router target Issues
status: planned
depends_on: [O-012]
release: R-006
---

# O-014: Give the Next area an Objective word and let the router target Issues

## Outcome

The Next line is the router's selection, stated in one structured,
copy-pasteable form — for example `In Progress O-014 · Build T-028 — <title>`.
The owner can paste that line into a fresh agent session to resume work. The
board, the non-TTY output, and the first line of `savepoint resume` print the
identical text. When the owner closes a Task on the board, the router moves to
the Objective's next Task in the same write, so the line stays right without
an agent hand-editing the router. The router's free-text `next_action` is
retired.

## Why

Design.md Section 8 says the Next area shows "its own Task or Objective — its
lifecycle word ... its identity, and its title." `internal/board/v2/next_panel.go`
builds a word only for a Task, so an Objective at the ready-for-Check rung
prints just its ID and title. This was noticed while handing O-012 to its
Full Objective Check: Next read "O-012 — Make Task cards easier to scan" with
nothing saying whether it was still building or ready for Check.

The router selection (`RouterSelectionV2`) can name a Release, Objective, or
Task, but not an Issue, so there is no way to point the board or
`savepoint resume` at one specific Issue.

Next ignores the router's selected Objective unless the router also names an
unfinished Task or the Objective is waiting on its integration Check. On
2026-09-23, after O-016 closed, the router selected O-018, which had no Tasks,
and Next still showed "Planned T-006 — Preserve existing project records"
from O-013. Task IDs are global, so planning O-018's Tasks would not have
fixed the order. The owner added this to O-014 on 2026-09-23.

The router and the board state two different next actions. The board's Next
is computed by `data.ResolveNext`; the router's `next_action` is hand-written
prose no V2 code renders, yet agents act on it because AGENTS.md sends them to
the router first. On 2026-09-23, after the owner closed T-020 and O-018
(C-910), the router still selected `objective: O-018, task: T-020` and asked
for O-018's Full Objective Check, while the board's Next silently fell through
to the lowest-ID in-progress Task in R-006. No role owns advancing the router
after an owner closure, and no surface says the router selects finished
records. The owner added this to O-014 on 2026-09-23.

## Confirmed Design Decisions

Confirmed by the owner on 2026-09-24:

- **Split.** Mandatory Goal context (router must name a declared Goal; every
  live Objective must reference one) moves to O-022, which applies it to every
  Savepoint project. O-014 keeps today's optional-Goal behavior and must not
  add new code paths that assume a Goal is absent.
- **Next is the router's selection; no record search.** (Revised
  2026-09-24, replacing an earlier in-Objective precedence ladder, to avoid
  over-engineering.) Next states exactly what the router selects: the
  Objective and, when selected, its Task, each with its word. The
  project-wide and Release-wide searches for "lowest-ID ready work" are no
  longer a source of Next. Gate state for the selected record (Check needed,
  owner wait, unmet dependency, replan) still comes from the existing
  resolvers. With no Objective selected, Next says nothing is selected. The
  Goal completion rungs for a selected Release whose Objectives are all done
  are unchanged.
- **Line format.** `<Objective word> O-### · <Task word> T-### — <Task
  title>`; with no Task selected, `<Objective word> O-### · Check — <Objective
  title>` when every owned Task is done, otherwise `<Objective word> O-### —
  <Objective title>`. Objective word: `Planned`, `In Progress`, or `Done`.
  Task word: the existing `Planned`/`Build`/`Test`/`Check`/`Done`. Plain text,
  no glyphs, so it pastes cleanly.
- **Owner acceptance wait reads `Check`.** An Objective whose Check is CLEAR
  but awaits owner acceptance still reads `· Check`; there is no extra word.
- **An Issue can be the selection.** (Revised 2026-09-24; replaces "Issue
  selection is context only".) The router may select an Issue on its own,
  with no Objective or Task; Next is then that Issue:
  `<Issue word> I-### — <title>`, Issue word `Open`, `In Progress`, or
  `Resolved` (Issues carry no stage). When an
  Objective or Task is also selected, their line wins and the Issue is shown
  as context (the Task is the repair). A selected resolved Issue gets the
  stale-selection warning. Board Issue lifecycle actions stay with O-015.
- **The board advances the router on closure; the Release is never
  blanked.** When the owner closes the selected Task on the board (Space to
  done, including the waiver and exception paths), the same action points the
  router at the Objective's lowest-ID Task that is not done (a blocked one
  shows its wait through its own gate); when every Task is done, it clears
  `task`, and Next reads `· Check`. Closing the
  selected Objective clears `objective` and `task`. Neither write changes
  `release:` or `issue:`. Skills set the selection in every other case —
  `savepoint-design` when it selects the next Objective, `savepoint-task` when
  it starts a Task — and the stale-selection diagnostic catches any miss,
  such as a Task closed by hand-editing its file.

## Success Conditions

- The board Next area, non-TTY output, and the first line of `savepoint
  resume` print the identical Next line in the format above; the board may
  colour the words but the text is the same.
- Next is derived only from the router's selection plus the existing gate
  resolvers for the selected records; no search picks other work. The
  2026-09-23 cases (router on O-018 showing T-006; router on done T-020
  showing another Task) cannot recur.
- With no Objective selected, Next says nothing is selected and resume says
  to select an Objective; it never picks one. The line tells the owner how:
  press `p` on the board, or ask the agent to "set router to O-### T-###".
- An owner request such as "set router to O-014 T-028" is a documented agent
  action: the agent confirms the records exist and the Task belongs to the
  Objective, edits only the `objective`/`task`/`issue` keys (an Issue may be selected
  alone, clearing `objective`/`task`), never
  changes `release:` unless the owner names a Release, and then runs
  `savepoint resume` to show the resulting Next line.
- The router's `## Current state` block accepts `issue: I-###`, decoded and
  validated with the same strictness and `none` sentinel as the Objective/Task
  selection, and written by `WriteRouterStateV2` without disturbing other
  keys. A router with no `issue:` key decodes exactly as before. An Issue that
  is not found produces a named selection diagnostic. An Issue selected
  alone is Next, in the line format above; alongside an Objective/Task it is
  shown as context.
- When the router selects a Task or Objective that is already `done`, a new
  `SelectionDiagnostic` kind (for example `SelectionDone`) names the record.
  The board's Next area, non-TTY output, `savepoint resume`, and doctor all
  state it; Next still shows the selection as-is and never substitutes
  other work.
- Closing the selected Task on the board moves the router to the
  Objective's lowest-ID unfinished Task, or clears `task` when all are done;
  closing the selected Objective clears `objective` and `task`. The router's
  `release:` is unchanged byte-for-byte, proven by a test for each closure.
  The board's `p` key also preserves `release:`.
- The router carries no `next_action`. An existing router that still has one
  keeps loading; the value is ignored and doctor reports it once as a retired
  field, never rendered as an instruction. The live router, the V2 router
  template, and `WriteRouterStateV2` stop writing it.
- AGENTS.md tells an agent that a pasted Next line names the Objective,
  Task, and skill (Task `Build`/`Test` → `savepoint-task`, `Check` →
  `savepoint-check`, `Planned` Objective → `savepoint-design`).
- AGENTS.md, the router template, and the phase skills (live and scaffold,
  byte-identical) name `savepoint resume` as the one CLI command agents may
  run, replace "act on the router's next_action" with "act on
  `savepoint resume`'s Next", say what to do when the binary is unavailable
  (read the router selection and report the missing tool; do not guess), and
  name who advances the router after an owner closure.
- Apart from the rules above, gate decisions for the selected records, Goal
  completion rungs, and non-TTY parity are unchanged.
- Design.md Section 8 states the Next line format, that Next is the router
  selection, the board's closure advance, the stale-selection line, and the single-Next
  contract. Section 1 and Section 4 drop router `next_action` and the "press
  `p`" handoff wording where it no longer holds.
- Focused board/data/resume/doctor tests and `git diff --check` pass during
  iteration; `make build && make test-fast` passes at each Task handoff; the
  Full Objective Check has current `make test-full` evidence.

## Architectural Considerations

- `internal/data` stays the sole owner of Objective completion, Issue
  lifecycle, selection, and the ladder. The Objective word is a presentation
  translation in `internal/board/v2` (`next_panel.go`/`badges.go`) from typed
  `data` values; if it needs a typed source (for example a resolved Objective
  phase on `Next`), that value is computed in `internal/data`, not re-derived
  in the board.
- `internal/data/next.go` stays the one source of Next, but it reads the
  router selection instead of searching: the global and Release ready
  searches and the "active Task elsewhere in the Release" rung stop being
  sources of Next. This is a deletion, not a new ladder. O-020's ranking may
  later suggest what to select; it does not select.
- The stale-selection diagnostic is a `SelectionDiagnostic` kind so every
  surface reports it from one value. `internal/resume` owns its wording; the
  board's Next area gains one diagnostic line, a deliberate change to Section
  8's one-line contract.
- Task planning must confirm each hand-off `next_action` used to carry (for
  example "request a Task Check or record a waiver") is already expressed by a
  `data.Next` rung or resume phrase, and add a phrase only where it is not. No
  free text returns.
- Allowing agents to run `savepoint resume` narrows the "never run savepoint
  commands" rule to one read-only command that performs no writes, so owner
  authority over status changes is unaffected.
- Router YAML stays backward compatible: `issue:` is optional and
  `next_action` is tolerated on read.

## Boundaries

**In scope:**

- The Objective lifecycle word in the Next area and non-TTY output.
- Router `issue:` selection: decode, validate, write, `none` sentinel, a
  not-found diagnostic, an Issue-only Next line, and context display.
- Next derived from the router selection in `internal/data/next.go`,
  replacing the record searches, with matching resume phrasing and the
  shared Next line.
- The stale-selection diagnostic in data, rendered by board, non-TTY, resume,
  and doctor.
- Board closure advancing the router to the Objective's next Task (or
  clearing it) while preserving `release:` byte-for-byte; the `p` key fix.
- Retiring `next_action` from the router model, writer, template, and this
  repository's router; doctor's retired-field report.
- AGENTS.md, scaffold, and phase-skill guidance for `savepoint resume` and
  router advancement; Design.md reconciliation.

**Out of scope:**

- Mandatory Goal context for routers and Objectives (O-022).
- Objective priority and ordering (O-020).
- Changing Issue lifecycle, severity, or type vocabulary; board actions
  that advance Issues (O-015).
- Changing the Task words (`Build`/`Test`/`Check`/`Planned`/`Done`).
- Any Task, Objective, Goal, or Issue gate or completion-policy change: the
  precedence change decides what Next shows, not what may start or complete.
