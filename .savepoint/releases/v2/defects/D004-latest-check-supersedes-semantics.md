---
id: v2/D004-latest-check-supersedes-semantics
release: v2
status: open
severity: blocker
title: "LatestCheck can disagree with explicit Check supersession"
---

# D004: LatestCheck can disagree with explicit Check supersession

## Symptom

`internal/data/project.go` sorts Checks by numeric `C###` ID and assigns the
highest ID to `V2Index.LatestCheck`. The same load path validates
`supersedes`, but does not use the supersession graph to choose the current
Check or reject multiple valid chain heads. `ResolveClearance()` then feeds
that choice into dependencies, gates, and `Next`.

## Expected Behavior

The explicit `supersedes` relationship must be authoritative, or loading a
project must enforce an invariant that makes numeric-ID order and
supersession order impossible to disagree.

## Reproduction

1. Create multiple Checks for one Task, Objective, or Release with a valid
   `supersedes` chain and a higher-numbered Check that is not part of that
   chain.
2. Load the V2 index and inspect `LatestCheck` for the scope.
3. Resolve clearance or a dependent gate and observe that the numeric-highest
   Check is selected even though the explicit supersession graph does not make
   it the chain's current head.

## Impact

A stale or unrelated Check can determine technical clearance and propagate an
incorrect result through dependency decisions, gates, and the project-wide
`Next` action.

## Fix Plan

Confirmed contract: fail closed at load. A scope's Checks must form one
ID-ordered supersession chain; a multi-head or stray high-ID Check that makes
numeric order and supersession order disagree is rejected with a named load
error. `LatestCheck` keeps its current semantics, so no downstream consumer
changes. Add coverage for multiple chain heads and for the downstream
`ResolveClearance`/gate/`Next` outcome.

## Acceptance Criteria

- [ ] `LatestCheck` and explicit `supersedes` semantics have one authoritative,
      enforced contract.
- [ ] A disagreement cannot silently load as a plausible current Check.
- [ ] Clearance, dependency, gate, and `Next` tests cover the repaired choice.

## Resolution Notes

Pending.
