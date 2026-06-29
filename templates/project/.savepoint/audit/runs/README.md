# Audit Runs

Append-only history of audit runs for this repo. Each run is an **immutable** record of one
audit: what was examined, by whom, against which prompt and commit, and what it found. The
mutable `../register.md` is derived from these records plus disposition work — never edit a
run after it is written; correct course in the next run instead. See `../prompt.md` for the
reconciliation rules a run must follow.

## Naming

- One run per file, named `YYYY-MM-DD-label.md` (e.g. `2026-06-29-pre-release-sweep.md`).
- The date is the day the run was performed; the `label` is a short slug describing the run.
- Multiple runs on the same day take distinct labels; never overwrite an existing file.

## Required frontmatter

```yaml
---
date: 0000-00-00               # day the run was performed
auditor: agent-or-person       # who ran the audit (person, or agent name)
model: model-or-na             # model used, or "n/a" for a human run
prompt_version: v1             # changelog version of ../prompt.md used this run
commit: 0000000                # commit SHA the audit was performed against
mode: full                     # full | incremental | targeted
coverage: examined-vs-skipped  # one-line coverage summary, expanded in the body
source_audits: []              # prior runs or external audits reconciled this run
net_new: 0                     # headline counts for this run
reopened: 0
verified: 0
deferred: 0
coverage_gaps: 0
---
```

## Required body sections

- **## Scope** — `mode`, the surfaces targeted, and why.
- **## Coverage** — **Examined** and **Unexamined** surfaces, with reasons for any skips.
  Unexamined surfaces are coverage gaps and feed the register convergence summary.
- **## Findings** — the `F###` IDs touched this run and their disposition (net-new,
  reopened, advanced, verified, deferred, duplicate). Full records live in `../findings/`.
- **## Reconciliation** — how prior findings and `source_audits` were carried forward,
  preserving stable IDs per `../prompt.md`.

## Headline counts

The frontmatter counts (`net_new`, `reopened`, `verified`, `deferred`, `coverage_gaps`)
summarize the run and must match the dispositions in the body. After the run, copy the
reconciled totals into the convergence summary in `../register.md`.
