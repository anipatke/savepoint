---
id: C-919
scope: {kind: objective, id: O-018}
result: CLEAR
checked_by: {role: checker, session: o018-recheck-20260924}
executed_session: o018-full-check-20260923
checked_at: '2026-09-24T10:32:54Z'
reviewed:
  base_commit: a1b55877d177060e83c9d867591e0ad8d2549630
  head_commit: 25e5b2dc56c2fe010a6d353254946e1979a59edc
  files:
    - 'internal/data/identity_v2.go#sha256=195febcbfe1738a6c8a8538246d54c0f489e24a999af21ddde110e51224f49b0'
    - 'internal/migrate/plan.go#sha256=239d24699e18d21bfa2b4c24a1808923811bc571de5d31c5d87ccb5522ced0bd'
    - 'internal/migrate/convert.go#sha256=9c85f77fe0984d4c4489c5ae1965f4750f4702e48e4aa4cd6e26da0b29860789'
    - 'scripts/hyphenate_record_ids.py#sha256=18894f8dacbe270b5f5903cb6ff6541b38c7dcca673d94ce3869e7cd3be593fc'
    - 'AGENTS.md#sha256=fec84b9cc073580463e7f6f264fba574e4159165c3a2000aa3f2f80d519b2a98'
    - '.savepoint/Design.md#sha256=039bba3deace91116bd0dba5a0237ce3bc92cdf62a4e82a7638c2b1c8bb2fd52'
    - '.savepoint/Guardrails.md#sha256=b2e548b5e8e460c6e31c2e044a726e576e3e7c2a6eca7876bc9fc2fcb898f93a'
    - 'README.md#sha256=89b833b1f0e10deb37c6e7600df57be58dd500d65f62b7e50fa915d998a169ae'
    - 'Makefile#sha256=6293168a7fde79200b748521238e6aa0a0274057dee0594f3471997dff11fe37'
    - 'scripts/git-hooks/pre-commit#sha256=312d6f2ac8822538b93686dd3d2b31219460b623fef9140a95ca414bad315fc6'
    - '.savepoint/objectives/O-018-hyphenate-numbered-record-identities/Objective.md#sha256=d416b90a39fb219201db2d06c710602b2a2c0693fedd1d736eb0070cb383a478'
    - '.savepoint/objectives/O-018-hyphenate-numbered-record-identities/tasks/T-018-map-the-identity-cutover.md#sha256=8a0121d56fdb40c83211497e80d157b88530bd6cfaef04fb534772c37694b415'
    - '.savepoint/objectives/O-018-hyphenate-numbered-record-identities/tasks/T-020-rename-this-projects-records.md#sha256=a9391250135d3f6c3af7f74ea77268ebbfc40705303e3fd9964cf4ab6a65dfc5'
    - '.savepoint/issues/I-022-hyphenate-numbered-record-identities.md#sha256=1bc82f4279349492038991ec80d6badc7d9a5dd7f6e4707a21aa95d9fe0cb574'
  dependencies: []
issues: []
supersedes: C-910
---

# C-919: O-018 Full Objective Recheck

## Verdict

`CLEAR`. This run rechecks C-910 after the I-033 and I-034 repairs. Both
repairs hold, every original matrix cell passes or is recorded as not
applicable, and a fresh `make build && make test-full` passes at the reviewed
head.

Closure map of the prior Issues:

| Issue | C-910 finding | Now |
|---|---|---|
| I-033 | Bare `X###` placeholders in `AGENTS.md` and `Design.md`; bare historical IDs in the Makefile, pre-commit hook, and V2 board comments | Closed: the I-033 proof search returns no hits |
| I-034 | Three bare-versus-hyphenated examples rewritten into self-comparisons; T-020 "No false positives" claim inaccurate | Closed: all three sentences reworded; T-020 carries a dated correction note |

I-032, I-033, and I-034 were already resolved as owner-`accepted` on
2026-09-23. This run leaves those resolutions unchanged. It adds the
technical proof that their repairs landed.

## Independence

This session built none of O-018 and made neither repair. The repairs were
made by session `o018-full-check-20260923`, which also wrote C-910. That is
why C-910's own author could not recheck them.

Reviewed state: HEAD `25e5b2d` plus an uncommitted router selection change
(`objective: O-019`, `task: T-010`), which is not O-018 code, tests, or
fixtures.

## Frozen Scope Lock

C-910's scope lock is reused unchanged: O-018 Success Conditions and
Boundaries; T-019 and T-020 Done When; the shared identity rule and every
decoder and router path that calls it; the V1→V2 allocator; the rename
script; active and scaffold guidance; this repository's live records. No axis
was added.

## Admission Ledger And Matrix

