---
id: C-928
scope: {kind: objective, id: O-024}
result: NEEDS WORK
checked_by: {role: checker, session: o024-objective-check-20260925}
executed_session: unrecorded-external-executor
checked_at: '2026-09-25T08:22:37Z'
reviewed:
  base_commit: 4080fe21665bc3abee36093bec740425dd51d16b
  files:
    - internal/data/project.go
    - internal/data/project_test.go
    - internal/board/v2/objectives.go
    - internal/board/v2/objectives_test.go
    - internal/board/v2/plain.go
    - internal/board/v2/update.go
    - internal/board/v2/actions_test.go
    - internal/doctor/v2_runtime_test.go
    - .savepoint/Design.md
  dependencies: []
issues: [I-055]
supersedes: null
---

# C-928: O-024 Full Objective Check

## Independence and Scope

This session did not build O-024. It started fresh after `/clear`. The work is
uncommitted on branch `v2` on top of `4080fe2`. O-024 owns one Task, T-044,
which is `done` with an owner Task-check waiver recorded by the board. The
waiver was inspected and is not treated as technical CLEAR. T-044 was
reviewed directly against its Done When list and O-024's Success Conditions.

### Scope lock

1. Criteria: O-024 Outcome, Success Conditions 1–7, Confirmed Design
   Decisions; T-044 Done When; guardrails ARCH-02, DATA-02, FS-01,
   TEST-01..04, TEST-08, STYLE-07; Design reconciliation (§1, §8).
2. Entry points: `data.OrderedObjectiveIDsForGoal`, `data.LoadV2Index`
   (`DuplicateObjectiveRanks`), `doctor.RunV2Checks`, the V2 board TUI
   (`renderSidebar`, `visibleObjectiveWindow`, `handleSidebarKey` →
   `sidebarObjectiveOrderChange`), and `renderPlain`.
3. Relied-on behavior: `data.WriteObjectiveGroupOrderV2` (writes only the IDs
   passed to it), board reload with cursor restore by ID, `ResolveNext`, and
   router selection (must be unchanged).
4. Matrix axes: status {planned, in_progress, done}; priority {each of 4,
   missing}; rank {ranked, unranked, duplicate among done, duplicate open vs
   done, duplicate among open}; group shape {all four groups + DONE, no done,
   only done, empty Goal}; keys {1–4, K, J, shift+↑, shift+↓} × row {done,
   first open, last open before DONE, middle open}; output {TUI, non-TTY};
   window {heights 6–24 × every cursor row}; router selection {open, done};
   reopen {done → planned with rank collision, then heal}.
   Not applicable: cross-Goal moves and key focus outside the sidebar
   (unchanged since C-926/C-927; no code in this diff touches them); external
   boundaries (none).
5. An Issue must violate a Success Condition, Done When item, guardrail, or
   the Objective Check's Design-reconciliation duty through a supported path.

## Coverage Results

