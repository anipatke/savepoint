---
id: I-136
title: Sanitize authored lane heading controls
type: defect
status: open
source:
  kind: check
  check: C-965
  actor: {role: checker, session: check-o033-20261004}
  at: '2026-10-04T03:26:27Z'
tasks: [T-104, T-105]
checks: [C-965]
history:
  - at: '2026-10-04T03:26:27Z'
    actor: {role: checker, session: check-o033-20261004}
    kind: observed
    check: C-965
    note: Initial independent Full Objective Check.
---

# I-136: Sanitize authored lane heading controls

## Summary

Sanitize authored lane heading controls before O-033 clearance.

## Evidence

T-105 DW5 forbids project record strings injecting terminal controls. O-033 promises safe board grouping; frozen M5 covers authored control text on the new heading surface.

In an enabled temporary board, author a lane title in YAML as "Core\x1b[2J\x07". Load and render the interactive board. Expected: visible safe text with no clear-screen/bell sequence. Actual Model.View contains raw ESC[2J and BEL from the lane title. These can clear the terminal and ring its bell.

`internal/board/v2/column.go:210–219` passes the raw title through wrapTitleLines and styled rendering; `lanes.go:82` carries the authored title without sanitation. The shared advice formatter sanitizes its own text, but headings bypass it. Independent TestO033HeadingControls fails. Existing tests sanitize Task instruction titles, not lane-heading controls.

## Proof Needed

Strip authored terminal controls before rendering headings while preserving generated style escapes and Unicode text. Replay the clear-screen/BEL heading through interactive board, detail, plain and resume; use existing eight text-class/width matrix and heading/focus/count checks. Frozen M5/M6; no generic terminal overhaul.
