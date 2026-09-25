---
name: savepoint-check
description: Runs an independent, fresh-session Check on an explicitly requested Task, a mandatory Objective, or a mandatory Goal when router state is check, applying the shared check method to write an immutable Check record and any Issues, and closing work only under the recorded conditional-acceptance rules.
---

# Savepoint Skill: Check

## Purpose

Turn recorded evidence into an independent, immutable verdict. A Task Check is an optional local review; the Full Objective Check is mandatory as the higher-level integration gate, and a Goal Check is mandatory for every Goal. This is the only role that can write those verdicts or close Issues as `verified`; the owner may explicitly close an Issue as `accepted`. The checker's authority is bounded on both sides: it can assess a requested Task Check without turning that local result into Objective or Goal completion, and it can never manufacture clearance by repairing the implementation, rewriting acceptance criteria to match what was built, or updating Design as a form of remediation. Correction always goes back to the planner or executor; this skill verifies the repair afterward, in a later run, and never edits its own prior record to do so.

## Goal Context

Every Savepoint project has at least one live Goal selected by the router, and every live Objective names exactly one Goal through `release:`. If the router Goal is missing, Next says `Choose a Goal`; use `g` to select a live Goal or follow `savepoint doctor`'s repair guidance to create one. An Objective missing `release:` remains loadable, but doctor names the Objective and the exact line to add. Fresh projects receive R-001 from `savepoint init`. `savepoint migrate` keeps the V1 router's live Goal and reuses a uniquely identifiable existing live Goal for active work when its selection is missing or unresolvable. If selected work belongs only to a historical Goal, migration creates a live continuation and moves that Objective into it; unresolved release lifecycle decisions remain in the preview.

## Trigger

Use this skill when router `state` is `check` for an explicitly requested
Task Check or for a mandatory Objective/Goal Check.

This session must be fresh: independent from the executor's conversation that built the work under review. The same model is allowed; the same session is not. Model names are optional — nothing here authenticates identity, it only refuses to treat the executor's own session, or the owner's or planner's self-report, as an independent Check.

## Next

If the owner pasted a `Next` line, act on it directly without re-running `savepoint resume`; otherwise run the read-only `savepoint resume` command and act on its `Next` line. This Check workflow uses `savepoint resume` only to resolve Next and validate project loading; Task creation remains planner-only under the narrow exception in AGENTS.md. If the binary is unavailable, follow AGENTS.md: read `.savepoint/router.md`, report the missing tool, and do not guess the next step. For an owner Task closure, use the board's router advance; see AGENTS.md's Router Selection section for the other selection owners and the `release:` rule.

## Read

- `.savepoint/router.md`
- The requested Task, Objective, or Goal under review; for an Objective
  Check, every Task it owns; for a Goal Check, every member Objective and
  its owned Tasks
- Applicable policy: the `.savepoint/Guardrails.md` rule IDs the scope names, when the project has that file
- Prior Checks and Issues linked to the scope
- The scoped source and test files the work actually changed
- `agent-skills/references/check-method.md` in full

Load `agent-skills/references/check-method.md` completely and apply it as written. It owns scope locks, coverage matrices, the adversarial pass, materiality, and re-audit convergence; this skill does not restate that method.

## Workflow

1. Confirm the session is fresh. If this session built the work under review, state that limitation; do not proceed as an independent Check unless the user explicitly asks to continue anyway.
2. Confirm the scope: a Task Check evaluates one Task's outcome and evidence (the Task Check itself is optional); an Objective Check does everything a Task Check does, plus integration across the Objective's owned Tasks and reconciliation against Design; a Goal Check uses the compatibility value `scope.kind: release` to evaluate integration across all member Objectives.
3. Apply `agent-skills/references/check-method.md` in full at the matching evidence mode — Quick for a requested Task Check, Full for the mandatory Objective Check or Goal Check.
4. Decide the result. Write one new, immutable Check record — never edit a prior one. A rerun gets a new `C-###` and names the run it replaces in `supersedes`. After writing the record, run `savepoint resume` to strict-load the complete V2 index, including the new Check.
5. On `NEEDS WORK`: record the Issues found, and hand remediation back to the executor or planner rather than repairing anything here. A Task Check's `NEEDS WORK` resumes the executor at `stage: build` inside that same Task. An Objective or Goal Check's `NEEDS WORK` must not retreat a Task that is already `done`; remediation is a direct repair under the recorded Issue by default, or new or newly selected work linked to the Objective only when the repair needs planning (see `agent-skills/references/issue-capture.md`, Out-Of-Scope Repair); every previously completed Task keeps its status.
6. On `CLEAR`: this alone does not close a Task or Objective. Apply the closure rules below to record whether the owner may complete the Task or accept the Objective/Goal outcome.
7. Record advisory observations as non-blocking, and fill the record's `## Code Style Review` checklist as `check-method.md` describes; neither changes the result.
8. Stop. Do not repair implementation, rewrite acceptance criteria, or update Design as part of this run.

## Verification Gates

A requested Task Check uses Quick evidence and remains optional; it does not replace handoff gates or mandatory integration evidence. Full Objective and Goal Checks use Full evidence and require current successful `make test-full` evidence. In this repository, CI runs `make ci`, which includes the full gate. A recorded full result is reusable only for a metadata-only correction with the original run documented and code, tests, fixtures, dependencies, and gate definitions proven unchanged since that run; otherwise require a fresh full run.

## Write Boundary

