---
id: I-096
title: Code Health map still denies implemented collection
type: drift
status: open
source:
  kind: check
  check: C-941
  actor: {role: checker, session: o028-check-20261001}
  at: '2026-10-01T08:38:18Z'
tasks: [T-057, T-058]
checks: [C-941]
guardrail_ids: [ARCH-04, TPL-02]
history:
  - at: '2026-10-01T08:50:00Z'
    actor: {role: executor, session: o028-repair-20261001}
    kind: repair_attempted
    note: "AGENTS.md Code Health map now describes Collect (runs tools, readers, one snapshot) and says production readers and a collection command are not shipped. No template carries the sentence. Gates: make build, make test-fast, make test-full (go1.26.2 linux/amd64) pass; awaiting a checker."
    check: C-941
---

# I-096: Code Health map still denies implemented collection

## Summary

The project's package responsibility map was updated for discovery and setup, but still explicitly denies collection despite T-057 implementing subprocess execution and snapshot persistence. ARCH-04 requires a responsibility-map update; TPL-02 requires guidance to describe current behavior.

## Evidence

C-941 frozen M12: `AGENTS.md:172` ends the Code Health description with `It performs no collection.` Actual public service `internal/codehealth/collect.go:90` runs configured instances and saves a snapshot at line 159. T-057's Drift Notes explicitly required updating the map to describe execution and snapshot saves, but T-058 retained the contradictory sentence.

Expected: the map accurately describes the implemented service and its separation from discovery/setup. Actual: agents following the map receive an incorrect package boundary. This is low materiality guidance drift, not a process-security failure.

## Proof Needed

Reconcile the Code Health map with actual collection responsibilities, checking live/scaffold guidance as applicable. Do not imply production readers or a collection CLI shipped in O-028. Verify scope M12 and template equality; a metadata-only recheck may reuse documented full evidence only under the repository's unchanged-input rule.
