---
id: I-145
title: CLAUDE.md with a lone BEGIN marker loses user text on the second upgrade
type: defect
status: open
source:
  kind: check
  check: C-974
  actor: {role: checker, session: check-o043-20261010}
  at: '2026-10-10T07:52:00Z'
tasks: [T-115]
checks: [C-974]
guardrail_ids: [FS-01, FS-02, TEST-03]
severity: medium
history:
  - at: '2026-10-10T07:52:00Z'
    actor: {role: checker, session: check-o043-20261010}
    kind: observed
    note: Found by the Full Objective Check of O-043.
    check: C-974
---

# I-145: CLAUDE.md with a lone BEGIN marker loses user text on the second upgrade

## Summary

`mergeClaudeGuide` (`internal/init/agents.go:100`) treats a `CLAUDE.md` without a complete marker pair as unmarked and appends a fresh managed block. When the file already holds a lone `<!-- SAVEPOINT:BEGIN -->` (for example a user deleted the `END` line), the first run appends a second BEGIN plus a block. On the next run `replaceManagedBlock` (`internal/init/agents.go:69`) pairs the user's orphan BEGIN with the appended END and replaces everything between them, deleting the user's text. Applies to `upgrade-assets` (`upgradeClaudeGuide`, `internal/init/upgrade.go:455`) and to `init` (`MergeClaudeGuide`), which share the merge.

Violates FS-01 and FS-02 (user bytes outside a managed region must stay byte-identical) and O-043 Success Condition 5 ("never overwrites a user's own `CLAUDE.md` content"). The AGENTS.md upgrade path is not affected, because it reports an unmarked file as a conflict instead of appending.

## Evidence

Reproduced with a virtual test file (`go test -overlay`, no tree change) calling `upgradeAssetsFromTree` twice on `CLAUDE.md` = `"# Mine\n<!-- SAVEPOINT:BEGIN -->\nKEEP ME: user notes\n"`:

- Run 1: action `merged`, `KEEP ME` present, file now has two BEGIN markers.
- Run 2: action `merged`, `KEEP ME` gone. Result is `# Mine` followed by only the managed block.

Expected: the user line survives every run, and repeat runs settle as `unchanged` (FS-04).
Missing test: `internal/init/claude_guide_test.go` covers no-marker, full-pair, already-imported, missing, and unreadable files, but not a half marker pair (lone BEGIN, lone END).

## Proof Needed

- A half marker pair in `CLAUDE.md` never causes deletion of user bytes on init or any number of upgrades; a reasonable fix is to report it as a conflict (as AGENTS.md does) or to refuse to pair markers that Savepoint did not write.
- Tests for lone BEGIN and lone END through both init and two consecutive upgrades, asserting user bytes are kept and the second run is `unchanged` or a stable conflict.
- `make build && make test-fast` passes; a re-check confirms.
