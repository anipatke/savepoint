# Prompt: Audit Reconciliation

<!-- AGENT: Produce one epic-local audit findings file. Prefer delta-only edits. Do not write separate proposal files. -->

You are reconciling an epic after its last task is done. The epic is `status: audit-pending`.

**Context gate:** If you built this epic in the current session, you **must not** audit it. Close this session and let the user start a fresh one.

## Input

- `.savepoint/releases/{release}/epics/{E##-slug}/E##-Detail.md` — the epic design (may have deltas from the original plan).
- `.savepoint/releases/{release}/epics/{E##-slug}/tasks/*.md` — task plans, ACs, and drift notes.
- `.savepoint/Design.md` — project architecture.
- `AGENTS.md` — agent guide and codebase map.
- Source and test files in the epic scope.

## Output

One file: `.savepoint/releases/{release}/epics/{E##-slug}/E##-Audit.md`

The audit file must contain these sections in order:

1. **Main Findings** — user-facing summary only: AC verification, important drift, and notable risks. Do not include per-file replacement blocks here.
2. **Code Style Review** — checklist against the 10 AGENTS.md code style rules.
3. **Proposed Changes** — admin/apply metadata only. This section may contain `### Target File`, `### Replace`, and `### With` blocks.

The TUI Audit tab renders only `## Main Findings` and `## Code Style Review`. Keep all file-specific proposal blocks under `## Proposed Changes` so the panel does not show stale admin details.

## File Format

```md
---
type: audit-findings
audited: YYYY-MM-DD
---
# Audit Findings: E## Epic Name

## Main Findings

## Code Style Review

## Proposed Changes
```

Use this shape for every proposed change inside `## Proposed Changes`:

````md
### Target File
path/to/file.md

### Replace
```
<exact old heading, marker, or section>
```

### With
```
<new content>
```
````

## Rules

1. **One audit file only.** Do not create `proposal-1.md`, `proposal-2.md`, `.savepoint/audit/...`, or any other separate audit files.
2. **Prefer delta-only edits.** Use `## Replace` anchored to exact existing text. Do not quote and replace entire large sections unless the whole section genuinely changed.
3. **Do not apply changes yourself.** Write the proposals; the user reviews them in the TUI before commit.
4. **Track context.** Count only intentional audit context reads, and keep notes short.
5. **Stop after proposals.** Do not update the router, do not mark the epic audited, and do not commit.
6. **Be honest about deviations.** If the implementation diverged from the design, document why in the epic-E##-Detail.md section.
7. **Never stop at chat-only findings.** Persist `.savepoint/releases/{release}/epics/{E##-slug}/E##-Audit.md` before you report the audit result back to the user.
8. **Carry-forwards must be actionable.** If an item is "Must Fix Before Next Epic," explain what the next epic's agent should verify. Vague carry-forwards are rejected.
9. **No CLI audit pipeline.** Savepoint audit is agent-led and skill-driven; do not invoke or design around a `savepoint audit` command.
