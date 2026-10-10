---
id: I-144
title: O-039 wording leftovers — "V2 Task Markdown" and "names a Release"
type: defect
status: resolved
source:
  kind: check
  check: C-971
  actor: {role: checker, session: check-o039-20261010}
  at: '2026-10-10T07:05:00Z'
tasks: [T-114]
checks: [C-971, C-972]
guardrail_ids: [TPL-01]
severity: low
resolution:
  disposition: verified
  check: C-972
  actor: {role: checker, session: recheck-o039-20261010}
  at: '2026-10-10T07:00:00Z'
  reason: Both phrases replaced live and scaffold, block parity identical, absence test pinned and mutation-proven, full gate passes.
history:
  - at: '2026-10-10T07:05:00Z'
    actor: {role: checker, session: check-o039-20261010}
    kind: observed
    note: Found by the Full Objective Check of O-039.
    check: C-971
  - at: '2026-10-10T07:00:00Z'
    actor: {role: checker, session: recheck-o039-20261010}
    kind: rechecked
    check: C-972
    note: Both phrases replaced live and scaffold, block parity identical, absence test pinned and mutation-proven, full gate passes.
---

# I-144: O-039 wording leftovers — "V2 Task Markdown" and "names a Release"

## Summary

T-114 set out to remove V2 qualifiers and to stop using "Release" for a Goal in prose. Two spots were missed, and the T-114 evidence claims one of them was removed.

1. `agent-skills/savepoint-design/SKILL.md:70` (and its scaffold copy) still says "Prepare the complete V2 Task Markdown as an ID-free draft." T-114 Technical Evidence lists "V2 Task Markdown" among removed qualifiers. Violates T-114 Done When 1 (evidence accuracy) and O-039 Success Condition 3.
2. The managed block's Router Selection (`templates/project-v2/AGENTS.md`, mirrored in this repository's `AGENTS.md`) still says "leaving `release:` byte-for-byte unchanged unless the owner names a Release as part of an explicit Goal choice." Here "Release" names a Goal in prose. Violates T-114 Done When 2 and O-039 Success Condition 4.

## Evidence

- `grep -n "V2 Task Markdown" agent-skills/savepoint-design/SKILL.md templates/project-v2/agent-skills/savepoint-design/SKILL.md` → one hit in each.
- `grep -n "names a Release" AGENTS.md templates/project-v2/AGENTS.md` → one hit in each.
- `TestActiveGuidanceKeepsNoHistoryAndOneTermPerRole` (`internal/init/template_freshness_test.go:543`) checks "V2 index", "V2 record", "the V2 workflow" but not "V2 Task Markdown", and has no check for "Release" naming a Goal, so the gate passes.

## Proof Needed

- Both phrases replaced (for example "Task Markdown" and "names a Goal"), live and scaffold copies byte-identical, this repository's managed block refreshed to match the template.
- `TestActiveGuidanceKeepsNoHistoryAndOneTermPerRole` asserts both old phrases are absent.
- `make build && make test-fast` passes; a re-check confirms.
