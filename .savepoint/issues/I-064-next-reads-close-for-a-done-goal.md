---
id: I-064
title: Next reads "Close" for a Goal that is already done
type: defect
status: resolved
source:
  kind: report
  actor: {role: owner, session: user}
  at: '2026-09-25T22:38:48Z'
severity: low
history:
  - at: '2026-09-25T22:38:48Z'
    actor: {role: owner, session: user}
    kind: observed
    note: >-
      After closing R-006 from the Goal selector (I-063), savepoint resume and
      the board still lead with "Close R-006" while the guidance line below
      says the Goal is already done.
  - at: '2026-09-25T22:39:38Z'
    actor: {role: executor, session: i064-repair-20260926}
    kind: repair_attempted
    note: >-
      resume.NextVerb now returns Done for NextReleaseReady when the selected
      Goal is status done, matching objectiveVerb; a ready Goal that is not
      done still reads Close. The board Next area and non-TTY output use the
      same verb. resume_test adds a "goal already done" case;
      TestNextReadsDoneForAClosedGoal closes a real Goal and asserts Close
      before and Done after in the board view and renderPlain, and fails
      with the fix reverted. gofmt, git diff --check, and make build && make
      test-fast passed. Live repo resume now reads "Done R-006". Issue remains
      open for independent verification.
  - at: "2026-09-25T22:52:38Z"
    actor:
      role: owner
      session: board-owner
    kind: owner_decision
    note: Moved from open to in_progress by the owner from the board.
  - at: '2026-09-25T23:03:15Z'
    actor: {role: owner, session: user}
    kind: owner_decision
    note: >-
      Owner instructed the executor to mark this Issue resolved once its fix
      was deployed. The fix shipped in savepoint 2.0.1 on npm (tag v2.0.1).
      No technical CLEAR is implied.
resolution:
  disposition: accepted
  actor: {role: owner, session: user}
  at: '2026-09-25T23:03:15Z'
  reason: Owner accepted the fix as deployed in savepoint 2.0.1.
---

# I-064: Next reads "Close" for a Goal that is already done

## Summary

With a Goal selected and no Objective, `data.ResolveNext` returns
`NextReleaseReady` whenever `ResolveReleaseCompletion` allows the Goal,
which stays true after the Goal is recorded `status: done`.
`resume.NextVerb` maps `NextReleaseReady` to `Close` without reading the
Goal's status, so the Next line says `Close R-###` for a finished Goal.
`releaseReadyNextActionPhrase` already checks `status: done` and says no
further action is required, so the headline and guidance disagree.
`objectiveVerb` already returns `Done` for a done Objective; Goals lack the
same check.

## Evidence

- `internal/data/next.go:476` `resolveReleaseCompletionRung` sets
  `NextReleaseReady` from the completion decision alone.
- `internal/resume/resume.go:86-90` `NextVerb` returns `Close` for
  `NextReleaseReady`.
- `internal/resume/resume.go:480` guidance returns "Goal %s is already done;
  no further Goal action is required."
- Live repo after R-006 closed: `savepoint resume` first line
  `Close R-006 — Savepoint V2 — Simple workflow, trustworthy completion`.

## Proof Needed

- A selected Goal with `status: done` reads `Done R-### — Title` in
  `savepoint resume`, the board Next area, and non-TTY output.
- A selected ready Goal that is not done still reads `Close R-### — Title`.
- Tests cover both cases; `make build && make test-fast` pass.
