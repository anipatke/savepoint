---
id: C-961
scope: {kind: objective, id: O-034}
result: NEEDS WORK
checked_by: {role: checker, session: check-o034-20261003}
executed_session: build-t100-20261003
checked_at: '2026-10-02T20:24:00Z'
health_snapshot: sha256:992d4e3d6373b0d75e0b1713e52990d227c71d932c208a4252115eaa9ad5d1a3
reviewed:
  base_commit: f654ca1
  head_commit: b27c17a
  files:
    - AGENTS.md
    - README.md
    - agent-skills/savepoint-idea/SKILL.md
    - agent-skills/savepoint-design/SKILL.md
    - agent-skills/savepoint-task/SKILL.md
    - agent-skills/savepoint-check/SKILL.md
    - agent-skills/references/check-method.md
    - agent-skills/references/issue-capture.md
    - agent-skills/references/commands-and-procedures.md
    - templates/project-v2/AGENTS.md
    - templates/project-v2/.savepoint/Design.md
    - templates/project-v2/.savepoint/Guardrails.md
    - templates/project-v2/agent-skills/savepoint-design/SKILL.md
    - templates/project-v2/agent-skills/savepoint-task/SKILL.md
    - templates/project-v2/agent-skills/savepoint-check/SKILL.md
    - templates/project-v2/agent-skills/references/issue-capture.md
    - internal/init/agent_skills_test.go
    - internal/init/template_freshness_test.go
    - internal/init/upgrade_test.go
  dependencies: []
issues: [I-123]
supersedes: null
---

# C-961: O-034 Full Objective Check

## Result and independence

NEEDS WORK solely for I-123: two new README statements promise more edited-skill
protection than supported pre-manifest upgrades provide. A wording repair is
sufficient; no implementation repair is requested. This is a fresh conversation
that executed none of T-099, T-100 or T-101. All three are owner-completed with
explicit board-owner waivers. Their statuses, router, implementation, acceptance
criteria and Design remain unchanged. Scope includes b27c17a and uncommitted
T-101 guidance/tests present at review, not just HEAD.

## Frozen scope lock

1. Test O-034's seven Success Conditions and all Done When criteria of T-099,
   T-100 and T-101; FS-02, TPL-01/02/04, POL-01/02, TEST-01/02/03/08/09.
   Full gate is `make build && make test-full`; collect official health after it.
2. Public surfaces: four skills, three references, their scaffold copies, both
   AGENTS guides, changed Design/Guardrails templates, README and init content,
   parity and upgrade tests. No production Go behavior changed.
3. Direct reliance: actual upgrade and manifest behavior for delivery/preservation;
   existing runtime dependency gates for waivers; owner authority, project gates,
   package ownership and independent Full Check remain mandatory.
4. Matrix R1–R10 below covers normal, missing/failure, bypass and repeat where
   applicable. R9 includes tracked clean, tracked edited and pre-manifest inputs
   crossed with dry-run/apply/repeat. R6 includes aim/watch/over-watch and blocking
   overrides. Written walkthroughs are the owner-approved scenario method.
5. Admit only supported violations introduced, touched or promised by O-034.
   Unrelated runtime internals, existing Design drift and unrelated platform
   evidence are outside the perimeter. No new renderer, network, text-width,
   parser, numerical limit or transaction behavior was implemented here.

Scope was recorded before adversarial probes in `/tmp/o034-check-scope.md`.
No amendment or additional blocking axis was introduced.

## Coverage matrix and adversarial pass

