---
id: C-969
scope: {kind: objective, id: O-044}
result: NEEDS WORK
checked_by: {role: checker, session: check-o044-recheck-20261010}
executed_session: executor-o044-unrecorded
checked_at: '2026-10-10T06:34:45Z'
health_snapshot: sha256:c8f2f5d0cb734cf243f5ba48af978f2067b9b2dc55d1de7822509713d497d2b7
unmet: [O-044-SC2, T-122-DW1]
reviewed:
  base_commit: 3ddef5d
  head_commit: 3ddef5d
  files:
    - internal/data/next.go
    - internal/resume/resume.go
    - internal/data/decision_gate_v2_test.go
    - .savepoint/Design.md
  dependencies:
    - .savepoint/checks/C-968-o044-objective-check.md
    - .savepoint/issues/I-142-uncovered-unmet-requirement-routes-to-check-not-owner.md
    - .savepoint/issues/I-143-design-next-verbs-omit-assess-and-carried-decisions.md
issues: [I-142, I-143]
supersedes: C-968
---

# C-969: O-044 Full Objective Re-check

NEEDS WORK, narrowly. The repair routes the uncovered-requirement case to the
owner (`Accept`) for both Objectives and Tasks, and its evidence line names
the uncovered IDs. Design is reconciled. One defect remains inside the frozen
cell M17: the action line for that route still tells the owner to "accept the
current Check", which is a NEEDS WORK Check the owner cannot accept and which
the board does not offer to accept.

Reviewed state: commit `3ddef5d` plus the uncommitted working tree (repair
diff to `internal/data/next.go`, `internal/resume/resume.go`,
`internal/data/decision_gate_v2_test.go`, `.savepoint/Design.md`). This
re-check used C-968's frozen scope lock and matrix; no new axes. The checker
session that wrote C-968 ran this re-check; it did not build or repair the
work.

## Closure map of prior Issues

| Issue | State | Evidence |
|---|---|---|
| I-142 | **Still open** (partly repaired) | Route and evidence line fixed; action line wrong (below) |
| I-143 | Repair proven; stays open until a CLEAR Check can record `verified` | `.savepoint/Design.md` now lists `Assess`, the widened `Accept` (including uncovered unmet requirements), and a "Carried owner decisions (O-044)" entry with `carried_forward`, `scope`, `unmet` and routing |

## Admission ledger

| Re-check item | Prior claim | Frozen cell | Allowed result | Result |
|---|---|---|---|---|
| P4 Objective: carried exception, `unmet [TEST-08, DESIGN-02]` | I-142 | M17 | blocking | Verb `Accept O-002` ✓; evidence "latest check C-006 lists unmet requirements the exception does not cover: DESIGN-02" ✓; **Next action "Ask the owner to accept the current Check." ✗** |
| P4 Task variant (I-142 Proof Needed) | I-142 | M17 | blocking | `Accept T-001`; same evidence ✓; **same wrong action line ✗** |
| P1 carried, covered | regression | M14 | blocking | `Close`, ready to close by exception ✓ |
| P2 unassessed | regression | M15 | blocking | `Assess` ✓ |
| P3 changed exception | regression | M16 | blocking | `Accept`, renew-the-decision action ✓ |
| P5 no `unmet` list | regression | M18 | blocking | `Close` ✓ |
| P6 control sequence | regression | M19 | blocking | stripped ✓ |
| Task variant, `unmet [TEST-08]` covered | regression | M14 | blocking | `Close T-001`, ready to close by exception ✓ |
| Design reconciliation | I-143 | M26 | blocking | proven ✓ |
| Board refusal for `exception_scope` | I-142 Proof Needed | M20 | blocking | text unchanged and still accurate; `a` not offered (acceptance only) ✓ |
| M1–M13, M20–M25, M27 | C-968 | same | blocking | unchanged code; full gate green ✓ |

## Remaining defect (I-142)

- Violated: O-044 SC2 ("asked again only for a decision a material change
  actually affected, and is told what changed") and T-122 DW1/DW3 (the action
  phrase must match the route).
- Reproduction: either P4 shape. `internal/resume/resume.go` `ActionPhrase`,
  case `data.NextOwnerValidationRequired`, only special-cases
  `GateBlockDecisionChanged`; `GateBlockExceptionScope` falls to "Ask the
  owner to accept the current Check."
- Expected: an owner instruction to widen or renew the exception to cover the
  named requirements, or to send the work back for repair.
- Missing tests: the repair adds only
  `TestResolveObjectiveIntegrationRung_uncoveredRequirementAsksTheOwner`
  (Next kind). No resume/render test covers the route's evidence or action
  line, and no Task-variant test exists.

## Test and command results

- `make test-full` on the repaired tree, go1.26.2 linux/amd64,
  2026-10-10 ~06:33Z: exit 0 (all packages; linux/darwin/windows builds).
- `TestResolveObjectiveIntegrationRung_uncoveredRequirementAsksTheOwner`: PASS.
- Independent probes P1–P6 and the Task variants (scratch copy of the
  repaired tree, temporary projects): results in the ledger.
- `git diff --check` clean; canonical vs scaffold skill and check-method
  `cmp` identical.
- Code Health: `savepoint health check O-044 .` → snapshot
  `sha256:c8f2f5d0…d2b7` (created); "Code Health does not block clearance."

## Materiality

| Issue | Likelihood | Impact | Materiality | Recommendation |
|---|---|---|---|---|
| I-142 (remaining): wrong action line on the uncovered-requirement route | Medium — every such case shows it | Low — verb and evidence are right; the instruction points at an action that does not exist | Low | Fix now: one action-phrase branch plus a resume test for Objective and Task variants |

## Owner decisions on the scope

None recorded on O-044 or its Tasks; nothing to carry forward.

## Owner validation still needed

- T-122 DW7: owner review of the Close, Assess and Accept routes, and
  recorded acceptance. Still not recorded.

## Observations (non-blocking)

- The repair did not append `repair_attempted` history to I-142 or I-143, as
  `issue-capture.md` Out-Of-Scope Repair asks.
- The C-968 observations still stand (no board exception renewal; T-122
  closed by waiver without acceptance; carried Check need not be later than
  the originating one).

## Convergence

This was the one full re-check. If the remaining action-line fix is made, a
single targeted re-check of cell M17 (Objective and Task variants) and the
full gate follows. After that, stop and ask the owner rather than start
another cycle.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [ ] STYLE-03 **Test branches** — the new `GateBlockExceptionScope` branches in `resume.go` `EvidenceLines` have no test.
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [x] STYLE-07 **One source of truth**
- [x] STYLE-08 **Comments explain why** — `decisionRung` comment states why the owner, not a Check.
- [x] STYLE-09 **Content lives in data**
- [x] STYLE-10 **Small diffs** — repair is a few lines.
