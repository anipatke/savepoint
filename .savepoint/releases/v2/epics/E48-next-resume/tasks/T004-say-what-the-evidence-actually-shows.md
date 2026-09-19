---
id: E48-next-resume/T004-say-what-the-evidence-actually-shows
title: Say what the evidence actually shows
status: done
objective: Render one resolved projection as readable narrative that distinguishes missing, unknown, and stale evidence and claims no verification it did not perform.
depends_on:
    - E48-next-resume/T003-decide-the-one-next-action
complexity_tier: medium
complexity_reason: A new presentation package whose correctness is mostly wording precision and deterministic non-TTY output.
---

# T004: Say what the evidence actually shows

## Problem

The projection knows the truth; this task is where the truth survives being written down. Three states are routinely collapsed into one word by tools that should know better. `missing` means no Check has ever been recorded for this target. `unknown` means a Check exists and cleared, but nobody assessed whether that clearance still holds. `stale` means someone assessed it and concluded it does not. These are three different things a user has to do next, and rendering all of them as "unverified" — or worse, as a red mark that reads like failure — destroys the distinction E43 built the evidence model to preserve.

The second wording rule is that resume never takes credit. It read records; it ran nothing. "Cleared by C012 on 2026-08-14, freshness assessed current by checker-session-7" is a report of what is written down. "Verified" is not, and neither is any phrasing that implies the command confirmed anything just now. This matters because resume is the surface a user reaches for when they have lost the thread and are most likely to trust whatever it says.

The renderer is also where the epic's shared-interpretation claim is kept structurally true. It takes the resolved projection as a value and renders it. It does not take a project root, it does not take an index to consult on the side, and it does not re-resolve anything the projection already decided. If a fact is not on the projection, the fix is to put it on the projection — where the board gets it too — not to look it up here.

Output must be complete and deterministic without a TTY, since this command will be read out of pipes, logs, and agent transcripts at least as often as it is read on a terminal.

## Context Files

- `internal/data/next.go`
- `internal/data/gate_v2.go`
- `internal/data/evidence_v2.go`
- `internal/data/issue_v2.go`
- `internal/board/model.go`
- `internal/styles/styles.go`
- `AGENTS.md`
- `.savepoint/releases/v2/epics/E48-next-resume/E48-Detail.md`

## Acceptance Criteria

- [x] A new `internal/resume` package renders a resolved `Next` value and its carried decisions to an `io.Writer`, taking no project root and performing no discovery.
- [x] Rendering covers: the selected Objective and Task with their titles, implementation state, technical clearance state, owner-wait state, recorded evidence, relevant Issues, and the next action.
- [x] `missing`, `unknown`, `stale`, `needs_work`, and `current` clearance each render in distinct wording, asserted by test on the rendered text.
- [x] Clearance rendering names the Check ID when one exists, and renders the freshness assessor, date, and recorded basis when a freshness assessment exists.
- [x] No rendered string asserts that resume verified, ran, checked, or confirmed anything; a test asserts the absence of those claims across every rung's output.
- [x] An owner wait renders as an owner wait, distinct from a technical block, and names which Check the owner has yet to accept when one is recorded.
- [x] A completion allowed by recorded exception renders as an exception with its reason and owner, never as clearance.
- [x] Issues relevant to the selected Task or Objective render with their ID, type, and status; a project with none renders no Issues section rather than an empty one.
- [x] A selection diagnostic from T002 renders the diagnostic and the available next action together.
- [x] Output is byte-identical across runs for the same projection, contains no ANSI escapes when color is disabled, and remains readable at 80 and at 40 columns.
- [x] The evidence and freshness phrasing lives in one file with one source per phrase; no phrase is duplicated across layout code (STYLE-07, STYLE-09).
- [x] The package performs no filesystem access, no subprocess, and no network access.
- [x] `AGENTS.md` gains an `internal/resume/` Codebase Map row (ARCH-04).
- [x] `make build && make test` passes.

## Implementation Plan

- [x] Create `internal/resume/resume.go` taking the projection and a writer, and `internal/resume/evidence.go` holding the phrasing.
- [x] Render each rung kind from the projection's carried decision, with no branch that consults an index.
- [x] Write the five clearance phrasings and the owner-wait, exception, and Issue phrasings as named values, not inline literals.
- [x] Add golden-style tests per rung asserting exact rendered text.
- [x] Add the negative-wording test scanning every rung's output for verification claims.
- [x] Add the determinism test (two renders of the same projection are byte-identical) and the narrow-width tests.
- [x] Propagate writer errors rather than swallowing them, so T005 can surface them.
- [x] Add the `internal/resume/` Codebase Map row.
- [x] Run `go test ./internal/resume/... ./internal/data/...`, then `make build && make test`.

