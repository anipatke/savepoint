---
id: C-926
scope: {kind: objective, id: O-020}
result: NEEDS WORK
checked_by: {role: checker, session: o020-objective-check-20260925}
executed_session: unrecorded-external-executor
checked_at: '2026-09-25T06:55:00Z'
reviewed:
  base_commit: 1651018d1625ec0adea2e3b0e1e906179492e000
  files:
    - internal/data/objective_v2.go
    - internal/data/project.go
    - internal/data/write.go
    - internal/doctor/v2_runtime.go
    - internal/board/v2/objectives.go
    - internal/board/v2/plain.go
    - internal/board/v2/update.go
    - internal/board/v2/io.go
    - internal/board/v2/help.go
    - agent-skills/savepoint-design/SKILL.md
    - templates/project-v2/agent-skills/savepoint-design/SKILL.md
    - .savepoint/Design.md
  dependencies: []
issues: [I-054]
supersedes: null
---

# C-926: O-020 Full Objective Check

## Independence and Scope

This session did not build O-020. It started fresh after `/clear`. The work is
uncommitted on branch `v2` on top of `1651018`. All three Tasks
(T-041..T-043) are `done` with owner Task-check waivers. The waivers were
inspected and are not treated as technical CLEAR. Each Task was reviewed
directly against its Done When list.

### Scope lock

1. Criteria: O-020 Outcome, Success Conditions 1–12, and Confirmed Design
   Decisions; Done When of T-041..T-043; guardrails FS-01, FS-04, DATA-01..04,
   TPL-01, ARCH-02..04, TEST-01..04, TEST-08; Design reconciliation.
2. Entry points: `data.DecodeObjectiveV2`, `data.LoadV2Index`
   (`DuplicateObjectiveRanks`), `data.OrderedObjectiveIDsForGoal`,
   `data.WriteObjectiveGroupOrderV2`, `doctor.RunV2Checks`, the V2 board TUI
   (`renderSidebar`, `handleSidebarKey`, `writeObjectiveGroupOrderCmd`,
   `renderHelp`), and `renderPlain`.
3. Relied-on behavior: `writeV2Record` (atomic replace, freshness check,
   frontmatter node patching), board reload with cursor restore by ID,
   `ResolveNext`, and router writes (must be untouched).
4. Matrix axes: priority value {missing, each of 4, quoted, wrong case,
   unknown, null, int, list}; rank {missing, 1, large, 0, negative, float,
   quoted, word, null, overflow, map}; order {4 groups, ranked/unranked,
   ties, legacy, unknown Goal, nil index}; writer {no-op, partial renumber,
   stale member, deleted member, duplicate ID, missing ID, wrong Goal,
   unknown Goal, bad priority, injected mid-write failure, CRLF, flow
   sequences, comments, quoted title}; keys {1–4, current priority, K, J,
   shift+↑, shift+↓, group edges} × focus {sidebar, columns, detail overlay,
   Issues overlay}; output {TUI, NO_COLOR, narrow, non-TTY}.
   Not applicable: a board action for cross-Goal moves (confirmed out of scope);
   external network boundaries (none).
5. An Issue must violate a Success Condition, Done When item, guardrail, or
   the Objective Check's Design-reconciliation duty through a supported path.

## Coverage Results

