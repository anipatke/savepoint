---
id: C-968
scope: {kind: objective, id: O-044}
result: NEEDS WORK
checked_by: {role: checker, session: check-o044-20261010}
executed_session: executor-o044-unrecorded
checked_at: '2026-10-10T06:29:46Z'
health_snapshot: sha256:ff249badef70c99f3f1f036e45401bab2758828c5c82a522fbd4ae4e24dc6c74
unmet: [O-044-SC2, O-044-SC6, T-122-DW1, T-122-DW2]
reviewed:
  base_commit: 75b329f
  head_commit: 3ddef5d
  files:
    - internal/data/carry_v2.go
    - internal/data/decision_gate_v2.go
    - internal/data/evidence_v2.go
    - internal/data/check_v2.go
    - internal/data/project.go
    - internal/data/write.go
    - internal/data/gate_v2.go
    - internal/data/objective_gate_v2.go
    - internal/data/dependency.go
    - internal/data/next.go
    - internal/resume/evidence.go
    - internal/resume/resume.go
    - internal/board/v2/actions.go
    - internal/board/v2/badges.go
    - internal/board/v2/detail_view.go
    - internal/board/v2/io.go
    - internal/board/v2/next_panel.go
    - agent-skills/savepoint-check/SKILL.md
    - agent-skills/references/check-method.md
    - templates/project-v2/agent-skills/savepoint-check/SKILL.md
    - templates/project-v2/agent-skills/references/check-method.md
    - AGENTS.md
    - templates/project-v2/AGENTS.md
    - README.md
    - internal/init/agent_skills_test.go
  dependencies:
    - .savepoint/Design.md
    - .savepoint/Guardrails.md
    - .savepoint/issues/I-141-recheck-invalidates-owner-acceptance-and-exceptions.md
issues: [I-142, I-143]
supersedes: null
---

# C-968: O-044 Full Objective Check

NEEDS WORK. Fresh checker conversation; it built none of T-120 to T-123.
The runtime core is sound: decisions carry forward, the TheShed O-002 shape
reads "ready to close by exception" with clearance still NEEDS WORK, and
unassessed and changed decisions route to `Assess` and `Accept`. Two gaps
remain: an exception that misses an unmet requirement still routes to
another Check without naming the requirement (I-142), and Design was not
reconciled with the new verb and fields (I-143).

Reviewed state: commit `3ddef5d` (T-120 to T-122 code) plus the uncommitted
working tree carrying T-123's guidance edits and the T-122/T-123 board
closures. No implementation, Design, router, Task or Objective status, or
owner decision was changed by this Check.

## Frozen scope lock

1. O-044 SC1–7; T-120 DW1–6; T-121 DW1–7; T-122 DW1–7; T-123 DW1–6;
   guardrails DATA-01/02/03, FS-01/02, TPL-01/02/04, ARCH-02/05, TEST-01/02/03/04/05/08/09,
   STYLE (advisory). Gate: fresh `make test-full`; official health snapshot.
2. Public surfaces: evidence decoder (`owner_validation.scope`,
   `carried_forward` on acceptance and exception), Check decoder (`unmet`),
   index reference checks, evidence writer, `assessDecision` and its callers
   (Task/Objective completion, Task `requires: accepted` and Objective
   dependency, Task/Objective consistency), `ResolveNext` rungs, `resume`
   verb/evidence/action lines, board `a`/`x` actions, badges, detail view,
   Next panel; savepoint-check Closure Rules, check-method re-check, AGENTS.md
   (both copies), README, scaffold copies.
3. Relied-on orchestration: `LoadV2Index` strict load → gates → `ResolveNext`
   → `resume.Render` / board projection; board accept → `writeOwnerAcceptanceCmd`
   → evidence writer → reload.
4. Matrix below; text-width/Unicode classes N/A (no truncation in scope);
   external-boundary matrix N/A (no server, subprocess or provider added).
