---
id: v2/D006-clear-evidence-invariant-undecided
release: v2
status: open
severity: high
title: "CLEAR Checks do not define whether reviewed evidence is required"
---

# D006: CLEAR Checks do not define whether reviewed evidence is required

## Symptom

`DecodeCheckV2()` permits `reviewed` to be absent or present with empty
entries, while `ResolveClearance()` relies on independent checker provenance
and current freshness rather than the Check's `reviewed` basis. The existing
`TestDecodeCheckV2_reviewedEmptyEntriesStayAbsent` test preserves this
ambiguity instead of defining whether a CLEAR result represents technical
verification without a substantive reviewed basis.

## Expected Behavior

The project must explicitly choose and enforce one contract: either CLEAR
requires an appropriate reviewed-evidence basis, or empty/absent reviewed
evidence is legitimate and that contract is documented consistently across the
decoder, clearance resolver, and user-facing wording.

## Reproduction

1. Record a `CLEAR` Check with a checker `checked_by` actor and no meaningful
   `reviewed` content, such as `reviewed: {files: [], dependencies: []}`.
2. Provide current freshness assessed by a checker.
3. Load the project and resolve clearance; the Check can count as current
   despite carrying effectively empty reviewed evidence.

## Impact

The system can present technical clearance whose evidence contract is unclear,
weakening the trust boundary between a recorded CLEAR result and the work
that was actually verified.

## Fix Plan

Decide the CLEAR evidence invariant, encode it in the Check decoder and
clearance semantics, update the documentation and wording, and retain or
replace the empty-evidence regression test to lock the decision.

## Acceptance Criteria

- [ ] The CLEAR evidence contract is documented in the V2 record and workflow
      guidance.
- [ ] Decoder and clearance behavior enforce the chosen contract.
- [ ] Tests cover absent, empty, and substantive reviewed evidence.

## Resolution Notes

Pending.
