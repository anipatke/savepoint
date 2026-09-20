---
id: E49-objective-board/T006-show-the-evidence-behind-the-badge
title: Show the evidence behind the badge
status: done
objective: Open a Task or Objective to its Outcome, its recorded evidence and freshness basis, and the full Check history for that target in recorded order.
depends_on:
    - E49-objective-board/T004-navigate-a-project-by-its-objectives
complexity_tier: medium
complexity_reason: Two detail surfaces over already-indexed records, with immutability and supersession to render honestly.
---

# T006: Show the evidence behind the badge

## Problem

A badge compresses a Check into a glyph. The moment a user disagrees with it — or is about to act on it — they need the record underneath: which Check, recorded by whom, in what session, covering what scope, and on what basis anyone called it current. This is the surface where V2's trust boundary becomes inspectable instead of asserted.

Checks are immutable and a rerun gets a new ID with `supersedes` pointing at the one it replaces. `index.ScopeChecks[target]` already holds them in recorded order and `index.LatestCheck[target]` names the one that counts. Rendering only the latest would hide the thing history is for: seeing that a target was CLEAR, then NEEDS WORK, then CLEAR again tells a user something no single record does. So the detail shows the chain, marks the latest, and marks a superseded entry as superseded rather than dropping it.

Freshness is the part most easily rendered dishonestly. `Freshness` carries who assessed it, when, and the recorded basis — and the release design is explicit that the correct phrasing is "recorded current as of …", not a claim that the current tree has been verified. A Check with no freshness assessment is `unknown`, which is a different statement from `stale` and from `missing`, and all three are different from a CLEAR Check whose assessment lacks independent checker provenance. Each gets its own words. The board reads these states from `ResolveClearance` and the evidence block; it does not compare Check IDs to work out which one applies.

The Task detail also carries what a person needs before picking the Task up: its Outcome and body content, its `depends_on` entries with each dependency's resolved state, a recorded replan reason when one is present, whether owner validation is required and which Check the owner accepted, and any recorded exception with its requirement IDs, reason, owner provenance, and time. The Objective detail carries the same evidence shape plus its owned Tasks and its Objective dependencies.

The body of these records is author-owned markdown. It is displayed and scrolled; it is not parsed for meaning, and nothing in this surface writes.

## Context Files

- `internal/board/v2/model.go`
- `internal/board/v2/update.go`
- `internal/board/v2/view.go`
- `internal/board/detail.go`
- `internal/board/audit_detail.go`
- `internal/data/check_v2.go`
- `internal/data/evidence_v2.go`
- `internal/data/gate_v2.go`
- `internal/data/objective_gate_v2.go`
- `internal/data/task_v2.go`
- `internal/data/objective_v2.go`
- `internal/data/project.go`
- `internal/data/parser.go`
- `internal/styles/styles.go`

## Acceptance Criteria

- [x] A focused Task opens a detail overlay showing its ID, title, status, stage, owning Objective by ID and title, and its record body.
- [x] A focused Objective in the sidebar opens a detail overlay showing its ID, title, status, body, owned Tasks, and Objective dependencies with each dependency's resolved state.
- [x] Task dependencies render with the required level (`clear` or `accepted`) and each dependency's resolved satisfaction, taken from the gate decision rather than by inspecting the dependency Task's fields.
- [x] Check history for the open record renders every Check in `index.ScopeChecks` in recorded order, showing ID, result, `checked_by` role and session, and `checked_at`.
- [x] The latest Check is marked as latest and a superseded Check is marked as superseded; no Check is omitted.
- [x] A record with no Check renders an explicit "no Check recorded" statement, not an empty section.
- [x] Clearance renders as one of missing, needs_work, stale, unknown, or current, in distinct words, with the state taken from `ResolveClearance`.
- [x] A current clearance names the assessing actor, the assessment time, and the recorded basis, phrased as recorded rather than as verified by the board.
- [x] A CLEAR Check lacking independent checker provenance is described in its own words, distinct from both `unknown` and `current`.
- [x] A recorded replan renders its reason; a recorded exception renders its requirement IDs, reason, owner provenance, time, and the Check it applies to.
- [x] Owner validation renders whether it is required and, when accepted, which Check was accepted.
- [x] Issues linked to the open record render as a list of IDs, titles, types, and statuses, read from the index link maps.
- [x] The overlay scrolls, closes back to the surface it was opened from, and restores that surface's focus and cursor.
- [x] The record body is displayed without being parsed for meaning, and long or wide body content does not wrap the surrounding frame.
- [x] Opening, scrolling, and closing any detail writes nothing, asserted by a byte-and-mtime snapshot around the sequence.
- [x] `make build && make test` passes.

