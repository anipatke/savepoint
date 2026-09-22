---
id: I005
title: Check guidance requires explicit reviewed evidence while runtime treats it as optional
type: defect
status: resolved
source:
    kind: migration
    actor:
        role: owner
        session: .savepoint/releases/v2/defects/D006-clear-evidence-invariant-undecided.md
    at: "2026-09-20T05:03:29Z"
severity: high
resolution:
    disposition: accepted
    actor:
        role: owner
        session: owner-request-2026-09-22
    at: "2026-09-22T08:46:43Z"
    reason: Owner accepted the implemented contract without an independent proof Check.
---
## Migrated from V1

Relocated verbatim from the V1 defect body at `.savepoint/releases/v2/defects/D006-clear-evidence-invariant-undecided.md` (release `v2`).

## V1 Body (verbatim)

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

Confirmed contract: empty or absent `reviewed` evidence is legitimate for a
CLEAR Check. CLEAR is established by independent checker provenance and current
freshness, not by reviewed content. Document the contract across the Check
decoder, `ResolveClearance`, the V2 record and workflow guidance, and the
user-facing wording; keep the empty-evidence regression test and extend it to
cover absent, empty, and substantive reviewed evidence.

## Acceptance Criteria

- [ ] The CLEAR evidence contract is documented in the V2 record and workflow
      guidance.
- [ ] Decoder and clearance behavior enforce the chosen contract.
- [ ] Tests cover absent, empty, and substantive reviewed evidence.

## Resolution Notes

Implementation records the confirmed contract: `reviewed` is optional scope
metadata, including for `CLEAR`; technical clearance comes from independent
checker provenance and current freshness. The decoder, resolver comments,
workflow guidance, scaffold guidance, and regression tests now state and cover
that contract. The owner accepted closure without an independent proof Check.
