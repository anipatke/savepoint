---
id: I-011
title: Committed Windows Zone.Identifier artifact is not ignored
type: defect
status: resolved
source:
    kind: migration
    actor:
        role: owner
        session: .savepoint/releases/v2/defects/D012-ignore-zone-identifier-artifacts.md
    at: "2026-09-20T05:03:29Z"
severity: low
resolution:
    disposition: accepted
    actor: {role: owner, session: user-review-20260922}
    at: "2026-09-22T09:11:47Z"
    reason: >-
        Owner accepts the archived Zone.Identifier artifact as a non-blocking
        observation. It is confined to byte-preserved V1 history and does not affect
        V2 runtime behavior; removing it would conflict with archive preservation.
history:
    - at: "2026-09-22T09:11:47Z"
      actor: {role: owner, session: user-review-20260922}
      kind: owner_decision
      note: >-
          Reclassified the archived metadata file as a non-blocking observation and
          accepted it without a repair or technical CLEAR.
---
## Migrated from V1

Relocated verbatim from the V1 defect body at `.savepoint/releases/v2/defects/D012-ignore-zone-identifier-artifacts.md` (release `v2`).

## V1 Body (verbatim)

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

Resolved by explicit owner acceptance after reclassification as a non-blocking
archive observation. This is not a repair or technical `CLEAR`.