## Implementation Plan

- [x] Add detail overlay state to the V2 model: which record is open, where it was opened from, and the scroll offset.
- [x] Resolve the per-record decisions and Check chain in the load command so the view performs no resolution (ARCH-02).
- [x] Add the Check history renderer, marking latest and superseded.
- [x] Add the evidence block renderer with the five clearance statements and the checker-provenance case, reusing the phrasing decision recorded in T005.
- [x] Add the Task detail: identity, lifecycle, dependencies with resolved state, replan, exception, owner validation, linked Issues, and body.
- [x] Add the Objective detail: identity, owned Tasks, Objective dependencies with resolved state, the same evidence block, and body.
- [x] Add open/close/scroll key handling with return-to-origin focus restoration.
- [x] Add tests for each clearance state, the supersession chain, the no-Check case, dependency rendering, exception and replan rendering, focus restoration, and the no-write snapshot.
- [x] Run `go test ./internal/board/... ./internal/data/...`, then `make build && make test`.

## Context Log

### Decision: the detail reuses `internal/resume`'s phrases, and seven more of them are exported

T005 settled that the board calls `internal/resume` rather than restating its evidence wording, and the board's own `TestEvidenceWordingHasOneSource` enforces it. The detail needs the same facts about one record rather than about a projection, so `ClearancePhrase`, `DependencyPhrase`, `ObjectiveDependencyPhrase`, `ExceptionPhrase`, `ReplanPhrase`, `IssueLine`, and `ActorLabel` were exported from their unexported originals. No new phrase was written: every sentence the overlay shows already existed for `savepoint resume`.

Two of them changed while being exported.

`ClearancePhrase` gained the sixth statement this task requires. `ResolveClearance` reports two different facts as `unknown` — a CLEAR Check nobody has assessed, and a CLEAR Check whose current assessment carries no independent checker session — and the second one carries the assessment it distrusts, which is how the two are told apart. The old wording said "no freshness assessment has ever been recorded" for both, which the E48 audit had already recorded as a latent wording bug (`E48-Audit.md`, non-blocking observations). It is still unreachable from a project on disk, because `DecodeCheckV2` and `decodeFreshnessV2` both refuse the shapes that produce it, so it is proven over a constructed `data.Clearance` rather than a fixture.

`ObjectiveDependencyPhrase` was restated about the dependency rather than about whoever waits on it: "Objective O002 is not done yet." rather than "The owning Objective is waiting on Objective O002, which is not done yet." A Task rung waits on its owning Objective's dependencies and an Objective's own detail waits on its own, and the reason is the same fact in both — so the subject belongs to the caller's line and the reason lives in one place. `resume`'s own line now reads `Blocked: The owning Objective is waiting: Objective O005 is not done yet.` and its assertion was updated.

### Decision: resolution and rendering are two files, not one

The plan said to resolve in the load command so the view performs no resolution. Resolving every record's detail at load, including bodies, for a board that shows one at a time, is work no reader asked for; resolving on open and again on every reload gives the same guarantee — the view is handed a settled value — without it. So `detail.go` resolves a `RecordDetail` and is the only place the overlay reaches the index, `detail_view.go` renders one and reaches nothing, and `checks.go` resolves the Check chain. That split is what makes the constraint checkable rather than merely intended: `TestDetailRenderingResolvesNothing` greps `detail_view.go` for the index and the resolvers, and `TestDetailRenderingReadsOnlyTheResolvedValue` removes the loaded index from a model with an overlay open and asserts the body renders byte-identically.

