---
id: C-917
scope: {kind: objective, id: O-014}
result: CLEAR
checked_by: {role: checker, session: o014-objective-recheck-20260924}
executed_session: o014-i047-repair-20260924
checked_at: '2026-09-24T09:49:09Z'
reviewed:
  base_commit: 0287690
  head_commit: d69d732d19eacb7bba4d03259fbce13fb679130c
  files:
    - 'internal/resume/resume.go#sha256=8bc115a7a6c179aea14bb5c2b9c076646f09144d5c4f7c62e44334b29700191c'
    - 'internal/data/next.go#sha256=40a19646b61c83882a30b3596db1dc49fb8d9ecf89a454c0700ccf85d705243f'
    - 'internal/board/v2/io.go#sha256=ff5b0b31642764d5ae9a5697acf986be522aec4dd4706e007065d09f2a7380a4'
    - 'internal/board/v2/model.go#sha256=198b04cb83596ba36a1dc585e706610ea35624d03cfd31bc741671b43ccfe18b'
    - '.savepoint/Design.md#sha256=579e23aed0856a28298a9d53893608acac4406c369312bc2fc4a2a98de34e3d3'
    - 'AGENTS.md#sha256=f31ab53b3560a35db865427db886708d24e32a6d0b2fe89aaaf7c924ebaa99d9'
  dependencies: []
issues: []
supersedes: C-916
---

# C-917: O-014 Full Objective recheck

## Closure Map

| Prior Issue | Recheck outcome | Proof |
| --- | --- | --- |
| I-044 | Closed as verified | Board start now writes a planned Objective to `in_progress`; doctor names an interrupted write; `Fix` labels open Issues. |
| I-045 | Closed as verified | Issue-only selection opens the unique linked Objective, or a labelled all-Objectives view when no unique owner exists. |
| I-046 | Closed as verified | Live and scaffold task/Issue guidance now advances an Issue-only router after `repair_attempted`, preserving `release:`. |
| I-047 | Closed as verified | Next lines lead with an action verb and use `Accept` or `Close` for a current Check's owner step. |
| I-043 | Closed as verified | Design section 4 says Issues carry no stage, as the decoder enforces. |

## Independence, Scope, and Admission Ledger

This checker session built none of O-014's Tasks or direct Issue repairs. The seven owned Tasks, T-028 through T-034, are recorded done with explicit owner Task-check waivers. This recheck uses C-916's scope lock: selected-record resolution, router decoding and writes, Next rendering, board selection and closure, doctor diagnostics, guidance and Design reconciliation. The owner amended O-014's confirmed line format and Issue repair guidance after C-916; these explicit decisions supersede the older examples inside T-028 and T-031. No network or external boundary, text-width rule, or migration output enters this scope.

| Recheck item | Prior claim | Frozen cell / owner amendment | Allowed result |
| --- | --- | --- | --- |
| Objective lifecycle write and warning | I-044 | C-916 Objective planned/build and all-done rows; owner status decision | Verify or retain Issue |
| Issue-only board filter and unfiltered label | I-045 | C-916 Issue-only selection and board surface; owner board-view decision | Verify or retain Issue |
| Router advance after direct Issue repair | I-046 | C-916 agent guidance and Issue-only selection; owner routing decision | Verify or retain Issue |
| Action-first line, owner waits, Goal rungs | I-047 | C-916 line parity and Goal rows; owner format decision | Verify or retain Issue |
| Unchanged selection, diagnostics, closure, preservation | O-014 original criteria | C-916 remaining selection and board-write rows | Verify or retain Issue |
| Issue stage in Design | I-043 | C-916 Design section 4 observation | Verify or retain Issue |

## Coverage Matrix and Workflow Effects

Every row was classified before the verdict. The source trace, named regression tests, and fresh Full gate below cover the original cells and the owner's amended cases. A live `savepoint resume` supplied an independent scenario beyond a repeated fixture: with the router on O-014 and all seven Tasks done, its first line is `Check O-014 — Give the Next area an Objective word and let the router target Issues`, with no unrelated Task substituted.

