---
id: T-117
title: Start each Claude session knowing the Next step
objective: O-043
status: done
depends_on: [{task: T-116, requires: clear}]
complexity_tier: medium
complexity_reason: First Node hook script and the safe settings.json install path that the later hook Task reuses.
owner_validation:
    required: true
    accepted_check: ""
planned_by: {role: planner, session: plan-o043-20261010}
check_waiver:
    task: T-117
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-10T07:37:12Z"
---

# Start each Claude session knowing the Next step

## Outcome

A Claude Code session-start hook runs `savepoint resume` and gives Claude only the Next line, so Claude starts already knowing what to do without spending a tool call. The hook is silent and harmless when `savepoint` is missing. This Task also builds the shared pieces the later hook Task reuses: the script that finds `savepoint`, and the safe way `.claude/settings.json` is installed.

## User Check

In a freshly initialised project, start Claude Code and ask "what's next?" before anything else. Claude should answer from the Next line without running a command. Then temporarily make `savepoint` unavailable and confirm the session starts normally with no error.

## Done When

- `.claude/hooks/` contains a Node helper that finds `savepoint` (PATH, then `node_modules/.bin`) and a session-start script that adds only the Next line, about 60 tokens or fewer.
- The script finishes quickly. It adds nothing and exits successfully when `savepoint` is missing, `resume` fails, or the project is not a Savepoint project.
- Init writes `.claude/settings.json` with the hook only when the file is absent. When it exists, init and upgrade leave it byte-identical and report the exact entry to add (FS-01, TEST-03).
- Hook commands use the project-dir variable and work on Windows (CFG-02, CFG-03).
- Upgrade installs and refreshes the hook scripts through the manifest, as for skills.
- `make build && make test-fast` passes.

## Context Files

`main.go`; `internal/init/scaffold.go`; `internal/init/upgrade.go`; `internal/init/upgrade_test.go`; `internal/init/manifest.go`; `internal/init/v2_scaffold_test.go`; `cmd/resume.go`; `internal/resume/resume.go`; `bin/savepoint.js`.

## Design References

O-043 Architectural Considerations; Design sections 1, 6 (CLI surface), and 12.

## Guardrails

FS-01, FS-02, FS-03, FS-04, FS-05, CFG-02, CFG-03, TPL-02, TPL-04, DEP-01, TEST-01, TEST-02, TEST-03, TEST-04, TEST-08.

## Implementation Plan

1. Confirm the current Claude Code SessionStart hook contract (input, `additionalContext` output) from its documentation; record the source.
2. Write the helper and session-start scripts with no npm dependencies (DEP-01).
3. Add the settings template and the absent-only install and report path.
4. Test with `savepoint` present, missing, failing, and outside a project; test settings absent, present, and dry run.
5. Install in this repo; record the session-start token weight.

## Boundaries

No change to `resume` output. No permission allowlist or post-compaction hook (deferred).

## Technical Verification

Focused tests during iteration; handoff gate from AGENTS.md's Verification Policy. Native Windows CI evidence for the hook command (CFG-03).

## Technical Evidence

Executor evidence (claims for a fresh Check to verify).

**Hook contract source.** Claude Code hooks documentation (code.claude.com/docs/en/hooks), fetched 2026-10-10: SessionStart exit 0 adds stdout as context; JSON shape `hookSpecificOutput.{hookEventName: "SessionStart", additionalContext}`; exec form `command` + `args` with `${CLAUDE_PROJECT_DIR}`; Windows uses Git Bash or PowerShell, and exec form needs a real executable, so the hook runs `node <script>`. Only the first 100000 of 262610 characters of the page were read.

**Per-criterion outcomes.**
- Helper and script exist: `.claude/hooks/savepoint-find.js` (PATH, then `node_modules/.bin`; `.exe`/`.cmd` on Windows) and `session-start.js` (adds only the `Next` line, prefixed `Savepoint `). Output for this repo: `Savepoint Next action: Advance Task T-117 to its next stage.` (~15 tokens, under the 60 limit). Met.
- Silent and fast: `TestSessionStartHook` covers savepoint missing, `resume` failing, no Next line, and not a Savepoint project; each exits 0 with no output. `resume` timeout is 5 s; whole script wrapped in try/catch. Met.
- Settings written only when absent: `Scaffold` skips an existing file even with `--force` (`TestScaffold_keepsExistingSettings`); upgrade installs when absent, leaves a user file byte-identical in dry-run and real runs, and reports the exact entry in an `info` note; `init` prints the same advice on stderr (`TestUpgrade_claudeSettings`, `TestClaudeSettingsAdvice_silentWhenAbsentOrWired`). Met (FS-01, TEST-03).
- Project-dir variable and Windows: command is `node` with arg `${CLAUDE_PROJECT_DIR}/.claude/hooks/session-start.js`, paths built with `path.join`/`filepath`; `TestFindSavepoint_windowsNames` exercises the `.cmd` lookup. Native Windows CI run not yet observed. Partly met (CFG-02, CFG-03).
- Upgrade installs and refreshes hooks through the manifest: `isManifestPath` covers `.claude/hooks/*.js`; `TestScaffold_installsTrackedHooksAndSettings`, `TestUpgrade_claudeHooksFollowSkillPolicy` (edited hook kept with `.new`, missing hook restored). Met.
- `make build && make test-fast` passes (exit 0). Met.

**Commands run.** `make build`; `make test-fast` (exit 0); focused `go test ./internal/init/ -run 'Claude|Hook|SessionStart|Settings'` during iteration.

**Files read (Context Files).** `main.go`, `internal/init/scaffold.go`, `upgrade.go`, `manifest.go`. Not read: `upgrade_test.go` beyond helper names, `v2_scaffold_test.go`, `cmd/resume.go`, `internal/resume/resume.go`, `bin/savepoint.js` beyond its head.
**Extra reads (logged).** `internal/init/agents.go` (diff, to reuse the T-115 managed-block code), `internal/init/lifecycle_test.go` and `main_test.go` (two existing tests hard-coded the managed-file list and had to be updated), `.github/workflows/ci.yml` (confirm Node is on CI runners), `.claude/settings.local.json` head (confirm no existing `settings.json`).
**Files changed.** New: `templates/project-v2/.claude/hooks/{savepoint-find,session-start}.js`, `templates/project-v2/.claude/settings.json`, `internal/init/claude_settings.go`, `internal/init/claude_hooks_test.go`, repo copies `.claude/hooks/*.js` and `.claude/settings.json`. Edited: `internal/init/{manifest,scaffold,upgrade}.go`, `main.go`, `internal/init/lifecycle_test.go`, `main_test.go`.

**Limitations.** No real Claude Code session was started, so the User Check (ask "what's next?") is unverified and needs the owner. Native Windows CI not yet run. The session-start token weight for a fresh project was not separately measured (same one-line format). This repo's `.claude/settings.json` now activates the hook for the owner's next session. Hook tests skip when `node` is absent.


## Drift Notes

None yet.
