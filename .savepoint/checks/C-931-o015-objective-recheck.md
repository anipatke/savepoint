---
id: C-931
scope: {kind: objective, id: O-015}
result: NEEDS WORK
checked_by: {role: checker, session: o015-objective-recheck-20260925}
executed_session: o015-issue-repair-20260925
checked_at: '2026-09-25T09:58:00Z'
reviewed:
  base_commit: 1691dbca4889791d011732424e46354f92a93f31
  files:
    - internal/data/write.go
    - internal/data/issue_v2_test.go
    - AGENTS.md
    - templates/project-v2/AGENTS.md
    - .savepoint/issues/I-056-scaffold-agents-omits-board-issue-keys.md
    - .savepoint/issues/I-057-issue-transition-reindents-existing-frontmatter.md
  dependencies: []
issues: []
supersedes: C-930
---

# C-931: O-015 Full Objective Recheck

## Independence and Scope

This is a fresh session. It did not plan or build O-015, and it did not
write the I-056 or I-057 repairs. This is the one full re-check after
remediation. It uses the frozen scope lock and coverage matrix from C-930.
It adds no new axes.

## Closure Map

| Issue | State after this recheck |
| --- | --- |
| I-056 Scaffold AGENTS.md omits board Issue keys | Repair proven. Stays open: `verified` requires a CLEAR Check |
| I-057 Issue moves re-indent the existing frontmatter | **Still open.** Partly repaired; two reproductions remain in its frozen cell |

## Admission Ledger

| Re-check item | Prior claim | Frozen cell (C-930) | Allowed result |
| --- | --- | --- | --- |
| Scaffold AGENTS.md passages | I-056 repair | Scaffold `templates/project-v2/AGENTS.md` | Pass / Issue |
| Scaffold skill copies unchanged | TPL-01 | Scaffold `agent-skills/` copies byte-identical | Pass / Issue |
| Existing history bytes kept on every move | I-057 repair | History append-only (bytes, T-045) | Pass / Issue |
| Original reproduction: live I-030 | I-057 Evidence | History append-only (bytes, T-045) | Pass / Issue |
| Task and Objective writes through the shared writer | I-057 Evidence names them | Adjacent case named in I-057 | Pass / Issue |
| Every other C-930 cell | unchanged code, rerun gate and probe | same cells | Pass / Issue |

## Coverage Matrix (re-run)

| Cell | Evidence | Result |
| --- | --- | --- |
| All C-930 transition, guard, freshness, selection, key, footer, help, and gate cells | `make test-full` exit 0; probe ran 3 moves on each of the 57 real Issues (171 writes) and strict-loaded the whole index after each one, with no errors | Pass |
| History append-only (decoded) | probe: every write strict-loads | Pass |
| History append-only (bytes, T-045) | probe byte diff on 171 Issue writes | **Issue I-057** (5 writes) |
| Scaffold `agent-skills/` copies byte-identical (TPL-01) | `cmp` on the six O-015 files: identical | Pass |
| Scaffold `templates/project-v2/AGENTS.md` | the four passages match `AGENTS.md:104`, `:112`, the stop rule, and `:123-128` | Pass (I-056 proven) |
| Task writes through `writeV2Record` (adjacent) | probe changed the status or stage of every real Task; only `status:`/`stage:` lines changed | Pass |

## I-057 Residual Reproductions

Probe harness: a scratch copy of the working tree. It advances or retreats
each real Issue three times and diffs the file before and after each
write. Expected: the only removed lines are `status:` and, when reopening,
the `resolution:`/`duplicate_of:`/`escalated_to:` block. 211 of 216
Issue and Task writes matched. Five did not:

1. **4-space files re-indent from their second board move.** The repair
   infers indent by counting key-to-value column widths
   (`internal/data/write.go:412`), and a tie goes to the smaller width
   (`:429`). A board-written entry in a 4-space file puts a block
   `actor:` mapping inside a sequence item (`actor:` at column 7, `role:`
   at column 9), so each one counts as width 2. In the live I-030, two
   entries at width 2 tie with `source` and `history` at width 4, and the
   tie picks 2. The next move re-indents every nested line of `source`
   and `history`. This is the exact reproduction I-057 names. The same
   happened on the second move of I-003 and I-006, after the probe's first
   move added a block-style entry.
2. **Folded multi-line notes are reflowed.** On the first move of I-031
   and I-037, plain multi-line history `note:` values (wrapped at about 75
   columns) are re-emitted on one line. Existing entries change in bytes
   but not in meaning.

`TestAdvanceIssueV2_preservesExistingHistoryBytes` passes. Its fixture has
only one nested width, a single-line note, and a single move, so it does
not cover either case.

## Success Conditions

| SC | Classification |
| --- | --- |
| SC1–SC5, SC7–SC9 | Proven (unchanged from C-930; gate and probe re-run) |
| SC6 One entry per move; existing entries not edited or reordered | Proven in meaning; T-045's byte-identical wording still not met (I-057) |
| SC10 AGENTS.md, skills/references, and scaffold copies | Proven (I-056 repair) |

## Commands

- `git diff --check`: clean.
- `make build && make test-full`: exit 0, 2026-09-25T09:56:13Z–09:56:30Z,
  `go1.26.2 linux/amd64`, HEAD `1691dbc` plus the working tree; linux,
  darwin, and windows cross-builds succeeded.
- Probe harness (`internal/probe/probe_test.go` in the scratch copy only):
  171 Issue writes and one write per Task, each followed by a strict index
  load. The probe was not added to the repository.

## Issues

| Issue | Likelihood | Impact | Materiality | Recommendation |
| --- | --- | --- | --- | --- |
| I-057 Issue moves still rewrite some existing frontmatter bytes | Medium (second move on a 4-space Issue; any wrapped plain note) | Low (whitespace and line breaks only; no data change; git diff noise) | Low | Convergence limit: owner accepts I-057 as is, or orders one targeted repair (keep untouched lines verbatim, or make the indent inference ignore sequence-item offsets and keep wrapped scalars), with a byte test over multi-move and wrapped-note fixtures |

## Observations (non-blocking)

- `agent-skills/bubbletea-tui-design/SKILL.md` differs from its scaffold
  copy. It is not tracked as changed in this working tree and is not one
  of the O-015 files. Outside scope.
- I-056's repair history lists three `repair_attempted` entries within
  seven minutes. The final wording is what this recheck compared.

## Owner Validation

T-046 declares `owner_validation.required: true`. When a later Check on
O-015 is CLEAR, the owner's acceptance must name that Check. If the owner
accepts I-057, a later recheck can then close I-056 as `verified` under a
CLEAR result.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [ ] STYLE-03 **Test branches** — the tie-break and sequence-item branches in `v2FrontmatterIndent` (`internal/data/write.go:401`) have no test.
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [x] STYLE-07 **One source of truth**
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data**
- [x] STYLE-10 **Small diffs**
