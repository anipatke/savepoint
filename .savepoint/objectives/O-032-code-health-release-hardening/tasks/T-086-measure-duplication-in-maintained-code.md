---
id: T-086
title: Measure duplication in maintained code
objective: O-032
status: in_progress
stage: audit
depends_on: [{task: T-085, requires: clear}]
complexity_tier: spike
complexity_reason: A global duplicate percentage cannot establish the maintained-code baseline or safe refactoring scope.
owner_validation: {required: true}
planned_by: {role: planner, session: planning-o032-20261002-owner-confirmed}
---

# Measure duplication in maintained code

## Outcome

A reproducible maintained-code duplication baseline and concrete clone-repair proposals separate real maintenance costs from preserved history and required template mirrors.

## User Check

Review which files are measured and which real duplicate blocks should be changed.

## Done When

- Record whole-repository baseline evidence and proposed production-Go/test-Go scopes, exact provider targets/exclusions/version and denominator; do not calculate a scoped percentage by filtering the capped evidence list.
- Use project-owned duplication tooling directly with temporary report paths when needed; leave confirmed config, canonical/template copies and archived history unchanged.
- Prove whether this provider version has per-file counts or needs explicitly scoped execution. Identify substantive production clone pairs and reusable test setup separately from required mirrors/fixtures.
- Deliver the named duplication-scope-and-repair decision, including before/after comparison rules, trend restart implications, proposed config change for owner review and bounded follow-up repairs. Keep 3/5 percent thresholds; target <=5 percent after justified repairs, with 3 percent an aim rather than an invented result.

## Context Files

`.savepoint/objectives/O-032-code-health-release-hardening/tasks/T-085-find-why-the-vulnerability-scan-fails.md` (dependency evidence); `.savepoint/health/config.json`; `.savepoint/health/report.md`; `internal/codehealth/config.go`; `internal/codehealth/reader_jscpd.go`; `internal/codehealth/reader_jscpd_test.go`; `internal/codehealth/reader_paths.go`; `internal/codehealth/testdata/readers/duplication/jscpd5.json`; `internal/codehealth/testdata/readers/duplication/versioned.json`.

## Design References

O-032 Confirmed Release-Hardening Design and Owner-Requested Health Analysis; Design sections 1, 6, 7, 9, 12 and 13 as applicable.

## Guardrails

FS-01, FS-04, FS-05, FS-06, ARCH-01, ARCH-02, ARCH-03, CFG-01, CFG-02, TEST-01, TEST-02, TEST-03, TEST-04, TEST-05, TEST-08, TEST-09, STYLE-01, STYLE-02, STYLE-07.

## Implementation Plan

1. Inspect saved top evidence and reader scope/denominator contract.
2. Run bounded maintained-code baseline scenarios with temporary output and document actual numerator/denominator.
3. Rank substantive clone pairs with exact locations and preservation risks.
4. Propose scoped measurement and narrow fixes; return to design when implementation details are not settled.

## Boundaries

Diagnosis and proposals only; no health config edit, automatic exclusion of tests, archive/template modification or threshold changes.

## Technical Verification

`make build && make test-fast` for ordinary handoff; focused `make test-focused TEST=... PKGS=...` only during iteration. Record named per-criterion evidence and limitations. The later independent Check applies `agent-skills/references/check-method.md`; the mandatory Full Objective Check runs its fresh full gate and official Code Health collection. This executor never writes a Check or runs savepoint health setup/report/check. Research completion requires the named decision deliverable; it does not claim any recommended repair has landed.

## Technical Evidence

Run 2026-10-02 ~03:35Z, jscpd 5.4.0 (the configured provider), go1.26.2, linux/WSL2. Reports and copied scope trees stayed in the session scratchpad; nothing committed. `savepoint health check/setup/report` not run; config, templates, archive and thresholds untouched.

### Duplication-scope-and-repair decision

**Baseline (whole repository, configured args, temp output):** 10.10% = 20,459 duplicated of 202,559 lines, 1,042 clones, 1,844 sources. The saved report said 9.8% from an earlier tree; the shape is unchanged. Per-format: Go 5.23% (4,909/93,933), markdown 10.03%, yaml 12.17%, json 67.5%.

