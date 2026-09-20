---
id: E46-agent-workflow-assets/T004-give-every-task-a-title-a-person-can-read
title: Give every task a title a person can read
status: done
objective: Add the Task artifact template with a separate human title and detailed Outcome, and the planner rule that forbids reusing one as the other.
depends_on:
    - E46-agent-workflow-assets/T003-plan-one-objective-at-a-time
complexity_tier: medium
complexity_reason: A template plus rule wording in one skill pair, but it fixes the field contract every later skill and board
---

# T004: Give every task a title a person can read

## Problem

In V1, `objective` doubles as the display name: `internal/data/parser.go` falls back to it when a title is missing, so the board shows a sentence written for an executor to a person trying to see what is going on. V2 separates them. `title` is short, plain, and for the owner. `objective` is an `O###` reference, and the detailed build outcome moves into the body as Outcome.

The failure this prevents is subtle, because it looks like it works: a planner that fills `title` by truncating the Outcome satisfies every structural check and still produces a board nobody can read. So the rule has to be written as a prohibition the planner can follow and a checker can cite — the title is not derived from the Outcome text — and the template has to make both fields required and obviously distinct.

Structural tests can prove the fields exist and differ. Whether a generated title is actually understandable is semantic, and E50 evaluates it with real agents; this task does not claim it from text assertions.

## Context Files

- `agent-skills/savepoint-design/SKILL.md`
- `templates/project/agent-skills/savepoint-design/SKILL.md`
- `internal/init/agent_skills_test.go`
- `internal/data/task_v2.go`
- `internal/data/objective_v2.go`
- `.savepoint/releases/v2/v2-Design.md`

## Acceptance Criteria

- [x] The Task artifact template lives in `agent-skills/savepoint-design/SKILL.md` as a named section, and the shipped copy stays byte-identical.
- [x] The template's frontmatter carries `id`, `title`, `objective: O###`, `status: planned`, `depends_on` entries in `{task: T###, requires: clear}` form, `owner_validation: {required: true|false}`, and `planned_by`, and the field names match what `internal/data/task_v2.go` decodes.
- [x] `title` and the detailed build outcome are separate required fields: `title` is a short plain-English phrase for the owner, and the detailed outcome lives in the body as Outcome.
- [x] `objective` is documented as an `O###` reference to the owning Objective and never as free text, and every Task belongs to exactly one Objective.
- [x] The planner rules state explicitly that `title` must not be the Outcome text, a truncation of it, or a restatement of the technical objective.
- [x] The template body carries Outcome, User Check, Done When, Context Files, Design References, Guardrails, Implementation Plan, Boundaries, Technical Verification, Technical Evidence, and Drift Notes.
- [x] The template states that Context Files name exact paths — no globs, no directory-only entries.
- [x] `internal/init/agent_skills_test.go` asserts that the template declares `title` and `objective` as distinct required fields, that the no-reuse rule text is present, and that the body sections above all appear.
- [x] The task file written from this template is readable end to end without opening another document to learn what the Task is for.
- [x] A note in the skill records that title readability is evaluated by agent scenarios in E50, not asserted by these tests.

## Implementation Plan

- [x] Read design section 13's complete Task example and `internal/data/task_v2.go` to confirm the exact frontmatter field names and shapes V2 already decodes.
- [x] Add the Task artifact template section to `agent-skills/savepoint-design/SKILL.md`, with frontmatter and body sections filled as a worked example rather than an empty skeleton.
- [x] Add the title rules: required, short, plain, owner-facing, and never derived from the Outcome or objective text.
- [x] Add the Context Files exactness rule and the one-Objective membership rule.
- [x] Mirror the file to the template tree.
- [x] Extend `internal/init/agent_skills_test.go` with the template field and section assertions, and the no-reuse rule assertion.
- [x] Run `go test ./internal/init/... ./internal/data/...`, then `make build && make test`.

## Context Log

**Files read:** `agent-skills/savepoint-design/SKILL.md`, `templates/project/agent-skills/savepoint-design/SKILL.md`, `internal/init/agent_skills_test.go`, `internal/init/skill_validation_test.go`, `internal/init/template_freshness_test.go`, `internal/data/task_v2.go`, `internal/data/objective_v2.go`, `.savepoint/releases/v2/v2-Design.md` (section 13), `.savepoint/releases/v2/epics/E46-agent-workflow-assets/E46-Detail.md`, `.savepoint/Guardrails.md`.

**Files edited:**
- `agent-skills/savepoint-design/SKILL.md` — added `## Task Artifact Template` section (worked example based on design section 13's T014/O008 resume Task) between the Objective template and the Design template, plus a `title` no-reuse rule bullet in `## Rules`.
- `templates/project/agent-skills/savepoint-design/SKILL.md` — mirrored byte-for-byte (verified with `diff`).
- `internal/init/agent_skills_test.go` — added `TestSavepointDesignSkillTaskTemplate` (frontmatter fields, `depends_on`/`owner_validation` shapes, all body section headings, Context Files exactness phrase) and `TestSavepointDesignSkillTaskTitleNoReuseRule` (no-reuse rule text, O### reference/one-Objective phrasing, E50 note).

**Quality gates:**
- `go test ./internal/init/... ./internal/data/...` — pass.
- `make build && make test` — pass (all packages ok).
- No `.savepoint/Health-Check.md` in this project, so the Quick check step is skipped; absence is not a finding per TPL-03.
