---
id: T-099
title: Review the skills against what the Code Health Goal taught us
objective: O-034
status: done
depends_on: []
complexity_tier: spike
complexity_reason: Evidence-gathering review whose findings decide what the next Task changes.
owner_validation:
    required: true
    accepted_check: ""
planned_by: {role: planner, session: plan-o034-20261003}
check_waiver:
    task: T-099
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-02T20:07:23Z"
---

# Review the skills against what the Code Health Goal taught us

## Outcome

A findings list in this Task's Technical Evidence classifies every candidate improvement to the four workflow skills, their three shared references, and AGENTS.md routing guidance. Each finding is marked **apply**, **follow-up Issue**, or **no change**, cites evidence from R-007's records, and gives its reason.

## User Check

Read the findings list. Each line names the skill or reference it concerns, the R-007 evidence behind it, and an apply, follow-up, or no-change decision you can agree or disagree with in one pass.

## Done When

- Every one of the seven skill and reference files and the AGENTS.md routing guidance was reviewed on six axes: clarity, duplication, context and token cost, routing, handoffs, and consistency with implemented behaviour. A file with no finding says so.
- Evidence is drawn from R-007 and cited by ID: Tasks that returned REPLAN REQUIRED, Objective Checks that returned NEEDS WORK and the rechecks they caused (especially O-031 and O-032), Issues I-084..I-121 where a workflow cause is visible, and the Objective's Lessons Carried In.
- Each finding records its decision. **Apply** findings name the exact file and the intended change in a sentence. **Follow-up** findings name why they need a separate product or architecture decision. **No change** findings give their reason.
- The Code Health advice rule from O-034's Confirmed Design is included as an apply finding with its proposed wording, and the report next-step wording is included as a follow-up.
- Every apply finding is checked against the boundary in the Objective's Success Conditions: owner authority, checker independence, acceptance evidence, mandatory Objective Checks, and project-owned gates are not weakened.
- The scenario walkthroughs named in O-034's Confirmed Design are listed with the expected behaviour of the current text, as the baseline for the next Task.
- No skill, reference, template, or code file is changed by this Task.

## Context Files

`.savepoint/objectives/O-034-skills-optimisation/Objective.md`; `AGENTS.md`; `agent-skills/savepoint-idea/SKILL.md`; `agent-skills/savepoint-design/SKILL.md`; `agent-skills/savepoint-task/SKILL.md`; `agent-skills/savepoint-check/SKILL.md`; `agent-skills/references/check-method.md`; `agent-skills/references/issue-capture.md`; `agent-skills/references/commands-and-procedures.md`. R-007 evidence is read by targeted search of `.savepoint/objectives/O-026-*` through `O-036-*` Task files, `.savepoint/checks/` records scoped to those Objectives, and `.savepoint/issues/I-084-*` through `I-121-*`. Read summaries and findings sections, not whole records.

## Design References

O-034 Confirmed Design and Lessons Carried In; Design section 1 (token-efficiency principle).

## Guardrails

TPL-02, POL-01, POL-02, STYLE-07.

## Implementation Plan

1. Tally R-007's REPLAN REQUIRED Tasks, NEEDS WORK Checks, and Issues. For each, note whether the skill text helped cause the miss.
2. Read each skill and reference once, recording findings on the six axes.
3. Merge the two lists into findings with decisions, then list the scenario baselines.

## Boundaries

Review only; nothing is edited. No new lifecycle states, Goal-level Check, or weaker verification is proposed as an apply finding.

## Technical Verification

Research Task: no gate run is needed because no code or asset changes. The findings list is the decision deliverable, reviewed by the owner before the next Task starts. A requested Task Check applies `agent-skills/references/check-method.md`.

## Technical Evidence

Executor evidence (not a Check, not CLEAR). Review only; no skill, reference, template or code file was changed. Files read: the Context Files (Objective, AGENTS.md, four skills, three references) and, by targeted search, R-007 Task frontmatter and Technical Evidence (O-026..O-036), Check records C-937..C-960 (result, closure map, observations), and Issue titles I-084..I-121. Extra reads: `templates/project-v2/AGENTS.md` and `internal/codehealth/report.go`, `dashboard_copy.go`, `.savepoint/health/report.md` (to confirm what the Code Health report tells an agent), `internal/init/*_test.go` names (to see what pins the skills). Files changed: this Task and the router/Objective status for starting it.

