---
id: I-041
title: Design.md says the board loads through data.LoadProject
type: drift
status: resolved
source:
  kind: check
  check: C-914
  actor: {role: checker, session: o021-objective-check-20260924}
  at: '2026-09-23T22:16:20Z'
tasks: [T-026, T-027]
checks: [C-914, C-915]
guardrail_ids: [TPL-02]
severity: low
resolution:
  disposition: verified
  check: C-915
  actor: {role: checker, session: o021-objective-recheck-20260924}
  at: '2026-09-23T23:04:15Z'
history:
  - at: '2026-09-23T22:16:20Z'
    actor: {role: checker, session: o021-objective-check-20260924}
    kind: observed
    check: C-914
    note: Found while reconciling T-026's Design edits against the live board load path.
  - at: '2026-09-23T22:59:28Z'
    actor: {role: executor, session: o021-remediation-20260924}
    kind: repair_attempted
    note: >-
      T-027 updated Design to describe CheckRuntimeSchema, LoadV2Index,
      router decoding, and ResolveNext. Fresh make build and make test-full
      passed. Awaiting independent Full Objective Check C-914 recheck.
  - at: '2026-09-23T23:04:15Z'
    actor: {role: checker, session: o021-objective-recheck-20260924}
    kind: rechecked
    check: C-915
    note: Independent re-check reproduced the Proof Needed; see C-915.
---

# I-041: Design.md says the board loads through data.LoadProject

## Summary

T-026 edited `.savepoint/Design.md`'s "Board persistence and refresh"
paragraph (line 182) to drop migration-state detection. The edited paragraph
still says every board startup and reload uses one load command built from
"`data.LoadProject`, V2 state and router decoding, and `data.ResolveNext`".

The board does not call `data.LoadProject`. `internal/board/board.go:53` runs
`data.CheckRuntimeSchema`, and `internal/board/v2/load.go:75` calls
`data.LoadV2Index`. So Design still names code that no command reaches (see
I-040). That misses O-021 Success Condition 9 ("`Design.md` describe[s] the
smaller shape") and TPL-02.

## Evidence

- `.savepoint/Design.md:182` (working tree, sha256 `426153dd…`).
- `grep -rn 'CheckRuntimeSchema\|LoadV2Index' internal/board main.go`.

## Proof Needed

Line 182 names the real load path (`CheckRuntimeSchema` and `LoadV2Index`,
then router decoding and `ResolveNext`), and no guidance names `LoadProject`
as a runtime path. The natural fix is to do this in the same change as I-040.
