---
id: E49-objective-board/T008-let-the-gates-decide-which-keys-exist
title: Let the gates decide which keys exist
status: done
objective: Offer a board write only where a gate decision grants owner authority, refuse everything else by naming the unmet requirement and whose authority it needs, and make both writes safe.
depends_on:
    - E49-objective-board/T005-put-the-one-next-action-on-the-board
    - E49-objective-board/T006-show-the-evidence-behind-the-badge
    - E49-objective-board/T007-one-place-for-every-follow-up
complexity_tier: high
complexity_reason: The epic's only write path, and the one place a rendering bug becomes a recorded project change.
---

# T008: Let the gates decide which keys exist

## Problem

V1's board lets any focused card be advanced or retreated with a keypress, and checks dependencies itself to decide whether to allow it. That is the rule-in-the-consumer pattern the whole release is built to remove: the board enforcing a policy it also defines means the policy has two homes, and section 1 of the release design names `internal/board/transitions.go` as evidence of exactly that.

In V2 the decision already exists and already says who may take it. `GateDecision.Actor` grants executor authority for start and advance, checker authority for a completion with no remaining blockers, and owner authority for a completion allowed by a recorded exception. The person sitting at the board is the owner. So the rule is simple and it is the epic's sharpest break from V1: **a key exists only where the decision names owner authority.** Starting a Task, advancing a stage, writing a Check, and resolving an Issue are displayed with their state and their blockers, and there is no key that performs them — not a key that refuses, and not a key that warns. The board reports whose session does that work.

The two writes the owner does own are both real and both needed. Recording acceptance of a current Check is unambiguously the owner's act: it is the one thing `GateBlockOwnerAcceptance` is waiting for, and `WriteTaskEvidenceV2` / `WriteObjectiveEvidenceV2` already exist to record it. Recording a selection — which Objective and Task the router points at — is navigation made durable, and `WriteRouterStateV2` from T001 is deliberately shaped so it can write only `objective` and `task`, never `state` and never `next_action`.

A refusal has to be worth reading. The blockers on a `GateDecision` are typed and carry detail — which dependency, which Check, which requirement — and a refusal message should name the unmet requirement and the authority that resolves it, not say the action is unavailable. This is content, so it lives in one mapping rather than scattered through key handlers (STYLE-09).

Both writes go through `tea.Cmd`s (ARCH-02) and both refuse first if `migrate.PendingOperation` reports an incomplete conversion — the same boundary the V1 board, `doctor`, and `upgrade-assets` already use, so the board does not grow a second migration policy. Both re-resolve the gate decision against freshly read records immediately before writing, because section 4 requires completion controls to resolve current evidence again before they write and because the board's in-memory index can be seconds stale. Both detect a file changed underneath them and refuse with a refresh-and-retry message rather than a partial overwrite. A write that would change nothing changes nothing, including mtime. After a successful write the board reloads through the single load command from T002, so the resolved `Next`, the badges, and the evidence all move together rather than one surface updating and the others lagging.

## Context Files

- `internal/board/v2/model.go`
- `internal/board/v2/update.go`
- `internal/board/v2/load.go`
- `internal/board/v2/view.go`
- `internal/board/io.go`
- `internal/board/transitions.go`
- `internal/board/help.go`
- `internal/data/gate_v2.go`
- `internal/data/objective_gate_v2.go`
- `internal/data/evidence_v2.go`
- `internal/data/write.go`
- `internal/data/errors.go`
- `internal/migrate/operation.go`
- `.savepoint/Guardrails.md`
- `.savepoint/releases/v2/v2-Design.md`

## Acceptance Criteria

- [x] An action key exists for a focused record only when the relevant gate decision names `ActorRoleOwner`; the set of available keys is derived from the decision, not from the record's status.
- [x] No key in the V2 package starts a Task, advances a stage, retreats a Task, writes a Check, or changes an Issue's status — not even one that refuses; the capability is absent.
- [x] A Task or Objective whose completion is allowed with checker authority is displayed as clear to close by a checker, with no board key offered.
- [x] Recording owner acceptance is offered exactly when clearance is current and `GateBlockOwnerAcceptance` is the unmet requirement, and it records the accepted Check through the existing evidence writers.
- [x] A completion allowed by a recorded exception offers the owner action and names the exception rather than presenting a CLEAR result.
- [x] Recording a selection writes only `objective` and `task` through `WriteRouterStateV2`, leaving `state`, `next_action`, unknown fields, and the body byte-identical.
- [x] A refused action names the unmet requirement from the decision's blockers and the authority that resolves it; refusal wording lives in one mapping file, not in key handlers.
- [x] Every blocker kind — replan, dependency, objective_dependency, clearance_missing, clearance_needs_work, clearance_stale, clearance_unknown, owner_acceptance_required, checker_authority, invalid_state — produces its own refusal statement.
- [x] Every write is a `tea.Cmd` returning a typed message; no filesystem access occurs in `Update` (ARCH-02).
- [x] A pending migration operation refuses both writes with the operation's recovery guidance and writes nothing.
- [x] Both writes re-resolve the gate decision against freshly read records immediately before writing, and a decision that no longer permits the action refuses without writing.
- [x] A record or router file changed underneath the board refuses with a refresh-and-retry message and leaves the file byte-identical, with no partial write (FS-01).
- [x] A write whose result equals the current content leaves bytes and mtime unchanged (FS-04).
- [x] A successful write triggers a reload through the single load command, after which the Next area, badges, and evidence all reflect the new state.
- [x] Repeated presses of any action key are idempotent and produce no second write.
- [x] The help overlay lists only the keys that exist, and states which actions belong to an executor or checker session.
- [x] A full navigation session that takes no action — sidebar, columns, details, Issues overlay, filters — leaves the whole project byte-identical and mtime-identical, asserted by a snapshot over the project tree.
- [x] `make build && make test` passes.