5. Admit only behaviour O-044 introduced or promises. Pre-existing gate
   semantics outside the decision rule are observations.

## Coverage matrix

| # | Surface | Cell | Result | Evidence |
|---|---|---|---|---|
| M1 | Carry decoder | valid acceptance + exception carries; scope | Passed | `TestDecodeEvidenceV2_ownerValidationScopeAndCarriedForward`, `_exceptionCarriedForward` |
| M2 | Carry decoder | missing check/applies/reason/assessed_at, bad role, owner `applies:false`, missing `material_change`, origin Check named | Passed | `TestDecodeEvidenceV2_carriedForwardRejections` (12 cases); source `carry_v2.go:49-85` |
| M3 | Carry decoder | out of order, duplicate per Check per role; checker+owner at same Check allowed | Passed | `_carriedForwardOrderAndDuplicates`; `compareV2CheckIDs` numeric (`project.go:537`) |
| M4 | Carry decoder | scope/carry without `accepted_check`; blank scope entry | Passed | `_scopeAndCarriedForwardRequireAcceptedCheck`, `_ownerValidationScopeBlankEntryRejected` |
| M5 | Check decoder | `unmet` valid; blank entry; on CLEAR | Passed | `TestDecodeCheckV2_unmet`, `_unmetRejections` |
| M6 | Index refs | dangling carried Check; cross-scope carried Check | Passed | `TestLoadV2Index_carriedForwardReferences` (6) |
| M7 | Writer | round trip Task + Objective, unknown fields/body kept; absent fields not written | Passed | `TestWriteEvidenceV2_roundTripsScopeAndCarriedForwardForTaskAndObjective`, `TestWriteTaskEvidenceV2_omitsScopeAndCarriedForwardWhenEmpty`; existing no-op byte tests in full gate |
| M8 | `assessDecision` | origin=latest; none at latest; only earlier; checker applies; checker changed; owner renewal; owner outranks checker; no latest | Passed | `TestAssessDecision` (8) |
| M9 | Objective completion | carried exception grants (owner authority, carry exposed); unassessed; changed exception only; owner renewal; uncovered unmet; no `unmet` list; unfinished Task still blocks | Passed | `decision_gate_v2_test.go` objective cases |
| M10 | Task completion | carried acceptance → checker authority; changed acceptance names change, owner renewal restores | Passed | `TestResolveTaskCompletion_acceptanceCarriedForward…`, `_changedAcceptanceNamesTheChange` |
| M11 | Dependencies | Task `requires: accepted` follows applicability; Objective dependency cleared-by-exception via carry | Passed | `TestResolveTaskDependencyV2_acceptedFollowsApplicability`, `TestObjectiveGates_carriedExceptionIsHonoured…` |
| M12 | Consistency | Task acceptance carried / unassessed / changed; Objective done by carried exception | Passed | `TestInspectTaskConsistency_carriedAcceptanceIsNotSuperseded`, `TestObjectiveGates_…` |
| M13 | Regression SC1 | TheShed O-002 shape from disk, with and without carries | Passed | `TestTheShedO002_*`; independent probe P1/P2 below |
| M14 | Next/resume | carried → `Close`, clearance and completion on separate lines, origin + carried Check + assessor + reason | Passed | probe P1; `TestRender_closeByExceptionNamesCarry` |
| M15 | Next/resume | unassessed → `Assess` naming decision and C-006; action says not a new Check | Passed | probe P2; `TestRender_assessAndAcceptRoutes` |
| M16 | Next/resume | changed exception → `Accept` naming decision and change | Passed | probe P3 |
| M17 | Next/resume | applicable exception, latest `unmet` has an uncovered ID | **Issue I-142** | probe P4: `Check O-002`, "Record the Objective O-002 integration Check.", no IDs |
| M18 | Next/resume | `unmet` absent → checker assessment relied on | Passed | probe P5 |
| M19 | ARCH-05 | control sequence in carry reason | Passed | probe P6: `"\e[31m"` stripped to `evil RED` |
| M20 | Board accept | renewal appends owner carry at latest, keeps `accepted_check`; refused when no longer offered, file unchanged | Passed | `TestOwnerAcceptanceRenewalAppendsCarryAndKeepsOrigin` |
| M21 | Board `x` | close by carried exception | Passed | same `ResolveObjectiveCompletion` decision as M9; existing exception-close tests in full gate |
| M22 | Board detail | scope and every carry entry for both decisions | Passed | `detail_view.go:313-345` source trace |
| M23 | Guidance | Closure Rules carry rule, shapes, re-check assessment, `Assess`, structured exception; old exact-Check wording gone | Passed | `TestSavepointCheckSkillClosureRules`; `grep` finds no old wording in skills/templates/README/AGENTS |
| M24 | Guidance | AGENTS.md both copies route `Assess`; one Verification Policy bullet | Passed | `git diff AGENTS.md templates/project-v2/AGENTS.md`; `TestAgentsGuideRoutesAssessToCheckSkill` |
| M25 | TPL-01 | canonical vs scaffold byte parity | Passed | `cmp` clean ×2; `TestV2SkillSetIsCompleteWithByteParity` |
| M26 | Design reconciliation | Next verbs and decision fields | **Issue I-143** | `.savepoint/Design.md:198` lacks `Assess`; no carry/unmet record |
| M27 | Prior Checks immutable; no decision inferred; no auto-close | Passed | source trace: gates read only; only writer is board `a` on owner action |

