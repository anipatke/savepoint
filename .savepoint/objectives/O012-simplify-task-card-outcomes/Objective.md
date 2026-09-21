---
id: O012
title: Make Task cards easier to scan
status: in_progress
release: R006
---

# O012: Make Task cards easier to scan

## Outcome

Task cards communicate implementation stage, one understandable review outcome,
and only the blockers that still require action. Completion is conveyed by the
Done column rather than repeated by another badge, while owner waiver and owner
risk acceptance remain visibly sufficient without being misrepresented as
independent technical clearance.

## Why

The current card combines stage, clearance freshness, completion provenance,
and gate blockers into a dense row that repeats facts and exposes distinctions
better suited to detail and Next surfaces. The owner has selected a smaller
vocabulary that preserves accountability while making the board easier to scan.

## Success Conditions

- Planned, In Progress, and Done remain the only Task columns; Build, Test, and
  Check remain the only displayed implementation stages.
- Each non-planned card shows at most one review outcome: `[ ] CHECK`,
  `[✓] CHECK`, `[!] NEEDS WORK`, `[!] REVIEW`, `[✓] WAIVED`, or
  `[✓] OWNER ACCEPTED`.
- Waiver and owner acceptance use the green success treatment but remain
  distinct from independent Check clearance in data, detail, Next, dependency,
  Objective, and Release behavior.
- Stale and unknown/unverified clearance collapse to `[!] REVIEW` on Task cards
  while their exact evidence remains available outside the card.
- Ordinary `DONE`, warning `DONE`, `BY WAIVER`, and `BY EXCEPTION` completion
  badges no longer appear on Task cards; the Done column and review outcome
  carry those facts without duplication.
- Replan, Task dependency, Objective dependency, and pending owner-action
  blockers remain visible when actionable; checker-authority detail is folded
  into `[!] REVIEW`, and an owner-accepted override does not also display the
  blocker it resolved.
- Interactive and non-TTY cards use the same labels and ordering, remain
  understandable without colour, and retain narrow-terminal wrapping behavior.
- The disposable O900 test-data Objective visibly exercises every retained
  outcome plus wait, replan, and owner blockers so the simplified treatment
  can be reviewed before release.
- Active architecture and repository guidance record the simplified vocabulary
  and explicitly preserve the underlying verification policy so a later Check
  cannot infer the retired presentation from older evidence.

## Architectural Considerations

- `internal/data` remains the sole owner of clearance, waiver, exception,
  dependency, and completion decisions. This Objective changes presentation,
  not gate semantics.
- `internal/board/v2/badges.go` remains the single typed-value-to-badge mapping;
  card rendering remains pure and performs no IO or gate re-evaluation.
- Existing palette roles remain authoritative: green means an outcome accepted
  for the Task, orange means attention, purple means waiting, and dim means no
  review evidence yet. Glyph and text continue to carry meaning without colour.
- Detail, Next, resume, doctor, and mandatory Objective/Release Checks retain
  their exact evidence wording and policy distinctions.
- Planned change stays in this Objective until the mandatory Full Objective
  Check reconciles the implemented result into `.savepoint/Design.md`.

## Boundaries

**In scope:**

- V2 Task-card badge composition, labels, accents, blocker suppression, plain
  output parity, focused rendering tests, and active documentation review.
- Documentation reconciliation in current architecture, repository guidance,
  visual policy, README, and the active V2 scaffold guide where applicable.

**Out of scope:**

- Changes to clearance states, Check freshness, waiver/exception records,
  completion resolvers, dependencies, Objective or Release gates, detail/Next/
  resume evidence wording, Objective sidebar badges, or Issue presentation.
- Rewriting historical Checks, Issues, or archived V1 material; those remain
  evidence rather than this Objective's documentation authority. O900 is the
  disposable live UI fixture and may be refreshed only to visualise this
  Objective's card outcomes.
