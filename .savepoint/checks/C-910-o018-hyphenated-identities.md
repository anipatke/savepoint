---
id: C-910
scope: {kind: objective, id: O-018}
result: NEEDS WORK
checked_by: {role: checker, session: o018-full-check-20260923}
executed_session: o018-t019-t020-20260923
checked_at: '2026-09-23T08:40:00Z'
reviewed:
  base_commit: a1b55877d177060e83c9d867591e0ad8d2549630
  head_commit: a1b55877d177060e83c9d867591e0ad8d2549630
  files:
    - internal/data/identity_v2.go
    - internal/data/identity_v2_test.go
    - internal/data/release_v2.go
    - internal/data/objective_v2.go
    - internal/data/task_v2.go
    - internal/data/check_v2.go
    - internal/data/issue_v2.go
    - internal/data/evidence_v2.go
    - internal/data/router_v2.go
    - internal/data/project.go
    - internal/data/write.go
    - internal/migrate/plan.go
    - internal/migrate/convert.go
    - internal/doctor/repairs.go
    - scripts/hyphenate_record_ids.py
    - agent-skills/savepoint-check/SKILL.md
    - agent-skills/savepoint-design/SKILL.md
    - agent-skills/references/issue-capture.md
    - templates/project-v2/AGENTS.md
    - AGENTS.md
    - README.md
    - .savepoint/Design.md
    - .savepoint/router.md
    - .savepoint/objectives/O-018-hyphenate-numbered-record-identities/Objective.md
    - .savepoint/objectives/O-018-hyphenate-numbered-record-identities/tasks/T-018-map-the-identity-cutover.md
    - .savepoint/objectives/O-018-hyphenate-numbered-record-identities/tasks/T-019-require-hyphenated-ids.md
    - .savepoint/objectives/O-018-hyphenate-numbered-record-identities/tasks/T-020-rename-this-projects-records.md
  dependencies: []
issues: [I-033, I-034]
supersedes: null
---

# C-910: O-018 Full Objective Check

## Verdict

`NEEDS WORK`. The code side of O-018 is sound:

- One shared rule refuses bare and malformed IDs on every decoder and router
  path.
- The Check and Issue writers mint `C-###`/`I-###`, and the V1→V2 allocator
  emits hyphenated IDs.
- The rename script reproduces the staged record graph exactly and makes no
  changes on a second run.
- A fresh `make test-full` passes.

Two documentation defects remain:

- I-033: `Design.md` and root `AGENTS.md`, both named by O-018, still document
  the identity and path format as bare `R###`/`O###`/`T###`/`C###`/`I###`.
- I-034: the rename script turned three bare-versus-hyphenated examples into
  sentences that compare an ID with itself.

O-018 is not ready for owner acceptance.

Reviewed state: this is the uncommitted staged index on top of `a1b5587`
(tree `95eb7c4`) plus T-020's unstaged frontmatter change. T-019 and T-020
have not been committed yet, per the same-commit boundary.

## Evidence Mode

Full Objective evidence. This run reviewed all three owned Tasks, T-018
through T-020, each carrying an owner Task-check waiver. It also covered
cross-Task integration (T-019's grammar with T-020's rename) and
reconciliation against Design, AGENTS.md, and Guardrails. This session is
fresh and did not build any of the work.

## Frozen Scope Lock

1. Acceptance and policy:
   - O-018's five Success Conditions and Boundaries.
   - T-019 and T-020 Done When. T-018 is reviewed only as the superseded
     inventory; its migration contract was replaced by the 2026-09-23 replan.
   - Guardrails FS-01, FS-03, FS-04, DATA-01..04, TPL-01..02, TEST-01..05,
     and TEST-08.
2. Changed behavior and public entry points:
   - `matchesV2Identity` and every decoder and validator that calls it
     (Release, Objective, Task, Check, Issue, evidence, router read,
     router-selection write).
   - `nextV2CheckID`/`nextV2IssueID` and `compareV2CheckIDs`.
   - `idAllocator` and `ownerObjectiveID` in `internal/migrate`.
   - `doctor` repair wording.
   - `scripts/hyphenate_record_ids.py`.
   - Active and scaffold guidance, and this repository's live records.