Workflow lock (board `a`): load index → gate re-evaluated → refusal returns
error with no write (M20) → evidence patch → atomic write → reload. Failure
before write leaves the file byte-identical (asserted). Other surfaces are
read-only projections.

## Independent probes

Scratch copy of the tree, external test in `internal/resume`, temporary
project (not TheShed's files), router `G-001`/`O-002`, two done Tasks,
acceptance and exception (`requirements: [TEST-08, CFG-03]`) at C-005,
C-005 and C-006 NEEDS WORK, rendered through `LoadV2Index → ResolveNext →
resume.Render`.

| Probe | Setup | Result |
|---|---|---|
| P1 | checker carries at C-006, `unmet [TEST-08]` | `Close O-002`; "Technical clearance: Check C-006 recorded NEEDS WORK."; "Completion: Ready to close by exception — … at Check C-005 …; carried to Check C-006: checker session chk confirmed it still applies — only Design text changed" |
| P2 | no carries | `Assess O-002`; "owner exception recorded at C-005 has not been assessed at latest check C-006"; action "… this is not a new Check." |
| P3 | exception `applies: false` | `Accept O-002`; names exception and change; acceptance unaffected |
| P4 | carries, `unmet [TEST-08, DESIGN-02]` | `Check O-002`; "Record the Objective O-002 integration Check." — Issue I-142 |
| P5 | carries, `unmet []` | `Close O-002`, by exception |
| P6 | `\e[31m` in reason | escape removed |

## Acceptance classification

- **O-044:** SC1 Proven; SC2 **Issue** (I-142: uncovered-requirement change
  is not routed to the owner or named); SC3 Proven; SC4 Proven; SC5 Proven
  (legacy decisions at latest Check unchanged; older ones read unassessed);
  SC6 **Issue** (I-142 surface mismatch with README/Design; I-143 Design);
  SC7 Proven.
- **T-120:** DW1–6 Proven.
- **T-121:** DW1–3, 5–7 Proven; DW4 Proven (failing-before shown by the five
  updated legacy tests and `TestTheShedO002_withoutCarryEntriesReportsUnassessed`).
- **T-122:** DW1 **Issue** and DW2 **Issue** (I-142: `exception_scope` has no
  rung and no evidence line); DW3–6 Proven; DW7 **Unverified — owner**: no
  owner acceptance of the three routes is recorded (`owner_validation.required: true`,
  `accepted_check: ""`).
- **T-123:** DW1–6 Proven.

## Test and command results

- `make test-full` at `3ddef5d` + working tree, go1.26.2 linux/amd64,
  2026-10-10 ~06:25Z: exit 0 (all packages, plus linux/darwin/windows builds).
- Focused re-runs of every test named above: PASS.
- `git diff --check`: clean. `cmp` canonical vs scaffold skill and
  check-method: identical.
- `quality_gates` in `.savepoint/config.yml` are all null; the repository's
  full gate is `make test-full` per AGENTS.md Build.
- Code Health: `savepoint health check O-044 .` → snapshot
  `sha256:ff249bad…c74` (created). "Code Health does not block clearance."
  All five instances optional, no blocking finding.

## Applicable Guardrails

DATA-01 (M7), DATA-02 (no new lifecycle vocabulary; new Next kind lives in
`internal/data`), DATA-03 (M2–M6), FS-01 (M20), FS-02/TPL-04 (existing
upgrade-assets freshness tests in full gate; guidance files are managed
assets), TPL-01 (M25), TPL-02 (README claims "you are asked again only for
the decision that changed"; holds except I-142), ARCH-02 (renderers do no IO;
write happens in a `tea.Cmd`), ARCH-05 (M19), TEST-01/02/03/04/05/08 met;
TEST-09: T-120–T-123 each carry an owner Task-check waiver.

## Materiality

| Issue | Likelihood | Impact | Materiality | Recommendation |
|---|---|---|---|---|
| I-142 exception missing an unmet ID routes to `Check` | Medium — any re-check that finds a new unmet requirement while an exception exists | Medium — owner is sent round the very re-check loop O-044 removes, and is not told which requirement is uncovered | Medium | Fix now: add the rung and evidence line, plus Next/resume/board tests |
| I-143 Design omits `Assess` and carried decisions | High (always) | Low — documentation only; runtime and skills are right | Low | Fix now alongside I-142 via `savepoint-design` |

## Owner decisions on the scope

No owner acceptance or exception is recorded on O-044 or its Tasks, so
there is nothing to carry forward. T-120–T-123 carry owner Task-check
waivers, which are not technical CLEAR.

## Owner validation still needed

- T-122 DW7: the owner reviews the resume output and board screens for the
  three routes and records acceptance. The probe outputs above are a ready
  reference. Best done after I-142 is repaired, since that changes one route.

## Observations (non-blocking)

- The board offers no way to renew an exception (`a` covers acceptance only),
  while resume's `Accept` action asks the owner to renew a changed
  exception. T-122 DW4 scoped `a` to acceptance, so this is follow-up, not a
  finding; today an agent records the renewal on the owner's instruction.
- The board let the owner close T-122 (which declares
  `owner_validation.required`) with a Task-check waiver and no acceptance.
  That is the pre-existing waiver path in `ResolveTaskCompletion`, not
  introduced by O-044.
- T-120 limitation stands: a carried Check is not required to be later than
  the originating Check (only different and same-scope).
- The two `.coverage.out.*` and `.go-test.json.*` untracked files predate
  this Check.

## Code Style Review

- [x] STYLE-01 **One job per file** — carry decoding (`carry_v2.go`) and applicability (`decision_gate_v2.go`) split out.
- [x] STYLE-02 **One job per function**
- [x] STYLE-03 **Test branches**
- [x] STYLE-04 **Types document intent** — `DecisionApplicability`, `DecisionKind`, typed blockers.
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [x] STYLE-07 **One source of truth** — `assessDecision` is the single rule; `acceptanceOffered` shared by board offer and write.
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data** — refusal templates in a map; wording in `internal/resume`.
- [ ] STYLE-10 **Small diffs** — commit `3ddef5d` adds ~3,200 lines across 47 files, mixing O-044 code with planning records for O-039–O-043 and a product spike.
