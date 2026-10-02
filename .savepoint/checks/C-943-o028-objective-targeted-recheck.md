---
id: C-943
scope: {kind: objective, id: O-028}
result: CLEAR
checked_by: {role: checker, session: o028-recheck-20261001}
executed_session: o028-repair-20261001
checked_at: '2026-10-01T08:58:32Z'
reviewed:
  base_commit: 01105a341ce76c396dc7e04d34f10ecbace47e0c
  head_commit: 01105a341ce76c396dc7e04d34f10ecbace47e0c
  files:
    - internal/codehealth/collect.go
    - internal/codehealth/collect_repair_test.go
  dependencies: []
issues: []
supersedes: C-942
---

# C-943: O-028 Targeted Re-check (I-093)

## Result and authority

**CLEAR.** This is the targeted remediation re-check that the convergence
limit allows after C-942. It uses C-941's frozen scope lock unchanged. The
checker session is the one that wrote C-942. It is independent of the
executor and repair sessions and repaired nothing.

The reviewed change is the uncommitted working-tree repair on top of 01105a3
(`collect.go` `ownHistory`, and a strengthened
`TestCollectKeepsTwoScopedInstancesSeparate`). The `head_commit` names that
base; the diff was reviewed as it stands in the working tree. No router, Task
or Objective status was changed.

## Closure map

| Issue | State | Basis |
| --- | --- | --- |
| I-092 | Closed (verified) | C-942 evidence. The C-941 M5 harness and the C-942 M5 delivery matrix re-ran on the repaired tree: all pass. |
| I-093 | Closed (verified) | See below. |
| I-094 | Closed (verified) | C-942 evidence. The M10 boundary probe re-ran: 27/28 accepted, 29/32 refused with zero tool calls. |
| I-095 | Closed (verified) | C-942 native Windows evidence. `collect.go`'s change is platform-neutral history filtering and does not touch the runner. |
| I-096 | Closed (verified) | C-942 evidence. `AGENTS.md` unchanged since. |

## Admission ledger and I-093 evidence

| Probe | Prior Issue | Frozen cell | Result |
| --- | --- | --- | --- |
| Two scoped instances (api=5, web=40), four official Collects | I-093 / C-942 sibling leak | M2 independent instances | PASS: `api: … Unchanged at 5 over 4 official checks.` and `web: … Unchanged at 40 over 4 official checks.`, with no "not compared" text. Matches 919544c. |
| Unnamed→named, named→named rename | I-093 | M2 rename/history | PASS: "Unchanged at 5 over 4 official checks." |
| Same name, changed scope | I-093 | M2 changed scope | PASS: starts over, "3 earlier official results not compared." |
| C-941 harness M2 | I-093 | M2 | PASS |
| C-941 harness M5, C-942 M5 matrix (21 cells), M10 boundaries | I-092, I-094 | M5, M8, M10 | PASS. C-941's original M10 subtest still "fails" only because a 32-instance config is now refused before any run, which is I-094's intended remedy. |

`ownHistory` keeps earlier results that share the instance's name, or that
share its comparison series. It drops a sibling's results, which have neither,
so they no longer reach `selectSeries`'s `incompatible` count. Reverting only
`collect.go` makes the strengthened repository test fail, so it guards the
regression.

## Evidence and gates

- Fresh `make test-full` on the repaired working tree, go1.26.2 linux/amd64, 2026-10-01T08:58Z: exit 0, including linux/darwin/windows builds.
- `git diff --check`: clean. `go vet ./internal/codehealth`: clean. go.mod/go.sum/Makefile unchanged.
- `go test -overlay <scratch>/overlay.json ./internal/codehealth -run '^TestO028(IndependentMatrix|Recheck)$' -count=1 -v`: TestO028Recheck PASS in full. IndependentMatrix fails only on the superseded M10 expectation explained above.
- With `collect.go` stashed: `go test ./internal/codehealth -run '^TestCollectKeepsTwoScopedInstancesSeparate$'` FAIL. Restored: PASS.
- Scratch harnesses stayed in the session scratchpad and are not in the repository.

## Closure readiness

Every owned Task (T-055–T-058) is `done` with an owner Task-check waiver. With
this CLEAR, no material Issue linked to the current Check remains open. Two
owner items are still outstanding before O-028 can close:

1. T-058 declares owner validation. The owner's scratch-project walkthrough
   and acceptance naming C-943 are still outstanding.
2. The I-093 repair is uncommitted. It should be committed as reviewed; any
   further code change makes this Check stale.

This Check records no owner acceptance and sets no status.

## Observations (non-blocking)

- If two instances swap names while keeping their scopes, each also matches
  the other's old results by name. Its own series is still found, but the
  sibling's results add to the "not compared" count. This is an exotic edit
  and not a frozen M2 cell.
- The C-942 note about hand-configured executed reports outside
  `.savepoint/health/reports/` not being cleared before a run still stands,
  for O-029/O-030.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [x] STYLE-03 **Test branches** — the sibling test now rejects both lost history and sibling counts.
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [x] STYLE-07 **One source of truth**
- [x] STYLE-08 **Comments explain why** — the `history()` comment now points to `ownHistory`.
- [x] STYLE-09 **Content lives in data**
- [x] STYLE-10 **Small diffs**
