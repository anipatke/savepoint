---
id: I-118
title: Snapshot header shortcut can hide the newest valid snapshot
type: defect
status: resolved
source:
  kind: check
  check: C-957
  actor: {role: checker, session: check-o032-20261002-independent}
  at: '2026-10-02T05:23:00Z'
tasks: [T-098]
checks: [C-957, C-958, C-959, C-960]
guardrail_ids: [CFG-01, TEST-02]
severity: medium
resolution:
  disposition: verified
  check: C-960
  actor: {role: checker, session: recheck-o032-20261003-independent}
  at: '2026-10-02T19:20:00Z'
  reason: "Against the owner-amended T-098: valid repeated created_at/origin agree with LoadSnapshots on newest, order and history; amended DW4 damage boundary passes 17/17 with no writes; single n=1,000 38 ms under the 250 ms budget, heavy cost measured as recorded."
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
  - at: '2026-10-02T06:19:18Z'
    actor: {role: checker, session: recheck-o032-20261002-independent}
    kind: rechecked
    check: C-958
    note: "Original repeated timestamp/origin and representation reproductions pass; duplicate detection falls back to full decode. Technically proven for frozen finding, open pending CLEAR Check proof."
  - at: '2026-10-02T06:33:20Z'
    actor: {role: executor, session: repair-o032-20261002b}
    kind: repair_attempted
    note: "Owner-directed simplification: the 512-byte head shortcut (parseHead) is removed. LoadWindow reads each file and takes created_at/origin with json.Unmarshal, the decoder's own last-wins rule, validating only the window. Changes the T-097 bounded-read design: older bodies are now read, not decoded; dashboard 33 ms at 1,000 normal snapshots (was 13) and 0.56 s at 1,000 maximum-size (was 34 ms; full load 1.2 s). User text now says Older history was not checked."
  - at: '2026-10-02T06:50:00Z'
    actor: {role: checker, session: recheck-o032-20261002c-independent}
    kind: rechecked
    check: C-959
    note: "Original finding proven: valid repeated created_at/origin (before, and same-value after) agree with LoadSnapshots on newest, order and history; representations and the I-117 matrix pass. Not met: a truncated or trailing-garbage file outside the window now fails the dashboard load (T-098 DW4), and load cost grows with snapshot size again (heavy n=1,000 561.6 ms / 168.0 MB vs 34.2 ms / 23.7 MB), contradicting T-098 Outcome/DW1/DW6, the T-097 decision and Design.md. No recorded owner exception or amended criterion. Convergence limit reached; owner decides: amend through savepoint-design, record an exception, or restore a bounded read."
  - at: '2026-10-02T07:05:00Z'
    actor: {role: owner, session: plan-o032-20261002-i118-amend}
    kind: owner_decision
    note: "Owner chose option 1 from C-959: keep the whole-file read. Planner amended T-098 (Outcome, Done When 1/3/4/6, User Check), the T-097 decision and Design.md to describe it, including the measured cost and that a non-JSON file outside the window fails the load. No code change. Needs a fresh independent Check for verified."
  - at: '2026-10-02T19:20:00Z'
    actor: {role: checker, session: recheck-o032-20261003-independent}
    kind: rechecked
    check: C-960
    note: "CLEAR. Against the owner-amended T-098: valid repeated created_at/origin agree with LoadSnapshots on newest, order and history; amended DW4 damage boundary passes 17/17 with no writes; single n=1,000 38 ms under the 250 ms budget, heavy cost measured as recorded."
---

# I-118: Snapshot header shortcut can hide the newest valid snapshot

## Summary

T-098 Done When 1 requires the window to use the same stored-time/identity order as `LoadSnapshots`. The header shortcut and full decoder disagree on repeated JSON fields, so a valid snapshot can be omitted and an older snapshot reported as newest.

## Evidence

C-957 `TestO032WindowDuplicateHeaderMatchesFullDecoder/created_at`: create 40 valid snapshots, prepend `"created_at":"2000-01-01T00:00:00Z",` immediately after the newest snapshot's opening `{`, preserving its real later `created_at` and every other field. `DecodeSnapshot` and `LoadSnapshots` accept it, retaining its original verified ID (JSON decoder's last value wins). `parseHead` stops once it has the first timestamp and origin (`internal/codehealth/window.go:173-202`), sorts that newest snapshot into old history at :69-75 and never decodes it. `LoadDashboard` reports sha256:2b52ef6a3e6ff6e88c4bbb298ce6dc537b311f3c80c198aa5dcc16d73c0ca9c1 rather than the actual newest sha256:363552ff89fc60ae03d7bc1b93c8c07606a5e662e34aa4647088a19a46a13456; its history order also differs. This is a fully valid body and digest, not the owner-approved damaged-body-outside-window case. Compact, reordered, 512-byte leading whitespace, wrong-type, trailing-delimiter and symlink cells were also completed. Existing window tests cover reordered fields, but not repeated header fields.

## Proof Needed

Make shortcut ordering agree with accepted full-decoder representations or clearly refuse ambiguous headers through the supported loader. Do not silently reinterpret the immutable file's timestamp/origin. Verify repeated timestamp/origin both before and after the normal header fields, order/tie/window boundary, newest identity and no-write behavior. Retain the documented old-body-damage boundary and measured cost.

Repair directly under this Issue by default; preserve all owner-completed Task statuses. This record grants no clearance or owner acceptance.
