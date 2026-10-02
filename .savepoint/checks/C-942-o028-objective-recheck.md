---
id: C-942
scope: {kind: objective, id: O-028}
result: NEEDS WORK
checked_by: {role: checker, session: o028-recheck-20261001}
executed_session: o028-repair-20261001
checked_at: '2026-10-01T08:53:46Z'
reviewed:
  base_commit: 919544c0da6e53f73e2b61b01fcb7e24098b72fa
  head_commit: 01105a341ce76c396dc7e04d34f10ecbace47e0c
  files:
    - AGENTS.md
    - internal/codehealth/collect.go
    - internal/codehealth/collect_repair_test.go
    - internal/codehealth/config.go
    - internal/codehealth/runner_alive_windows_test.go
    - internal/codehealth/runner_test.go
    - internal/codehealth/runner_windows.go
    - internal/codehealth/snapshot.go
  dependencies: []
issues: [I-093]
supersedes: C-941
---

# C-942: O-028 Full Objective Re-check

## Result and authority

**NEEDS WORK.** Four of five C-941 Issues are proven repaired. I-093 is still
open: the rename repair keeps history across a rename, but it introduces a
cross-instance leak in the same frozen M2 cell ("two independent instances")
that I-093's own Proof Needed names.

This session is independent of the executors and of the repair session. It
did not build O-028 or write the I-092..I-096 repairs. It also did not write
C-941. That record and repair commit 01105a3 carry the same Claude session
link, so C-941's own independence from the repair is not something this
record relies on. This re-check reproduced every C-941 Issue from scratch.
No router, Task or Objective status was changed. No implementation was
repaired.

This is the one full re-check under the convergence limit. One targeted
remediation and re-check of I-093 remains before the owner decides.

## Closure map

| Issue | State | Basis |
| --- | --- | --- |
| I-092 | Closed by evidence | Frozen M5 C-941 harness passes for Lizard, jscpd, OSV. The Discover→Plan→Apply→Collect probe below covers ok/silent/nonzero/missing/malformed/cancel/stale for all three; the report is never left in the project. |
| I-093 | **Still open** | Rename continuity passes (unnamed→named, named→named; changed scope correctly starts over). Two scoped instances now each report the other's results as "N earlier official results not compared"; this did not happen at 919544c. |
| I-094 | Closed by evidence | Validation and Collect agree at the capacity boundary 27/28 accepted, 29/32 refused, with zero tool calls on refusal. C-941's original 32-instance probe now fails validation before any run, which is the remedy I-094 allowed. |
| I-095 | Closed by evidence | Native Windows/amd64: C-941's cancel probe and a new deadline probe both PASS (child wait=0). The repository's un-skipped children test, ExecRunner and every TestCollect pass natively. |
| I-096 | Closed by evidence | AGENTS.md map describes Collect, and says production readers and a collection command are not shipped. No template carries the old sentence. |

"Closed by evidence" is not a `verified` resolution. Under issue-capture.md
that needs a CLEAR Check, so these four Issues get `rechecked` history entries
and stay open until a CLEAR re-check.

## Admission ledger

| Probe | Prior Issue / claim | Frozen cell | Allowed result |
| --- | --- | --- | --- |
| C-941 harness M5 + recheck M5 delivery/cleanup matrix | I-092 | M5 (delivery, missing, malformed, cleanup) and adjacent M8 nonzero/cancel | blocking |
| C-941 harness M2 + rename variants | I-093 | M2 rename/history integration | blocking |
| Two scoped instances collected together | I-093 Proof Needed "two independent instances" | M2 independent instances | blocking |
| C-941 harness M10 + capacity boundaries | I-094 | M10 max capacity plus placeholders | blocking |
| Native Windows cancel, deadline, children, collect | I-095 | M7 Windows, M8 cancel/temp cleanup | blocking |
| Map text | I-096 | M12 docs | blocking |
| Report path prefix vs stale project report | none | none | observation only |

## I-093 remaining failure

Scenario: config with `lizardInstance("api", …, "api/**")` and
`lizardInstance("web", …, "web/**")`; four official Collects; api reads 5, web
reads 40.

- Expected (and actual at 919544c): `api: … Unchanged at 5 over 4 official checks.`
- Actual at 01105a3: `api: … Unchanged at 5 over 4 official checks. 3 earlier official results not compared.` and the same for web.

Cause: `collect.go:186` now groups history by `historyKey()` (capability +
provider), so each instance's history includes its siblings. `selectSeries`
(`history.go:66-68`) counts every measured official result in a different
series as `incompatible`, and `describe` (`history.go:157`) writes it into the
saved, immutable summary. Values and trend stay separate. Only the persisted
explanation is wrong, and the count grows with every instance and every run.

