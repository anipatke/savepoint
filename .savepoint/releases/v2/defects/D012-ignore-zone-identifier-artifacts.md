---
id: v2/D012-ignore-zone-identifier-artifacts
release: v2
status: open
severity: low
title: "Committed Windows Zone.Identifier artifact is not ignored"
---

# D012: Committed Windows Zone.Identifier artifact is not ignored

## Symptom

The repository tracks `.savepoint/Savepoint V2 Refactor Prompt.md:Zone.Identifier`,
and `.gitignore` has no `*:Zone.Identifier` pattern. The file is Windows
download metadata rather than project content.

## Expected Behavior

The committed metadata artifact should be removed, and future Windows
`Zone.Identifier` streams should be ignored by Git.

## Reproduction

1. Run `git ls-files | rg 'Zone\.Identifier$'` from the repository root.
2. Observe the tracked `.savepoint/Savepoint V2 Refactor Prompt.md:Zone.Identifier`
   entry.
3. Inspect `.gitignore` and observe that it lacks `*:Zone.Identifier`.

## Impact

Platform-specific download metadata adds repository noise and can reappear in
future changes unless the ignore rule is present.

## Fix Plan

Confirmed approach: remove the committed `Zone.Identifier` artifact from the
index (`git rm --cached`) and add `*:Zone.Identifier` to `.gitignore`. Verify
the tracked-file list is clean and unrelated ignore rules are unchanged.

## Acceptance Criteria

- [ ] No `Zone.Identifier` artifact remains tracked.
- [ ] `.gitignore` contains `*:Zone.Identifier`.
- [ ] The ignore rule does not hide ordinary project files.

## Resolution Notes

Pending.
