---
type: audit-findings
audited: 2026-07-26
---

# Audit Findings: E39 Code-Style Ownership and Policy-Asset Upgrade

## Main Findings

### Verdict

CLEAR — both findings applied. The audit raised two guidance defects: the audit skeleton still embedded two STYLE rule IDs despite E39 making those IDs project-owned, and the project architecture still described the pre-E39 upgrade boundary. Both proposed changes are applied and re-verified below. The implementation and all gates pass, and the repository is ready to commit.

### What Was Applied

#### The audit skeleton no longer hardcodes project-owned rule IDs

The skeleton's `## Code Style Review` list now reads `{STYLE rule ID} {rule wording from .savepoint/Guardrails.md}` plus a repeat-per-rule instruction, so a project that renumbers or replaces its style rules is never shown Savepoint's labels. Applied identically in `agent-skills/savepoint-audit-epic/SKILL.md` and `templates/project/agent-skills/savepoint-audit-epic/SKILL.md`; the two copies remain byte-identical (`diff` clean), satisfying TPL-01 and T001's "no hardcoded rule labels" criterion. No `STYLE-0*` literal remains anywhere under `agent-skills/` or `templates/project/agent-skills/`.

#### The architecture now states the policy-asset upgrade boundary

`.savepoint/Design.md:26` now reads that `upgrade-assets` "skips project-owned state except for two non-destructive paths": missing `.savepoint/Guardrails.md` and `.savepoint/Health-Check.md` are installed without overwriting existing copies, alongside the pre-existing migrations-archive exception. The canonical design record now matches `internal/init/upgrade.go` and README, closing T002's recorded drift note and satisfying TPL-02.

### Materiality Summary

| Finding | Likelihood | Impact | Materiality | Outcome |
|---|---|---|---|---|
| Audit skeleton hardcodes `STYLE-01` and `STYLE-02` | High — every scaffolded audit skill receives the examples | Medium — project-owned policy can be misrepresented in generated audit guidance | Medium | Applied — rule-neutral placeholder in both byte-identical skill copies. |
| Design omits the policy-asset upgrade exception | High — the current architecture is always stale for this path | Low — runtime is correct, but maintainers receive an inaccurate ownership boundary | Low | Applied — architecture sentence now names both install-if-missing policy assets. |

### What Is Proven / Not Proven

Proven:

- Both AGENTS.md copies point to project-owned STYLE rules and no longer restate the ten-rule body.
- The repository Guardrails file and the project template define `STYLE-01..10` as advisory Guideline rules; the build and audit skills do not make style blocking.
- Both canonical/template skill pairs are byte-identical, and absence of Guardrails has an explicit non-failing fallback.
- The exact two-file policy allowlist installs missing Guardrails and Health-Check assets, while other `.savepoint/` files remain skipped.
- Existing policy files remain byte-identical, including locally modified bytes and when `force` is true; repeated upgrades report unchanged.
- Dry-run reports installs without writing, policy templates are interpolated consistently with fresh scaffold, audit assets retain their existing action behavior, and the installed AGENTS.md pointer resolves to the installed STYLE rules.
- Missing or unwritable targets fail clearly; the independent unwritable-directory scenario left no policy file behind.

Previously not proven, now resolved by the applied changes:

- T001's “no hardcoded rule labels” criterion is satisfied: no `STYLE-0*` literal remains in either audit-skill copy, and the pair is byte-identical.
- Design drift is reconciled, so TPL-02 holds for the architecture description.

Re-verified after apply: `make build && make test` pass (all packages ok).

### Audit Evidence

