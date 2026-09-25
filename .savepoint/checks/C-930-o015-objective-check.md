---
id: C-930
scope: {kind: objective, id: O-015}
result: NEEDS WORK
checked_by: {role: checker, session: o015-objective-check-20260925}
executed_session: o015-build-20260925
checked_at: '2026-09-25T09:30:30Z'
reviewed:
  base_commit: 1691dbca4889791d011732424e46354f92a93f31
  files:
    - internal/data/write.go
    - internal/data/errors.go
    - internal/data/issue_v2_test.go
    - internal/board/v2/io.go
    - internal/board/v2/issues.go
    - internal/board/v2/update.go
    - internal/board/v2/help.go
    - internal/board/v2/view.go
    - internal/board/v2/issues_test.go
    - internal/board/v2/footer_test.go
    - internal/init/agent_skills_test.go
    - AGENTS.md
    - templates/project-v2/AGENTS.md
    - agent-skills/references/issue-capture.md
    - agent-skills/savepoint-check/SKILL.md
    - agent-skills/savepoint-design/SKILL.md
    - agent-skills/savepoint-task/SKILL.md
    - .savepoint/Design.md
  dependencies: []
issues: [I-056, I-057]
supersedes: null
---

# C-930: O-015 Full Objective Check

## Independence and Scope

This is a fresh session. It did not plan or build O-015. The executor
sessions are not named in T-045 or T-046, so `executed_session` records the
build as a whole. Full evidence mode. Both owned Tasks, T-045 and T-046, are
`done` with owner Task-check waivers (09:11:53Z and 09:26:33Z), so this Check
reviews both from the code and tests.

The working tree also holds uncommitted O-024 changes (`objectives.go`,
`plain.go`, `project.go`, and the done-row guards in `update.go`). They are
out of scope here. C-929 covers them.

## Scope Lock

1. Criteria: the ten O-015 Success Conditions (SC1–SC10), the T-045 and T-046
   Done When lists, and guardrails FS-01, FS-04, DATA-01..03, ARCH-02,
   TPL-01, TPL-02, TEST-01..04, TEST-08.
2. Entry points: `data.AdvanceIssueV2`, `data.RetreatIssueV2`,
   `writeIssueTransitionV2` (via `writeV2Record`), `handleIssuesKey`,
   `writeIssueTransitionCmd`, the `actionMsg` handler and reload restore
   (`clampIssueCursor`), `hints()`, `renderHelp`, and the guidance files
   above plus their scaffold copies.
3. Runtime: `freshV2Index` → transition → `writeV2Record` freshness check →
   `loadCmd` reload → `refreshIssues`.
4. Matrix axes: status {open, in_progress, resolved} × direction {advance,
   retreat}; resolved disposition {accepted, verified, escalated} on reopen;
   selection {selected, empty column, filtered empty, no SelectedID, detail
   open, help open}; actor {owner, owner with no session, checker}; time
   {set, zero}; file state {fresh, changed on disk, unwritable}; history
   shape {block list, flow `[]`, missing}; key event {KeyRunes " ",
   KeySpace, KeyBackspace}; guidance {live, scaffold}.
5. Materiality boundary: supported board and data paths on a valid V2
   project. Hand-malformed Issue files are out of scope.

## Coverage Matrix

| Cell | Evidence | Result |
| --- | --- | --- |
| open → in_progress, `owner_decision` | repo tests; probe on real I-030 and I-031 | Pass |
| in_progress → resolved, accepted, owner, time, fixed reason, no `check` | repo tests; probe output on real I-030 | Pass |
| resolved → in_progress, `reopened` naming the disposition, closure fields removed | probe on real I-001 (accepted), I-018 (verified), I-023 (escalated) | Pass |
| in_progress → open | repo tests; probe | Pass |
| Space on resolved, Backspace on open, unknown ID, nil index → named error, no write | repo tests; probe byte compare | Pass |
| non-owner actor, owner with empty session, zero time → refused, no write | probe | Pass |
| History append-only (decoded) | repo tests; probe compares every prior entry's time and note | Pass |
| History append-only (bytes, T-045) | live I-030 diff; probe | **Issue I-057** |
| Body and unknown frontmatter kept | repo tests; probe body compare on six real Issues | Pass |
| Validate before replace; strict index reload after every write | `writeV2Record` callback; probe strict-loads the whole real index after each step | Pass |
| Stale file refused, external bytes kept | repo test; probe on real I-037 | Pass |
| Unwritable file → status line | `TestIssuesWriteFailureAppearsInStatusLine` | Pass |
| Same index reused for two writes | probe (source refreshed; second write succeeds) | Pass |
| History missing / flow `[]` | probe: both append and strict-load | Pass |
| Empty filtered column, no SelectedID, detail open → no command | `TestIssuesTransitionDoesNothingWithoutASelectionOrInDetail` | Pass |
| Help open → Space does nothing | probe | Pass |
| Real `tea.KeySpace` event → transition; focus follows into new column | probe; `TestIssuesSpaceAndBackspaceMoveSelectedIssueAndRetainFocus` | Pass |
| Footer lists both keys | `TestFooterOmitsDeadTaskIssuesHint`; `view.go:385` | Pass |
| Help lists both keys (labels) | probe matched both help labels | Pass |
| Gates unchanged (Check, Task, Objective, Goal) | no gate code in diff; only `accepted` resolution is written, and it names no Check | Pass |
| Live guidance: AGENTS.md, skills, issue-capture, Design §1/§8 | read | Pass |
| Scaffold `agent-skills/` copies byte-identical (TPL-01) | `cmp` on all six files | Pass |
| Scaffold `templates/project-v2/AGENTS.md` | read; no diff | **Issue I-056** |

