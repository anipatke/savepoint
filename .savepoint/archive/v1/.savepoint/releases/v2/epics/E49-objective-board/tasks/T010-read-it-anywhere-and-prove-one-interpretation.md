---
id: E49-objective-board/T010-read-it-anywhere-and-prove-one-interpretation
title: Read it anywhere, and prove one interpretation
status: done
objective: Make the V2 board correct at narrow widths, with wide characters, without colour, and without a TTY, then prove board and resume answer identically across every rung.
depends_on:
    - E49-objective-board/T009-stay-current-without-losing-your-place
complexity_tier: high
complexity_reason: Cross-cutting rendering correctness plus the epic's closing proof matrix over real projects.
---

# T010: Read it anywhere, and prove one interpretation

## Problem

Two obligations close this epic, and both are about whether what was built survives contact with reality.

**The first is that the board has to be readable where it is actually used.** Terminal geometry is fixed-width cell math, and the visual identity treats accidental wrapping as a bug, not a cosmetic flaw. A Task title is author-written prose that can be long, contain CJK or emoji, or contain combining marks — so width has to be measured, not counted in bytes or runes. A narrow terminal has to degrade by dropping or collapsing surfaces in a defined order, not by overflowing them. Focus has to change colour and glyphs without changing geometry, everywhere, not just in the columns where T003 asserted it. And colour is an enhancement: with a forced monochrome profile every badge and every state distinction must still be legible, because the badge vocabulary is the board's primary encoding and the user may be on a 16-colour terminal or piping through something that strips ANSI.

Without a TTY the board does not start at all; it prints. That output leads with Next, because a user who pipes the board is asking what is going on, and the answer is the same projection the TUI shows. Then the three columns with titles and badges, then a one-line Issues summary by type. V1's plain output opens with a scan for `## Proposed Changes` across `releases/*/epics/*/E##-Audit.md`; that has no V2 analogue and is not carried over — open Issues are the V2 signal that a project has outstanding follow-up, and the summary already carries them. The output must be deterministic across runs and free of ANSI escapes.

**The second is the proof the epic's central claim is true.** E49 asserts the board and resume share one interpretation. Every earlier task asserted a piece of that; this task assembles the matrix over real projects on disk, driven through the built binary rather than through in-package fixtures, so what is proven is the product and not the test harness. For each rung a board user can reach, the same project must produce the same rung, the same selected records, and the same action from both surfaces. Where they disagree, the projection is right and the board is wrong — that is what "the board derives nothing" means, and the test should be written so a divergence points at the board.

The matrix also closes the two claims that are easiest to erode late: that browsing writes nothing, and that no V1 interpretation leaks into the V2 board. Both were asserted in pieces; both are properties of the finished package and are checked here against the whole of it.

## Context Files

- `internal/board/v2/view.go`
- `internal/board/v2/model.go`
- `internal/board/v2/update.go`
- `internal/board/v2/load.go`
- `internal/board/plain.go`
- `internal/board/plain_test.go`
- `internal/board/layout.go`
- `internal/board/theme.go`
- `internal/board/util.go`
- `internal/board/board.go`
- `internal/styles/styles.go`
- `internal/styles/palette.go`
- `internal/resume/resume.go`
- `internal/data/next.go`
- `main.go`
- `main_test.go`
- `internal/testutil/fs.go`
- `internal/testutil/fixture.go`
- `.savepoint/visual-identity.md`
- `.savepoint/releases/v2/v2-Design.md`

## Acceptance Criteria