### R-007 tally

- **REPLAN REQUIRED:** two Tasks. T-050 (O-026) stopped because its dependency was met but the provider catalogue was still awaiting owner confirmation; T-054 (O-027) stopped because default thresholds, the material-decline rule and severity inputs were undefined. Both are decisions the owner had not made, not skill-text failures. The skill's stop-and-hand-back rule worked as written. A replan recorded only in the body does not move the router (T-050 says so in its body); the skill already states this.
- **Objective Checks:** Nine Objectives (O-026..O-032, O-035, O-036) needed initial and recheck chains; 15 Checks returned NEEDS WORK before the final CLEAR (C-937, 939, 941, 942, 944, 946, 948, 949, 950, 952, 953, 955, 957, 958, 959). Longest chains: O-031 four Checks, O-032 four, O-028 three, O-035 three. Issues I-084..I-121 are 38; none came from a skill-text cause. Most are product defects the Full Check found in code the owner had accepted without a Task Check.
- **Task Checks:** 48 of 53 R-007 Tasks carry an owner waiver. The optional Quick mode is almost never used; the Full Objective Check carried the verification.
- **Repair rounds:** five direct-repair rounds failed their first recheck because the repair broke a neighbouring case or missed part of the Issue: I-093 (rename repair leaked between two instances, C-942), I-106 (Linux repaired, native Windows unverified, C-949/C-950), I-107 and I-112 (still open, C-953), I-121 (truncation still exposed the credential, C-958), I-118 (replacement repair lost the documented damage boundary and the measured cost, C-959). C-959 reached the convergence limit and stopped for the owner, as `check-method.md` requires.
- **Pasted Next versus router state:** C-941 and C-957 both note the router said `task` with no selected Task while the owner pasted `Check O-0xx`, and both resolved it by following the pasted line.
- **Extra-read noise:** at least four Task records (O-031, O-035, O-032) log `AGENTS.md` or `agent-skills/savepoint-task/SKILL.md` as extra reads, the latter because the skill tool reported it unknown and AGENTS.md says to read the file directly.
- **Complexity lesson:** O-032's Lessons Carried In; C-957 observation "most complexity improvement is production refactoring; residual aggregate CCN includes test functions and is not itself an opt-in health blocker"; C-960 closed with the Watch band still populated. The report today says "Investigate each signal ... propose a fix, then apply it" and "aim for 10 or less" (`internal/codehealth/report.go`, `dashboard_copy.go`).

### Findings

Decision key: **apply** (T-100 edits it), **follow-up** (Issue, needs a product or architecture decision), **no change**. Every apply finding was checked against owner authority, checker independence, acceptance evidence, mandatory Objective Checks and project-owned gates; none weakens them, and the one that touches Checks (A2) only narrows what a checker may raise as an Issue, never what it may clear.

