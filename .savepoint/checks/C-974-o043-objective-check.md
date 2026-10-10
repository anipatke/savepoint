---
id: C-974
scope: {kind: objective, id: O-043}
result: NEEDS WORK
checked_by: {role: checker, session: check-o043-20261010}
executed_session: executor-o043-unrecorded
checked_at: '2026-10-10T07:53:00Z'
health_snapshot: sha256:4733d3666ef82425586ed5bc6e97ba93bc7b517853ffc7cc6d16d58c6daca1ba
reviewed:
  base_commit: 0e95269
  head_commit: 0e95269
  files:
    - main.go
    - main_test.go
    - CLAUDE.md
    - AGENTS.md
    - internal/init/agents.go
    - internal/init/scaffold.go
    - internal/init/upgrade.go
    - internal/init/manifest.go
    - internal/init/claude_settings.go
    - internal/init/lifecycle_test.go
    - internal/init/claude_guide_test.go
    - internal/init/claude_skills_test.go
    - internal/init/claude_hooks_test.go
    - internal/init/claude_guard_test.go
    - templates/project-v2/CLAUDE.md
    - templates/project-v2/.claude/settings.json
    - templates/project-v2/.claude/hooks/
    - templates/project-v2/.claude/skills/
    - .claude/settings.json
    - .claude/hooks/
    - .claude/skills/
  dependencies: []
issues: [I-145, I-146, I-147]
unmet: [FS-01, FS-02, CFG-02, CFG-03, O-043-SC5, T-117-DW4]
supersedes: null
---

# C-974: O-043 Full Objective Check

NEEDS WORK. Fresh session (started after `/clear`; did not build T-115 to T-118). Mode: Full. All O-043 work is uncommitted in the working tree on top of `0e95269`. The four Tasks are `done` with owner Task-check waivers (board-owner, 2026-10-10T07:27:09Z, 07:30:21Z, 07:37:12Z, 07:40:02Z). This Check reviews all four. The working tree also carries unrelated planning edits (O-041 deletion, O-042, G-002 boundary widening); only the G-002 widening was read, as context.

Most of the Objective works. Two defects and one missing piece of evidence block clearance: a user-text loss in a rare `CLAUDE.md` state (I-145), an untested and fragile Windows launch path in the session-start hook (I-146), and no native Windows CI run (I-147).

## Full Check Progress Checklist

- [x] Establish Scope
- [x] Freeze The Check Scope
- [x] Turn Acceptance Into Invariants
- [x] Build The Mandatory Coverage Matrix
- [x] Finite External-Boundary Matrix — subprocesses in scope: `savepoint resume` (session start) and `git rev-parse` (guard).
- [x] Workflow And Side-Effect Check Lock — init and upgrade writes of `CLAUDE.md`, pointers, hooks, settings.
- [x] Matrix Completion Lock
- [x] Parallel Planning Advice — not a feature of this scope; no finding.
- [x] Perform The Adversarial Pass
- [x] Re-check After Remediation — not applicable: initial Check.
- [x] Verify File Reality
- [x] Verify Evidence And Gates
- [x] Collect Code Health Evidence
- [x] Complete The Issues Pass
- [x] Summarize Materiality
- [x] Review Code Style

## Scope Lock

