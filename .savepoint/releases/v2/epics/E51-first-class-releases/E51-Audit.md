---
type: audit-findings
audited: 2026-09-20
---

# Audit Findings: E51 Make releases first-class in V2

## Main Findings

### Verdict

NEEDS WORK. Two repository issues remain: the Release selector can persist a cross-Release router combination that the canonical selection resolver immediately rejects, and a malformed Objective `release` value is silently accepted when the project has no Release records. No owner-run evidence or waiver is otherwise outstanding, but the repository is **NOT READY TO COMMIT/PUSH** until these two contract failures are resolved and re-audited.

### What Needs Attention

#### 1. Switching Releases leaves the router in a rejected context

The board visibly switches from `R001` to `R002`, but its persistence command keeps the old `O001`/`T001` selection even though that Objective belongs to `R001`. After reload, `ResolveSelection` correctly reports `release_mismatch`, so `ResolveNext` abandons the selected Release ladder and returns a global action. The board filter and the Next/resume interpretation can therefore describe different delivery contexts immediately after the normal `r` → Enter workflow.

The switch should clear Objective and Task selection when they do not belong to the new Release, preserve them only when they remain valid, and validate the complete Release/Objective/Task combination before writing. The regression test should assert a diagnostic-free `data.Next` and parity after reload, not only the visible cards and `router.release` field.

Evidence: `go test ./internal/board/v2 -run '^TestReleaseSelectionFiltersIndexedObjectivesAndPersistsOnlyRouterContext$' -count=1 -v` passes while [internal/board/v2/io.go](../../../../../internal/board/v2/io.go) lines 263–267 copy the old selections, lines 282–302 omit the Objective-to-Release check, and [internal/data/next.go](../../../../../internal/data/next.go) lines 147–160 classifies the resulting state as `release_mismatch` before lines 395–402 fall back to the global ladder.

#### 2. Invalid packaging text bypasses the typed Release reference contract

E51 says that a present Objective `release` value is an indexed `R###` reference and that malformed or mismatched references fail closed with a named diagnostic. The loader instead accepts any non-`R###` value whenever the project contains zero Release records. That makes the same Objective valid or invalid depending on whether an unrelated Release record exists, preserves the superseded free-text model as current data, and leaves no diagnostic telling the owner to remove or migrate the value.

The index should always reject a non-empty, non-`R###` reference. If transitional fixtures need compatibility, that conversion belongs in the migration boundary or must produce a named repair diagnostic; it should not silently weaken the live V2 index.

Evidence: `go test ./internal/data -run 'TestDecodeObjectiveV2_preservesTransitionalPackagingText|TestLoadV2Index_rejectsLegacyPackagingTextWhenReleasesExist' -count=1 -v` confirms the deliberate bypass; [internal/data/project.go](../../../../../internal/data/project.go) lines 204–209 skip validation when `len(index.Releases) == 0`, contrary to T001 acceptance criteria 4 and 7.

### Materiality Summary

| Finding | Likelihood | Impact | Materiality | Recommendation |
|---|---|---|---|---|
| Release switch persists a rejected cross-Release selection | High | High | High | Fix now: clear invalid Objective/Task selections, validate Release ownership, and assert board/Next/resume parity after the switch. |
| Invalid free-text Release metadata is accepted in Release-free projects | Medium | Medium | Medium | Fix now: enforce `R###` whenever `release` is present and route transitional conversion through migration or a named diagnostic. |

### What Is Proven / Not Proven

Proven: Release identity/discovery and derived membership for valid records; Release-scoped Check decoding; completion refusals for incomplete Objectives, missing/NEEDS WORK/stale/unknown evidence, unresolved linked Issues, stale acceptance, and missing owner acceptance; historical completion separation; source-preserving managed writes; doctor diagnostics; optional no-Release scaffolds; skill/template parity; migration preview/apply/interruption/conflict/retry/no-op behavior; temporary repository-copy accountability; read-only Release detail; selector navigation/cancel/conflict handling; deterministic ANSI-free plain output; and board/resume wording over valid projections.

Not proven because of the findings: a successful Release switch that leaves one valid canonical router/Next context across board, plain output, and resume; and the universal rule that every present Objective `release` value is a typed indexed reference. T007 and T008 also retain unchecked acceptance boxes despite `status: done`; that bookkeeping did not substitute for implementation evidence.

### Audit Evidence

