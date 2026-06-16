---
id: E21-document-template-optimisation/T001-agents-md-template-cleanup
status: done
objective: Audit and tighten the scaffolded AGENTS.md template so it reflects the E20 lifecycle contract and contains no stale or duplicate guidance.
complexity_tier: low
complexity_reason: Single-file edit with a clear target; alignment check against live AGENTS.md.
---

# T001: AGENTS.md Template Cleanup

## Problem

The scaffolded `templates/project/AGENTS.md` was last meaningfully updated in E20 alongside the live file, but may still carry sections that duplicate phase-skill prompts, reference stale terminology (`phase`, `stage: implementation`), or include guidance that belongs in skill files rather than the entry-point agent guide.

## Context Files

- `templates/project/AGENTS.md`
- `AGENTS.md`
- `internal/init/template_freshness_test.go`

## Acceptance Criteria

- [x] Scaffolded `AGENTS.md` contains no references to stale lifecycle terms (`phase`, `stage: implementation`, legacy status values).
- [x] Guidance that duplicates phase-skill prompt instructions is removed; the file routes to skills, it does not re-state them.
- [x] Live `AGENTS.md` and scaffolded `AGENTS.md` are consistent on all lifecycle and terminology guidance.
- [x] Template freshness test passes without modification (or is updated only to reflect intentional structural changes).
- [x] `make build && make test` passes.

## Implementation Plan

- [x] Read `templates/project/AGENTS.md` and `AGENTS.md` side by side; list every divergence and every section that duplicates skill-level prompt content.
- [x] Remove or condense duplicate and stale sections from the scaffolded template.
- [x] Mirror any lifecycle or terminology corrections to the live `AGENTS.md` where applicable.
- [x] Run `go test ./internal/init` to verify freshness tests.
- [x] Run `make build && make test`.

## Context Log

### Divergences found

1. Template `## Codebase Map` table had stray `Epic` column with no rows; live uses `Module | Purpose`.
2. Template was missing `## Context Budget` section entirely.
3. Template `## Audit` was a single section; live split it into `## Audit Handoff` + `## Audit File Structure`.
4. Template `## Implementation` had 5 steps; live had 6 (live added "Stop. Prompt user before continuing").
5. Template `## Build` was a single line `Build gate: ...`; live used a fenced bash block.
6. Both files still echoed the build-task and savepoint-audit skill prompt content in the `## Implementation` and `## Audit` sections, violating the file's own principle ("do not duplicate phase-by-phase prompt instructions").

### Changes applied

- `templates/project/AGENTS.md`:
  - Added `## Context Budget` section (mirroring live).
  - Replaced the strayed `## Codebase Map` `Module | Epic | Purpose` header with `Module | Purpose`.
  - Replaced `## Build` plain line with a fenced `bash` block (`make build && make test`).
  - Condensed `## Implementation` to a routing line that points to `savepoint-build-task` plus the user-prompt rule.
  - Condensed `## Audit` to a routing line that points to `savepoint-audit` plus the two required file/section anchors.
- `AGENTS.md` (live):
  - Condensed `## Implementation` to match the template.
  - Merged `## Audit Handoff` + `## Audit File Structure` into a single `## Audit` section that routes to the skill, matching the template.

### Verification

- `go test ./internal/init/...` passes; all canonical strings (stage required, never `stage: implementation`, "Only the user may set a task to `status: done`", "make build && make test", "During audit apply/close…") remain present in both files.
- `make build && make test` passes.