| Cell | Normal | Failure / boundary / bypass | Repeat / handoff | Classification |
|---|---|---|---|---|
| R1 review and decisions | Seven files plus AGENTS, six axes, R-007 evidence and apply/follow-up/no-change present | No-activity idea/migration files explicitly say no change; product wording deferred to I-122 | Apply list matches diff; deviations recorded; baseline scenarios retained | Passed |
| R2 Next routing | Start/Build/Test → task; Plan/Replan → design; Check → check | Pasted Check with task router (this session) selects check; stale/unnamed selections remain AGENTS-owned | No substitute selection or router write | Passed |
| R3 replan | Frontmatter replan routes to planner | Body-only note insufficient; partial work and lifecycle preserved | Re-entry step 3, confirmation/readiness before detailing | Passed |
| R4 waiver integration | Three owner Task waivers present | clear dependency met, accepted dependency not met; waiver cannot clear an Objective or authorize executor self-review | All owned Tasks inspected at Full; owner closes work | Passed |
| R5 direct Issue repair | Proof Needed and neighbouring cases before repair_attempted | Unrun platform/case explicitly unverified; no executor verified closure | Single linked Objective advances; zero/multiple cannot be guessed; Goal preserved | Passed |
| R6 health advice | Needs Attention brought to watch line, prioritizing risky production code | Aim-to-watch band observation; scope narrowing owner-only; score-only branch moves discouraged; blocking verdict still prevents CLEAR | Residual watch result recorded; health command Full-only | Passed |
| R7 retrospective | Final Objective added after other Objectives planned; planner owns recorded outcome | No-change conclusion supported; downstream packaged skills become Issue suggestions rather than edits | Goal completion remains member-Objective completion; no new state/field/Goal Check | Passed |
| R8 parity | Seven pairs byte-identical | Existing parity/content tests detect absent or stale phrases | Both trees tested | Passed |
| R9 upgrade | Tracked clean skill receives real revised text; design receives retrospective | Tracked edit remains live with .new; pre-manifest authored skill replaced with .bak contradicts README | Dry-run preserves bytes; apply/repeat deterministically updated/unchanged or conflict | Issue I-123 |
| R10 policy/token weight | Duplicate AGENTS sections and Design policy removed | Policy remains project-owned; existing Design/Guardrails files not overwritten; fresh-init-only changes disclosed | 82,210 → 83,819 → 84,726 bytes for seven assets; named additions justified | Passed |

Environment/input axes not applicable: color, redirected TTY, Unicode width,
nonfinite numbers, serialization-model mutation and remote provider response
classes have no changed surface. External boundary is upgrade filesystem IO:
missing/edited/legacy inputs and dry-run/apply/repeat are covered; network,
redirect, connection refusal, cancellation and secrets are N/A to these edits.
No new partial-write contract is claimed.

## Workflow and side effects

| Order | Real operation | Effect / failure ownership | Independent oracle |
|---|---|---|---|
| 1 | Review Goal evidence | Findings, no implementation edits; missing evidence unverified | R-007 records |
| 2 | Apply agreed findings | Canonical assets changed; unsupported edits fail scope | T-099 findings versus diff |
| 3 | Mirror templates | Package bytes updated; mismatch blocks | Byte comparison |
| 4 | Upgrade dry-run | No asset/manifest write; errors belong to upgrade | Existing dry-run matrix and scratch original bytes |
| 5 | Upgrade apply | Clean replacement, edited .new, legacy .bak; primary IO failure returned | Original/incoming exact bytes and existing failure tests |
| 6 | Repeat upgrade | Clean unchanged, unresolved conflict same sidecar | Existing idempotence tests and independent matrix |
| 7 | Task handoff | Owner waivers; executor cannot grant clearance | Three Task records and full runtime test suite |
| 8 | Full Check/health | New immutable Check and official snapshot | Gate output and health verdict |

Existing upgrade failure-before-write, backup-before-replacement, sidecar failure
and partial-write/manifest reporting are retained and exercised by the full suite.
No new write operation or cleanup policy is introduced by these asset changes.

## Acceptance coverage

Criterion numbers below follow each record's displayed order.

- **T-099 1–7: Proven.** Per-file six-axis findings cover every asset; independent
  targeted evidence confirms T-050/T-054 replans, 15 NEEDS WORK Checks and 48 owner
  waivers, cited repair chains and convergence. Every finding has a disposition;
  Code Health advice and product follow-up are included. Review-only Task preceded
  implementation and recorded five baseline scenarios. No live agent run required.
- **T-100 1–8: Proven for its specified outcomes.** A1–A8/T1/T2/T4 match assets
  and documented deviations; I-122 exists; seven pairs match; named content tests
  and real-template clean/edited upgrade tests pass; five walkthroughs retain the
  intended behavior; byte weights independently reproduced; fresh full gate
  exceeds the ordinary build/fast requirement. Tracked edited-project test proves
  its own case, not the additional absolute README claim (I-123).