- Scope lock: T001 and T002 acceptance criteria; the two AGENTS guides, two build-skill copies, two audit-skill copies, Guardrails/Health-Check policy assets, Design/README guidance, and the `UpgradeProjectAssets` public workflow with its direct scaffold interpolation and atomic-write dependencies.
- Coverage and workflow result: classified both policy assets across missing, pristine, user-modified, repeat, dry-run, force, rendered-template, and unwritable states; classified non-allowlisted `.savepoint/` paths and existing audit assets; traced target validation → template walk → exact action selection → existence guard → dry-run return or render/write → sorted report. All runtime cells passed.
- File reality and drift: every file named in both task context logs exists, including root `agent_skills_test.go`; canonical skill pairs compare byte-identical. T001's Design update is present. T002's recorded Design drift remains and is Finding 2.
- Gates: focused upgrade tests passed (`go test ./internal/init -run 'TestUpgrade(ProjectAssets|DeliversPolicyAssets|Report)' -count=1`); independent policy-upgrade harness passed; `go vet ./internal/init/`, `git diff --check`, `make build`, and `make test` all passed.

### Guardrails Verification

- Rule IDs checked: FS-01 through FS-06; TPL-01 through TPL-04; TEST-01 through TEST-04, TEST-06 through TEST-08; POL-01 and POL-02.
- Health check mode: Full. This repository has no `.savepoint/Health-Check.md`, so the project-specific Full checklist was skipped as required; acceptance, guardrail, file-reality, drift, focused, and full-gate evidence was still completed.
- Evidence: user-file preservation, exact allowlisting, dry-run behavior, idempotency, path construction, failure handling, template parity, absent-policy fallback, and named test outcomes are summarized above.
- File reality evidence: no unexplained phantom files or discarded scratch files were found.
- Waivers or unresolved findings: no waivers. TPL-02 remains unresolved through Finding 2; T001 acceptance remains unresolved through Finding 1. No Blocker-severity guardrail failed.

### Non-Blocking Observations

The router still names T002 under `task-building` and the epic detail remains `status: planned` even though both task files are `done`. The explicit request permits this completed-epic audit, and lifecycle records are owner-controlled, so this does not change the verdict or remediation scope.

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

## Proposed Changes

### Target File
agent-skills/savepoint-audit-epic/SKILL.md

### Replace
```md
- [ ] STYLE-01 rule wording from `.savepoint/Guardrails.md`
- [ ] STYLE-02 rule wording from `.savepoint/Guardrails.md`
- [ ] One checkbox per remaining STYLE rule, in file order
```

### With
```md
- [ ] {STYLE rule ID} {rule wording from `.savepoint/Guardrails.md`}
- [ ] Repeat once per STYLE rule, in file order
```

### Target File
templates/project/agent-skills/savepoint-audit-epic/SKILL.md

### Replace
```md
- [ ] STYLE-01 rule wording from `.savepoint/Guardrails.md`
- [ ] STYLE-02 rule wording from `.savepoint/Guardrails.md`
- [ ] One checkbox per remaining STYLE rule, in file order
```

### With
```md
- [ ] {STYLE rule ID} {rule wording from `.savepoint/Guardrails.md`}
- [ ] Repeat once per STYLE rule, in file order
```

### Target File
.savepoint/Design.md

### Replace
```md
- **Upgrade-assets command** (`savepoint upgrade-assets [dir] [--dry-run] [--force]`) refreshes package-owned `agent-skills/**/SKILL.md` files, shared references under `agent-skills/references/`, and the managed block in the root agent guide from embedded templates for existing Savepoint projects. It skips project-owned state, except that a retired generic audit skill is preserved under the non-triggerable `.savepoint/migrations/` archive before its active copy is removed.
```

### With
```md
- **Upgrade-assets command** (`savepoint upgrade-assets [dir] [--dry-run] [--force]`) refreshes package-owned `agent-skills/**/SKILL.md` files, shared references under `agent-skills/references/`, and the managed block in the root agent guide from embedded templates for existing Savepoint projects. It skips project-owned state except for two non-destructive paths: missing `.savepoint/Guardrails.md` and `.savepoint/Health-Check.md` policy assets are installed without overwriting existing copies, and a retired generic audit skill is preserved under the non-triggerable `.savepoint/migrations/` archive before its active copy is removed.
```
