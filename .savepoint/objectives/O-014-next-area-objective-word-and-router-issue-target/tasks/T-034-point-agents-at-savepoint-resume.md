---
id: T-034
title: Point agents at savepoint resume for the next step
objective: O-014
planned_by: {role: planner, session: o014-design-20260924}
status: done
complexity_tier: medium
complexity_reason: "Guidance-only, but spans AGENTS.md, the router template, four live skills with byte-identical scaffold copies, and Design reconciliation for the whole Objective."
depends_on: [{task: T-028, requires: clear}, {task: T-033, requires: clear}]
owner_validation: {required: false}
check_waiver:
    task: T-034
    reason: not needed
    actor: {role: owner, session: owner-chat}
    recorded_at: "2026-09-24T08:48:30Z"
---

# T-034: Point agents at savepoint resume for the next step

## Outcome

Agent guidance names `savepoint resume` as the one CLI command agents may
run and the source of the next action, says who updates the router selection
after a closure, and Design.md describes what O-014 implemented.

## User Check

Paste a Next line such as `In Progress O-014 · Build T-029 — ...` into a
fresh session: the agent opens T-029 with `savepoint-task` without asking
what the line means.

## Done When

- AGENTS.md (live and `templates/project-v2/AGENTS.md`), the router
  template, and the phase skills replace "act on the router's next_action"
  with "run `savepoint resume` and act on its Next". CLI Rules allow exactly
  `savepoint resume` (read-only) and keep every other command human-only.
- AGENTS.md's Workflow starts with the Next line: a pasted or `resume`-printed
  line names the Objective and Task; its word picks the skill (Task
  `Planned`/`Build`/`Test` → `savepoint-task`, Task or Objective `Check` →
  `savepoint-check`, Objective with no Task → `savepoint-design`; an
  `Open` or `In Progress` Issue → repair under `savepoint-task` per
  `issue-capture.md`, leaving closure to a checker or the owner).
- AGENTS.md documents the owner request "set router to O-### [T-###]
  [I-###]": confirm each record exists and the Task belongs to the
  Objective (refuse and say why otherwise), edit only those selection keys,
  leave `release:` unchanged unless the owner names a Release, and finish by
  running `savepoint resume` and showing its Next line.
- When the binary is unavailable, guidance says to read the router
  selection, report the missing tool, and not guess the next step.
- Router-advance ownership is stated once in AGENTS.md and referenced by
  skills: the board advances the router when the owner closes a Task; `savepoint-design`
  selects the next Objective; `savepoint-task` selects the Task it starts;
  nobody blanks or changes `release:` except an explicit Goal choice.
- AGENTS.md (live and template, identical wording) keeps the rule that a
  recorded owner Task-check waiver satisfies a `requires: clear` dependency
  but never `requires: accepted`, and the rule that `savepoint resume` and
  the board's `Blocked:` lines (from `ResolveTaskDependencyV2`) decide
  whether a dependency blocks — agents do not re-derive it from prose.
- Skills stop writing `next_action` (`savepoint-design` step 11,
  `savepoint-idea`). Live and scaffold skills stay byte-identical.
- Design.md section 8 states the Next line format, that Next is the router
  selection, the board's closure advance, the diagnostic line, and Issue
  context; sections 1 and 4 drop `next_action` and outdated `p` hand-off
  wording. Section 4's Issue sentence says Issues carry no stage (repairs
  I-043; record `repair_attempted` in the Issue, leave closure to a checker).
- Template freshness and agent-skill tests pass; `git diff --check`,
  `make build && make test-fast` pass.

## Context Files

`AGENTS.md`, `templates/project-v2/AGENTS.md`,
`templates/project-v2/.savepoint/router.md`, `.savepoint/router.md`,
`agent-skills/savepoint-idea/SKILL.md`,
`agent-skills/savepoint-design/SKILL.md`,
`agent-skills/savepoint-task/SKILL.md`,
`agent-skills/savepoint-check/SKILL.md`,
`templates/project-v2/agent-skills/savepoint-idea/SKILL.md`,
`templates/project-v2/agent-skills/savepoint-design/SKILL.md`,
`templates/project-v2/agent-skills/savepoint-task/SKILL.md`,
`templates/project-v2/agent-skills/savepoint-check/SKILL.md`,
`.savepoint/Design.md`, `internal/init/template_freshness_test.go`,
`internal/init/agent_skills_test.go`, `agent_skills_test.go`.

## Design References

Design sections 1, 4, 6, and 8.

## Guardrails

TPL-01, TPL-02, TPL-04, POL-02, TEST-06, TEST-08.

## Implementation Plan

1. Update AGENTS.md workflow, CLI Rules, and router-advance ownership.
2. Update skills in both trees together; diff to confirm byte identity.
3. Reconcile Design.md sections 1, 4, 6, and 8 against T-028..T-033.
4. Run template and skill tests.

