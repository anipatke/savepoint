---
id: v1.2/D013-agent-invalid-task-completion-values
release: v1.2
status: resolved
severity: high
title: "Agent-written task values can break board updates"
---

# D013: Agent-Written Task Values Can Break Board Updates

## Symptom

After implementation, an AI agent may write a completed task with `status: complete` instead of the canonical `status: done`. The task file then fails lifecycle validation when the user tries to move or update the task in the board.

Separately, tasks can contain an overlong complexity reason. The product is supposed to enforce a word limit, but agent-written content can exceed that limit. When the user later moves that task across the board, validation rejects the file and the board update fails.

## Expected Behavior

Agent-authored task metadata should not be able to strand a user in an invalid board state.

- Task completion must use only `status: done`.
- Common agent synonyms such as `complete` should be caught early with a clear, repairable diagnostic or normalized at a controlled write boundary.
- Complexity reason length limits should be enforced before data is persisted or should produce an actionable repair suggestion that does not block unrelated board movement.
- The board should guide the user to a fix instead of failing during a move with a generic validation error.

## Reproduction

1. Complete a task through an AI agent session.
2. Let the agent update task frontmatter with `status: complete`.
3. Open the board and try to update or move the task.
4. Observe lifecycle validation fail because `complete` is not a valid task status.
5. Create or keep a task with a complexity reason that exceeds the intended word limit.
6. Try to move that task across the board.
7. Observe the board update fail because the overlong reason is validated only at move/update time.

## Impact

The user can be blocked from normal board operations by agent-authored metadata mistakes. The failure is especially disruptive because the invalid values are not new user intent during the board move; they are stale invalid task fields discovered only when the board attempts to write another update.

## Fix Plan

Strategic repair should make lifecycle and complexity validation a single, boundary-owned contract rather than relying on agent compliance.

1. Keep `done` as the only canonical completed task status in `internal/data`.
2. Add a compatibility normalization layer for persisted or incoming task frontmatter that maps `complete` and `completed` to `done` only when reading legacy or agent-authored files, while still writing canonical `done`.
3. Emit a doctor warning or repair suggestion when normalization occurs so the file can be cleaned up rather than silently carrying non-canonical vocabulary forever.
4. Centralize the complexity reason word-limit constant and validator in `internal/data`; remove duplicated or prompt-only enforcement.
5. Apply complexity reason validation at creation/update boundaries before writing task files, including board async update commands.
6. Prefer a repairable board error for existing invalid files: show the specific field, current value, allowed value or word limit, and suggested replacement.
7. Add tests that cover `status: complete` normalization/diagnostics, canonical writeback to `done`, overlong complexity reason rejection, and board movement behavior when an unrelated update touches a task with stale invalid metadata.
8. Update agent-facing instructions only as a supporting control, not as the primary fix. The product should remain robust when an agent writes a predictable synonym or too many words.

## Implementation Plan

- [x] Add data-layer compatibility for task status synonyms `complete` and `completed`, while keeping `done` as the only canonical write value.
- [x] Surface non-canonical completion statuses through doctor diagnostics or repair suggestions so compatibility does not hide stale file content forever.
- [x] Centralize the complexity reason word-limit constant and validation helper in `internal/data`.
- [x] Apply the centralized complexity reason validation before task persistence through board/update write paths.
- [x] Improve board-facing errors for known task metadata failures so the message names the field, current value or count, limit, and suggested repair.
- [x] Wrap long board status/error messages within the terminal width so repair guidance remains readable.
- [x] Add focused data, doctor, and board tests for status synonym handling, canonical writeback, complexity reason limits, and board movement/update behavior.
- [x] Run `make build && make test`.

## Acceptance Criteria

- [x] A task file containing `status: complete` no longer prevents the board from loading or moving the task.
- [x] Any compatibility handling writes back or recommends canonical `status: done`.
- [x] Task writes never emit `status: complete`.
- [x] Complexity reason word-limit rules live in one data-layer source of truth.
- [x] Overlong complexity reasons are rejected before persistence or reported with a concrete repair suggestion.
- [x] Board movement is not blocked by vague validation errors for these two known agent-authored mistakes.
- [x] Long board status/error messages wrap within the visible terminal width instead of running off screen.
- [x] Tests cover status synonym handling, canonical writeback, complexity reason limits, and board update behavior.

## Resolution Notes

Resolved. Task lifecycle loading now maps agent-written `status: complete` and `status: completed` to canonical `done`, while task writes remain strict and only emit canonical statuses. Doctor reports non-canonical completion aliases with repair guidance so stale files can be cleaned up.

Complexity reason validation now lives in `internal/data` as a word-count rule using `MaxComplexityReasonWords`; parse and write paths use the same validator. Board write errors for complexity failures include a repair hint, and long status/error messages wrap across bounded footer lines instead of running past the terminal width.

Verification: `go test ./internal/data ./internal/doctor ./internal/board`, `make build`, and `make test` passed.
