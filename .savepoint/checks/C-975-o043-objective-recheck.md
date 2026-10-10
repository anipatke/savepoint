---
id: C-975
scope: {kind: objective, id: O-043}
result: NEEDS WORK
checked_by: {role: checker, session: check-o043-20261010}
executed_session: executor-o043-remediation-unrecorded
checked_at: '2026-10-10T07:59:00Z'
reviewed:
  base_commit: 0e95269
  head_commit: 0fb8c37
  files:
    - internal/init/agents.go
    - internal/init/claude_guide_test.go
    - internal/init/claude_hooks_test.go
    - internal/init/claude_guard_test.go
    - templates/project-v2/.claude/hooks/session-start.js
    - templates/project-v2/.claude/hooks/guard.js
    - .claude/hooks/
  dependencies: []
issues: [I-148]
unmet: [CFG-03, T-118-DW2]
supersedes: C-974
---

# C-975: O-043 Full Objective re-check

NEEDS WORK. Re-check of C-974 under its frozen scope lock. This checker session did not build the O-043 work or the remediation; the remediation was done in a separate executor session and committed as `0fb8c37` (clean working tree, pushed). Both fixes are verified. The native Windows run that I-147 asked for has now happened, and it fails on one frozen cell: the router guard on the main checkout (I-148).

No `health_snapshot`: the remediation is code-only, and C-974's official snapshot (`sha256:4733d3666ef82425586ed5bc6e97ba93bc7b517853ffc7cc6d16d58c6daca1ba`) was collected before it. A clearing re-check should collect a fresh one.

## Closure Map

"Fix verified" means the remediation passes its frozen cells. The Issue stays `open` because the runtime accepts a `verified` resolution only from a `CLEAR` Check; the clearing re-check closes it.

| Prior Issue | Status | Evidence |
|---|---|---|
| I-145 lone BEGIN marker loses user text | Fix verified | Original reproduction re-run via `go test -overlay` on `0fb8c37`: lone BEGIN, lone END and reversed markers each stay byte-identical across three upgrades (all `unchanged`) and two `MergeClaudeGuide` init merges; plain user file still merges once, then `unchanged`. `TestUpgradeProjectAssets_claudeGuideHalfMarkerPairKeepsUserText` (3 subtests) passes on Linux and native Windows. |
| I-146 Windows `.cmd` spawn untested; path with space | Fix verified | Windows skip removed; `fakeSavepoint` writes a `.cmd` on Windows. Run 38036069019 `windows-tests`: all 7 `TestSessionStartHook` subtests PASS, including `runs savepoint from a directory with a space`; `TestFindSavepoint_windowsNames` PASS. Linux probe: quoted path with `shell: true` → status 0. |
| I-147 native Windows evidence | Still open | Run 38036069019 at `0fb8c37`: `ci` success, `windows-tests` failure (I-148). |

## Admission Ledger

| Re-check item | Prior claim | Frozen cell (C-974) | Allowed result |
|---|---|---|---|
| Half-marker merge | I-145 fix | SC5 / T-115 merge: half pair, lone BEGIN and lone END | Pass or Issue |
| Windows `.cmd` spawn | I-146 fix | T-117 Windows: `.cmd` spawn; path with space | Pass or Issue |
| Native Windows suite | I-147 | CFG-03 native Windows: full Go suite on `windows-latest` | Pass or Issue |
| Router guard, main checkout, Windows | none (new result from that suite) | SC4 / T-118 router rule: main × platform Windows (CFG-03 axis) | Issue I-148 |
| Stray BEGIN above a complete pair | none | no exact cell (duplicate markers were not an axis) | Observation only |
| Prior owner decisions | Task-check waivers only | no acceptance or exception recorded on O-043 | No `carried_forward` entry needed |

## Re-run Of The Frozen Matrix

All other C-974 matrix rows were re-run against `0fb8c37`, with these results:

- The fresh full gate passed (below), and every `internal/init` Claude/hook/guard test passes on Linux.
- On native Windows, every `internal/init` test passed except `TestGuard_routerInLane`.
- The repo hook copies are byte-identical to the templates (`cmp`).
- The live session evidence from C-974 still stands: the guide import, the Next line at session start, the four skills, and the done and CLI guard blocks.

## Gates

- `make build && make test-full` at `0fb8c37`: fresh, 2026-10-10T07:56:15Z to 07:56:32Z, go1.26.2 linux/amd64, `-count=1`, cross-builds for linux, darwin, windows. Exit 0.
- GitHub Actions run 38036069019 (head `0fb8c3777884e4bd7e9f0b61db881aa9dfa3991e`): `ci` (`make ci`, ubuntu) success; `windows-tests` failure, `TestGuard_routerInLane` at `claude_guard_test.go:151`.

## Issues

| Issue | Likelihood | Impact | Materiality | Recommendation |
|---|---|---|---|---|
| I-148 router guard blocks main checkout on Windows | Medium: fails in Windows CI; whether real Windows sessions hit it depends on the path form Claude Code passes, which is unverified | Medium: CFG-03 gate fails; a Windows owner's agent could be wrongly refused router edits on main | Medium | Fix now, narrow: normalise both git paths, then confirm in Windows CI |

## Convergence

This is the one full re-check after remediation. Next comes one targeted fix of I-148 and one targeted re-check limited to the I-148 and I-147 cells. If that still fails, stop and ask the owner to choose.

## Owner Validation Still Needed

Unchanged from C-974: T-115 and T-117 User Checks (`owner_validation.required: true`, no acceptance recorded).

## Observations (non-blocking)

- A stray BEGIN marker placed above a complete Savepoint pair in `CLAUDE.md` is still paired with the real END, and the text between them is replaced. The shared `replaceManagedBlock` treats the text between the first BEGIN and the first END as the managed region, and AGENTS.md behaves the same way. This was not in the frozen lock.
- A half-marker `CLAUDE.md` is now left alone and reported `unchanged`, so the user gets no import and no prompt to repair the marker. A conflict or info note would be clearer.
- C-974's other observations still apply.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [ ] STYLE-03 **Test branches** — the Windows branch of `inLane` fails its test (I-148).
- [ ] STYLE-04 **Types document intent** — `isClaudeSkillPointer` still also matches hook scripts (`internal/init/upgrade.go`).
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [ ] STYLE-07 **One source of truth** — settings entries in `internal/init/claude_settings.go` still restate the template.
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data**
- [x] STYLE-10 **Small diffs**
