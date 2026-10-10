---
id: C-972
scope: {kind: objective, id: O-039}
result: CLEAR
checked_by: {role: checker, session: recheck-o039-20261010}
executed_session: executor-o039-unrecorded
checked_at: '2026-10-10T07:00:00Z'
health_snapshot: sha256:4373077fd83519565b9938bba847cd41672f622b390ce691c0bcc06fc198dfa8
reviewed:
  base_commit: 33c70a7
  head_commit: 33c70a7
  files:
    - AGENTS.md
    - templates/project-v2/AGENTS.md
    - agent-skills/savepoint-design/SKILL.md
    - templates/project-v2/agent-skills/savepoint-design/SKILL.md
    - internal/init/template_freshness_test.go
    - agent-skills/
    - templates/project-v2/agent-skills/
  dependencies:
    - .savepoint/checks/C-971-o039-objective-check.md
    - .savepoint/issues/I-144-o039-wording-leftovers-v2-task-markdown-and-release.md
issues: []
supersedes: C-971
---

# C-972: O-039 Full Objective Re-check

CLEAR. Fresh session (started after `/clear`; did not build T-112..T-114, did not write C-971, did not make the I-144 repair). Mode: Full re-check under C-971's frozen scope lock; no new axes. All O-039 work, including the I-144 repair, is uncommitted in the working tree on top of `33c70a7`.

## Closure Map

| Prior Issue | Result | Evidence |
|---|---|---|
| I-144 item 1: "V2 Task Markdown" in design skill | Closed | `savepoint-design/SKILL.md:70` live and scaffold now read "Prepare the complete Task Markdown as an ID-free draft"; `grep "V2 Task Markdown"` over both skill trees → no hits. |
| I-144 item 2: "names a Release" in Router Selection | Closed | `AGENTS.md:80` and `templates/project-v2/AGENTS.md:30` read "unless the owner names a Goal as part of an explicit Goal choice". |
| I-144 proof: tests pin both absences | Closed | `TestActiveGuidanceKeepsNoHistoryAndOneTermPerRole` (`template_freshness_test.go:553`, `:561`) asserts both. Mutation probe in a scratch copy: reinserting "V2 Task Markdown" into the scaffold design skill → FAIL at :554; reinserting "names a Release" into `AGENTS.md` → FAIL at :561. Scratch copy deleted. |
| C-971 O3 (duplicate "Converted V1 Releases") | Folded in | Now one occurrence in each AGENTS file (`AGENTS.md:126`). |

## Admission Ledger

| Re-check item | Prior claim | Frozen cell | Allowed result |
|---|---|---|---|
| Design skill phrase | I-144 item 1 | history/term phrases × design skill (live, scaffold) | Issue if present |
| Router Selection phrase | I-144 item 2 | history/term phrases × AGENTS block (live, template) | Issue if present |
| Test pinning | I-144 proof | test pinning (absence) | Issue if absent |
| Live vs scaffold parity | SC5 | live vs scaffold parity | Issue if differs |
| Every other C-971 row | C-971 Proven | same rows | Issue only on regression |

## Gate Evidence

- `make build && make test-full`: exit 0, 2026-10-10T06:58:32Z–06:58:49Z, go1.26.2 linux/amd64, HEAD 33c70a7 plus working tree. No FAIL lines; cross-platform builds ran.
- Focused: `TestActiveGuidanceKeepsNoHistoryAndOneTermPerRole`, `TestSharedVerificationRulesHaveOneHome` → PASS.
- `git diff --check`: clean. `./savepoint doctor`: ALL CLEAN. `./savepoint resume`: loads; Next was "Record the Objective O-039 integration Check."
- Code Health: `savepoint health check O-039 .` → snapshot `sha256:4373077f…dfa8` (created); "Code Health does not block clearance"; all five instances optional, no finding.

## Coverage (C-971 rows re-checked)

| Criterion | Result | Evidence |
|---|---|---|
| SC1 / T-112 | Proven | Unchanged since C-971; managed block inner text byte-identical to template (`diff` → identical). |
| SC2 / T-113 | Proven | `TestSharedVerificationRulesHaveOneHome` passes; repair touched no home/pointer text. |
| T-113 upgrade delivery | Proven | Full gate passes, including `TestUpgradeDeliversRevisedSkillsToUneditedProjects`. |
| SC3 / T-114 DW1 | Proven | I-144 item 1 closed; T-114 evidence's "V2 Task Markdown removed" claim is now accurate. |
| SC4 / T-114 DW2 | Proven | I-144 item 2 closed. |
| SC5 / TPL-01 | Proven | `diff -r agent-skills templates/project-v2/agent-skills` → only live-only `bubbletea-tui-design`; block parity identical; tests pin both phrases. |
| SC6 | Proven | Repair is a two-word removal and a one-sentence fold; no increase over C-971's measured −13%. |
| O-034 scenario re-walk | Proven | Repair text does not touch the walked scenarios. |

Adversarial pass (remediation paths only): the phrase change in Router Selection keeps the rule's meaning (`release:` changes only on an explicit owner Goal choice); "Task Markdown" still tells the planner what draft to prepare. No rule change.

No materiality actions are required.

## Owner Decisions

No owner acceptance or exception is recorded on O-039 or its Tasks, so no `carried_forward` entry is written. T-113 carries an owner Task-check waiver; the repair did not touch T-113's text, so it is unaffected.

## Owner Validation Still Needed

- T-113 declares `owner_validation.required: true` with no recorded acceptance (`accepted_check: ""`). The owner may accept against this Check or decide it is not needed; until then the Objective cannot close on T-113's acceptance.
- I-144 is closed `verified` by this Check.

## Observations (non-blocking)

- O1 (carried from C-971): "Reporting to the Owner" still appears above and inside the managed block in this repository's `AGENTS.md`.
- O2: I-144 has no `repair_attempted` history entry from the executor who made the repair; the repair is proven here regardless.
- O3: C-971's `checked_at` (07:05:00Z) is later than this re-check's real time; the prior record's time appears to be estimated. Not edited (records are immutable).
- O4 (carried from C-971): the "run `savepoint resume` after creating a record" rule is repeated in four places; candidate for O-040.

## Code Style Review

- [x] STYLE-01 **One job per file**
- [x] STYLE-02 **One job per function**
- [x] STYLE-03 **Test branches**
- [x] STYLE-04 **Types document intent**
- [x] STYLE-05 **Build only what is needed**
- [x] STYLE-06 **Handle errors at boundaries**
- [ ] STYLE-07 **One source of truth** — "Reporting to the Owner" still duplicated in `AGENTS.md`; see O1.
- [x] STYLE-08 **Comments explain why**
- [x] STYLE-09 **Content lives in data**
- [x] STYLE-10 **Small diffs**