3. Relied-on orchestration: `LoadProject` → `LoadV2Index` → discovery
   filename-prefix checks → decoders → cross-link resolution. Also
   `CreateCheckV2`/`CreateIssueV2` writing create-only files, and the
   script's `git mv` plus in-place content rewrite.
4. Matrix axes:
   - ID shape: valid, bare, lowercase, double hyphen, two digits, 4+ digits,
     underscore, en dash, non-ASCII digits, leading/trailing whitespace,
     wrong kind.
   - Field family: id, owner, dependency, release, supersedes, router,
     evidence.
   - Mixed old/new graph.
   - Generators: live index and fresh allocation.
   - Script: first run, repeat run, dry run, non-identity siblings, archive.
   - Guidance: canonical versus scaffold parity, active prose.
5. Materiality boundary: a finding must violate an O-018/T-019/T-020
   criterion or a named Guardrail through a supported load, generator,
   script, or guidance path. Script collision and partial-state handling are
   out of scope by the Objective's "simple rename" boundary.

## Coverage Matrix

| # | Cell | Evidence | Result |
|---|---|---|---|
| 1 | Shared rule `^[ROTCI]-[0-9]{3,}$` plus kind byte replaces five regexes | `identity_v2.go`; staged diff removes `release/objective/task/check/issueIDPattern*`; no other ID regex remains in `internal/data` | Proven |
| 2 | Bare ID refused: Task id, depends_on, objective owner, Objective release, Check supersedes, router objective | Overlay probe on copies of the live `.savepoint`: each load refused with a named diagnostic (`ErrV2InvalidID` / dependency / ownership / release reference) | Proven |
| 3 | Malformed shapes: lowercase, `--`, 2 digits, `_`, en dash, Arabic-Indic digits, trailing `\n`, leading space | Overlay probe: all refused; `C-1000` accepted | Proven |
| 4 | Mixed old/new graph refused (one bare reference in an otherwise hyphenated graph) | Probe row "bare dep ref (mixed)" refused | Proven |
| 5 | Generators emit hyphenated IDs over the live graph | Probe: `CreateCheckV2` → `C-910` (`checks/C-910.md`), `CreateIssueV2` → `I-033` (`issues/I-033-probe-issue.md`); both reload | Proven |
| 6 | Check ordering and next-ID parse the `-` suffix | `project.go:472`, `write.go:1107,1331` trim `C-`/`I-`; generator probe allocated after C-909/I-032 | Proven |
| 7 | V1→V2 allocator emits `X-###`; converter finds the owning Objective from a hyphenated directory | `plan.go:422`; `convert.go` `ownerObjectiveID`; golden fixtures regenerated; `internal/migrate` passes in the full gate | Proven |
| 8 | Skill/template examples hyphenated; canonical and scaffold copies byte-identical (TPL-01) | `cmp` of all 7 V2 skill/reference pairs: identical; no bare IDs or `X###` placeholders in `templates/project-v2`, active V2 skills, README, `cmd`, doctor, or resume | Proven |
| 9 | Live repository strict load and router selection | Probe `LoadProject(.savepoint)`: 6 R / 10 O / 20 T / 5 C / 32 I; `ObjectiveTasks[O-018] = [T-018 T-019 T-020]`; router `R-006/O-018/T-020` | Proven |
| 10 | Script renames only identity-named entries and leaves `archive/` and V1 `releases/v*` untouched | Replay on `git archive HEAD` copy: output matches the staged tree except for unrelated owner edits (see Observations); `archive/` and `releases/v*` unchanged | Proven |
| 11 | Script repeat run changes nothing (FS-04) | Second replay run: "0 changed by content", no renames | Proven |
| 12 | Script dry run writes nothing (FS-03) | Code path: `git_mv`/`rewrite_references` return before any write when `dry_run` is set | Proven (inspection) |
| 13 | Check records change only in identity references | Word-diff of all 5 renamed Checks with IDs normalised: removed and added token sets are identical | Proven |
| 14 | Active prose (`AGENTS.md`, `Design.md`, `Guardrails.md`) has no unhyphenated IDs | `Guardrails.md` clean; `AGENTS.md:64,88` and 12 `Design.md` placeholders still bare; bare historical refs remain in the Makefile, the pre-commit hook, and V2 board comments | **Issue I-033** |
| 15 | Script rewrites identity references only, and false positives are fixed by hand | O-018 `Objective.md:19`, I-022 history note, T-018 evidence line 68 are now self-comparisons | **Issue I-034** |
| 16 | `git diff --check` | Staged and unstaged: no output | Proven |
| 17 | Fresh `make build && make test-full` (TEST-08) | 2026-09-23T08:23:15Z–08:25:26Z, go1.26.2 linux/amd64, over tree `95eb7c4` plus the unstaged T-020 frontmatter change; full suite plus Linux/Darwin/Windows builds; EXIT=0 | Proven |
| 18 | Board TUI and CLI doctor/resume output | Not run: AGENTS.md forbids agents running `savepoint` commands. Load path covered by row 9; the renderers only format IDs taken from the index, and their fixtures pass in row 17 | N/A (policy) |
| 19 | T-018 inventory deliverable | Present and traceable. Its migration contract was superseded by the replan; its failed `test-fast` is superseded by row 17 | Proven (historical) |