This skill may write: the Check record, Issues, evaluation metadata, and
authorized closure evidence when the closure rules below allow it. It may
record that a Task or Objective is ready for owner closure and whether the
owner has accepted a Goal Check, but it never sets Task/Objective status to
`done`; it never records owner acceptance on the owner's behalf.

It must never: repair implementation, edit acceptance criteria to match a result, or update Design as a form of remediation. A correction is always routed back to the planner or executor, and a later Check — a new record, never an edit to this one — verifies that the repair actually landed.

## Check Artifact Template

Write each Check record with this structure:

```yaml
id: C-###
scope: {kind: task|objective|release, id: T-###, O-###, or R-###}
result: CLEAR|NEEDS WORK
checked_by: {role: checker, session: review-001}
executed_session: build-001
checked_at: '2026-09-19T00:00:00Z'
reviewed:
  base_commit: optional-base-sha
  head_commit: optional-head-sha
  files: []
  dependencies: []
issues: []
supersedes: null
```

Goal Checks retain the existing serialized `scope.kind: release` and `R-###`
identity. Do not rename those compatibility fields when writing a Check.

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

Each run writes a new record with a new `C-###`. A recheck never edits the superseded record; it sets its own `supersedes` and leaves the prior run intact as history.

## Closure Rules

- A checker may complete a technical Task's optional Check — one with no `owner_validation.required` — once its clearance is current and no unexcepted material blocker remains; only the owner may set the Task's `status: done`.
- A Task with no requested Task Check may be owner-closed only when its implementation evidence is complete and an explicit Task-check waiver names the Task, reason, actor, and time. The waiver skips only the optional local Check: it is not technical `CLEAR`, and it does not waive any acceptance criterion, guardrail, Objective Check, or Goal Check. It satisfies a downstream Task dependency that requires `clear` — the owner's own completion decision stands in there — but never one that requires `accepted`, since there is no Check for the owner to have accepted.
- A Task declaring `owner_validation.required` additionally needs the owner's recorded acceptance naming this same current Check; acceptance naming a Check a later run has superseded does not count.
- An Objective closes only after every Task it owns is done, the mandatory Objective integration Check is current, and every material Issue linked to that current Check is resolved (including explicit owner acceptance recorded as an Issue resolution), with the same conditional owner-acceptance rule applied at the Objective level. The Full Objective Check reviews every owned Task, including waived Task Checks. An unfinished owned Task is never excused by an Objective-level exception — cross-Task repair goes back through Tasks, and no Objective Check ever closes a Task directly.
- A Goal Check is mandatory for every Goal. It reviews cross-Objective integration for the existing `R-###` identity, reuses ordinary Issues for material findings, and does not invent a parallel audit. Goal `done` requires at least one member Objective, every member Objective complete, current CLEAR integration evidence, resolved or explicitly excepted material Issues, and the owner's acceptance of that exact current Check. The checker never supplies that acceptance, and Goal `done` does not mean published or deployed.
- A record lacking sufficient scope or evidence cannot support completion. A freshness assessment is optional: record one only to mark the latest Check `stale` or `unknown` (for example, when code changed after it). That blocks normal completion until a new Check runs.
- A recorded owner exception can grant completion despite an unmet requirement, but it is reported as completion by exception, never as a `CLEAR` result or as current clearance, and it applies only to the Check it names.

## Issue Capture

Enter Issue capture as an entry from this workflow when a Check finds
something that blocks the verdict. A `NEEDS WORK` Check records the Issues it
finds; see `agent-skills/references/issue-capture.md` for the artifact
template and rules. A failed Goal Check creates or reuses ordinary Issues;
it does not create Goal-only findings. This skill may close an Issue as `verified` after
verifying Check proof. The owner may close an Issue as `accepted` by an explicit decision recorded with reason, actor, and time; that does not create a `CLEAR` Check.

## Rules

- This session must be independent from the executor's conversation under review; state the limitation plainly when it is not, and do not call the result independent unless the user explicitly says to continue anyway.
- Load and apply `agent-skills/references/check-method.md` in full; do not restate its scope-lock, coverage-matrix, adversarial-pass, materiality, or convergence mechanics here.
- Write only the Check record, Issues, evaluation metadata, and authorized closure. Never repair implementation, edit acceptance criteria, or update Design as remediation — route corrections back to the planner or executor.
- Every Check run is a new immutable `C-###` record; a recheck sets `supersedes` and never edits a prior run.
- Apply Quick evidence only for a requested Task Check. Apply Full evidence for the mandatory Objective Check and Goal Check; an Objective Check additionally covers cross-Task integration and Design reconciliation, and a Task-only Check never substitutes for either mandatory integration gate.
- Apply the Goal Check scope when `scope.kind: release`: inspect all member Objectives and their cross-Objective integration, then reuse ordinary Issues and hand owner acceptance back to the owner.
- A `NEEDS WORK` result records Issues and hands remediation to the executor or planner. A Task Check's `NEEDS WORK` resumes the executor at `stage: build` inside that same Task; an Objective or Goal Check's `NEEDS WORK` instead routes remediation to a direct repair under the Issue, or to new or newly selected work linked to the Objective when the repair needs planning, and must never retreat a Task that is already `done`.
- Apply the closure rules above exactly; do not use stale, unknown, or missing clearance to satisfy an invoked Task Check or a mandatory Objective/Goal Check, and do not treat a Task-check waiver or other exception as a `CLEAR` result.
- Treat advisory observations, including `STYLE` guardrail rules, as non-blocking; record them, but never let them change the result on their own.
- Use `state` only for router phase, Task `status` only for Task lifecycle, and `stage` only when the Task is `in_progress`.
