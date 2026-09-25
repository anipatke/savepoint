---
id: C-925
scope: {kind: objective, id: O-022}
result: CLEAR
checked_by: {role: checker, session: o022-objective-check-20260925}
executed_session: o022-repair-followup-20260925
checked_at: '2026-09-25T05:10:19Z'
reviewed:
  base_commit: 35ce82dffbcfa02a88121f99496d668c67d69db9
  files:
    - .savepoint/objectives/O-022-require-goal-context/Objective.md
    - .savepoint/objectives/O-022-require-goal-context/tasks/T-036-flag-work-that-has-no-goal.md
    - .savepoint/objectives/O-022-require-goal-context/tasks/T-037-keep-the-board-inside-one-goal.md
    - .savepoint/objectives/O-022-require-goal-context/tasks/T-038-start-new-projects-with-a-goal.md
    - .savepoint/objectives/O-022-require-goal-context/tasks/T-039-give-migrated-projects-a-live-goal.md
    - .savepoint/objectives/O-022-require-goal-context/tasks/T-040-say-a-goal-is-required-everywhere.md
    - internal/data/project.go
    - internal/doctor/checks.go
    - internal/doctor/v2_runtime.go
    - internal/board/v2/next_panel.go
    - internal/migrate/plan.go
    - internal/migrate/convert_docs.go
    - internal/migrate/testdata/golden/v1-router-missing.yml
    - internal/migrate/testdata/golden/v1-router-unresolvable.yml
    - internal/migrate/testdata/golden/v1-router-archived.yml
    - AGENTS.md
    - templates/project-v2/AGENTS.md
  dependencies: []
issues: []
supersedes: C-924
---

# C-925: O-022 Final Objective Recheck

## Why This Record Exists

The owner asked for a recheck after C-924 recorded CLEAR. This session wrote
C-923 and did not implement any repair. It reran C-923's own reproductions
against the repaired working tree. C-924's verdict is confirmed. This record
supersedes C-924 for one reason: C-924's `reviewed.files` names two Task paths
that do not exist (`T-036-report-a-missing-router-goal.md` and
`T-038-create-a-goal-on-init.md`). The real files are
`T-036-flag-work-that-has-no-goal.md` and
`T-038-start-new-projects-with-a-goal.md`. The verdict itself is unaffected.

The C-923 scope lock is unchanged. The owner's I-052 decision is now recorded
in O-022's Architectural Considerations and in the revised T-039 Done When.
It changes the expected migration fallback: reuse an identifiable live Goal
for active work, create a continuation only for selected work that sits in a
historical Goal and move that work into it, and keep an unresolved lifecycle
decision in preview. These criteria were checked against that recorded
decision.

## Closure Map (C-923 reproductions rerun)

| Issue | Reproduction rerun on a fresh binary | Result |
| --- | --- | --- |
| I-050 | Fresh init, R-001 deleted, router still `release: R-001`: doctor `[router-goal-missing] … does not exist …` with `Create a Goal first, then choose it with g on the board.`, exit 1. Board: `No live Goals exist; run savepoint doctor.` Archived-only Goal with router on it: same doctor repair, same board pointer. | Fixed |
| I-051 | Fresh `savepoint init` then `savepoint doctor`: `ALL CLEAN (exit code 0)`. | Fixed |
| I-052 | Materialized regenerated goldens (`schema_version: 2`). Missing/unresolvable: router R-001/O-001/T-001, resume `Build T-001 …`, board 1 in-progress card in Goal R-001. Archived: continuation R-002, O-001 moved to `release: R-002`, resume `Build T-001 …`, board 1 in-progress card. | Fixed |
| I-053 | Temporary probe (removed) with the real V1 router shape, `audited` and unrecognized status, active epic: preview lists `release_completion_evidence:` / `release_lifecycle:` as unresolved, `Appliable: false`, Apply refused. Decision `in_progress` selects live R-001. Decision `done` creates R-002 and moves the active Objective. | Fixed |

The remaining C-923 matrix rows (router none/blank/absent, Issue-only
selection, Objectives missing `release:`, malformed references, the Goal
selector, init contents) are covered by the fresh full suite. The repairs
touched none of their code paths except the shared `HasLiveGoal` predicate,
which the I-050 probes exercised.

## Gates

- `make test-full`: passed fresh at 2026-09-25T05:10Z on go1.26.2 linux/amd64,
  exit 0, including linux/darwin/windows builds.
- `git diff --check`: passed.
- Live vs scaffold: all four skills and three shared references are
  byte-identical; the AGENTS.md Required Goal Context sections match.

## Materiality

No Issues. No materiality actions are required.

## Observations (non-blocking)

- On a migrated project with active work, doctor still reports
  `v2-release-objective-incomplete` for the selected in-progress Goal. The same
  finding appears on this repository's own live Goal. It is Goal-readiness
  reporting that predates O-022. It was present in C-923's probes and is
  outside the locked Issue set, so the planner may want a separate Issue.
- C-924's `reviewed.files` phantom paths are described above. C-924 stays
  intact as history.
- C-923's other observations still apply: global counts in plain board
  output, absolute paths in doctor, the O-019 and O-022 working tree is
  uncommitted and mixed, and waived Tasks show as missing clearance in doctor.

## Owner Validation Still Needed

All five Tasks declare `owner_validation.required`. The owner accepts this
current CLEAR Check and then records O-022 as done. Goal R-006 still needs its
mandatory Goal Check.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [x] STYLE-03 **Test branches**
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [x] STYLE-07 **One source of truth** — board and doctor now share `HasLiveGoal`.
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data**
- [ ] STYLE-10 **Small diffs** — still uncommitted and interleaved with O-019.
