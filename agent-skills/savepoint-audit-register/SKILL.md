---
name: savepoint-audit-register
description: Runs register-backed Savepoint audit work when .savepoint/audit/ exists, reconciling prior findings against the register, preserving stable F### IDs, recording an immutable run, and requiring named proof before verified.
---

# Savepoint Skill: Audit Register

## Purpose

Run audit work that converges on the existing Audit Register instead of restarting from a cold scan: open a run, reconcile every prior finding, preserve stable finding IDs, and record the reconciled state.

## Trigger

Use this skill whenever `.savepoint/audit/` exists and audit work is starting — an epic audit in `audit-pending`, or a user-requested repo audit. The `savepoint-audit-epic` skill hands off here for register work; it still owns the epic `E##-Audit.md` handoff file. A read-only `savepoint-audit-task` review writes no register records.

## Read (in order)

1. `.savepoint/audit/prompt.md` — canonical audit prompt and reconciliation rules; note its current changelog version
2. `.savepoint/audit/register.md` — the current reconciled state
3. `.savepoint/audit/findings/` — full finding records behind the register rows (`findings/README.md` defines the record shape)
4. `.savepoint/audit/runs/` — run history, newest first, for how the register got here (`runs/README.md` defines the run shape)
5. Work items linked from findings — the release, epic, task, or defect files named in `work_item` fields this run may touch
6. The source and test files under audit

## Run history vs register state

- `runs/` is **immutable history**: one file per audit run, never edited after it is written. Correct mistakes in the next run, not by rewriting an old one.
- `register.md` is **current state**: a mutable index derived from run history plus disposition work. The register says what is true now; runs say what happened when.
- Never treat an old run as current state, and never update the register without recording the run: every audit writes one new run file and an updated register.

## Workflow

1. Complete the reads above.
2. Audit the in-scope surfaces.
3. Reconcile every prior register finding using the classification below. No finding may be silently dropped from the register.
4. Update finding records in `findings/`: refresh `last_seen` and append dated `## History` notes on existing records; add a new `F###-slug.md` file only for a net-new finding.
5. Write one immutable run file at `runs/YYYY-MM-DD-label.md` per `runs/README.md`, including examined and unexamined coverage.
6. Update `register.md`: one row per finding, and convergence-summary counts that match the run's dispositions.
7. For an epic audit, also write the `E##-Audit.md` handoff file per the `savepoint-audit-epic` skill.
8. Stop for user review.

## Reconciliation classification

Classify each finding the run touches as exactly one of:

- **Existing** — a prior finding still present: keep its `F###`, refresh `last_seen`, advance the status only if work actually advanced.
- **New** — genuinely net-new: mint the next unused `F###` and create its finding file.
- **Regression** — a `fixed` or `verified` finding observed again: keep the ID, reopen it, and count it as reopened in the run.
- **Duplicate** — the same issue as another record: keep the lower-numbered ID canonical and point the other at it via `duplicate_of` with `status: duplicate`.
- **Deferred** — the user chose not to act now: `status: deferred` with `deferral_reason` recorded.
- **Waived** — the user accepted the issue as-is: `status: waived` with the reason recorded.
- **Owner-decision** — blocked on a call only the user can make: `status: owner_decision` with the open question recorded.

## Stable finding IDs

- An `F###` ID is permanent: once assigned it never changes and is never reused, even after `verified`, `waived`, or `duplicate`.
- Reuse the existing ID when a finding is seen again; mint a new ID only for a genuinely net-new finding.
- When in doubt whether a finding is new or existing, prefer matching it to the existing record and noting the uncertainty in its `## History`.

## Proof and closure

- `fixed` means the repair landed; `verified` means named proof exists that the repair holds.
- A finding reaches `verified` only with named proof — preferably a passing regression test, otherwise an explicit manual verification note recorded in the finding's `## Proof` section.
- Never set `status: verified` while `proof_needed` is unmet or the `## Proof` section is empty.

## Authority rules

- Register work never changes task lifecycle authority: agents do not mark tasks `status: done` and do not retreat task status.
- `waived` and `owner_decision` are the user's dispositions. An agent may propose them with a recommendation, but records them only after an explicit user decision, with the reason captured.
- Agents do not close coverage gaps by assertion; an unexamined surface stays a recorded gap until a run examines it.

## Rules

- Never edit a run file after it is written.
- Never duplicate a finding that already has a stable `F###` ID.
- Do not write product code during audit work.
- Keep the register, finding files, and run counts consistent — the convergence summary must match the latest run's dispositions.
- Use `state` only for router phase, task `status` only for task lifecycle, and finding lifecycle values only in audit files.