1. Criteria: O-043 Success Conditions SC1 to SC6; Done When of T-115, T-116, T-117, T-118. Guardrails named by those Tasks: FS-01 to FS-06, DATA-05, TPL-01, TPL-02, TPL-04, ARCH-04, CFG-02, CFG-03, DEP-01, TEST-01 to TEST-04, TEST-08.
2. Entry points: `savepoint init` (`Scaffold`, `ClaudeSettingsAdvice`), `savepoint upgrade-assets` (`upgradeClaudeGuide`, `upgradeClaudeSettings`, manifest path rules), shipped `session-start.js`, `guard.js`, `savepoint-find.js`, `settings.json`, skill pointers, the `main.go` embed.
3. Relied-on runtime: Claude Code's CLAUDE.md `@` import, skill discovery, SessionStart and PreToolUse hook contracts (exec form `command` + `args`); Node `child_process`; `git rev-parse`.
4. Axes: file state (absent, user-only, full marker pair, half pair, already-imported, unreadable, CRLF); run sequence (init, upgrade, repeat upgrade, dry run); manifest state (missing, unedited-old, edited); hook environment (savepoint on PATH, in `node_modules/.bin`, missing, failing, not a project, Windows `.cmd`, path with space); guard inputs (Edit, Write, MultiEdit, Bash, other tools; done vs other status; lane vs main; allowed vs owner-only commands; prefixed and chained commands); platform (Linux observed, Windows required by CFG-03).
5. Admission: a supported path that violates a named criterion or guardrail. Deliberate bypasses (`bash -c`, `sed -i`, `env`/`sudo` prefixes) are outside scope by T-118's stated aim (honest slips) and are observations.

## Coverage Matrix

| Row | Cells | Result | Evidence |
|---|---|---|---|
| SC1 guide loads with no prompt | this repo, live session | Proven (this repo) | This Check session received AGENTS.md through `CLAUDE.md`'s `@AGENTS.md` import with no pasted prompt. Fresh-project form is the T-115 owner User Check. |
| SC1 / T-115 merge | absent, user-only, full pair, already-imported, unreadable, dry run, repeat, CRLF | Proven | `TestScaffold_writesClaudeGuideWithImport`, `_keepsUserClaudeText`, `_noSecondImport`, `TestUpgradeProjectAssets_claudeGuide{Lifecycle,RefreshesOnlyBlock,AlreadyImportedUnchanged,MissingIsCreated,WriteFailureReported}`; independent overlay probe: CRLF user file, second run `unchanged`. |
| SC5 / T-115 merge | half pair: lone BEGIN | **Issue I-145** | Overlay probe: second upgrade deletes the user line after the orphan BEGIN. |
| SC5 / T-115 merge | half pair: lone END | Passed | Reasoned from `replaceManagedBlock`: END before BEGIN is unmarked, appended block then carries `@AGENTS.md`, so later runs return the file unchanged. |
| SC2 / T-116 pointers | fresh install, parity, body, only four, repo byte-identical, missing, unedited-old, edited, dry run | Proven | `TestClaudeSkillPointers_*` (4), `TestScaffold_installsTrackedClaudeSkillPointers`, `TestUpgrade_claudeSkillPointers` (4 subtests); `cmp` repo vs template, identical; live: this session's skill list showed all four `savepoint-*` skills and `savepoint-check` launched through its pointer. |
| SC3 / T-117 hook | Next line only, ≤60 tokens | Proven | Live: this session started with `Savepoint Next action: Record the Objective O-043 integration Check.` (~15 tokens). `TestSessionStartHook/adds only the Next line`. |
| SC3 / T-117 hook | missing, resume fails, no Next, not a project, `node_modules/.bin` | Proven (non-Windows) | `TestSessionStartHook` subtests, all PASS, none skipped on Linux. |
| T-117 Windows | `.cmd` lookup | Proven | `TestFindSavepoint_windowsNames`. |
| T-117 Windows | `.cmd` spawn; path with space | **Issue I-146** | Windows subtests skip with a false reason; `shell: true` path with space fails (Linux reproduction of the mechanism). |
| T-117 settings | absent → installed; present → byte-identical + exact entry; dry run; `--force` on init | Proven | `TestScaffold_keepsExistingSettings`, `TestUpgrade_claudeSettings`, `TestClaudeSettingsAdvice_silentWhenAbsentOrWired`, `TestSettingsNote_namesOnlyMissingEntries`. |
| T-117 hooks via manifest | install, edited kept with `.new`, missing restored | Proven | `TestScaffold_installsTrackedHooksAndSettings`, `TestUpgrade_claudeHooksFollowSkillPolicy`, `TestRepoHooksAreByteIdenticalToTemplates`. |
| SC4 / T-118 done rule | Edit, Write (new and existing), MultiEdit, quoted, extra spaces, CRLF, relative path, body-only text, already-done file | Proven | Live: an Edit setting a scratch Task to `status: done` was denied with the one-line reason. Independent probe script (29 cases) plus `TestGuard_taskDone`. |
| SC4 / T-118 router rule | lane, main, Write to lane path from main cwd, outside repo | Proven | Independent probe with a real scratch `git worktree`: lane BLOCK, main allow, cross-cwd Write BLOCK; `TestGuard_routerInLane`. |
| SC4 / T-118 CLI rule | board, doctor, init, upgrade-assets, health setup, quoted subcommand, `&&`, `$()`, `./savepoint`; allows resume (with flags or dir), create-task, health check, non-savepoint words | Proven | Live: `savepoint doctor` denied with the one-line reason. Probe script and `TestGuard_savepointCommands`. |
| T-118 fail open | empty, invalid JSON, missing input, unknown tool | Proven | `TestGuard_failsOpen`. |
| SC5 upgrade parity | `CLAUDE.md`, pointers, hooks reach existing projects; settings advised only | Proven, except I-145 | Tests above; settings handling matches the Objective's Architectural Considerations. |
| SC6 token weight | this repo, fresh project | Proven (estimate) | T-115 records bytes before and after for both (est. bytes/4, no tokenizer); T-117 records the ~15-token Next line. This Check adds: the four skill names and descriptions are 1,269 bytes (~320 tokens), loaded at start by design (SC2). |
| CFG-03 native Windows | full Go suite on `windows-latest` | **Unverified, I-147** | No CI run covers the uncommitted work. |

