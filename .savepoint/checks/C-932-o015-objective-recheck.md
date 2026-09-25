---
id: C-932
scope: {kind: objective, id: O-015}
result: CLEAR
checked_by: {role: checker, session: o015-surgical-recheck-20260925}
executed_session: o015-i057-repair-20260925
checked_at: '2026-09-25T10:14:41Z'
reviewed:
  base_commit: 1691dbca4889791d011732424e46354f92a93f31
  files:
    - internal/data/write.go
    - internal/data/write_splice.go
    - internal/data/parser.go
    - internal/data/issue_v2_test.go
    - AGENTS.md
    - templates/project-v2/AGENTS.md
    - .savepoint/issues/I-056-scaffold-agents-omits-board-issue-keys.md
    - .savepoint/issues/I-057-issue-transition-reindents-existing-frontmatter.md
  dependencies: []
issues: []
supersedes: C-931
---

# C-932: O-015 Surgical Objective Recheck

## Independence and Scope

This is a fresh session. It did not plan or build O-015, and it did not write
the I-056 or I-057 repairs (`o015-i057-repair-20260925`). This is a surgical
re-check limited to exactly the cells named in the invoking task, all of them
drawn from C-931's frozen scope lock and admission ledger: History
append-only (bytes, T-045); Task writes through the shared writer; the same
in-memory index reused for two writes; the scaffold `templates/project-v2/AGENTS.md`
passages (I-056); and the gates. It adds no new axes. A supported-path splice
review (CRLF, blank lines, comment lines, a removed-field span) was performed
as an out-of-cell code-review pass, per the invoking instructions, and is
recorded as a non-blocking observation, not a matrix cell.

## Closure Map

| Issue | State after this recheck |
| --- | --- |
| I-056 Scaffold AGENTS.md omits board Issue keys | Repair proven again. This Check is CLEAR, so the Issue is resolved below as `verified`, naming C-932 |
| I-057 Issue moves re-indent the existing frontmatter | Both C-931 residual reproductions (I-030 second move; I-031/I-037 wrapped notes) are now clean. Repair proven across the full matrix. This Check is CLEAR, so the Issue is resolved below as `verified`, naming C-932 |

## Admission Ledger

| Re-check item | Prior claim | Frozen cell (C-930/C-931) | Allowed result |
| --- | --- | --- | --- |
| Live I-030, further moves past the two already on disk (4-space file, board-written block-style `actor:` entries) | C-931 residual reproduction #1: second move re-indents `source`/`history` | History append-only (bytes, T-045) | Pass / Issue |
| I-031, I-037 wrapped plain multi-line `note:` values, first move onward | C-931 residual reproduction #2: folded notes reflowed onto one line | History append-only (bytes, T-045) | Pass / Issue |
| Every other real Issue (57 total; 3 moves each = 171 writes) | I-057 repair claim: "216 of 216 real writes changed only status/appended lines/resolution block" | History append-only (bytes, T-045) | Pass / Issue |
| Task writes through `writeV2Record` (only `status:`/`stage:` change) | I-057 Evidence names this the shared adjacent case | Adjacent case named in I-057 | Pass / Issue |
| Same in-memory index reused for a second write with no reload | I-057 repair: "writer re-parses after a write ... a second write in the same process splices against correct line positions" | Adjacent case named in I-057 | Pass / Issue |
| Scaffold `templates/project-v2/AGENTS.md` four passages | I-056 repair: final wording alignment | Scaffold `templates/project-v2/AGENTS.md` | Pass / Issue |
| Scaffold `agent-skills/` copies byte-identical | TPL-01 | Scaffold `agent-skills/` copies byte-identical (C-930 cell, re-run for currency) | Pass / Issue |
| `git diff --check`, `make build && make test-full` | Required on every recheck | Verification Gates | Pass / Issue |

## Independent Harness

Built and run only in a scratch copy of the working tree
(`tar --exclude=./.git -cf - . | tar -xf - -C <scratch>`), never against the
real `.savepoint/`. A throwaway Go command, `cmd/o015probe` (not part of the
repository; discarded scratch work, never committed), chdirs into the
scratch `.savepoint` directory and calls `data.LoadV2Index(".")`, matching
how `Source.Path` values resolve.

