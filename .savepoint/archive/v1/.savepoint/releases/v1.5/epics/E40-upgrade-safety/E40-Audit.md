---
type: audit-findings
audited: 2026-07-26
---

# Audit Findings: E40 Upgrade Safety and Backward Compatibility

## Main Findings

### Verdict

CLEAR. All three findings from the initial audit are closed inside the frozen
E40 scope, with no waiver or owner-run evidence outstanding. The epic's
repository changes and required gates are clear. **CLEAR TO COMMIT/PUSH.**

### Prior Finding Closure Map

| Prior finding | Status | Re-audit result |
|---|---|---|
| 1. Unchanged upgrades replace the manifest | Closed | Identical manifest content is not rewritten, and the manifest-writability guard now runs lazily only before the first real write. A no-op upgrade never invokes it and leaves the project tree unchanged. |
| 2. Late write failure can leave partial upgrade state | Closed | Failures retain recoverable ordering and return a truthful `failed` entry naming any backup already written; the CLI prints the partial report, and probe close/removal failures are preserved. |
| 3. Architecture records omit E40 provenance and conflict behavior | Closed | Design.md and the AGENTS.md Codebase Map describe the manifest, conflict and backup policy, router compatibility, and recoverable write ordering. |

### What Was Proven

The manifest has a stable schema, exact-byte hashing, forward-slash keys, the
declared skill-only scope, named malformed-file errors, and idempotent saves.
Fresh scaffolds record the skill bytes actually written. Dry runs neither create
nor change the manifest.

The complete skill decision set passes: missing, identical, tracked-outdated,
customized, pre-manifest, and forced files; fixed `.new` refresh; backup-before-
replace; manifest recording; repeated conflicts; and second-run idempotency.
Agent guides pass absent, marked, unmarked, half-marked, reversed-marker, force,
dry-run, and casing-variant paths. Conflicts and failures lead the report and
name their sidecars.

The lifecycle tests enforce every shipped template's install and upgrade
ownership with force on and off. Project-owned files and install-if-missing
policy/audit assets remain byte-identical, package-owned assets refresh under
their declared rules, and an unchanged full-project upgrade touches neither
files nor directories.

Legacy router, agent-guide, customized-skill, and task-`phase` fixtures remain
readable and follow the migration contract. Unknown router keys remain
tolerated, optional fields default, structural anchors fail clearly when
missing, and the shipped template retains those anchors. README guidance
matches the single-command, optional-dry-run, conflict, backup, force, and
failed-write behavior.

The write workflow is proven at the originally frozen failure timings: before
any write, after a backup but before replacement, after a live replacement, and
while committing the manifest. Applied files retain matching provenance where
possible, user content remains recoverable, failed attempts are visible, and
cleanup failures preserve their primary error.

### Materiality Summary

No materiality actions are required. All prior findings are closed.

### Audit Evidence

- Scope lock: unchanged from the initial audit—every T001-T004 criterion,
  relevant filesystem/data/template/architecture/testing guardrails, the
  declared upgrade entry points and direct filesystem dependencies, and the
  epic quality gates.
- Admission ledger: lazy no-write guard behavior maps exactly to prior Finding
  1's second-run idempotency check; failed backup visibility and probe cleanup
  map exactly to prior Finding 2's named backup/cleanup checks; documentation
  reconciliation maps to prior Finding 3. No new blocking check was admitted.
- Coverage result: the focused E40 init/data suites and CLI partial-report test
  passed uncached. The full original behavior set was rerun rather than only the
  remediation examples.
- File reality and drift: every scoped source, test, fixture, task, and
  documentation path exists. Design.md and AGENTS.md reconcile the recorded
  drift. No phantom file remains.
- Gates: `make build`, full uncached `go test ./... -count=1`, and
  `git diff --check` passed. `internal/init` compiled for Windows amd64 and
  arm64. One earlier uncached run hit the known unrelated
  `TestCopyToClipboard_multipleCalls` transient; the complete uncached rerun
  passed. A separate retry was needed when the sandbox denied access to the Go
  build cache; the host-cache run passed.

### Guardrails Verification

- Rule IDs checked: FS-01..06, DATA-02..03, TPL-01..04, ARCH-01,
  ARCH-03..04, DEP-01..02, TEST-01..08, REL-01.
- Health check mode: Full. `.savepoint/Health-Check.md` is absent, so there are
  no additional project-specific Full commands.
- Result: all mapped rules pass. User-authored content is preserved; dry-run
  and unchanged paths write nothing; replacement order remains recoverable;
  failure and cleanup results are visible; lifecycle vocabulary and legacy
  parsing remain compatible; shipped guidance matches behavior; architecture
  records are current; and the required build/test and portability evidence is
  green.
- Waivers or unresolved findings: none.

### Non-Blocking Observations

`make build-all` still builds Linux and Darwin only although Design says it
includes Windows. E40's changed package compiled directly for both declared
Windows architectures, so this pre-existing Makefile/Design mismatch remains
outside the frozen E40 finding set.

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