## Context Log

**Files read:** this task file, `E48-Detail.md`, `AGENTS.md` Codebase Map, completed `T002`/`T003` task files, `internal/data/next.go`, `internal/data/next_test.go`, `internal/data/gate_v2.go`, `internal/data/objective_gate_v2.go`, `internal/data/evidence_v2.go`, `internal/data/issue_v2.go`, `internal/data/task_v2.go`, `internal/data/objective_v2.go`, `internal/data/router_v2.go`, `internal/data/project.go` (`V2Index`'s `TaskIssues`/`CheckIssues`/`ScopeChecks` link maps), `internal/data/dependency.go`, `internal/data/check_v2.go`, `internal/data/lifecycle.go`, `internal/styles/styles.go` (confirmed no TUI styling belongs in this plain-text package), `go.mod` (module path).

**Files edited:**
- `internal/data/next.go` — added `Next.Issues []*IssueV2` and the unexported `relevantIssues` helper. The Problem statement's rule — "If a fact is not on the projection, the fix is to put it on the projection" — meant relevant Issues had to live on `Next` itself, since resume's contract (AC 1) is that it renders only the projection value, with no index of its own. `relevantIssues` reads only maps `LoadV2Index` already built (`TaskIssues`, `ObjectiveTasks`, `ScopeChecks`, `CheckIssues`); it derives no new relationship. A selected Task's Issues come straight from `TaskIssues[task.ID]`. A selected Objective (no Task) has no such direct map — `IssueV2` carries no Objective reference — so its Issues are the union of every owned Task's `TaskIssues` and any Issue linked to a Check scoped to the Objective itself, sorted and deduplicated. Returns nil (not empty) when nothing applies, so rendering can omit the section entirely.
- `internal/data/next_test.go` — three new tests: Issues linked to a selected Task, the Objective union case (owned-task issue + objective's-own-check issue), and the nil-not-empty case.
- `internal/resume/evidence.go` (new) — every phrase resume uses, held in one file per STYLE-07/09: `clearancePhrase` (five distinct states, naming the Check and, via `freshnessBasisPhrase`, the assessor/date/basis), `ownerWaitPhrase`, `exceptionPhrase`, `replanPhrase`, `taskDependencyPhrase`/`objectiveDependencyPhrase`, `issueLine`, `implementationPhrase`/`objectiveStatusPhrase`, and `selectionDiagnosticPhrase`. None of these words include "verified", "confirmed", or "checked" as a claim — `grep`-checked directly.
- `internal/resume/resume.go` (new) — `Render(w io.Writer, next data.Next) error`, building the full text via unexported `renderText` (so exact-text/determinism tests never need an `io.Writer`) and writing it in one `io.WriteString` call, propagating its error unchanged. Dispatches on `next.Kind` with no branch that reads an index; a completion allowed by exception is rendered before any clearance check so it can never fall through to "Technical clearance:" wording. Plain text only — no lipgloss/ANSI — which trivially satisfies "no ANSI escapes when color is disabled" and stays readable at any width since nothing is column-aligned.
- `internal/resume/resume_test.go` (new) — one golden exact-text test per reachable rung (migration, replan, dependency — both Task- and Objective-level —, execute both plain and exception-allowed, check-needed, the five-way clearance-state distinctness table, owner-validation-required, objective-integration both plain and owner-wait, ready both Task- and Objective-only, plan-objective, and the two selection-diagnostic variants), an Issues-section test and a no-Issues test, and four cross-cutting tests run over a shared `allRungFixtures()` set: no-verification-claims, determinism, no-ANSI, and narrow-width (no box-drawing/tab glyphs). Added a writer-error-propagation test and an import-boundary test (`go/build.ImportDir`) proving the package imports none of `os`, `os/exec`, `net`, `net/http`, `syscall`.
- `AGENTS.md` — extended the `internal/data/` row with the Issues-on-`Next` responsibility, and added the `internal/resume/` Codebase Map row.

**Quality gates:** `go build ./...` (clean), `gofmt -l` over every changed/new file (clean), `go vet ./internal/resume/... ./internal/data/...` (clean), `go test ./internal/resume/... ./internal/data/... -v` (all pass), `make build && make test` (pass, all packages). No `.savepoint/Health-Check.md` in this project, so the Quick check step is skipped per the build-task skill.

No drift: `internal/resume/resume.go` and `internal/resume/evidence.go` are exactly the two files `E48-Detail.md`'s Components and files table names for this task, plus the `Next.Issues` extension the task's own Problem statement required and the `AGENTS.md` rows it explicitly requires.
