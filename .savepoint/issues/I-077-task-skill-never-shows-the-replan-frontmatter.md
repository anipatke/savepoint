---
id: I-077
title: Task skill never shows the replan frontmatter, so REPLAN REQUIRED does not route
type: defect
status: open
source:
  kind: report
  actor: {role: owner, session: user}
  at: '2026-09-26T05:18:49Z'
severity: medium
history:
  - at: '2026-09-26T05:18:49Z'
    actor: {role: owner, session: user}
    kind: observed
    note: >-
      Raised by an independent review of the packaged Savepoint skills; claim verified against the code in a follow-up review session before capture.
---

# I-077: Task skill never shows the replan frontmatter, so REPLAN REQUIRED does not route

## Summary

`savepoint-task` tells the executor to "set the replan reason with the
handoff evidence" but never shows the `replan:` frontmatter block the runtime
reads. `savepoint resume` shows `Replan` only when that block exists, so an
executor that writes its REPLAN REQUIRED in prose leaves the router pointing
at the Task as if nothing broke.

## Evidence

- `agent-skills/savepoint-task/SKILL.md:57` and `:71` describe the replan in
  words only; no skill or reference contains a `replan:` example.
- `internal/data/evidence_v2.go:117-121` requires `reason`,
  `recorded_by: {role, session}`, and `recorded_at`.
- T-010 Drift Notes (O-019): the executor returned REPLAN REQUIRED but
  "`savepoint resume` shows `Replan` only when the Task carries a `replan:`
  frontmatter block, and none was recorded."

## Proof Needed

- `savepoint-task` shows the exact `replan:` block with all required fields.
- A Task written from that example loads and `savepoint resume` reads
  `Replan`.
- Packaged template copies stay byte-identical.
