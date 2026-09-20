---
id: E49-objective-board/T003-show-the-work-not-the-lifecycle
title: Show the work, not the lifecycle
status: done
objective: Render three columns of Task cards labelled by their human title, with implementation, clearance, owner, dependency, replan, and exception state carried by badges from one mapping.
depends_on:
    - E49-objective-board/T002-open-the-board-a-v2-project-already-has
complexity_tier: high
complexity_reason: The epic's core rendering surface and the one vocabulary every later surface reuses.
---

# T003: Show the work, not the lifecycle

## Problem

This is where the epic's title promise becomes visible. The columns stay Planned, In Progress, Done — the release design is explicit that implementation, check, owner-wait, and replan distinctions become badges rather than new columns, because five columns would turn three lifecycle states into five competing completion concepts. The card's label is `TaskV2.Title`, the sentence a person wrote about the outcome. The `objective` field is an `O###` identity reference and is never display language; V1's title-falls-back-to-objective behavior has no counterpart here and must not be recreated.

The badges carry the part a user actually needs before deciding what to open: is this Task being built, tested, or ready for a Check; does it have a Check at all; is that Check CLEAR, NEEDS WORK, stale, or unrecorded; is it waiting on the owner; is it blocked by a dependency; does it carry a replan flag; was it closed only by a recorded exception. Two of those have specific presentation obligations from the release design: done-by-exception must be visibly distinct from an ordinary done, and a stale completion must read as needing attention rather than as finished.

Every one of those states is a typed value that already exists — `ProgressStage`, `ClearanceState`, `GateBlockKind`, `FreshnessState`, and the presence of `Evidence.Exception` and `Evidence.Replan`. The board picks a glyph and an accent for `ClearanceStale`; it does not decide that clearance is stale. Concretely: the card asks `ResolveClearance` and the relevant one of `ResolveTaskStart` / `ResolveTaskAdvance` / `ResolveTaskCompletion` for the focused state, and renders what comes back. No comparison of Check IDs, no reading of `freshness.state` to draw a conclusion, no re-derivation of whether a dependency is satisfied. If a badge needs a fact no resolver returns, that is a signal to stop and raise it, not to compute it in a render path.

The mapping from typed value to glyph, label, and accent lives in one file. That is what keeps a second copy of the vocabulary from appearing in the sidebar, the detail overlay, and the plain renderer as the epic proceeds — each of those reuses this mapping rather than writing its own (STYLE-07, STYLE-09). The badge introduces no word of its own into the lifecycle vocabulary `internal/data` owns (DATA-02).

Colour alone cannot carry any of it. The palette is Atari-Noir's existing accents in `internal/styles`; every badge must remain distinguishable by glyph and text when colour is absent, which the visual identity requires and the monochrome and non-TTY paths later in this epic depend on.

## Context Files

- `internal/board/v2/model.go`
- `internal/board/v2/view.go`
- `internal/board/card.go`
- `internal/board/column.go`
- `internal/board/layout.go`
- `internal/board/status.go`
- `internal/board/util.go`
- `internal/styles/styles.go`
- `internal/styles/palette.go`
- `internal/styles/styles_test.go`
- `internal/data/task_v2.go`
- `internal/data/gate_v2.go`
- `internal/data/evidence_v2.go`
- `internal/data/lifecycle.go`
- `.savepoint/visual-identity.md`
- `agent-skills/bubbletea-tui-design/SKILL.md`

## Acceptance Criteria

- [x] Three columns render — Planned, In Progress, Done — grouped by `TaskV2.Status`, with no fourth column and no lifecycle flag promoted to a column.
- [x] Each card's primary label is `TaskV2.Title`; the `O###` owner reference is never used as display language, and no title fallback exists in the V2 package.
- [x] Each card shows its `T###` identity alongside the title.
- [x] A single badge mapping file translates `ProgressStage`, `ClearanceState`, `GateBlockKind`, `FreshnessState`, exception presence, and replan presence into glyph, label, and style; no other file in the package restates any of those values.
- [x] Every clearance state renders distinctly: missing, needs_work, stale, unknown, and current.
- [x] Stage renders distinctly for build, test, and audit on an in-progress Task, and is absent on planned and done Tasks.
- [x] An owner wait, a Task dependency wait, an Objective dependency wait, and a replan flag each render as their own badge, read from the gate decision rather than from the Task's raw fields.
- [x] A Task done by recorded exception is visibly distinct from an ordinary done Task.
- [x] A done Task whose clearance is stale reads as needing attention rather than as finished.
- [x] Every badge remains distinguishable with colour disabled, by glyph and text alone.
- [x] No card rendering path compares Check IDs, inspects `freshness.state`, or evaluates a dependency; every readiness statement is a value returned by `ResolveClearance`, `ResolveTaskStart`, `ResolveTaskAdvance`, or `ResolveTaskCompletion`, asserted by a test that stubs or fixtures the resolver inputs and observes the rendered result.
- [x] Focus changes colour and glyphs only; column width, card width, padding, and border thickness are identical focused and unfocused.
- [x] A column with more cards than fit scrolls, with a scroll indicator, and no content wraps.
- [x] New style tokens live in `internal/styles` and use the existing palette; no new colour value is introduced outside it.
- [x] `make build && make test` passes.

