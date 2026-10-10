---
id: T-118
title: Stop Claude breaking owner-only rules
objective: O-043
status: planned
depends_on: [{task: T-117, requires: clear}]
complexity_tier: medium
complexity_reason: One pre-tool-use guard script with three rules, each needing a block and an allow case, including worktree-lane detection.
owner_validation: {required: false}
planned_by: {role: planner, session: plan-o043-20261010}
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

Pending execution.

## Drift Notes

None yet.
