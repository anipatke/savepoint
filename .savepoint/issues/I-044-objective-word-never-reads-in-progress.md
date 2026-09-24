---
id: I-044
title: The Next line's Objective word never reads In Progress
type: defect
status: open
source:
  kind: check
  check: C-916
  actor: {role: checker, session: o014-objective-check-20260924}
  at: '2026-09-24T08:58:19Z'
tasks: [T-028]
checks: [C-916]
severity: medium
history:
  - at: '2026-09-24T08:58:19Z'
    actor: {role: checker, session: o014-objective-check-20260924}
    kind: observed
    check: C-916
    note: >-
      The owner reported that the live line "Planned O-014 · Check" reads in
      the wrong order. It was confirmed during the O-014 Full Objective Check.
---

# I-044: The Next line's Objective word never reads In Progress

## Summary

`resume.ObjectiveWord` translates the Objective's recorded `status`. Nothing
ever writes an Objective's `status: in_progress`: no board action, no
`internal/data` writer, and no skill step does it. So an Objective reads
`Planned` from its first Task's Build through its integration Check, then
jumps to `Done`. Every Objective in this repository is `planned` or `done`
(7 and 6); none is `in_progress`.

## Violated requirement

O-014 Outcome and Success Conditions: the Next line states the Objective's
lifecycle word, and the Outcome's example is
`In Progress O-014 · Build T-028`. In practice that example can't happen, and
the owner sees `Planned O-014 · Build T-028` and `Planned O-014 · Check`.
The Why says the word must say whether the Objective is still building or
ready for Check.

## Reproduction

1. Select an Objective whose status is `planned` and a Task at `stage: build`.
   Next reads `Planned O-001 · Build T-001 — …`.
2. Mark every owned Task `done`. Next reads `Planned O-001 · Check — …`.
3. Live repository, 2026-09-24: `./savepoint resume` prints
   `Planned O-014 · Check — Give the Next area an Objective word and let the
   router target Issues`.

Expected: `In Progress` once any owned Task has started. Actual: `Planned`.

## Evidence

- `internal/resume/resume.go:99` `ObjectiveWord` maps `objective.Status` only.
- No `ColumnInProgress` write for an Objective exists in `internal/`.
  `agent-skills/savepoint-task/SKILL.md:33` starts the Task but never the
  Objective.
- C-916 permutation harness: every `obj planned` row reads `Planned`, whatever
  stage the Task is at.
- The existing tests build fixtures with the Objective already set to
  `in_progress`, so no test covers the lifecycle actually shipped.

## Repair

Owner decision, 2026-09-24: record the status (option 2), not a derived
word, so every surface agrees from the file. It is repaired directly under
this Issue (issue-capture.md, Out-Of-Scope Repair), with no new Task.

Options considered: Either derive the word in `internal/data`: an
Objective reads `In Progress` when any owned Task is `in_progress` or `done`
and the Objective is not `done` (O-014's Architectural Considerations already
allows a typed value on `Next`). Or have the board and `savepoint-task` write
`status: in_progress` on the Objective when its first Task starts. Add a test
that starts from a `planned` Objective.

## Repair Scope

Context: `internal/board/v2/io.go`, `internal/data/objective_gate_v2.go`,
`internal/doctor/checks.go`, `agent-skills/savepoint-task/SKILL.md` (and its
scaffold copy), `AGENTS.md` (and its scaffold copy), `.savepoint/Design.md`,
and O-014's `Objective.md`. Guardrails: DATA-01, DATA-02, DATA-05, FS-04,
TPL-01, TPL-02, TEST-01..03, TEST-05, TEST-08. Gate: `make build && make
test-fast`; the O-014 re-check (superseding C-916) verifies it.

- Board Space on a `planned` Task (planned → in_progress/build) also sets the
  owning Objective to `in_progress` when it is `planned`, in the same action,
  through `data.WriteObjectiveV2`. The Task write happens first. If the
  Objective write fails, the Task stays started and the action reports the
  error; the doctor warning below then names it.
- An Objective already `in_progress` or `done` is never rewritten. Retreating
  a Task never moves the Objective back. Only the owner sets `done`, as today.
- `savepoint-task` step 2 (live and scaffold, byte-identical) also sets the
  owning Objective's `status: in_progress` when it is `planned`. AGENTS.md's
  lifecycle rules (live and scaffold) allow agents to do exactly that, and
  nothing more, for Objectives.
- `data.InspectObjectiveConsistency` gains a kind for a `planned` Objective
  that owns an `in_progress` or `done` Task. Doctor reports it as a `!`
  pending-review warning (not `✗`) with a repair that says to set the
  Objective's `status: in_progress`.
- This repository's O-014 is set to `in_progress`, so the new warning does not
  fire here.
- No gate result changes: no resolver reads `planned` vs `in_progress` for an
  Objective (checked 2026-09-24 in `gate_v2.go` and `objective_gate_v2.go`).
- Design sections 4 and 8 say who moves an Objective to `in_progress`.
- Tests: board start from a `planned` Objective writes both records (the
  case every existing fixture skips); an already `in_progress` Objective is
  untouched; an Objective write failure keeps the Task started and reports
  the error; retreat leaves the Objective alone; the doctor warning fires and
  stays silent for a `planned` Objective with only planned Tasks; the resume
  and board Next line reads `In Progress O-### · Build T-###` after a start.
- `make build && make test-fast` pass.
