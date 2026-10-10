---
id: C-971
scope: {kind: objective, id: O-039}
result: NEEDS WORK
checked_by: {role: checker, session: check-o039-20261010}
executed_session: executor-o039-unrecorded
checked_at: '2026-10-10T07:05:00Z'
health_snapshot: sha256:923e11ac5301f02309c548bfc2f068f0f234dac5bd7f87ec570adf009af10217
unmet: [O-039-SC3, O-039-SC4, T-114-DW1, T-114-DW2]
reviewed:
  base_commit: 33c70a7
  head_commit: 33c70a7
  files:
    - AGENTS.md
    - templates/project-v2/AGENTS.md
    - agent-skills/savepoint-idea/SKILL.md
    - agent-skills/savepoint-design/SKILL.md
    - agent-skills/savepoint-task/SKILL.md
    - agent-skills/savepoint-check/SKILL.md
    - agent-skills/references/check-method.md
    - agent-skills/references/issue-capture.md
    - agent-skills/references/commands-and-procedures.md
    - templates/project-v2/agent-skills/
    - internal/init/agent_skills_test.go
    - internal/init/template_freshness_test.go
    - internal/init/upgrade_test.go
    - internal/init/v2_scaffold_test.go
    - .savepoint/.upgrade-manifest.yml
  dependencies: []
issues: [I-144]
supersedes: null
---

# C-971: O-039 Full Objective Check

Session is fresh (started after `/clear`; did not build T-112..T-114). Mode: Full. All O-039 work is uncommitted in the working tree on top of `33c70a7`; this Check reviews that working tree.

## Scope Lock

1. Criteria: O-039 Success Conditions SC1–SC6; T-112, T-113, T-114 Done When; Guardrails TPL-01, TPL-02, TPL-04, FS-02, POL-01, POL-02, TEST-01, TEST-02, TEST-03, TEST-08, STYLE-07, STYLE-10.
2. Changed files: listed in `reviewed.files`. Public entry points: the always-loaded `AGENTS.md`, the four skills, three references, the scaffold copies, and `upgrade-assets` delivery of them.
3. Relied-on behavior: `upgrade-assets` refreshing the managed block and unedited skills while keeping edited ones (FS-02).
4. Matrix axes: per shared rule × each of the 8 guidance files (home / pointer / restated); live vs scaffold parity; history and term phrases × each file; test pinning (positive + absence). External-boundary, side-effect, Unicode, numeric axes: not applicable (text-only change, no runtime code changed).
5. Admission boundary: text an agent loads in a session; paraphrases count only where the Task evidence claims removal.

## Gate Evidence

- `make build && make test-full`: exit 0, 2026-10-10T06:53:44Z–06:54:01Z, go1.26.2 linux/amd64, HEAD 33c70a7 plus working tree. No FAIL lines.
- `git diff --check`: clean. `./savepoint doctor`: ALL CLEAN. `./savepoint resume`: loads, Next was this Check.
- Code Health: `savepoint health check O-039 .` → snapshot `sha256:923e11ac…0217` (created); "Code Health does not block clearance"; all five instances optional, no finding.

## Coverage

