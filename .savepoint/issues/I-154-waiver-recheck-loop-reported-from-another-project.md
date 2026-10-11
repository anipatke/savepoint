---
id: I-154
title: Waiver-driven re-check loop reported from another project
type: other
status: open
source:
  kind: report
  actor: {role: executor, session: owner-report-20261011}
  at: '2026-10-11T04:43:10Z'
history:
  - at: '2026-10-11T04:43:10Z'
    actor: {role: executor, session: owner-report-20261011}
    kind: observed
    note: Owner pasted another project's analysis of a re-check loop (its O-009, I-018) and asked for it to be noted. Recorded as reported; none of it has been reproduced here.
---

# I-154: Waiver-driven re-check loop reported from another project

## Summary

The owner reports that another Savepoint project is still forced into a
re-check loop around a waiver. This overlaps I-141 (escalated to O-044) and
I-142 (verified by C-970), which already cover carry-forward of owner
decisions, "Close by exception" routing, and the `Assess` verb. Either that
project runs a release older than O-044, or part of the loop survives it.
That is the first thing to establish. Items 2 to 5 below were not covered by
I-141 or I-142 as far as their text shows.

## Evidence

Reported analysis (unverified here; the other project's own file paths and
rule quotes are its claims):

1. A waiver can never become CLEAR, so every later Check stays NEEDS WORK
   even when the fix is clean. (Covered in intent by I-141 / O-044.)
2. `resume` kept saying "Record the Objective integration Check" after the
   waiver, the confirmed fix and every sign-off were in place. (Same symptom
   as I-141 and I-142; confirm against the current release.)
3. An Issue whose own fix was proven could not be closed `verified`, because
   `verified` needs a Check that recorded CLEAR and an unrelated waiver kept
   the overall result NEEDS WORK. (New: Resolution Dispositions in
   issue-capture.md.)
4. check-method.md reads as self-contradicting: "a material unverified
   criterion means NEEDS WORK" against "Return CLEAR only when every criterion
   is proven". A low-importance, owner-acknowledged gap was rated unmet.
   (New.)
5. A Task was handed off "Partially met" for a requirement that was impossible
   from the start (the baseline code was never saved) instead of returning
   REPLAN REQUIRED, so the first Check met an unmeetable requirement. (New:
   savepoint-task hand-off and Closure Rules.)
6. Owner task sign-offs and an Issue acceptance were missing from the saved
   records, so the gate saw them as absent. Cause unknown. (New; may be a
   board save failure.)

Suggested changes in the report: resume routes to close-by-exception when
every unmet item has a waiver that still applies; an Issue may be `verified`
when a Check proves that Issue's own fix whatever the overall result; a minor
acknowledged unverified item becomes a pass with a note; an impossible
requirement goes back to planning or is waived at Task close; a waiver-only
or sign-off-only change uses `Assess`, not a new Check.

## Proof Needed

- Reproduce, or rule out, the loop on the current release with the shape in
  I-141 plus a waiver on a low-importance criterion; record whether resume
  reaches `Close`.
- For items 3 to 5, the planner decides whether each is a wanted change; the
  Savepoint skill files are overwritten on upgrade, so any change belongs in
  the shipped `agent-skills/` and templates.
- For item 6, find out whether the board or any command can drop an owner
  acceptance or Task sign-off without an error.
