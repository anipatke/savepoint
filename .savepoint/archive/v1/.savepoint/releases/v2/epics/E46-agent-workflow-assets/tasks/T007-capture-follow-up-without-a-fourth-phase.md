---
id: E46-agent-workflow-assets/T007-capture-follow-up-without-a-fourth-phase
title: Capture follow-up without a fourth phase
status: done
objective: Add the shared Issue capture guidance the three working skills enter from, keeping Defect as a user word and adding no public phase.
depends_on:
    - E46-agent-workflow-assets/T006-check-in-a-fresh-session-and-record-the-proof
complexity_tier: medium
complexity_reason: A new shared reference plus entry wording in three skill pairs; the risk is duplicating Issue rules in four places.
---

# T007: Capture follow-up without a fourth phase

## Problem

V1 makes defect capture its own skill and its own router state. V2 does not: an Issue is a record, not a phase. Follow-up gets captured from wherever it is noticed — a planner spotting drift, an executor hitting something out of scope, a checker recording a blocker — and the entry has to read the same from all three without each skill carrying its own copy of the rules.

The rules that matter are the ones that stop an Issue from being a lie. Search before creating, so the same problem does not acquire three identities. Record what proof would close it, so `resolved` means something. Distinguish `verified` from `accepted` from `duplicate`, because two of those are not repairs. Append history rather than editing it. And route out-of-scope repair into a Task instead of letting the Issue become a parallel backlog.

"Defect" stays a word the user says. It maps to `type: defect`; it does not resurrect a separate record kind.

## Context Files

- `agent-skills/references/issue-capture.md`
- `templates/project/agent-skills/references/issue-capture.md`
- `agent-skills/savepoint-design/SKILL.md`
- `agent-skills/savepoint-task/SKILL.md`
- `agent-skills/savepoint-check/SKILL.md`
- `templates/project/agent-skills/savepoint-design/SKILL.md`
- `templates/project/agent-skills/savepoint-task/SKILL.md`
- `templates/project/agent-skills/savepoint-check/SKILL.md`
- `agent-skills/savepoint-create-defect/SKILL.md`
- `internal/init/agent_skills_test.go`
- `internal/data/issue_v2.go`
- `.savepoint/releases/v2/v2-Design.md`

## Acceptance Criteria

- [x] `agent-skills/references/issue-capture.md` exists with `type: issue-capture-reference`, `triggerable: false`, no `name` key, and a byte-identical shipped copy.
- [x] It carries the Issue artifact template matching `internal/data/issue_v2.go`: `id`, `title`, `type: defect|drift|guardrail|verification|other`, `status: open|in_progress|resolved`, `source`, `tasks`, `checks`, `guardrail_ids`, optional severity, `resolution`, `duplicate_of`, `history`; body sections Summary, Evidence, Proof Needed.
- [x] It states that `type` is descriptive and never sufficient on its own to block a Task.
- [x] It states the search-first rule: look for an existing Issue by symptom, location, violated requirement, and linked work before allocating an ID, and that no automatic deduplication is assumed.
- [x] It defines the resolution dispositions — `verified` requires proof from a Check, `accepted` is an explicit owner decision and not a repair, `duplicate` names the canonical Issue and proves nothing — and states that reopening reuses the same ID with dated evidence.
- [x] It states that history is an append-only frontmatter list of `{at, actor, kind, note, check}` entries and that a write which shortens, reorders, or edits a recorded entry is refused.
- [x] It states that deferral is a dated history entry on an open Issue, not a fourth lifecycle state.
- [x] It states the role boundaries: executors report repair evidence without closing, checkers verify proof and close, owners decide acceptance.
- [x] It states that out-of-scope repair becomes a new bounded Task in an Objective rather than work performed inside the Issue.
- [x] It states that "defect" is a user-facing word mapping to `type: defect`, and that no separate defect record, state, or public phase exists in V2.
- [x] Each of `savepoint-design`, `savepoint-task`, and `savepoint-check` names Issue capture as an entry from its own workflow and points at this reference instead of restating the rules; the rule prose appears in exactly one file.
- [x] `internal/init/agent_skills_test.go` asserts the reference frontmatter contract, the template fields, the disposition and history rules, the three skill entry points, and that the Issue rules are not duplicated into the skills.
- [x] `agent-skills/savepoint-create-defect/SKILL.md` is unchanged and still present in both trees.

## Implementation Plan

- [x] Read design section 6 and `internal/data/issue_v2.go` to match field names, dispositions, and history entry shape exactly.
- [x] Write `agent-skills/references/issue-capture.md` with the template, the search-first rule, dispositions, history rules, role boundaries, and the repair-routing rule.
- [x] Add a short Issue capture entry to each of the three working skills, pointing at the reference by path.
- [x] Mirror all four changed files to the template tree.
- [x] Extend `internal/init/agent_skills_test.go` with the reference contract case and the no-duplication assertion across the three skills.
- [x] Run `go test ./internal/init/... ./internal/data/...`, then `make build && make test`.

## Context Log

**Files read:** `.savepoint/releases/v2/epics/E46-agent-workflow-assets/E46-Detail.md`, `.savepoint/releases/v2/v2-Design.md` (section 6), `internal/data/issue_v2.go`, `agent-skills/savepoint-design/SKILL.md`, `agent-skills/savepoint-task/SKILL.md`, `agent-skills/savepoint-check/SKILL.md`, `agent-skills/savepoint-create-defect/SKILL.md`, `agent-skills/references/check-method.md`, `internal/init/agent_skills_test.go`, `internal/init/skill_validation_test.go`, `internal/init/template_freshness_test.go`. No reads beyond the Context Files plus these test-helper files (needed to reuse `skillRoots`/`frontmatterField`/`sectionBody`/`assertFileMatches` without duplicating them).

**Files written:** `agent-skills/references/issue-capture.md` (new), `agent-skills/savepoint-design/SKILL.md` (added `## Issue Capture` entry), `agent-skills/savepoint-task/SKILL.md` (added `## Issue Capture` entry), `agent-skills/savepoint-check/SKILL.md` (added `## Issue Capture` entry), mirrored byte-identical to `templates/project/agent-skills/...`, and `internal/init/agent_skills_test.go` extended with the reference contract, artifact template, disposition/history, role-boundary, search-first, defect-word-mapping, live/template-match, three-skill entry-point, no-duplication, and create-defect presence tests.

**Commands run:** `go test ./internal/init/... ./internal/data/...` — pass. `make build && make test` — pass (all packages ok).

**Limitations:** No manual byte-diff check was needed beyond the automated `assertFileMatches` tests, which cover live/template parity for all four changed files.