`reopenDetail` is the other half: every load re-resolves an open overlay, so a Check recorded underneath it appears without the reader closing and reopening, and a record the load no longer holds closes it rather than leaving a copy older than the project.

### Decision: the overlay takes the columns' region, and `v` opens an Objective

The detail replaces the three columns rather than the whole screen, so the header, the Next answer, and the status bar stay where they were. `renderBoard` was split into `boardChrome` / `boardBodyHeight` / `renderBody` so `Update` and `View` size the overlay through the same functions — that is what lets a scroll key clamp against exactly the window the renderer will draw, rather than moving state the view then ignores.

`enter` opens the focused Task's detail, because `enter` has no other meaning on the columns. On the sidebar `enter` already selects an Objective, so `v` is the detail key there; it is bound on the columns too, so one key works on both surfaces.

### Files read

`AGENTS.md`, `.savepoint/router.md`, `.savepoint/Guardrails.md`, `agent-skills/savepoint-build-task/SKILL.md`, `E49-Detail.md`, this task file and T005's, `internal/board/v2/{model,update,view,card,column,objectives,badges,load,run,next_panel}.go`, `internal/board/v2/{boundary,fixture,view,columns_view}_test.go`, `internal/board/detail.go`, `internal/resume/{resume,evidence}.go`, `internal/data/{check_v2,evidence_v2,gate_v2,objective_gate_v2,task_v2,objective_v2,project,dependency,issue_v2}.go`, `internal/styles/styles.go`.

### Files edited

- `internal/board/v2/detail.go` (new) — `RecordDetail`, `RecordRef`, `DependencyEntry`, `ObjectiveDependencyEntry`, `newTaskDetail`, `newObjectiveDetail`, `reopenDetail`, `linkedIssues`.
- `internal/board/v2/checks.go` (new) — `CheckEntry` and `checkHistory`, marking latest from `index.LatestCheck` and superseded from the `supersedes` links the Checks themselves carry.
- `internal/board/v2/detail_view.go` (new) — `renderDetail`, the section renderers, the wrap with hanging indent, and the window and scroll-limit functions.
- `internal/board/v2/model.go` — `Detail`, `DetailOffset`, `DetailOrigin`, and the `detailOrigin` type.
- `internal/board/v2/update.go` — overlay-first key dispatch, `handleDetailKey`, `openDetail`, `detailUnderCursor`, `closeDetail`, `scrollDetail`, `clampDetailScroll`, `refreshDetail`; `clampObjectiveCursorToRows` split out of `clampObjectiveCursor` so a restore puts the sidebar cursor back where the reader left it rather than snapping it to the selection.
- `internal/board/v2/view.go` — `boardChrome`, `boardBodyHeight`, `renderBody`, `detailViewport`; overlay hints.
- `internal/resume/evidence.go`, `internal/resume/resume.go` — the seven exports and the two wording changes above.
- `internal/board/v2/{detail,fixture,boundary}_test.go`, `internal/resume/resume_test.go`.
- `AGENTS.md` — Codebase Map rows for `internal/board/v2/` and `internal/resume/`.

### Named evidence (TEST-01, TEST-02, TEST-06)

All in `internal/board/v2/detail_test.go` unless noted.

