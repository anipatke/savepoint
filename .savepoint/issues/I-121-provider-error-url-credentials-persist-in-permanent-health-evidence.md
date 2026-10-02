---
id: I-121
title: Provider error URL credentials persist in permanent health evidence
type: defect
status: open
source:
  kind: check
  check: C-957
  actor: {role: checker, session: check-o032-20261002-independent}
  at: '2026-10-02T05:23:00Z'
tasks: [T-085, T-091]
checks: [C-957, C-958]
guardrail_ids: [TEST-02, TEST-05]
severity: high
history:
  - at: '2026-10-02T05:23:00Z'
    actor: {role: checker, session: check-o032-20261002-independent}
    kind: observed
    check: C-957
    note: Initial independent Full O-032 Check; frozen scope and reproducible harness in C-957.
  - at: '2026-10-02T05:27:50Z'
    actor: {role: executor, session: repair-o032-20261002}
    kind: repair_attempted
    note: 'cleanLine drops URL userinfo before bounding; Collect also sanitizes runner stderr. Added TestFailureReasonsDropURLCredentials (ExecRunner-style sanitizers and Collect). A credential cut off from its scheme by the stderr tail limit is not detected. T-085 follow-up line removed from CHANGELOG.'
  - at: '2026-10-02T06:19:18Z'
    actor: {role: checker, session: recheck-o032-20261002-independent}
    kind: rechecked
    check: C-958
    note: "Normal URI credentials redact, but actual ExecRunner plus Collect at 65548 stderr bytes truncates the scheme before sanitization and persists synthetic URI credential in result reason and immutable snapshot. Below/exact 64KiB cases pass. Same oversized-stderr proof cell; remains material and open. Harness and output in C-958."
  - at: '2026-10-02T06:21:27Z'
    actor: {role: executor, session: repair-o032-20261002b}
    kind: repair_attempted
    note: "Targeted repair for C-958: when ExecRunner's stderr tail dropped earlier bytes, tailBuffer.Text removes a leading user:password@ fragment (scheme cut away) before sanitizing. Added TestExecRunnerDropsCredentialsWhoseSchemeWasTruncated, sweeping cut positions. make test-full passed locally."
---

# I-121: Provider error URL credentials persist in permanent health evidence

## Summary

T-091 Done When 2 requires safe diagnostics without source/secrets leakage. A credential-bearing URL in a provider error survives sanitization and is saved into immutable health evidence. Naming redaction as a future follow-up does not supply an explicit owner exception to this accepted criterion.

## Evidence

C-957 `TestO032ProviderErrorCredentialPersistence` uses only a synthetic credential. A configured local Lizard instance receives exit 127 and stderr `proxy connection failed https://alice:o032-secret-for-probe@proxy.example.invalid`. The string passes through the production `sanitizeTail`, then supported `Collect` with a deterministic fake ToolRunner. Both returned result.Reason and the on-disk snapshot loaded by `LoadSnapshots` contain the credential. `internal/codehealth/runner.go:143-182` removes controls and bounds text, but does not redact URL userinfo; `collect.go:339-345` appends stderr to the persisted reason. The post-save dashboard/report also uses that reason. Existing `TestCollectFailureReasonsLeakNeitherReportNorEnvironment` sets a secret environment variable but uses stderr "boom"; it cannot prove safety when the provider itself includes a credential. T-085 and CHANGELOG name URL redaction as unresolved follow-up, with no recorded requirement-specific owner exception. This finding contains no actual secret and no production snapshot was modified by the probe.

## Proof Needed

Use the smallest bounded protection for credentials in provider error URLs before persistence/display, retaining useful final error text, UTF-8 validity, limits and control sanitization. Verify normal and oversized stderr, URI userinfo, and failure paths through ExecRunner and Collect. If the owner wants to accept this exposure, record the exact permitted requirement exception explicitly; do not infer acceptance from a follow-up list.

Repair directly under this Issue by default; preserve all owner-completed Task statuses. This record grants no clearance or owner acceptance.
