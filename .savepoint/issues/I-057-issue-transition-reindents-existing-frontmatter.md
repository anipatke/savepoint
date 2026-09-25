---
id: I-057
title: Board Issue moves re-indent the existing frontmatter, so old history is not byte-identical
type: defect
status: resolved
source:
  kind: check
  check: C-930
  actor: {role: checker, session: o015-objective-check-20260925}
  at: '2026-09-25T09:30:30Z'
tasks: [T-045]
checks: [C-930, C-931, C-932]
severity: low
resolution:
  disposition: verified
  check: C-932
  actor: {role: checker, session: o015-surgical-recheck-20260925}
  at: '2026-09-25T10:14:41Z'
  reason: An independent harness ran 171 real-Issue writes (57 Issues x 3 moves, including the live I-030 second-move reproduction and the I-031/I-037 wrapped-note reproduction) and 90 real-Task writes through the shared writer with 0 byte violations against a from-scratch oracle, and the Full Objective gate (make build && make test-full) passed.
history:
  - at: '2026-09-25T09:30:30Z'
    actor: {role: checker, session: o015-objective-check-20260925}
    kind: observed
    check: C-930
    note: Each board Issue move re-serializes the whole frontmatter at 4-space indent, so existing history entries change in whitespace though not in meaning. T-045 promises they stay byte-identical.
  - at: '2026-09-25T09:40:04Z'
    actor: {role: executor, session: o015-issue-repair-20260925}
    kind: repair_attempted
    note: Shared V2 writes now serialize using indentation inferred from retained YAML source nodes. Added a byte-level regression test for a 2-space Issue history block; focused tests and make build && make test-fast passed. Issue remains open for independent verification.
  - at: '2026-09-25T09:58:00Z'
    actor: {role: checker, session: o015-objective-recheck-20260925}
    kind: rechecked
    check: C-931
    note: Partly repaired. The first move on a 2-space or 4-space Issue now keeps existing bytes, but two cases still fail. The indent inference counts block actor maps inside sequence items as width 2, and a tie picks the smaller width, so the live I-030 and the second move on other 4-space Issues re-indent the whole frontmatter. Wrapped plain history notes (I-031, I-037) are reflowed onto one line.
  - at: '2026-09-25T10:01:54Z'
    actor: {role: executor, session: o015-i057-repair-20260925}
    kind: repair_attempted
    note: Owner-directed targeted repair after C-931. Managed V2 writes now keep the source text of every top-level field no patch touches, and an appended history entry is added after the existing lines at the sequence's own indent. If the spliced text does not decode to the patched document, the writer falls back to a full encode. Indent inference now counts only top-level fields. A new multi-move byte test covers a 4-space file with a block-style entry and a wrapped note. make build && make test-full passed. Issue remains open for independent verification.
  - at: '2026-09-25T10:14:41Z'
    actor: {role: checker, session: o015-surgical-recheck-20260925}
    kind: rechecked
    check: C-932
    note: Both C-931 residual reproductions are clean. Live I-030's further moves kept every earlier line byte-identical; I-031 and I-037's wrapped notes stayed folded and unchanged. An independent from-scratch oracle found 0 violations across 171 real Issue writes and 90 real Task writes, including the resolution/escalated_to removed-field case (I-023, I-024) and the same-index-reused-with-no-reload case. C-932 is CLEAR, so this Issue is resolved as verified.
---

# I-057: Board Issue moves re-indent the existing frontmatter

## Summary

T-045 Done When: "Each transition appends exactly one history entry ... after
the existing entries, which stay byte-identical." They stay identical in
meaning, but not in bytes. `writeV2Record` re-encodes the patched frontmatter
with the YAML library's default 4-space indent, so an Issue written with
2-space indent (every Issue this repo has) has every nested line of
`source`, `resolution`, and `history` re-indented on its first board move.
New entries are also written block-style with double-quoted times, while the
existing ones use flow-style actors and single quotes.

## Evidence

- Live repo: `git diff .savepoint/issues/I-030-windows-full-suite-has-platform-failures.md`
  shows every `source:` and prior `history:` line removed and re-added after
  two owner board moves at 2026-09-25T09:26:40Z and 09:26:42Z.
  `git diff --no-index -w` of the same file against `HEAD` shows only the 12
  added lines, so the change to existing content is whitespace only.
- Checker probe (repo copy): the same happened on I-001, I-018, I-023, I-030,
  I-031, and I-037. All decoded history entries matched before and after.
- `internal/data/issue_v2_test.go` compares decoded entries, not bytes, so
  the test does not cover the byte-identical promise.
- This is the shared behavior of `writeV2Record`. Task and Objective board
  writes re-indent the same way (for example, T-045's own `owner_validation`
  block).

## Repair Attempt

`writeV2Record` now serializes with the indentation inferred from the retained
source YAML nodes, preserving the existing file's nested block indentation.
`TestAdvanceIssueV2_preservesExistingHistoryBytes` verifies that the original
2-space history block remains byte-identical after an Issue transition. The
focused regression tests and `make build && make test-fast` passed. This Issue
remains open for independent verification; no resolution disposition is
recorded.

### Targeted Repair After C-931

C-931 found two remaining cases: 4-space Issues re-indented from their
second board move, and wrapped plain notes reflowed onto one line. Both came
from re-encoding the whole document. The owner ordered one more repair.

- `V2SourceDocument` keeps the frontmatter text it was parsed from.
  `spliceV2Frontmatter` (`internal/data/write_splice.go`) copies the source
  lines of every top-level field no patch names. A patched block sequence
  that only gained items keeps its lines and gets the new items at its own
  indent. Every other patched or new field is encoded on its own.
- `writeV2Record` uses the spliced text only when it decodes to the same
  document as the patched node. Otherwise it encodes the whole document as
  before. After a write it re-parses the written text, so a second write in
  the same process splices against correct line positions.
- `v2FrontmatterIndent` counts only top-level fields, so board-written block
  entries in a 4-space file no longer tip the choice to 2.
- `TestIssueTransitionsKeepEarlierFrontmatterBytes` runs open → in_progress
  → resolved → reopen → open, with one reload from disk. It checks that all
  earlier frontmatter lines survive byte-for-byte (apart from the status
  line and a removed resolution block). It fails with the splice disabled.
- Probe on a scratch copy: 3 moves on each of the 57 real Issues and one
  write to every Task. All 216 writes changed only `status:`/`stage:`,
  appended lines, and the resolution block.
- `make build && make test-full` passed (exit 0, including cross-builds).

## Consequence

No data is lost. The cost is git noise: the first board move on an Issue
shows the whole frontmatter as changed, which hides the one appended entry
in review.

## Proof Needed

Pick one:

- Owner accepts it as is. The Objective's "never edits existing entries" is
  met in meaning, and `writeV2Record`'s formatting is shared with every other
  board write.
- A narrow repair keeps untouched frontmatter lines as they were (or matches
  the file's indent), plus a byte-level test on a 2-space Issue. A recheck
  confirms it.
