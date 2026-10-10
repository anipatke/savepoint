---
id: T-116
title: Let Claude Code find the Savepoint skills
objective: O-043
status: done
depends_on: [{task: T-115, requires: clear}]
complexity_tier: medium
complexity_reason: Four generated pointer files tracked by the upgrade manifest, plus a parity test against the canonical skills.
owner_validation:
    required: false
    accepted_check: ""
planned_by: {role: planner, session: plan-o043-20261010}
check_waiver:
    task: T-116
    reason: Owner completed this Task via the board without requesting a Task Check.
    actor:
        role: owner
        session: board-owner
    recorded_at: "2026-10-10T07:30:21Z"
---

# Let Claude Code find the Savepoint skills

## Outcome

Each of the four workflow skills has a thin pointer at `.claude/skills/<skill>/SKILL.md` that Claude Code discovers natively. Only its name and description sit in context until it is used, and its body sends Claude to the canonical `agent-skills/<skill>/SKILL.md`. Rules still have one home.

## User Check

In a freshly initialised project, start Claude Code and type `/` or ask "what skills do you have?". The four `savepoint-*` skills should be listed.

## Done When

- Init writes four pointer files. Each carries the canonical skill's exact `name` and `description` and a body of one or two sentences pointing to the canonical file. No rule text is copied.
- A content test fails when a pointer's `name` or `description` differs from its canonical skill, so O-040 description edits cannot drift.
- Upgrade installs missing pointers, refreshes unedited ones, and keeps and reports an owner-edited pointer through the manifest, as for skills (FS-02, TEST-03).
- Shared references get no pointer; they stay reachable from the skills.
- This repo's `.claude/skills/` pointers are generated from the templates and are byte-identical to them.
- `make build && make test-fast` passes.

## Context Files

`main.go`; `internal/init/scaffold.go`; `internal/init/upgrade.go`; `internal/init/upgrade_test.go`; `internal/init/manifest.go`; `internal/init/manifest_test.go`; `internal/init/agent_skills_test.go`; `internal/init/template_freshness_test.go`; `agent-skills/savepoint-idea/SKILL.md`; `agent-skills/savepoint-design/SKILL.md`; `agent-skills/savepoint-task/SKILL.md`; `agent-skills/savepoint-check/SKILL.md`.

## Design References

O-043 Architectural Considerations; Design section 1 (V2 agent workflow assets, upgrade-assets).

## Guardrails

TPL-01, TPL-02, TPL-04, FS-02, FS-03, FS-04, FS-05, TEST-01, TEST-02, TEST-03, TEST-04, TEST-08.

## Implementation Plan

1. Add `templates/project-v2/.claude/skills/<skill>/SKILL.md` pointers. Make sure the embed includes the hidden `.claude` directory.
2. Track pointers in the upgrade manifest alongside skills.
3. Add the name/description parity test.
4. Test fresh install, missing, unedited-outdated, edited, and dry-run upgrade cases.
5. Generate this repo's pointers.

## Boundaries

No changes to canonical skill text. No pointers for `bubbletea-tui-design` or references.

## Technical Verification

Focused tests during iteration; handoff gate from AGENTS.md's Verification Policy.

## Technical Evidence

Executor evidence (claims for a fresh Check to verify):

Criteria:
- Init writes four pointers: met. `templates/project-v2/.claude/skills/<skill>/SKILL.md` for idea, design, task, check; frontmatter lines copied verbatim from the canonical skill; body is one sentence naming `agent-skills/<skill>/SKILL.md`. No rule text copied. `main.go` now embeds `all:templates/project-v2` so the hidden `.claude` directory ships.
- Content test fails on name/description drift: met. `TestClaudeSkillPointers_matchCanonicalNameAndDescription` compares both fields against the canonical skill; `..._bodyOnlyPointsAtCanonicalSkill` checks the body.
- Upgrade installs missing, refreshes unedited, keeps and reports edited: met. `isManifestPath` and `isPackageSkillAsset` now cover the pointers; `TestUpgrade_claudeSkillPointers` covers missing, unedited-outdated, edited (conflict, `.new` written, file kept) and dry-run; `TestScaffold_installsTrackedClaudeSkillPointers` covers fresh install provenance.
- No pointer for references: met. `..._onlyWorkflowSkillsHavePointers` asserts exactly four pointer directories (also none for bubbletea-tui-design).
- This repo's pointers byte-identical to templates: met. `.claude/skills/` copied from the templates; `..._repoCopiesAreByteIdenticalToTemplates` enforces it.
- `make build && make test-fast`: passed (exit 0).

Commands run: `go test ./internal/init/` (iteration), `make build && make test-fast`.

Files changed: `main.go`, `internal/init/manifest.go`, `internal/init/upgrade.go`, `internal/init/lifecycle_test.go`, `main_test.go`; added `internal/init/claude_skills_test.go`, `templates/project-v2/.claude/skills/*/SKILL.md`, `.claude/skills/*/SKILL.md`.

Extra reads (outside Context Files): `internal/init/lifecycle_test.go` and `main_test.go`, because existing ownership and manifest-scope expectations failed once pointers became manifest-tracked and had to be updated. `internal/init/claude_guide_test.go`/`agents.go` read for T-115 helper names.

Limitations: not verified inside a live Claude Code session (the User Check: `/` lists the four skills); full gate (`make test-full`) not run, as this is not migration/platform-sensitive. Whether Claude Code accepts the exact description text as-is was not tested.

## Drift Notes

None yet.