| # | C-910 cell | Recheck evidence | Result |
|---|---|---|---|
| 1 | One shared rule replaces five regexes | `identity_v2.go`: `^[ROTCI]-[0-9]{3,}$` plus kind byte; the only other `regexp` in `internal/data` are the audit-finding and completion-hash patterns | Proven |
| 2 | Bare IDs refused in each field family | Overlay probe on a copy of the live `.savepoint`: bare Task id, Objective `depends_on`, Objective `release`, Check `supersedes`, and Check `scope.id` each refused by `LoadV2Index` with a named diagnostic. `RouterSelectionV2.validate` refuses bare release/objective/task/issue; `ReadStateV2` refuses a bare router objective | Proven |
| 3 | Malformed shapes | Probe: lowercase, `--`, two digits, `_`, en dash, Arabic-Indic digits, trailing `\n`, leading space, trailing space, empty, `T-`, `T-01a`, wrong kind all refused; `T-000`, `T-1000`, `C-1000` accepted; `compareV2CheckIDs` orders `C-999` < `C-1000` | Proven |
| 4 | Mixed old/new graph refused | Row 2 overlays each put one bare reference into an otherwise hyphenated live graph | Proven |
| 5 | `CreateCheckV2`/`CreateIssueV2` emit hyphenated IDs | Both functions were deleted as unreachable by O-021 (`89ed584`, cleared by C-915). No supported Check/Issue generator remains in `internal/data` | N/A (removed) |
| 6 | Check ordering and next-ID parse the `-` suffix | `normalizedV2CheckNumber` trims `C-` (`project.go:424`); ordering proven in row 3. Next-ID half removed with row 5 | Proven / N/A |
| 7 | V1→V2 allocator emits `X-###`; converter reads hyphenated owners | `plan.go:421` `%s-%03d`; `convert.go:368` `ownerObjectiveID` requires `O-` plus digits; `internal/migrate` golden and end-to-end tests pass in row 17 | Proven |
| 8 | Skill/template examples hyphenated; canonical and scaffold copies identical | No bare IDs or `X###` placeholders in `templates/project-v2`, `agent-skills`, `README.md`, `cmd`, `internal/doctor`, `internal/resume`; byte parity covered by the TPL-01 parity tests in row 17 | Proven |
| 9 | Live strict load and router | `LoadV2Index(.savepoint)`: 6 R / 13 O / 34 T / 14 C / 48 I, every ID matches its kind; `ObjectiveTasks[O-018] = [T-018 T-019 T-020]`; live router decodes | Proven |
| 10 | Script renames only identity-named entries; archive and V1 `releases/v*` untouched | Replay on a clone at `a1b5587`: the renamed path set equals the delivered tree at `bb03eaf` with no rename differences. Every content difference is a later hand edit (I-033/I-034 repairs, T-019/T-020 evidence) or a record created after `a1b5587` (C-910, I-032..I-034), or an owner edit C-910 already listed | Proven |
| 11 | Repeat run changes nothing (FS-04) | Second replay run and a run on a clone of HEAD: "0 changed by content", clean `git status` on HEAD clone | Proven |
| 12 | Dry run writes nothing (FS-03) | `--dry-run` on the `a1b5587` clone: `git status` empty | Proven |
| 13 | Check records change only in identity references | Row 10 replay: no pre-existing Check differs from the delivered tree | Proven |
| 14 | Active prose has no unhyphenated IDs (I-033) | `grep -nE '\b[ROTCI](###\|[0-9]{3,})\b' AGENTS.md .savepoint/Design.md .savepoint/Guardrails.md Makefile scripts/git-hooks/pre-commit README.md`: no hits; `internal/board/v2` non-test sources: no hits | Proven |
| 15 | False positives fixed by hand (I-034) | O-018 `Objective.md` Why, I-022 observed note, and T-018 line 68 now describe the old form in words; T-020 has "Correction after C-910 (2026-09-23)" | Proven |
| 16 | `git diff --check` | No output | Proven |
| 17 | Fresh `make build && make test-full` (TEST-08) | 2026-09-24T10:30:53Z–10:31:00Z (warm build cache; tests run with `-count=1`), go1.26 linux/amd64, HEAD `25e5b2d`; full suite plus Linux/Darwin/Windows builds; EXIT=0 | Proven |
| 18 | Board TUI and CLI doctor/resume | `savepoint resume` runs on this repository and resolves `O-019`/`T-010` (AGENTS.md now permits the read-only resume). Doctor loads the project; its findings are unrelated pre-existing clearance diagnostics | Proven |
| 19 | T-018 inventory | Unchanged apart from the I-034 rewording | Proven (historical) |

The probe harness was a temporary `internal/data` test, removed after the run.

## Criteria Classification

- Shared identity rule refusing unhyphenated IDs: Proven (rows 1–4).
- Generators and V1→V2 allocator emit hyphenated IDs: Proven for the
  allocator (row 7); the V2 Check/Issue generators no longer exist (row 5).
- Skill, template, and fixture examples hyphenated; copies identical: Proven
  (row 8).
- Records, directories, filenames, cross-references, router, and active prose
  renamed by the in-repo script: Proven (rows 9–15).
- Strict load, board, doctor, and resume succeed; targeted search clean;
  `make build` and `make test-full` pass: Proven (rows 9, 14, 17, 18).

## Guardrails

FS-03, FS-04 (rows 11–12), DATA-01..04 (rows 2–4, 9), TPL-01..02 (rows 8,
14), TEST-01..05 and TEST-08 (row 17): satisfied.

## Materiality

No Issues. No materiality actions are required.

## Observations (non-blocking)

- Bare IDs remain in Issue bodies that quote the defect they record
  (I-032, I-033, I-034) and in the V1 source manifest
  `.savepoint/migrations/v1-to-v2.yml`. Both are historical evidence of the
  old form, not active identity references.
- C-910's observations still stand: Check files are written as
  `checks/C-###-slug.md` by hand while no generator remains to settle the
  naming, and `scripts/hyphenate_record_ids.py` has no collision guard.

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
- [ ] STYLE-10 **Small diffs** — O-018 landed as one large commit (`bb03eaf`) because T-019 and T-020 could not load apart; the Objective recorded that trade-off.

## Owner Validation Still Needed

T-020 declares `owner_validation.required: true` with an empty
`accepted_check`. Its owner acceptance must name this Check, C-919. The
checker does not record it on the owner's behalf.
