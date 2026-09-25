---
id: E49-objective-board/T005-put-the-one-next-action-on-the-board
title: Put the one next action on the board
status: done
objective: Render a prominent Next area from the same data.ResolveNext call resume makes, and prove the two surfaces report the same answer.
depends_on:
    - E49-objective-board/T003-show-the-work-not-the-lifecycle
complexity_tier: medium
complexity_reason: Small rendering surface whose value is entirely in the parity obligation it must satisfy.
---

# T005: Put the one next action on the board

## Problem

Section 11 of the release design says board and resume use the same interpretation. E48 built that interpretation and shipped one consumer. Until a second one calls it, "shared" describes an intention rather than a property of the code, and the natural failure — a board that quietly re-derives what to do next from the records in front of it — is exactly the divergence E48 existed to prevent.

So the Next area is not a board feature that happens to agree with resume. It is one `data.ResolveNext` call over a `NextInput` assembled from the loaded index, the decoded V2 router, and `migrate.PendingOperation` — the same three values `main.go` assembles for `resume`. The board formats the returned `Next`; it does not consult the index to second-guess it, does not read the router's `next_action` prose to decide anything, and does not filter the answer by the Objective the sidebar happens to have selected. The projection answers for the project, and the user needs that answer most when they are looking at the wrong part of it.

The rendering obligations follow from what a `Next` carries. The rung determines the sentence. The selected Objective and Task are named by ID and title. A `SelectionDiagnostic` is shown alongside the available next action, never instead of it — a router line that rotted still leaves a project with work to do, and substituting a similarly numbered record for a missing one is the specific behavior this model forbids. `Next.Issues` are the follow-ups relevant to the selection, already resolved; the Next area summarises them and T007's overlay shows them. A pending migration outranks everything and says so.

Wording carries the same honesty constraint the resume renderer holds: missing, unknown, and stale clearance are three different statements, a completion allowed by a recorded exception is named as an exception rather than as a CLEAR result, and nothing on this panel claims the board verified anything it merely read. `internal/resume/evidence.go` already holds that phrasing for one surface. Whether the board reuses that package or restates the phrasing is a judgement to make while building: `internal/resume` is a rendering package with its own contract, and importing it into the board may be cleaner than a second copy of the wording, or may drag a narrative layout into a panel that needs a compact one. Decide it deliberately and record which and why — a second uncontrolled copy of the evidence vocabulary is the outcome to avoid either way (STYLE-07, STYLE-09).

## Context Files

- `internal/board/v2/model.go`
- `internal/board/v2/view.go`
- `internal/board/v2/load.go`
- `internal/resume/resume.go`
- `internal/resume/evidence.go`
- `internal/resume/resume_test.go`
- `internal/data/next.go`
- `internal/data/gate_v2.go`
- `internal/data/issue_v2.go`
- `internal/migrate/operation.go`
- `main.go`
- `internal/styles/styles.go`

## Acceptance Criteria

- [x] The Next area renders from exactly one `data.ResolveNext` call, whose `NextInput` is built from the loaded index, the decoded V2 router, and `migrate.PendingOperation`.
- [x] No code in the V2 package inspects the index, a Task, an Objective, or a Check to alter, override, or re-derive what the Next area says, asserted by a test that fixtures a project and observes that the rendered action tracks only the projection.
- [x] The Next area is unaffected by the sidebar's Objective selection: changing selection does not change the reported next action.
- [x] Every rung renders a distinct, readable statement: pending migration, replan, dependency, execute, check needed, owner validation required, objective integration, ready, and plan objective.
- [x] The selected Objective and Task are named by ID and by title when the rung names one.
- [x] A `SelectionDiagnostic` renders alongside the available next action, naming the record the router named and why it did not resolve; no similarly numbered record is shown in its place.
- [x] Missing, unknown, and stale clearance render in distinct words, naming the Check and the recorded freshness basis when one exists.
- [x] A completion allowed by a recorded exception is named as an exception, not as a CLEAR result or current clearance.
- [x] Nothing in the Next area asserts a verification the board performed; recorded evidence is described as recorded.
- [x] `Next.Issues` are summarised by count and type without the overlay being open.
- [x] A pending migration outranks every other rung and reports the operation.
- [x] A fresh V2 project with no Objectives renders the planning next action with no error.
- [x] For a shared set of fixture projects covering the execute, check-needed, owner-validation, dependency, replan, objective-integration, ready, and plan rungs, the board's Next area and `savepoint resume` report the same rung, the same selected Objective and Task, and the same action.
- [x] The decision on reusing versus restating the resume evidence phrasing is recorded in the Context Log with its reason.
- [x] `make build && make test` passes.

