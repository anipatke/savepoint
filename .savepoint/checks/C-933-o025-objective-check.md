---
id: C-933
scope: {kind: objective, id: O-025}
result: CLEAR
checked_by: {role: checker, session: o025-objective-check-20260925}
executed_session: o025-t047-t048-execution-20260925
checked_at: '2026-09-25T20:59:13Z'
reviewed:
  base_commit: 5ce6b014bc2d7a887c809b795b684682901ca9a8
  files:
    - internal/data/release_gate_v2.go
    - internal/data/next.go
    - internal/resume/resume.go
    - internal/resume/evidence.go
    - internal/doctor/checks.go
    - internal/board/v2/detail.go
    - internal/board/v2/detail_view.go
    - internal/data/release_gate_v2_test.go
    - internal/data/next_test.go
    - internal/resume/resume_test.go
    - internal/doctor/checks_test.go
    - internal/board/v2/detail_test.go
    - internal/board/v2/releases_test.go
    - main_resume_matrix_test.go
    - internal/init/agent_skills_test.go
    - .savepoint/Design.md
    - .savepoint/Guardrails.md
    - .savepoint/router.md
    - AGENTS.md
    - README.md
    - agent-skills/savepoint-check/SKILL.md
    - agent-skills/savepoint-design/SKILL.md
    - agent-skills/savepoint-task/SKILL.md
    - agent-skills/references/check-method.md
    - agent-skills/references/commands-and-procedures.md
    - agent-skills/references/issue-capture.md
    - templates/project-v2/.savepoint/Design.md
    - templates/project-v2/.savepoint/Guardrails.md
    - templates/project-v2/.savepoint/router.md
    - templates/project-v2/AGENTS.md
    - templates/project-v2/agent-skills/savepoint-check/SKILL.md
    - templates/project-v2/agent-skills/savepoint-design/SKILL.md
    - templates/project-v2/agent-skills/savepoint-task/SKILL.md
    - templates/project-v2/agent-skills/references/check-method.md
    - templates/project-v2/agent-skills/references/commands-and-procedures.md
    - templates/project-v2/agent-skills/references/issue-capture.md
  dependencies: []
issues: []
supersedes: null
---

# C-933: O-025 Objective Check

## Independence and Scope

This fresh checker session did not plan or build O-025. Full evidence covers
both owned Tasks, their cross-Task integration, and reconciliation with Design.

## Scope Lock

1. Goal completion requires at least one member Objective and every member to
   pass `ResolveObjectiveCompletion`; Goal-scoped Check state, linked Issues,
   exceptions, and owner acceptance do not affect that decision. Migrated
   legacy completion remains supported.
2. Data, Next/resume, doctor, and board are the supported runtime surfaces.
   Check decoding remains the compatibility surface for `scope.kind: release`.
3. Active repository guidance and its new-project scaffold describe only the
   optional Task Check and mandatory Full Objective Check.
4. Matrix axes are Goal membership (empty, incomplete, complete), completion
   representation (ordinary and legacy), Goal Check evidence (missing, current,
   stale, unknown, NEEDS WORK, and linked unresolved Issue), and surface (gate,
   Next/resume, doctor, board, decoder, live guidance, scaffold). Network,
   timeout, cancellation, numeric, Unicode, and mutable-input cells are not
   applicable because the scoped behavior has no such input or boundary.
5. A material Issue must violate an O-025 success condition through one of the
   supported surfaces above.

## Coverage Matrix

| Cell | Independent evidence | Result |
| --- | --- | --- |
| Complete members, no Goal Check | Gate and on-disk resume matrix tests; code trace returns owner-ready | Pass |
| Complete members with stale, unknown, or NEEDS WORK Goal Check | Parameterized gate tests and NEEDS WORK resume fixture | Pass |
| Goal Check with unresolved linked Issue or exception | Gate and doctor fixtures ignore both for Goal completion | Pass |
| Empty Goal | Gate retains `release-no-objectives`; doctor retains membership diagnostic | Pass |
| Incomplete member Objective | Gate composes `ResolveObjectiveCompletion`; Next and doctor retain the blocker | Pass |
| Migrated legacy completion, with and without live members | Legacy gate tests preserve historical behavior | Pass |
| Existing `scope.kind: release` record | Decoder and strict-loading fixtures accept it | Pass |
| Resume output | Ready action uses member completion and contains no clearance or owner-wait copy | Pass |
| Board Goal detail | Membership readiness remains; clearance and owner-validation sections are absent | Pass |
| Active guidance and scaffold | Targeted search leaves only Design's compatibility note; parity tests pass | Pass |
| Full repository gate | Fresh `make test-full`, including Linux, macOS, and Windows builds | Pass |

## Acceptance Classification

- Goal completion from member Objective completion only: **Proven**.
- Goal-ready Next action with obsolete Goal Check and acceptance kinds removed:
  **Proven**.
- Doctor reports membership problems only: **Proven**.
- Board Goal detail omits Check clearance and owner acceptance: **Proven**.
- Release-scoped Check records still load and do not affect decisions:
  **Proven**.
- Live and scaffold guidance consistently describe two Check kinds: **Proven**.
- Repository resume and doctor behavior are covered by on-disk and package
  integration tests: **Proven**.

## Commands

- `go test ./internal/data ./internal/resume ./internal/doctor ./internal/board/v2 ./internal/init . -count=1` — passed.
- `git diff --check` — passed.
- Targeted active-guidance and runtime reference searches — passed; the only
  active `scope.kind: release` wording is the required compatibility note.
- `make test-full` — passed at 2026-09-25T20:59Z with all Go packages and
  Linux, macOS, and Windows builds successful.

## Issues and Materiality

No materiality actions are required. No in-scope Issue was found.

## Owner Validation

Neither owned Task requires owner validation. This CLEAR Check makes O-025
technically ready for the owner to accept and close; the Check does not change
the Objective's status.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [x] STYLE-03 **Test branches**
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [x] STYLE-07 **One source of truth**
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data**
- [x] STYLE-10 **Small diffs**
