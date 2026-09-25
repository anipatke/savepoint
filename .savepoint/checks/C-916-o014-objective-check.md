---
id: C-916
scope: {kind: objective, id: O-014}
result: NEEDS WORK
checked_by: {role: checker, session: o014-objective-check-20260924}
executed_session: unrecorded-o014-executor-sessions
checked_at: '2026-09-24T08:58:19Z'
reviewed:
  base_commit: 1f82cf2
  head_commit: 1e35d1fb8120453f4a4d503184ac55a43d253029
  files:
    - 'internal/resume/resume.go#sha256=c01c9b01d9ac296dea50e0656d71d10f08d063b2deb9324a37325ff4e64978dc'
    - 'internal/data/next.go#sha256=40a19646b61c83882a30b3596db1dc49fb8d9ecf89a454c0700ccf85d705243f'
    - 'internal/data/write.go#sha256=115002e56443bfa47c2b8480fb06cbbbeb7ac2efb4b120eefc52887131a57b0d'
    - 'internal/board/v2/io.go#sha256=f2adcace86e8b7e5d46cac306c7dd9b131b539faa696552cae49a52ed6153192'
    - 'internal/doctor/v2_runtime.go#sha256=6b7d738626bf0f663fc775b4421a5a2841820590e3706ad6e27f1d6b37be3b97'
    - '.savepoint/Design.md#sha256=824f868fb6dcafcabc44db11b267280ffa220df30773b29750eed4aa876c236f'
    - 'AGENTS.md#sha256=b3903382a29cfed64f72a99f826353d51f756370b7aa72d13ef0e887710edb16'
  dependencies: []
issues: [I-044]
supersedes: null
---

# C-916: O-014 Full Objective Check

## Independence And Scope

This session did not build any O-014 Task. In this session it planned O-023
and edited `check-method.md`, and neither touches O-014's code. The reviewed
state is HEAD `1e35d1f` plus the uncommitted T-030 to T-034 working tree
(`git diff` sha256 `a17a0522…c23e`). That hash was identical before and after
the full gate, so nothing changed during review.

Scope lock:

1. Criteria: O-014's Success Conditions, the Confirmed Design Decisions, and
   Guardrails DATA-01, DATA-03, ARCH-02, TPL-01, TPL-02, TEST-01, TEST-02,
   TEST-08.
2. Surfaces: `data.ResolveSelection`/`ResolveNext`, `resume.NextLine` and
   `Render`, board Next area and plain table, board Space/exception closure
   and `p`, `ReadStateV2`/`WriteRouterStateV2`,
   `RouterSelectionAfterClosureV2`, doctor V2 runtime, AGENTS.md, phase skills
   and scaffold copies, and Design sections 1, 4 and 8.
3. Not applicable: network and external boundaries (none); text-width
   classes (the line is plain text, and wrapping is pre-existing board code).
4. Out of scope: O-022 Goal enforcement, O-020 ranking, O-015 Issue actions,
   migration output.

## Next Line Matrix

A temporary harness (`zz_o014_permutations_test.go`, run in a scratch copy of
the tree and not committed) wrote each project to disk and ran the real
`CheckRuntimeSchema` → `LoadV2Index` → `ReadStateV2` → `ResolveNext` path.

| Router selection | Objective status | Line | Warning |
| --- | --- | --- | --- |
| Task planned / build / test / audit | planned | `Planned O-001 · Planned\|Build\|Test\|Check T-001 — …` | — |
| same | in_progress | `In Progress O-001 · … T-001 — …` | — |
| Objective only, work left | either | `<word> O-001 — …` | — |
| Objective only, no Tasks | either | `<word> O-001 — …` (plan Tasks) | — |
| Objective only, all Tasks done | either | `<word> O-001 · Check — …` | — |
| Done Task, all Tasks done | either | `<word> O-001 · Check — …` | finished Task |
| Done Task, work left | either | `<word> O-001 — …` | finished Task |
| All done, CLEAR Objective Check | either | `<word> O-001 · Check — …` (owner records done) | — |
| Objective done | done | `Done O-001 · Check — …` | finished Objective/Task |
| Nothing | — | `Nothing selected` | — |
| Issue alone: open / in_progress / resolved | — | `Open\|In Progress\|Resolved I-001 — …` | resolved Issue |
| Issue + Task at build | planned | Task line + `Issue: Open I-001 — …` | — |
| Issue or Task not found | — | `Nothing selected` | names the missing ID |

No row picks work the router didn't select. The live line
`Planned O-014 · Check — …` + `Warning: router still selects finished Task
T-034.` is the "done Task, all Tasks done" row. T-034 was marked done by
editing the file (waiver session `owner-chat` 08:48Z, router last written
08:38Z), not on the board, so the warning is the intended backstop.

## Acceptance Coverage