`TestCollectKeepsTwoScopedInstancesSeparate` checks only for "Not enough
comparable", so it passes with the leak. Violates T-055 Done When 5 (collected
history separation) and O-028 Success Condition 8 (separate scoped instance
values), within frozen M2.

Materiality: Likelihood High (any project with two instances of one
capability, the normal monorepo setup). Impact Low–Medium (no wrong
classification, but every saved summary misstates its history permanently).
Materiality Medium. Recommendation: fix now as the one targeted remediation,
for example by counting only prior results of this instance (same name or
same series) as `incompatible`. Strengthen the test to reject sibling counts.

## Evidence and gates

- Fresh `make test-full` at head 01105a3, go1.26.2 linux/amd64, 2026-10-01: exit 0, including linux/darwin/windows builds.
- `git diff --check 919544c 01105a3`: clean. `go vet ./internal/codehealth ./cmd`: clean. go.mod/go.sum/Makefile unchanged.
- C-941 Linux harness restored by Go overlay (`go test -overlay … -run '^TestO028(IndependentMatrix|Recheck)$' -count=1 -v`): M5 ×3 PASS, M2 rename PASS. M10 "fails" only because the 32-instance config is now refused by validation before any tool runs, which is the intended remedy.
- Re-check harness `TestO028Recheck` (overlay, scratch only): M5 21/21 cells PASS. Lizard and jscpd nonzero = failed; OSV exit 1 = available, per its findings code. Stale pre-existing report is removed before the run and reported "wrote no report". M10 PASS. M2 rename ×3 PASS. Sibling cell FAIL as above. The same sibling probe at 919544c (temporary worktree): PASS.
- Native Windows: `GOOS=windows GOARCH=amd64 go test -overlay … -c` run from PowerShell with a Windows-local TEMP. TestO028WindowsChildren, TestO028WindowsDeadlineChildren, TestExecRunner*, all TestCollect* PASS (two skip by design). A first run with TEMP on the WSL UNC path failed git-based tests on git's "dubious ownership" check. That is environmental, not product behaviour.
- File reality: every file in the repair commit exists. Scratch harnesses lived only in the session scratchpad and were never added to the repository.

## Observations (non-blocking)

- `discardOwnReport` clears and recreates the report only under `.savepoint/health/reports/`. A hand-configured executed tool whose `report` points elsewhere in the project is not cleared first, so a leftover file from an earlier run could be read as this run's output. Discovery never proposes such a path, and the C-941 lock has no cell for hand-configured executed report locations, so this is outside the blocking perimeter. Worth noting for O-029/O-030.
- Windows `Cancel` ignores a `taskkill` failure and falls back to `Process.Kill`, which is appropriate. `taskkill` resolves from `%SystemRoot%`, not PATH.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [ ] STYLE-03 **Test branches** — `TestCollectKeepsTwoScopedInstancesSeparate` asserts only the absence of one phrase, so it misses the sibling leak.
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [x] STYLE-07 **One source of truth** — capacity now derives from `Capabilities()` in one place.
- [ ] STYLE-08 **Comments explain why** — `history()` comment claims "Assess then narrows to the comparison series, which separates scopes", but Assess still counts siblings in its explanation.
- [x] STYLE-09 **Content lives in data**
- [x] STYLE-10 **Small diffs**

## Key re-check probe

```go
// Overlay into internal/codehealth; scratch only.
run := func(t *testing.T, root string, cfg Config, n int, val func(ReportInput) float64) Snapshot {
	tools := &fakeTools{t: t, behavior: map[string]func(context.Context, ToolSpec) (ToolResult, error){"fake": stdout("report")}}
	readers := Readers{ProviderLizardCSV: readerFunc(func(_ context.Context, in ReportInput) (Reading, error) { return goodReading(val(in), UnitCCN), nil })}
	got, err := Collect(context.Background(), CollectRequest{Root: root, Origin: OriginOfficial, Config: cfg, Readers: readers, Runner: tools,
		Clock: func() time.Time { return testClock().Add(time.Duration(n) * time.Second) }})
	if err != nil { t.Fatal(err) }
	snaps, _ := NewStore(root).LoadSnapshots()
	for _, s := range snaps { if s.ID == got.SnapshotID { return s } }
	return Snapshot{}
}
root := project(t)
cfg := cfgOf(lizardInstance("api", "fake", "api/**"), lizardInstance("web", "fake", "web/**"))
val := func(in ReportInput) float64 { if in.Scope[0] == "web/**" { return 40 }; return 5 }
var s Snapshot
for i := 0; i < 4; i++ { s = run(t, root, cfg, i, val) }
for _, sm := range s.Summary.Capabilities {
	if sm.Capability == CapabilityComplexity && strings.Contains(sm.Explanation, "not compared") {
		t.Errorf("%s reports sibling instance results as its own earlier results", sm.Name)
	}
}
```
