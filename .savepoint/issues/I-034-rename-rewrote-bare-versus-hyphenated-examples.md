---
id: I-034
title: Rename rewrote bare-versus-hyphenated examples into nonsense
type: defect
status: resolved
source:
  kind: check
  check: C-910
  actor: {role: checker, session: o018-full-check-20260923}
  at: '2026-09-23T08:40:00Z'
tasks: [T-020]
checks: [C-910]
severity: low
resolution:
  disposition: accepted
  actor: {role: owner, session: user}
  at: '2026-09-23T09:05:00Z'
  reason: >-
    Owner accepted the repair (three self-comparing sentences reworded, T-020 evidence corrected) without a fresh-session recheck. The repair was made by the same session that wrote C-910. This is owner acceptance, not a technical verified closure.
history:
  - at: '2026-09-23T08:40:00Z'
    actor: {role: checker, session: o018-full-check-20260923}
    kind: observed
    check: C-910
    note: >-
      Three records contrasted the old bare form with the new hyphenated form.
      The rename script hyphenated the bare side too, so each sentence now
      compares an ID with itself. T-020 evidence says the dry run had no false
      positives.
  - at: '2026-09-23T08:55:00Z'
    actor: {role: executor, session: o018-full-check-20260923}
    kind: repair_attempted
    note: >-
      At the owner's direction, rephrased the three sentences without a literal
      bare ID (O-018 Objective.md Why, the I-022 observed note, T-018 evidence
      line 68), so the targeted bare-ID search stays clean. Added a correction
      note to T-020's evidence. This repair was made by the same session that
      wrote C-910, so it needs a fresh-session recheck before closure.
---

# I-034: Rename rewrote bare-versus-hyphenated examples into nonsense

## Violated Requirement

T-020 Done When: the script "rewrites identity references". These were
examples of the old format, not references to records. T-020's plan step 2
says to "review the diff, and fix any false positives by hand". T-020's
Technical Evidence says "No false positives", which is inaccurate.

## Scenario

Before (HEAD `a1b5587`) → after (staged):

- `.savepoint/objectives/O-018-hyphenate-numbered-record-identities/Objective.md:19`
  — "bare numbered identities such as `O015` and `T012` are harder to scan
  than `O-015` and `T-012`" became "...such as `O-015` and `T-012` are harder
  to scan than `O-015` and `T-012`".
- `.savepoint/issues/I-022-hyphenate-numbered-record-identities.md:24` — the
  same sentence in the `observed` history note.
- `.savepoint/objectives/O-018-hyphenate-numbered-record-identities/tasks/T-018-map-the-identity-cutover.md:68`
  — "prefix insertion only (`R006` → `R-006`)" became "(`R-006` → `R-006`)".

Expected: each example still shows the bare form next to the hyphenated form.
Actual: each now compares an ID with itself, so the original reason for the
change can't be read from these records.

## Proof Needed

Each of the three sentences again shows the old and new forms side by side.
There are two ways to do that:

- Restore the bare examples. A record body is prose, not a parsed identity,
  so strict load still succeeds. The Objective's targeted search would then
  list these three lines as documented exceptions.
- Rephrase them without a literal bare ID, e.g. "a letter followed directly
  by digits".

T-020's "No false positives" claim is corrected in a later evidence note, not
by editing history.
