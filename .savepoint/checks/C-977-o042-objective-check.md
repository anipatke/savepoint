---
id: C-977
scope: {kind: objective, id: O-042}
result: CLEAR
checked_by: {role: checker, session: check-o042-20261010}
executed_session: executor-t126-unrecorded
checked_at: '2026-10-10T08:24:00Z'
health_snapshot: sha256:e935874630868e27ad470371dc4a32776ea65ee087e89b84358757c8e4ff1742
reviewed:
  base_commit: 3596b2c
  head_commit: 3596b2c (uncommitted working tree)
  files:
    - agent-skills/savepoint-design/SKILL.md
    - templates/project-v2/agent-skills/savepoint-design/SKILL.md
    - internal/init/agent_skills_test.go
    - internal/init/upgrade_test.go
    - .savepoint/releases/G-002-skills-optimisation/Release.md
    - .savepoint/objectives/O-042-skills-optimisation-workflow-retrospective/Objective.md
    - .savepoint/objectives/O-042-skills-optimisation-workflow-retrospective/tasks/T-126-write-the-end-of-goal-review-as-a-note-not-an-objective.md
  dependencies: []
issues: []
supersedes: null
---

# C-977: O-042 Full Objective Check

CLEAR. Technical clearance only. O-042 still needs the owner's T-126 User Check acceptance before it can close (see Owner Validation Still Needed).

**Independence.** This session did not build T-126; the skill, scaffold and test changes were made in a separate executor session that recorded no session ID. Disclosure: this same session earlier acted as planner for O-042 (Objective text, T-126 draft) and wrote G-002's `## Workflow Review` note in `Release.md`. That note is a planner record, not T-126 output; it is judged here only for existence and loadability, not for the quality of its own conclusions. The owner may ask for a fresh-session recheck if they want the note's content reviewed independently.

## Full Check Progress Checklist

- [x] Establish Scope
- [x] Freeze The Check Scope
- [x] Turn Acceptance Into Invariants
- [x] Build The Mandatory Coverage Matrix
- [x] Finite External-Boundary Matrix — not applicable: no server, subprocess, provider or network boundary in scope.
- [x] Workflow And Side-Effect Check Lock — not applicable beyond the existing `upgrade-assets` delivery path, covered by row M5.
- [x] Matrix Completion Lock
- [x] Parallel Planning Advice — none recorded; not a feature of this scope.
- [x] Perform The Adversarial Pass
- [x] Re-check After Remediation — not applicable: first run.
- [x] Verify File Reality
- [x] Verify Evidence And Gates
- [x] Collect Code Health Evidence
- [x] Complete The Issues Pass
- [x] Summarize Materiality
- [x] Review Code Style

## Scope Lock

1. Criteria: O-042 Success Conditions 1–4; T-126 Done When 1–5. Guardrails TPL-01, TPL-02, TPL-04. Gate: `make test-full`.
2. Changed files: the seven listed in `reviewed.files`. Public entry points: the packaged `savepoint-design` skill text, `init` scaffolding, `upgrade-assets` refresh of managed skills.
3. Relied-on behaviour: the managed-asset manifest refresh in `internal/init` (unchanged), the Goal loader reading `Release.md` (unchanged).
4. Matrix rows M1–M8 below; no runtime state machine or text-width cells apply (guidance-only change, no rendered terminal output).
5. Materiality boundary: an Issue must violate a criterion above through the shipped skill text, scaffold, tests or upgrade path. Historical records (O-038, past Checks) and non-shipped maintenance notes are out of scope per O-042 Boundaries.

## Coverage Matrix

