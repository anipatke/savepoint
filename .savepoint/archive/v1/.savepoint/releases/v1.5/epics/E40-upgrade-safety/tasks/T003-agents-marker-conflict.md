---
id: E40-upgrade-safety/T003-agents-marker-conflict
status: done
objective: Stop upgrade appending a second managed block to an AGENTS.md that has no markers — leave the file untouched, write the merged result as AGENTS.md.new, and report a conflict.
depends_on: []
complexity_tier: low
complexity_reason: One fallback branch in replaceManagedBlock plus its upgrade caller and tests.
---

# T003: AGENTS.md Marker Conflict

## Problem

`replaceManagedBlock` (`agents.go:62-71`) swaps the managed block in place when `<!-- SAVEPOINT:BEGIN -->` and `<!-- SAVEPOINT:END -->` are both present. When they are absent it falls through to `:69-70` and appends the block to the end of the existing file.

So a hand-written `AGENTS.md`, or one written before the markers existed, ends up holding two sets of workflow instructions: the user's original above, unmarked and never refreshed again, and the managed block below. The two drift further apart with every release. For the file whose only job is instructing agents, silently duplicated and contradictory guidance is the worst available outcome — and it is reported as a routine `merged`.

## Context Files

- `internal/init/agents.go`
- `internal/init/upgrade.go` (the `AGENTS.md` branch, lines 137-157 dry-run and 207-233 write)
- `internal/init/agents_test.go`
- `internal/init/upgrade_test.go`

## Acceptance Criteria

- [x] `AGENTS.md` absent → managed block written as the whole file, reported `updated`. Unchanged behaviour.
- [x] `AGENTS.md` present with both markers → block replaced in place, all content outside the markers byte-identical afterwards, reported `merged` or `unchanged` as today. Unchanged behaviour.
- [x] `AGENTS.md` present with **no** markers → the file is left byte-identical, the merged result is written to `AGENTS.md.new`, and the result is reported `conflict`.
- [x] A file containing only one marker of the pair is treated as having no markers, and takes the conflict path rather than corrupting the file.
- [x] `--force` on the no-marker case appends the block as today, after saving the original to `AGENTS.md.bak`, reported `merged`.
- [x] Dry run reports the same action for every case above and writes nothing.
- [x] The casing variants resolved by `FindAgentGuide` (`agents.go:15-21`) keep working — the conflict path writes `<actual-filename>.new`, preserving on-disk casing.
- [x] `make build && make test` passes.

## Implementation Plan

- [x] Change `replaceManagedBlock` to report whether markers were found rather than silently appending — return the merged string plus a found flag, or split the lookup into its own small function.
- [x] Update `MergeAgentGuide` and both `AGENTS.md` branches in `upgrade.go` to take the conflict path when markers are absent and `force` is false.
- [x] Keep the append behaviour reachable under `--force`, so there is still a way to adopt an existing file deliberately.
- [x] Extend `agents_test.go` and `upgrade_test.go` for the no-marker, half-marker, force, and dry-run cases.

## Notes

Conflict rather than auto-adopt is deliberate. Guessing which part of an unmarked `AGENTS.md` is Savepoint-owned means pattern-matching headings, and a wrong guess wraps the user's own prose in markers that the next upgrade then overwrites. Writing `.new` and letting the user place the markers is the only option that cannot destroy their content.

This is the same governing rule as T002 — keep theirs, offer ours — applied to the other shared file.
