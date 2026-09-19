---
id: E51-first-class-releases/T008-show-one-release-answer-on-board-and-resume
status: done
objective: Present the same Release promise, evidence, readiness, and next action across board, plain output, and resume.
depends_on:
    - E51-first-class-releases/T007-switch-releases-with-the-existing-r-key
complexity_tier: high
complexity_reason: Coordinates interactive detail, deterministic output, shared wording, width safety, and cross-command parity.
---

# T008: Show one Release answer on board and resume

## Problem

Release selection is useful only if every surface explains the same delivery promise and readiness. The board must show detail and evidence, while plain output and resume must use the exact same `data.Next` interpretation rather than independently querying the Release index.

## Context Files

- `.savepoint/releases/v2/epics/E51-first-class-releases/E51-Detail.md`
- `internal/data/next.go`
- `internal/data/release_gate_v2.go`
- `internal/board/v2/detail.go`
- `internal/board/v2/detail_test.go`
- `internal/board/v2/detail_view.go`
- `internal/board/v2/checks.go`
- `internal/board/v2/next_panel.go`
- `internal/board/v2/next_panel_test.go`
- `internal/board/v2/plain.go`
- `internal/board/v2/run_test.go`
- `internal/board/v2/width.go`
- `internal/board/v2/width_test.go`
- `internal/resume/evidence.go`
- `internal/resume/resume.go`
- `internal/resume/resume_test.go`
- `main_board_next_parity_test.go`
- `main_resume_matrix_test.go`

## Acceptance Criteria

- [ ] Release detail shows title/outcome, member Objective progress, latest and superseded Release Checks, relevant Issues, freshness basis, exceptions, and owner acceptance.
- [ ] Current CLEAR, owner-wait, stale/unknown/NEEDS WORK, done-by-exception, and historical-completion states are textually distinct without relying on colour.
- [ ] Board Next, non-TTY board output, and `savepoint resume` report the same selected Release, readiness state, evidence basis, owner wait, and action from one `data.Next`.
- [ ] Renderers do not consult `V2Index` to reinterpret Release readiness after receiving the projection.
- [ ] A project with no Releases retains the existing output shape and wording except where an additive empty context is explicitly required.
- [ ] Missing, archived, and Objective-mismatched Release selections use the shared selection diagnostic and still expose an available global action.
- [ ] Release detail and selector-related status render without wrapping or geometry changes at narrow widths and with CJK, emoji, and combining marks.
- [ ] Non-TTY output is deterministic and ANSI-free; browsing Release detail, Checks, Issues, and filters changes no project bytes or mtimes.
- [ ] Parity coverage includes selected Release execution, Release Check needed, Release owner validation, Release ready, missing selection, and no Release.
- [ ] Existing Task/Objective parity rungs and V1 plain-board behavior remain unchanged.

## Implementation Plan

- [x] Add Release detail resolution and rendering over indexed identity plus canonical gate/evidence values.
- [x] Extend the Next panel and plain output with compact Release context from the projection.
- [x] Add shared Release wording to `internal/resume` and consume it from the board where the same fact is stated.
- [x] Keep historical proof, current technical clearance, and owner acceptance visually and textually separate.
- [x] Apply existing cell-width, truncation, monochrome, and stable-geometry helpers to every new line.
- [x] Extend binary-level board/resume parity fixtures across Release rungs and diagnostics.
- [x] Add no-write browse snapshots and deterministic repeated-output assertions.
- [x] Run focused board/resume tests and the unchanged no-Release matrix.

## Context Log

- Read the E51 detail, this task's context files, `.savepoint/Guardrails.md`, and the active task-building skill before implementation.
- Updated `internal/board/v2` Release detail/selector/Next rendering, `internal/resume` Release identity/evidence/action wording, and board/resume parity fixtures. Existing no-Release and V1 paths remain covered by the existing matrix.
- Added coverage for Release detail evidence/history/issues/owner boundary, read-only selector browsing, Release rung and selection-diagnostic resume cases, and board/plain/TUI/built-command parity.
- Quality gates passed: `go test ./internal/resume ./internal/board/v2 .`, `make build && make test`, and `git diff --check`.
- No `.savepoint/Health-Check.md` is present, so the repository health-check step was skipped per the phase workflow.