| Row | Invariant | Evidence | Result |
|---|---|---|---|
| M1 | Design skill no longer asks for a retrospective Objective | Diff of `agent-skills/savepoint-design/SKILL.md` §Goal Workflow Review; repo-wide search outside `.savepoint/`, `.git` and tests finds no shipped mention | Proven |
| M2 | New section names `## Workflow Review` in the Goal's `Release.md`, says when it is written, "not an Objective, Task or Check", "no change, because…", real change → Objective or Issue | Section text, lines 96–98 | Proven |
| M3 | Scope and evidence list, package-receiving rule, and "Goal is still complete when every member Objective is complete" retained | Unchanged sentences in §Goal Workflow Review and following paragraph | Proven |
| M4 | Scaffold copy byte-identical (TPL-01) | `cmp` identical; `TestV2SkillSetIsCompleteWithByteParity` PASS | Proven |
| M5 | `upgrade-assets` delivers the new rule to an existing project (TPL-04) | `TestUpgradeDeliversGoalWorkflowReviewRuleToDesignSkill` PASS (stale managed copy → `ActionUpdated`, new heading present) | Proven |
| M6 | Tests pin new wording and forbid the old | `TestSkillReviewRulesArePinned/design_writes_the_Goal_workflow_review_as_a_note` PASS. Independent mutation: appending "Add a retrospective Objective too." to the live skill made that subtest fail ("still contains \"retrospective Objective\"") and the parity test fail; file restored, SHA-1 `86c0753f…` before and after, parity re-confirmed | Proven |
| M7 | G-002 review recorded in `Release.md` and the Goal still loads | `## Workflow Review — 2026-10-10` present; `savepoint resume` strict-loads the index with it | Proven (existence and loadability only; see Independence) |
| M8 | Guidance describes actual behaviour (TPL-02): zero-Task Objectives never close | `resolveSelectedObjective` and `resolveObjectiveIntegrationRung` in `internal/data/next.go` route a zero-Task Objective to `NextPlanObjective`; no runtime change claimed or made | Proven |

## Adversarial Pass

- Bypass via another shipped copy: searched `agent-skills/`, `templates/project-v2/`, `AGENTS.md`, `CLAUDE.md` and non-test Go sources; no other statement of the old rule.
- Trigger reachability: after a Goal's last Objective closes, Next becomes the owner's Goal step; the next planning session happens when the owner picks or plans further work, and "or any session the owner asks" covers a project with no further planning. Reachable; see Observation 2.
- Loader tolerance: an extra `##` section in `Release.md` does not break strict loading (M7).
- Representation: scaffold and live copies cannot diverge silently (M4, M6 mutation).

## File Reality

Every file named in T-126 evidence exists with the described content. `.go-test.json.*` files at the repo root are pre-existing scratch output, not part of this scope.

## Evidence And Gates

- `git diff --check`: clean.
- `make build`: exit 0.
- `make test-full`: exit 0, run fresh 2026-10-10T08:20:56Z–08:21:13Z on `go1.26.2 linux/amd64`, base `3596b2c` plus the uncommitted O-042 changes; no FAIL lines; Linux, Darwin and Windows cross-builds completed.
- Focused: the three tests in M4–M6, `-count=1`, PASS.
- `quality_gates` in `.savepoint/config.yml` are all null; AGENTS.md Build defines the gates used above.
- Native Windows CI: not required; no platform-sensitive code touched (O-042 Readiness and Verification).

## Code Health

`savepoint health check O-042 .` created official snapshot `sha256:e935874630868e27ad470371dc4a32776ea65ee087e89b84358757c8e4ff1742`: "Code Health does not block clearance." Complexity, dependency vulnerabilities and duplication: no finding. Tests and coverage (both optional): stale evidence. That is a collection-freshness statement, not a code finding, and does not block.

## Issues Pass

| Criterion | Classification |
|---|---|
| O-042 SC1 / T-126 DW1 | Proven (M1, M2, M3) |
| O-042 SC2 | Proven (M3) |
| O-042 SC3 / T-126 DW2, DW3 | Proven (M4, M5, M6) |
| O-042 SC4 | Proven (M7) |
| T-126 DW4 | Proven (M1) |
| T-126 DW5 | Proven (gates) |
| TPL-01, TPL-02, TPL-04 | Satisfied |

No Issues. No materiality actions are required.

## Observations (non-blocking)

1. `GUIDANCE-COMPARISON.md` (repo root, tracked, not shipped) still cites "Design / Goal Workflow Retrospective" in its point-in-time comparison table. It describes itself as a historical maintenance comparison, so a stale name there is expected; update or leave as history at the owner's choice.
2. The design skill's Read section ends "Read nothing else", while the new review needs the finished Goal's Tasks, Checks and Issues. The review section names that evidence, which reads as permission, but a planner could skip it. Same exposure as the old rule; consider a pointer in Read if a future review is skipped.
3. Health tests and coverage reported stale after `make test-full`; the health instances may read report paths the full gate did not refresh.

## Owner Validation Still Needed

T-126 declares `owner_validation.required: true` and has no acceptance yet. The owner's User Check: open `agent-skills/savepoint-design/SKILL.md`, read §Goal Workflow Review, and compare with the G-002 note. Then record acceptance against C-977 and close O-042. Its Task-check waiver (recorded 2026-10-10T08:16:17Z) is not technical clearance; this Check covers T-126.

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