- **T-101 1–5: Proven.** Rule names timing/owner/outcome/downstream ownership,
  preserves Goal semantics, parity and real upgrade delivery, includes walkthrough
  and passes gate/content tests. Final retrospective order is a planning instruction,
  not an added runtime Goal completion gate.
- **O-034 1–6: Proven.** O-032 is done before this final Goal retrospective;
  findings cite Goal evidence; agreed improvements and follow-up exist; authority
  and mandatory verification preserved; parity/upgrade and recurrence supported.
  Project-owned template policy changes are fresh-init-only as explicitly disclosed.
- **O-034 7: Issue.** Scenarios and gates pass, but revised shipped README guidance
  contradicts supported upgrade behavior under TPL-02. All other scenario cells pass.

## Evidence and gates

- 2026-10-02 UTC: `make build && make test-full`, Go 1.26.2 linux/amd64, exit 0.
  Full host suite passed; linux, darwin and windows binaries built. Native Windows
  execution is not newly claimed by this documentation-only Objective.
- `git diff --check`: exit 0. All changed named paths exist; no phantom file.
- Scratch `/tmp/o034-probes.py` and Go overlay `/tmp/o034-overlay.json` run
  `go test -overlay /tmp/o034-overlay.json ./internal/init -run
  'TestO034IndependentUpgradeMatrix|TestSkillReview|TestUpgradeDelivers|TestUpgradeKeeps|TestScaffoldedPolicy|TestRepoAgents'
  -count=1 -v`: exit 0. Independent matrix reproduces the README counterexample
  while verifying backup, tracked preservation, dry-run and repeat behavior.
  An initial scratch compilation failed for a missing harness import; corrected
  in scratch only and rerun successfully. No repository test was changed.
- Official health after full gate: first sandbox collection saved
  `sha256:0a0b4600e02ce1059744e9b0ce65bc105eb5d7743a649af1d0ae87d94d750a22`,
  with optional OSV DNS/socket failure. Authorized network retry saved the
  frontmatter snapshot and reports **Code Health does not block clearance**.
  Collection failure is not a measured vulnerability or bad code.
- Canonical seven-asset weight: base 82,210 B; T-100 83,819 B (+1,609);
  current with retrospective 84,726 B (+907 more). Necessary health/repair/
  platform/read/routing/recurrence guidance accounts for growth; bootstrap
  AGENTS and template Design duplication reductions offset part of it.

## Finding and materiality

I-123 violates TPL-02 within frozen R9: README:153 and 436–438 claim edited
skills remain untouched without force, while upgrade.go:487–495 replaces a
supported pre-manifest skill with a recoverable backup. Existing legacy tests
and independent matrix agree. Documentation repair only; retain runtime policy.

| Issue | Likelihood | Impact | Materiality | Recommendation |
|---|---|---|---|---|
| I-123 | Medium: older projects without provenance | Medium: owner unexpectedly loses active custom instructions, recoverable in .bak | Medium | Qualify both README claims; direct Issue repair |

## Observations and reconciliation

- Complexity remains Watch at CCN 20; coverage is about 86.6%, duplication about
  2.8%, tests have zero failures. Watch is an observation, not a blocker.
- I-122 is the intentionally separate generated-report wording follow-up; it is
  not a blocker of this review. No report product code was changed.
- Template CODE-01 says bring a Watch signal to its watch line, whereas task advice
  says Needs Attention. This advisory example is imprecise; it neither authorizes
  metric chasing nor blocks clearance. Record for future copy cleanup.
- Existing Design init prose still says R-001; README now correctly says G-001.
  That unrelated architectural prose predates this Objective and is not remediation
  required by this scope. Workflow additions fit the existing asset architecture;
  no new runtime component, lifecycle field or Goal completion semantics.
- Pre-existing untracked temporary coverage/JSON files were left alone. No owner
  validation or acceptance was invented. Owner completion follows a later CLEAR
  Check; completed Tasks remain done during the direct Issue repair.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [x] STYLE-03 **Test branches** — upgrade happy/conflict cases and independent legacy branch covered.
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [x] STYLE-07 **One source of truth** — canonical/scaffold byte parity; deliberate independent skill policy copies retained.
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data**
- [x] STYLE-10 **Small diffs**

## Handoff

Repair I-123 directly without retreating any completed Task, then request a fresh
independent recheck of this frozen scope. No further implementation change is
required by this Check.
