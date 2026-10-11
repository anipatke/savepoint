---
id: C-981
scope: {kind: objective, id: O-048}
result: CLEAR
checked_by: {role: checker, session: check-o048-recheck-20261011}
executed_session: separate-executor-sessions-T-137-T-138-and-I-155-repair
checked_at: '2026-10-11T05:55:00Z'
health_snapshot: sha256:fb70b59f6409b104bc0c29fb985ead4e561e6ae63168108a4159f89a3deca634
reviewed:
  base_commit: 12e63a766ac24215358b9ed3eaae659a69bcd573
  head_commit: 6145746682c05d4f1be7c3cec88e7e0757ea47b7
  files:
    - .savepoint/objectives/O-048-animate-tiles/Objective.md
    - .savepoint/issues/I-156-unplanned-savepoint-task-skill-change-in-o048.md
    - internal/board/v2/update.go
    - internal/board/v2/spinner_test.go
  dependencies: []
issues: []
supersedes: C-980
---

# C-981: O-048 Full Objective Re-check

Re-check of C-980 against its frozen scope lock. The work C-980 reviewed was committed as 778f6e5; the remediation since then is 29e42c5 (I-156 owner acceptance and Objective boundary) and 6145746 (I-155 fix). This checker did not write either remediation commit. C-980 is unchanged and superseded by this immutable run.

**Result: CLEAR** (technical clearance). T-137's owner validation (its User Check) is still to be recorded by the owner.

## Closure Map

| Prior Issue | Result | Proof |
|---|---|---|
| I-155 spinner freezes after scope change | Verified | `Model.Update` now wraps `update` and calls `spinnerCmd()` after any non-tick message that returns no command; `applyLoad` returns nil so loads take the same path. Independent probes (temporary, deleted): esc back to the whole Goal restarts the timer; sidebar down/up walk through a quiet Objective stops then restarts it; Goal picker R-001→R-002→R-001(+tick)→R-002 starts, stops, restarts; an ordinary key on a running board schedules no second tick; keys on an idle board return no command. Regression test `TestSpinnerRestartsWhenAnObjectiveSwitchShowsALiveTileAgain` added. |
| I-156 unplanned skill change | Resolved by owner (`accepted`, 2026-10-11T06:30:00Z) | Owner decision recorded on the Issue; O-048 Boundaries now name the savepoint-task stage-writing step as in scope. Not a technical verification. |
| I-157 native Windows CI missing | Verified | CI run 38115852437, job 114400474006 `windows-tests`: success on head commit 6145746. |

## Admission Ledger

| Re-check item | Prior Issue / claim | Frozen cell | Allowed result |
|---|---|---|---|
| Esc back to live view | I-155 | Timer: scope change (Objective select / esc) | Passed |
| Sidebar up/down walk | I-155 | Timer: scope change (Objective select) | Passed |
| Goal picker switch | I-155 | Timer: scope change (Goal select) | Passed |
| No double timer | I-155 remediation | Timer: open/reload with live (one timer only) | Passed |
| Idle board stays idle | I-155 remediation | Idle board / Timer: reload with no live | Passed |
| Remaining original cells | C-980 matrix | all rows | Passed (gates and unchanged code) |
| Windows CI | I-157 | Gates: native Windows CI | Passed |
| Boundary | I-156 | Boundaries: workflow skills | Owner-accepted, boundary amended |

No probe ran outside the frozen lock.

## Full Check Progress Checklist

- [x] Establish Scope
- [x] Freeze The Check Scope (C-980 lock, unchanged)
- [x] Turn Acceptance Into Invariants
- [x] Build The Mandatory Coverage Matrix (C-980 matrix re-run)
- [x] Finite External-Boundary Matrix (not applicable, as in C-980)
- [x] Workflow And Side-Effect Check Lock
- [x] Matrix Completion Lock
- [x] Parallel Planning Advice (not a feature here)
- [x] Perform The Adversarial Pass
- [x] Re-check After Remediation
- [x] Verify File Reality
- [x] Verify Evidence And Gates
- [x] Collect Code Health Evidence
- [x] Complete The Issues Pass
- [x] Summarize Materiality
- [x] Review Code Style

## Gates And Evidence

| Command | Time (UTC) | Toolchain | Result |
|---|---|---|---|
| `git diff --check 12e63a7 HEAD` | 2026-10-11T05:50Z | — | exit 0 |
| `make build` | 2026-10-11T05:50Z | go1.26.2 linux/amd64 | exit 0 |
| `make test-full` (fresh, at 6145746) | 2026-10-11T05:50Z | go1.26.2 linux/amd64 | exit 0; golden files unmodified |
| CI `ci` job 114400474078 (`make ci`) | 2026-10-11T05:46Z | GitHub Actions | success |
| CI `windows-tests` job 114400474006 | 2026-10-11T05:48Z | GitHub Actions Windows | success |
| Focused spinner/tile tests + probes | 2026-10-11T05:49Z | go1.26.2 | pass |

File reality: all files named in the remediation exist at 6145746.

## Code Health

`savepoint health check O-048 .` → snapshot `sha256:fb70b59f…deca634` (created). "Code Health does not block clearance." Supporting evidence only.

## Acceptance Classification

| Criterion | Class |
|---|---|
| O-048 SC1 / T-137 DW1 | Proven |
| O-048 SC2 | Proven |
| O-048 SC3 / T-137 DW2 | Proven |
| O-048 SC4 / T-137 DW3 | Proven |
| O-048 SC5 / T-138 DW1-DW5 | Proven |
| O-048 SC6 / T-138 DW6 | Proven (local full gate + native Windows CI on 6145746) |
| T-137 DW4 | Proven |
| T-137 DW5 | Proven (scope-change test added) |
| T-137 DW6 | Proven for gates; owner validation pending |
| O-048 Boundaries | Satisfied (boundary amended by owner decision, I-156) |

## Materiality

No open Issues; no materiality actions required.

## Owner Decisions

No structured `exception` or `owner_validation` acceptance exists on O-048 or its Tasks, so there is nothing to carry forward. The Task-check waivers on T-137 and T-138 are unchanged. I-156's owner acceptance is an Issue resolution and stands.

## Owner Validation Still Needed

T-137 `owner_validation.required: true`, `accepted_check: ""`. The owner's User Check (spinner at build and test, still at check, idle board unchanged) can now be given against C-981.

## Observations (non-blocking)

- Design §8 "Task tile run animation" still reads "planned, not built" and says O-045's run code still ships, and does not mention the savepoint-task stage-writing step accepted under I-156. The planner should reconcile it.
- The Update wrapper calls `hasLiveTile()` (a scan of the shown cards) after every non-tick message while no timer runs; trivial cost at board sizes, no IO.
- Narrow single-column layout may tick while the live tile's column is off screen (T-137 Limitations).

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [x] STYLE-03 **Test branches**
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [x] STYLE-07 **One source of truth**
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data**
- [x] STYLE-10 **Small diffs**
