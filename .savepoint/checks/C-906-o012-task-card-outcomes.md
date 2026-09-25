---
id: C-906
scope: {kind: objective, id: O-012}
result: NEEDS WORK
checked_by: {role: checker, session: o012-full-check-20260922}
executed_session: task-card-simplification-20260921
checked_at: '2026-09-22T09:28:27Z'
reviewed:
  base_commit: dcb4c7d736dd6d7c7581c7db5ac58d30d6356cc5
  head_commit: 84e9e502985aeb9c9602f3e08c8dbddd21233e4c
  files:
    - internal/board/v2/card.go
    - internal/board/v2/badges.go
    - internal/board/v2/plain.go
    - internal/board/v2/card_test.go
    - internal/board/v2/badges_test.go
    - internal/board/v2/columns_view_test.go
    - internal/board/v2/run_test.go
    - .savepoint/objectives/O-900-test-objective-for-ui-checks/Objective.md
    - .savepoint/objectives/O-900-test-objective-for-ui-checks/tasks/T-900.md..T-911.md
    - .savepoint/Design.md
    - .savepoint/Guardrails.md
    - .savepoint/visual-identity.md
    - AGENTS.md
    - README.md
    - templates/project-v2/AGENTS.md
  dependencies: []
issues: [I-018]
supersedes: null
---

# C-906: O-012 Full Objective Check

## Verdict

`NEEDS WORK`. The Task-card renderer, policy boundary, documentation, and
quality gates are sound, but the disposable O-900 visual spread does not render
the pending `[ ] CHECK` outcome that O-012 and T-005 require it to exercise.
I-018 records the failed integration cell. O-012 is not ready for owner closure.

## Evidence Mode

Full Objective evidence. This run reviewed all three owned Tasks (T-003, T-004,
and T-005), including their owner-waived optional Task Checks, plus cross-Task
integration and reconciliation against Design.

## Frozen Scope Lock

1. Acceptance and policy: O-012's nine Success Conditions; every Done-When
   criterion in T-003-T-005; DATA-02, DATA-03, TPL-02, ARCH-02, ARCH-04,
   TEST-01, TEST-02, TEST-04, TEST-08, and TEST-09; STYLE-03 and STYLE-07..10
   as advisory rules.
2. Changed behavior and public surfaces: `newTaskCard`, `TaskCard.badges`,
   `TaskCard.showsReviewOutcome`, `taskReviewOutcomeBadge`, `blockerBadge`,
   `renderCard`, `renderPlain`, the live O-900 records, and the active guidance
   named by T-004. Detail, Next, resume, doctor, Objective sidebar semantics,
   gate resolvers, and Release completion remain outside the presentation
   change except for confirming that they were not changed.
3. Relied-on orchestration: typed values from `internal/data` flow into one
   pure badge mapping and then the interactive and non-TTY renderers. There is
   no runtime network, subprocess, persistence, or external-service boundary;
   the finite external-boundary and side-effect matrices are not applicable.
4. Matrix axes: planned/in-progress/done; build/test/audit; missing/current/
   needs-work/stale/unknown clearance; waiver and exception precedence; task,
   Objective, replan, and owner blockers; interactive and plain output; colour
   and no-colour output; widths 20/28/40/72 and board widths 80/100/120/160;
   retained versus retired labels; and the twelve live O-900 fixture roles.
5. Materiality boundary: a finding must violate an O-012/T-003-T-005 criterion or
   named Guardrail through the supported renderer or live fixture. Historical
   archives, unrelated Issue edits, and alternative compact-copy preferences
   are outside this Check.

## Coverage Matrix

