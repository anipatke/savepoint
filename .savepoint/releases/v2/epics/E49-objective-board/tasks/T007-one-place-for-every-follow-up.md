---
id: E49-objective-board/T007-one-place-for-every-follow-up
title: One place for every follow-up
status: done
objective: Replace the separate defect and audit-register overlays with one Issues view over index.Issues, filtered by type, with Defect as a type rather than a collection.
depends_on:
    - E49-objective-board/T004-navigate-a-project-by-its-objectives
complexity_tier: medium
complexity_reason: One overlay over an already-validated record set, with filtering, history, and duplicate links to render.
---

# T007: One place for every follow-up

## Problem

V1 keeps follow-up in three places: defects under a release, findings in an audit register, and drift notes appended to task files. Each has its own overlay, its own vocabulary, and its own idea of what "open" means. V2 merged them into one Issue model with a `type` field, and section 14 is explicit that finding and defect infrastructure becomes Issues and that register bookkeeping becomes derived listings. This task is where that merge becomes visible: one overlay, five types, and no separate defect surface anywhere in the V2 package.

That the type is descriptive matters for how it is presented. A `defect` is not automatically blocking and a `guardrail` Issue is not automatically advisory — section 6 says type is never sufficient by itself to block a Task. So the filter organises the list; it does not encode severity, and the overlay must not imply an ordering of importance that the records do not carry.

An Issue's identity survives repair, recheck, deferral, and reopening — that is the point of the global ID — and its `history` is an append-only record of exactly those events. Rendering history in recorded order is what makes a reopened Issue legible: the same ID, observed, repaired, rechecked, reopened, with dates and actors. A resolved Issue carries a disposition that means three different things: `verified` was proven fixed by a Check, `accepted` was a decision to live with it, and `duplicate` points at a canonical Issue and is not proof of anything. Collapsing those into "closed" would discard the distinction the model exists to preserve.

The links are already indexed. `index.TaskIssues` and `index.CheckIssues` were built at load, and `LoadV2Index` has already proven every reference resolves and that Check-to-Issue pairing agrees. The overlay reads those maps; it does not walk Issues to rebuild a relation, and a duplicate's canonical target is a navigable link rather than a printed string.

Nothing here writes. Resolving an Issue is a checker's act on evidence, not a keypress on a board.

## Context Files

- `internal/board/v2/model.go`
- `internal/board/v2/update.go`
- `internal/board/v2/view.go`
- `internal/board/defect_overlay.go`
- `internal/board/defect_detail.go`
- `internal/board/audit_overlay.go`
- `internal/board/audit_backlinks.go`
- `internal/data/issue_v2.go`
- `internal/data/project.go`
- `internal/data/check_v2.go`
- `internal/styles/styles.go`
- `.savepoint/releases/v2/v2-Design.md`

## Acceptance Criteria

- [x] One Issues overlay lists every Issue in `index.Issues` in stable `I###` order, each row showing ID, title, type, status, and severity when recorded.
- [x] The overlay filters by each of the five types — defect, drift, guardrail, verification, other — and back to unfiltered, with the active filter visible.
- [x] Defect appears as a type within this overlay; no separate defect overlay, defect column, defect card marker, or release-scoped defect listing exists in the V2 package.
- [x] No audit-register overlay, finding list, or register summary exists in the V2 package.
- [x] The filter conveys no severity or blocking ordering; rows are not sorted or styled to imply that a type is inherently more urgent.
- [x] Opening an Issue shows its summary body, origin (kind, Check when present, actor, time), linked Tasks, linked Checks, guardrail IDs, and resolution when present.
- [x] A resolved Issue renders its disposition distinctly for `verified`, `accepted`, and `duplicate`, and `accepted` and `duplicate` are not presented as proof of repair.
- [x] A duplicate Issue shows its canonical target and can navigate to it.
- [x] An Issue's `history` renders in recorded order with each entry's time, actor, kind, note, and Check; a reopened Issue shows its full prior history under the same ID.
- [x] Linked Tasks and Checks come from `index.TaskIssues`, `index.CheckIssues`, and the Issue's own reference lists; no relation is rebuilt by walking records.
- [x] The overlay can be opened scoped to the focused Task, showing that Task's linked Issues, and returns to the board with focus restored.
- [x] A project with no Issues shows an explicit empty state, not an error.
- [x] The overlay scrolls at any viewport height and wraps nothing.
- [x] Opening, filtering, navigating, and closing the overlay writes nothing, asserted by a byte-and-mtime snapshot around the sequence.
- [x] `make build && make test` passes.

## Implementation Plan

- [x] Add Issues overlay state: the filter, the list cursor, the open Issue, the detail scroll offset, and the origin surface.
- [x] Build the ordered Issue list and its per-Issue link sets in the load command.
- [x] Add list rendering with type, status, and severity columns, and the filter indicator.
- [x] Add filter key handling cycling the five types and the unfiltered view.
- [x] Add Issue detail rendering: origin, links, guardrail IDs, resolution by disposition, and history in recorded order.
- [x] Add duplicate-target navigation and return.
- [x] Add the Task-scoped entry point from a focused card, with focus restoration on close.
- [x] Add tests: ordering, each filter, empty state, each disposition, reopened history, duplicate navigation, Task-scoped open, focus restoration, and the no-write snapshot.
- [x] Run `go test ./internal/board/... ./internal/data/...`, then `make build && make test`.

## Context Log

Implemented the read-only V2 Issues surface in the sibling board package. The
load command now prepares stable I### rows and Task-scoped link sets from the
V2 index; the reducer handles filter/list/detail state and duplicate navigation;
rendering shows the recorded Issue data without doing IO or deriving severity.

Named evidence:

- `TestIssuesListUsesStableIdentityOrderAndShowsRecordedSeverity`
- `TestIssuesFilterCyclesThroughEveryTypeAndBackToAll`
- `TestIssueDetailShowsOriginLinksGuardrailsResolutionAndHistory`
- `TestTaskScopedIssuesUseIndexedLinksAndRestoreBoardFocus`
- `TestEmptyIssuesOverlayIsExplicitAndScrollableDetailDoesNotWrap`
- `TestIssuesNavigationAndFilteringAreReadOnly`
- Existing V2 boundary tests continue to prove no separate defect or audit
  surface is reachable from the package.

Files read: `.savepoint/router.md`, the E49 epic detail, this task, and
`.savepoint/Guardrails.md`; the task Context Files listed in the task were
reviewed before implementation.

Files edited: `internal/board/v2/load.go`, `internal/board/v2/model.go`,
`internal/board/v2/update.go`, `internal/board/v2/view.go`, new
`internal/board/v2/issues.go`, new `internal/board/v2/issues_view.go`, and
new `internal/board/v2/issues_test.go`.

Quality gates:

- `go test ./internal/board/... ./internal/data/...` — passed.
- `make build && make test` — passed for all packages.
- `.savepoint/Health-Check.md` is absent, so no Quick evidence block applies.