### External-Boundary Matrix

| Cell | session-start → `savepoint resume` | guard → `git rev-parse` |
|---|---|---|
| target discovery | PATH then `node_modules/.bin`; tested | `git` on PATH |
| unavailable | silent, exit 0; tested | fails open (not in lane); tested outside a repo |
| non-success | silent; tested | not in lane; allow |
| timeout | 5 s `spawnSync` timeout → silent; not separately tested, accepted by reading | none set; `rev-parse` is local and fast |
| malformed output | no `Next` line → silent; tested | path resolve; `realpathSync` error caught → allow |
| redirect, retry, partial effects | not applicable: read-only, no network | not applicable |
| secret-safe output | only the Next line is emitted | only fixed messages |
| Windows `.cmd` | **I-146** | not applicable |

### Workflow And Side-Effect Lock

| Order | Operation | Side effect | Failure | Final state | Oracle |
|---|---|---|---|---|---|
| 1 | read existing `CLAUDE.md` | none | unreadable → `failed` + error | file untouched | `..._WriteFailureReported` |
| 2 | merge block | in memory | none | — | byte comparison in tests and overlay probe |
| 3 | dry run | none | — | no write | `..._claudeGuideLifecycle` |
| 4 | atomic write | file replaced | write error → `failed` | prior file kept (atomic) | existing `AtomicWrite` tests |
| 5 | settings | write only when absent | read error → `failed` | user file byte-identical | `TestUpgrade_claudeSettings` |
| 6 | pointers and hooks | manifest-tracked write, `.new` on edit | as for skills | edited file kept | `TestUpgrade_claudeSkillPointers`, `..._claudeHooksFollowSkillPolicy` |
| 7 | init advice | stderr line | none | — | `TestClaudeSettingsAdvice_*` |

## Gates

- `make build && make test-full` — fresh, 2026-10-10T07:48:27Z to 07:48:45Z, go1.26.2 linux/amd64, `-count=1`, cross-builds for linux, darwin, windows. Exit 0.
- `go test -count=1 -v ./internal/init/ -run 'Claude|Hook|SessionStart|Settings|Guard|FindSavepoint|Lifecycle|Scaffold_'` — PASS, no skips on Linux.
- `git diff --check` — clean.
- Node v22.22.2 for hook probes.
- Native Windows CI: not run (I-147).

## Code Health

`savepoint health check O-043` after the full gate: snapshot `sha256:4733d3666ef82425586ed5bc6e97ba93bc7b517853ffc7cc6d16d58c6daca1ba` (created). "Code Health does not block clearance." All five instances optional, no blocking finding. Supporting evidence only.