| Criterion | Result | Evidence |
|---|---|---|
| SC1 / T-112: no repeated section above the block | Proven, with observation O1 | Above-block sections are Build, Codebase Map, Runtime Gates, Building `savepoint`, Reporting to the Owner. Removed sections re-diffed: all now live only in the managed block. |
| T-112: managed block bytes unchanged during T-112 | Proven (later refreshed by T-113/T-114 via `upgrade-assets`, as planned) | Live block inner text is byte-identical to `templates/project-v2/AGENTS.md` (`diff` → identical). |
| SC2 / T-113: one home for Goal Context and waiver/dependency/Objective-coverage rules | Proven, with observation O2 | Idea, Design, Task, Check say "Apply AGENTS.md's Required Goal Context" plus own consequence; waiver `clear`/`accepted` meaning only in Verification Policy; `TestSharedVerificationRulesHaveOneHome` asserts home and absence. |
| T-113: `health report` human-only in template; repo block matches | Proven | Template CLI Rules carries it; block parity identical. |
| T-113: upgrade delivers revised skills, keeps edited ones | Proven | `TestUpgradeDeliversRevisedSkillsToUneditedProjects` (`upgrade_test.go:1572`) plus existing edited-skill test; manifest hashes updated for four skills. |
| SC3 / T-114 DW1: no V1/V2 narrative outside Required Goal Context; removed passages listed accurately | **Issue (I-144)** | `savepoint-design/SKILL.md:70` still says "V2 Task Markdown", which T-114 evidence lists as removed. Remaining `migrate`/V1 text in `commands-and-procedures.md` and idea's one-line "Legacy input requires `savepoint migrate` first" are live procedure, not narrative — accepted. |
| SC4 / T-114 DW2: "owner" for decision-maker, "Goal" for record | **Issue (I-144)** | Managed block Router Selection: "unless the owner names a Release as part of an explicit Goal choice". "user" hits remaining are the product's users, `## User Check` headings, and "user-authored files" — correct. "Converted V1 Releases" names the V1 record kind — acceptable. |
| SC5 / TPL-01: live and scaffold byte-identical; tests pin; gate passes | Proven | `diff -r agent-skills templates/project-v2/agent-skills` → only live-only `bubbletea-tui-design`. Full gate passes. Tests do not pin the two I-144 phrases. |
| SC6: token weight before/after | Proven | Recorded per Task. Independent re-measure of AGENTS.md + 4 skills + 3 refs (`wc`): before (HEAD) 16400 words / 110619 bytes, after 14213 / 96365 (−13%). No increase. |
| O-034 scenario re-walk (T-113, T-114) | Proven by reading | Scenarios 1–5 recorded per Task; my read of the revised text agrees. |

Adversarial pass: a skill that now only points to AGENTS.md could lose a role consequence. Checked each pointer: Check and Task keep "never picks or creates a Goal"; Check keeps "Never write a Check to stand in for a waiver"; Task keeps the `check_waiver` shape; Design keeps "never substitute a Task-only result or waiver". No rule change found.

## Materiality

| Issue | Likelihood | Impact | Materiality | Recommendation |
|---|---|---|---|---|
| I-144 two wording leftovers | High (always loaded) | Low (wording only; no wrong action results) | Low | Fix now as one small direct repair: two phrase edits, mirror, refresh block, add two absence assertions. |

## Owner Validation Still Needed

- T-113 declares `owner_validation.required: true` but records no owner acceptance (`accepted_check` is empty); it was closed with a Task-check waiver. Under Closure Rules this is the owner's call; the Objective cannot close by acceptance on T-113 until the owner records it or decides it is not needed.

## Observations (non-blocking)

- O1: "Reporting to the Owner" above the block is byte-identical to the same section inside the managed block, so it loads twice. The owner-confirmed design listed it as repository-specific, so it is not an Issue; dropping the above-block copy saves ~60 words.
- O2: Objective Check coverage ("covers every owned Task, including waived …") is still phrased in check-method Depth and in the Check skill's Closure Rules, besides the Verification Policy. These read as the Check role's own job, so not counted as restatement.
- O3: Required Goal Context says "Converted V1 Releases keep their R-### identities" twice (two paragraphs apart). Could be folded into the I-144 repair.
- O4: "After creating or renaming an identity-bearing record, run `savepoint resume`" appears in CLI Rules, Design, issue-capture and Check step 4 — outside O-039's named rules; candidate for O-040.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [x] STYLE-03 **Test branches**
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [ ] STYLE-07 **One source of truth** — "Reporting to the Owner" is duplicated in `AGENTS.md` (above and inside the block); see O1.
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data**
- [x] STYLE-10 **Small diffs**
