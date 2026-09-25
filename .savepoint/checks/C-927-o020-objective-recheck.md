---
id: C-927
scope: {kind: objective, id: O-020}
result: CLEAR
checked_by: {role: checker, session: o020-objective-recheck-20260925}
executed_session: unrecorded-external-executor
checked_at: '2026-09-25T07:29:07Z'
reviewed:
  base_commit: 1651018d1625ec0adea2e3b0e1e906179492e000
  files:
    - .savepoint/Design.md
    - .savepoint/Guardrails.md
    - .savepoint/objectives/O-020-rank-release-objectives/Objective.md
    - .savepoint/objectives/O-020-rank-release-objectives/tasks/T-041-give-each-objective-a-priority-and-a-place-in-line.md
    - .savepoint/objectives/O-020-rank-release-objectives/tasks/T-042-show-objectives-in-priority-groups.md
    - .savepoint/objectives/O-020-rank-release-objectives/tasks/T-043-reorder-objectives-from-the-keyboard.md
    - internal/data/objective_v2.go
    - internal/data/objective_v2_test.go
    - internal/data/project.go
    - internal/data/project_test.go
    - internal/data/write.go
    - internal/data/write_test.go
    - internal/doctor/v2_runtime.go
    - internal/doctor/v2_runtime_test.go
    - internal/board/v2/objectives.go
    - internal/board/v2/objectives_test.go
    - internal/board/v2/plain.go
    - internal/board/v2/releases_test.go
    - internal/board/v2/update.go
    - internal/board/v2/io.go
    - internal/board/v2/help.go
    - internal/board/v2/actions_test.go
    - agent-skills/savepoint-design/SKILL.md
    - templates/project-v2/agent-skills/savepoint-design/SKILL.md
  dependencies: []
issues: []
supersedes: C-926
---

# C-927: O-020 Full Objective Recheck

## Closure Map

| Prior Issue | Status | Evidence |
| --- | --- | --- |
| I-054 | Closed as verified | Design §8 now describes the grouped sidebar and plain list, the order keys, and the Objective order writer; current sources and tests match that description. |

## Independence and Scope

This is a fresh checker session, independent from C-926 and the implementation
session. I did not implement or edit the Design correction. This is a Full
Objective recheck of O-020, including T-041, T-042, T-043, cross-Task
integration, and reconciliation against Design.

All three Tasks remain `done` with their recorded owner Task-check waivers.
Each still has `owner_validation.required: true` and an empty
`accepted_check`; those owner decisions remain pending and are not supplied by
this Check.

## Frozen Scope Lock

The lock below is unchanged from C-926. No criterion, entry point, matrix axis,
dependency layer, or interpretation was added during this recheck.

1. **Criteria:** O-020 Outcome, Success Conditions 1–12, and Confirmed Design
   Decisions; Done When of T-041..T-043; guardrails FS-01, FS-04, DATA-01..04,
   TPL-01, ARCH-02..04, TEST-01..04, TEST-08; Design reconciliation.
2. **Entry points:** `data.DecodeObjectiveV2`, `data.LoadV2Index`
   (`DuplicateObjectiveRanks`), `data.OrderedObjectiveIDsForGoal`,
   `data.WriteObjectiveGroupOrderV2`, `doctor.RunV2Checks`, the V2 board TUI
   (`renderSidebar`, `handleSidebarKey`, `writeObjectiveGroupOrderCmd`,
   `renderHelp`), and `renderPlain`.
3. **Relied-on behavior:** `writeV2Record` (atomic replace, freshness check,
   frontmatter node patching), board reload with cursor restore by ID,
   `ResolveNext`, and router writes (must be untouched).
4. **Matrix axes:** priority value {missing, each of 4, quoted, wrong case,
   unknown, null, int, list}; rank {missing, 1, large, 0, negative, float,
   quoted, word, null, overflow, map}; order {4 groups, ranked/unranked,
   ties, legacy, unknown Goal, nil index}; writer {no-op, partial renumber,
   stale member, deleted member, duplicate ID, missing ID, wrong Goal,
   unknown Goal, bad priority, injected mid-write failure, CRLF, flow
   sequences, comments, quoted title}; keys {1–4, current priority, K, J,
   shift+↑, shift+↓, group edges} × focus {sidebar, columns, detail overlay,
   Issues overlay}; output {TUI, NO_COLOR, narrow, non-TTY}. Not applicable:
   a board action for cross-Goal moves (confirmed out of scope); external
   network boundaries (none).
5. **Issue boundary:** an Issue must violate a Success Condition, Done When
   item, guardrail, or the Objective Check's Design-reconciliation duty
   through a supported path.

## Re-check Admission Ledger

