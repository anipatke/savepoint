---
id: T-118
title: Stop Claude breaking owner-only rules
objective: O-043
status: done
depends_on: [{task: T-117, requires: clear}]
complexity_tier: medium
complexity_reason: One pre-tool-use guard script with three rules, each needing a block and an allow case, including worktree-lane detection.
owner_validation:
    required: false
    accepted_check: ""
planned_by: {role: planner, session: plan-o043-20261010}
check_waiver:
    task: T-118
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-10T07:40:02Z"
---

# Stop Claude breaking owner-only rules

## Outcome

A Claude Code pre-tool-use hook blocks three owner-authority breaches before they happen, each with a one-line reason: an edit that sets a Task to `status: done` (DATA-05), an edit to `.savepoint/router.md` inside a worktree lane, and an agent run of a `savepoint` command other than `resume`, `create-task`, or `health check`. Nothing is added to context unless a rule fires.

## User Check

In a scratch project, ask Claude to mark a planned Task done. The edit should be refused with a short reason, and you can still change it yourself on the board.

## Done When

- File edits that would leave a Task file with `status: done` are blocked. Edits to other fields of a done Task are allowed.
- Router edits are blocked only when `git rev-parse --git-dir` differs from `--git-common-dir`, and allowed on the main checkout.
- `savepoint` / `npx savepoint` / `./savepoint` commands other than `resume`, `create-task`, and `health check` are blocked for the agent. Owner commands typed with `!` are not affected.
- Each block message is one line naming the rule and who may do it.
- The script fails open (allows) on its own errors, so it never wedges a session. Every rule has block and allow tests (TEST-02).
- Prose rules in AGENTS.md and the skills are unchanged.
- `make build && make test-fast` passes.

## Context Files

`internal/init/upgrade.go`; `internal/init/upgrade_test.go`; `internal/init/v2_scaffold_test.go`; `templates/project-v2/AGENTS.md`.

## Design References

O-043 Architectural Considerations; AGENTS.md Terminology, Worktree Lanes, and CLI Rules.

## Guardrails

DATA-05, FS-01, FS-02, CFG-02, CFG-03, TPL-02, DEP-01, TEST-01, TEST-02, TEST-04, TEST-08.

## Implementation Plan

1. Confirm the current PreToolUse hook contract (tool input shape for Edit, Write, and Bash; deny output) from documentation; record the source.
2. Write the guard script, reusing the T-117 helper.
3. Register it in the settings template and in the upgrade report snippet.
4. Test each rule's block and allow paths, the lane and main checkouts, and fail-open behaviour.

## Boundaries

No guard on Check record writes (a Check session cannot be told apart reliably). No removal of prose rules.

## Technical Verification

Focused tests during iteration; handoff gate from AGENTS.md's Verification Policy. Native Windows CI evidence (CFG-03).

## Technical Evidence

**Hook contract (plan step 1).** Source: https://code.claude.com/docs/en/hooks (fetched 2026-10-10). PreToolUse stdin carries `tool_name`, `tool_input`, `cwd`; Bash input has `command`. Deny is exit 0 with `hookSpecificOutput.permissionDecision: "deny"` plus `permissionDecisionReason`. The page excerpt did not show the Edit/Write/MultiEdit schemas; `file_path`, `content`, `old_string`/`new_string`/`replace_all` and `edits[]` are from Claude Code's tool definitions, not the page. The script reads them defensively and allows on anything unexpected.

**Extra reads (logged).** `.claude/hooks/*.js`, `templates/project-v2/.claude/settings.json`, `internal/init/claude_settings.go`, `claude_hooks_test.go`, `manifest.go`, `main_test.go` (T-117 wiring the guard had to reuse); the Objective; `agent-skills/savepoint-task/SKILL.md`. Context Files `upgrade_test.go` and `v2_scaffold_test.go` were not needed.

**Changed.** New `templates/project-v2/.claude/hooks/guard.js` (copied to `.claude/hooks/guard.js`); `templates/project-v2/.claude/settings.json` (PreToolUse entry, matcher `Edit|Write|MultiEdit|Bash`; copied to this repo's `.claude/settings.json`, which was identical to the old template); `internal/init/claude_settings.go` (advice and upgrade note now name only the missing entries: session-start and/or guard); new `internal/init/claude_guard_test.go`; `claude_hooks_test.go` and `main_test.go` (guard.js is a tracked hook). Hooks are tracked by the existing `.claude/hooks/*` manifest rule, so no manifest code changed. AGENTS.md and skills prose untouched.

**Per criterion.**
- Done-block: `TestGuard_taskDone` blocks Edit, Write, MultiEdit and a quoted `'done'` that would leave `status: done`; allows other statuses, a body line containing the words, other files, and edits to an already-done Task. Pass.
- Router-in-lane: `TestGuard_routerInLane` uses a real git worktree: blocked in the lane, allowed on the main checkout and outside a repo. Pass.
- CLI: `TestGuard_savepointCommands` (21 cases) blocks `savepoint`, `./savepoint`, `npx savepoint`, `.exe`, env prefixes, chained and newline commands for board, init, upgrade-assets, migrate, doctor, health setup/report and a bare call; allows resume, create-task, health check and non-savepoint commands such as `echo savepoint board`. Owner `!` commands do not go through the tool hook, so are unaffected by construction (not testable here). Pass.
- One-line reasons naming the rule and who may act: asserted single-line in `runGuard`; wording checked by the tests above. Pass.
- Fails open: `TestGuard_failsOpen` (empty and invalid input, missing tool input, unknown tool, missing file, no command) all exit 0 with no output. Pass.
- Prose rules unchanged. Pass.
- Settings upgrade: `TestSettingsNote_namesOnlyMissingEntries` plus the existing settings tests; a user's settings file stays byte-identical. Pass.
- `make build && make test-fast`: both passed (run 2026-10-10). Pass.

**Limitations.** Native Windows CI evidence (CFG-03) not gathered here; the script uses `path`, `node:child_process` and no shell, but was run only on Linux (WSL). The guard cannot see a command hidden inside `bash -c "..."`, `eval`, or a script, nor a Task edited through Bash (e.g. `sed -i`); it targets honest slips, not a determined bypass. The Edit/Write field names were not confirmed on the docs page. Full gate (`make test-full`) was not run; not required for ordinary handoff.

## Drift Notes

None yet.