## Implementation Plan

- [x] Confirm the `NextInput` the load command builds matches `main.go`'s `runResume` field for field, and reuse the same assembly if one can be shared.
- [x] Decide reuse versus restatement of the evidence phrasing; record the decision before writing the panel.
- [x] Add the Next panel renderer taking one `data.Next` value and producing a compact block.
- [x] Add the per-rung sentences, the selection-diagnostic line, and the Issues summary.
- [x] Place the panel prominently in the layout, above or beside the columns, without changing the geometry T003 established.
- [x] Add per-rung render tests over fixture projects.
- [x] Add the parity test invoking both the board's projection path and the resume renderer over the same fixtures and comparing rung, selection, and action.
- [x] Add the test proving selection changes do not move the Next answer.
- [x] Run `go test ./internal/board/... ./internal/resume/... ./internal/data/...`, then `make build && make test`.

## Context Log

### Decision: the board reuses `internal/resume`'s evidence phrasing rather than restating it

`internal/resume` now exports `EvidenceLines`, `ActionPhrase`, and `SelectionPhrase`; `internal/board/v2/next_panel.go` calls them and supplies its own compact layout. The reason is that restatement was not a neutral alternative here. The board's own boundary tests already confine the clearance, blocker, and freshness vocabulary to `badges.go`, so a second set of prose phrases would have had to translate `ClearanceStale`, `GateBlockOwnerAcceptance`, and `Exception` a second time, in a second file, with no mechanism keeping the two wordings in step — the drift E48 built one projection to prevent, reintroduced one layer up (STYLE-07, STYLE-09). Importing the package costs nothing structurally: `internal/resume` performs no IO, holds no state, and imports only `internal/data`.

The narrative-layout risk the task named is avoided by importing at the phrase level rather than calling `resume.Render`. The board builds its own lines — a rung heading, identity lines, evidence, selection diagnostic, an Issues count, and the action — and lays them out indented under the rung. `internal/resume` keeps its own narrative shape, including the `Implementation:` and per-Issue lines the panel does not carry.

A consequence worth stating: `internal/resume`'s recorded purpose widened from "render the projection for resume" to "own the evidence vocabulary for every surface reporting a `data.Next`". The AGENTS.md Codebase Map row was updated to say so, and the Drift Note below records it.

### Decision: `NextInput` is not shared with `runResume`, and parity is proven instead

The two assemblies are field for field the same — index, decoded V2 router, `migrate.PendingOperation` — but they cannot be reduced to one call today. `internal/data` cannot import `internal/migrate` (that package imports `internal/data`), so the shared assembly would need a new package, which E49 does not scope. The obligation was met by proof rather than by construction: `main_board_next_parity_test.go` runs both surfaces over the same project and fails if they disagree.

### Files read

`AGENTS.md`, `.savepoint/router.md`, `.savepoint/Guardrails.md`, `agent-skills/savepoint-build-task/SKILL.md`, `E49-Detail.md`, this task file, `internal/board/v2/{model,view,load,plain,update,run,badges}.go`, `internal/board/v2/{fixture,view,boundary,run,objectives,columns_view}_test.go`, `internal/resume/{resume,evidence}.go`, `internal/resume/resume_test.go`, `internal/data/{next,issue_v2}.go`, `internal/styles/styles.go`, `main.go` (`runResume`), `main_board_test.go`, `main_resume_matrix_test.go`.

### Files edited