- [x] Every surface — sidebar, columns, Next area, detail overlay, Issues overlay, help, status bar — renders without wrapping or overflow at a defined narrow width and at several widths above it.
- [x] Display width is measured by terminal cell width, not byte or rune count; titles containing CJK, emoji, and combining marks render with correct alignment and truncation.
- [x] Truncation is visible as truncation and never splits a character.
- [x] Below the narrow threshold, surfaces collapse in a defined, tested order rather than overflowing.
- [x] Borders stay stable and focus changes colour and glyphs only, with identical geometry focused and unfocused, on every surface.
- [x] Under a forced monochrome or 16-colour profile, every badge, clearance state, focus indicator, and status distinction remains legible without colour.
- [x] Without a TTY the V2 board prints instead of starting, leading with the Next action, then the three columns with titles and badges, then a one-line Issues summary by type.
- [x] Non-TTY output contains no ANSI escape sequences and is byte-identical across repeated runs over an unchanged project.
- [x] Non-TTY output carries no scan for V1 audit proposals and no release or epic vocabulary.
- [x] A V2 project that fails to load prints its diagnostic to stderr and exits nonzero without a TTY.
- [x] A matrix over real on-disk projects, driven through the built binary, shows the TUI's Next area, the non-TTY output, and `savepoint resume` reporting the same rung, the same selected Objective and Task, and the same action for each of: pending migration, replan, task dependency, objective dependency, execute, check needed, owner validation required, objective integration, ready, and plan objective.
- [x] The parity assertions are written so a divergence identifies the board as the incorrect surface.
- [x] A full browse of the finished board — every surface, every filter, every overlay, a reload — leaves the project byte-identical and mtime-identical.
- [x] No file in `internal/board/v2` references a V1 board type, `data.Task`, `data.Defect`, `data.ReleaseInfo`, `data.EpicInfo`, or `data.AuditRegisterSet`, and the V2 package contains no release, epic, defect-collection, or audit-register vocabulary.
- [x] Every existing `internal/board` V1 test passes unchanged, and a V1 project's board behavior is unaltered.
- [x] Every fixture is a temporary project; the live `.savepoint/` project is never used as a mutable fixture (TEST-04).
- [x] `make build && make test` passes.

## Implementation Plan

- [x] Audit every V2 render path for width handling and replace any byte or rune width with cell-width measurement.
- [x] Define the narrow-width thresholds and the collapse order, and implement the degradation.
- [x] Sweep every surface for focus-induced geometry changes and fix them.
- [x] Add the forced monochrome profile path and verify each badge distinction survives it.
- [x] Implement V2 plain output: Next first, columns with badges, Issues summary; confirm no V1 audit-proposal scan remains.
- [x] Add the non-TTY failure path with stderr diagnostic and nonzero exit.
- [x] Build the fixture projects for the ten rungs as real directories, reusing the shapes earlier tasks established.
- [x] Add the matrix test driving the built binary for board plain output and resume, comparing rung, selection, and action.
- [x] Add the whole-board browse snapshot and the package boundary assertion over the finished package.
- [x] Run the full V1 board suite to confirm it is unchanged, then `go test ./...`, then `make build && make test`.

## Context Log

Implemented and verified while retaining `status: in_progress` for owner review.

Files read/edited: `internal/board/v2/{view.go,model.go,update.go,load.go,card.go,column.go,detail_view.go,help.go,issues_view.go,next_panel.go,objectives.go,plain.go,run.go,watch.go}`, `internal/board/{board.go,plain.go,plain_test.go,layout.go,theme.go,util.go}`, `internal/styles/{styles.go,palette.go}`, `internal/resume/resume.go`, `internal/data/next.go`, `main.go`, `internal/testutil/{fs.go,fixture.go}`, `.savepoint/visual-identity.md`, `.savepoint/releases/v2/v2-Design.md`, and the focused V2/main test fixtures and boundary tests.

Implementation: added cell/grapheme-safe truncation and fitted non-wrapping chrome; defined sidebar and compact-column collapse thresholds; kept focused/unfocused geometry stable; completed V2 plain output with card titles, badge text, and deterministic all-Issues summary; routed pending migrations through the V2 board so its highest-priority Next rung matches resume; added temporary-project tests for every V2 surface, Unicode widths, compact layout, monochrome badge readability, full browse/reload no-write behavior, package boundaries, TUI parity, and subprocess board/resume parity.

Quality gates: `go test ./...` passed; `make build && make test` passed. `.savepoint/Health-Check.md` is absent, so no Quick health-check evidence block was applicable.