The oracle is a from-scratch line scanner independent of
`internal/data/write_splice.go`: it re-splits each frontmatter into top-level
fields by scanning for a `^key:` line at column 0 (no `yaml.Node`, no indent
inference, no knowledge of the splice/indent code under test), then requires
every field outside an explicit allow-list to be byte-identical and to keep
its relative order across the write. `status` may change value in place;
`resolution`/`duplicate_of`/`escalated_to` may appear or disappear as a
whole block; `history`'s prior lines must be an exact, unmodified prefix of
the new lines; `stage` (Task writes) may be added or removed. Every write is
also independently re-decoded with `DecodeIssueV2`/`DecodeTaskV2` as a
strict-load sanity check.

### Cell: History append-only (bytes, T-045)

All 57 real Issues, 3 moves each (a full open/in_progress/resolved/reopen
cycle in whichever order the Issue's starting status requires), 171 writes
total, single clean run against a freshly extracted scratch copy: **0
violations**.

- **Live I-030** (the exact C-931 reproduction: 4-space indent, two
  board-written block-style `owner_decision` entries already on disk from
  the real session): 3 further moves applied cleanly. Diffed against the
  committed `HEAD` copy, the only changes across all 5 accumulated moves
  (2 pre-existing + 3 from this probe) are the `status:` value and appended
  `history:` entries; every earlier line, including the two previously-added
  block-style entries and the original flow-style ones, is byte-identical
  and correctly indented at the file's own 4-space level. No re-indentation
  occurred on the second, third, fourth, or fifth move.
- **I-031** (`.savepoint/issues/I-031-reload-errors-for-v2-issue-records.md`,
  5 wrapped `note: >-` folded blocks): after a resolve/reopen cycle, all 5
  folded note blocks kept their exact line numbers, wrapping, and bytes;
  diff shows only the `status:` line and 3 appended history entries.
- **I-037** (`.savepoint/issues/I-037-doctor-names-unsupported-schema-version-malformed.md`,
  1 wrapped `note: >-` block): same result — the folded note is untouched;
  diff shows only `status:` and 3 appended entries.
- **I-023, I-024** (the only two real Issues carrying both `resolution:` and
  `escalated_to:`): reopening removed both fields as a whole block with no
  corruption of the following `history:` field — this is the "wrong span
  when a field is removed" case named in the review instructions, and it is
  clean on the only two real fixtures that exercise it.

### Cell: Task writes through the shared writer

All 45 real Tasks, 2 writes each (a status/stage round trip: to
`in_progress`/`build` or `done` and back to the original state), 90 writes
total: **0 violations**. Every write changed only the `status:` and/or
`stage:` lines; `stage:` was independently confirmed to add or remove
cleanly with no other field disturbed, on both directions of the round trip.

### Cell: Same index reused for two writes (no reload)

Every Issue (2–3 writes) and every Task (2 writes) in the above loops was
written from the same in-memory record, with no reload between writes,
across a single probe process. All subsequent writes on the same object
spliced correctly against the previous write's output, confirming the
repair's re-parse-after-write step keeps the retained `frontmatterText` and
`Frontmatter` node positions correct for the next call.

### Out-of-cell review: supported-path splice correctness

Three synthetic fixtures (not real project files; discarded scratch work),
each given one `AdvanceIssueV2` move:

- **CRLF** (2-space indent, all `\r\n`): output kept `\r\n` throughout with
  no mixed line endings, and every untouched line matched exactly.
- **Blank line** between two top-level fields (`type:` and `status:`):
  preserved verbatim; the blank line stayed where it was.
- **Comment line** (`# a hand-written comment above history`) immediately
  before the touched `history:` field: preserved verbatim ahead of the
  appended entries.

All three passed with 0 violations. This is a non-blocking observation, not
a frozen matrix cell — no in-scope Issue is opened from it.

### Scaffold `templates/project-v2/AGENTS.md` (I-056)

Exact-text diff of the four passages named in I-056's evidence:

- `AGENTS.md:104` vs `templates/project-v2/AGENTS.md:101` (owner Space/Backspace
  authority + direct-instruction rule) — identical.
- `AGENTS.md:112` vs `templates/project-v2/AGENTS.md:109` (Issue Capture
  checker/owner/planner resolution sentence) — identical.
- `AGENTS.md:118` vs `templates/project-v2/AGENTS.md:115` (the task-status stop
  rule) — identical.
- `AGENTS.md:123-128` vs `templates/project-v2/AGENTS.md:120-125` (Check
  section's Issues panel wording) — identical.

### Scaffold `agent-skills/` copies (TPL-01)

`cmp -s` on `agent-skills/references/issue-capture.md`,
`agent-skills/savepoint-check/SKILL.md`, `agent-skills/savepoint-design/SKILL.md`,
and `agent-skills/savepoint-task/SKILL.md` against their
`templates/project-v2/agent-skills/...` copies: byte-identical, all four.

## Coverage Matrix (surgical re-run)

| Cell | Evidence | Result |
| --- | --- | --- |
| History append-only (bytes, T-045) — all 57 Issues, 171 writes | independent harness, single clean run, 0 violations | Pass |
| History append-only — live I-030 (C-931 residual repro #1) | before/after diff on the real working-tree file | Pass |
| History append-only — I-031, I-037 wrapped notes (C-931 residual repro #2) | before/after diff on the real working-tree files | Pass |
| History append-only — resolution/escalated_to removal on reopen (I-023, I-024) | before/after diff; no adjacent-field corruption | Pass |
| Task writes through `writeV2Record` (adjacent) | 45 Tasks, 90 writes, 0 violations | Pass |
| Same in-memory index reused for two writes, no reload (adjacent) | every Issue/Task write in the two loops above | Pass |
| Scaffold `templates/project-v2/AGENTS.md` (I-056) | 4-passage exact diff | Pass |
| Scaffold `agent-skills/` copies byte-identical (TPL-01) | `cmp -s`, 4 files | Pass |

## Success Conditions

| SC | Classification |
| --- | --- |
| SC1–SC5, SC7–SC9 | Proven (unchanged from C-930/C-931; not in this surgical scope, no code affecting them changed since C-931) |
| SC6 One entry per move; existing entries not edited or reordered | **Proven**, including in bytes: both C-931 residual reproductions are now clean, and 171 further real-Issue writes plus 90 real-Task writes show 0 violations |
| SC10 AGENTS.md, skills/references, and scaffold copies | Proven (I-056 repair, re-confirmed) |

## Commands

- `git diff --check`: clean (2026-09-25T10:13:50Z).
- `make build && make test-full`: exit 0, 2026-09-25T10:14:35Z–10:14:41Z,
  `go1.26.2 linux/amd64`, HEAD `1691dbc` plus the working tree; linux,
  darwin, and windows cross-builds succeeded.
- Independent probe harness (`cmd/o015probe` in the scratch copy only, never
  committed): 171 Issue writes + 90 Task writes + 3 synthetic splice-review
  fixtures, each followed by an independent byte-level oracle check and a
  strict production-decoder re-load. All passed.

## Issues

No materiality actions required. No new in-scope defect was found; both
carried-forward Issues (I-056, I-057) are resolved below.

## Observations (non-blocking)

- `agent-skills/bubbletea-tui-design/SKILL.md` still differs from its
  scaffold copy, as C-931 noted. It is not tracked as changed in this
  working tree and is not one of the O-015 files. Outside scope; unchanged
  since C-931.
- The supported-path splice review (CRLF, blank line, comment line, and the
  resolution+escalated_to removed-field span) found no defect on the cases
  tried; see "Out-of-cell review" above. This is a review pass, not a
  frozen matrix cell, so it does not expand the Objective's coverage
  requirement — it is recorded for the owner's information only.
- `v2FrontmatterIndent` (`internal/data/write.go:418`) now counts only
  top-level field widths (the sequence-item branch C-931 flagged is gone),
  but its tie-break branch (`count == frequency && width < indent`) still
  has no direct unit test of its own; it is only exercised indirectly
  through the full-cycle byte tests. Advisory only (STYLE-03 below).

## Owner Validation

T-046 declares `owner_validation.required: true`. This Check is CLEAR, but
that alone does not close the Task or the Objective. The owner's acceptance
must name this exact Check, C-932, not an earlier one.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [ ] STYLE-03 **Test branches** — `v2FrontmatterIndent`'s tie-break branch (`internal/data/write.go:418`) has no direct unit test; only indirect coverage through the full move-cycle tests.
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [x] STYLE-07 **One source of truth**
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data**
- [x] STYLE-10 **Small diffs**
