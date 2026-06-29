# Audit Findings

Durable, one-file-per-finding records for this repo. Each file is the full record behind a
row in `../register.md`. Findings persist across audit runs so the register can converge
instead of restarting from a cold scan — see `../prompt.md` for the reconciliation rules.

## Naming

- One finding per file, named `F###-slug.md` (e.g. `F001-token-leak-in-logs.md`).
- The `F###` number is the **stable ID**. Once assigned it never changes and is never
  reused, even after the finding is `verified`, `waived`, or `duplicate`.

## Stable IDs and duplicates

- New findings take the next unused `F###`. Mint a new ID only for a genuinely net-new
  finding; an existing finding seen again keeps its ID and refreshes `last_seen`.
- When two records describe the same issue, keep the lower-numbered ID as canonical, set
  the other to `status: duplicate`, and point it at the canonical ID via `duplicate_of`.

## Required frontmatter

```yaml
---
id: F###                       # stable finding ID, matches the filename
title: Short specific summary
status: open                   # lifecycle value, see below
severity: medium               # critical | high | medium | low
confidence: medium             # high | medium | low
source_auditor: agent-or-person
work_item: E00-example/T000    # linked release/epic/task/defect, or omit
guardrail_ids: []              # implicated code-style/guardrail rules
locations: []                  # file paths with line ranges, or surface under review
first_seen: 0000-00-00
last_seen: 0000-00-00
proof_needed: What evidence moves this to verified
duplicate_of:                  # canonical F### when status is duplicate
deferral_reason:               # required for deferred / owner_decision / waived
---
```

### Lifecycle status values

`open`, `triaged`, `mapped`, `in_progress`, `fixed`, `verified`, `deferred`,
`owner_decision`, `waived`, `duplicate`.

## Required body sections

- **## Summary** — what the finding is and why it matters.
- **## Evidence** — concrete locations and observations that show the issue is real.
- **## Proof** — the named evidence required to reach `verified` (and, once verified, the
  passing regression test or explicit manual verification note that satisfied it).
- **## History** — dated reconciliation notes appended each audit run that touched it.

## Proof requirement

A finding reaches `verified` only with named proof — preferably a passing regression test,
otherwise an explicit manual verification note. Never set `status: verified` while
`proof_needed` is unmet or the `## Proof` section is empty.
