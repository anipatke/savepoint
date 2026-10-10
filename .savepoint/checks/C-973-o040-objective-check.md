---
id: C-973
scope: {kind: objective, id: O-040}
result: CLEAR
checked_by: {role: checker, session: check-o040-20261010}
executed_session: executor-o040-unrecorded
checked_at: '2026-10-10T07:25:00Z'
health_snapshot: sha256:4a8798f0799b9086d0d6e17bcdcbf4392981f406b1f2aea5a114d3ea47fb77dc
reviewed:
  base_commit: a2e343b
  head_commit: a2e343b
  files:
    - AGENTS.md
    - templates/project-v2/AGENTS.md
    - agent-skills/references/check-method.md
    - agent-skills/references/issue-capture.md
    - agent-skills/savepoint-idea/SKILL.md
    - agent-skills/savepoint-design/SKILL.md
    - agent-skills/savepoint-task/SKILL.md
    - agent-skills/savepoint-check/SKILL.md
    - templates/project-v2/agent-skills/
    - internal/init/reference_navigation_test.go
    - internal/init/agent_skills_test.go
    - internal/init/template_freshness_test.go
    - internal/init/v2_scaffold_test.go
    - .savepoint/.upgrade-manifest.yml
  dependencies: []
issues: []
supersedes: null
---

# C-973: O-040 Full Objective Check

CLEAR. Fresh session (started after `/clear`; did not build T-124 or T-125). Mode: Full. All O-040 work is uncommitted in the working tree on top of `a2e343b`. Both Tasks are `done` with owner Task-check waivers (board-owner, 2026-10-10T07:13:17Z and 07:15:21Z); this Check reviews both.

## Full Check Progress Checklist

- [x] Establish Scope
- [x] Freeze The Check Scope
- [x] Turn Acceptance Into Invariants
- [x] Build The Mandatory Coverage Matrix
- [x] Finite External-Boundary Matrix — not applicable: no server, subprocess, or provider in scope (text and content tests only; upgrade-assets code unchanged).
- [x] Workflow And Side-Effect Check Lock — upgrade-assets delivery only; code path unchanged, covered by the two upgrade tests.
- [x] Matrix Completion Lock
- [x] Parallel Planning Advice — not a feature of this scope; no finding.
- [x] Perform The Adversarial Pass
- [x] Re-check After Remediation — not applicable: initial Check.
- [x] Verify File Reality
- [x] Verify Evidence And Gates
- [x] Collect Code Health Evidence
- [x] Complete The Issues Pass
- [x] Summarize Materiality
- [x] Review Code Style

## Scope Lock

1. O-040 Success Conditions SC1–SC7; T-124 and T-125 Done When; Guardrails TPL-01, TPL-02, TPL-04, FS-02, POL-01, POL-02, TEST-01, TEST-02, TEST-03, TEST-08.
2. Changed files listed in `reviewed.files`; public entry points: the four skills, two references, the managed AGENTS.md block, `savepoint upgrade-assets`.
3. Relied-on behaviour: upgrade-assets refresh of managed skills, references and block (code unchanged).
4. Matrix axes below; external-boundary and re-check cells not applicable.
5. Issue admission: violation of a named criterion or guardrail, reproducible through the shipped text or tests.

## Coverage Matrix

| Criterion | Check | Independent evidence | Result |
|-----------|-------|----------------------|--------|
| SC1 one-hop audit | Grepped every file named by the three references | Only AGENTS.md (always loaded), `.savepoint/Guardrails.md`, `config.yml`, `Health-Check.md` (project files) and `check-method.md` (named directly by Design and Check). No two-hop chain. Audit table present in T-124. | Proven |
| SC2 contents lists | Python oracle comparing `**Contents**` items with fence-aware `##` headings | check-method 16/16, issue-capture 7/7, equal and in order, before first `##` | Proven |
| SC3 Full checklist | Python oracle: every step names an existing `##`/`###` heading, in file order; only `Quick Check Procedure` and the checklist heading omitted | Ordered and complete; Full bullet points to it | Proven |
| SC4 write → resume → fix | Read CLI Rules; grep for restatements | One paragraph in CLI Rules (live and template); each skill ends with the pointer; Design and Check restatements removed | Proven (see observation 1) |
| SC5 Next words | Compared descriptions with Workflow line | design `Plan`/`Replan`; task `Start`/`Build`/`Test`/`Pick a Task in`/`Fix`; check `Check`/`Assess`; idea keeps "router state is idea"; lengths 236/342/332/216 < 1024 | Proven |
| SC6 parity, block, upgrade, tests | `diff -r` skills vs template; managed block vs template; upgrade tests | Byte-identical (only `bubbletea-tui-design`, not a scaffold skill, differs); block matches template; upgrade tests pass | Proven |
| SC7 scenarios and weight | Read both Tasks' evidence | Five scenarios re-walked per Task; weight +258 and +72 words, each justified by name | Proven |
| Rule content unchanged | Read full diff | Additions plus the Full-bullet pointer, two removed restatements now covered by CLI Rules, and descriptions | Proven |

Mutation probes (in a scratch copy): deleting one contents item, removing `Fix` from task's description, and adding a restatement to idea each made the matching new test fail.

## Gates

- `git diff --check`: exit 0.
- `make build`: ok. `make test-full`: exit 0 (go1.26.2 linux/amd64, 2026-10-10T07:16Z), including the three cross-platform builds.
- `savepoint health check O-040`: snapshot `sha256:4a8798f0…77dc` created; "Code Health does not block clearance"; all five instances optional, no blocking finding.

## Guardrails

TPL-01 parity confirmed; TPL-02 text matches behaviour; TPL-04 and TEST-03 upgrade tests deliver revised files to unedited projects, upgrade code unchanged so existing preservation tests still apply; FS-02 AGENTS.md diff lies inside the managed block; POL-01/02 no new rule defined; TEST-01/02 named tests with mutation-confirmed failure paths; TEST-08 full gate passed.

## Issues

None. No materiality actions are required.

## Owner Validation Still Needed

Both T-124 and T-125 declare `owner_validation.required: true` with no `accepted_check`. The owner's User Checks (contents lists and checklist; one CLI Rules rule, pointers, and description words) are still to be recorded against this Check before the Objective closes.

## Observations (non-blocking)

1. `agent-skills/references/issue-capture.md:76` (and its template copy) still carries a narrower older form of the rule: "After creating or renaming an Issue or another identity-bearing record … run `savepoint resume` …". It does not conflict with CLI Rules, and no criterion covers references (T-125 tested skills only), but it keeps a second home for the rule against the Goal's "each rule lives in one place". Candidate one-line follow-up: replace it with the CLI Rules pointer.
2. T-124 evidence says `issue-capture.md` has 9 `##` sections; it has 7. The contents list itself is correct.
3. `TestSkillDescriptionsNameTheirNextWords` uses a hardcoded word list, so a new Workflow word would not be caught automatically; its `len(want) == 0` guard can never fire.
4. O-040 `status` is still `planned` though both Tasks are `done`.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [x] STYLE-03 **Test branches**
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [ ] STYLE-07 **One source of truth** — `issue-capture.md:76` still restates the resume rule (observation 1).
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data**
- [x] STYLE-10 **Small diffs**