**Provider contract:** jscpd 5.4.0 writes only per-format aggregates and an integer `sources` count; there are no per-file line counts. So `jscpdScopedCounts` finds nothing and an instance scope cannot rebuild a percentage from a whole-repo run. A scoped percentage needs a separate jscpd run over only the scoped files (denominator = that run's own `lines`). I did not filter the capped evidence list to compute any percentage.

**Where the 10.1% comes from (clone lines by first file; lines are counted per clone, so the sum is 20,774, not 20,459):** `.savepoint` records 8,403; archived V1 5,659; Go tests 3,911; template/agent-skills mirrors 1,637; testdata 671; other 280; production Go 213. About 90% of the headline is history, project records, or required mirrors.

**Scoped runs (exact file lists copied to a temp tree, jscpd defaults otherwise):**

| Scope | Files | Lines | Duplicated | % | Clones |
|---|---|---|---|---|---|
| Production Go (`*.go` minus tests, archive, templates) | 133 | 35,536 | 213 | 0.60 | 22 |
| Test Go (`*_test.go`, same exclusions) | 145 | 56,237 | 4,000 | 7.11 | 428 |
| Maintained Go (both) | 278 | 91,773 | 4,213 | 4.59 | 450 |
| Maintained code and docs (Go, md, yml, json, sh, js, ts, py outside archive, templates, testdata, `.savepoint/`, dist) | 310 | 96,917 | 4,381 | 4.52 | 467 |

Both maintained scopes are already under the 5% watch line and above the 3% aim. Test Go is the only part over 5%.

**Production clone pairs (all small, max 19 lines):** `data/router.go:48<->77` (19), `data/gate_v2.go:279<->330` (14) and `gate_v2.go:350<->objective_gate_v2.go:44` (12), `migrate/convert_docs.go:182<->224` (14), `codehealth/reader_coverage.go:240<->371` (13), `data/parser.go:129<->260` (13), `migrate/convert.go:138<->221` and `<->convert_releases.go:45` (11), `migrate/convert_issues.go:154<->213` (11/10). None is a cost worth a Task on its own; `router.go` and `gate_v2.go` are the only ones with a plausible shared helper.

**Reusable test setup (the real cost):** `init/agent_skills_test.go` (1,955 clone lines, 105 clones, repeated 12-14 line blocks at :66/:242/:553/:815), `data/write_test.go` (996; 25-line blocks at :87<->:1120, :328<->:865), `init/upgrade_test.go` (943; 21 lines at :89<->:545), `data/next_test.go` (373), `data/project_test.go` (272). Largest single pair: `data/migration_history_test.go:194<->migration_source_test.go:99` (46 lines, the fixture inventory walk, a clear candidate for one helper).

**Required mirrors/fixtures (leave alone):** every `agent-skills/**` file against `templates/project-v2/agent-skills/**` (e.g. `check-method.md` 391 lines, `savepoint-design/SKILL.md` 285 lines); `internal/migrate/testdata/golden/*.yml`; archived V1 PRD/Audit files; Check pairs C-948/C-949.

**Proposed config change (owner review, NOT applied):** keep the whole-repo `duplication` instance as is for history, and replace it with, or add beside it, scoped instances: (a) production Go, `scope: ["**/*.go"]` plus exclusions of `**/*_test.go`, and (b) test Go; each with its own jscpd run over that target (jscpd `--pattern` or a target file list), since the provider cannot rescope a total. Exclude `.savepoint/archive/**`, `templates/**`, `**/testdata/**` and `.savepoint/**` records from the maintained run; do not exclude tests wholesale from the gate. Thresholds stay 3/5.

**Before/after comparison rules:** compare only runs with the same provider version, target list, exclusions and denominator. A scoped run's percentage is never comparable with the 10.1% whole-repo number. Changing scope restarts the trend for that instance (new comparison series); the old whole-repo series stays as history. Report numerator and denominator with each number.

**Bounded follow-up repairs (proposals, none landed):** (1) one fixture-inventory helper for the migration_history/source test pair (-46 lines of clone, first pair); (2) shared setup helpers in `data/write_test.go` and `init/upgrade_test.go`; (3) a shared skill-file loader in `init/agent_skills_test.go`; (4) optional small helpers in `data/router.go` and `gate_v2.go`. Target: maintained scope <=5% after repairs (it is 4.59% now, so the target is already met by measurement; the repairs aim toward 3%, which is an aim, not a promised result). Estimated reachable: test Go roughly 7.1% -> 4-5% if the first three repairs remove about half of the 4,000 clone lines; this is an estimate, not measured.

### Owner-directed config change and re-measurement (scope extension)

At 2026-10-02 the owner instructed, in-session, to exclude Savepoint docs from the duplication stat, then to add `**/testdata/**` as well and record both here. This overrides the Boundaries line "no health config edit" for this one change only; Done When is unchanged and the Boundaries text is left as planned. Not a Check or waiver.

- Changed `.savepoint/health/config.json`, `duplication` instance only: `--ignore` gains `.savepoint/**,templates/**,**/testdata/**`, and `exclusions` gains the same three globs so a reader of the config sees one scope. Thresholds (3/5), provider, report path and every other instance are unchanged. `go test ./internal/codehealth -run Config` passes.
- Re-measured with jscpd 5.4.0 and the new `--ignore`, output in the scratchpad (no `savepoint health check`): **4.51% = 4,381 of 97,040 lines, 467 clones, 312 sources** (was 10.10% = 20,459/202,559). Go 4.59% (4,213/91,773); markdown 3.13% (143/4,573); yaml 6.97% (14/201); text 20% (11/55, 55 lines only).
- Trend: this starts a new comparison series for the instance; the whole-repo 9.8-10.1% numbers are not comparable. The drop is a scope change, not a code improvement.
- Risks accepted: future copy-paste inside Check/Issue/Task records, template drift against `agent-skills/` (45 clone lines remain visible there) and fixture duplication under `testdata` are no longer measured. Undo is one config string.

### What remains (4.51%)

96% of the remaining duplicated lines (4,213 of 4,381) are Go test code; production Go is 213. 388 of the 467 clones are repeats inside a single file (table-style setup repeated in one test file), not copies across files. Sizes: 9 clones of 20+ lines, 144 of 10-19, 314 under 10. Concentration by file (clone lines): `init/agent_skills_test.go` 982, `data/write_test.go` 489, `init/upgrade_test.go` 456, `data/next_test.go` 187, `data/project_test.go` 137, `data/gate_v2_test.go` 108, `migrate/convert_issues_test.go` 84, `data/migration_history_test.go` 70; the three largest files hold about 44% of all remaining duplication. By package: `internal/init` 1,744, `internal/data` 1,393, `internal/migrate` 413, `board/v2` 222. Non-Go remainder is 168 lines: `agent-skills/references/check-method.md` against `examples.md` (20), `.github/workflows/ci.yml` against itself (14), and overlapping root-level `project-audit/*.md` reports (about 13 each). The 3% aim therefore needs test-setup helpers in the three largest test files; repairs to production code would not move the number.

### Owner-directed product change: default exclusions (scope extension)

At 2026-10-02 the owner instructed, in-session, to remediate within this Task so new projects get Savepoint docs and templates excluded by default, and to update docs. This extends the earlier owner-directed config change beyond "diagnosis and proposals only"; Done When is unchanged and the Boundaries text is left as planned. Not a Check or waiver.

- `internal/codehealth/discovery_catalogue.go`: `defaultExclusions` gains `.savepoint/**` (every tool); new `capabilityExclusions` adds `templates/**` and `**/testdata/**` for duplication only.
- `internal/codehealth/discovery.go`: `component.exclusions` now takes the capability and merges both lists; a directory pattern containing a wildcard (`**/testdata/**`) is always proposed, since it names no single entry to look up. Both callers pass the capability.
- `internal/codehealth/discovery_test.go`: default-exclusion and tool-argument tests cover `.savepoint/**` for complexity, duplication and OSV, and the duplication-only extras. Existing projects keep their confirmed config (setup does not rewrite confirmed entries).
- Docs: `.savepoint/Design.md` (`health setup` row) and `AGENTS.md` (`internal/codehealth` map row) state the defaults and the reason.
- Evidence: `go test ./internal/codehealth ./internal/healthcheck` passes; `make build && make test-fast` exit 0, 0 FAIL lines.
- Limitations: the setup preview still prints one merged "Default exclusions for new tools" list, so it does not say that `templates/**` and `**/testdata/**` apply to duplication only. `**/testdata/**` is proposed even when no testdata directory exists (harmless). Not run: `make test-full`, `savepoint health setup` against a fresh project. README has no health section, so nothing was added there. T-096 (user-facing explanation) may want to mention the defaults.

### Per-criterion outcomes

1. Baseline and scopes: whole-repo 10.10% with exact args, version 5.4.0 and denominator recorded; production, test, combined and non-record scopes measured by separate runs, none by filtering the capped list. Met.
2. Project-owned jscpd run directly with temp output; config, mirrors, archive untouched. Met.
3. Provider has no per-file counts in 5.4.0, so scoped execution is required; production clone pairs and test setup clones listed apart from mirrors/fixtures. Met.
4. Decision above: config proposal, comparison rules, trend restart, bounded repairs, thresholds kept. Met, with the 3% aim unquantified.

### Commands

- `jscpd --reporters json --output <scratchpad>/base --workers 1 --ignore 'dist/**,**/*.pb.go,**/*_generated.*,**/*.min.js,**/.git/**' .` (0.8s)
- jscpd defaults over four copied scope trees (`prod`, `test`, `gomaint`, `code`), python3 aggregation of the JSON reports
- `make build && make test-fast`: not rerun for this Task, no repo code changed by T-086 (the runner change in the working tree belongs to T-085).

### Files

Read: router.md, AGENTS.md, savepoint-task skill, T-085 and T-086, health/config.json, health/report.md, internal/codehealth/reader_jscpd.go. Extra reads: sections of `internal/data/migration_source_test.go` and `migration_history_test.go` to confirm the largest test clone. Not read though listed: `reader_paths.go`, `reader_jscpd_test.go`, the two testdata JSON files (their contract was clear from the reader). Changed: this Task (status/stage and evidence) `.savepoint/health/config.json`, `internal/codehealth/discovery{,_catalogue,_test}.go`, `.savepoint/Design.md` and `AGENTS.md` (owner-directed, see above).

### Limitations

Percentages depend on jscpd 5.4.0 defaults (min tokens/lines) and on my file lists; the owner should confirm the exact exclusion list before a config change. Repair estimates are unmeasured. `.savepoint` record duplication (8,403 lines) is excluded as project history; the owner may want a view on whether Check/Issue records should count.

## Drift Notes

Record implemented responsibility or interface changes for planner reconciliation before the mandatory Full Objective Check. If evidence requires a material unknown repair, return REPLAN REQUIRED rather than inventing another scope.
