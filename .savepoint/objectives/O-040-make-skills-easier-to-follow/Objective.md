---
id: O-040
title: Make the skills easier to follow
status: done
depends_on: [O-039]
release: G-002
priority: medium
rank: 2
---

# O-040: Make the skills easier to follow

## Outcome

An agent reaches every rule it must apply in one step from its skill or from the always-loaded guide. It can see the full scope of a long reference from its first lines, and it tracks a Full Objective Check against a copyable progress checklist. After writing a record, it validates and fixes the record, and the skill descriptions name the `Next` words that select each skill.

## Why

Some reference chains are two or three hops deep: Check → `check-method.md` → AGENTS.md, and Task → AGENTS.md Worktree Lanes → `issue-capture.md`. `check-method.md` is 362 lines with 20 headings and no contents list. The Full Check's required steps are easy to skip without a checklist.

## Success Conditions

- Every reference a workflow skill must apply is named directly in that skill or in the always-loaded AGENTS.md managed block. A recorded audit lists each skill → reference → further-file pointer and how it resolved. No rule needs a chain of two or more non-always-loaded files.
- `check-method.md` and `issue-capture.md`, the references over 100 lines, open with a contents list naming every `##` section in order.
- `check-method.md` has a copyable Full Objective Check progress checklist next to the Quick Check Procedure. It covers every section that Full mode applies, in method order.
- The managed AGENTS.md block holds one write → `savepoint resume` → fix rule for any `.savepoint/` record edit, widened from the existing CLI Rules sentence. Each skill that writes records points to it by section name.
- Each workflow skill's `description` names the `Next` words that select it, matching AGENTS.md's Workflow routing. Idea has no `Next` word and keeps its router-state trigger.
- What each rule requires is unchanged. Live and scaffold copies are byte-identical, and this repository's managed block matches the template. Content tests pin each addition, existing projects receive the changes through `savepoint upgrade-assets`, and `make build && make test-fast` passes.
- The O-034 scenarios are re-walked against the revised text, and the token weight before and after is recorded. Any increase is justified by name.

## Architectural Considerations

AGENTS.md is always loaded through CLAUDE.md, so a pointer to one of its sections costs no extra read. O-039 already turned most of the chains named in Why into direct pointers: the Task skill now names `issue-capture.md` Out-Of-Scope Repair directly. The one-hop work is therefore an audit with fixes only where a chain remains. The managed block comes from `templates/project-v2/AGENTS.md`, and the skills and references have scaffold copies under `templates/project-v2/agent-skills/` (TPL-01). Content tests in `internal/init/agent_skills_test.go`, `skill_validation_test.go`, and `template_freshness_test.go` pin text and parity. Skill `description` frontmatter is what an agent's skill tool matches, so trigger words there improve selection, not routing logic. The Idea state has no `Next` verb in `internal/data`.

## Confirmed Design

Owner-confirmed on 2026-10-10 (planning session plan-o040-20261010):

- **Full Check checklist home.** The copyable checklist lives in `check-method.md`, beside the Quick Check Procedure, so every Check step stays in one file. The Check skill does not repeat it.
- **Write → resume → fix rule home.** The rule is stated once in the managed AGENTS.md block, by widening the CLI Rules sentence about running `savepoint resume` after identity-bearing record changes. The skills keep a one-line pointer. Check's existing step 4 becomes that pointer.
- **Contents lists.** These apply to references over 100 lines (`check-method.md`, `issue-capture.md`), not to the skills.
- **Tasks.** There are two sequential Tasks: (1) reference structure: the one-hop audit, contents lists, and the Full Check checklist; (2) guide and skill wording: the write → resume → fix rule and the trigger words. They share content-test files, so they do not run as parallel lanes.

## Boundaries

**In scope:** one-hop references, contents lists for references over 100 lines, a Full Check progress checklist, a write → `savepoint resume` → fix loop after record edits, and trigger words in skill descriptions.

**Out of scope:** new commands or validators, changed rule content, and work already done by O-039.
