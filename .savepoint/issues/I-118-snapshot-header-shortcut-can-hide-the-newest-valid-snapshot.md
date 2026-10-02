---
id: I-118
title: Snapshot header shortcut can hide the newest valid snapshot
type: defect
status: open
source:
  kind: check
  check: C-957
  actor: {role: checker, session: check-o032-20261002-independent}
  at: '2026-10-02T05:23:00Z'
tasks: [T-098]
checks: [C-957]
guardrail_ids: [CFG-01, TEST-02]
severity: medium
history:
  - at: '2026-10-02T05:23:00Z'
    actor: {role: checker, session: check-o032-20261002-independent}
    kind: observed
    check: C-957
    note: Initial independent Full O-032 Check; frozen scope and reproducible harness in C-957.
  - at: '2026-10-02T05:27:50Z'
    actor: {role: executor, session: repair-o032-20261002}
    kind: repair_attempted
    note: 'parseHead now scans the whole 512-byte head and refuses (full decode instead) a repeated created_at or origin. Added TestLoadWindowAgreesWithTheDecoderOnRepeatedHeaderFields. A repeat beyond the head stays inside the documented old-body boundary.'
---

# I-118: Snapshot header shortcut can hide the newest valid snapshot

## Summary

T-098 Done When 1 requires the window to use the same stored-time/identity order as `LoadSnapshots`. The header shortcut and full decoder disagree on repeated JSON fields, so a valid snapshot can be omitted and an older snapshot reported as newest.

## Evidence

C-957 `TestO032WindowDuplicateHeaderMatchesFullDecoder/created_at`: create 40 valid snapshots, prepend `"created_at":"2000-01-01T00:00:00Z",` immediately after the newest snapshot's opening `{`, preserving its real later `created_at` and every other field. `DecodeSnapshot` and `LoadSnapshots` accept it, retaining its original verified ID (JSON decoder's last value wins). `parseHead` stops once it has the first timestamp and origin (`internal/codehealth/window.go:173-202`), sorts that newest snapshot into old history at :69-75 and never decodes it. `LoadDashboard` reports sha256:2b52ef6a3e6ff6e88c4bbb298ce6dc537b311f3c80c198aa5dcc16d73c0ca9c1 rather than the actual newest sha256:363552ff89fc60ae03d7bc1b93c8c07606a5e662e34aa4647088a19a46a13456; its history order also differs. This is a fully valid body and digest, not the owner-approved damaged-body-outside-window case. Compact, reordered, 512-byte leading whitespace, wrong-type, trailing-delimiter and symlink cells were also completed. Existing window tests cover reordered fields, but not repeated header fields.

## Proof Needed

Make shortcut ordering agree with accepted full-decoder representations or clearly refuse ambiguous headers through the supported loader. Do not silently reinterpret the immutable file's timestamp/origin. Verify repeated timestamp/origin both before and after the normal header fields, order/tie/window boundary, newest identity and no-write behavior. Retain the documented old-body-damage boundary and measured cost.

Repair directly under this Issue by default; preserve all owner-completed Task statuses. This record grants no clearance or owner acceptance.