- Scope lock: T001–T009 acceptance and E51 quality gates; Release/Object/Check/Issue decoders and indexers; completion/cutover/Next resolvers; router and Release writers; migration inventory/plan/preview/apply/recovery; doctor; V2 board selector/detail/TTY/plain paths; resume; agent skills, scaffold assets, and public documentation.
- Coverage and workflow result: exercised absent/one/multiple Releases; unassigned/valid/dangling/malformed/historical membership; planned/in-progress/done lifecycles; missing, NEEDS WORK, stale, unknown, current, and superseded evidence; owner wait/acceptance; unresolved/resolved/excepted Issues; valid/missing/mismatched selections; migration dry-run/apply/interruption/conflict/retry/no-op; Unicode/narrow/read-only rendering and non-TTY output. The two failed cells are the cross-Release Enter workflow and non-`R###` metadata with zero Release records. Network, service, database, authentication, billing, deployment, and publishing cells are not applicable to this local file-only epic.
- Side-effecting workflows: migration ordering covers source inventory, decisions, target/archive construction, preview, staged writes, archive verification, protected replacement/removal, schema activation, manifest/recovery, conflict refusal, retry, and unchanged second apply. Selector ordering covers open, navigate, cancel/detail, optimistic filter, guarded router write, reload, conflict rollback, and removal diagnostics; its successful cross-Release write is the failed workflow row described above.
- File reality and drift: all task context files and new Release source/test files exist; canonical and scaffolded Idea/Design/Check skills are byte-identical. The implemented Release boundary is recorded in Design, v2-Design, README, E50, and the Codebase Map. No unexplained phantom file was found.
- Gates: focused data/doctor/board/resume/init suites passed; the full migration suite passed in 129.883s; `make build && make test` passed, including `go test ./...` and migration in 134.343s; `git diff --check 61cf510..HEAD` passed. `.savepoint/Health-Check.md` is absent, so the optional Full procedure file was skipped.

### Guardrails Verification

- Rule IDs checked: FS-01, FS-03–06; DATA-01–04; TPL-01–04; ARCH-01–04; CFG-01–02; DEP-01–02; TEST-01–08; POL-01–02. REL-01–03 are outside this epic because it changes no distribution target, archive, checksum, or version path.
- Health check mode: Full.
- Evidence: file-write preservation/idempotence and failure paths passed in data and migration suites; dry-run/recovery/conflict tests passed; canonical/template skill parity passed; no new third-party dependency was added; renderer I/O remains behind commands; all required build/test gates passed.
- File reality evidence: scoped paths were checked against the task files and current tree; new `release_v2`, `release_gate_v2`, `release_cutover`, `convert_releases`, and board `releases` source/test files are present.
- Waivers or unresolved findings: no waiver is recorded. Findings 1 and 2 remain unresolved and violate T007/T008 cross-surface selection promises and T001's typed-reference/fail-closed requirements respectively.

### Non-Blocking Observations

The router still names E51/T009 in `task-building`, E51 detail remains `status: planned`, and T009's independent-audit handoff checkbox is open. Those are expected to be reconciled only after the audit is approved and applied; they did not add to the finding count.

## Code Style Review

- [x] STYLE-01 **One job per file** — split files when responsibilities mix.
- [x] STYLE-02 **One job per function** — small, named, testable units.
- [ ] STYLE-03 **Test branches** — cover meaningful conditionals and edge cases. The selector test does not assert the canonical mismatch it creates, and the compatibility test codifies the invalid zero-Release branch without reconciling the acceptance contract.
- [x] STYLE-04 **Types document intent** — prefer explicit types over comments.
- [x] STYLE-05 **Build only what is needed** — no speculative abstractions.
- [x] STYLE-06 **Handle errors at boundaries** — validate inputs, APIs, IO, and external data.
- [x] STYLE-07 **One source of truth** — no duplicated rules, constants, state, or config.
- [x] STYLE-08 **Comments explain why** — not what the code already says.
- [x] STYLE-09 **Content lives in data** — keep copy/config out of logic.
- [x] STYLE-10 **Small diffs** — minimal, reviewable, behaviour-preserving changes.

## Proposed Changes

### Target File
internal/board/v2/io.go

### Replace
```go
		selection := data.RouterSelectionV2{
			Release:   release,
			Objective: router.Objective,
			Task:      router.Task,
		}
```

### With
```go
		selection := data.RouterSelectionV2{
			Release:   release,
			Objective: router.Objective,
			Task:      router.Task,
		}
		if objective := index.Objectives[selection.Objective]; objective == nil || objective.Release != release {
			selection.Objective = ""
			selection.Task = ""
		}
```

Also extend `validateSelectionAgainstIndex` to reject an Objective whose `release` does not equal the selected Release, and update `TestReleaseSelectionFiltersIndexedObjectivesAndPersistsOnlyRouterContext` to require cleared stale selections, no `SelectionDiagnostic`, and Release-scoped Next/board/resume parity after reload.

### Target File
internal/data/project.go

### Replace
```go
		if !releaseIDPatternV2.MatchString(ref) {
			if len(index.Releases) == 0 {
				continue
			}
			return fmt.Errorf("%w: %s: objective %s release %q must match R###", ErrV2InvalidReleaseReference, objective.Source.Path, objective.ID, ref)
		}
```

### With
```go
		if !releaseIDPatternV2.MatchString(ref) {
			return fmt.Errorf("%w: %s: objective %s release %q must match R###", ErrV2InvalidReleaseReference, objective.Source.Path, objective.ID, ref)
		}
```

Replace the transitional free-text acceptance test with an index-level rejection test for a Release-free project. If compatibility with pre-E51 V2 fixtures is still required, handle it explicitly in migration/preview with a named owner decision rather than in the live V2 index.