## Boundaries

No production code. No Goal-is-mandatory wording (O-022).

## Technical Verification

`make build && make test-fast` at handoff; the Full Objective Check then
runs a fresh `make test-full`.

## Technical Evidence

Extra read: `.savepoint/Guardrails.md` to check the Task's named rules TPL-01, TPL-02, TPL-04, POL-02, TEST-06, and TEST-08 before editing guidance and tests. These rules are advisory or required as stated in that file.
Extra read: `Makefile` to identify the executable path produced by the required build gate, because `savepoint` was not available on `PATH` and a runtime `resume` check is needed.
Extra read: `internal/buildtool` build-output implementation to locate the host binary from `make build`, so the permitted read-only `savepoint resume` command can verify the selected Task and dependency gate.
Extra read: `internal/data/evidence_v2.go` to record the explicit owner Task-check waiver in the runtime-supported evidence format after the owner's direction.
Extra read: `agent-skills/references/issue-capture.md` and the Issue index to check whether the resume wording conflict is already tracked before recording a follow-up.

### Acceptance evidence

1. **Resume and CLI guidance — met.** Live and scaffolded AGENTS guides, both router files, and all four live/scaffolded phase skills name the read-only `savepoint resume` command as the source of Next. The CLI rules permit only that Savepoint command to agents.
2. **Next routing — met.** AGENTS maps Task `Planned`/`Build`/`Test`, Task or Objective `Check`, Objective with no Task (except its Check rung), and standalone `Open`/`In Progress` Issue lines to their skills; a co-selected Issue remains context.
3. **Owner router request — met.** AGENTS gives the requested `set router to O-### [T-###] [I-###]` form, validates records and Task ownership, limits edits to selection keys, preserves `release:`, and finishes with `savepoint resume`.
4. **Missing binary — met.** AGENTS and skills direct agents to read the router selection, report the missing tool, and not guess the next step. The initial `savepoint resume` invocation confirmed the command was not on `PATH` (exit 127); the guidance still allows a supplied Next line as an explicit selection.
5. **Selection ownership — met.** AGENTS states that the board advances after owner Task closure, `savepoint-design` selects the next Objective, and `savepoint-task` selects the Task it starts; the skills refer to that section. The active router now selects O-014/T-034 and retains `release: R-006`.
6. **Waiver and dependency guidance — met.** Both AGENTS copies retain the same owner-waiver and runtime dependency-gate rules.
7. **Skill field and parity — met.** The idea and design skills no longer write `next_action`; all four live/scaffold pairs compare byte-identically.
8. **Design reconciliation — met.** Design Sections 1, 4, 6, and 8 describe resume-based Next, selection ownership, the agent CLI boundary, line formats, Issue context and stale diagnostics, board closure advancement, and preserved Release selection. Section 4 says Issues carry no `stage` and records `repair_attempted` on the Issue.
9. **Verification — met.** The template-freshness and agent-skill tests pass within the final fast gate; `git diff --check` passes.

### Commands and results

- `savepoint resume` — unavailable before build (exit 127: command not found).
- First `make build && make test-fast` — build succeeded; fast gate failed on three stale guidance assertions. Restored the required Idea Objective-handoff wording, kept the Check freshness rule in its Trigger section, and removed the scaffold guide's stale term.
- Final `make build && make test-fast` — passed (exit 0).
- `./savepoint resume` at start — passed (exit 0); Next was `Planned O-014 · Build T-034 — Point agents at savepoint resume for the next step`; output reported `Ready: It may advance to its next stage.` and no `Blocked:` line.
- `./savepoint resume` after the fast gate — passed (exit 0); Next was `Planned O-014 · Test T-034 — Point agents at savepoint resume for the next step`; output again reported `Ready: It may advance to its next stage.`
- `git diff --check` — passed.
- Four `cmp -s` live/scaffold skill comparisons — passed. A targeted `rg` for `next_action` across those eight skill files returned no matches.

### Files read and changed

Read the router, O-014 Objective, T-034, Guardrails, both AGENTS guides, both router templates, all four live skills and their four scaffold copies, Design.md, `internal/init/template_freshness_test.go`, `internal/init/agent_skills_test.go`, and `agent_skills_test.go`. Extra reads are listed above.

Changed for T-034: both AGENTS guides, both router files, all four live/scaffold skill pairs, `.savepoint/Design.md`, and this Task record. No production code or tests were changed.

### Limitations and handoff

This is an ordinary guidance Task; no `make test-full` run was required. The owner explicitly waived the optional Task Check with reason `not needed` (actor `owner`, session `owner-chat`; recorded `2026-09-24T08:48:30Z`). This waiver is not technical `CLEAR` and routes the evidence to the mandatory Full Objective Check. The Task remains `in_progress` at `stage: audit`; only the owner may mark it `done`.

## Drift Notes

None expected.
