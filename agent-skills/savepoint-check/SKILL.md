---
name: savepoint-check
description: Runs an independent, fresh-session Check on an explicitly requested Task or a mandatory Objective when router state is check, applying the shared check method to write an immutable Check record and any Issues, and closing work only under the recorded conditional-acceptance rules.
---

# Savepoint Skill: Check

## Purpose

Write an independent, immutable verdict for an optional requested Task Check or mandatory Full Objective Check. Remediation goes to the planner or executor and is verified in a later run; this role cannot repair work to manufacture clearance.

## Goal Context

Every Savepoint project has at least one live Goal selected by the router, and every live Objective names exactly one Goal through `release:`. Goals come from `savepoint init` (G-001), `savepoint migrate`, and the planner, never from this skill. If Next says `Choose a Goal`, or `savepoint doctor` reports a missing Goal or `release:`, report it to the owner; do not pick or create a Goal yourself.

## Trigger

Use this skill when the `Next` line starts with `Check`, or router `state` is `check`, for an explicitly requested
Task Check or for a mandatory Objective Check. Also use it when the `Next` line starts with `Assess`, which is the applicability step in Closure Rules and writes no Check. A pasted `Next` line is the selection.

This session must be fresh: independent from the executor's conversation that built the work under review. The same model is allowed; the same session is not. Model names are optional; names or owner/planner self-report do not establish independence.

## Next

Start from the `Next` line as AGENTS.md's Workflow describes; if `savepoint` is unavailable, follow AGENTS.md rather than guessing. AGENTS.md's Router Selection section says who changes the router. This Check workflow uses `savepoint resume` only to resolve Next and validate project loading.

## Read

- `.savepoint/router.md`
- The requested Task or Objective under review; for an Objective Check,
  every Task it owns
- Applicable policy: the `.savepoint/Guardrails.md` rule IDs the scope names, when the project has that file
- Prior Checks and Issues linked to the scope
- The scoped source and test files the work actually changed
- `agent-skills/references/check-method.md` in full

Load `agent-skills/references/check-method.md` in full and apply its selected mode; this skill does not restate that method.

## Workflow

1. Confirm the session is fresh. If this session built the work under review, stop and hand off to a fresh session. An owner-requested self-review may provide observations, but never writes a Check record or satisfies the independent Check gate.
2. Confirm the scope: a Task Check evaluates one Task's outcome and evidence (the Task Check itself is optional); an Objective Check does everything a Task Check does, plus integration across the Objective's owned Tasks and reconciliation against Design.
3. Apply `agent-skills/references/check-method.md` in full at the matching evidence mode — Quick for a requested Task Check, Full for the mandatory Objective Check.
4. Decide the result. Write one new, immutable Check record — never edit a prior one. A rerun gets a new `C-###` and names the run it replaces in `supersedes`. After writing the record, run `savepoint resume` to strict-load the complete V2 index, including the new Check.
5. On `NEEDS WORK`: record the Issues found, and hand remediation back to the executor or planner rather than repairing anything here. A Task Check's `NEEDS WORK` resumes the executor at `stage: build` inside that same Task. An Objective Check's `NEEDS WORK` must not retreat a Task that is already `done`; remediation is a direct repair under the recorded Issue by default, or new or newly selected work linked to the Objective only when the repair needs planning (see `agent-skills/references/issue-capture.md`, Out-Of-Scope Repair); every previously completed Task keeps its status.
6. On `CLEAR`: this alone does not close a Task or Objective. Apply the closure rules below to record whether the owner may complete the Task or accept the Objective outcome.
7. Treat advisory observations, including `STYLE` guardrail rules, as non-blocking. Fill the `## Code Style Review` checklist as the method describes; neither changes the result.
8. Stop. Do not repair implementation, rewrite acceptance criteria, or update Design as part of this run.

Lane groupings and read/write manifests are advisory. See check-method.md's Parallel Planning Advice; ignoring them is never a finding.

## Verification Gates

Follow AGENTS.md's Verification Policy and the selected method mode. Quick is optional and does not replace handoff gates; Full requires current successful full-gate evidence. Apply the shared metadata-only reuse requirements; otherwise run the full gate fresh.

## Code Health Evidence

Only a Full Objective Check runs `savepoint health check O-### [dir]`, after the full gate, and records its official snapshot ID in `health_snapshot`. Apply `check-method.md`, Collect Code Health Evidence, for absent configuration, verdict blockers, optional failures, stale reports, and Issue admission. Manual snapshots never count; `savepoint health setup` and `savepoint health report` are human-only.

A signal between the aim and the watch line in `.savepoint/health/config.json` is an observation, never an Issue and never a reason for `NEEDS WORK`; list what remains in the Check body. The health verdict supports evidence; it never grants clearance.

## Check Artifact Template

Write each Check record with this structure:

```yaml
id: C-###
scope: {kind: task, id: T-001}
result: CLEAR|NEEDS WORK
checked_by: {role: checker, session: review-001}
executed_session: build-001
checked_at: '2026-09-19T00:00:00Z'
health_snapshot: optional-snapshot-id
reviewed:
  base_commit: optional-base-sha
  head_commit: optional-head-sha
  files: []
  dependencies: []
issues: []
supersedes: null
```

`health_snapshot` is optional and set only by a Full Objective Check that ran `savepoint health check`; it holds the official snapshot ID and nothing else.

