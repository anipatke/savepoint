---
id: C-911
scope: {kind: objective, id: O-013}
result: NEEDS WORK
checked_by: {role: checker, session: o013-full-check-20260923}
executed_session: unrecorded
checked_at: '2026-09-23T09:45:00Z'
reviewed:
  base_commit: d9f90f4aff7b05a28b2d621991659b92a2f7a478
  head_commit: d9f90f4aff7b05a28b2d621991659b92a2f7a478
  files:
    - internal/board/v2/releases.go
    - internal/board/v2/releases_test.go
    - internal/board/v2/detail_view.go
    - internal/board/v2/help.go
    - internal/board/v2/io.go
    - internal/board/v2/plain.go
    - internal/board/v2/update.go
    - internal/board/v2/view.go
    - internal/board/v2/footer_test.go
    - internal/board/v2/run_test.go
    - internal/resume/evidence.go
    - internal/resume/resume.go
    - internal/resume/resume_test.go
    - main_resume_matrix_test.go
    - internal/data/release_v2_test.go
    - internal/init/agent_skills_test.go
    - README.md
    - AGENTS.md
    - .savepoint/Design.md
    - .savepoint/Guardrails.md
    - .savepoint/router.md
    - templates/project-v2/AGENTS.md
    - templates/project-v2/.savepoint/Design.md
    - templates/project-v2/.savepoint/Guardrails.md
    - templates/project-v2/.savepoint/router.md
    - agent-skills/savepoint-idea/SKILL.md
    - agent-skills/savepoint-design/SKILL.md
    - agent-skills/savepoint-task/SKILL.md
    - agent-skills/savepoint-check/SKILL.md
    - agent-skills/references/check-method.md
    - agent-skills/references/issue-capture.md
    - agent-skills/references/commands-and-procedures.md
    - .savepoint/objectives/O-013-rename-release-to-goals/Objective.md
    - .savepoint/objectives/O-013-rename-release-to-goals/tasks/T-006-preserve-existing-project-records.md
    - .savepoint/objectives/O-013-rename-release-to-goals/tasks/T-007-present-goals-in-the-v2-board.md
    - .savepoint/objectives/O-013-rename-release-to-goals/tasks/T-008-align-active-guidance-with-goals.md
  dependencies: []
issues: [I-035, I-036]
supersedes: null
---

# C-911: O-013 Full Objective Check

## Verdict

`NEEDS WORK`. The code change is sound:

- The board and resume use Goal wording.
- `g` opens the selector. `r` is an identical, hidden alias.
- The stored `R-###` / `release:` contract is unchanged, and no data,
  migration, doctor, or gate code changed.
- A fresh `make test-full` passes.

Two in-scope wording gaps remain:

- **I-035** — the shipped V2 scaffold (`templates/project-v2/.savepoint/`
  Design, Guardrails, router) and this repository's live `router.md` still
  teach "Release Check" and "optional Releases". Every new V2 project inherits
  that.
- **I-036** — when an open Goal detail's record is removed, the reload status
  reads "RELEASE R-### no longer exists; detail closed."

Reviewed state: uncommitted working tree on top of `d9f90f4`. All three owned
Tasks are `done` with owner Task-check waivers.

## Evidence Mode

Full Objective evidence covering T-006, T-007, and T-008, all waived. Scope
includes cross-Task integration (board and resume wording against the
guidance docs) and Design reconciliation. This session is fresh and built
none of the work.

## Frozen Scope Lock

1. Acceptance and policy:
   - O-013's seven Success Conditions and its Boundaries.
   - Done When for T-006, T-007, and T-008.
   - Guardrails DATA-02, DATA-03, TPL-01, TPL-02, TPL-04, ARCH-02, ARCH-03,
     ARCH-04, TEST-01, TEST-02, TEST-04, TEST-06, and TEST-08.
2. Changed behavior and entry points:
   - V2 board: selector, header, footer hints, help, detail overlay, selection
     write and refusal messages, reload status, and plain (non-TTY) output.
   - `internal/resume` evidence, action, and selection phrases.
   - Active docs, and V2 skills with their scaffold copies.
3. Relied-on orchestration: `handleKey` dispatch order (help → selector →
   quit → detail → issues → record actions → `g`/`r`); `loadCmd` →
   `applyLoad` reload diagnostics; `LoadV2Index` Release membership;
   `ResolveSelection`; the Release completion and cutover resolvers.
4. Matrix axes:
   - Surface: selector, header, hints, help, detail, status/error, plain,
     resume.
   - Key: `g`, `r`, and collision with action keys.
   - State: none selected, selected, removed selection, removed open detail,
     empty Goal, cross-Goal refusal.
   - Width: narrow and normal.
   - Storage: existing `R-###` load, router selection, completion, cutover,
     migration.
   - Docs: active repository guidance, V2 skills and references, V2 scaffold
     files, parity, archive untouched.
5. Materiality boundary: a finding must violate an O-013 or Task criterion,
   or a named Guardrail, through a supported board, resume, or scaffold path.
   Doctor wording is outside this Objective's in-scope list (see
   Observations). Record bodies authored with "Release" prose are data, not UI.

## Coverage Matrix

