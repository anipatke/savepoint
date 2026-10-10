---
id: T-117
title: Start each Claude session knowing the Next step
objective: O-043
status: planned
depends_on: [{task: T-116, requires: clear}]
complexity_tier: medium
complexity_reason: First Node hook script and the safe settings.json install path that later hook and status line Tasks reuse.
owner_validation: {required: true}
planned_by: {role: planner, session: plan-o043-20261010}
---

# Start each Claude session knowing the Next step

## Outcome

A Claude Code session-start hook runs `savepoint resume` and gives Claude only the Next line, so Claude starts already knowing what to do without spending a tool call. The hook is silent and harmless when `savepoint` is missing. This Task also builds the shared pieces the later hook and status line Tasks reuse: the script that finds `savepoint`, and the safe way `.claude/settings.json` is installed.

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

Pending execution.

## Drift Notes

None yet.
