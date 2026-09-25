---
id: T-048
title: Stop describing a Goal Check in guidance
objective: O-025
status: done
depends_on: [{task: T-047, requires: clear}]
owner_validation:
    required: false
    accepted_check: ""
planned_by: {role: planner, session: o025-plan-20260925}
check_waiver:
    task: T-048
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-09-25T20:58:09Z"
---

# Stop describing a Goal Check in guidance

## Outcome

Every active document and the new-project scaffold describe two Checks: the
mandatory Full Objective Check and the optional Task Check. A Goal is
complete when all its Objectives are complete.

## User Check

Search the active guidance for "Goal Check", "Release Check", and
`scope.kind: release`: only the compatibility note that old
`scope.kind: release` records still load remains. Read AGENTS.md's
Verification Policy and Check sections: they name only the Objective and Task
Checks.

## Done When

- `.savepoint/Design.md` describes Goal completion as all member Objectives
  complete, removes the Goal Check step from the verification order, and
  notes that `scope.kind: release` records still load but no decision reads
  them.
- `.savepoint/Guardrails.md` TEST-08, TEST-09, and the Check summary no
  longer mention a Goal Check.
- `AGENTS.md`, `.savepoint/router.md`, and `README.md` drop the Goal Check.
- `agent-skills/savepoint-check`, `savepoint-design`, `savepoint-task`, and
  the three shared references drop the Goal Check, including savepoint-design's
  Goal `done` rule.
- The `templates/project-v2` copies of every file above match the live
  wording where they are meant to.
- `internal/init/template_freshness_test.go`,
  `internal/init/agent_skills_test.go`, and migration goldens are updated
  only where their expected text changes.
- `make build && make test-fast` pass.

## Context Files

`.savepoint/Design.md`, `.savepoint/Guardrails.md`, `.savepoint/router.md`,
`AGENTS.md`, `README.md`,
`agent-skills/savepoint-check/SKILL.md`,
`agent-skills/savepoint-design/SKILL.md`,
`agent-skills/savepoint-task/SKILL.md`,
`agent-skills/references/check-method.md`,
`agent-skills/references/commands-and-procedures.md`,
`agent-skills/references/issue-capture.md`,
`templates/project-v2/.savepoint/Design.md`,
`templates/project-v2/.savepoint/Guardrails.md`,
`templates/project-v2/.savepoint/router.md`,
`templates/project-v2/AGENTS.md`,
`templates/project-v2/agent-skills/savepoint-check/SKILL.md`,
`templates/project-v2/agent-skills/savepoint-design/SKILL.md`,
`templates/project-v2/agent-skills/savepoint-task/SKILL.md`,
`templates/project-v2/agent-skills/references/check-method.md`,
`templates/project-v2/agent-skills/references/commands-and-procedures.md`,
`templates/project-v2/agent-skills/references/issue-capture.md`,
`internal/init/template_freshness_test.go`,
`internal/init/agent_skills_test.go`,
`.savepoint/objectives/O-025-remove-the-goal-check/Objective.md`.

## Design References

Design Section 1 (Goal boundary, Check workflow), the Check table, and
Section 10 verification order.

## Guardrails

TEST-08, TEST-09, TPL-02, TEST-01..04.

## Implementation Plan

1. Confirm the runtime Task landed; return REPLAN REQUIRED if Goal
   completion still reads a Goal Check.
2. Edit the live guidance files, then mirror the scaffold copies.
3. Update template and skill tests and any migration goldens whose text
   changed.
4. Search for leftover Goal Check wording in active guidance.
5. Run `make build && make test-fast`.

## Boundaries

No runtime code changes, and no edits to historical Checks, Issues,
Objectives, or Tasks.

## Technical Verification

Focused template and skill tests during iteration; `make build &&
make test-fast` at handoff; the Full Objective Check requires
`make test-full`.

## Technical Evidence

Started 2026-09-26 from the owner-supplied router selection `Start T-048`;
the selection confirmed the required T-047 `clear` dependency. The owning
Objective O-025 was already `in_progress`.

Acceptance criteria:

- **Design completion and verification order — proven:** Design now derives
  Goal completion from all member Objectives, lists only Task and Objective
  Checks, and retains one compatibility note that old `scope.kind: release`
  records load but are not read by completion decisions.
- **Guardrails — proven:** TEST-08, TEST-09, and the Check summary describe
  only Task and Objective Checks.
- **AGENTS.md, router, README — proven:** these guides describe only the
  optional Task Check and mandatory Full Objective Check; Goal completion is
  based on member Objectives.
- **Skills and references — proven:** the Check, Design, and Task skills and
  all three shared references no longer describe a Goal Check. The Design
  skill's Goal completion rule derives completion from its Objectives.
- **Scaffold copies — proven:** corresponding scaffold wording was updated;
  skill and reference parity checks pass.
- **Tests and migration goldens — proven:** `agent_skills_test.go` now expects
  Objective-only Check guidance and a valid Task-scoped artifact example.
  `template_freshness_test.go` needed no expectation change; migration
  goldens needed no update.
- **Handoff gate — proven:** final `make build && make test-fast` passed.

Guidance search: no active guidance contains “Goal Check” or “Release Check”;
the sole `scope.kind: release` mention is the Design compatibility note.
AGENTS.md's Verification Policy and Check sections name only Task and
Objective Checks.

Commands and results:

- `make build && make test-fast` — final run passed all packages. Earlier
  iteration runs exposed stale assertions and an invalid placeholder scope in
  the Check example; both were corrected before the passing run.
- `go test ./internal/init -count=1` — an intermediate diagnostic run caught
  the same invalid placeholder scope; the final fast gate passed after it was
  replaced with a valid Task scope.
- `git diff --check` — passed after all content edits.
- Targeted `rg` guidance searches — passed; only the Design compatibility
  note retains `scope.kind: release`.

Files read: `.savepoint/router.md`, this Task, `.savepoint/objectives/O-025-remove-the-goal-check/Objective.md`, `.savepoint/Design.md`, `.savepoint/Guardrails.md`, `AGENTS.md`, `README.md`, the live Check/Design/Task skills, the three live shared references, all listed `templates/project-v2` copies, `internal/init/template_freshness_test.go`, and `internal/init/agent_skills_test.go`. The required `savepoint-task` workflow skill was also read. The five files under `internal/migrate/testdata/golden/` were scanned for changed guidance phrases; none matched.

Files changed: `.savepoint/Design.md`, `.savepoint/Guardrails.md`,
`.savepoint/router.md`, `AGENTS.md`, `README.md`, the live Check/Design/Task
skills and three shared references, the corresponding scaffold files, this
Task's evidence, and `internal/init/agent_skills_test.go`.

Extra reads: `internal/migrate/testdata/golden/v1-basic.yml`,
`v1-history.yml`, `v1-router-archived.yml`, `v1-router-missing.yml`, and
`v1-router-unresolvable.yml` were searched for the changed guidance phrases.
No matches were found, so their expected text did not change. This targeted
read was needed because the implementation plan allows updating migration
goldens when their guidance text changes and those paths are outside the
Context Files.

Limitations: no optional independent Task Check or owner waiver has been
recorded. The Task is at `audit` for the owner's handoff decision.

## Drift Notes

None expected.