## Adversarial Pass

- Bypass via another entry point: the router write path
  (`RouterSelectionV2` validation in `write.go`) and the router read path
  both use the shared rule. No decoder keeps a private regex.
- Discovery prefix check (`discover.go:368,435,551`) is unchanged and still
  requires the basename to start with the declared ID. It interacts correctly
  with the extra hyphen.
- The script, given a sibling that is already hyphenated, would `git mv` the
  bare directory into it. This falls outside the "simple rename" boundary;
  recorded as an observation.

## Issues

- **I-033** — `Design.md` and `AGENTS.md` still document bare-format IDs and
  paths. Violates O-018 Success Conditions 4–5, T-020 Done When, and TPL-02.
- **I-034** — the script turned three before/after examples into
  self-comparisons, and T-020's "No false positives" evidence is inaccurate.
  Violates T-020 Done When ("identity references").

| Issue | Likelihood | Impact | Materiality | Recommendation |
|---|---|---|---|---|
| I-033 | Medium: agents read AGENTS.md every session and Design.md during planning | Medium: a hand-minted `I###-slug.md` or `C###-slug.md` fails strict load; the failure is loud but blocks the board | Medium | Fix now; the edit is a small hyphen insertion |
| I-034 | High: already present in three records | Low: meaning of historical prose lost; no load or behavior impact | Low | Fix together with I-033 |

## Remediation Routing

Hand remediation to the executor or planner as new or newly selected work
linked to O-018, I-033, and I-034. T-018, T-019, and T-020 are `done` and
must not be retreated. The later re-check reuses this frozen scope lock.

## Owner Validation Still Needed

T-020 declares `owner_validation.required: true` with an empty
`accepted_check`. Owner acceptance must name the current CLEAR Check once one
exists; this `NEEDS WORK` run cannot satisfy it.

## Related Issue

I-032 (README placeholders): a targeted search of `README.md` finds no bare
IDs or `X###` placeholders, so the repair holds. It stays open: a
`verified` closure needs a CLEAR Check, and this run is `NEEDS WORK`. The
re-check can close it.

## Observations (non-blocking)

- `CreateCheckV2` writes `checks/C-###.md` without a slug, while Design and
  AGENTS.md document `C-###-slug.md`. T-018 asked T-019/T-020 to settle this
  naming convention, but the replanned Objective dropped it and nobody
  decided. It predates O-018; worth a follow-up decision.
- `scripts/hyphenate_record_ids.py` has no collision guard. When a
  hyphenated sibling already exists, `git mv` nests the bare directory inside
  it, and the plain-rename fallback can overwrite a file. Content rewriting
  hits any word-bounded `[ROTCI]\d{3,}` in prose, which is how I-034
  happened. The owner scoped this as a simple rename; other projects should
  run `--dry-run` first.
- The pending staged change also carries edits outside O-018:
  - O-014's mandatory-Release decision
  - O-016's closure note and `status: done`
  - `release: R-006` added to O-016, O-017, and O-019
  - I-031 title and history
  - the new `.gitattributes`

  None of these were reviewed as O-018 work. Consider committing them
  separately.
- The router still says `state: task` and names the finished T-020 audit
  step. Updating it is outside the checker's write boundary.
