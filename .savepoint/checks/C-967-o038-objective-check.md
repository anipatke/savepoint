---
id: C-967
scope: {kind: objective, id: O-038}
result: CLEAR
checked_by: {role: checker, session: check-o038-20261004-codex}
executed_session: exec-t111-2026-10-04
checked_at: '2026-10-04T05:37:39Z'
health_snapshot: sha256:f68d6f08e0ca2f5b7e82e980f3978437e63d8d6be5c9eb9afc8ebe03c01e84cf
reviewed:
  base_commit: 2527047
  head_commit: 80b8657
  files:
    - .github/workflows/ci.yml
    - .savepoint/Guardrails.md
    - .savepoint/objectives/O-038-worktree-goal-workflow-retrospective/Objective.md
    - .savepoint/objectives/O-038-worktree-goal-workflow-retrospective/tasks/T-111-run-ci-on-every-pushed-branch.md
  dependencies:
    - .github/workflows/publish.yml
    - Makefile
    - agent-skills/savepoint-design/SKILL.md
    - agent-skills/references/check-method.md
    - AGENTS.md
    - templates/project-v2/AGENTS.md
issues: []
supersedes: null
---

# C-967: O-038 Full Objective Check

CLEAR. This is a fresh checker conversation that implemented neither T-111 nor the retrospective. All owned Tasks (T-111) were reviewed, including its owner-recorded optional Check waiver. No implementation, Design, router, Task status, Objective status or owner acceptance was changed.

## Frozen scope lock

O-038 frozen initial scope lock
1. O-038 SC1–3; T-111 DW1–4; CFG-03, TEST-01/02/08/09, TPL-01/02, owner authority and advisory boundaries. Fresh full gate and official health.
2. Public surfaces: CI push and pull_request event filters, unchanged jobs/permissions, publish tag ownership; recorded retrospective and ARCH-05 policy addition. No runtime terminal overhaul promised.
3. GitHub event filtering -> job scheduling -> checkout/setup -> make ci or native full Go tests -> read-only results. Remote runner failures fail job; no changes to retry/cleanup or secret handling.
4. M1 branch names main, v2.20, never-listed, feature/deep/name, Unicode; tag exclusion; PR filter equality; entire workflow structural equality outside push branches. M2 retrospective Task/Check/Issue lesson attribution, no-replan claim, decisions/no-change rationale, canonical/scaffold equality, Design reconciliation. M3 waiver/file reality/full gate/native CI/current SHA/health.
5. Admit only changed or explicitly promised behavior. Existing terminal surfaces outside the lane fix are follow-up observations, not a retrospective implementation defect. GitHub unavailable evidence is unverified, not failed code.
Workflow: event parse/filter (no local write) -> jobs checkout/setup (runner workspace) -> gates (reports/builds; fatal errors) -> result (success only after all steps); independent oracle GitHub run API and workflow diff. Malformed workflow, refusal, timeout/cancel/non-success are unchanged GitHub-owned behavior, source trace only. Redirects and secret-safe errors N/A: no new client/network implementation. Terminal modes/text-width/state machine/serialization mutation N/A for one-line declarative event filter.

## Coverage matrix and invariants

| Row | Applicable cells | Expected / actual and evidence | Classification |
|---|---|---|---|
| M1 event filter | ordinary, previously listed, never-listed, nested and Unicode branch names | `branches: ['**']` includes all branch names, including `/`; documented GitHub filter semantics cover main, v2.20, never-listed, feature/deep/name and 実験/新しい. Independent YAML structural comparison proves the exact universal filter. No new branch was pushed as a probe. | Passed by declarative review and official oracle |
| M1 negative/bypass event | tag push; PR to listed/unlisted target | Branches-only filtering excludes tags; PR filter remains master/v2/v2.1. Publish still owns v* tags. Removing all filters would admit tags and using '*' would miss nested branches; neither bypass is present. | Passed |
| M1 complete workflow | jobs, runners, permissions, setup and gate commands | Python YAML comparison against f905f04 changes only push.branches before asserting entire workflow equality. Git diff is one line; both jobs execute their prior gates. | Passed |
| M2 evidence retrospective | T102–107, T108–110; C963–966; I131–140, I115, I074, I083 | Task evidence/drift supports no actual REPLAN REQUIRED; initial Checks show the defects and missing native proof; rechecks and resolved Issues support closure. Repeated branch-list restrictions and I136 terminal-control reproduction support the two decisions. | Passed |
| M2 decisions and no-change conclusions | planner contract, routing, advice/health/authority, policy addition | Final retrospective is present; bounded CI fix and ARCH-05 are recorded. No new lifecycle, field, Goal Task ownership or Code Health switch; optional advice and independent verification remain explicit. Existing routing explains the repaired Start-line lesson. | Passed |
| M2 shipped guidance and Design | canonical/scaffold Markdown pairs, root/scaffold guide, Design settings/parallel sections | All existing canonical/scaffold agent-skills Markdown pairs compare byte-identically. Guides retain owner/worktree boundaries; wording difference below is nonblocking. Design records settings and advice as separate from gates; CI change alters no runtime architecture. | Passed with observation |
| M3 Task evidence and file reality | DW1–4, waiver, owner completion | Named workflow/Objective files exist; commit diff confirms no other product/skill/template changes. T111 is owner-done with explicit task/reason/actor/time waiver. Its recorded build/fast evidence and native run are present. | Passed |
| M3 gates and attributable native proof | fresh full/build, whitespace, exact CI revision | Local fresh build/full exits 0; GitHub API reports push on v2.20 at 0ad598ce73ca39f15cf9fce610280691a66e9631 with both jobs successful. Current workflow and all code/test/gate inputs are unchanged from that revision; later changes are records only. | Passed |
| M3 health | official collection after full gate | Official snapshot named above: all five instances collected, no clearance blocker. | Passed |