| Success Condition / criterion | Evidence | Result |
| --- | --- | --- |
| SC1 four priorities; legacy loads Medium in ID order | Probe decode matrix (21 cells); `OrderedObjectiveIDsForGoal` vs independent hand-written oracle over 10 Objectives | Proven |
| SC2 group → rank → ID; invalid/duplicate diagnosed | Oracle order matched; invalid values fail with path, ID, field, and allowed values; two duplicate facts reported deterministically; doctor test names Goal, priority, rank, Objectives | Proven |
| SC3 fields on Objective only; no Goal member list | `project.go` derives from `ReleaseObjectives`; Release file untouched | Proven |
| SC4 keys `1`–`4`, `K`/`J`, shift arrows | Probe: 3 Objectives moved to High in order, then K / shift+↑ / shift+↓ matched oracle lists at each step; current-priority key is a no-op | Proven |
| SC5 cursor/selection follow; persist; reload | Probe cursor stayed on O-003 across every step; existing tests reopen the board | Proven |
| SC6 preserve content; atomic; conflict before write; partial heals | Probe: quoted title, flow `depends_on`, inline comment, CRLF, body all preserved; refusals (duplicate, missing, wrong/unknown Goal, bad priority, nil index, deleted member) wrote nothing; existing stale-member and chmod mid-write tests | Proven |
| SC7 completed/blocked stay in place | Probe render: done O-001 and waiting O-003 kept rank position with `[✓] Check` and `→ WAITS O-002` | Proven |
| SC8 cross-Goal hand-edit | No board action (by decision); writer refuses a non-member; `project_test.go` ranks per Goal with an R-002 member | Proven |
| SC9 TUI and non-TTY agree | Probe compared plain-output line order with sidebar rows after reorder; real `./savepoint board` output grouped | Proven |
| SC10 status line removed; signals legible | Probe render shows no Planned/In Progress/Done line; badges, `▸`/`●` present; existing NO_COLOR and narrow tests passed | Proven |
| SC11 Next unchanged | Existing tests compare router bytes, `resume.NextLine`, status, and cards across reorder | Proven |
| SC12 tests; guidance aligned; gates | `cmp` of both savepoint-design SKILL copies passed (TPL-01); gates below | Proven |
| Order keys ignored off-sidebar | Probe: keys 1/2/4/K/J with columns, detail overlay, and Issues overlay focused changed no file | Proven |
| Design reconciliation | §1 updated. §8 Layout, Keybindings, and Board persistence not updated, though T-042/T-043 Drift Notes promised it | **Issue I-054** |

## Commands

- `make build && make test-full` — exit 0, finished 2026-09-25T06:50:11Z,
  `go1.26.2 linux/amd64`, on the working tree reviewed here. Cross-builds for
  linux, darwin, and windows succeeded.
- `git diff --check` — clean.
- `cmp agent-skills/savepoint-design/SKILL.md templates/project-v2/agent-skills/savepoint-design/SKILL.md` — identical.
- Temporary probe tests in `internal/data` and `internal/board/v2` — all
  passed; files removed after the run.
- `./savepoint board` (non-TTY) and `./savepoint doctor` on this repository —
  grouped list rendered; no rank or priority diagnostics.

## Issues

- **I-054** (drift, low): Design §8 does not describe the priority-grouped
  sidebar, the plain Objective list, the `1`–`4`/`K`/`J` keys, or the
  Objective order write.

## Materiality

| Issue | Likelihood | Impact | Materiality | Recommendation |
| --- | --- | --- | --- | --- |
| I-054 | High (every reader of §8) | Low (the code, help overlay, and Objective record are correct; only the architecture document lags) | Low | Fix now: a planner-only Design edit, then a targeted recheck |

## Owner Validation Still Needed

Each Task declares `owner_validation.required: true` with an empty
`accepted_check`. After a CLEAR recheck, the owner's acceptance should name
that Check. The User Check scenarios (real terminal, NO_COLOR, narrow width)
were covered by tests and rendered probes here, not by an interactive terminal.

## Observations (non-blocking)

- Two presses of `K` delivered before the first write's reload collapse into
  one move, and both report "moved up". The on-disk order stays valid.
  Writes plus reload take milliseconds, so this is unlikely in practice.
- `writeObjectiveGroupOrderCmd` passes the model's live index into the command
  goroutine. The writer updates `Priority`, `Rank`, and `Source` on those
  shared records while the view may read them. Other board writes load a
  fresh index inside the command. The effect is benign today because a reload
  replaces the index immediately afterwards.
- A member file deleted after load returns a plain `stat` error rather than
  `ErrV2SourceConflict`. Nothing is written, and the board shows the error.
- `savepoint doctor` reports `v2-done-without-clearance` for T-041..T-043,
  the same pre-existing behavior as T-030 for waived Tasks with
  `owner_validation.required`. This is not introduced by O-020.
- Design §8 Keybindings already lists V1-era keys (`d`, `D`, `A`, `R`). The
  planner may want to tidy these while fixing I-054. This is outside O-020.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [x] STYLE-03 **Test branches**
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [ ] STYLE-07 **One source of truth** — `internal/board/v2/update.go` `sidebarPriorityKey`, `sidebarRowPriority`, and `objectivePriorityLabel` restate priority strings as literals instead of the `data.ObjectivePriority*` constants.
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data**
- [ ] STYLE-10 **Small diffs** — the uncommitted O-020 change spans 24 files and about 1,600 added lines, mostly tests, in one working tree.
