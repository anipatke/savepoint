# Prompt: Audit Reconciliation

<!-- AGENT: Produce one patch-shaped proposal bundle. Prefer delta-only edits. Do not write separate proposal files. -->

You are reconciling an epic after its last task is done. The epic is `status: audit-pending`.

**Context gate:** If you built this epic in the current session, you **must not** audit it. Close this session and let the user start a fresh one.

## Input

- `.savepoint/audit/{release}/{E##-slug}/snapshot.md` — what changed during the epic. If this file is missing while the audit CLI is still unavailable, create one manual snapshot from the known epic scope once; do not search broadly for replacement inputs.
- `.savepoint/releases/{release}/epics/{E##-slug}/E##-Detail.md` — the epic design (may have deltas from the original plan).
- `.savepoint/Design.md` — project architecture.
- `AGENTS.md` — agent guide and codebase map.
- Only the source and test files listed in the snapshot.

## Output

One file: `.savepoint/audit/{release}/{E##-slug}/proposals.md`

The proposals bundle must contain these sections in order:

1. **Design.md section** — merge only the epic's architectural delta into the project-level `.savepoint/Design.md`.
2. **AGENTS.md section** — refresh the Codebase Map table with new or changed modules; preserve existing rows.
3. **epic-E##-Detail.md section** — add "Implemented as:" notes and deviations from the original plan.
4. **Quality Review section** — semantic-review findings against the project's code style rules.

## Proposal Format

Use this shape for every proposed change:

```md
## Target File

`path/to/file.md`

## Replace

<exact old heading, marker, or section>

## With

<new content>
```

## Quality Review Format

```md
## Must Fix Before Close

## Must Fix Before Next Epic

## Carry Forward

## Already Fixed
```

## Rules

1. **One proposals file only.** Do not create `proposal-1.md`, `proposal-2.md`, or any other separate files.
2. **Prefer delta-only edits.** Use `## Replace` anchored to exact existing text. Do not quote and replace entire large sections unless the whole section genuinely changed.
3. **Do not apply changes yourself.** Write the proposals; the user reviews them in the TUI before commit.
4. **Track context.** Count only intentional audit context reads, and keep notes short.
5. **Stop after proposals.** Do not update the router, do not mark the epic audited, and do not commit.
6. **Be honest about deviations.** If the implementation diverged from the design, document why in the epic-E##-Detail.md section.
7. **Never stop at chat-only findings.** Persist `.savepoint/audit/{release}/{E##-slug}/proposals.md` before you report the audit result back to the user.
8. **Carry-forwards must be actionable.** If an item is "Must Fix Before Next Epic," explain what the next epic's agent should verify. Vague carry-forwards are rejected.
