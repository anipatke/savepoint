---
id: T-008
title: Align active guidance with Goals
objective: O-013
planned_by: {role: planner, session: goals-terminology-20260921}
status: done
complexity_tier: medium
complexity_reason: "The active repository guide, architecture, public README, workflow skills, and shipped V2 copies must change together while historical evidence remains byte-preserved."
depends_on: [{task: T-007, requires: clear}]
owner_validation:
    required: false
    accepted_check: ""
check_waiver:
    task: T-008
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-09-23T09:37:49Z"
---

# T-008: Align active guidance with Goals

## Outcome

Active repository guidance and shipped V2 workflow documentation describe the
delivery context as Goals, explain the preserved R###/Release storage boundary,
and remain consistent with the implemented V2 board.

## User Check

Read the reconciled README and repository guide alongside the V2 board wording.
Confirm that a new user learns “Goals” as the public concept, while the
compatibility note makes existing `R###` and `release:` files understandable.
No additional product choice is required for this documentation Task.

## Done When

- `README.md`, `AGENTS.md`, `.savepoint/Design.md`, and
  `.savepoint/Guardrails.md` use Goal/Goals for public V2 terminology and
  retain only necessary compatibility references to Release.
- `templates/project-v2/AGENTS.md` and every active V2 skill/reference copy
  describe the same Goal vocabulary, shortcut, optional boundary, and check
  semantics.
- Documentation does not imply that Goals publish, deploy, tag, or create
  changelogs, and does not turn owner waiver/acceptance into technical
  clearance.
- Canonical active skills and their `templates/project-v2/` copies remain
  byte-identical where the template contract requires it.
- Targeted searches distinguish allowed internal, migration, and historical
  Release references from stale public UI instructions.
- Archived V1 material, immutable Checks, Issues, and current historical
  evidence are not rewritten.

## Context Files

`README.md`, `AGENTS.md`, `.savepoint/Design.md`, `.savepoint/Guardrails.md`, `templates/project-v2/AGENTS.md`, `agent-skills/savepoint-idea/SKILL.md`, `agent-skills/savepoint-design/SKILL.md`, `agent-skills/savepoint-task/SKILL.md`, `agent-skills/savepoint-check/SKILL.md`, `agent-skills/references/check-method.md`, `agent-skills/references/issue-capture.md`, `agent-skills/references/commands-and-procedures.md`, `templates/project-v2/agent-skills/savepoint-idea/SKILL.md`, `templates/project-v2/agent-skills/savepoint-design/SKILL.md`, `templates/project-v2/agent-skills/savepoint-task/SKILL.md`, `templates/project-v2/agent-skills/savepoint-check/SKILL.md`, `templates/project-v2/agent-skills/references/check-method.md`, `templates/project-v2/agent-skills/references/issue-capture.md`, `templates/project-v2/agent-skills/references/commands-and-procedures.md`.

## Design References

Design sections 2, 3, 6, 7, 8, 10, and 13.

## Guardrails

TPL-01, TPL-02, TPL-04, ARCH-04, TEST-01, TEST-06, TEST-08, STYLE-07,
STYLE-09, and STYLE-10.

## Implementation Plan

1. Compare each named active document with the completed board vocabulary and
   the compatibility evidence from T-006.
2. Update public Goal terminology, shortcut guidance, hierarchy descriptions,
   and compatibility notes without changing policy semantics.
3. Apply the same edits to the canonical V2 skills and their shipped copies;
   leave V1/archive documentation untouched.
4. Search for stale public Release wording, verify template parity, and run
   `git diff --check`, `make build`, and `make test`.

## Boundaries

No production behavior, persisted schema, migration archive, immutable Check,
Issue, historical V1, or independent Release/Goal policy rewrite.

## Technical Verification

Exact-path active-document terminology searches, template freshness/parity
checks, focused board tests, `git diff --check`, `make build`, and `make test`.
The mandatory Full Objective Check must verify documentation and board agreement
before O-013 can close.

## Technical Evidence