`reviewed` is optional scope metadata, including for a `CLEAR` Check; omit the
whole block when there is no review scope metadata to record. When the block is
present, `reviewed.files` and `reviewed.dependencies` use real path or
content-hash entries, or explicit `[]` when that category is absent. This
metadata does not establish technical clearance: `CLEAR` depends on the
independent checker recorded on the Check. A CLEAR Check is current on its
own; no separate freshness assessment is needed.
`issues` lists the `I-###` references this run opened. `supersedes` names the
prior `C-###` this run replaces, or stays empty on a first run. The record body
carries outcome coverage, test and command results, negative and boundary
probes, applicable Guardrails, the `## Code Style Review` checklist when
Guardrails defines `STYLE-*` rules, owner validation still needed, and
nonblocking observations.

A recheck never edits the superseded record; its new `C-###` sets `supersedes`, leaving the prior run intact.

A `NEEDS WORK` Check may add `unmet: [<requirement IDs>]` naming the requirements it found unmet; `unmet` is not allowed on `CLEAR`. When present, an exception grants completion only if its `requirements` cover every listed ID, and each uncovered ID is a named blocker. Without the list, rely on your own assessment of what the exception covers.

## Closure Rules

- A checker may complete a technical Task's optional Check — one with no `owner_validation.required` — once its clearance is current and no unexcepted material blocker remains; only the owner may set the Task's `status: done`.
- A Task with no requested Task Check may be owner-closed only when its implementation evidence is complete and an explicit Task-check waiver names the Task, reason, actor, and time. The waiver skips only the optional local Check: it is not technical `CLEAR`, and it does not waive any acceptance criterion, guardrail, or Objective Check. It satisfies a downstream Task dependency that requires `clear` — the owner's own completion decision stands in there — but never one that requires `accepted`, since there is no Check for the owner to have accepted.
- A Task declaring `owner_validation.required` additionally needs the owner's recorded acceptance. Acceptance is recorded as `owner_validation: {required: true, accepted_check: C-###, scope: [<accepted behavior or criterion IDs>], accepted_by: {role: owner, session: <session>}}`; `accepted_check` is the Check it was first given against and stays as provenance, and `scope` is optional.
- An owner exception is recorded as a structured `exception` only on the owner's explicit instruction — never as Issue history or prose alone: `exception: {requirements: [<criterion or rule IDs>], reason: ..., owner: <owner>, recorded_at: '2026-09-19T00:00:00Z', check: C-###}`. The owner's waiver of evidence is such an exception; `check` is the originating Check.
- A decision carries forward. An acceptance or exception keeps applying through later Checks while the behavior, requirements and scope it covers are unchanged. Applicability is recorded on the decision, append-only, in `carried_forward`: `carried_forward: [{check: C-###, applies: true, assessed_by: {role: checker, session: <session>}, assessed_at: '2026-09-19T00:00:00Z', reason: ..., material_change: <required when applies is false>}]`. An owner entry in the same list renews the decision and always applies. Never edit the originating Check or an earlier entry.
- A re-check, or an `Assess` step, must assess every prior owner decision on the Objective: compare each decision's scope with what changed since its last entry, and record one `carried_forward` entry per decision at the latest Check. Unchanged scope applies, with the reason stated; do not ask the owner to renew it. Only a material change to a decision's own scope sets `applies: false`, naming the change, and then the owner is asked to renew that decision alone.
- `Assess` is applicability only: write the `carried_forward` entries and no new Check record. It never creates, renews or infers an owner decision, and never changes the Check's result.
- Waived evidence is unproven, so a Check stays `NEEDS WORK` while a requirement rests on an exception. A carried decision can make work ready to close by exception, which is completion eligibility, never `CLEAR` or current clearance.
- An Objective closes only after every Task it owns is done, the mandatory Objective integration Check is current, and every material Issue linked to that current Check is resolved (including explicit owner acceptance recorded as an Issue resolution), with the same conditional owner-acceptance rule applied at the Objective level. The Full Objective Check reviews every owned Task, including waived Task Checks. An unfinished owned Task is never excused by an Objective-level exception — cross-Task repair goes back through Tasks, and no Objective Check ever closes a Task directly.
- A Goal is complete when every member Objective is complete; no Goal-level clearance or acceptance is required. Goal completion does not mean published or deployed.
- A record lacking sufficient scope or evidence cannot support completion. A freshness assessment is optional: record one only to mark the latest Check `stale` or `unknown` (for example, when code changed after it). That blocks normal completion until a new Check runs.
- A recorded owner exception can grant completion despite an unmet requirement while it applies, but it is reported as completion by exception, never as a `CLEAR` result or as current clearance.

## Issue Capture

Enter Issue capture as an entry from this workflow when a Check finds
something that blocks the verdict. A `NEEDS WORK` Check records the Issues it
finds; see `agent-skills/references/issue-capture.md` for the artifact
template and rules. This skill may close an Issue as
`verified` after verifying Check proof. The owner may resolve an Issue as
`accepted` or reopen a resolved one; these owner actions do not create a
`CLEAR` Check.

## Rules

This skill may write: the Check record, Issues, evaluation metadata, and authorized closure evidence under Closure Rules. It may record readiness for owner closure, but never sets Task/Objective `status: done` or records owner acceptance on the owner's behalf.

Never repair implementation, edit acceptance criteria to match a result, or update Design as a form of remediation. Route correction back to the planner or executor; verify it in a later immutable record. Apply AGENTS.md's lifecycle terminology.
