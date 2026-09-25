---
id: v1.2/D001-router-priority-cross-epic
release: v1.2
status: resolved
severity: medium
title: "Router priority icon appears on matching task numbers across epics"
---

# D001: Router Priority Icon Appears Across Epics

## Symptom

The green task priority icon appears on every task with the same task number across epics in a release.

## Expected Behavior

The green task priority icon should appear only on the current router task in the router epic and router release.

## Reproduction

1. Open a release that has multiple epics with the same task number, such as `T001`.
2. Set router priority to one task in one epic.
3. Navigate to another epic in the same release that also has a task with that number.
4. Observe that the green priority icon appears on the matching task number even though it is not the router epic task.

## Impact

The board can suggest the wrong active task, which weakens the router priority signal and makes cross-epic work ambiguous.

## Fix Plan

- Compare router priority using release, epic, and task identity together.
- Avoid matching only the short task number when router epic/release context is available.
- Add a regression test with duplicate task numbers across epics in the same release.

## Acceptance Criteria

- [x] A router-priority task in one epic does not mark a same-numbered task in another epic.
- [x] The intended task still shows the green priority icon.
- [x] Matching remains scoped to the router release.
- [x] Regression tests cover duplicate task numbers across epics.

## Resolution Notes

`isRouterPriority` in `internal/board/card.go` already compared release, epic (via `shortID`), and task — the cross-epic leak never made it to ship. Regression tests added: `TestRenderCard_routerSameTaskNumberDifferentEpicNoMatch` (cross-epic) and `TestRenderCard_routerSameTaskNumberDifferentReleaseNoMatch` (cross-release).