| Area | Expected invariant | Evidence | Classification |
| --- | --- | --- | --- |
| Columns and stages | Exactly Planned/In Progress/Done; Build/Test/Check only while in progress | `TestStageBadgeIsPresentOnlyWhileInProgress`, `TestBoardDrawsExactlyThreeColumns` | Proven |
| Outcome mapping | At most one outcome with exception → waiver → clearance precedence and six exact labels | `TestTaskReviewOutcomeBadgeRendersEveryPathAtExactText`, precedence and distinctness tests, independent live-card count | Proven |
| Open-card suppression | Missing/current stay off open cards; actionable outcomes and owner outcomes remain | `TestRenderCardInProgressOmitsCheckBadgeWhenNotActionable`, actionable and owner-blocker cases | Proven |
| Retired vocabulary | No completion badge or stale/unverified wording on Task cards | card and board negative assertions plus targeted source/document search | Proven |
| Blockers | Wait, Objective wait, replan, and unresolved owner action remain; accepted exception suppresses the resolved owner blocker | blocker exhaustiveness tests and card regression cases | Proven |
| Policy boundary | Waiver/exception stay distinct from technical clearance and do not change data, dependency, Objective, or Release gates | no `internal/data` change in the reviewed range; Design/Guardrails/AGENTS/template reconciliation | Proven |
| Output parity and layout | Interactive/plain use the same card badges; no-colour remains distinct; narrow output does not overflow | board, plain, colour-disabled, card-width, and board-width tests | Proven |
| O-900 outcome spread | Live fixture visibly renders all six retained outcomes and wait/replan/owner blockers | Independent live O-900 harness; blockers and five outcomes passed, `[ ] CHECK` absent on card and plain paths | Issue — I-018 |
| O-900 integrity | Twelve IDs and ownership remain valid with 4/4/4 column spread; strict load succeeds | live V2 index load and card projection | Proven |
| Active guidance | Design and active guidance match the renderer and preserve verification policy | targeted search and direct reconciliation | Proven |
| Task evidence | T-003-T-005 are owner-done with explicit Task-Check waivers; waivers do not substitute for this Check | frontmatter and TEST-09 review | Proven |

Every mandatory cell was classified. No external-boundary or side-effect cell
applies because the scoped behavior is a pure projection and render over
already-loaded records.

## Independent Adversarial Probe

A temporary, uncommitted package-level Check harness loaded the live
`.savepoint` index, projected O-900 through `groupTaskCardsFor`, counted review
outcomes per card, and rendered the same selection through `renderPlain`.
Every card had at most one outcome. It observed:

```text
[✓] CHECK: T-907
[!] NEEDS WORK: T-905
[!] REVIEW: T-909
[✓] WAIVED: T-910
[✓] OWNER ACCEPTED: T-908
→ WAITS T-900: T-901
⚠ REPLAN: T-904
! AWAITS OWNER: T-906
```

`[ ] CHECK` had no card, and the plain output also lacked it. The harness was
removed after the run; `git status --short` confirmed no temporary file
remained.

## Issue

- I-018 — O-900 no longer renders the pending Check outcome. T-005 says T-903 is
  the pending example, but T-003's final open-card rule intentionally suppresses
  missing clearance on every in-progress card. The fixture and its recorded
  evidence were not reconciled after that refinement.

## Commands

- `go test ./internal/board/v2 -run TestO012CheckLiveO900OutcomeMatrix -count=1 -v`
  — intentional independent harness failed only for absent `[ ] CHECK` on the
  live card matrix and plain output.
- `go test ./internal/board/v2/... ./internal/data/...` — pass.
- `git diff --check` — pass.
- `make build` — pass.
- `make test` — pass; all packages passed, including `internal/migrate`
  (121.009s).

The pre-existing modified Issue records I-002, I-004, I-008-I-011, and I-013-I-015
were outside this Check and were not changed by it.

## Materiality

| Issue | Likelihood | Impact | Materiality | Recommendation |
| --- | --- | --- | --- | --- |
| I-018 | High — every O-900 render omits the promised cell | Medium — owner visual verification cannot inspect one retained outcome and recorded evidence is false | Medium | Fix now in T-005, add durable fixture-level regression evidence, then run a fresh Full Objective re-check |

## Non-blocking Observation

`internal/board/v2/columns_view_test.go:48-50` still says pending `[ ] CHECK`
is covered by `TestRenderCardAuditAlwaysShowsCheckBadgeEvenWhenMissing`, but
that test was replaced when the final open-card rule began suppressing missing
clearance at audit too. This is advisory STYLE-08 drift; it does not create a
second Issue or change the verdict.

## Owner Action

The owner must retreat T-005 from `done` to `in_progress` at `stage: build`
before an executor may repair I-018. The checker did not alter Task lifecycle,
repair the fixture, or mark O-012 complete.
