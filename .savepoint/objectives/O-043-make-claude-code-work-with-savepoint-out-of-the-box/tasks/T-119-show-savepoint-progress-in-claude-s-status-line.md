---
id: T-119
title: Show Savepoint progress in Claude's status line
objective: O-043
status: planned
depends_on: [{task: T-117, requires: clear}]
complexity_tier: low
complexity_reason: One status line script reusing the T-117 helper and settings install path.
owner_validation: {required: true}
planned_by: {role: planner, session: plan-o043-20261010}
---

# Show Savepoint progress in Claude's status line

## Outcome

Claude Code's status line shows the selected Savepoint record and a short form of the Next action, for example `▸ T-115 build · O-043`. When the Next line needs an owner decision (`Accept`, `Close`, `Choose`, `Done`, `Blocked`), it says so plainly, for example `⚑ Your call: Accept O-043`. The status line is never sent to Claude, so it costs no tokens.

## User Check

In a scratch project, open Claude Code and look below the prompt. Then move the router to a state waiting on you and confirm the status line changes to the "your call" form.

## Done When

- A Node status line script renders one short line from `savepoint resume`, with control sequences removed (ARCH-05).
- It shows nothing when `savepoint` is missing or the folder is not a Savepoint project, and renders fast enough for frequent refresh.
- The settings template sets `statusLine` only when the settings file is absent. A user's existing status line is never replaced; upgrade reports the entry instead (FS-01).
- Tests cover a normal selection, an owner decision, a missing tool, and a non-project folder.
- `make build && make test-fast` passes.

## Context Files

`internal/init/upgrade.go`; `internal/init/upgrade_test.go`; `internal/init/v2_scaffold_test.go`; `internal/resume/resume.go`.

## Design References

O-043 Architectural Considerations; Design section 6 (CLI surface).

## Guardrails

ARCH-05, FS-01, CFG-02, CFG-03, DEP-01, TEST-01, TEST-02, TEST-04, TEST-08.

## Implementation Plan

1. Confirm the current status line contract (stdin JSON, one line of output) from documentation; record the source.
2. Write the script, reusing the T-117 helper.
3. Add it to the settings template and the upgrade report snippet.
4. Test the four cases.

## Boundaries

No change to `resume` output and no new CLI command.

## Technical Verification

Focused tests during iteration; handoff gate from AGENTS.md's Verification Policy.

## Technical Evidence

Pending execution.

## Drift Notes

None yet.