| Success Condition | Result | Evidence |
| --- | --- | --- |
| Identical line on board, non-TTY and resume | Proven | `TestBoardNextAndResumeReportTheSameAnswer`, `TestDoneTaskSelectionWarningIsSharedAcrossSurfaces`, `TestRouterSelectedIssueFlowsAcrossBoardAndResume`; live `./savepoint resume` |
| Next is the selection only; 2026-09-23 cases can't recur | Proven | Matrix above; `matrixBuildSelectedTasklessObjective`, `matrixBuildDoneSelectedTaskRelease` |
| No Objective → nothing selected + how to select | Proven | Matrix "Nothing" row; action names `p` and "set router to" |
| "set router to" agent action | Proven | AGENTS.md Router Selection |
| Router `issue:` decode, validate, `none`, write, not found | Proven | `TestReadStateV2_*Issue*`, `TestWriteRouterStateV2_setsIssueAndPreservesEveryOtherByte`, `…_addsIssueAlone`; matrix not-found row |
| SelectionDone on every surface | Proven | Matrix; live `savepoint doctor` prints `! … router still selects finished Task T-034` |
| Board closure advance; `release:` byte-for-byte; `p` | Proven | `TestTaskCompletionMovesRouterToLowestUnfinishedBlockedTask`, `TestTaskWaiverCompletionAdvancesRouterAndKeepsReleaseAndIssue`, `TestLastTaskCompletionClearsRouterTaskAndShowsObjectiveCheck`, `TestObjectiveExceptionCompletionClearsOnlyObjectiveAndTaskSelection`, `TestRecordSelectionPreservesReleaseAndIssueForTaskAndObjective` |
| `next_action` retired; old routers load; doctor reports it | Proven | `TestReadStateV2_recordsEmptyRetiredNextActionPresence`, `TestWriteRouterStateV2_doesNotAddRetiredNextAction`, `TestRunV2ChecksReportsRetiredNextActionWithoutItsValue`; no `next_action` in the live or template router |
| AGENTS.md maps the pasted line to a skill | Proven | AGENTS.md Workflow step 2 |
| Guidance names `savepoint resume`; missing binary; who advances; live and scaffold identical | Proven | AGENTS.md; `cmp` of all four skill pairs identical |
| Objective word states the Objective's lifecycle (Outcome example `In Progress O-014 · Build T-028`) | **Issue** | I-044 |
| Gates, Goal rungs, parity unchanged | Proven | Release matrix rows in `main_resume_matrix_test.go`; full gate |
| Design sections 1, 4, 8 reconciled | Proven | Design diff; section 4 also fixes I-043 |
| Gates | Proven | Below |

## Commands

- `make test-full`: exit 0, no `FAIL`, all three cross-builds (go1.26.2 linux/amd64).
- `git diff --check`: clean.
- `./savepoint resume` and `./savepoint doctor` on the live repository (binary built 08:44Z from this tree).

## Issues

- **I-044**: the Objective word never reads `In Progress`, because nothing
  writes an Objective's `in_progress`. The code maps status correctly; the
  recorded status is never kept current.

| Issue | Likelihood | Impact | Materiality | Recommendation |
| --- | --- | --- | --- | --- |
| I-044 | High: every Objective, every day | Low–Medium: misleading word on the line the owner pastes into agents; no gate effect | Medium | Fix now, as one Task under O-014; the owner picks derived-word or written-status |

## Observations (non-blocking, outside the frozen scope)

1. `savepoint migrate` still writes `next_action` into every V2 router
   (`internal/migrate/convert_docs.go:53`, no `omitempty`), so a freshly
   migrated project gets the retired-field doctor warning straight away.
   T-033 left migration alone on purpose. Worth a follow-up.
2. `· Check` shows in three different cases: the Check is needed, the Check
   is CLEAR and the owner only has to record done, and the Objective is
   already `Done`. AGENTS.md maps `Check` → `savepoint-check`, so an agent
   given only the line could re-run a Check that already passed. This is the
   confirmed design ("owner acceptance wait reads Check"); the owner may want
   to revisit it for the ready and done cases.
3. The doctor repair for a finished selection always says "Use p on an
   unfinished Task", which doesn't fit an all-done Objective or a resolved
   Issue. It also prints an absolute path, where other findings use relative
   paths.
4. Doctor lists every waived done Task as `✗ done but clearance is missing`,
   which contradicts the waiver policy in AGENTS.md. This predates O-014
   (O-001's Tasks show it too).
5. I-043's repair (Design section 4: Issues carry no `stage`) is present and
   correct. Close it `verified` on the CLEAR re-check.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [ ] STYLE-03 **Test branches**: every lifecycle fixture creates the Objective already `in_progress` or `planned` and never moves it, so the branch where a Task starts under a `planned` Objective is untested (I-044).
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [x] STYLE-07 **One source of truth**: board, non-TTY, resume and doctor all read `resume.NextLine` / `SelectionPhrase`.
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data**
- [ ] STYLE-10 **Small diffs**: T-030 to T-034 are one uncommitted ~2,000-line change across 39 files.

## Owner Validation Still Needed

After I-044 is repaired: a re-check, then the owner records O-014 as done.
