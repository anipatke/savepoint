---
id: O-039
title: Cut duplicated guidance
status: planned
depends_on: []
release: G-002
priority: medium
rank: 1
---

# O-039: Cut duplicated guidance

## Outcome

Every session loads each workflow rule once. This repository's AGENTS.md keeps only repository-specific rules above the managed Savepoint block. Shared rules (Goal Context, Task-check waivers, closure, and dependency rules) have one home, and the skills and references point to it. Active guidance no longer narrates V1/V2 history, and it uses one term for each role and record.

## Why

This repository's AGENTS.md is a hand-written guide (lines 1–178) followed by the managed block that upgrade merged in (commit 2527047). The two copies already differ: only the managed block has Existing Codebase Adoption, and only the repository copy has the `health report` human-only rule. The Goal Context paragraph is repeated in Idea, Task, Check, Design, and AGENTS.md. The waiver and closure rules appear in AGENTS.md Verification Policy, the Check skill's Closure Rules, and `check-method.md`. Repeated text costs context in every session and drifts.

## Success Conditions

- This repository's AGENTS.md has no section that repeats the managed block. What remains above the block is repository-specific: make gate commands, the Codebase Map, runtime gate code pointers, the build-from-source rule for `savepoint`, and Reporting to the Owner.
- Goal Context and the Verification Policy rules (waiver, `clear`/`accepted` dependencies, Objective Check coverage) each have one home in the managed AGENTS.md block. The skills and references keep only role-specific consequences and cite the section by name.
- Active skills, references, and the managed block contain no V1/V2 migration narrative beyond what `savepoint migrate` users need, and that lives in one place.
- Prose uses "owner" for the human decision-maker and "Goal" for the record; `release:` appears only as the field name.
- Live and scaffold copies stay byte-identical. Content tests pin each moved rule at its new single home, and `make build && make test-fast` passes.
- The token weight of the always-loaded guide plus the four skills and three references is recorded before and after. Any increase is justified by name.

## Architectural Considerations

AGENTS.md is always loaded through CLAUDE.md, so a pointer from a skill to an AGENTS.md section costs no extra read. The managed block (`<!-- SAVEPOINT:BEGIN/END -->`) is package-owned and comes from `templates/project-v2/AGENTS.md`. The text above it is project-owned. Content tests in `internal/init/agent_skills_test.go` and `template_freshness_test.go` pin rules and parity, so moving a rule moves its assertion.

## Confirmed Design

Owner-confirmed on 2026-10-06 (planning session plan-o039-20261006):

- **This repository's AGENTS.md.** Above the managed block, keep only repository-specific rules: make gate commands, the Codebase Map, runtime gate code pointers (`ResolveTaskDependencyV2`, `CheckWaiver`), the build-from-source rule for `savepoint`, and Reporting to the Owner. Remove every section that repeats the managed block.
- **One home per shared rule.** Goal Context and the Verification Policy rules (waiver, `clear`/`accepted` dependencies, Objective Check coverage) live once in the managed AGENTS.md block. AGENTS.md is always loaded, so skills point to it by section name at no read cost.
- **Migration history.** The V1/V2 migration narrative stays only in the managed block's Required Goal Context.
- **Drift repair.** The `savepoint health report` human-only CLI rule moves from this repository's half into the managed block template, so every project receives it.
- **Terms.** Use "owner" for the human decision-maker and "Goal" for the record; `release:` appears only as the field name.
- **Tasks.** Three Tasks run in sequence: trim AGENTS.md, give shared rules one home, then remove history and settle terms. Each records token weight before and after and passes `make build && make test-fast`.

## Boundaries

**In scope:** this repository's AGENTS.md above the managed block; the managed block template; the four workflow skills, three references, and scaffold copies; the matching content tests.

**Out of scope:** changing what any rule requires; restructuring references or adding checklists (O-040); scenario sets (O-041); changing how upgrade merges a managed block into an existing guide (record as an Issue if needed).
