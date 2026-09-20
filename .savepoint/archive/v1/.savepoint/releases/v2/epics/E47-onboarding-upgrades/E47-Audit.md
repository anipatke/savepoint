---
type: audit-findings
audited: 2026-09-19
---

# Audit Findings: E47 Start and refresh coherent V2 projects

## Main Findings

### Verdict

REMEDIATION APPLIED, OWNER REVIEW PENDING. Findings 1–5 below are closed by repository changes and regression tests. Finding 6 remains open because T007 is still `in_progress` and the owner must review it before the router can advance to `audit-pending`. This audit does not mark the task done or advance the router.

### Prior Finding Closure Map

| Prior finding | Re-audit result | Evidence |
|---|---|---|
| Adoption guidance contradicts managed `AGENTS.md` writes | Closed | V2 guidance now states the managed-block exception and byte-preservation boundary; `v2_scaffold_test.go` asserts both phrases. |
| Retired-but-already-absent skill leaves a stale manifest hash | Closed | Retirement forgets every successfully processed retired path, including a confirmed-absent file; `TestRetireV1Skills_missingFileForgetsStaleManifestEntry` covers the stale hash. |
| Retirement reports do not name recovery locations | Closed | Retirement entries retain the exact archive path, including numbered conflicts; the lifecycle matrix asserts formatted output for archive and `.new` recovery paths. |
| The claimed V1 compatibility baseline is self-referential | Closed | V1 migration uses the frozen pre-E47 README bytes, and dispatch is checked against `internal/init/testdata/legacy/migrations-readme-pre-e47.md`, not a second current code path. |
| Existing migration documentation stays stale after V2 retirement | Closed | V2 upgrades replace only the exact stock pre-E47 README and preserve edited documentation; focused stock/edit regression cases cover the sequence. |
| E47 handoff remains in progress and uncommitted | Still open | T007 remains `in_progress`; the owner must review the remediation, mark the task done, and advance the router to `audit-pending`. |

### What Needs Attention

#### 1. Adoption guidance — applied

The V2 adoption guide still says adoption modifies no user-authored file and writes only below `.savepoint/`. In reality, init deliberately adds or refreshes the Savepoint-managed block in an existing root `AGENTS.md`. That managed-region behavior preserves surrounding bytes and is well tested, but the current promise is false and violates TPL-02.

Applied: the guide now states the managed-guide exception explicitly while retaining the guarantee for content outside the markers and for all other authored files.

Evidence: `templates/project-v2/AGENTS.md:70`; `internal/init/scaffold.go:43-49`; `internal/init/agents.go:43-58`.

#### 2. A missing retired skill — applied

When a retired V1 file is already absent but its manifest hash remains, retirement returns no entry and therefore does not call `manifest.Forget`. The saved manifest can continue claiming ownership of a nonexistent path, contrary to T006's explicit reconciliation rule.

Applied: retirement forgets a path after successful removal or confirmed absence, while read, archive, and removal failures still retain the entry; the missing-file regression is covered.

Evidence: `internal/init/retire_v1_skills.go:55-65,77-80`.

#### 3. Retirement reports — applied

T007 requires a combined retirement-and-conflict run to name both affected paths and both recovery locations. Conflict entries carry a `.new` note, but retirement entries carry no note or archive path, so formatted output names only the deleted source. The test comment repeats the requirement but checks archive existence rather than `report.Format()`, allowing the acceptance failure to pass.

Applied: retirement retains the exact selected archive path, includes it in every retirement entry, and tests formatted combined output plus numbered archive conflicts.

Evidence: `internal/init/retire_v1_skills.go:85-92`; `internal/init/upgrade.go:66-71,140-145`; `internal/init/lifecycle_matrix_test.go:561-607`.

#### 4. V1 compatibility baseline — applied

E47 promises that a V1 upgrade produces the same actions and bytes as before, apart from one informational line. The baseline test instead compares schema dispatch with the current direct upgrade function; both paths share all E47 changes. In particular, E47 expanded the migration README, so an older V1 project retiring the generic audit skill now receives different bytes and the test cannot detect it.

Applied: V1 keeps the pre-E47 migration README bytes, and the dispatch regression compares the result with the frozen fixture rather than another current implementation path.

Evidence: `internal/init/migrate_audit_skill.go:23-47`; `internal/init/upgrade_schema_test.go:391-430`; T005 acceptance criterion at `T005-refresh-assets-that-match-the-projects-own-version.md:44`.

#### 5. Existing migration documentation — applied

A project that previously retired the generic audit skill already has `.savepoint/migrations/README.md` with the older explanation. After it migrates to V2, `writeArchive` sees that README exists and never refreshes it, even while archiving the nine V1 skills. The resulting V2 project fails T006's requirement that the README explain those retired assets.

Applied: the V2 path upgrades only the exact stock legacy README and preserves local edits; the pre-existing-stock-README sequence is tested.

Evidence: `internal/init/migrate_audit_skill.go:125-132`; T006 acceptance criterion at `T006-retire-the-old-instructions-without-losing-edits.md:46`.

#### 6. Epic handoff is still incomplete

The explicit re-audit request allows review before the normal phase transition, but it does not grant authority to complete T007. The task remains `in_progress`, the router remains `task-building`, and the new lifecycle matrix is untracked.

Next: after the repository findings are fixed and rechecked, the user must review T007, mark it done, and move the router to `audit-pending`.

Evidence: `T007-prove-the-whole-onboarding-loop-protects-user-files.md:4-5`; `.savepoint/router.md:16-19`; `git status --short`.