| # | File | Evidence | Decision | Finding and reason |
|---|------|----------|----------|--------------------|
| A1 | `savepoint-task` | O-032 Lessons; C-957/C-960 observations | **apply** | Add a short section "Acting On A Code Health Report" after Verification Gates: bring a signal that Needs Attention back to the watch line, not the aim; prefer production code that is risky or often changed over tests; leave flat dispatch tables alone; record what remains in the Task evidence instead of silently stopping or continuing; narrowing the measured scope (for example excluding test files) is the owner's decision. Cites `.savepoint/health/config.json` thresholds by name, restates no Guardrail (POL-02). |
| A2 | `savepoint-check` | C-957 and C-960 observations; ~280 signals above the aim | **apply** | Add two sentences to Code Health Evidence: a signal between the aim and the watch line is an observation, never an Issue or a reason for `NEEDS WORK`; list what remains in the Check body. A blocking verdict still prevents `CLEAR`. |
| A3 | `savepoint-task` | four Task evidence records | **apply** | Read section: say the router, `AGENTS.md`, this skill (read directly when the skill tool cannot find it) and the Guardrail IDs the Task names are not extra reads. Workflow step 5 says "editing" where Extra Reads says "reading"; make step 5 say "reading or editing". |
| A4 | `agent-skills/references/issue-capture.md` | C-942, C-949, C-953, C-958, C-959: five repairs failed first recheck | **apply** | In Out-Of-Scope Repair, before `repair_attempted` is recorded, require the executor to re-run the Issue's Proof Needed scenario and the neighbouring frozen cells the repair could change, record each result in the history, and mark any platform or cell it could not run as unverified. The Check still independently proves the repair; this does not claim `verified` (checker independence kept). |
| A5 | `savepoint-check`, `savepoint-task`, `savepoint-design` Trigger | C-941, C-957 | **apply** | Trigger lines say "router `state` is `check`" while AGENTS.md says the `Next` line's first word picks the skill and a pasted line is the selection. Change each Trigger to name the Next line (`Check`; `Start`/`Build`/`Test`; `Plan`/`Replan`) first, router state second. |
| A6 | `savepoint-design` step 10 | C-949, C-950, C-957 (I-115): native Windows evidence was the sole blocker for two extra rechecks | **apply** | When a Task or Objective touches processes, paths, signals or file replacement, name in its Technical Verification the platform evidence the Full Check needs (for example the native Windows CI job) and who produces it, so it exists before the Check starts. It names the evidence, and the project's own gates still decide. |
| A7 | `savepoint-design`, `savepoint-task`, `savepoint-check`, `AGENTS.md` | key names appear ten times | **apply (low)** | Drop the board key names (Space, Backspace) from the skills; they carry no agent action. Keep one mention in `issue-capture.md` and in AGENTS.md. Saves roughly eight lines. |
| A8 | repo `AGENTS.md` | not shipped; loaded every session | **apply** | Remove the duplicate "V2 Routing" table (same table as Skill Activation, plus a statement that E47/E50 "ships/activates") and the "Legacy V1 compatibility (not active)" section. `templates/project-v2/AGENTS.md` has neither, so parity is unaffected. About 1.9 KB per session. |
| F1 | `internal/codehealth/report.go`, `dashboard_copy.go` | O-032 Lessons; report text above | **follow-up** | The report brief says "Investigate each signal ... propose a fix, then apply it" and the complexity line says "aim for 10 or less", so an agent chases the aim. Naming the watch line as the stopping point is generated-report product work, outside this review. |
| F2 | `savepoint-check` / `check-method.md` | 38 Issues; 15 NEEDS WORK Checks | **no change** | The review looked for a skill-text cause behind the long recheck chains and found none: each NEEDS WORK reproduced real defects, and C-959 applied the convergence limit correctly. Longer chains came from repairs failing (A4), not from Check procedure. |
| N1 | `check-method.md` (405 lines) | at most 5 of 53 Tasks lacked a waiver | **no change** | Largest context cost (about 21 KB), loaded in full for every Check. Full uses nearly every section, and the Full-only sections are already labelled for Quick. Splitting it would change the "loaded in full" rule and every test and skill that pins it, which is a structure decision. |
| N2 | `savepoint-design` Verification Contract, Required Goal Context; waiver paragraph repeated in four files | no R-007 case where the copies diverged | **no change** | Duplication is real (the waiver rule appears in five places) but every copy agrees today and tests pin them. The tradeoff favours leaving them: each skill is loaded alone, and a missing copy would hide a rule. Revisit only if a divergence appears. |
| N3 | `savepoint-idea` | R-007 had no Idea-state work | **no change** | Not exercised; the file is short (86 lines) and consistent with AGENTS.md. |
| N4 | `references/commands-and-procedures.md` | read only for migration reconciliation; none in R-007 | **no change** | Names `internal/doctor` and `internal/migrate` functions that downstream projects lack, but it is read only while reconciling a migrated project. |
| N5 | `savepoint-task` REPLAN REQUIRED | T-050, T-054 | **no change** | Stopped correctly on missing owner decisions; the frontmatter-block rule is already stated. |
| N6 | `bubbletea-tui-design` | — | **out of scope** | Not part of the review per the Confirmed Design. |

The recurring-review mechanism (planner adds a final retrospective Objective to each Goal) belongs to T-101 and is not an apply finding here.

### Scope extension: templated project documents (owner request)

The owner asked to include the documents `savepoint init` scaffolds. Reviewed: `templates/project-v2/.savepoint/Design.md`, `Guardrails.md`, `Idea.md`, `config.yml`, `router.md`, `releases/G-001-first-goal/Release.md` and the template `AGENTS.md`. Extra read, logged. `Idea.md`, `config.yml`, `Release.md` and the template `router.md`: no finding.

