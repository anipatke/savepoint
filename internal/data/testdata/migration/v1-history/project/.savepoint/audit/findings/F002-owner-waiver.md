---
id: F002
title: Health-Check custom procedure omits Deep mode
status: waived
severity: low
confidence: medium
source_auditor: agent
work_item: E01-example/T001-shared
releases:
  - v1.1
guardrail_ids:
  - POL-02
locations:
  - .savepoint/Health-Check.md
first_seen: 2026-06-22
last_seen: 2026-07-01
proof_needed: "n/a; owner accepted the gap"
waiver_reason: Owner decided Deep mode is unnecessary for this example project; only Quick is ever exercised here, and this frozen fixture will never reach production launch.
---

## Summary

The project's custom `Health-Check.md` only defines a Quick procedure; it never adopted
Full or Deep modes.

## Evidence

`.savepoint/Health-Check.md` in this fixture defines only `## Quick Check`; there is no
`## Full Check` or `## Deep Check` section.

## Proof

No proof is required: this is an owner disposition, not a technical fix. The `##
Summary` gap remains exactly as described; only the decision to accept it changed.

## History

- 2026-06-22: Opened while reviewing the custom Health-Check procedure.
- 2026-07-01: Owner waived; the waiver reason is the only recorded disposition, and no
  proof section is required or claimed for a `waived` finding.
