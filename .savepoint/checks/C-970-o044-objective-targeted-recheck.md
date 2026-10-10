---
id: C-970
scope: {kind: objective, id: O-044}
result: CLEAR
checked_by: {role: checker, session: check-o044-final-20261010}
executed_session: executor-o044-unrecorded
checked_at: '2026-10-10T06:38:20Z'
health_snapshot: sha256:d79917ae7be4fe01b1d083446401e3d1de5911cb6998a3910bb1ffcf557bda64
reviewed:
  base_commit: 3ddef5d
  head_commit: 3ddef5d
  files:
    - internal/resume/resume.go
    - internal/resume/resume_test.go
    - .savepoint/objectives/O-044-keep-owner-decisions-across-rechecks/tasks/T-122-show-close-assess-and-accept-steps-on-the-status-line-and-board.md
  dependencies:
    - .savepoint/checks/C-968-o044-objective-check.md
    - .savepoint/checks/C-969-o044-objective-recheck.md
issues: []
supersedes: C-969
---

# C-970: O-044 Targeted Re-check

CLEAR. This is the targeted re-check that C-969's convergence note allowed: frozen cell M17
(Objective and Task variants) and the full gate, against C-968's scope lock. No new axes. The
checker session that wrote C-968 and C-969 ran it; it built and repaired none of the work.

Reviewed state: commit `3ddef5d` plus the uncommitted working tree, including the action-phrase
fix in `internal/resume/resume.go` and `TestRender_uncoveredRequirementAsksToWidenTheException`.

## Closure map of prior Issues

| Issue | State | Evidence |
|---|---|---|
| I-142 | Closed — verified by this Check | Route `Accept`, evidence names uncovered IDs, action line now "Ask the owner to widen or renew the exception to cover the named requirements, or send the work back for repair." for Objective and Task |
| I-143 | Closed — verified by this Check | Design reconciliation proven in C-969; Design unchanged since |

## Admission ledger

| Re-check item | Prior claim | Frozen cell | Result |
|---|---|---|---|
| P4 Objective, `unmet [TEST-08, DESIGN-02]` | I-142 | M17 | `Accept O-002`; evidence names DESIGN-02; correct action line ✓ |
| P4 Task variant | I-142 | M17 | `Accept T-001`; evidence names DESIGN-02; correct action line ✓ |
| P1, P2, P3, P5, P6; Task covered variant | regression | M14–M16, M18, M19 | `Close` / `Assess` / `Accept` / `Close` / sanitised / `Close` — unchanged ✓ |
| Full gate | TEST-08 | gate | `make test-full` exit 0 ✓ |

## Acceptance classification

All O-044 Success Conditions 1–7 and T-120–T-123 Done When criteria are Proven, building on
C-968 (unchanged cells) and C-969 (I-143, Design). T-122 DW7: the owner's review and approval of
the three routes after the final wording fix is recorded in T-122's Technical Evidence.

## Test and command results

- `make test-full` on the repaired tree, go1.26.2 linux/amd64, 2026-10-10 ~06:37Z: exit 0, no
  failing tests; linux/darwin/windows builds succeed.
- `TestRender_uncoveredRequirementAsksToWidenTheException`: PASS.
- Independent probes (scratch copy, temporary projects): ledger above.
- `git diff --check` clean; canonical vs scaffold skill and check-method `cmp` identical.
- Code Health: snapshot `sha256:d79917ae…da64` (created); "Code Health does not block clearance."

## Materiality

No Issues remain; no materiality actions required.

## Owner decisions on the scope

None recorded on O-044 or its Tasks; nothing to carry forward.

## Owner validation still needed

- T-122 declares `owner_validation.required` but its frontmatter still reads
  `accepted_check: ""`. The approval exists only as prose in its Technical Evidence. If the owner
  wants it in structured form, they record acceptance against this Check (C-970). This Check does
  not record it on their behalf.
- Closing O-044 is the owner's decision.

## Observations (non-blocking)

- The C-968 observations stand as follow-up: no board action renews an exception; T-122 was closed
  by a Task-check waiver despite `owner_validation.required` (pre-existing waiver path); a carried
  Check need not be later than the originating Check.
- The repairs did not add `repair_attempted` history to I-142 or I-143.
- The work, Checks, Issues and health snapshots are uncommitted.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [x] STYLE-03 **Test branches** — the exception-scope branch in resume now has a render test.
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [x] STYLE-07 **One source of truth**
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data**
- [x] STYLE-10 **Small diffs** — the repair adds three lines and one test.
