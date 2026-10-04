---
id: I-138
title: Include a routable Start line in session instructions
type: defect
status: resolved
source:
  kind: check
  check: C-965
  actor: {role: checker, session: check-o033-20261004}
  at: '2026-10-04T03:26:27Z'
tasks: [T-105, T-107]
checks: [C-965, C-966]
resolution:
  disposition: verified
  check: C-966
  actor: {role: checker, session: check-o033-recheck-20261004}
  at: '2026-10-04T03:46:04Z'
  reason: Copied T-004 block now includes standalone exact Start selection while shared router still names T-002.
history:
  - at: '2026-10-04T03:26:27Z'
    actor: {role: checker, session: check-o033-20261004}
    kind: observed
    check: C-965
    note: Initial independent Full Objective Check.
  - at: '2026-10-04T03:46:04Z'
    actor: {role: checker, session: check-o033-recheck-20261004}
    kind: rechecked
    check: C-966
    note: Copied T-004 block now includes standalone exact Start selection while shared router still names T-002.
---

# I-138: Include a routable Start line in session instructions

## Summary

Include a routable Start line in session instructions before O-033 clearance.

## Evidence

O-033 SC8 promises complete copyable instructions; T-105 DW3 requires the current Start line. AGENTS Workflow routes a pasted Next line by its first word; without one, savepoint-task resolves the shared router through resume. Real worktree instructions intentionally leave that router untouched.

With router selecting T-002 and a recommended parallel T-004, copy the generated T-004 block. Its only selection-like line is "Use savepoint-task. Start T-004 — Core work (O-001)." No line begins Start or Next: Start. Expected: an explicit standalone generated Start selection for T-004, so the fresh session follows T-004 despite shared router T-002. Actual: the supplied block does not meet the canonical pasted-Next routing contract and relies on informal reinterpretation of prose. No actual AI agent was spawned; this is a deterministic instruction-contract check, not a claim that every model chooses the wrong Task.

`internal/resume/concurrency.go:147` prefixes the Start line with another instruction and appends punctuation. Independent TestO033FreshSessionStartLine compares against resume.NextLine for T-004 (plain title, no controls) and finds no standalone line. Existing tests assert the combined prose verbatim rather than testing router-independent selection.

## Proof Needed

Emit a standalone Start/Next selection line for the intended Task plus the skill instruction separately. Test router T-002 versus copied T-004; name Goal/Objective/Task, scopes and prerequisite/worktree restrictions unchanged. Safe sanitation remains required. Replay eligible group/active, blocked/unknown/singleton and all formatter surfaces from frozen M6. No router write in the worktree.

## Independent Recheck — C-966

Copied T-004 block now includes standalone exact Start selection while shared router still names T-002.

Verified within the original C-965 scope by CLEAR Check C-966 on 19a6807. Original evidence remains above.
