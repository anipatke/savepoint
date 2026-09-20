---
type: audit-findings
audited: 2026-09-19
applied: 2026-09-19
re_audit_of: 2026-09-18
---

# Audit Findings: E45 Convert active work without losing history

## Main Findings

### Verdict

**CLEAR — all three approved remediation findings and the follow-up FS-01 probe-collision blocker are closed. CLEAR TO COMMIT/PUSH.** The migration now protects late-created additive destinations, recovers from the recorded operation across separate command invocations, keeps preview and recovery-report paths read-only, and makes the apply-only writeability probe collision-safe. The required focused and repository gates pass; no waiver or owner-run evidence is needed to close these repository findings.

### Prior Finding Closure Map

| Prior finding | Applied outcome | Evidence |
|---|---|---|
| 1. Late create destination could overwrite user content | **Closed** | `ActionCreate` now uses atomic no-replace publication on Unix and Windows, and occupied destinations return `ErrCreateDestinationExists` without changing their bytes. |
| 2. Cross-invocation recovery could not reproduce the pending operation | **Closed** | `operation.yml` records the original plan and all outputs are staged before live installation; a separate `RunCommand` with different clock and operation-ID inputs resumes the recorded operation. |
| 3. Preview performed a write probe | **Closed** | Target resolution is read-only; writeability probing runs only for apply. Read-only default preview and recovery-report tests detect that the probe is not invoked. |
| 4. Apply writeability probe could overwrite a pre-existing sentinel | **Closed** | Apply-only preflight uses `os.CreateTemp` for an exclusive unique path, removes only that path, and `TestRunCommand_applyWriteabilityProbePreservesPreExistingSentinelAndPrefix` preserves exact bytes for the legacy sentinel and a same-prefix collision while detecting leaked transient files. |

### Closure Status

Nothing remains requiring remediation for E45 or the follow-up FS-01 probe-collision blocker.

### Materiality Summary

No unresolved materiality actions remain. The former FS-01, recovery, FS-03, and TEST-05 concerns, plus the follow-up FS-01 probe-collision blocker, are covered by the applied code and regression evidence below.

### What Is Proven / Not Proven

**Proven:**

- Plan-time collision checks are backed by installation-time create-only checks for target records, relocated documents, and archive entries; unrelated bytes survive a late destination appearing after Plan.
- The operation journal persists the original plan, including decisions, generated time, operation ID, source hashes, and ordered entries. All complete outputs are staged before the first live install, and recovery consumes those staged bytes instead of replanning.
- Separate command calls cover interrupted apply followed by `Recover+Write`; changing a source between calls returns `ErrSourceConflict` and preserves the edit.
- Default preview, explicit dry-run, and recovery-report paths perform no writeability probe and work on a readable `0555` project. Apply-only writeability validation remains named and tested.
- Apply-mode writeability probing creates an exclusive unique temporary file, cleans up only that returned path, and leaves pre-existing legacy sentinel and same-prefix files byte-identical; the regression also detects a leaked transient probe file.

**No repository criterion remains unproven for the approved findings and follow-up blocker.** The Windows path is covered by the platform-specific `MoveFileExW` no-replace implementation and the required `GOOS=windows GOARCH=amd64` build; host-specific Windows execution remains outside this Linux test environment.

### Audit Evidence

- **Scope lock:** reused the 2026-09-18 E45 scope: all twelve task criteria, the public Plan/Apply/RunCommand/ResolveTarget surfaces, operation lifecycle, command preview/apply/recover forms, and the mapped filesystem, architecture, data, testing, release, and style rules.
- **Coverage:** `TestApply_lateCreateDestinationsAreCreateOnly` covers target, document, and archive classes; `TestOperation_createInstallRejectsOccupiedDestinationAndPreservesBytes` covers the operation boundary; `TestRunCommand_recoverApplyAcrossInvocationsUsesRecordedOperation` proves separate-call recovery with original volatile fields; `TestRunCommand_recoverApplyAcrossInvocationsRejectsEditedSource` proves user-edit conflict; `TestRunCommand_previewAndRecoveryReportAreReadOnlyOnReadable0555Project` proves no apply-only probe in read-only modes; `TestRunCommand_applyWriteabilityProbePreservesPreExistingSentinelAndPrefix` proves apply-mode probe collision safety and cleanup. Existing migration, fixture, resume, archive, and platform replacement tests also pass.
- **Workflow and side effects:** fresh apply backs up replacements/removals, stages every output, then installs in the recorded order; create publication cannot replace a late destination. The apply-only writeability preflight creates and removes only an exclusive unique temporary file. Recovery validates untouched sources and installed outputs before continuing.
- **File reality and drift:** all scoped files exist, and `internal/migrate` is already present in the Codebase Map. No router, E45 task lifecycle, or design record was changed. The pre-existing router, `internal/data/write_test.go`, D001 defect, and E46 worktree files were preserved.
- **Gates:** `GOCACHE=/tmp/savepoint-gocache go test ./internal/migrate -count=1`, `GOCACHE=/tmp/savepoint-gocache go vet ./...`, `GOCACHE=/tmp/savepoint-gocache GOOS=windows GOARCH=amd64 go build ./...`, `GOCACHE=/tmp/savepoint-gocache make build`, `GOCACHE=/tmp/savepoint-gocache make test`, and `git diff --check` all pass.

### Guardrails Verification

- **Rule IDs checked:** FS-01, FS-03, FS-04, FS-05, FS-06, DATA-01, DATA-03, ARCH-01, ARCH-03, ARCH-04, CFG-02, DEP-01, TEST-01, TEST-02, TEST-03, TEST-04, TEST-05, TEST-06, TEST-08, REL-01, TPL-01, STYLE-01..10.
- **Health check mode:** Full. `.savepoint/Health-Check.md` and `.savepoint/audit/` are absent, so their optional procedures were skipped.
- **Result:** FS-01 (including the follow-up probe-collision blocker), FS-03, FS-04, FS-06, ARCH-03, CFG-02, TEST-01, TEST-02, TEST-03, TEST-05, TEST-08, and REL-01 have direct remediation or gate evidence above. Other mapped rules retain the prior pass evidence. No blocker or required waiver remains.

## Code Style Review

- [x] STYLE-01 **One job per file** — split files when responsibilities mix.
- [x] STYLE-02 **One job per function** — small, named, testable units.
- [x] STYLE-03 **Test branches** — cover meaningful conditionals and edge cases.
- [x] STYLE-04 **Types document intent** — prefer explicit types over comments.
- [ ] STYLE-05 **Build only what is needed** — no speculative abstractions. The pre-existing exported `WriteManifestCreateOnly` helper remains production-dead; this is advisory and does not block closure.
- [x] STYLE-06 **Handle errors at boundaries** — validate inputs, APIs, IO, and external data.
- [x] STYLE-07 **One source of truth** — no duplicated rules, constants, state, or config.
- [x] STYLE-08 **Comments explain why** — not what the code already says.
- [x] STYLE-09 **Content lives in data** — keep copy/config out of logic.
- [x] STYLE-10 **Small diffs** — minimal, reviewable, behaviour-preserving changes.