| Criterion | Test |
|---|---|
| Task identity, lifecycle, owning Objective, body | `TestTaskDetailNamesTheRecordAndItsLifecycle`, `TestTaskDetailReportsAPlannedTaskHasNoStage` |
| Objective identity, body, owned Tasks, Objective dependencies with state | `TestObjectiveDetailNamesItsTasksAndItsDependencies` |
| Dependency level and resolved satisfaction, from the resolver | `TestTaskDetailDependenciesCarryTheirLevelAndResolvedState` |
| Whole chain in recorded order, with provenance and timestamp | `TestCheckHistoryShowsTheWholeChainInRecordedOrder`, `TestObjectiveCheckHistoryShowsItsOwnIntegrationChain` |
| Latest and superseded marked, nothing omitted | same two, plus `TestReloadRefreshesAnOpenDetailAndClosesADeletedOne` |
| No Check recorded stated explicitly | `TestARecordWithNoCheckSaysSoRatherThanShowingNothing` |
| Six clearance statements, all distinct, incl. the provenance case | `TestDetailClearanceStatesReadDistinctly`; `TestRender_clearanceStatesAreDistinct` (`internal/resume/resume_test.go`) |
| Current names actor, time, basis, as recorded | `TestDetailClearanceStatesReadDistinctly`, `TestDetailClaimsNoVerification` |
| Replan, exception with every recorded fact | `TestTaskDetailRendersAReplanByItsRecordedReason`, `TestTaskDetailRendersAnExceptionWithEveryRecordedFact` |
| Owner validation required, accepted, and not required | `TestTaskDetailReportsOwnerValidationAndWhatWasAccepted` |
| Linked Issues by ID, type, status, title, listed once | `TestDetailListsTheIssuesLinkedToTheRecord` |
| Scrolls, clamps at both ends, indicator | `TestDetailScrollsAndClampsAtBothEnds` |
| Closes to origin with focus, cursor, and selection intact | `TestClosingTheDetailRestoresTheSurfaceItWasOpenedFrom`, `TestTheOpenOverlayHoldsTheKeys` |
| Body displayed, not parsed; frame never widened | `TestTheRecordBodyIsDisplayedAndNotParsed`, `TestLongAndWideBodyContentDoesNotWidenTheFrame` |
| Writes nothing, byte-and-mtime snapshot | `TestOpeningScrollingAndClosingADetailWritesNothing` |
| Rendering resolves nothing | `TestDetailRenderingReadsOnlyTheResolvedValue`, `TestDetailRenderingResolvesNothing` (`boundary_test.go`) |

Failure-path cases (TEST-02): an unsatisfied dependency, a record with no Check, a `NEEDS WORK` rerun superseding a CLEAR one, owner validation required but unaccepted, completion by exception, a reload that removed the open record, and the detail key pressed over a project with no records. Every fixture is a temporary directory built by `writeEvidenceProject`; the live Savepoint project is never read (TEST-04).

### Quality gates

`go test ./internal/board/... ./internal/data/... ./internal/resume/...` — pass. `make build && make test` — pass, whole module, 2026-09-19.

## Drift Notes

Two departures from `E49-Detail.md`'s component table, both recorded rather than silent.

**A third file.** The table names `internal/board/v2/detail.go, checks.go` for this surface; it shipped as three — `detail.go` resolving, `detail_view.go` rendering, `checks.go` for the Check chain. The epic allows refining identifiers during task breakdown, and the split is what turns "rendering resolves nothing" into a structural assertion rather than a convention (STYLE-01).

**`internal/resume`'s exported surface widened again.** T005 moved that package from "render the projection for resume" to "own the evidence vocabulary for every surface reporting a `data.Next`". This task widens it once more: `ClearancePhrase`, `DependencyPhrase`, `ObjectiveDependencyPhrase`, `ExceptionPhrase`, `ReplanPhrase`, `IssueLine`, and `ActorLabel` are used by a surface that reports one record's evidence and no projection at all. Nothing about the package's boundary changed — no IO, no state, no index, `TestPackage_noFilesystemNetworkOrSubprocessImports` still holds — and the AGENTS.md row was updated (ARCH-04). T005's open question stands and is now sharper: the vocabulary's home may belong in a package of its own rather than in `internal/resume`, which E50 can settle when it promotes `internal/board/v2`. Worth confirming at epic audit.