All O-038 Success Conditions 1–3 and T-111 Done When 1–4 are Proven. The independent scenario is the structural YAML comparison against the predecessor, crossed with documented positive/nested and negative/tag event cases, rather than repeating an implementation test. This declarative change adds no Go parser, renderer, numeric boundary, mutable model or runtime state transition: those matrix axes are not applicable. Empty/malformed branch names are not supported Git refs, and malformed workflow handling is unchanged GitHub infrastructure.

## External boundary and workflow inventory

| Order | Operation / effect | Failure ownership and final state | Independent evidence |
|---|---|---|---|
| 1 | GitHub parses event and branch filter; schedules jobs | Unmatched tags/PR targets schedule no CI; malformed YAML is platform refusal. No local persistence. | Exact YAML plus official documented semantics |
| 2 | Each job checks out pushed SHA and sets up Go/Node | Startup/unavailable dependency/setup failure fails job before gate; runner workspace only. | Unchanged source ordering; successful run API steps |
| 3 | Linux make ci; Windows native full Go suite | Nonzero gate fails job; reports/build artifacts remain runner-local. No new retry or partial publication. | Both jobs and all steps successful |
| 4 | Publish run conclusion | Success visible only after completed steps; cleanup stays platform-owned. | ci completed 05:13:08Z; windows-tests 05:14:36Z |

Configured and actual targets match ubuntu-latest/windows-latest and the pushed SHA. Discovery/reuse, refusal, timeout/cancellation, malformed responses, retry and cleanup code are unchanged external-platform responsibilities, source-traced rather than newly fault-injected. Redirect behavior, secret handling and partial publication are not introduced by this CI trigger; publish.yml is unchanged and outside implementation scope. Read-only GitHub query initially failed under network restriction and succeeded after approved access; no failure was mistaken for bad code.

## Adversarial pass

Checked nested branch exclusion (`*` versus `**`), tag-admission bypass (omitting branch filters), PR broadening, job/permission drift, false native attribution and accidental scaffold changes. None occurs. The retrospective does not turn advisory manifests into gates or waive checks. ARCH-05 is a durable policy decision responding to I136, not a promise that this Objective repairs every preexisting terminal surface. C966's existing Task-card control observation remains separate follow-up, without widening this retrospective into a terminal overhaul.

## Evidence and gates

- Fresh `make build && make test-full`, 2026-10-04 approximately 05:29–05:34 UTC, Go go1.26.2 linux/amd64: exit 0; full host tests and linux/darwin/windows builds. Transcript: /tmp/o038-full-gate.log (scratch evidence, not a project artifact). No prior full result reused.
- `git diff --check`: exit 0. Configured lint/typecheck/test entries are null; project-owned Makefile gates above apply.
- Python/PyYAML structural equality and canonical/scaffold pair comparison: passed. Initial root-guide equality assertion exposed the observation below; inspection classified the difference rather than claiming exact guide equality.
- Independently verified [CI run 37179190078](https://github.com/anipatke/savepoint/actions/runs/37179190078): event push, branch v2.20, exact implementation SHA above, ci success and native windows-tests success. Workflow/code/tests unchanged since that run.
- Filter interpretation checked against [GitHub workflow syntax](https://docs.github.com/en/actions/reference/workflows-and-actions/workflow-syntax#onpushbranchestagsbranches-ignoretags-ignore) and its filter pattern table: universal double-star includes slash-bearing names; branches-only excludes tag events.
- `./savepoint health check O-038` after full gate first saved official b80adf1f783c1f154671fead5488cb15ba20e2ec9f91c64b53cf01764251a593 with optional OSV network collection failure. Required retry with approved network access saved the frontmatter snapshot; all five instances collected. Both immutable snapshots are retained; the successful one is Check evidence.

## Observations and materiality

No Issues; no materiality actions required.

The live CI run uses v2.20, which was previously listed. A never-listed branch was not pushed for this Check; universal coverage rests on the exact filter and official semantics, not a claimed observed remote event.

The root AGENTS.md managed block still says the planner shapes Tasks for lanes where practical; the scaffold says suggestions are optional when enabled. Root project rules and canonical skills explicitly preserve optional advice, so this wording difference does not contradict current behavior or block clearance. No guide repair was made during the Check.

## Code Style Review

Rules without a Go implementation change are satisfied by the bounded declarative change; no speculative code or abstraction was added.

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

## Owner handoff

Technical integration is CLEAR; the owner decides Objective closure. T-111's owner completion/waiver stands; this Check invents no owner validation or acceptance. No Check superseded and no Issue resolution needed.