## Implementation Plan

- [x] Write the availability rule first: a function from a focused record and its gate decision to the set of offered actions, with owner authority as the only admitting condition.
- [x] Write the refusal mapping from `GateBlockKind` to a statement naming the requirement and the authority, in its own file.
- [x] Add the owner-acceptance command: migration guard, fresh re-read, re-resolve, evidence write, typed result.
- [x] Add the selection command: migration guard, fresh re-read, `WriteRouterStateV2`, typed result.
- [x] Add mtime-conflict handling for both, returning a refresh-and-retry message rather than retrying blind.
- [x] Chain a reload through the T002 load command on every successful write.
- [x] Add key handling and the help overlay entries, and confirm no removed V1 transition key survived into the V2 package.
- [x] Add tests: availability per decision, absence of executor and checker capabilities, each blocker's refusal wording, both writes' field preservation, pending-migration refusal, stale-decision refusal, mtime conflict, no-op idempotence, reload after write, repeated keypresses, and the browse-writes-nothing snapshot.
- [x] Run `go test ./internal/board/... ./internal/data/... ./internal/migrate/...`, then `make build && make test`.

## Context Log

Implemented the V2 owner-action boundary. Focused records now expose only
owner-authorized capabilities: `p` records the router selection, `a` records
acceptance of a current Check when the completion gate reports the owner wait,
and `x` completes a record only when the gate grants owner authority through a
recorded exception. Executor/checker transitions, Check creation, and Issue
changes have no V2 board key. Refusal copy is centralized and names the
resolver's requirement and the role whose session owns it.

Named evidence:

- `TestActionsOnlyExposeOwnerAuthorityAndSelection`
- `TestEveryGateBlockerHasDistinctRefusalWithAuthority`
- `TestOwnerAcceptanceCommandRereadsAndReloads`
- `TestOwnerActionKeyReturnsTypedCommandAndReloadsThroughLoad`
- `TestExceptionCompletionCommandWritesOnlyWhenOwnerAuthorityRemains`
- `TestSelectionCommandOnlyChangesRouterSelectionAndReloads`
- `TestWritesRefusePendingMigrationWithoutChangingFiles`
- `TestHelpListsOnlyFocusedOwnerCapabilities`
- Existing `TestUpdateDoesNotPerformIO`, `TestNavigationWritesNothingToTheProject`,
  `TestWriteTaskEvidenceV2_refusesStaleSourceWithoutOverwritingUserEdit`,
  `TestWriteRouterStateV2_refusesStaleMtimeAndLeavesFileUntouched`, and the
  V2 package boundary tests.

Files read: `.savepoint/router.md`, `.savepoint/Guardrails.md`, the E49 epic
detail, this task, `agent-skills/savepoint-build-task/SKILL.md`, the task
Context Files, the existing V2 board/action-adjacent files, and the V2 design
sections governing authority, writes, and the board.

Files edited:

- `internal/board/v2/actions.go` — owner capability derivation, selection
  target resolution, centralized refusal statements, and focused action text.
- `internal/board/v2/io.go` — typed Bubble Tea commands for acceptance,
  exception completion, and selection, with migration/freshness guards and
  reload results.
- `internal/board/v2/help.go` — focused capability/help rendering.
- `internal/board/v2/model.go`, `update.go`, `view.go` — help state, typed
  action dispatch, reload chaining, and action/refusal hints.
- `internal/board/v2/badges.go` — shared typed gate/clearance predicates used
  by the action layer.
- `internal/board/v2/actions_test.go` — focused action, refusal, write,
  migration, reload, idempotence, and help coverage.

Quality gates:

- `go test ./internal/board/v2` — passed.
- `go test ./internal/board/... ./internal/data/... ./internal/migrate/...` — passed.
- `make build && make test` — passed for all packages.
- `go vet ./...` — passed.
- Changed Go files are gofmt-clean.
- `.savepoint/Health-Check.md` is absent, so the Quick health-check evidence
  block does not apply.

The task remains `status: in_progress` for user review; only the user may mark
it `done`.