| Re-check item | Prior finding or repair claim | Exact frozen cell | Allowed result |
| --- | --- | --- | --- |
| §8 Layout and plain Objective list | I-054: Design omitted the priority headings, removed lifecycle line, and grouped plain list | Design reconciliation; output {TUI, NO_COLOR, narrow, non-TTY}; O-020 SC2, SC9, SC10; T-042 Done When | Clear only if Design describes the current display and the source/tests agree; otherwise Issue or Unverified. |
| §8 persistence and refresh | I-054: Design omitted the Objective order write and its failure behavior | Design reconciliation; writer {stale member, partial renumber, deleted member, injected mid-write failure, comments, quoted title}; O-020 SC6; T-041/T-043 Done When | Clear only if Design matches the current writer's preflight, per-file replace, and partial-write recovery; otherwise Issue or Unverified. |
| §8 keybindings | I-054: Design omitted `1`–`4`, `K`/`J`, and shifted arrows | Design reconciliation; keys {1–4, current priority, K, J, shift+↑, shift+↓, group edges} × focus {sidebar, columns, detail overlay, Issues overlay}; O-020 SC4; T-043 Done When | Clear only if keys, actions, and focus scope match dispatch and help; otherwise Issue or Unverified. |
| Remaining O-020 behavior | C-926's Full Objective evidence covered the unchanged implementation and its acceptance matrix | The remaining exact priority, rank, order, writer, key/focus, and output cells in the frozen matrix above | Retain only where current source/tests are unchanged and the current full gate passes; otherwise Issue or Unverified. |
| Full Objective gate | C-926 recorded a successful full run | O-020 SC12; TEST-08 | Require current successful `make build && make test-full`; a failed or unavailable gate means Unverified and NEEDS WORK. |

## Re-check Coverage

The updated §8 text at `.savepoint/Design.md:179,206,210` now resolves all
three parts of I-054:

- **Layout:** §8 names the non-empty Critical/High/Medium/Low groups, priority
  then rank then ID order, the matching grouped Objective list in plain
  output, and removal of the Planned / In Progress / Done line. This matches
  `objectiveRowsForRelease`, `renderSidebar`, `renderObjectiveRow`, and
  `renderPlain` in `internal/board/v2/objectives.go` and `plain.go`.
- **Persistence:** §8 says order actions run through
  `data.WriteObjectiveGroupOrderV2`, validates the full write set against the
  selected Goal before writing, patches only priority and rank, atomically
  replaces each changed file, and describes the non-transactional partial
  write and next-move recovery. `WriteObjectiveGroupOrderV2` performs the
  membership and freshness preflight; `writeV2Record` writes through
  `replaceV2File` using a same-directory temporary file and rename. The
  existing writer and board tests cover no-op/preservation, stale-member
  refusal, partial failure, loadability, and recovery.
- **Keys:** §8 names `1`–`4` for priority and `K`/`J` or shifted arrows for
  within-group movement, append and group-edge behavior, cursor following,
  and sidebar focus. `handleSidebarKey`, `sidebarObjectiveOrderChange`, and
  `renderHelp` implement that behavior; the full suite includes the priority,
  movement, group-edge, cursor, conflict, and focused-help cases.

The rest of the frozen acceptance matrix remains supported by the unchanged
implementation and its named tests: Objective decode diagnostics and defaults
(`TestDecodeObjectiveV2_rejectsInvalidPriorityAndRankWithNamedDiagnostics`),
canonical ordering and duplicate-rank reporting
(`TestOrderedObjectiveIDsForGoalUsesPriorityRankAndID`,
`TestLoadV2IndexReportsDuplicateObjectiveRanksByGoalAndPriority`),
cross-Goal writer validation, TUI/plain ordering and badges
(`TestSidebarAndPlainOutputUsePriorityAndRankOrder`), group visibility,
no-colour/narrow rendering, keyboard behavior and cursor retention
(`TestSidebarPriorityKeysAppendAndKeepSelection`,
`TestSidebarGroupMovementRenumbersLegacyRowsAndClamps`,
`TestSidebarGroupMovementDoesNotCrossPriorityGroups`), conflict behavior
(`TestSidebarReorderConflictKeepsOrderAndExternalEdit`), and partial-write
recovery (`TestSidebarOrderHealsAfterInjectedMidWriteFailure`). The current
full gate reran every package test. The scoped implementation and test files'
recorded modification times predate C-926's 06:50:11Z gate; the worktree has no
dependency or gate-definition edits. C-926's temporary probe files had been
removed after that run, so I relied on its recorded probe outcomes only after
checking the unchanged scoped sources and rerunning the full gate.

## Full Gate Evidence

- `make build && make test-full` — passed in this recheck on
  `go1.26.2 linux/amd64`. `make test-full` ran `go test -json -count=1 ./...`
  and successfully built Linux, Darwin, and Windows targets.
- `git diff --check` — passed.
- `cmp agent-skills/savepoint-design/SKILL.md
  templates/project-v2/agent-skills/savepoint-design/SKILL.md` — passed.

## Issues

No material Issues remain in the frozen O-020 scope. I-054 is closed as
verified after this CLEAR Check. No new Issue was found.

## Materiality

No materiality actions are required; the sole material finding from C-926 was
the Design drift now repaired and verified.

## Owner Validation Still Needed

The owner must still record the required validation for T-041, T-042, and
T-043, naming the current Check C-927 as directed by the Task evidence. This
Check does not set their `accepted_check` fields, change Task status, or accept
the Objective on the owner's behalf. Objective closure remains an owner
decision after those Task validations.

## Observations

The non-blocking observations recorded by C-926 remain outside the frozen
blocking perimeter and do not change this verdict: rapid repeated `K` input
can collapse before reload; the order command shares the loaded index while
the writer updates it; a deleted member after load reports a plain `stat`
error; doctor reports the pre-existing `v2-done-without-clearance` condition
for the waived Tasks; and §8 retains older `d`, `D`, `A`, and `R` key entries
outside O-020's promised additions. No new observation changes the result.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [x] STYLE-03 **Test branches**
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [ ] STYLE-07 **One source of truth** — `internal/board/v2/update.go` uses
  string literals for priority values instead of `data.ObjectivePriority*`.
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data**
- [ ] STYLE-10 **Small diffs** — the pre-existing O-020 implementation diff
  spans 24 files and approximately 1,600 added lines, mostly tests.
