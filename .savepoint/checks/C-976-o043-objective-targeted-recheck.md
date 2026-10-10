---
id: C-976
scope: {kind: objective, id: O-043}
result: CLEAR
checked_by: {role: checker, session: check-o043-recheck2-20261010}
executed_session: executor-o043-i148-unrecorded
checked_at: '2026-10-10T08:10:00Z'
health_snapshot: sha256:aa05e1fd094901f1a42a2677d76cd254391ba180b6c804369ee8f851cb998a13
reviewed:
  base_commit: 0fb8c37
  head_commit: 3596b2c
  files:
    - templates/project-v2/.claude/hooks/guard.js
    - .claude/hooks/guard.js
    - internal/init/claude_guard_test.go
  dependencies: []
issues: []
supersedes: C-975
---

# C-976: O-043 targeted re-check (I-148, I-147)

CLEAR. This is the one targeted re-check C-975's Convergence section allowed, limited to the I-148 and I-147 cells of C-974's frozen lock. This checker session is fresh: it did not build O-043, the C-975 remediation, or the I-148 fix (committed as `3596b2c` in a separate session; working tree clean, `HEAD` equals `origin/2.3.0-skills-optimisation`).

Technical clearance only. O-043 still needs the owner's T-115 and T-117 User Check acceptances before it can close (see Owner Validation Still Needed).

## Full Check Progress Checklist

- [x] Establish Scope — diff `0fb8c37..3596b2c` outside `.savepoint/` touches only the two `guard.js` copies (+12/−1 each).
- [x] Freeze The Check Scope — C-974 lock reused unchanged; no new axis.
- [x] Turn Acceptance Into Invariants — T-118 router rule: block iff `--git-dir` ≠ `--git-common-dir`.
- [x] Build The Mandatory Coverage Matrix — frozen; re-run cells below.
- [x] Finite External-Boundary Matrix — `git rev-parse` subprocess cells from C-974 unchanged; fallback path covered below.
- [x] Workflow And Side-Effect Check Lock — guard is read-only; unchanged from C-974.
- [x] Matrix Completion Lock
- [x] Parallel Planning Advice — not applicable.
- [x] Perform The Adversarial Pass — symlinked main checkout and subfolder probes (below).
- [x] Re-check After Remediation
- [x] Verify File Reality — every file in `reviewed.files` exists; repo and template `guard.js` are byte-identical (`cmp`).
- [x] Verify Evidence And Gates
- [x] Collect Code Health Evidence
- [x] Complete The Issues Pass
- [x] Summarize Materiality
- [x] Review Code Style

## Closure Map

| Prior Issue | Status | Evidence |
|---|---|---|
| I-148 router guard blocks main checkout on Windows | Closed (verified) | Run 38036497392 at `3596b2c`: `windows-tests` success; `internal/init` passed with no FAIL line; `TestGuard_routerInLane` is not in the job's 28 skipped tests (it skips only without git). |
| I-147 native Windows evidence | Closed (verified) | Same run: full Go suite on `windows-latest` success; `ci` (`make ci`, ubuntu) success. |
| I-145 lone BEGIN marker loses user text | Closed (verified) | Fix verified in C-975; `TestUpgradeProjectAssets_claudeGuideHalfMarkerPairKeepsUserText` passes again on Linux (fresh gate) and Windows (run 38036497392). Code unchanged since. |
| I-146 Windows `.cmd` spawn; path with space | Closed (verified) | Fix verified in C-975 (run 38036069019); `session-start.js` unchanged since, and `TestSessionStartHook` passes again in run 38036497392. |

## Admission Ledger

| Re-check item | Prior claim | Frozen cell (C-974) | Result |
|---|---|---|---|
| Router guard, main checkout, Windows | I-148 fix | SC4 / T-118 router rule: main × Windows | Pass |
| Router guard, lane, Windows | I-148 fix must keep it | SC4 / T-118 router rule: lane × Windows | Pass (same test, lane case) |
| Router guard, outside a repo | named adjacent case in I-148 | T-118 router rule: not a repo | Pass |
| Router guard, Linux main/lane | I-148 Proof Needed | SC4 / T-118 router rule × Linux | Pass |
| Native Windows suite | I-147 | CFG-03 | Pass |
| Prior owner decisions | Task-check waivers only | no acceptance or exception recorded on O-043 | No `carried_forward` entry needed |

## Independent Probes

Real git repo and `git worktree` in a scratch directory, running `templates/project-v2/.claude/hooks/guard.js` with Node against an `Edit` of `.savepoint/router.md`:

| Case | Expected | Actual |
|---|---|---|
| main checkout root | allow | allow (no output) |
| main checkout subfolder | allow | allow |
| main reached through a symlink | allow | allow |
| worktree lane | block | deny: "do not edit .savepoint/router.md in a worktree lane" |
| not a git repo | allow | allow |

The new `canonical` helper falls back to `path.resolve` when `realpathSync.native` throws, and case-folds only on `win32`, so Linux comparison stays case-sensitive.

## Gates

- `make build && make test-full` at `3596b2c`: fresh, 2026-10-10T08:05:53Z to 08:06:09Z, go1.26.2 linux/amd64, cross-builds for linux, darwin, windows. Exit 0.
- `git diff --check 0fb8c37 3596b2c`: clean.
- GitHub Actions run 38036497392 (https://github.com/anipatke/savepoint/actions/runs/38036497392), head `3596b2c3be8c4ccef93ae4996dbf5e84fc625889`: `windows-tests` success, `ci` success, read with `gh run view` by this checker.

## Code Health

Snapshot `sha256:aa05e1fd094901f1a42a2677d76cd254391ba180b6c804369ee8f851cb998a13` (created after the full gate). Verdict: "Code Health does not block clearance." All five instances are optional and report no blocking finding.

## Issues

None. No materiality actions are required.

## Owner Validation Still Needed

T-115 and T-117 User Checks (`owner_validation.required: true`, no acceptance recorded). O-043 cannot close until the owner records them.

## Observations (non-blocking)

C-975's observations still apply (stray BEGIN above a complete pair; silent `unchanged` for half-marker `CLAUDE.md`).

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [x] STYLE-03 **Test branches** — the Windows branch of `inLane` now passes its test natively.
- [ ] STYLE-04 **Types document intent** — `isClaudeSkillPointer` still also matches hook scripts (`internal/init/upgrade.go`).
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [ ] STYLE-07 **One source of truth** — settings entries in `internal/init/claude_settings.go` still restate the template.
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data**
- [x] STYLE-10 **Small diffs**