- `internal/resume/resume.go` — `EvidenceLines` split out of `rungLines` (which now appends the narrative's section blank); `nextActionPhrase` → `ActionPhrase`.
- `internal/resume/evidence.go` — `selectionDiagnosticPhrase` → `SelectionPhrase`; package doc records the shared-vocabulary role.
- `internal/board/v2/next_panel.go` (new) — `nextLines`, the rung labels, identity lines, `issuesSummary`, and the styled `renderNext`.
- `internal/board/v2/view.go` — rung labels, `renderNext`, `nextSummary`, and `selectionDiagnosticSummary` removed; `renderSelection` is now about the sidebar's filter alone, with the selection diagnostic reported in full by the Next area instead of twice in two wordings.
- `internal/board/v2/plain.go` — non-TTY output prints the same `nextLines`; the selected-Objective line is `Selected:` so it cannot be mistaken for the projection's own `Objective:` line, and the pending-migration guidance is `Recovery:`.
- `internal/board/v2/next_panel_test.go` (new), `boundary_test.go` (three structural assertions), `main_board_next_parity_test.go` (new).
- Existing assertions updated for the new wording: `internal/board/v2/{run,objectives,columns_view}_test.go`, `internal/board/dispatch_test.go`, `main_board_test.go`. Three of those also raised their fixture terminal height, because the Next area is four to six lines where it was one and the assertions require every card or Objective to be on screen at once.
- `AGENTS.md` — Codebase Map rows for `internal/resume/` and `internal/board/v2/`.

### Named evidence (TEST-01, TEST-02, TEST-06)

| Criterion | Test |
|---|---|
| One `ResolveNext` call, in the load command | `TestProjectionIsResolvedOnlyInTheLoadCommand` (`boundary_test.go`) |
| Nothing re-derives the answer | `TestNextPanelDerivesNothing`, `TestEvidenceWordingHasOneSource` (`boundary_test.go`), `TestNextAreaTracksOnlyTheProjection` (`next_panel_test.go`) |
| Unaffected by sidebar selection | `TestNextAreaIgnoresTheSidebarSelection` |
| Nine distinct rungs | `TestNextPanelRendersEveryRungDistinctly` |
| Records named by ID and title | `TestNextPanelNamesTheSelectedRecordsByIDAndTitle`, `TestNextAreaReportsTheRungTheLoadResolved` |
| Diagnostic beside the action, no substitute | `TestNextPanelSelectionDiagnosticAccompaniesTheAction`, `TestNextPanelSelectionMismatchNamesBothRecords`, `TestRouterNamingAMissingObjectiveOpensTheBoardAnyway` |
| Clearance states distinct, Check and basis named | `TestNextPanelClearanceStatesReadDistinctly` |
| Exception named as an exception | `TestNextPanelNamesAnExceptionAsAnException` |
| No verification claim | `TestNextPanelClaimsNoVerification` |
| Issues by count and type; none renders no line | `TestNextPanelSummarisesIssuesByCountAndType` |
| Pending migration outranks and names the operation | `TestNextPanelPendingMigrationOutranksAndNamesTheOperation` |
| Fresh project plans, with no error | `TestNextAreaOpensAFreshProjectWithNoError`, `TestViewEmptyTemplateProjectOpensWithThreeEmptyColumns` |
| Board Next and resume agree, every rung | `TestBoardNextAndResumeReportTheSameAnswer` (`main_board_next_parity_test.go`), over the twelve on-disk `resumeMatrixCases()` fixtures |

Failure-path cases (TEST-02): unresolved and mismatched router selections, a project held back by a pending migration, a project with no Objectives, and the four `invalidProjectCases()` refusals, which still draw no board. All fixtures are temporary directories; the live project is never read (TEST-04).

### Quality gates

`make build && make test` — pass, whole module, 2026-09-19.

## Drift Notes

`internal/resume` gained an exported responsibility E49-Detail's component table did not anticipate: it is now the evidence vocabulary for every surface that reports a `data.Next`, not only the renderer behind `savepoint resume`. `internal/board/v2` imports it, which is a package dependency the epic's architecture section did not name — it described `internal/styles` as the only package shared with the V1 board, and said nothing about sharing with `resume`.

Nothing about the boundary changed: `internal/resume` still performs no IO, still consults no project root or index, and `TestPackage_noFilesystemNetworkOrSubprocessImports` still holds. The AGENTS.md Codebase Map rows for `internal/resume/` and `internal/board/v2/` were updated to record it (ARCH-04). Worth confirming at epic audit that this is the intended home for the vocabulary rather than a shared package of its own, which E50 could revisit when it promotes `internal/board/v2`.