| Criterion | Evidence | Result |
| --- | --- | --- |
| SC1 open rows under priority headings (priority, rank, ID); one `DONE` heading, done rows by ID; no `TO DO`; empty groups hidden | Probe: 6 open rows across all 4 priorities (incl. in_progress, unranked) + 4 done rows with mixed priority, colliding ranks, unranked, out-of-ID-order index — sidebar sequence matched a hand-written oracle exactly; empty Goal drew no heading; existing no-done and only-done tests | Proven |
| SC2 piped output same rows, headings, order | Same probe compared plain sequence to the same oracle; real `./savepoint board \| cat` shows CRITICAL, HIGH, then DONE with O-001…O-023 in ID order | Proven |
| SC3 order keys no-op on done rows; K/J only among open rows of the group; priority change appends after open rows | Probe via full `Model.Update`, including `tea.KeyShiftUp`/`KeyShiftDown` messages (existing tests cover shift arrows only at helper level): all 8 keys × 4 done rows wrote nothing and kept order. J/shift+↓ on last open row before DONE and K on first row: no-op. Moving LOW O-010 to CRITICAL landed after the two open critical rows at rank 3 while two done critical rank-1 rows existed | Proven |
| SC4 reorder never rewrites a done file; duplicate-rank warning open-only | Probe: every open move left all 4 done files byte-identical; done/done and open/done rank collisions produced no fact; an open/open collision still produced exactly one fact naming only the open IDs; doctor test `TestRunV2ChecksIgnoresDuplicateRanksOnDoneObjectives`; real `./savepoint doctor` shows no rank finding | Proven |
| SC5 reopen returns to recorded rank; collision warned then healed | `TestReopenedDoneObjectiveReturnsToRankedGroupAndHealsCollision` rewrites the file on disk, reloads, and heals via the writer; code path reviewed | Proven |
| SC6 cursor, selection, badges, markers, Next, router unchanged; router-selected done row shows `●` under DONE | Probe: router selecting done O-009 renders `●` on its row below the `DONE` heading; plain prints `Selected: O-009`. No change to `ResolveNext`, router, or badge code in the diff | Proven |
| SC7 tests listed; `git diff --check`; build and fast gate | Listed cases present in `project_test.go`, `objectives_test.go`, `actions_test.go`, `v2_runtime_test.go`; gates below | Proven |
| T-044 windowing counts every heading | Probe: heights 6–24 × every cursor row — never more lines than the height, cursor row and its own heading (including `DONE`) always visible | Proven |
| Architecture: data owns order; board reads heading from status | `objectiveRowsForRelease` renders `OrderedObjectiveIDsForGoal` without sorting; `objectiveRowHeading` reads `Status`, not position | Proven |
| Design reconciliation | §1 line 29, §8 Layout (line 179) and Keybindings (line 210) still describe the O-020 order and key scope; T-044 Drift Notes promised the update | **Issue I-055** |

## Commands

- `git diff --check` — clean.
- `make build` then `make test-full` — exit 0 on the working tree reviewed
  here, `go1.26.2 linux/amd64`, 2026-09-25 ~08:20Z. Cross-builds for linux,
  darwin, and windows succeeded.
- Temporary probe test `internal/board/v2/zz_o024_check_probe_test.go`
  (8 tests) — all passed; file removed after the run.
- `./savepoint board | cat` and `./savepoint doctor` on this repository —
  grouped list with DONE section rendered; no rank diagnostics.

## Issues

- **I-055** (drift, low): Design §1 and §8 do not describe the `DONE`
  section, done-last ID order, done rows excluded from duplicate-rank
  detection, or order keys being no-ops on done rows.

## Materiality

| Issue | Likelihood | Impact | Materiality | Recommendation |
| --- | --- | --- | --- | --- |
| I-055 | High (every reader of §1/§8) | Low (code, tests, and the Objective record are correct; only the architecture document lags, and §8 now states an order the board no longer uses) | Low | Fix now: a planner-only Design edit, then a targeted recheck |

## Owner Validation Still Needed

T-044 declares `owner_validation.required: true` with an empty
`accepted_check`. After a CLEAR recheck, the owner's acceptance should name
that Check. The T-044 User Check (interactive keys in a real terminal) was
covered here by the board's `Update` path and rendered output, not by an
interactive terminal.

O-024's Confirmed Design Decisions flag one planner choice for owner review:
done Objectives are left out of reorder writes and duplicate-rank detection.
The implementation follows that choice. No record shows the owner confirming
it.

## Observations (non-blocking)

- T-044's Technical Evidence ends "no owner waiver has been recorded", but the
  board later added a `check_waiver` block to the same record. The line is
  stale, not wrong at the time it was written.
- `savepoint doctor` still reports `v2-done-without-clearance` for
  T-041..T-043 (O-020), the pre-existing behavior for waived Tasks with
  `owner_validation.required`. Not introduced by O-024.
- The help overlay does not say order keys are ignored on done rows. Pressing
  one does nothing and shows no message. Not required by any criterion.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [x] STYLE-03 **Test branches**
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [ ] STYLE-07 **One source of truth** — `internal/board/v2/objectives.go` `objectiveRowHeading` repeats the missing-priority → Medium default that `update.go` `sidebarRowPriority` already owns.
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data**
- [x] STYLE-10 **Small diffs**
