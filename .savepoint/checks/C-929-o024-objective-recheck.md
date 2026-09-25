---
id: C-929
scope: {kind: objective, id: O-024}
result: CLEAR
checked_by: {role: checker, session: o024-objective-recheck-20260925}
executed_session: o024-design-reconciliation-20260925
checked_at: '2026-09-25T08:45:28Z'
reviewed:
  base_commit: 4080fe21665bc3abee36093bec740425dd51d16b
  files:
    - .savepoint/Design.md
    - .savepoint/objectives/O-024-move-done-objectives-below-open-work/Objective.md
    - internal/data/project.go
    - internal/data/write.go
    - internal/board/v2/objectives.go
    - internal/board/v2/plain.go
    - internal/board/v2/update.go
  dependencies: []
issues: []
supersedes: C-928
---

# C-929: O-024 Full Objective Recheck

## Independence and Scope

This session wrote C-928 and did not build O-024 or the I-055 repair. The
repair was made in planner session `o024-design-reconciliation-20260925`. The
recheck uses C-928's frozen scope lock and matrix without adding any axis.

## Closure Map

| Issue | Status | Evidence |
| --- | --- | --- |
| I-055 | Closed as verified | Design §1 and §8 now describe the implemented order, `DONE` section, duplicate-rank scope, write-set exclusion, and done-row key behavior; sources and tests match. |

## Admission Ledger

| Re-check item | Prior claim | Frozen cell | Allowed result |
| --- | --- | --- | --- |
| §1 Objective planning order | I-055 | Design reconciliation × `project.go` order and duplicate facts | Clear / Issue |
| §8 Layout (sidebar and plain) | I-055 | Design reconciliation × output {TUI, non-TTY} × group shape | Clear / Issue |
| §8 Keybindings | I-055 | Design reconciliation × keys × row {done, open} | Clear / Issue |
| §8 Board persistence | I-055 repair touched it | Design reconciliation × reorder never writes done files | Clear / Issue |
| Code unchanged since C-928 | C-928 SC1–SC7 Proven | All C-928 matrix cells | Clear / Issue |
| Planner choice owner decision | C-928 Owner Validation note | Confirmed Design Decisions | Clear / Observation |

## Results

- **§1** (`Design.md:29`): open Objectives by priority, ranked before
  unranked, rank, then ID; finished Objectives after them in ID order;
  duplicate ranks reported among unfinished Objectives only. Matches
  `OrderedObjectiveIDsForGoal` and `duplicateObjectiveRankFacts`. Clear.
- **§8 Layout** (`Design.md:179`): one `DONE` heading after the priority
  headings, finished rows in ID order without priority subgroups, empty
  headings omitted, plain output in the same headings and order. Matches
  `objectiveRowHeading`, `renderSidebar`, and `renderPlain`. Clear.
- **§8 Keybindings** (`Design.md:210`): priority and reorder keys do nothing
  on a finished row; reorder moves only unfinished rows within their group;
  priority change appends to the destination group's end. Matches
  `sidebarObjectiveOrderChange` and `sidebarObjectiveIDsInPriority`. Clear.
- **§8 Board persistence** (`Design.md:206`): finished Objectives are excluded
  from reorder write sets and keep their files unchanged. Matches the board's
  write-set construction and C-928's no-write probes. Clear.
- **Code unchanged:** the `internal/` diff has the same per-file line counts
  as the tree C-928 reviewed; only `.savepoint/` records and `Design.md`
  changed since. Gates were rerun fresh regardless (below). Clear.
- **Owner decision:** O-024 now records the owner's approval of excluding
  finished Objectives from reorder writes and duplicate-rank checks
  (2026-09-25T08:44:06Z, owner, chat). Clear.

## Commands

- `git diff --check` — clean.
- `make build` — passed.
- `make test-full` — exit 0, finished 2026-09-25T08:45:28Z,
  `go1.26.2 linux/amd64`; cross-builds for linux, darwin, and windows
  succeeded.

## Issues

None. No materiality actions are required.

## Owner Validation Still Needed

T-044 declares `owner_validation.required: true` with an empty
`accepted_check`. The owner's acceptance should name C-929. The Objective
can then be closed by the owner.

## Observations (non-blocking)

- §8 Board persistence says the writer "patches only `priority` and `rank` on
  unfinished Objectives". The exclusion is enforced by the board when it
  builds the write set, not by `data.WriteObjectiveGroupOrderV2`, which would
  accept a finished Objective's ID if a caller passed one. No supported path
  does so. A future Design tidy could attribute the rule to the board.
- C-928's observations (stale "no owner waiver" line in T-044 evidence,
  pre-existing `v2-done-without-clearance` for T-041..T-043, help overlay
  silent on done-row keys) still stand and remain non-blocking.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [x] STYLE-03 **Test branches**
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [ ] STYLE-07 **One source of truth** — `internal/board/v2/objectives.go` `objectiveRowHeading` repeats the missing-priority → Medium default that `update.go` `sidebarRowPriority` already owns (unchanged since C-928).
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data**
- [x] STYLE-10 **Small diffs**
