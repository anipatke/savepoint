---
type: epic-design
status: audited
---

# E39: Code-Style Ownership and Policy-Asset Upgrade

## Purpose

Make code style project-owned policy instead of Savepoint-owned boilerplate, and make `.savepoint/` policy files reach existing projects on upgrade.

Today the ten code-style rules live in the managed AGENTS.md block, which `upgrade-assets` overwrites wholesale (`internal/init/agents.go:62-71`). A project that tailors them loses its edits on the next upgrade. The same rules are restated as a hardcoded checklist in the audit skeleton, which the v1.5 Guardrails template forbids: POL-01 says no skill may fail compliance on a rule not defined in `Guardrails.md`, and POL-02 says skill docs may reference guardrails but must not redefine them.

Separately, `upgrade-assets` skips the entire `.savepoint/` subtree (`internal/init/upgrade.go:118-121`), so the new E37 and E38 templates would reach fresh `init` projects only and never an upgraded one. That is a defect independent of code style, and it becomes harmful once AGENTS.md points at a file an upgraded project does not have.

## What this epic adds

- `templates/project/AGENTS.md` and the repo's own `AGENTS.md`: the `## Code Style` body becomes a one-line pointer to the STYLE rules in `.savepoint/Guardrails.md`.
- `savepoint-build-task` reads the STYLE rules during build, closing the gap that no builder skill references code style today.
- The `savepoint-audit-epic` skeleton's `## Code Style Review` sources one checkbox per STYLE rule from `Guardrails.md` instead of hardcoding ten labels.
- A generalized install-if-missing upgrade path so `.savepoint/Guardrails.md` and `.savepoint/Health-Check.md` reach existing projects without overwriting user edits.
- Savepoint's own `.savepoint/Guardrails.md`, dogfooding the E37 template.

## Components and files

| Module | Purpose |
|--------|---------|
| `templates/project/AGENTS.md` and `AGENTS.md` | `## Code Style` body replaced with a pointer to the STYLE rules |
| `agent-skills/savepoint-build-task/SKILL.md` and template copy | Read STYLE rules during build (graceful when `Guardrails.md` is absent) |
| `agent-skills/savepoint-audit-epic/SKILL.md` and template copy | `## Code Style Review` sourced from STYLE rule IDs |
| `internal/init/upgrade.go` | Generalize the install-if-missing gate from `.savepoint/audit/` to a policy-asset allowlist |
| `internal/init/upgrade_test.go` | Install-if-missing, unchanged-when-present, dry-run, and idempotency coverage |
| `.savepoint/Guardrails.md` | Savepoint's own dogfooded policy file |

## Architectural delta

Code-style rules move from Savepoint-owned managed content to project-owned `.savepoint/` policy, matching how Guardrails and Health-Check already work. AGENTS.md retains routing and discoverability only, which is what its own preamble says it is for.

The upgrade change generalizes an existing, tested pattern rather than inventing one. `upgradeAuditAsset` (`internal/init/upgrade.go:260-284`) already installs a missing file, reports `ActionUnchanged` when the file exists, honors dry-run, and never overwrites user content. This epic widens its gate from the single `.savepoint/audit/` prefix to an explicit policy-asset allowlist. The rest of the `.savepoint/` subtree stays skipped.

## Boundaries

**In scope:**

- Replacing the AGENTS.md code-style body with a pointer, in both copies
- Build-task and audit-epic references to STYLE rules
- The install-if-missing upgrade path for named policy assets, with tests
- Savepoint's own `.savepoint/Guardrails.md`

**Out of scope:**

- Authoring the STYLE rules themselves (E37 adds the `STYLE-01..10` category to the Guardrails template)
- Changing the ten rules' wording or adding new ones
- Making code style blocking; it stays Guideline severity and advisory, per `Design.md`
- Installing any other `.savepoint/` file on upgrade
- Rewriting historical release records

## Quality gates

- `make build && make test` passes.
- An upgraded project that lacks `.savepoint/Guardrails.md` or `.savepoint/Health-Check.md` receives them; one that already has them is left byte-identical.
- Dry-run reports the policy-asset installs without touching the filesystem.
- No live AGENTS.md, skill, or template restates the ten style rules; each references the STYLE rule IDs instead.
- Fresh scaffold and upgraded projects both end with code-style guidance reachable from AGENTS.md.

## Open decisions

None. Savepoint dogfoods the Guardrails template rather than keeping its rules inline, so the shipped template is exercised by its own repository.