### Materiality Summary

| Finding | Likelihood | Impact | Materiality | Recommendation |
|---|---|---|---|---|
| Adoption guidance contradicts managed `AGENTS.md` writes | High | Medium | High | Fix the shipped guidance to state the managed-block exception explicitly. |
| Retired-but-already-absent skill leaves a stale manifest hash | Medium | Medium | Medium | Fix manifest cleanup and add the missing-file regression case. |
| Retirement report omits each archive recovery location | High | Medium | Medium | Name the exact archive in retirement entries and assert formatted combined output. |
| V1 compatibility baseline cannot detect E47 byte changes | Medium | Medium | Medium | Restore V1 README compatibility and replace the self-comparison with a frozen pre-E47 fixture. |
| Pre-existing migration README remains stale on V2 retirement | Medium | Medium | Medium | Safely upgrade the stock legacy README on the V2 path and test the migration sequence. |
| E47 handoff remains in progress and uncommitted | High | High | High | Complete repository fixes, then obtain user review and advance to `audit-pending`. |

### What Is Proven / Not Proven

Proven after remediation:

- The V2 scaffold, prompt, router, per-tree parity, init wiring, schema dispatch, conflict/backup behavior, retirement archives, optional-file behavior, repeat upgrades, dry runs, malformed-version refusals, and user-byte preservation cases exercised by the current matrix pass.
- Focused lifecycle and retirement tests, uncached `go test ./...`, `go vet ./...`, `git diff --check`, and `make build && make test` all pass.

Still pending:

- T007 has not received the owner-controlled completion handoff.

### Audit Evidence

- Scope lock: unchanged from the initial audit—E47 criteria and guardrails over init, `upgrade-assets`, schema selection, V1/V2 templates, prompt/guidance, manifest ownership, retirement/archive behavior, recovery reporting, repeat/dry-run behavior, and their focused tests.
- Re-audit admission: the three additional findings map respectively to the already-recorded combined recovery-report cell (T007), V1 exact-byte baseline cell (T005), and existing migration-documentation/retirement cell (T006); no new axis or dependency layer was added.
- Coverage and workflow result: all original public paths and lifecycle cases were rerun or traced. Writes still follow validation/version read, retirement archive, deletion, asset refresh, then manifest save. The five repository findings are closed by the applied changes and named regression cases.
- File reality and drift: scoped files exist or are explicitly recorded as moved; the T007 matrix and this audit are present for owner review. No unrecorded architecture drift was found. `.savepoint/Health-Check.md` and `.savepoint/audit/` remain absent and were correctly skipped.
- Gates: focused E47 tests pass; `go test ./... -count=1`, `go vet ./...`, `git diff --check`, and `make build && make test` pass. `gofmt -l` still reports only the pre-existing `internal/init/manifest_test.go` formatting debt.

### Guardrails Verification

- Rule IDs checked: FS-01, FS-02, FS-03, FS-04, FS-06, DATA-02, TPL-01, TPL-02, TPL-03, TPL-04, CFG-01, ARCH-01, ARCH-04, TEST-01, TEST-02, TEST-03, TEST-05, TEST-06, and TEST-08.
- Health check mode: Full. No project `Health-Check.md` exists, so the repository's Full evidence consists of the mapped guardrails, acceptance matrix, focused probes, and full build/test gates.
- Blockers/required rules: TPL-02, TEST-01/05/06, and the related acceptance evidence are remediated by the changes and tests above. The owner handoff remains open; no waiver is recorded.
- File safety: no destructive audit action was taken and no product/planning file was edited. The only audit write is this handoff file.

### Non-Blocking Observations

`gofmt -l` continues to report the pre-existing `internal/init/manifest_test.go` formatting debt. It is advisory here and does not account for any finding.

## Code Style Review

- [x] STYLE-01 **One job per file** — split files when responsibilities mix.
- [x] STYLE-02 **One job per function** — small, named, testable units.
- [x] STYLE-03 **Test branches** — cover meaningful conditionals and edge cases.
- [x] STYLE-04 **Types document intent** — prefer explicit types over comments.
- [x] STYLE-05 **Build only what is needed** — no speculative abstractions.
- [x] STYLE-06 **Handle errors at boundaries** — validate inputs, APIs, IO, and external data.
- [x] STYLE-07 **One source of truth** — no duplicated rules, constants, state, or config.
- [x] STYLE-08 **Comments explain why** — not what the code already says.
- [x] STYLE-09 **Content lives in data** — keep copy/config out of logic.
- [x] STYLE-10 **Small diffs** — minimal, reviewable, behaviour-preserving changes.

## Applied Remediation

- `templates/project-v2/AGENTS.md` now documents the managed `AGENTS.md` block exception and the byte-preservation boundary; `v2_scaffold_test.go` locks both phrases.
- `retire_v1_skills.go` clears stale manifest provenance for confirmed-absent retired files and records the exact selected archive path, including numbered conflicts.
- V1 and V2 migration README payloads are lifecycle-specific. V1 keeps the frozen pre-E47 bytes; V2 upgrades only an exact stock legacy README and preserves edited documentation.
- `upgrade_schema_test.go`, `retire_v1_skills_test.go`, and `lifecycle_matrix_test.go` cover the frozen compatibility fixture, missing provenance, stock/edit README behavior, numbered archive notes, and formatted combined recovery output.
- The repository remains in `task-building` with T007 `in_progress` until the owner reviews these changes and advances the handoff.