| # | File | Evidence | Decision | Finding and reason |
|---|------|----------|----------|--------------------|
| T1 | template `Design.md` | the waiver rule now has six copies | **apply** | The "Verification Contract" section is policy in a file the skill says describes implemented reality, and it is not among the sections `savepoint-design` lists for Design. It also carries board key detail (Space auto-records the waiver). Replace it with one line pointing to AGENTS.md "Verification Policy", which the template already carries. Owners who edited it keep their bytes through the upgrade path. |
| T2 | template `Guardrails.md` | O-032 Lessons; A1 | **apply** | Add one Guideline-severity example rule, for example CODE-01, "Complexity: fix a Watch signal back to its watch line; the aim is not a target", so the project owns the stopping point (Guardrails owns severity and exceptions) and A1/A2 can cite the rule ID instead of restating prose. The advice stays advisory and cannot block a Check. |
| T3 | template `Guardrails.md` TEST-04 | duplicates the waiver rule | **no change** | It is a worked example the owner is told to replace, and it is one line. |
| T4 | template `AGENTS.md` | A3, A5, A7 | **apply** | Mirror the key-name trim (A7) and keep it free of the repo-only sections; it already lacks V2 Routing and Legacy. |
| T5 | template `Design.md` Current Technical State | adoption guidance | **no change** | Clear and strict about verified-only claims. |

Boundary check: T1 removes a duplicate, not the rule; T2 adds an advisory rule only. Neither weakens owner authority, Checks or gates. T-100 grows from seven parity assets to include these template files; the Objective Confirmed Design review scope was amended to match.

### Per-file review (six axes: clarity, duplication, token cost, routing, handoffs, consistency)

- `savepoint-idea`: no finding (N3).
- `savepoint-design`: routing (A5), handoffs (A6), duplication (A7, N2). Largest skill at 290 lines, but its template sections are used while planning.
- `savepoint-task`: clarity and token cost (A1, A3), handoffs (A4 by pointer), routing (A5).
- `savepoint-check`: routing (A5), clarity (A2, A7). Closure rules are long but consistent with AGENTS.md and the runtime gate.
- `check-method.md`: token cost only (N1); no clarity, routing or consistency finding.
- `issue-capture.md`: handoffs (A4).
- `commands-and-procedures.md`: no change (N4).
- `AGENTS.md`: duplication (A8); routing text otherwise matches the skills, and `Worktree Lanes` is consistent with the skills' references to it.

### Scenario baselines (behaviour of the current text; T-100 re-walks these)

1. **Next to skill routing.** `savepoint resume` prints `Next: Check O-0xx` while router `state: task`. AGENTS.md routes by the first word to `savepoint-check`; the check skill's Trigger says state must be `check`, and the pasted line is accepted per AGENTS.md step 1. Current outcome: works, via two readings that need reconciling (A5).
2. **REPLAN REQUIRED re-entry.** Executor writes the `replan:` frontmatter block, leaves `status` and `stage`, and stops. Resume prints `Replan`; `savepoint-design` resumes at step 3, updates Design or the Objective, and restores the Task. A body-only note leaves the router on the Task. Current outcome: correct.
3. **Waived Task Check reaching the Full Objective Check.** Task evidence carries a `check_waiver:` block; the owner marks it done; a `requires: clear` dependency is satisfied and `requires: accepted` is not; the Full Check reviews every owned Task including waived ones. Current outcome: correct and stated consistently in all five places.
4. **Issue-only repair.** Router selects `I-###` alone; the executor repairs, records `repair_attempted`, moves the router to the single linked Objective (or clears `issue` only), preserves `release:` and runs `savepoint resume`. Current outcome: correct; the repair does not yet require reproducing the Issue's Proof Needed (A4).
5. **Acting on an advisory Code Health report.** Executor or checker sees a Watch signal at the aim-to-watch band. Today the skills say only that warnings never create Issues automatically; the generated report tells the agent to fix each signal. Current outcome: an agent may chase the aim (A1, A2, F1).

### Limitations

No live agent run; findings rest on the written records. Counts of Task evidence mentions were by targeted search, not a full read, so the extra-read figure is a lower bound. No owner review has happened yet.

## Drift Notes

None expected.