Start evidence: T-007 is `done` with an explicit Task-check waiver recorded at
2026-09-23T09:13:53Z; its `requires: clear` dependency is satisfied. O-013 is
`in_progress`, so T-008 may begin.

Extra read before start:

- T-007 task record — verify the declared `requires: clear` dependency because
  the router selects T-008 but its `next_action` text still describes T-006.
- T-006 task record — use its completed preservation evidence as the
  compatibility baseline called for by T-008's implementation plan.
- `.savepoint/objectives/O-013-rename-release-to-goals/Objective.md` — check
  O-013 boundaries/readiness before beginning its final Task.
- `.savepoint/Design.md` sections 2, 3, 6, 7, 8, 10, and 13, and the named
  rule rows in `.savepoint/Guardrails.md` — apply T-008's design references and
  documentation/template policies.

Extra reads during verification:

- `internal/init/agent_skills_test.go` — read only the failing Goal/Release
  contract assertions after the first `make test` showed that active V2 skill
  expectations still required the retired public Release wording. Updated
  those assertions to require Goal wording and the preserved serialized scope.
- `internal/init/v2_scaffold_test.go` — read only the failing scaffold-guide
  assertions after the same run; adjusted the active template wording to retain
  the load-bearing absence statement and avoid the legacy `phase` vocabulary.

Criterion evidence and document inventory:

- `README.md`, root `AGENTS.md`, `.savepoint/Design.md`, and
  `.savepoint/Guardrails.md` now use Goal/Goals for public V2 context, board,
  and integration-check terminology. Compatibility notes retain `R-###`,
  `.savepoint/releases/<slug>/Release.md`, `release:` references, internal
  Release resolver names, and `scope.kind: release` without implying a schema
  migration.
- `templates/project-v2/AGENTS.md`, the four active V2 skills, and three shared
  references describe the optional Goal boundary and mandatory Goal Check
  semantics. The board key is `g`; `r` is documented only as a hidden
  compatibility alias. Owner waiver/acceptance remains distinct from technical
  `CLEAR`; Goal completion does not publish, deploy, tag, or create changelogs.
- All seven live workflow skill/reference files compare byte-identically with
  their `templates/project-v2/` copies (`cmp -s`: seven matches).
- Exact-path stale-public-wording search found no remaining Release Check,
  optional Release, or Release-selector instructions. Remaining Release hits
  are compatibility serialization (`release:`, `scope.kind: release`),
  internal resolver/model names, and generic release-gate wording.
- No archived V1 material, immutable Check, Issue, or historical evidence was
  edited. No production behavior or persisted schema was changed. The
  documentation contract assertions in `internal/init/agent_skills_test.go`
  now assert Goal terminology while retaining the serialization contract.

Verification:

- `make build` — passed.
- Initial `make test` — failed only on active-doc assertions that still pinned
  the old Release vocabulary and a scaffold phrase split across lines. Updated
  the V2 documentation assertions and restored the scaffold's exact
  load-bearing absence sentence.
- Focused contract tests via `go test ./internal/init -run
  'TestSavepointCheckSkillScopes|TestSavepointCheckSkillClosureRules|TestV2SkillsTeachOptionalGoalWorkflow|TestV2ScaffoldAgentsGuideIsLiveAndUsesV2Vocabulary|TestV2AgentsGuideCarriesExistingCodebaseAdoptionSection' -count=1`
  — passed. Elevated access was required for Go's build cache; no repository
  files were changed by that access.
- Final `make build` and full `make test` (`./...`) — passed; migration was
  the slowest package at 1m55.641s.
- `git diff --check` — passed.

Handoff: implementation and audit evidence are complete. T-008 remains
`in_progress` at `stage: audit`; awaiting the owner's choice of an independent
optional Task Check or an explicit Task-check waiver. Only the owner may mark
the Task `done`.

## Drift Notes

Any documentation change that requires a persisted schema migration or changes
Goal/Release gate semantics is REPLAN REQUIRED and returns to design.