## Guardrails

- FS-01, FS-02: Issue I-145 (lone BEGIN marker). Otherwise satisfied: user `CLAUDE.md` bytes, `settings.json`, edited pointers and hooks are kept.
- FS-03, FS-04: satisfied (dry run writes nothing; repeat upgrade unchanged, including CRLF).
- FS-05: satisfied (`filepath.Join`, `path.join`).
- FS-06: satisfied for the new paths (unreadable `CLAUDE.md` fails clearly).
- DATA-05: satisfied by the guard as a convenience; prose rules unchanged.
- TPL-01: canonical skills untouched. TPL-02: guidance matches code. TPL-04: every new file has an upgrade path; settings are advice-only by design.
- ARCH-04: Codebase Map row for `internal/init/` updated for `CLAUDE.md`; it does not mention the Claude skill pointers, hooks, or settings (observation).
- CFG-02, CFG-03: Issues I-146 and I-147.
- DEP-01: no new dependencies (Node built-ins only).
- TEST-01 to TEST-04: satisfied except the missing half-marker and Windows `.cmd` tests named in I-145 and I-146. TEST-08: fresh full gate passed.

## Issues

| Issue | Likelihood | Impact | Materiality | Recommendation |
|---|---|---|---|---|
| I-145 lone BEGIN marker deletes user text in `CLAUDE.md` | Low: needs a hand-damaged marker pair | High: silent loss of user-authored text | Medium | Fix now; narrow (report a half pair as a conflict) plus two tests |
| I-146 Windows `.cmd` spawn untested; path with space fails | Medium: Windows npm installs use `.cmd` shims; spaces in profile names are common | Low: Next line silently missing, nothing breaks | Medium | Fix now; quote the path and run the Windows subtests with a `.cmd` fake |
| I-147 no native Windows CI | Certain until pushed | Medium: CFG-03 Blocker gate | Medium | Commit, push, record a passing `windows-tests` run after I-146 |

`unmet` IDs: `O-043-SC5` is Success Condition 5; `T-117-DW4` is T-117 Done When 4 ("work on Windows").

## Owner Validation Still Needed

- T-115 and T-117 declare `owner_validation.required: true` with no `accepted_check`. Their User Checks (fresh-project "what should I do next?" with no pasted prompt; `savepoint` made unavailable) are not yet recorded. This repo's own live session supports both.
- After remediation: a Full Objective re-check.

## Observations (non-blocking)

- The guard does not see `env`, `sudo`, `command`, or `time` prefixes, `go run . <cmd>`, `bash -c`, or Bash edits such as `sed -i` to a Task or the router. This matches T-118's stated aim (honest slips) and its Limitations.
- `savepoint --version` is blocked as an owner-only command; harmless but stricter than needed.
- A Task file in YAML flow style (`{status: done}`) is not detected; Savepoint does not write that form.
- Objective `status: done` edits are not guarded; out of the Objective's stated three rules.
- `isClaudeSkillPointer` (`internal/init/upgrade.go:629`) also matches hook scripts; the name says less than it does.
- The exact settings entries in `claude_settings.go` duplicate `templates/project-v2/.claude/settings.json`; a change to one could drift from the other. No test ties them.
- T-115 notes the `@AGENTS.md` import does not resolve for a project whose guide is `agents.md` on a case-sensitive disk, though `FindAgentGuide` accepts that casing.
- The session-start hook's 5 s timeout is not covered by a test.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [ ] STYLE-03 **Test branches** — half-marker branch of `replaceManagedBlock` for `CLAUDE.md` and the Windows `.cmd` spawn branch in `session-start.js` are untested (I-145, I-146).
- [ ] STYLE-04 **Types document intent** — `isClaudeSkillPointer` (`internal/init/upgrade.go:629`) also returns true for hook scripts.
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [ ] STYLE-07 **One source of truth** — settings entries in `internal/init/claude_settings.go:24` restate `templates/project-v2/.claude/settings.json`.
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data**
- [x] STYLE-10 **Small diffs**