## Implementation Plan

- [x] Enumerate the badge states from the typed `data` values and write the mapping table first, in its own file, before any rendering code.
- [x] Add the style tokens the mapping needs to `internal/styles`, reusing existing accents.
- [x] Add card rendering that takes a Task plus its resolved decisions and returns a line, with no resolver calls inside the renderer itself.
- [x] Add column rendering, grouping, scrolling, and the scroll indicator, reusing the geometry approach in the V1 `column.go` rather than inventing a second one.
- [x] Wire columns into the V2 view alongside the space reserved for the sidebar and Next area.
- [x] Add golden-style render tests over fixture projects covering every badge state, including the exception and stale-done cases.
- [x] Add the focus-geometry test asserting identical widths focused and unfocused.
- [x] Add the colour-disabled test asserting every badge is still distinguishable.
- [x] Add the test proving no readiness is computed in the render path.
- [x] Run `go test ./internal/board/... ./internal/styles/...`, then `make build && make test`.

## Context Log

### Files read

`.savepoint/router.md`, `E49-Detail.md`, this task, `.savepoint/Guardrails.md`, and every file in `## Context Files`. Also read for targeted verification: `T004`'s Problem and Acceptance Criteria (to leave the sidebar, Objective filtering, and sidebar-to-column focus where they belong), `internal/data/dependency.go` and `internal/data/objective_gate_v2.go` for the typed dependency blocks a wait badge names, `internal/data/check_v2.go` for the Check frontmatter the fixtures write, and `internal/board/detail.go` for the V1 stage wording.

### Files added

- `internal/board/v2/badges.go` — the one mapping: `Badge`, the glyph constants, and the translation of stage, clearance, gate blockers, completion, and exception into glyph, label, and accent.
- `internal/board/v2/card.go` — `TaskCard` (a Task plus its resolved clearance and decision), `newTaskCard`/`groupTaskCards` resolution, and the pure `renderCard`.
- `internal/board/v2/column.go` — column frame, header, focus-anchored window, scroll indicators, and the shared geometry helpers.
- `internal/board/v2/update.go` — the reducer, moved here from `model.go`, plus column navigation.
- `internal/board/v2/{badges,card,column,columns_view}_test.go`.

### Files edited

`internal/board/v2/model.go` (cards and focus state; reducer moved out), `internal/board/v2/view.go` (height budgeting, column layout, reserved sidebar space), `internal/board/v2/boundary_test.go` (the single-mapping and derive-nothing structural assertions, and comment-free source scanning), `internal/board/v2/fixture_test.go`, `internal/styles/styles.go`, `internal/styles/styles_test.go`, `AGENTS.md`.

### Decisions worth recording