Not applicable: text-width classes (fixed ASCII copy only); network and
subprocess boundaries (none).

## Workflow and Side Effects

| Order | Operation | Side effect | Failure → final state |
| --- | --- | --- | --- |
| 1 | `freshV2Index` | none | load error → status line, no write |
| 2 | transition guard (status, actor, time) | none | named error → status line, no write |
| 3 | build patches, encode nodes | none | encode error → no write |
| 4 | `writeV2Record`: freshness check, validate, replace | file replaced | stale, invalid, or I/O error → file unchanged, status line |
| 5 | in-memory index update | memory only | n/a |
| 6 | `actionMsg` → `SelectedID` set → `loadCmd` | reload | load error → existing board load handling |

Observed: on success the "I-### moved to …" message is cleared by the reload
(probe: status empty after reload). The criterion asks only that failures
show, so this is an observation.

## Success Conditions

| SC | Classification |
| --- | --- |
| SC1 Space/Backspace move one status; focus follows | Proven |
| SC2 No-op cases write nothing | Proven |
| SC3 open ↔ in_progress change only status + one entry | Proven |
| SC4 Resolve as accepted, owner, time, fixed reason, no Check | Proven |
| SC5 Reopen removes resolution/duplicate_of/escalated_to; note names disposition | Proven |
| SC6 One entry per move; existing entries not edited or reordered | Proven in meaning; T-045's byte-identical wording not met (I-057) |
| SC7 Body/unknown kept; validate; stale refusal; failures in status; reload | Proven |
| SC8 Footer and help | Proven |
| SC9 Focused tests; build + test-fast | Proven (full gate also run) |
| SC10 AGENTS.md, skills/references, and scaffold copies | **Issue** (I-056): scaffold AGENTS.md not updated |

## Commands

- `git diff --check`: clean.
- `make build && make test-full`: exit 0, finished about 2026-09-25T09:28Z,
  `go1.26.2 linux/amd64`, HEAD `1691dbc` plus the working tree; linux,
  darwin, and windows cross-builds succeeded.
- Probe harness in a scratch copy of the working tree
  (`internal/probe/probe_test.go`, `internal/board/v2/zz_probe_test.go`):
  every probe listed above passed except the byte-level history comparison.
  The probes were not added to the repository.

## Issues

| Issue | Likelihood | Impact | Materiality | Recommendation |
| --- | --- | --- | --- | --- |
| I-056 Scaffold AGENTS.md omits board Issue keys | High (every new project) | Low (agents in new projects are told only a Check closes Issues) | Low | Fix now: four short passages, no code |
| I-057 Issue moves re-indent the existing frontmatter | High (first move on any 2-space Issue) | Low (whitespace only; git diff noise) | Low | Owner accepts as is, or a narrow writer fix with a byte test |

## Observations (non-blocking)

- The success status message does not survive the reload. It matches how
  the other board actions handle reloads.
- New history entries are block-style with double-quoted times. Older ones
  are flow-style with single quotes. Both decode the same.
- The live I-030 carries two owner board entries from 09:26:40Z and 09:26:42Z
  (open → in_progress → open). That looks like owner testing. The Issue is
  valid and still open.

## Owner Validation

T-046 declares `owner_validation.required: true`. When a later Check on
O-015 is CLEAR, the owner's acceptance must name that Check.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [x] STYLE-03 **Test branches**
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [x] STYLE-07 **One source of truth**: the reason string and board session are each defined once
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data**
- [x] STYLE-10 **Small diffs**
