---
id: I-060
title: Remove the Goal Check; only the Objective Check and optional Task Check remain
type: drift
status: open
source:
  kind: report
  actor: {role: owner, session: owner-chat-20260925}
  at: '2026-09-25T11:12:00Z'
severity: medium
history:
  - at: '2026-09-25T11:12:00Z'
    actor: {role: owner, session: owner-chat-20260925}
    kind: owner_decision
    note: >-
      Owner: "No such thing as goal check. Mandatory Obj check and optional
      task check only." The Goal Check is too expensive and is descoped.
---

# I-060: Remove the Goal Check

## Summary

Design, AGENTS.md, the skills, and the runtime make a Full Goal Check
mandatory before a Goal can close. The owner has decided there is no Goal
Check: the only Checks are the mandatory Full Objective Check and the
optional Task Check. A Goal is complete when all its member Objectives are
complete.

## Evidence

- Live `savepoint doctor` reports `[v2-release-clearance-missing]` for R-006
  although all 13 member Objectives are done.
- `savepoint resume` Next: `Record a fresh Goal Check for R-006.`
- `internal/data/release_gate_v2.go` `resolveReleaseCompletionForRecord`
  requires current `scope.kind: release` clearance after member Objectives.
- The rule appears in about 40 files: `internal/data` (gate, next),
  `internal/doctor`, `internal/resume`, `internal/board/v2` detail view,
  `.savepoint/Design.md`, `.savepoint/Guardrails.md`, `.savepoint/router.md`,
  `AGENTS.md`, `README.md`, every skill and shared reference (live and
  `templates/project-v2`), and migration goldens.

## Proof Needed

Goal completion depends only on member Objective completion. Resume, board,
and doctor no longer ask for a Goal Check. Active guidance and the scaffold no
longer describe one. Existing `scope.kind: release` Check records still load.
Live doctor reports no Goal Check error for R-006.
