---
id: v2/D001-router-next-action-invalid-yaml
release: v2
status: resolved
severity: medium
title: "Board cannot parse router next_action containing an unquoted colon"
---

# D001: Board cannot parse router next_action containing an unquoted colon

## Symptom

`npx savepoint board` fails with `failed to parse router YAML: yaml: line 4: mapping values are not allowed in this context`.

The active `.savepoint/router.md` has a `next_action` value containing `status: planned` without quoting the full scalar.

## Expected Behavior

The board should load the router and display the current project state when the `next_action` prose contains a colon.

## Reproduction

1. Run `npx savepoint board` from the repository root.
2. Observe the router YAML parse error.
3. Inspect `.savepoint/router.md` line 18; the `next_action` value contains `status: planned` as an unquoted scalar.

## Impact

The board is unavailable until the router document is manually corrected.

## Fix Plan

Quote router `next_action` values containing YAML-significant punctuation and add regression coverage for an unquoted-colon router input or for the writer/template path that produced it.

## Acceptance Criteria

- [x] `npx savepoint board` loads the current router successfully.
- [x] Router handling remains valid when `next_action` contains a colon.
- [x] Regression coverage prevents emitting or accepting this malformed router shape.

## Resolution Notes

Quoted the active router `next_action` value and restored the router to `epic-task-breakdown`.

Added `TestWriteRouterState_quotesNextActionWithColon` to verify writer/parser round-tripping when `next_action` contains colons.

Verification passed: focused router tests, `make build`, and `make test`.