| # | Cell | Evidence | Result |
|---|---|---|---|
| 1 | Existing `R-###` records load; membership comes from `release:` | `TestDecodeReleaseV2_valid`, `TestLoadV2Index_releasesDeriveObjectiveMembership`; live-project probe loaded R-006 and filtered O-013 | Proven |
| 2 | Router `release:` selection resolves; cross-Goal selection refused | `TestLoadV2Index_routerReleaseSelectionResolvesExistingRecords` (new), `TestSelectionWriteRejectsCrossReleaseObjective` | Proven |
| 3 | No Goal-specific gate, parser, writer, or migration change | `git diff --stat -- internal/data internal/migrate internal/doctor cmd main.go`: only `release_v2_test.go` (+47 test lines) | Proven |
| 4 | Selector heading and empty-state wording | `SELECT GOAL`, "(no Goals in this project)"; `releases.go:88,93`; tests pass | Proven |
| 5 | Header `GOAL:` line; plain `Selected Goal:` | `view.go:270`, `plain.go:47`; live-project probe printed `Selected Goal: R-006` | Proven |
| 6 | `g` canonical; `r` hidden alias with the same behavior | `update.go:113`; live-project probe: `g` and `r` render identical views; footer shows only `g:goals`; help shows only `g` | Proven |
| 7 | `g` does not collide with action keys | `actions.go:22-24` (`p`, `a`, `x`) are the only action keys; no other `"g"` literal in the package | Proven |
| 8 | Detail heading and sections | `GOAL DETAIL`, `GOAL PROMISE`, `GOAL READINESS`, "(no Objectives in this Goal)"; tests pass | Proven |
| 9 | Selection write status and errors | `io.go`: "Goal selection", "Goal selected", "belongs to Goal"; tests pass | Proven |
| 10 | Reload after the selected Goal is removed | "Goal R-002 no longer exists; selection cleared." asserted | Proven |
| 11 | Reload after an open Goal detail's record is removed | Probe: status `RELEASE R-003 no longer exists; detail closed.` (`update.go:587`) | **Issue I-036** |
| 12 | Narrow widths | `TestGoalSelectorFitsNarrowBoardWidths` (32, 48); detail width tests | Proven |
| 13 | Resume and board Next share one phrase source | `internal/resume` phrases changed once; board imports them; `main_resume_matrix_test.go` updated | Proven |
| 14 | Rendering does no new IO or gate evaluation (ARCH-02) | Changes are string-only in render paths; no new calls | Proven |
| 15 | V2 skills and references byte-identical to scaffold copies (TPL-01) | `cmp -s` on all 7 pairs: identical | Proven |
| 16 | Active repository guidance uses Goals (README, AGENTS, Design, Guardrails) | Remaining Release hits are the storage and compatibility boundary, internal names, or the unrelated "Release And Distribution" heading | Proven |
| 17 | Shipped V2 scaffold `.savepoint/` files and live router use Goals | Scaffold `Design.md:55`, `Guardrails.md:47`, `router.md:40,57`; live `router.md:31,45,48` still say Release Check or Releases | **Issue I-035** |
| 18 | Goals are not described as publish, deploy, tag, or changelog; waiver is not CLEAR | Skills, AGENTS.md, README, and Design state both | Proven |
| 19 | Archive, immutable Checks, and prior Issues untouched | `git status` shows no changes under `.savepoint/archive`, `checks/`, or `issues/` before this run | Proven |
| 20 | V1 board unchanged | No diff under `internal/board/*.go` (V1) | Proven |
| 21 | `git diff --check` | No output | Proven |
| 22 | Fresh `make test-full` (TEST-08) | 2026-09-23, go1.26.2 linux/amd64, working tree above; all packages plus Linux/Darwin/Windows builds; EXIT 0 | Proven |
| 23 | Interactive TUI by hand | Not run: agents may not run `savepoint` commands. Covered by the model-level probes in rows 6 and 11 | N/A (policy) |

## Adversarial Pass

- Alias bypass: `r` goes through the same `openReleaseSelector`, and no help
  or hint text exposes it.
- A detail overlay suppresses the selector key: `m.Detail != nil` returns
  before the `g`/`r` check. Existing behavior, consistent for both keys.
- Kind-string leak: `DetailKind` "RELEASE" is used raw in two places. The
  heading was fixed; the reload status was not (I-036).
- Removing a Goal that Objectives still reference makes the index refuse the
  project, which shows the load diagnostic. That is correct, not a wording
  path.

## Issues

| Issue | Likelihood | Impact | Materiality | Recommendation |
|---|---|---|---|---|
| I-035 | High: every `savepoint init` V2 project and this repo's router carry it | Medium: agents and owners learn "Release Check" while the board and skills say "Goal Check"; no behavior impact | Medium | Fix now; a handful of wording edits |
| I-036 | Low: needs an open detail on an empty Goal whose file is deleted | Low: one status line uses the old term | Low | Fix together with I-035 |

## Remediation Routing

Hand remediation to the executor or planner as new or newly selected work
linked to O-013, I-035, and I-036. T-006, T-007, and T-008 are `done` and
must not be retreated. The re-check reuses this frozen scope lock.

## Owner Validation Still Needed

T-006 and T-007 declare `owner_validation.required: true` with an empty
`accepted_check`. Owner acceptance must name the current CLEAR Check once
one exists; this run cannot satisfy it.

## Observations (non-blocking)

- `internal/doctor` still reports "Release Check" and "Release evidence" in
  its readiness diagnostics (`checks.go:272-318`, `report.go:222`). Doctor
  is not in O-013's in-scope list, but its wording now differs from resume
  and the board. Worth a follow-up.
- The router's `next_action` text went stale during T-007/T-008; it still
  described T-006. It has since been updated to route to this Check.
- Task evidence does not record the executor session IDs, so
  `executed_session` is `unrecorded`.