- **A card's decision is chosen exactly as `data.ResolveNext` chooses one** — start for planned, advance for build/test, completion at audit — so a card and the Next area can never report two judgements about the same Task. A test pins the card's decision to the projection's for the selected Task.
- **A done Task carries no decision.** `ResolveTaskCompletion` over a done Task returns `invalid_state`, which says nothing a reader needs; a closed Task's state is its clearance and how it closed. Its exception therefore comes from the recorded `Evidence.Exception`, and an open Task's from `GateDecision.AllowedByException` — one `ByException` field, resolved before rendering.
- **Clearance blockers get no badge of their own.** `clearanceBadge` already states missing, needs_work, stale, and unknown from the same resolved value; `blockerBadge` reports them as not-stated so the card never prints one fact twice in two wordings. `checker_authority` does get its own badge — it is a provenance fact the clearance state does not carry.
- **Four accents, no new colour.** `BadgeClear`, `BadgeAttention`, `BadgeWaiting`, and `BadgeNeutral` reuse the palette's green, orange, purple, and dim. Every badge also carries a glyph and a label, which is what the monochrome tests assert against.
- **The card frame is bordered in both focus states.** V1's `Card`/`CardFocused` pair adds a border on focus, which moves the layout; `CardBox`/`CardBoxFocused` differ in border colour alone, and a styles test asserts their border and padding match.
- **Column navigation landed here** rather than in T004: a focus that renders but cannot move leaves the focus-geometry and scrolling criteria unverifiable in use. T004 adds the sidebar surface above these keys in the same reducer, as the epic's file map lays out.
- **The sidebar's space is held open at ≥120 columns** so column geometry does not move when T004 fills it.
- **Structural assertions scan code with comments stripped** (`go/parser` + `go/printer`), after a comment mentioning `data.ResolveNext` tripped the derive-nothing scan. The assertion is about what the package does, not about a sentence explaining it.
- **Task status and stage are deliberately outside the single-mapping scan.** `card.go` reads them to pick which gate decision governs a Task — the same dispatch `data.ResolveNext` performs — which is not a second copy of the badge vocabulary.
- Not touched, and owned elsewhere: non-TTY output still reports counts (T010), and the header's `1 objectives` grammar is T002's wording, left as shipped.

### Evidence

- `go test ./internal/board/... ./internal/styles/... ./cmd/... .` — pass. Named cases:
  - `badges_test.go`: `TestClearanceBadgeRendersEveryStateDistinctly`, `TestStageBadgeIsPresentOnlyWhileInProgress`, `TestBlockerBadgeCoversEveryKind`, `TestBlockerBadgeNamesTheWaitTarget`, `TestCompletionBadgeDistinguishesExceptionAndStaleFromFinished`, `TestBadgesStayDistinctWithColorDisabled`, `TestBadgesCarryGlyphAndLabelUnderColor`.
  - `card_test.go`: `TestRenderCardLabelsWithTheTitleAndCarriesTheIdentity`, `TestRenderCardReadsOnlyResolvedValues` (the derive-nothing proof: a stale clearance naming no Check and a wait on a Task in no project), `TestRenderCardOmitsBlockersTheClearanceBadgeAlreadyStates`, `TestRenderCardDistinguishesDoneByException`, `TestRenderCardStageAbsentOffAnInProgressTask`, `TestRenderCardFocusChangesColorNotGeometry`, `TestRenderCardNeverExceedsItsWidth`, `TestGroupTaskCardsGroupsByRecordedStatus`, `TestGroupTaskCardsResolvesTheSameDecisionTheProjectionDoes`, `TestCardReportsAnObjectiveLevelWait`, `TestGroupTaskCardsHandlesAProjectWithNoIndex`.
  - `column_test.go`: `TestRenderColumnHeadsWithItsLabelAndCount`, `TestRenderColumnEmptySaysSo`, `TestRenderColumnScrollsWithIndicators`, `TestRenderColumnKeepsItsBudget`, `TestRenderColumnNeverExceedsItsWidth`, `TestRenderColumnFocusChangesColorNotGeometry`, `TestVisibleWindowKeepsTheFocusedCardVisible`.
  - `columns_view_test.go`: `TestBoardShowsTasksInTheColumnTheirStatusNames` (every badge state over one fixture project), `TestBoardDrawsExactlyThreeColumns`, `TestBoardLinesNeverExceedTheTerminalWidth`, `TestBoardEmptyProjectStillDrawsThreeEmptyColumns`, `TestNavigationMovesFocusAndClamps`, `TestNavigationOverAnEmptyColumnHoldsAValidCursor`, `TestReloadClampsFocusIntoTheCardsThatRemain`.
  - `boundary_test.go`: `TestBadgeVocabularyLivesInOneFile`, `TestRenderingResolvesNothing`, plus T002's `TestPackageDoesNotImportTheV1Board`, `TestPackageReferencesNoV1RecordType`, and `TestUpdateDoesNotPerformIO` (now following the reducer to `update.go`).
  - `internal/styles/styles_test.go`: `TestBadgeStyles_usePaletteAccents`, `TestCardBoxStyles_differOnlyInAccent`.
  - Every `internal/board` (V1) test unchanged and passing; no V1 board file was touched by this task.
- `make build && make test` — pass.
- Rendered and read at 80×24, 100×20, and 120×40 with `NO_COLOR=1`: three aligned columns of equal height, titles wrapping inside their cards, badges packed onto their own lines, scroll indicators where cards fall below the fold, and no line wider than the terminal.

`.savepoint/Health-Check.md` is absent, so no Quick health-check evidence block applies.