| Row / axes | Normal and boundary cells | Failure, bypass, and representation cells | Result |
| --- | --- | --- | --- |
| Router read and write | Objective, Task, Issue-only, Goal, `none`, absent optional `issue` | malformed/wrong-type/missing ID; stale mtime; CRLF, unknown fields, retired `next_action` | Proven by `TestReadStateV2_*Issue*`, `TestWriteRouterStateV2_*`, `TestResolveSelection_*` |
| Selected Next resolution | planned/build/test/audit/done Task; Taskless Objective with no Tasks, unfinished Tasks, or all done; Issue open/in_progress/resolved; Goal check/owner/ready | missing or mismatched selection, blocked dependency, replan, done selection; no project-wide fallback | Proven by `TestResolveNext_selectedObjectiveDoesNotChooseUnrelatedReadyTask`, `TestResolveNext_selectedTasklessObjectiveWinsWhenOtherObjectiveHasActiveTask`, `TestResolveSelection_issueAloneBecomesNextAndResolvedIssueIsStale`, Full gate |
| Text output | board Next, non-TTY, resume first line; Task/Objective/Issue/Goal verbs | missing Objective record, stale diagnostic, Issue context, no selection, `NO_COLOR` rendering | Proven by `TestNextLineFormatsEverySelectionShape`, `TestBoardNextAndResumeReportTheSameAnswer`, `TestBuiltBoardAndResumeReportTheSameAnswer`, `TestRenderNextAccentsOnlyTheVerb` and live resume |
| Board startup and reload | selected Objective; Issue with one linked Objective | Issue with no or multiple Objective links, unknown Issue, unfiltered label | Proven by `TestIssueOnlySelectionOpensOnTheIssuesObjective`, `TestIssueOnlySelectionWithUnknownIssueOpensUnfiltered`, `TestFilteredViewHasNoAllObjectivesLabel` |
| Board start and closure | planned Task start writes Task then Objective; normal, waiver, exception Task closure; Objective closure; `p` | Objective write failure, no rewrite of already started/done Objective, retreat, router mtime conflict, nonmatching selection, blocked next Task | Proven by `TestTaskStartMovesPlannedObjectiveToInProgress`, `TestTaskStartLeavesInProgressAndDoneObjectivesUntouched`, `TestStartedTaskActionReportsObjectiveWriteFailure`, `TestTaskRetreatLeavesObjectiveStatusAlone`, `TestTaskCompletionMovesRouterToLowestUnfinishedBlockedTask`, `TestTaskWaiverCompletionAdvancesRouterAndKeepsReleaseAndIssue`, `TestRouterConflictAfterCompletionKeepsRecordAndReportsStaleSelection`, `TestObjectiveExceptionCompletionClearsOnlyObjectiveAndTaskSelection` |
| Doctor and guidance | planned Objective with started Task; done selection; retired field; Issue repair handoff | planned-only Objective has no warning; repair preserves Goal selection; live/scaffold skills match | Proven by `TestRunV2ChecksWarnsWhenPlannedObjectiveHasStartedTask`, `TestRunV2ChecksWarnsWhenRouterSelectsDoneTask`, `TestRunV2ChecksReportsRetiredNextActionWithoutItsValue`, source review and byte comparison |

The write sequence is input/selection validation, record write, reload of router and index, then router write. Record-write failure stops before router mutation; router conflict after completion leaves the completed record and surfaces a stale-selection error. Starting a Task writes its Objective next; an Objective-write failure leaves the Task started and gives doctor a named warning. No cleanup, network, subprocess, timeout, or secret-bearing external boundary applies to these paths. The independent oracle for writes is the persisted file bytes and reload, not only the action message. The closure tests verify `release:` preservation and the no-op and conflict cases.

## Acceptance Coverage

All O-014 Success Conditions are **Proven**. T-028 to T-034 outcomes are covered by the matrix, including waived Task Checks; no Task is reopened. In particular, `resume.NextVerb` produces the owner's amended action-first format; `routerObjective` uses linked Task and Objective Check owners for Issue-only board startup; the task skill and Issue-capture reference instruct the same post-repair router transition. Design sections 1, 4, and 8 describe the shipped behavior. The live and scaffold copies of all four skills and three references are byte-identical. The live and template AGENTS guides carry the same O-014 routing rule; their other project-specific sections differ intentionally.

## Gates

- `make test-full` passed on 2026-09-24 with Go 1.26.2 linux/amd64: all package tests, then Linux, Darwin, and Windows builds; exit 0.
- `make build` passed after the Full gate; `git diff --check` was clean.
- The working tree was clean before writing this Check and the Issue-resolution metadata. Those metadata writes do not change code, tests, fixtures, dependencies, or gate definitions.

## Materiality and Observations

No verdict Issue remains, so no materiality action is required. C-916's observations about migration's retired `next_action`, doctor repair wording, and waived-Task diagnostics remain outside this Check's frozen scope. The older Task records still quote the superseded status-first line; O-014's later owner-confirmed decision, Design, guidance, code, and tests consistently use the action-first line.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [x] STYLE-03 **Test branches**
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [x] STYLE-07 **One source of truth**
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data**
- [ ] STYLE-10 **Small diffs** — the reviewed O-014 work spans 59 files from C-916's base; this is advisory.

## Owner Decision Still Needed

The owner may now record O-014 `done`. This Check neither changes the Objective status nor accepts the Goal outcome.
